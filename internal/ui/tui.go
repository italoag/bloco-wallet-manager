package ui

import (
	"blocowallet/internal/blockchain"
	"blocowallet/internal/wallet"
	"blocowallet/pkg/config"
	"context"
	"fmt"
	"log"
	"math"
	"math/big"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/digitallyserviced/tdfgo/tdf"
	zone "github.com/lrstanley/bubblezone"
)

// NewModel creates a new TUI model
func NewModel(walletService *wallet.Service, cfg *config.Config) Model {
	// Initialize text inputs
	nameInput := textinput.New()
	nameInput.Placeholder = "Enter wallet name..."
	nameInput.Width = 40

	passwordInput := textinput.New()
	passwordInput.Placeholder = "Enter password..."
	passwordInput.EchoMode = textinput.EchoPassword
	passwordInput.Width = 40

	mnemonicInput := textinput.New()
	mnemonicInput.Placeholder = "Enter 12-word mnemonic phrase..."
	mnemonicInput.Width = 60

	privateKeyInput := textinput.New()
	privateKeyInput.Placeholder = "Enter private key (without 0x prefix)..."
	privateKeyInput.EchoMode = textinput.EchoPassword
	privateKeyInput.Width = 60

	// Private key password input for viewing keystore keys
	privateKeyPassword := textinput.New()
	privateKeyPassword.Placeholder = "Enter keystore password to view private key..."
	privateKeyPassword.EchoMode = textinput.EchoPassword
	privateKeyPassword.Width = 50

	// Wallet authentication password input for accessing wallet details
	walletAuthPassword := textinput.New()
	walletAuthPassword.Placeholder = "Enter wallet password to access details..."
	walletAuthPassword.EchoMode = textinput.EchoPassword
	walletAuthPassword.Width = 50

	// RPC input
	rpcInput := textinput.New()
	rpcInput.Placeholder = "https://..."
	rpcInput.Width = 60

	// Add network inputs
	networkNameInput := textinput.New()
	networkNameInput.Placeholder = "Enter network name..."
	networkNameInput.Width = 40

	chainIDInput := textinput.New()
	chainIDInput.Placeholder = "Enter chain ID (e.g., 137)..."
	chainIDInput.Width = 40

	rpcEndpointInput := textinput.New()
	rpcEndpointInput.Placeholder = "https://..."
	rpcEndpointInput.Width = 60

	// Preparar items dos menus
	settingsItems := []string{"Network Configuration", "Language", "Back to Main Menu"}

	var networkItems []string
	networkKeys := cfg.GetAllNetworkKeys()
	for _, key := range networkKeys {
		if network, exists := cfg.GetNetworkByKey(key); exists {
			status := ""
			if network.IsActive {
				status = " (Active)"
			}
			networkItems = append(networkItems, fmt.Sprintf("%s%s", network.Name, status))
		}
	}
	networkItems = append(networkItems, "Add Network", "Back to Settings")

	var languageItems []string
	langCodes := cfg.GetLanguageCodes()
	for _, code := range langCodes {
		name := config.SupportedLanguages[code]
		status := ""
		if cfg.Language == code {
			status = " (Current)"
		}
		languageItems = append(languageItems, fmt.Sprintf("%s%s", name, status))
	}
	languageItems = append(languageItems, "Back to Settings")

	// Create the model
	model := Model{
		walletService:       walletService,
		config:              cfg,
		wallets:             []*wallet.Wallet{},
		selected:            0,
		loading:             false,
		currentView:         SplashView,
		nameInput:           nameInput,
		passwordInput:       passwordInput,
		mnemonicInput:       mnemonicInput,
		privateKeyInput:     privateKeyInput,
		privateKeyPassword:  privateKeyPassword,
		extractedPrivateKey: "",
		privateKeyError:     "",
		needsWalletAuth:     false,
		walletAuthPassword:  walletAuthPassword,
		walletAuthError:     "",
		inputFocus:          0,
		rpcInput:            rpcInput,
		networkNameInput:    networkNameInput,
		chainIDInput:        chainIDInput,
		rpcEndpointInput:    rpcEndpointInput,
		addNetworkFocus:     0,
		addingNetwork:       false,
		selectedNetworkKey:  "",
		networkSuggestions:  []blockchain.NetworkSuggestion{},
		showingSuggestions:  false,
		selectedSuggestion:  0,
		chainListService:    blockchain.NewChainListService(),
		settingsSelected:    0,
		networkSelected:     0,
		languageSelected:    0,
		editingRPC:          false,
		settingsItems:       settingsItems,
		networkItems:        networkItems,
		languageItems:       languageItems,
	}

	// Load a default TDF font
	model.loadDefaultFont()

	return model
}

