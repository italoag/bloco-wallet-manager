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

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/table"
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

	// Add network inputs with autocomplete
	networkNameInput := textinput.New()
	networkNameInput.Placeholder = "Enter network name..."
	networkNameInput.Width = 40
	networkNameInput.ShowSuggestions = true

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
			customTag := ""
			if network.IsCustom {
				customTag = " [Custom]"
			}
			networkItems = append(networkItems, fmt.Sprintf("%s%s%s", network.Name, status, customTag))
		}
	}
	networkItems = append(networkItems, "Add Custom Network", "Back to Settings")

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

	// Initialize spinner
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("86"))

	// Create the model
	model := Model{
		walletService: walletService,
		config:        cfg,
		wallets:       []*wallet.Wallet{},
		selected:      0,
		loading:       false,
		currentView:   SplashView,

		// Initialize components
		splashComponent:           NewSplashComponent(),
		mainMenuComponent:         NewMainMenuComponent(),
		settingsMenuComponent:     NewSettingsMenuComponent(),
		networkListComponent:      NewNetworkListComponent(cfg),
		addNetworkComponent:       NewAddNetworkComponent(),
		languageMenuComponent:     NewLanguageMenuComponent(cfg),
		walletListComponent:       NewWalletListComponent(),
		balanceComponent:          NewBalanceComponent(),
		createWalletComponent:     NewCreateWalletComponent(),
		importMnemonicComponent:   NewImportMnemonicComponent(),
		importPrivateKeyComponent: NewImportPrivateKeyComponent(),

		// Loading components
		loadingSpinner: s,
		isLoading:      false,
		loadingText:    "",

		// Legacy fields (to be removed gradually)
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

	// Setup wallet table
	model.setupWalletTable()

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

	// Update components based on current view - for form views, let components handle ALL messages first
	var cmd tea.Cmd
	switch m.currentView {
	case MenuView:
		updatedComponent, componentCmd := m.mainMenuComponent.Update(msg)
		m.mainMenuComponent = *updatedComponent
		cmd = componentCmd
	case SettingsView:
		updatedComponent, componentCmd := m.settingsMenuComponent.Update(msg)
		m.settingsMenuComponent = *updatedComponent
		cmd = componentCmd
	case NetworkConfigView:
		updatedComponent, componentCmd := m.networkListComponent.Update(msg)
		m.networkListComponent = *updatedComponent
		cmd = componentCmd
	case LanguageView:
		updatedComponent, componentCmd := m.languageMenuComponent.Update(msg)
		m.languageMenuComponent = *updatedComponent
		cmd = componentCmd
	case AddNetworkView:
		updatedComponent, componentCmd := m.addNetworkComponent.Update(msg)
		m.addNetworkComponent = *updatedComponent
		cmd = componentCmd
	case CreateWalletView:
		updatedComponent, componentCmd := m.createWalletComponent.Update(msg)
		m.createWalletComponent = *updatedComponent
		cmd = componentCmd
	case ImportWalletView:
		updatedComponent, componentCmd := m.importMnemonicComponent.Update(msg)
		m.importMnemonicComponent = *updatedComponent
		cmd = componentCmd
	case ImportPrivateKeyView:
		updatedComponent, componentCmd := m.importPrivateKeyComponent.Update(msg)
		m.importPrivateKeyComponent = *updatedComponent
		cmd = componentCmd
	}

	switch msg := msg.(type) {
	case tea.MouseMsg:
		// Update bubblezone with mouse events
		return zone.AnyInBoundsAndUpdate(m, msg)

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		// Update components with new size
		m.mainMenuComponent.SetSize(msg.Width, msg.Height)
		m.settingsMenuComponent.SetSize(msg.Width, msg.Height)
		m.networkListComponent.SetSize(msg.Width, msg.Height)
		m.addNetworkComponent.SetSize(msg.Width, msg.Height)
		m.languageMenuComponent.SetSize(msg.Width, msg.Height)
		m.createWalletComponent.SetSize(msg.Width, msg.Height)
		m.importMnemonicComponent.SetSize(msg.Width, msg.Height)
		m.importPrivateKeyComponent.SetSize(msg.Width, msg.Height)
		// Update table size when window resizes
		m.setupWalletTable()
		m.updateWalletTable()
		return m, cmd

	case MenuItemSelectedMsg:
		return m.handleMenuSelection(msg)

	case SettingsItemSelectedMsg:
		return m.handleSettingsSelection(msg)

	case BackToSettingsMsg:
		m.currentView = SettingsView
		return m, nil

	case NetworkAddRequestMsg:
		m.currentView = AddNetworkView
		m.addNetworkComponent.Reset()
		return m, m.addNetworkComponent.Init()

	case AddNetworkRequestMsg:
		return m.handleAddNetworkRequest(msg)

	case BackToNetworkListMsg:
		m.currentView = NetworkConfigView
		m.networkListComponent.RefreshNetworks()
		return m, nil

	case NetworkToggleMsg:
		return m.handleNetworkToggle(msg)

	case NetworkEditMsg:
		return m.handleNetworkEdit(msg)

	case NetworkDeleteMsg:
		return m.handleNetworkDelete(msg)

	case NetworkSelectedMsg:
		return m.handleNetworkSelected(msg)

	case LanguageSelectedMsg:
		return m.handleLanguageSelected(msg)

	case CreateWalletRequestMsg:
		return m.handleCreateWalletRequest(msg)

	case ImportMnemonicRequestMsg:
		return m.handleImportMnemonicRequest(msg)

	case ImportPrivateKeyRequestMsg:
		return m.handleImportPrivateKeyRequest(msg)

	case BackToMenuMsg:
		m.currentView = MenuView
		return m, cmd

	case tea.KeyMsg:
		// Only call handleKeyMsg for views that don't use modern components
		// or for keys that components don't handle
		if m.currentView == MenuView || m.currentView == SettingsView ||
			m.currentView == NetworkConfigView || m.currentView == LanguageView ||
			m.currentView == AddNetworkView || m.currentView == CreateWalletView ||
			m.currentView == ImportWalletView || m.currentView == ImportPrivateKeyView {
			// These views use modern components that handle their own keys
			// Only override for escape key and other special cases
			if msg.String() == "esc" || msg.String() == "q" || msg.String() == "ctrl+c" {
				return m.handleKeyMsg(msg)
			}
			// For other keys, let the component command execute
			return m, cmd
		}
		return m.handleKeyMsg(msg)

	case walletsLoadedMsg:
		m.wallets = []*wallet.Wallet(msg)
		m.loading = false
		m.err = nil
		m.stopLoading()
		m.updateWalletTable()
		return m, nil

	case balanceLoadedMsg:
		m.currentBalance = (*wallet.Balance)(msg)
		return m, nil

	case multiBalanceLoadedMsg:
		m.currentMultiBalance = (*wallet.MultiNetworkBalance)(msg)
		m.stopLoading()
		return m, nil

	case walletCreatedMsg:
		newModel := m
		newModel.currentView = WalletListView
		newModel.selected = 0
		// Reset components
		newModel.createWalletComponent.Reset()
		newModel.importMnemonicComponent.Reset()
		newModel.importPrivateKeyComponent.Reset()
		return newModel, newModel.loadWalletsCmd()

	case errorMsg:
		m.err = fmt.Errorf("%s", string(msg))
		m.loading = false
		// Pass error to current component
		switch m.currentView {
		case CreateWalletView:
			m.createWalletComponent.SetError(m.err)
		case ImportWalletView:
			m.importMnemonicComponent.SetError(m.err)
		case ImportPrivateKeyView:
			m.importPrivateKeyComponent.SetError(m.err)
		}
		return m, nil

	case networkSuggestionsMsg:
		m.networkSuggestions = []blockchain.NetworkSuggestion(msg)
		m.showingSuggestions = len(m.networkSuggestions) > 0
		m.selectedSuggestion = 0
		m.stopLoading()

		// Update autocomplete suggestions
		suggestions := make([]string, len(m.networkSuggestions))
		for i, suggestion := range m.networkSuggestions {
			suggestions[i] = suggestion.Name
		}
		m.networkNameInput.SetSuggestions(suggestions)
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
		network, ok := msg.network.(config.Network)
		if !ok {
			m.err = fmt.Errorf("invalid network type")
			return m, nil
		}
		m.config.AddCustomNetwork(msg.key, network)
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

	case spinner.TickMsg:
		if m.isLoading {
			var cmd tea.Cmd
			m.loadingSpinner, cmd = m.loadingSpinner.Update(msg)
			return m, cmd
		}
	}

	return m, cmd
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
			// Reset component state instead of legacy input handling
			newModel.createWalletComponent.Reset()
			newModel.importMnemonicComponent.Reset()
			newModel.importPrivateKeyComponent.Reset()
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
						startLoadingCmd := newModel.startLoading("Searching networks...")
						return newModel, tea.Batch(cmd, startLoadingCmd, newModel.searchNetworksCmd(newValue))
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
					if newModel.currentView == ImportWalletView {
						newModel.mnemonicInput, cmd = newModel.mnemonicInput.Update(msg)
					} else if newModel.currentView == ImportPrivateKeyView {
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
			// Reset component state instead of legacy input handling
			m.createWalletComponent.Reset()
			m.importMnemonicComponent.Reset()
			m.importPrivateKeyComponent.Reset()
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
		if m.currentView == WalletListView {
			if len(m.wallets) > 0 {
				var cmd tea.Cmd
				m.walletTable, cmd = m.walletTable.Update(msg)
				return m, cmd
			}
		} else if m.currentView == SettingsView {
			if m.settingsSelected > 0 {
				m.settingsSelected--
			}
		} else if m.currentView == NetworkConfigView {
			if !m.editingRPC && m.networkSelected > 0 {
				m.networkSelected--
			}
		} else if m.currentView == LanguageView {
			if m.languageSelected > 0 {
				m.languageSelected--
			}
		} else if m.currentView == AddNetworkView {
			if !m.addingNetwork && m.addNetworkFocus > 0 {
				m.addNetworkFocus--
				m.updateAddNetworkFocus()
			}
		}
		return m, nil

	case "down", "j":
		maxItems := 0
		if m.currentView == WalletListView {
			maxItems = len(m.wallets) - 1
		} else if m.currentView == SettingsView {
			maxItems = len(m.settingsItems) - 1
		} else if m.currentView == NetworkConfigView {
			if !m.editingRPC {
				maxItems = len(m.networkItems) - 1
			}
		} else if m.currentView == LanguageView {
			maxItems = len(m.languageItems) - 1
		} else if m.currentView == AddNetworkView {
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
		} else if m.currentView == WalletListView && len(m.wallets) > 0 {
			var cmd tea.Cmd
			m.walletTable, cmd = m.walletTable.Update(msg)
			return m, cmd
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

	case "r":
		// Refresh functionality
		if m.currentView == WalletListView {
			return m, tea.Batch(m.startLoading("Refreshing wallets..."), m.loadWalletsCmd())
		} else if m.currentView == WalletDetailsView && m.selectedWallet != nil {
			return m, tea.Batch(m.startLoading("Refreshing balance..."), m.getMultiBalanceCmd(m.selectedWallet.Address))
		}
		return m, nil
	}

	return m, nil
}

// handleMenuSelection handles menu item selection messages
func (m Model) handleMenuSelection(msg MenuItemSelectedMsg) (tea.Model, tea.Cmd) {
	switch msg.Index {
	case 0: // View Wallets
		m.currentView = WalletListView
		return m, m.loadWalletsCmd()
	case 1: // Create New Wallet
		m.currentView = CreateWalletView
		m.createWalletComponent.Reset()
		return m, m.createWalletComponent.Init()
	case 2: // Import Wallet from Mnemonic
		m.currentView = ImportWalletView
		m.importMnemonicComponent.Reset()
		return m, m.importMnemonicComponent.Init()
	case 3: // Import Wallet from Private Key
		m.currentView = ImportPrivateKeyView
		m.importPrivateKeyComponent.Reset()
		return m, m.importPrivateKeyComponent.Init()
	case 4: // Settings
		m.currentView = SettingsView
		return m, nil
	case 5: // Exit
		return m, tea.Quit
	}
	return m, nil
}

// handleSettingsSelection handles settings menu item selection messages
func (m Model) handleSettingsSelection(msg SettingsItemSelectedMsg) (tea.Model, tea.Cmd) {
	switch msg.Index {
	case 0: // Network Configuration
		m.currentView = NetworkConfigView
		return m, nil
	case 1: // Language
		m.currentView = LanguageView
		return m, nil
	case 2: // Back to Main Menu
		m.currentView = MenuView
		return m, nil
	}
	return m, nil
}

// handleCreateWalletRequest handles wallet creation requests
func (m Model) handleCreateWalletRequest(msg CreateWalletRequestMsg) (tea.Model, tea.Cmd) {
	m.createWalletComponent.SetCreating(true)
	return m, m.createWalletCmd(msg.Name, msg.Password)
}

// handleImportMnemonicRequest handles mnemonic import requests
func (m Model) handleImportMnemonicRequest(msg ImportMnemonicRequestMsg) (tea.Model, tea.Cmd) {
	m.importMnemonicComponent.SetImporting(true)
	return m, m.importWalletCmd(msg.Name, msg.Password, msg.Mnemonic)
}

// handleImportPrivateKeyRequest handles private key import requests
func (m Model) handleImportPrivateKeyRequest(msg ImportPrivateKeyRequestMsg) (tea.Model, tea.Cmd) {
	m.importPrivateKeyComponent.SetImporting(true)
	return m, m.importWalletFromPrivateKeyCmd(msg.Name, msg.Password, msg.PrivateKey)
}

// handleAddNetworkRequest handles add network requests
func (m Model) handleAddNetworkRequest(msg AddNetworkRequestMsg) (tea.Model, tea.Cmd) {
	m.addNetworkComponent.SetAdding(true)
	return m, m.addNetworkCmd(msg.Name, msg.ChainID, msg.RPCEndpoint)
}

// handleNetworkToggle handles network toggle requests
func (m Model) handleNetworkToggle(msg NetworkToggleMsg) (tea.Model, tea.Cmd) {
	err := m.config.ToggleNetworkActive(msg.Key)
	if err == nil {
		// Refresh multi-provider with updated network configuration
		m.walletService.RefreshMultiProvider(m.config)
		// Save configuration
		_ = m.config.Save()
		// Refresh network list
		m.networkListComponent.RefreshNetworks()
	}
	return m, nil
}

// handleNetworkEdit handles network edit requests
func (m Model) handleNetworkEdit(msg NetworkEditMsg) (tea.Model, tea.Cmd) {
	// For now, we don't have a separate edit view, so this is a placeholder
	// In the future, this could open an edit form similar to add network
	return m, nil
}

// handleNetworkDelete handles network delete requests
func (m Model) handleNetworkDelete(msg NetworkDeleteMsg) (tea.Model, tea.Cmd) {
	m.config.RemoveCustomNetwork(msg.Key)
	// Save configuration
	_ = m.config.Save()
	// Refresh network list
	m.networkListComponent.RefreshNetworks()
	return m, nil
}

// handleNetworkSelected handles network selection for details view
func (m Model) handleNetworkSelected(msg NetworkSelectedMsg) (tea.Model, tea.Cmd) {
	switch msg.Key {
	case "add-network":
		m.currentView = AddNetworkView
		m.addNetworkComponent.Reset()
		return m, m.addNetworkComponent.Init()
	case "back-to-settings":
		m.currentView = SettingsView
		return m, nil
	default:
		// Handle regular network selection - for now just show details
		// In the future this could open a network details view
		return m, nil
	}
}

// handleLanguageSelected handles language selection
func (m Model) handleLanguageSelected(msg LanguageSelectedMsg) (tea.Model, tea.Cmd) {
	switch msg.Code {
	case "back-to-settings":
		m.currentView = SettingsView
		return m, nil
	default:
		// Change language
		m.config.Language = msg.Code
		if err := m.config.Save(); err == nil {
			// Refresh language menu to show new selection
			m.languageMenuComponent.RefreshLanguages()
		}
		// Go back to settings
		m.currentView = SettingsView
		return m, nil
	}
}

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
		baseView = m.mainMenuComponent.View()
	case SettingsView:
		baseView = m.settingsMenuComponent.View()
	case WalletListView:
		baseView = m.renderWalletList()
	case WalletAuthView:
		baseView = m.renderWalletAuth()
	case WalletDetailsView:
		baseView = m.renderWalletDetails()
	case CreateWalletView:
		baseView = m.createWalletComponent.View()
	case ImportWalletView:
		baseView = m.importMnemonicComponent.View()
	case ImportPrivateKeyView:
		baseView = m.importPrivateKeyComponent.View()
	case NetworkConfigView:
		baseView = m.networkListComponent.View()
	case LanguageView:
		baseView = m.languageMenuComponent.View()
	case AddNetworkView:
		baseView = m.addNetworkComponent.View()
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
		projectInfo := fmt.Sprintf("%s v%s", "BLOCO Wallet Manager", "0.2.0")

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

	b.WriteString(subtitleStyle.Render("Your Blockchain Wallet"))
	b.WriteString("\n\n")

	instructionStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("86")).
		Align(lipgloss.Center).
		Width(m.width)

	b.WriteString(instructionStyle.Render("Press any key to continue..."))

	return b.String()
}

