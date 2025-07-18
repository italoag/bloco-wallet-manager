package ui

import (
	"blocowallet/internal/constants"
	"blocowallet/pkg/config"
	"blocowallet/pkg/localization"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

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

			// Get the network to edit
			network, exists := m.currentConfig.Networks[key]
			if !exists {
				m.networkListComponent.SetError(fmt.Errorf("network not found"))
				return m, nil
			}

			// Initialize add network component for editing
			m.addNetworkComponent = NewAddNetworkComponent()

			// Pre-fill the form with existing network data
			m.addNetworkComponent.nameInput.SetValue(network.Name)
			m.addNetworkComponent.chainIDInput.SetValue(strconv.FormatInt(network.ChainID, 10))
			m.addNetworkComponent.symbolInput.SetValue(network.Symbol)
			m.addNetworkComponent.rpcEndpointInput.SetValue(network.RPCEndpoint)

			// Store the key for updating later
			m.editingNetworkKey = key

			// Set the current view to add network (which will function as edit)
			m.currentView = constants.AddNetworkView
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

			// Save the configuration to file
			err := m.saveConfigToFile()
			if err != nil {
				m.networkListComponent.SetError(fmt.Errorf("falha ao salvar configuração: %v", err))
				return m, nil
			}

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

// saveConfigToFile saves the current configuration to the config file
func (m *CLIModel) saveConfigToFile() error {
	if m.currentConfig == nil {
		return fmt.Errorf("no configuration to save")
	}

	// Get the config file path
	configPath := filepath.Join(m.currentConfig.AppDir, "config.toml")

	// Create backup of existing config if it exists
	if _, err := os.Stat(configPath); err == nil {
		backupPath := configPath + ".bak"
		if err := copyFile(configPath, backupPath); err != nil {
			return fmt.Errorf("failed to create backup: %w", err)
		}
	}

	// Ensure the directory exists
	dir := filepath.Dir(configPath)
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	// Criar o conteúdo do arquivo TOML manualmente
	var sb strings.Builder

	// Seção [app]
	sb.WriteString("[app]\n")
	sb.WriteString(fmt.Sprintf("language = %q\n", m.currentConfig.Language))
	sb.WriteString(fmt.Sprintf("app_dir = %q\n", m.currentConfig.AppDir))
	sb.WriteString(fmt.Sprintf("wallets_dir = %q\n", m.currentConfig.WalletsDir))
	sb.WriteString(fmt.Sprintf("database_path = %q\n", m.currentConfig.DatabasePath))
	sb.WriteString(fmt.Sprintf("locale_dir = %q\n", m.currentConfig.LocaleDir))
	sb.WriteString("\n")

	// Seção [database]
	sb.WriteString("[database]\n")
	sb.WriteString(fmt.Sprintf("type = %q\n", m.currentConfig.Database.Type))
	sb.WriteString(fmt.Sprintf("dsn = %q\n", m.currentConfig.Database.DSN))
	sb.WriteString("\n")

	// Seção [security]
	sb.WriteString("[security]\n")
	sb.WriteString(fmt.Sprintf("argon2_time = %d\n", m.currentConfig.Security.Argon2Time))
	sb.WriteString(fmt.Sprintf("argon2_memory = %d\n", m.currentConfig.Security.Argon2Memory))
	sb.WriteString(fmt.Sprintf("argon2_threads = %d\n", m.currentConfig.Security.Argon2Threads))
	sb.WriteString(fmt.Sprintf("argon2_key_len = %d\n", m.currentConfig.Security.Argon2KeyLen))
	sb.WriteString(fmt.Sprintf("salt_length = %d\n", m.currentConfig.Security.SaltLength))
	sb.WriteString("\n")

	// Seção [fonts]
	sb.WriteString("[fonts]\n")
	sb.WriteString("available = [")
	for i, font := range m.currentConfig.Fonts {
		if i > 0 {
			sb.WriteString(", ")
		}
		sb.WriteString(fmt.Sprintf("%q", font))
	}
	sb.WriteString("]\n\n")

	// Seção [networks]
	if len(m.currentConfig.Networks) > 0 {
		sb.WriteString("[networks]\n")

		// Adicionar cada rede como uma subseção
		first := true
		for key, network := range m.currentConfig.Networks {
			if !first {
				sb.WriteString("\n")
			}
			first = false

			// Sanitizar a chave para garantir que seja válida para TOML
			sanitizedKey := sanitizeNetworkKey(key)

			sb.WriteString(fmt.Sprintf("[networks.%s]\n", sanitizedKey))
			sb.WriteString(fmt.Sprintf("name = %q\n", network.Name))
			sb.WriteString(fmt.Sprintf("rpc_endpoint = %q\n", network.RPCEndpoint))
			sb.WriteString(fmt.Sprintf("chain_id = %d\n", network.ChainID))
			sb.WriteString(fmt.Sprintf("symbol = %q\n", network.Symbol))
			sb.WriteString(fmt.Sprintf("explorer = %q\n", network.Explorer))
			sb.WriteString(fmt.Sprintf("is_active = %t\n", network.IsActive))
		}
	}

	// Escrever o conteúdo no arquivo
	err := os.WriteFile(configPath, []byte(sb.String()), 0644)
	if err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	// Validar o arquivo após a escrita
	_, err = config.LoadConfig(m.currentConfig.AppDir)
	if err != nil {
		// Se houver erro na validação, restaurar o backup
		backupPath := configPath + ".bak"
		if _, statErr := os.Stat(backupPath); statErr == nil {
			if restoreErr := copyFile(backupPath, configPath); restoreErr != nil {
				return fmt.Errorf("failed to validate config and restore backup: %w (original error: %v)", restoreErr, err)
			}
		}
		return fmt.Errorf("failed to validate config file: %w", err)
	}

	return nil
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
		// Replace spaces and special characters with underscores to create a valid TOML key
		sanitizedName := strings.ReplaceAll(msg.Name, " ", "_")
		sanitizedName = strings.Map(func(r rune) rune {
			if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') || r == '_' {
				return r
			}
			return '_'
		}, sanitizedName)
		key := fmt.Sprintf("custom_%s_%d", sanitizedName, chainID)

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

		// Save the configuration to file
		err = m.saveConfigToFile()
		if err != nil {
			m.addNetworkComponent.SetError(fmt.Errorf("failed to save configuration: %v", err))
			return m, nil
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
