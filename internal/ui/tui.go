package ui

import (
	"fmt"
	"math"
	"math/big"
	"strings"

	"blocowallet/internal/wallet"
	"blocowallet/pkg/config"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// View constants
const (
	SplashView         = "splash"
	MenuView           = "menu"
	WalletListView     = "wallet_list"
	WalletDetailsView  = "wallet_details"
	CreateWalletView   = "create_wallet"
	ImportWalletView   = "import_wallet"
	SettingsView       = "settings"
	NetworkConfigView  = "network_config"
	LanguageView       = "language"
	AddNetworkView     = "add_network"
	NetworkDetailsView = "network_details"
)

// Model represents the TUI application state
type Model struct {
	walletService  *wallet.Service
	wallets        []*wallet.Wallet
	selected       int
	loading        bool
	err            error
	width          int
	height         int
	currentView    string
	selectedWallet *wallet.Wallet
	currentBalance *wallet.Balance
	currentMultiBalance *wallet.MultiNetworkBalance

	// Configuration
	config *config.Config

	// Input fields for wallet creation/import
	nameInput     textinput.Model
	passwordInput textinput.Model
	mnemonicInput textinput.Model
	inputFocus    int

	// Settings fields
	settingsSelected int
	networkSelected  int
	languageSelected int
	editingRPC       bool
	rpcInput         textinput.Model

	// Add network fields
	networkNameInput   textinput.Model
	chainIDInput       textinput.Model
	rpcEndpointInput   textinput.Model
	addNetworkFocus    int
	addingNetwork      bool
	selectedNetworkKey string

	// Settings menu items
	settingsItems []string
	networkItems  []string
	languageItems []string
}

// Message types
type walletsLoadedMsg []*wallet.Wallet
type balanceLoadedMsg *wallet.Balance
type multiBalanceLoadedMsg *wallet.MultiNetworkBalance
type walletCreatedMsg struct{}
type errorMsg string

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

	return Model{
		walletService:      walletService,
		config:             cfg,
		wallets:            []*wallet.Wallet{},
		selected:           0,
		loading:            false,
		currentView:        SplashView,
		nameInput:          nameInput,
		passwordInput:      passwordInput,
		mnemonicInput:      mnemonicInput,
		inputFocus:         0,
		rpcInput:           rpcInput,
		networkNameInput:   networkNameInput,
		chainIDInput:       chainIDInput,
		rpcEndpointInput:   rpcEndpointInput,
		addNetworkFocus:    0,
		addingNetwork:      false,
		selectedNetworkKey: "",
		settingsSelected:   0,
		networkSelected:    0,
		languageSelected:   0,
		editingRPC:         false,
		settingsItems:      settingsItems,
		networkItems:       networkItems,
		languageItems:      languageItems,
	}
}

// Init initializes the TUI
func (m Model) Init() tea.Cmd {
	return tea.EnterAltScreen
}