// loadDefaultFont loads a default TDF font for the splash screen
func (m *Model) loadDefaultFont() {
	log.Println("Starting TDF font loading...")

	// Get current working directory
	cwd, err := os.Getwd()
	if err != nil {
		log.Printf("Error getting cwd: %v", err)
		return
	}
	log.Printf("Current working directory: %s", cwd)

	// Try different possible paths for fonts directory
	fontBasePaths := []string{
		"fonts",                // If running from project root
		"../fonts",             // If running from bin directory
		"../../fonts",          // If running from nested directory
		"~/.blocowallet/fonts", // User's home directory
		filepath.Join(os.Getenv("HOME"), ".blocowallet", "fonts"), // User's home directory
	}

	fontNames := []string{
		"dynasty.tdf",
		"carbonx.tdf",
		"eleuthix.tdf",
		"grandx.tdf",
		"portal.tdf",
	}

	for _, basePath := range fontBasePaths {
		for _, fontName := range fontNames {
			fontPath := filepath.Join(basePath, fontName)
			log.Printf("Attempting to load font: %s", fontPath)

			// Check if file exists first
			if _, err := os.Stat(fontPath); os.IsNotExist(err) {
				log.Printf("Font file does not exist: %s", fontPath)
				continue
			}

			// Try absolute path as well
			absPath, _ := filepath.Abs(fontPath)
			log.Printf("Absolute path: %s", absPath)

			fontInfo := &tdf.FontInfo{
				Path:    fontPath,
				File:    filepath.Base(fontPath),
				BuiltIn: false,
			}

			if fontFile, err := tdf.LoadFont(fontInfo); err == nil && len(fontFile.Fonts) > 0 {
				m.selectedFont = &fontFile.Fonts[0] // Use the first font in the file
				m.fontName = filepath.Base(fontPath)
				log.Printf("Successfully loaded TDF font: %s (contains %d fonts)", fontPath, len(fontFile.Fonts))
				return
			} else {
				log.Printf("Failed to load font %s: %v", fontPath, err)
			}
		}
	}

	// If no font loads successfully, selectedFont will remain nil
	// and renderSplash will fall back to text rendering
	log.Println("No TDF fonts could be loaded, using fallback text rendering")
}

// Init initializes the TUI
func (m Model) Init() tea.Cmd {
	// Initialize bubblezone manager for mouse support
	zone.NewGlobal()
	return tea.EnterAltScreen
}

// Update handles user input and state changes
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// Se há um diálogo de exclusão ativo, processar suas mensagens primeiro
	if m.deleteDialog != nil {
		if keyMsg, ok := msg.(tea.KeyMsg); ok {
			updatedDialog, cmd := m.deleteDialog.Update(keyMsg)
			if dialog, ok := updatedDialog.(DeleteWalletDialog); ok {
				m.deleteDialog = &dialog
			}
			return m, cmd
		}
		if mouseMsg, ok := msg.(tea.MouseMsg); ok {
			updatedDialog, cmd := m.deleteDialog.Update(mouseMsg)
			if dialog, ok := updatedDialog.(DeleteWalletDialog); ok {
				m.deleteDialog = &dialog
			}
			return m, cmd
		}
		if winMsg, ok := msg.(tea.WindowSizeMsg); ok {
			updatedDialog, _ := m.deleteDialog.Update(winMsg)
			if dialog, ok := updatedDialog.(DeleteWalletDialog); ok {
				m.deleteDialog = &dialog
			}
		}
	}

	switch msg := msg.(type) {
	case tea.MouseMsg:
		// Update bubblezone with mouse events
		return zone.AnyInBoundsAndUpdate(m, msg)

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case tea.KeyMsg:
		return m.handleKeyMsg(msg)

	case walletsLoadedMsg:
		m.wallets = []*wallet.Wallet(msg)
		m.loading = false
		m.err = nil
		return m, nil

	case balanceLoadedMsg:
		m.currentBalance = (*wallet.Balance)(msg)
		return m, nil

	case multiBalanceLoadedMsg:
		m.currentMultiBalance = (*wallet.MultiNetworkBalance)(msg)
		return m, nil

	case walletCreatedMsg:
		newModel := m
		newModel.currentView = WalletListView
		newModel.selected = 0
		// Reset inputs manually
		newModel.nameInput.SetValue("")
		newModel.passwordInput.SetValue("")
		newModel.mnemonicInput.SetValue("")
		newModel.privateKeyInput.SetValue("")
		newModel.inputFocus = 0
		newModel.nameInput.Blur()
		newModel.passwordInput.Blur()
		newModel.mnemonicInput.Blur()
		newModel.privateKeyInput.Blur()
		return newModel, newModel.loadWalletsCmd()

	case errorMsg:
		m.err = fmt.Errorf("%s", string(msg))
		m.loading = false
		return m, nil

	case networkSuggestionsMsg:
		m.networkSuggestions = []blockchain.NetworkSuggestion(msg)
		m.showingSuggestions = len(m.networkSuggestions) > 0
		m.selectedSuggestion = 0
		return m, nil

	case chainInfoLoadedMsg:
		// Auto-fill fields when chain info is loaded by Chain ID
		chainInfo := msg.chainInfo
		if chainInfo != nil {
			// Fill network name if empty
			if strings.TrimSpace(m.networkNameInput.Value()) == "" {
				m.networkNameInput.SetValue(chainInfo.Name)
			}
			// Fill RPC endpoint if empty
			if strings.TrimSpace(m.rpcEndpointInput.Value()) == "" {
				m.rpcEndpointInput.SetValue(msg.rpcURL)
			}
		}
		return m, nil

	case ConfirmDeleteMsg:
		// Confirmar a exclusão da wallet
		if m.deleteDialog != nil && m.selectedWallet != nil {
			if err := m.walletService.DeleteWalletByAddress(context.Background(), m.selectedWallet.Address); err != nil {
				m.err = fmt.Errorf("failed to delete wallet: %w", err)
			} else {
				// Reset state and go back to wallet list
				m.currentView = WalletListView
				m.selectedWallet = nil
				m.deleteDialog = nil
				m.showSensitiveInfo = false
				m.extractedPrivateKey = ""
				// Reload wallets
				return m, tea.Cmd(func() tea.Msg {
					wallets, err := m.walletService.GetAllWallets(context.Background())
					if err != nil {
						return errorMsg(err.Error())
					}
					return walletsLoadedMsg(wallets)
				})
			}
		}
		return m, nil

	case CancelDeleteMsg:
		// Cancelar a exclusão da wallet
		m.deleteDialog = nil
		return m, nil

	case networkAddedMsg:
		// Add the custom network to config
		m.config.AddCustomNetwork(msg.key, msg.network)
		if err := m.config.Save(); err != nil {
			m.err = err
		} else {
			// Refresh network items and go back to network config
			m.refreshNetworkItems()
			m.currentView = NetworkConfigView
			m.networkSelected = 0
			m.addingNetwork = false
			// Reset inputs
			m.networkNameInput.SetValue("")
			m.chainIDInput.SetValue("")
			m.rpcEndpointInput.SetValue("")
			m.addNetworkFocus = 0
		}
		return m, nil

	case networkErrorMsg:
		m.err = fmt.Errorf("%s", string(msg))
		m.addingNetwork = false
		return m, nil
	}

	return m, nil
}

