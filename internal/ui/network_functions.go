package ui

import (
	"blocowallet/internal/constants"
	"blocowallet/pkg/config"
	"blocowallet/pkg/localization"
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	tea "github.com/charmbracelet/bubbletea"
)

// initNetworkList initializes the network list view
func (m *CLIModel) initNetworkList() {
	// Initialize the network list component if it hasn't been initialized yet
	m.networkListComponent = NewNetworkListComponent()

	// Load the current configuration if it hasn't been loaded yet
	if m.currentConfig == nil {
		// Load the current configuration
		appDir := filepath.Join(os.Getenv("HOME"), ".wallets")
		cfg, err := config.LoadConfig(appDir)
		if err != nil {
			m.err = fmt.Errorf("failed to load configuration: %v", err)
			m.currentView = constants.DefaultView
			return
		}

		// Store the current configuration
		m.currentConfig = cfg
	}

	// Initialize Networks map if it's nil
	if m.currentConfig.Networks == nil {
		m.currentConfig.Networks = make(map[string]config.Network)
	}

	// Update the network list with the current networks
	m.networkListComponent.UpdateNetworks(m.currentConfig)

	// Set the current view to the network list view
	m.currentView = constants.NetworkListView
}

// initAddNetwork initializes the add network view
func (m *CLIModel) initAddNetwork() {
	// Initialize the add network component if it hasn't been initialized yet
	m.addNetworkComponent = NewAddNetworkComponent()

	// Load the current configuration if it hasn't been loaded yet
	if m.currentConfig == nil {
		// Load the current configuration
		appDir := filepath.Join(os.Getenv("HOME"), ".wallets")
		cfg, err := config.LoadConfig(appDir)
		if err != nil {
			m.err = fmt.Errorf("failed to load configuration: %v", err)
			m.currentView = constants.DefaultView
			return
		}

		// Store the current configuration
		m.currentConfig = cfg
	}

	// Initialize Networks map if it's nil
	if m.currentConfig.Networks == nil {
		m.currentConfig.Networks = make(map[string]config.Network)
	}

	// Set the current view to the add network view
	m.currentView = constants.AddNetworkView
}

// viewNetworkList renders the network list view
func (m *CLIModel) viewNetworkList() string {
	// Update the component size
	// Only call SetSize if the component is initialized and has rows
	rows := m.networkListComponent.table.Rows()
	if rows != nil && len(rows) > 0 {
		m.networkListComponent.SetSize(m.width, m.height)
	}

	// Render the component
	return m.networkListComponent.View()
}

// viewAddNetwork renders the add network view
func (m *CLIModel) viewAddNetwork() string {
	// Update the component size
	m.addNetworkComponent.SetSize(m.width, m.height)

	// Render the component
	return m.addNetworkComponent.View()
}