// Update handles user input and state changes
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
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
		newModel.inputFocus = 0
		newModel.nameInput.Blur()
		newModel.passwordInput.Blur()
		newModel.mnemonicInput.Blur()
		return newModel, newModel.loadWalletsCmd()

	case errorMsg:
		m.err = fmt.Errorf("%s", string(msg))
		m.loading = false
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
	// Handle input fields in create/import views, RPC editing, and add network
	if m.currentView == CreateWalletView || m.currentView == ImportWalletView ||
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
			}
			newModel.currentView = MenuView
			newModel.selected = 0
			// Reset inputs manually instead of calling resetInputs
			newModel.nameInput.SetValue("")
			newModel.passwordInput.SetValue("")
			newModel.mnemonicInput.SetValue("")
			newModel.inputFocus = 0
			newModel.nameInput.Blur()
			newModel.passwordInput.Blur()
			newModel.mnemonicInput.Blur()
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
				return m.handleAddNetworkSubmit()
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
					newModel.networkNameInput, cmd = newModel.networkNameInput.Update(msg)
				case 1:
					newModel.chainIDInput, cmd = newModel.chainIDInput.Update(msg)
				case 2:
					newModel.rpcEndpointInput, cmd = newModel.rpcEndpointInput.Update(msg)
				}
			} else {
				switch newModel.inputFocus {
				case 0: // Name input
					newModel.nameInput, cmd = newModel.nameInput.Update(msg)
				case 1: // Password input
					newModel.passwordInput, cmd = newModel.passwordInput.Update(msg)
				case 2: // Mnemonic input (only in import view)
					if newModel.currentView == ImportWalletView {
						newModel.mnemonicInput, cmd = newModel.mnemonicInput.Update(msg)
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
		case WalletDetailsView:
			m.currentView = WalletListView
		case CreateWalletView, ImportWalletView:
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
		if m.currentView == MenuView {
			if m.selected > 0 {
				m.selected--
			}
		} else if m.currentView == WalletListView {
			if m.selected > 0 {
				m.selected--
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
		if m.currentView == MenuView {
			maxItems = 4 // 5 menu items (0-4)
		} else if m.currentView == WalletListView {
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
		} else if m.selected < maxItems {
			m.selected++
		}
		return m, nil

	case "a", "A":
		// Activate/deactivate network
		if m.currentView == NetworkConfigView && !m.editingRPC {
			networkKeys := m.config.GetAllNetworkKeys()
			if m.networkSelected < len(networkKeys) {
				key := networkKeys[m.networkSelected]
				if network, exists := m.config.GetNetworkByKey(key); exists {
					// Deactivate all networks first
					for k, net := range m.config.Networks {
						net.IsActive = false
						m.config.Networks[k] = net
					}
					for k, net := range m.config.CustomNetworks {
						net.IsActive = false
						m.config.CustomNetworks[k] = net
					}

					// Activate selected network
					network.IsActive = true
					if network.IsCustom {
						m.config.CustomNetworks[key] = network
					} else {
						m.config.Networks[key] = network
					}

					// Refresh network items to show updated status
					m.refreshNetworkItems()

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

	case "d", "D":
		// Delete custom network (only for custom networks)
		if m.currentView == NetworkConfigView && !m.editingRPC {
			networkKeys := m.config.GetAllNetworkKeys()
			if m.networkSelected < len(networkKeys) {
				key := networkKeys[m.networkSelected]
				if network, exists := m.config.GetNetworkByKey(key); exists && network.IsCustom {
					// Remove custom network
					m.config.RemoveCustomNetwork(key)

					// Refresh network items
					m.refreshNetworkItems()

					// Adjust selection if needed
					if m.networkSelected >= len(m.networkItems)-2 { // -2 for "Add Custom Network" and "Back to Settings"
						m.networkSelected = len(m.networkItems) - 3 // Select last actual network
						if m.networkSelected < 0 {
							m.networkSelected = 0
						}
					}

					// Save configuration
					_ = m.config.Save()
				}
			}
		}
		return m, nil

	case "r":
		if m.currentView == WalletListView {
			return m, m.loadWalletsCmd()
		} else if m.currentView == WalletDetailsView && m.selectedWallet != nil {
			return m, m.getMultiBalanceCmd(m.selectedWallet.Address)
		}
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

	switch m.currentView {
	case SplashView:
		return m.renderSplash()
	case MenuView:
		return m.renderMenu()
	case WalletListView:
		return m.renderWalletList()
	case WalletDetailsView:
		return m.renderWalletDetails()
	case CreateWalletView:
		return m.renderCreateWallet()
	case ImportWalletView:
		return m.renderImportWallet()
	case SettingsView:
		return m.renderSettings()
	case NetworkConfigView:
		return m.renderNetworkConfig()
	case LanguageView:
		return m.renderLanguage()
	case AddNetworkView:
		return m.renderAddNetwork()
	case NetworkDetailsView:
		return m.renderNetworkDetails()
	default:
		return m.renderWalletList()
	}
}

// Render methods
func (m Model) renderSplash() string {
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

	headerStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("86")).
		MarginBottom(2)

	b.WriteString(headerStyle.Render("🏦 BlocoWallet - Main Menu"))
	b.WriteString("\n\n")

	menuItems := []string{
		"📋 View Wallets",
		"➕ Create New Wallet",
		"📥 Import Wallet",
		"⚙️  Settings",
		"🚪 Exit",
	}

	for i, item := range menuItems {
		if i == m.selected {
			selectedStyle := lipgloss.NewStyle().
				Background(lipgloss.Color("86")).
				Foreground(lipgloss.Color("232")).
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
					Background(lipgloss.Color("86")).
					Foreground(lipgloss.Color("232")).
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

func (m Model) renderWalletDetails() string {
	var b strings.Builder

	headerStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("86")).
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

	// Render multi-network balances
	if m.currentMultiBalance != nil {
		balanceStyle := lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("86"))
		
		networkStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("240"))
		
		errorStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("196"))

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
	b.WriteString(footerStyle.Render("r: refresh balance • esc: back • q: quit"))

	return b.String()
}

func (m Model) renderCreateWallet() string {
	var b strings.Builder

	headerStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("86")).
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
		Foreground(lipgloss.Color("86")).
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