// handleKeyMsg handles keyboard input
func (m Model) handleKeyMsg(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	// Handle input fields in create/import views, RPC editing, wallet auth, and add network
	if m.currentView == CreateWalletView || m.currentView == ImportWalletView || m.currentView == ImportPrivateKeyView ||
		m.currentView == WalletAuthView ||
		(m.currentView == NetworkConfigView && m.editingRPC) ||
		(m.currentView == AddNetworkView && !m.addingNetwork) {
		switch msg.String() {
		case "esc":
			newModel := m
			if newModel.currentView == NetworkConfigView && newModel.editingRPC {
				// Cancel RPC editing
				newModel.editingRPC = false
				newModel.rpcInput.Blur()
				return newModel, nil
			} else if newModel.currentView == WalletAuthView {
				// Cancel wallet authentication and go back to wallet list
				newModel.currentView = WalletListView
				newModel.needsWalletAuth = false
				newModel.walletAuthPassword.SetValue("")
				newModel.walletAuthPassword.Blur()
				newModel.walletAuthError = ""
				newModel.selectedWallet = nil
				return newModel, nil
			}
			newModel.currentView = MenuView
			newModel.selected = 0
			// Reset inputs manually instead of calling resetInputs
			newModel.nameInput.SetValue("")
			newModel.passwordInput.SetValue("")
			newModel.mnemonicInput.SetValue("")
			newModel.privateKeyInput.SetValue("")
			newModel.inputFocus = 0
			newModel.nameInput.Blur()
			newModel.passwordInput.Blur()
			newModel.mnemonicInput.Blur()
			newModel.privateKeyInput.Blur()
			return newModel, nil
		case "tab", "shift+tab", "up", "down":
			if m.currentView == AddNetworkView {
				return m.handleAddNetworkNavigation(msg)
			}
			return m.handleInputNavigation(msg)
		case "enter":
			if m.currentView == NetworkConfigView && m.editingRPC {
				// Save RPC endpoint
				networkKeys := m.config.GetAllNetworkKeys()
				if m.networkSelected < len(networkKeys) {
					key := networkKeys[m.networkSelected]
					if network, exists := m.config.GetNetworkByKey(key); exists {
						network.RPCEndpoint = strings.TrimSpace(m.rpcInput.Value())
						if network.IsCustom {
							m.config.CustomNetworks[key] = network
						} else {
							m.config.Networks[key] = network
						}
						// Save configuration to file
						_ = m.config.Save()
					}
				}
				m.editingRPC = false
				m.rpcInput.Blur()
				return m, nil
			} else if m.currentView == AddNetworkView {
				// Check if we're showing suggestions first
				if m.showingSuggestions && m.addNetworkFocus == 0 {
					// Handle suggestion selection
					return m.handleAddNetworkNavigation(msg)
				}
				// Otherwise handle form submission
				return m.handleAddNetworkSubmit()
			} else if m.currentView == WalletAuthView {
				// Handle wallet authentication
				return m.handleInputSubmit()
			}
			return m.handleInputSubmit()
		default:
			// Update the focused input field with the key
			newModel := m
			var cmd tea.Cmd

			if newModel.currentView == NetworkConfigView && newModel.editingRPC {
				// Handle RPC input
				newModel.rpcInput, cmd = newModel.rpcInput.Update(msg)
			} else if newModel.currentView == AddNetworkView {
				// Handle add network inputs
				switch newModel.addNetworkFocus {
				case 0:
					oldValue := newModel.networkNameInput.Value()
					newModel.networkNameInput, cmd = newModel.networkNameInput.Update(msg)
					newValue := newModel.networkNameInput.Value()

					// Trigger search if value changed and has at least 2 characters
					if oldValue != newValue && len(strings.TrimSpace(newValue)) >= 2 {
						return newModel, tea.Batch(cmd, newModel.searchNetworksCmd(newValue))
					}

					// Hide suggestions if input is too short
					if len(strings.TrimSpace(newValue)) < 2 {
						newModel.showingSuggestions = false
						newModel.networkSuggestions = []blockchain.NetworkSuggestion{}
					}
				case 1:
					oldValue := newModel.chainIDInput.Value()
					newModel.chainIDInput, cmd = newModel.chainIDInput.Update(msg)
					newValue := newModel.chainIDInput.Value()

					// Auto-load chain info when Chain ID is entered
					if oldValue != newValue && strings.TrimSpace(newValue) != "" {
						if chainID, err := strconv.Atoi(strings.TrimSpace(newValue)); err == nil {
							return newModel, tea.Batch(cmd, newModel.loadChainInfoByIDCmd(chainID))
						}
					}
				case 2:
					newModel.rpcEndpointInput, cmd = newModel.rpcEndpointInput.Update(msg)
				}
			} else if newModel.currentView == WalletAuthView {
				// Handle wallet authentication password input
				newModel.walletAuthPassword, cmd = newModel.walletAuthPassword.Update(msg)
			} else {
				switch newModel.inputFocus {
				case 0: // Name input
					newModel.nameInput, cmd = newModel.nameInput.Update(msg)
				case 1: // Password input
					newModel.passwordInput, cmd = newModel.passwordInput.Update(msg)
				case 2: // Mnemonic input (only in import view) or Private Key input (only in import private key view)
					switch newModel.currentView {
					case ImportWalletView:
						newModel.mnemonicInput, cmd = newModel.mnemonicInput.Update(msg)
					case ImportPrivateKeyView:
						newModel.privateKeyInput, cmd = newModel.privateKeyInput.Update(msg)
					}
				}
			}

			return newModel, cmd
		}
	}

	// Regular navigation for other views
	switch msg.String() {
	case "ctrl+c", "q":
		return m, tea.Quit

	case "esc":
		switch m.currentView {
		case WalletAuthView:
			m.currentView = WalletListView
			m.needsWalletAuth = false
			m.walletAuthPassword.Reset()
			m.walletAuthError = ""
		case WalletDetailsView:
			m.currentView = WalletListView
		case CreateWalletView, ImportWalletView, ImportPrivateKeyView:
			m.currentView = MenuView
			m.selected = 0
			m.resetInputs()
		case SettingsView:
			m.currentView = MenuView
			m.selected = 0
		case NetworkConfigView, LanguageView:
			m.currentView = SettingsView
			m.settingsSelected = 0
		case AddNetworkView:
			m.currentView = NetworkConfigView
			m.networkSelected = 0
			m.addingNetwork = false
			// Reset inputs
			m.networkNameInput.SetValue("")
			m.chainIDInput.SetValue("")
			m.rpcEndpointInput.SetValue("")
			m.addNetworkFocus = 0
		case NetworkDetailsView:
			m.currentView = NetworkConfigView
			m.networkSelected = 0
		case MenuView:
			m.currentView = SplashView
		case WalletListView:
			m.currentView = MenuView
			m.selected = 0
		}
		return m, nil

	case " ", "enter":
		return m.handleEnterKey()

	case "up", "k":
		switch m.currentView {
		case MenuView:
			if m.selected > 0 {
				m.selected--
			}
		case WalletListView:
			if m.selected > 0 {
				m.selected--
			}
		case SettingsView:
			if m.settingsSelected > 0 {
				m.settingsSelected--
			}
		case NetworkConfigView:
			if !m.editingRPC && m.networkSelected > 0 {
				m.networkSelected--
			}
		case LanguageView:
			if m.languageSelected > 0 {
				m.languageSelected--
			}
		case AddNetworkView:
			if !m.addingNetwork && m.addNetworkFocus > 0 {
				m.addNetworkFocus--
				m.updateAddNetworkFocus()
			}
		}
		return m, nil

	case "down", "j":
		maxItems := 0
		switch m.currentView {
		case MenuView:
			maxItems = 5 // 6 menu items (0-5)
		case WalletListView:
			maxItems = len(m.wallets) - 1
		case SettingsView:
			maxItems = len(m.settingsItems) - 1
		case NetworkConfigView:
			if !m.editingRPC {
				maxItems = len(m.networkItems) - 1
			}
		case LanguageView:
			maxItems = len(m.languageItems) - 1
		case AddNetworkView:
			maxItems = 2 // 3 input fields (0-2)
		}

		if m.currentView == SettingsView && m.settingsSelected < maxItems {
			m.settingsSelected++
		} else if m.currentView == NetworkConfigView && !m.editingRPC && m.networkSelected < maxItems {
			m.networkSelected++
		} else if m.currentView == LanguageView && m.languageSelected < maxItems {
			m.languageSelected++
		} else if m.currentView == AddNetworkView && !m.addingNetwork && m.addNetworkFocus < maxItems {
			m.addNetworkFocus++
			m.updateAddNetworkFocus()
		} else if m.selected < maxItems {
			m.selected++
		}
		return m, nil

	case "a", "A":
		// Toggle network active/inactive status (allows multiple active networks)
		if m.currentView == NetworkConfigView && !m.editingRPC {
			networkKeys := m.config.GetAllNetworkKeys()
			if m.networkSelected < len(networkKeys) {
				key := networkKeys[m.networkSelected]

				// Toggle network active status
				err := m.config.ToggleNetworkActive(key)
				if err == nil {
					// Refresh network items to show updated status
					m.refreshNetworkItems()

					// Refresh multi-provider with updated network configuration
					m.walletService.RefreshMultiProvider(m.config)

					// Save configuration
					_ = m.config.Save()
				}
			}
		}
		return m, nil

	case "e", "E":
		// Edit RPC endpoint
		if m.currentView == NetworkConfigView && !m.editingRPC {
			networkKeys := m.config.GetAllNetworkKeys()
			if m.networkSelected < len(networkKeys) {
				key := networkKeys[m.networkSelected]
				if network, exists := m.config.GetNetworkByKey(key); exists {
					m.editingRPC = true
					m.rpcInput.Focus()
					m.rpcInput.SetValue(network.RPCEndpoint)
				}
			}
		}
		return m, nil

	case "s":
		if m.currentView == WalletDetailsView && m.selectedWallet != nil {
			// Toggle sensitive information visibility
			newModel := m
			newModel.showSensitiveInfo = !newModel.showSensitiveInfo
			return newModel, nil
		}
		return m, nil

	case "d":
		if m.currentView == WalletDetailsView && m.selectedWallet != nil {
			// Criar o diálogo de exclusão
			if m.deleteDialog == nil {
				dialog := NewDeleteWalletDialog(m.selectedWallet.Name, m.selectedWallet.Address)
				m.deleteDialog = &dialog
			}
			return m, nil
		}
		return m, nil

	case "y":
		// A lógica de confirmação de exclusão agora é feita pelo diálogo
		return m, nil

	case "n":
		// A lógica de cancelamento agora é feita pelo diálogo
		return m, nil
	}

	return m, nil
}

