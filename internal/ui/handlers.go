package ui

import (
	"fmt"
	"strings"

	"blocowallet/pkg/config"

	tea "github.com/charmbracelet/bubbletea"
)

// handleInputNavigation handles navigation between input fields
func (m Model) handleInputNavigation(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "tab", "down":
		m.inputFocus = (m.inputFocus + 1) % 3
		if m.currentView == CreateWalletView {
			m.inputFocus = m.inputFocus % 2 // Only 2 fields in create view
		}
	case "shift+tab", "up":
		m.inputFocus--
		if m.inputFocus < 0 {
			if m.currentView == CreateWalletView {
				m.inputFocus = 1
			} else {
				m.inputFocus = 2
			}
		}
	}

	m.updateInputFocus()

	return m, nil
}

// handleInputSubmit handles form submission in input views
func (m Model) handleInputSubmit() (tea.Model, tea.Cmd) {
	if m.currentView == CreateWalletView {
		name := strings.TrimSpace(m.nameInput.Value())
		password := strings.TrimSpace(m.passwordInput.Value())

		if name == "" || password == "" {
			m.err = fmt.Errorf("name and password are required")
			return m, nil
		}

		return m, m.createWalletCmd(name, password)
	} else if m.currentView == ImportWalletView {
		name := strings.TrimSpace(m.nameInput.Value())
		password := strings.TrimSpace(m.passwordInput.Value())
		mnemonic := strings.TrimSpace(m.mnemonicInput.Value())

		if name == "" || password == "" {
			m.err = fmt.Errorf("name and password are required")
			return m, nil
		}

		if mnemonic == "" {
			m.err = fmt.Errorf("mnemonic phrase is required")
			return m, nil
		}

		return m, m.importWalletCmd(name, password, mnemonic)
	}

	return m, nil
}

// handleEnterKey handles the enter key press for navigation
func (m Model) handleEnterKey() (tea.Model, tea.Cmd) {
	switch m.currentView {
	case SplashView:
		m.currentView = MenuView
		m.selected = 0
		return m, nil

	case MenuView:
		switch m.selected {
		case 0: // View Wallets
			m.currentView = WalletListView
			m.selected = 0
			return m, m.loadWalletsCmd()
		case 1: // Create New Wallet
			m.currentView = CreateWalletView
			m.resetInputs()
			m.nameInput.Focus()
			return m, nil
		case 2: // Import Wallet
			m.currentView = ImportWalletView
			m.resetInputs()
			m.nameInput.Focus()
			return m, nil
		case 3: // Settings
			m.currentView = SettingsView
			m.settingsSelected = 0
			return m, nil
		case 4: // Exit
			return m, tea.Quit
		}

	case WalletListView:
		if len(m.wallets) > 0 && m.selected < len(m.wallets) {
			m.selectedWallet = m.wallets[m.selected]
			m.currentView = WalletDetailsView
			return m, m.getMultiBalanceCmd(m.selectedWallet.Address)
		}

	case SettingsView:
		switch m.settingsSelected {
		case 0: // Network Configuration
			m.currentView = NetworkConfigView
			m.networkSelected = 0
			return m, nil
		case 1: // Language
			m.currentView = LanguageView
			m.languageSelected = 0
			return m, nil
		case 2: // Back to Main Menu
			m.currentView = MenuView
			m.selected = 0
			return m, nil
		}

	case NetworkConfigView:
		if m.editingRPC {
			// Save RPC change
			m.editingRPC = false
			m.rpcInput.Blur()
			return m.saveRPCEndpoint()
		} else {
			networkKeys := m.config.GetAllNetworkKeys()
			if m.networkSelected < len(networkKeys) {
				// View network details
				key := networkKeys[m.networkSelected]
				m.selectedNetworkKey = key
				m.currentView = NetworkDetailsView
				return m, nil
			} else if m.networkSelected == len(networkKeys) {
				// Add Custom Network
				m.currentView = AddNetworkView
				m.addNetworkFocus = 0
				m.networkNameInput.Focus()
				m.networkNameInput.SetValue("")
				m.chainIDInput.SetValue("")
				m.rpcEndpointInput.SetValue("")
				return m, nil
			} else if m.networkSelected == len(networkKeys)+1 {
				// Back to Settings
				m.currentView = SettingsView
				m.settingsSelected = 0
				return m, nil
			}
		}

	case LanguageView:
		langCodes := m.config.GetLanguageCodes()
		if m.languageSelected < len(langCodes) {
			// Change language
			langCode := langCodes[m.languageSelected]
			m.config.Language = langCode
			return m.saveLanguageChange()
		} else if m.languageSelected == len(langCodes) {
			// Back to Settings
			m.currentView = SettingsView
			m.settingsSelected = 0
			return m, nil
		}
	}

	return m, nil
}