func (m Model) renderWalletList() string {
	var contentParts []string

	// Header
	headerStyle := HeaderStyle
	contentParts = append(contentParts, headerStyle.Render("📋 Your Wallets"))
	contentParts = append(contentParts, "")

	// Main content
	if len(m.wallets) == 0 {
		contentParts = append(contentParts, InfoStyle.Render("No wallets found. Create a new wallet to get started."))
	} else {
		// Use table view for wallets
		contentParts = append(contentParts, m.walletTable.View())
	}

	// Error display
	if m.err != nil {
		contentParts = append(contentParts, "")
		contentParts = append(contentParts, ErrorStyle.Render("Error: "+m.err.Error()))
	}

	// Create main content
	mainContent := lipgloss.JoinVertical(lipgloss.Left, contentParts...)

	// Footer
	footerText := "↑/↓: navigate • enter: view details • r: refresh • esc: back • q: quit"
	footer := FooterStyle.Render(footerText)

	// Loading status (appears at bottom)
	var loadingStatus string
	if m.loading || m.isLoading {
		loadingText := m.loadingText
		if loadingText == "" {
			loadingText = "Loading wallets..."
		}
		loadingStatus = LoadingStyle.Render(fmt.Sprintf("%s %s", m.loadingSpinner.View(), loadingText))
	}

	// Combine all parts with proper spacing
	var finalParts []string
	finalParts = append(finalParts, mainContent)

	// Add spacing before footer/loading
	availableHeight := m.height - lipgloss.Height(mainContent) - 4
	if availableHeight > 0 {
		padding := strings.Repeat("\n", availableHeight)
		finalParts = append(finalParts, padding)
	}

	// Add loading status if present
	if loadingStatus != "" {
		finalParts = append(finalParts, loadingStatus)
	}

	// Add footer
	finalParts = append(finalParts, footer)

	return lipgloss.JoinVertical(lipgloss.Left, finalParts...)
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

func (m Model) renderWalletDetails() string {
	var contentParts []string

	if m.selectedWallet == nil {
		contentParts = append(contentParts, HeaderStyle.Render("💳 Wallet Details"))
		contentParts = append(contentParts, "")
		contentParts = append(contentParts, InfoStyle.Render("No wallet selected."))
		return lipgloss.JoinVertical(lipgloss.Left, contentParts...)
	}

	// Header
	contentParts = append(contentParts, HeaderStyle.Render("💳 "+m.selectedWallet.Name))
	contentParts = append(contentParts, "")

	// Basic info
	info := fmt.Sprintf("Address: %s\nCreated: %s",
		m.selectedWallet.Address,
		m.selectedWallet.CreatedAt.Format("2006-01-02 15:04:05"))
	contentParts = append(contentParts, ValueStyle.Render(info))
	contentParts = append(contentParts, "")

	// Sensitive information section
	contentParts = append(contentParts, SensitiveStyle.Render("🔐 Sensitive Information:"))

	// Private Key
	if m.showSensitiveInfo {
		if m.selectedWallet.EncryptedMnemonic != "" {
			password, hasPassword := m.walletService.GetWalletPassword(m.selectedWallet.Address)
			if hasPassword {
				mnemonic, err := m.walletService.GetMnemonicFromWallet(m.selectedWallet, password)
				if err != nil {
					contentParts = append(contentParts, ErrorStyle.Render(fmt.Sprintf("Private Key: Error decrypting mnemonic - %v", err)))
				} else {
					privateKey, err := wallet.DerivePrivateKey(mnemonic)
					if err != nil {
						contentParts = append(contentParts, ErrorStyle.Render(fmt.Sprintf("Private Key: Error deriving key - %v", err)))
					} else {
						contentParts = append(contentParts, ValueStyle.Render(fmt.Sprintf("Private Key: %s", privateKey)))
					}
				}
			} else {
				contentParts = append(contentParts, InfoStyle.Render("Private Key: Authentication required"))
			}
		} else if m.selectedWallet.KeyStorePath != "" {
			if m.extractedPrivateKey != "" {
				contentParts = append(contentParts, ValueStyle.Render(fmt.Sprintf("Private Key: %s", m.extractedPrivateKey)))
			} else {
				if cachedPassword, hasPassword := m.walletService.GetWalletPassword(m.selectedWallet.Address); hasPassword {
					privateKey, err := m.walletService.LoadPrivateKeyFromKeyStoreV3(m.selectedWallet.KeyStorePath, cachedPassword)
					if err != nil {
						contentParts = append(contentParts, ErrorStyle.Render(fmt.Sprintf("Private Key: Error - %s", err.Error())))
					} else {
						privateKeyHex := fmt.Sprintf("%x", privateKey.D.Bytes())
						m.extractedPrivateKey = privateKeyHex
						contentParts = append(contentParts, ValueStyle.Render(fmt.Sprintf("Private Key: %s", privateKeyHex)))
					}
				} else {
					contentParts = append(contentParts, InfoStyle.Render("Private Key: Authentication required"))
				}
			}
		} else {
			contentParts = append(contentParts, InfoStyle.Render("Private Key: Not available (no keystore or mnemonic found)"))
		}
	} else {
		contentParts = append(contentParts, InfoStyle.Render("Private Key: ********************************"))
	}

	// Mnemonic
	if m.showSensitiveInfo {
		if m.selectedWallet.EncryptedMnemonic != "" {
			password, hasPassword := m.walletService.GetWalletPassword(m.selectedWallet.Address)
			if hasPassword {
				mnemonic, err := m.walletService.GetMnemonicFromWallet(m.selectedWallet, password)
				if err != nil {
					contentParts = append(contentParts, ErrorStyle.Render(fmt.Sprintf("Mnemonic: Error decrypting - %v", err)))
				} else {
					contentParts = append(contentParts, ValueStyle.Render(fmt.Sprintf("Mnemonic: %s", mnemonic)))
				}
			} else {
				contentParts = append(contentParts, InfoStyle.Render("Mnemonic: Authentication required"))
			}
		} else {
			contentParts = append(contentParts, InfoStyle.Render("Mnemonic: Not available (imported from private key)"))
		}
	} else {
		contentParts = append(contentParts, InfoStyle.Render("Mnemonic: *** *** *** *** *** *** *** *** *** *** *** ***"))
	}

	contentParts = append(contentParts, "")

	// Render multi-network balances
	if m.currentMultiBalance != nil {
		contentParts = append(contentParts, BalanceStyle.Render("🌐 Network Balances:"))
		contentParts = append(contentParts, "")

		for _, networkBalance := range m.currentMultiBalance.NetworkBalances {
			if networkBalance.Error != nil {
				networkLine := fmt.Sprintf("  %s: ", networkBalance.NetworkName)
				errorLine := fmt.Sprintf("Error - %v", networkBalance.Error)
				contentParts = append(contentParts, NetworkStyle.Render(networkLine)+ErrorStyle.Render(errorLine))
			} else {
				balanceFloat := new(big.Float).SetInt(networkBalance.Amount)
				divisor := new(big.Float).SetFloat64(math.Pow(10, float64(networkBalance.Decimals)))
				balanceFloat.Quo(balanceFloat, divisor)

				balanceStr := balanceFloat.Text('f', 6)
				balanceStr = strings.TrimRight(balanceStr, "0")
				balanceStr = strings.TrimRight(balanceStr, ".")

				networkLine := fmt.Sprintf("  %s: ", networkBalance.NetworkName)
				balanceLine := fmt.Sprintf("%s %s", balanceStr, networkBalance.Symbol)
				contentParts = append(contentParts, NetworkStyle.Render(networkLine)+ValueStyle.Render(balanceLine))
			}
		}

		updateTime := m.currentMultiBalance.UpdatedAt.Format("2006-01-02 15:04:05")
		contentParts = append(contentParts, "")
		contentParts = append(contentParts, NetworkStyle.Render(fmt.Sprintf("Last updated: %s", updateTime)))
	} else {
		contentParts = append(contentParts, InfoStyle.Render("Balance: Loading..."))
	}

	// Error display
	if m.err != nil {
		contentParts = append(contentParts, "")
		contentParts = append(contentParts, ErrorStyle.Render("Error: "+m.err.Error()))
	}

	// Create main content
	mainContent := lipgloss.JoinVertical(lipgloss.Left, contentParts...)

	// Footer
	var footerText string
	if m.showSensitiveInfo {
		footerText = "r: refresh balance • s: hide sensitive info • d: delete wallet • esc: back • q: quit"
	} else {
		footerText = "r: refresh balance • s: show sensitive info • d: delete wallet • esc: back • q: quit"
	}
	footer := FooterStyle.Render(footerText)

	// Loading status (appears at bottom)
	var loadingStatus string
	if m.isLoading {
		loadingText := m.loadingText
		if loadingText == "" {
			loadingText = "Loading..."
		}
		loadingStatus = LoadingStyle.Render(fmt.Sprintf("%s %s", m.loadingSpinner.View(), loadingText))
	}

	// Combine all parts with proper spacing
	var finalParts []string
	finalParts = append(finalParts, mainContent)

	// Add spacing before footer/loading
	availableHeight := m.height - lipgloss.Height(mainContent) - 4
	if availableHeight > 0 {
		padding := strings.Repeat("\n", availableHeight)
		finalParts = append(finalParts, padding)
	}

	// Add loading status if present
	if loadingStatus != "" {
		finalParts = append(finalParts, loadingStatus)
	}

	// Add footer
	finalParts = append(finalParts, footer)

	return lipgloss.JoinVertical(lipgloss.Left, finalParts...)
}

// refreshNetworkItems updates the network items list
func (m *Model) refreshNetworkItems() {
	var networkItems []string
	networkKeys := m.config.GetAllNetworkKeys()
	for _, key := range networkKeys {
		if network, exists := m.config.GetNetworkByKey(key); exists {
			status := ""
			if network.IsActive {
				status = " (Active)"
			}
			customTag := ""
			if network.IsCustom {
				customTag = " [Custom]"
			}
			networkItems = append(networkItems, fmt.Sprintf("%s%s%s", network.Name, status, customTag))
		}
	}
	networkItems = append(networkItems, "Add Custom Network", "Back to Settings")
	m.networkItems = networkItems
}

// startLoading starts the loading spinner with a message
func (m *Model) startLoading(text string) tea.Cmd {
	m.isLoading = true
	m.loadingText = text
	return m.loadingSpinner.Tick
}

// stopLoading stops the loading spinner
func (m *Model) stopLoading() {
	m.isLoading = false
	m.loadingText = ""
}

// setupWalletTable initializes the wallet table
func (m *Model) setupWalletTable() {
	// Calculate column widths dynamically based on screen width
	availableWidth := m.width - 4 // Account for borders and padding
	if availableWidth < 80 {
		availableWidth = 80 // Minimum width
	}

	// Distribute width: Name(25%), Address(45%), Type(15%), Created(15%)
	nameWidth := int(float64(availableWidth) * 0.25)
	addressWidth := int(float64(availableWidth) * 0.45)
	typeWidth := int(float64(availableWidth) * 0.15)
	createdWidth := availableWidth - nameWidth - addressWidth - typeWidth

	columns := []table.Column{
		{Title: "Name", Width: nameWidth},
		{Title: "Address", Width: addressWidth},
		{Title: "Type", Width: typeWidth},
		{Title: "Created", Width: createdWidth},
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithRows([]table.Row{}),
		table.WithFocused(true),
		table.WithHeight(10),
		table.WithWidth(availableWidth),
	)

	// Apply consistent styling using existing style definitions
	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("#874BFD")).
		BorderBottom(true).
		Bold(true).
		Foreground(lipgloss.Color("86"))
	s.Selected = s.Selected.
		Foreground(lipgloss.Color("#FFF")).
		Background(lipgloss.Color("#874BFD")).
		Bold(true)
	s.Cell = s.Cell.
		Foreground(lipgloss.Color("250"))
	t.SetStyles(s)

	m.walletTable = t
}

// updateWalletTable updates the wallet table with current wallets
func (m *Model) updateWalletTable() {
	var rows []table.Row
	for _, wallet := range m.wallets {
		walletType := "Mnemonic"
		if wallet.EncryptedMnemonic == "" {
			walletType = "Private Key"
		}

		rows = append(rows, table.Row{
			wallet.Name,
			wallet.Address,
			walletType,
			wallet.CreatedAt.Format("2006-01-02 15:04"),
		})
	}

	m.walletTable.SetRows(rows)
}