// handleInputNavigation handles navigation between input fields
// View renders the TUI
func (m Model) View() string {
	if m.width == 0 {
		return "Loading..."
	}

	var baseView string
	switch m.currentView {
	case SplashView:
		baseView = m.renderSplash()
	case MenuView:
		baseView = m.renderMenu()
	case WalletListView:
		baseView = m.renderWalletList()
	case WalletAuthView:
		baseView = m.renderWalletAuth()
	case WalletDetailsView:
		baseView = m.renderWalletDetails()
	case CreateWalletView:
		baseView = m.renderCreateWallet()
	case ImportWalletView:
		baseView = m.renderImportWallet()
	case ImportPrivateKeyView:
		baseView = m.renderImportPrivateKey()
	case SettingsView:
		baseView = m.renderSettings()
	case NetworkConfigView:
		baseView = m.renderNetworkConfig()
	case LanguageView:
		baseView = m.renderLanguage()
	case AddNetworkView:
		baseView = m.renderAddNetwork()
	case NetworkDetailsView:
		baseView = m.renderNetworkDetails()
	default:
		baseView = m.renderWalletList()
	}

	// Se há um diálogo de exclusão ativo, renderizá-lo sobre a view base
	if m.deleteDialog != nil {
		dialog := m.deleteDialog.View()

		// Centralizar o diálogo
		centeredDialog := lipgloss.Place(
			m.width, m.height,
			lipgloss.Center, lipgloss.Center,
			dialog,
		)

		return centeredDialog
	}

	return baseView
}