// updateNetworkList handles updates to the network list view
func (m *CLIModel) updateNetworkList(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "a":
			// Add a new network
			m.initAddNetwork()
			return m, nil

		case "e":
			// Edit the selected network
			key := m.networkListComponent.GetSelectedNetworkKey()
			if key == "" {
				m.networkListComponent.SetError(fmt.Errorf(localization.Labels["no_network_selected"]))
				return m, nil
			}

			// TODO: Implement network editing
			return m, nil

		case "d":
			// Delete the selected network
			key := m.networkListComponent.GetSelectedNetworkKey()
			if key == "" {
				m.networkListComponent.SetError(fmt.Errorf(localization.Labels["no_network_selected"]))
				return m, nil
			}

			// Ensure currentConfig is initialized
			if m.currentConfig == nil {
				// Load the current configuration
				appDir := filepath.Join(os.Getenv("HOME"), ".wallets")
				cfg, err := config.LoadConfig(appDir)
				if err != nil {
					m.err = fmt.Errorf("failed to load configuration: %v", err)
					m.currentView = constants.DefaultView
					return m, nil
				}

				// Store the current configuration
				m.currentConfig = cfg
			}

			// Ensure Networks map is initialized
			if m.currentConfig.Networks == nil {
				m.currentConfig.Networks = make(map[string]config.Network)
				m.networkListComponent.SetError(fmt.Errorf(localization.Labels["no_network_selected"]))
				return m, nil
			}

			// Remove the network from the configuration
			delete(m.currentConfig.Networks, key)

			// Update the network list
			m.networkListComponent.UpdateNetworks(m.currentConfig)

			return m, nil

		case "esc", "backspace":
			// Return to the network menu
			m.menuItems = NewNetworkMenu()
			m.selectedMenu = 0
			m.currentView = constants.NetworkMenuView
			return m, nil
		}

	case BackToNetworkListMsg:
		// Return to the network list view
		m.currentView = constants.NetworkListView

		// Ensure currentConfig is initialized
		if m.currentConfig == nil {
			// Load the current configuration
			appDir := filepath.Join(os.Getenv("HOME"), ".wallets")
			cfg, err := config.LoadConfig(appDir)
			if err != nil {
				m.err = fmt.Errorf("failed to load configuration: %v", err)
				m.currentView = constants.DefaultView
				return m, nil
			}

			// Store the current configuration
			m.currentConfig = cfg
		}

		// Ensure Networks map is initialized
		if m.currentConfig.Networks == nil {
			m.currentConfig.Networks = make(map[string]config.Network)
		}

		// Update the network list
		m.networkListComponent.UpdateNetworks(m.currentConfig)

		return m, nil
	}

	// Update the network list component
	networkList, cmd := m.networkListComponent.Update(msg)
	m.networkListComponent = *networkList

	return m, cmd
}

// updateAddNetwork handles updates to the add network view
func (m *CLIModel) updateAddNetwork(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case BackToNetworkMenuMsg:
		// Return to the network menu
		m.menuItems = NewNetworkMenu()
		m.selectedMenu = 0
		m.currentView = constants.NetworkMenuView
		return m, nil
	case AddNetworkRequestMsg:
		// Add the network to the configuration
		chainID, err := strconv.ParseInt(msg.ChainID, 10, 64)
		if err != nil {
			m.addNetworkComponent.SetError(fmt.Errorf(localization.Labels["invalid_chain_id"]))
			return m, nil
		}

		// Create a unique key for the network
		key := fmt.Sprintf("custom_%s_%s", msg.Name, msg.ChainID)

		// Ensure currentConfig is initialized
		if m.currentConfig == nil {
			// Load the current configuration
			appDir := filepath.Join(os.Getenv("HOME"), ".wallets")
			cfg, err := config.LoadConfig(appDir)
			if err != nil {
				m.err = fmt.Errorf("failed to load configuration: %v", err)
				m.currentView = constants.DefaultView
				return m, nil
			}

			// Store the current configuration
			m.currentConfig = cfg
		}

		// Ensure Networks map is initialized
		if m.currentConfig.Networks == nil {
			m.currentConfig.Networks = make(map[string]config.Network)
		}

		// Add the network to the configuration
		m.currentConfig.Networks[key] = config.Network{
			Name:        msg.Name,
			RPCEndpoint: msg.RPCEndpoint,
			ChainID:     chainID,
			Symbol:      msg.Symbol,
			IsActive:    true,
		}

		// Initialize the network list component if it hasn't been initialized yet
		if m.networkListComponent.table.Rows() == nil {
			m.networkListComponent = NewNetworkListComponent()
		}

		// Update the network list
		m.networkListComponent.UpdateNetworks(m.currentConfig)

		// Return to the network list view
		m.currentView = constants.NetworkListView

		return m, nil
	}

	// Update the add network component
	addNetwork, cmd := m.addNetworkComponent.Update(msg)
	m.addNetworkComponent = *addNetwork

	return m, cmd
}