// saveRPCEndpoint saves the edited RPC endpoint
func (m Model) saveRPCEndpoint() (tea.Model, tea.Cmd) {
	networkKeys := m.config.GetAllNetworkKeys()
	if m.networkSelected < len(networkKeys) {
		key := networkKeys[m.networkSelected]
		if network, exists := m.config.GetNetworkByKey(key); exists {
			network.RPCEndpoint = strings.TrimSpace(m.rpcInput.Value())
			m.config.UpdateNetwork(key, network)

			// Save configuration to file
			if err := m.config.Save(); err != nil {
				m.err = err
			}

			// Update network items
			m.refreshNetworkItems()
		}
	}
	return m, nil
}

// saveLanguageChange saves the language change and updates UI
func (m Model) saveLanguageChange() (tea.Model, tea.Cmd) {
	// Save configuration to file
	if err := m.config.Save(); err != nil {
		m.err = err
	}

	// Update language items to reflect the change
	m.languageItems = m.buildLanguageItems()

	// Go back to settings
	m.currentView = SettingsView
	m.settingsSelected = 0
	return m, nil
}

// buildNetworkItems creates the network items list
func (m Model) buildNetworkItems() []string {
	var items []string
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
			items = append(items, fmt.Sprintf("%s%s%s", network.Name, status, customTag))
		}
	}
	items = append(items, "Add Custom Network", "Back to Settings")
	return items
}

// buildLanguageItems creates the language items list
func (m Model) buildLanguageItems() []string {
	var items []string
	langCodes := m.config.GetLanguageCodes()
	for _, code := range langCodes {
		name := config.SupportedLanguages[code]
		status := ""
		if m.config.Language == code {
			status = " (Current)"
		}
		items = append(items, fmt.Sprintf("%s%s", name, status))
	}
	items = append(items, "Back to Settings")
	return items
}

// updateInputFocus updates which input field has focus
func (m *Model) updateInputFocus() {
	m.nameInput.Blur()
	m.passwordInput.Blur()
	m.mnemonicInput.Blur()

	switch m.inputFocus {
	case 0:
		m.nameInput.Focus()
	case 1:
		m.passwordInput.Focus()
	case 2:
		if m.currentView == ImportWalletView {
			m.mnemonicInput.Focus()
		}
	}
}

// resetInputs resets all input fields to their default state
func (m *Model) resetInputs() {
	m.nameInput.SetValue("")
	m.passwordInput.SetValue("")
	m.mnemonicInput.SetValue("")
	m.inputFocus = 0
	m.nameInput.Focus()
	m.passwordInput.Blur()
	m.mnemonicInput.Blur()
	m.err = nil
}

// handleAddNetworkNavigation handles navigation in add network view
func (m Model) handleAddNetworkNavigation(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "tab", "down":
		m.addNetworkFocus = (m.addNetworkFocus + 1) % 3
	case "shift+tab", "up":
		m.addNetworkFocus--
		if m.addNetworkFocus < 0 {
			m.addNetworkFocus = 2
		}
	}

	m.updateAddNetworkFocus()
	return m, nil
}

// updateAddNetworkFocus updates which input field has focus in add network view
func (m *Model) updateAddNetworkFocus() {
	m.networkNameInput.Blur()
	m.chainIDInput.Blur()
	m.rpcEndpointInput.Blur()

	switch m.addNetworkFocus {
	case 0:
		m.networkNameInput.Focus()
	case 1:
		m.chainIDInput.Focus()
	case 2:
		m.rpcEndpointInput.Focus()
	}
}

// resetAddNetworkInputs resets the inputs in the add network view
func (m *Model) resetAddNetworkInputs() {
	m.networkNameInput.SetValue("")
	m.chainIDInput.SetValue("")
	m.rpcEndpointInput.SetValue("")
	m.addNetworkFocus = 0
	m.networkNameInput.Focus()
	m.chainIDInput.Blur()
	m.rpcEndpointInput.Blur()
	m.err = nil
}

// handleAddNetworkSubmit handles the submission of add network form
func (m Model) handleAddNetworkSubmit() (tea.Model, tea.Cmd) {
	name := strings.TrimSpace(m.networkNameInput.Value())
	chainIDStr := strings.TrimSpace(m.chainIDInput.Value())
	rpcEndpoint := strings.TrimSpace(m.rpcEndpointInput.Value())

	if name == "" || chainIDStr == "" || rpcEndpoint == "" {
		m.err = fmt.Errorf("all fields are required")
		return m, nil
	}

	m.addingNetwork = true
	return m, m.addNetworkCmd(name, chainIDStr, rpcEndpoint)
}