// Render methods
func (m Model) renderSplash() string {
	// Check if TDF font is available
	if m.selectedFont != nil {
		// Initialize string renderer for the selected font
		fontString := tdf.NewTheDrawFontStringFont(m.selectedFont)

		// Render the "bloco" logo using TDF font
		renderedLogo := fontString.RenderString("bloco")
		renderedLogo = strings.TrimSpace(renderedLogo)

		// Project info
		projectInfo := fmt.Sprintf("%s v%s", "BLOCO Wallet Manager", "1.0.0")

		// Center the project info text
		projectInfoStyled := lipgloss.NewStyle().
			Align(lipgloss.Center).
			Foreground(lipgloss.Color("86")).
			Render(projectInfo)

		// Create the splash screen content
		splashContent := lipgloss.JoinVertical(
			lipgloss.Center,
			renderedLogo,
			"",
			projectInfoStyled,
			"",
			lipgloss.NewStyle().
				Foreground(lipgloss.Color("241")).
				Align(lipgloss.Center).
				Render("Press any key to continue..."),
		)

		// Use lipgloss.Place to center horizontally and vertically
		finalSplash := lipgloss.Place(
			m.width,
			m.height,
			lipgloss.Center,
			lipgloss.Center,
			splashContent,
		)

		return finalSplash
	}

	// Fallback to original text-based splash if no TDF font is available
	var b strings.Builder

	logoStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("86")).
		Align(lipgloss.Center).
		Width(m.width).
		MarginTop(m.height / 4)

	b.WriteString(logoStyle.Render("🏦 BlocoWallet"))
	b.WriteString("\n\n")

	subtitleStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("241")).
		Align(lipgloss.Center).
		Width(m.width)

	b.WriteString(subtitleStyle.Render("Your Ethereum Wallet Manager"))
	b.WriteString("\n\n")

	instructionStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("86")).
		Align(lipgloss.Center).
		Width(m.width)

	b.WriteString(instructionStyle.Render("Press any key to continue..."))

	return b.String()
}

func (m Model) renderMenu() string {
	var b strings.Builder

	headerStyle := HeaderStyle

	b.WriteString(headerStyle.Render("🏦 BlocoWallet - Main Menu"))
	b.WriteString("\n\n")

	menuItems := []string{
		"📋 View Wallets",
		"➕ Create New Wallet",
		"📥 Import Wallet (Mnemonic)",
		"🔑 Import Wallet (Private Key)",
		"⚙️  Settings",
		"🚪 Exit",
	}

	for i, item := range menuItems {
		if i == m.selected {
			selectedStyle := lipgloss.NewStyle().
				Padding(0, 1)
			b.WriteString(selectedStyle.Render("→ " + item))
		} else {
			b.WriteString("  " + item)
		}
		b.WriteString("\n")
	}

	b.WriteString("\n")
	footerStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("241"))
	b.WriteString(footerStyle.Render("↑/↓: navigate • enter: select • q: quit"))

	return b.String()
}

func (m Model) renderWalletList() string {
	var b strings.Builder

	headerStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("86")).
		MarginBottom(2)

	b.WriteString(headerStyle.Render("📋 Your Wallets"))
	b.WriteString("\n\n")

	if m.loading {
		b.WriteString("Loading wallets...")
		return b.String()
	}

	if len(m.wallets) == 0 {
		b.WriteString("No wallets found. Create a new wallet to get started.\n\n")
	} else {
		for i, wallet := range m.wallets {
			if i == m.selected {
				selectedStyle := lipgloss.NewStyle().
					Padding(0, 1)
				b.WriteString(selectedStyle.Render(fmt.Sprintf("→ %s (%s)", wallet.Name, wallet.Address[:10]+"...")))
			} else {
				b.WriteString(fmt.Sprintf("  %s (%s)", wallet.Name, wallet.Address[:10]+"..."))
			}
			b.WriteString("\n")
		}
		b.WriteString("\n")
	}

	if m.err != nil {
		errorStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("196"))
		b.WriteString(errorStyle.Render("Error: " + m.err.Error()))
		b.WriteString("\n\n")
	}

	footerStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("241"))
	b.WriteString(footerStyle.Render("↑/↓: navigate • enter: view details • r: refresh • esc: back • q: quit"))

	return b.String()
}

func (m Model) renderWalletAuth() string {
	var b strings.Builder

	headerStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("86")).
		MarginBottom(2)

	b.WriteString(headerStyle.Render("🔐 Wallet Authentication"))
	b.WriteString("\n\n")

	if m.selectedWallet != nil {
		infoStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("240"))

		b.WriteString(infoStyle.Render(fmt.Sprintf("Wallet: %s", m.selectedWallet.Name)))
		b.WriteString("\n")
		b.WriteString(infoStyle.Render(fmt.Sprintf("Address: %s", m.selectedWallet.Address)))
		b.WriteString("\n\n")
	}

	b.WriteString("Enter wallet password to access details:\n\n")

	// Password input
	b.WriteString("Password:\n")
	b.WriteString(m.walletAuthPassword.View())
	b.WriteString("\n\n")

	// Show error if any
	if m.walletAuthError != "" {
		errorStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("196"))
		b.WriteString(errorStyle.Render("❌ " + m.walletAuthError))
		b.WriteString("\n\n")
	}

	footerStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("241"))
	b.WriteString(footerStyle.Render("enter: authenticate • esc: back • q: quit"))

	return b.String()
}

func (m *Model) renderWalletDetails() string {
	var b strings.Builder

	headerStyle := lipgloss.NewStyle().
		Bold(true).
		MarginBottom(2)

	if m.selectedWallet == nil {
		b.WriteString(headerStyle.Render("💳 Wallet Details"))
		b.WriteString("\n\nNo wallet selected.")
		return b.String()
	}

	b.WriteString(headerStyle.Render("💳 " + m.selectedWallet.Name))
	b.WriteString("\n\n")

	info := fmt.Sprintf("Address: %s\nCreated: %s",
		m.selectedWallet.Address,
		m.selectedWallet.CreatedAt.Format("2006-01-02 15:04:05"))

	b.WriteString(info)
	b.WriteString("\n\n")

	// Add sensitive information section
	sensitiveStyle := lipgloss.NewStyle().
		Bold(true)

	b.WriteString(sensitiveStyle.Render("🔐 Sensitive Information:"))
	b.WriteString("\n")

	// Private Key
	if m.showSensitiveInfo {
		if m.selectedWallet.EncryptedMnemonic != "" {
			// Get cached password for this wallet
			password, hasPassword := m.walletService.GetWalletPassword(m.selectedWallet.Address)
			if hasPassword {
				// Decrypt mnemonic and derive private key
				mnemonic, err := m.walletService.GetMnemonicFromWallet(m.selectedWallet, password)
				if err != nil {
					b.WriteString(fmt.Sprintf("Private Key: Error decrypting mnemonic - %v\n", err))
				} else {
					privateKey, err := wallet.DerivePrivateKey(mnemonic)
					if err != nil {
						b.WriteString(fmt.Sprintf("Private Key: Error deriving key - %v\n", err))
					} else {
						b.WriteString(fmt.Sprintf("Private Key: %s\n", privateKey))
					}
				}
			} else {
				b.WriteString("Private Key: Authentication required\n")
			}
		} else if m.selectedWallet.KeyStorePath != "" {
			// Handle keystore-based wallets (both mnemonic-based and private key imports)
			if m.extractedPrivateKey != "" {
				// Private key has been successfully extracted
				b.WriteString(fmt.Sprintf("Private Key: %s\n", m.extractedPrivateKey))
			} else {
				// Try to extract using cached password
				if cachedPassword, hasPassword := m.walletService.GetWalletPassword(m.selectedWallet.Address); hasPassword {
					privateKey, err := m.walletService.LoadPrivateKeyFromKeyStoreV3(m.selectedWallet.KeyStorePath, cachedPassword)
					if err != nil {
						b.WriteString(fmt.Sprintf("Private Key: Error - %s\n", err.Error()))
					} else {
						// Cache the extracted private key for subsequent renders
						privateKeyHex := fmt.Sprintf("%x", privateKey.D.Bytes())
						m.extractedPrivateKey = privateKeyHex
						b.WriteString(fmt.Sprintf("Private Key: %s\n", privateKeyHex))
					}
				} else {
					// This shouldn't happen since we require auth before entering details
					b.WriteString("Private Key: Authentication required\n")
				}
			}
		} else {
			// Neither mnemonic nor keystore path available - this should be rare
			b.WriteString("Private Key: Not available (no keystore or mnemonic found)\n")
		}
	} else {
		b.WriteString("Private Key: ********************************\n")
	}

	// Mnemonic
	if m.showSensitiveInfo {
		if m.selectedWallet.EncryptedMnemonic != "" {
			// Get cached password for this wallet
			password, hasPassword := m.walletService.GetWalletPassword(m.selectedWallet.Address)
			if hasPassword {
				// Decrypt mnemonic
				mnemonic, err := m.walletService.GetMnemonicFromWallet(m.selectedWallet, password)
				if err != nil {
					b.WriteString(fmt.Sprintf("Mnemonic: Error decrypting - %v\n", err))
				} else {
					b.WriteString(fmt.Sprintf("Mnemonic: %s\n", mnemonic))
				}
			} else {
				b.WriteString("Mnemonic: Authentication required\n")
			}
		} else {
			b.WriteString("Mnemonic: Not available (imported from private key)\n")
		}
	} else {
		b.WriteString("Mnemonic: *** *** *** *** *** *** *** *** *** *** *** ***\n")
	}

	b.WriteString("\n")

	// Render multi-network balances
	if m.currentMultiBalance != nil {
		balanceStyle := lipgloss.NewStyle().
			Bold(true)

		networkStyle := lipgloss.NewStyle()
		errorStyle := lipgloss.NewStyle()

		b.WriteString(balanceStyle.Render("🌐 Network Balances:"))
		b.WriteString("\n\n")

		for _, networkBalance := range m.currentMultiBalance.NetworkBalances {
			if networkBalance.Error != nil {
				b.WriteString(networkStyle.Render(fmt.Sprintf("  %s: ", networkBalance.NetworkName)))
				b.WriteString(errorStyle.Render(fmt.Sprintf("Error - %v", networkBalance.Error)))
			} else {
				// Convert wei to ether for display
				balanceFloat := new(big.Float).SetInt(networkBalance.Amount)
				divisor := new(big.Float).SetFloat64(math.Pow(10, float64(networkBalance.Decimals)))
				balanceFloat.Quo(balanceFloat, divisor)

				balanceStr := balanceFloat.Text('f', 6)
				// Remove trailing zeros
				balanceStr = strings.TrimRight(balanceStr, "0")
				balanceStr = strings.TrimRight(balanceStr, ".")

				b.WriteString(networkStyle.Render(fmt.Sprintf("  %s: ", networkBalance.NetworkName)))
				b.WriteString(fmt.Sprintf("%s %s", balanceStr, networkBalance.Symbol))
			}
			b.WriteString("\n")
		}

		updateTime := m.currentMultiBalance.UpdatedAt.Format("2006-01-02 15:04:05")
		b.WriteString("\n")
		b.WriteString(networkStyle.Render(fmt.Sprintf("Last updated: %s", updateTime)))
	} else {
		b.WriteString("Balance: Loading...")
	}

	b.WriteString("\n\n")

	if m.err != nil {
		errorStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("196"))
		b.WriteString(errorStyle.Render("Error: " + m.err.Error()))
		b.WriteString("\n\n")
	}

	footerStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("241"))
	if m.showSensitiveInfo {
		b.WriteString(footerStyle.Render("r: refresh balance • s: hide sensitive info • d: delete wallet • esc: back • q: quit"))
	} else {
		b.WriteString(footerStyle.Render("r: refresh balance • s: show sensitive info • d: delete wallet • esc: back • q: quit"))
	}

	return b.String()
}

func (m Model) renderCreateWallet() string {
	var b strings.Builder

	headerStyle := lipgloss.NewStyle().
		Bold(true).
		MarginBottom(2)

	b.WriteString(headerStyle.Render("➕ Create New Wallet"))
	b.WriteString("\n\n")

	b.WriteString("Fill in the details to create a new wallet:\n\n")

	// Name input
	b.WriteString("Wallet Name:\n")
	b.WriteString(m.nameInput.View())
	b.WriteString("\n\n")

	// Password input
	b.WriteString("Password:\n")
	b.WriteString(m.passwordInput.View())
	b.WriteString("\n\n")

	if m.err != nil {
		errorStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("196"))
		b.WriteString(errorStyle.Render("Error: " + m.err.Error()))
		b.WriteString("\n\n")
	}

	footerStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("241"))
	b.WriteString(footerStyle.Render("tab: next field • enter: create wallet • esc: back • q: quit"))

	return b.String()
}

func (m Model) renderImportWallet() string {
	var b strings.Builder

	headerStyle := lipgloss.NewStyle().
		Bold(true).
		MarginBottom(2)

	b.WriteString(headerStyle.Render("📥 Import Wallet"))
	b.WriteString("\n\n")

	b.WriteString("Fill in the details to import an existing wallet:\n\n")

	// Name input
	b.WriteString("Wallet Name:\n")
	b.WriteString(m.nameInput.View())
	b.WriteString("\n\n")

	// Password input
	b.WriteString("Password:\n")
	b.WriteString(m.passwordInput.View())
	b.WriteString("\n\n")

	// Mnemonic input
	b.WriteString("Mnemonic Phrase (12 words):\n")
	b.WriteString(m.mnemonicInput.View())
	b.WriteString("\n\n")

	if m.err != nil {
		errorStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("196"))
		b.WriteString(errorStyle.Render("Error: " + m.err.Error()))
		b.WriteString("\n\n")
	}

	footerStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("241"))
	b.WriteString(footerStyle.Render("tab: next field • enter: import wallet • esc: back • q: quit"))

	return b.String()
}

func (m Model) renderImportPrivateKey() string {
	var b strings.Builder

	headerStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("86")).
		MarginBottom(2)

	b.WriteString(headerStyle.Render("🔑 Import Wallet from Private Key"))
	b.WriteString("\n\n")

	b.WriteString("Fill in the details to import an existing wallet from a private key:\n\n")

	// Name input
	b.WriteString("Wallet Name:\n")
	b.WriteString(m.nameInput.View())
	b.WriteString("\n\n")

	// Password input
	b.WriteString("Password:\n")
	b.WriteString(m.passwordInput.View())
	b.WriteString("\n\n")

	// Private Key input
	b.WriteString("Private Key (without 0x prefix):\n")
	b.WriteString(m.privateKeyInput.View())
	b.WriteString("\n\n")

	if m.err != nil {
		errorStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("196"))
		b.WriteString(errorStyle.Render("Error: " + m.err.Error()))
		b.WriteString("\n\n")
	}

	footerStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("241"))
	b.WriteString(footerStyle.Render("tab: next field • enter: import wallet • esc: back • q: quit"))

	return b.String()
}

// refreshNetworkItems updnetates the network items list
func (m *Model) refreshNetworkItems() {
	var networkItems []string
	networkKeys := m.config.GetAllNetworkKeys()
	for _, key := range networkKeys {
		if network, exists := m.config.GetNetworkByKey(key); exists {
			status := ""
			if network.IsActive {
				status = " (Active)"
			}
			networkItems = append(networkItems, fmt.Sprintf("%s%s", network.Name, status))
		}
	}
	networkItems = append(networkItems, "Add Network", "Back to Settings")
	m.networkItems = networkItems
}
