package ui

import (
	"fmt"
	"strconv"
	"strings"

	"blocowallet/internal/blockchain"
	"blocowallet/pkg/config"

	tea "github.com/charmbracelet/bubbletea"
)

// NetworkOperations contém operações específicas de rede que podem ser reutilizadas
type NetworkOperations struct {
	config           *config.Config
	chainListService *blockchain.ChainListService
}

// NewNetworkOperations cria uma nova instância de NetworkOperations
func NewNetworkOperations(cfg *config.Config, chainListService *blockchain.ChainListService) *NetworkOperations {
	return &NetworkOperations{
		config:           cfg,
		chainListService: chainListService,
	}
}

// GetAllNetworks retorna todas as redes configuradas
func (no *NetworkOperations) GetAllNetworks() []string {
	var networkItems []string
	networkKeys := no.config.GetAllNetworkKeys()
	for _, key := range networkKeys {
		if network, exists := no.config.GetNetworkByKey(key); exists {
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
	return networkItems
}

// ToggleNetworkActive alterna o status ativo de uma rede
func (no *NetworkOperations) ToggleNetworkActive(networkKey string) error {
	return no.config.ToggleNetworkActive(networkKey)
}

// UpdateNetworkRPC atualiza o endpoint RPC de uma rede
func (no *NetworkOperations) UpdateNetworkRPC(networkKey, rpcEndpoint string) error {
	if network, exists := no.config.GetNetworkByKey(networkKey); exists {
		network.RPCEndpoint = strings.TrimSpace(rpcEndpoint)
		if network.IsCustom {
			no.config.CustomNetworks[networkKey] = network
		} else {
			no.config.Networks[networkKey] = network
		}
		return no.config.Save()
	}
	return fmt.Errorf("network not found: %s", networkKey)
}

// AddCustomNetwork adiciona uma rede personalizada
func (no *NetworkOperations) AddCustomNetwork(name, chainID, rpcEndpoint string) (string, error) {
	// Validate input
	if err := ValidateNetworkInput(name, chainID, rpcEndpoint); err != nil {
		return "", err
	}

	// Parse chain ID
	chainIDInt, err := strconv.Atoi(strings.TrimSpace(chainID))
	if err != nil {
		return "", fmt.Errorf("invalid chain ID: %s", chainID)
	}

	// Create network key
	key := fmt.Sprintf("custom_%d", chainIDInt)

	// Create network configuration
	network := config.Network{
		Name:        strings.TrimSpace(name),
		ChainID:     int64(chainIDInt),
		RPCEndpoint: strings.TrimSpace(rpcEndpoint),
		Symbol:      "ETH", // Default symbol
		IsActive:    true,  // New networks are active by default
		IsCustom:    true,
	}

	// Add to config
	no.config.AddCustomNetwork(key, network)
	return key, no.config.Save()
}

// SearchNetworks busca redes por nome
func (no *NetworkOperations) SearchNetworks(query string) ([]blockchain.NetworkSuggestion, error) {
	if no.chainListService == nil {
		return nil, fmt.Errorf("chain list service not available")
	}
	return no.chainListService.SearchNetworksByName(query)
}

// GetChainInfoByID obtém informações de uma chain pelo ID
func (no *NetworkOperations) GetChainInfoByID(chainID int) (*blockchain.ChainInfo, string, error) {
	if no.chainListService == nil {
		return nil, "", fmt.Errorf("chain list service not available")
	}
	return no.chainListService.GetChainInfoWithRetry(chainID)
}

// ValidateNetworkInput valida os dados de entrada para adição de rede
func ValidateNetworkInput(name, chainID, rpcEndpoint string) error {
	name = strings.TrimSpace(name)
	chainID = strings.TrimSpace(chainID)
	rpcEndpoint = strings.TrimSpace(rpcEndpoint)

	if name == "" {
		return fmt.Errorf("network name is required")
	}
	if chainID == "" {
		return fmt.Errorf("chain ID is required")
	}
	if rpcEndpoint == "" {
		return fmt.Errorf("RPC endpoint is required")
	}

	// Validate chain ID is numeric
	if _, err := strconv.Atoi(chainID); err != nil {
		return fmt.Errorf("chain ID must be a number")
	}

	// Validate RPC endpoint format
	if !strings.HasPrefix(rpcEndpoint, "http://") && !strings.HasPrefix(rpcEndpoint, "https://") {
		return fmt.Errorf("RPC endpoint must start with http:// or https://")
	}

	return nil
}

// GetLanguages retorna os idiomas suportados
func (no *NetworkOperations) GetLanguages() []string {
	var languageItems []string
	langCodes := no.config.GetLanguageCodes()
	for _, code := range langCodes {
		name := config.SupportedLanguages[code]
		status := ""
		if no.config.Language == code {
			status = " (Current)"
		}
		languageItems = append(languageItems, fmt.Sprintf("%s%s", name, status))
	}
	return languageItems
}

// ChangeLanguage altera o idioma da aplicação
func (no *NetworkOperations) ChangeLanguage(languageCode string) error {
	no.config.Language = languageCode
	return no.config.Save()
}

// NetworkMessages define types para operações de rede
type NetworkAddedMsg struct {
	Key     string
	Network config.Network
}
type NetworkErrorMsg string

func SearchNetworksCmd(no *NetworkOperations, query string) tea.Cmd {
	return tea.Cmd(func() tea.Msg {
		suggestions, err := no.SearchNetworks(query)
		if err != nil {
			return NetworkErrorMsg(err.Error())
		}
		return networkSuggestionsMsg(suggestions)
	})
}

func LoadChainInfoByIDCmd(no *NetworkOperations, chainID int) tea.Cmd {
	return tea.Cmd(func() tea.Msg {
		chainInfo, rpcURL, err := no.GetChainInfoByID(chainID)
		if err != nil {
			return NetworkErrorMsg(err.Error())
		}
		return chainInfoLoadedMsg{
			chainInfo: chainInfo,
			rpcURL:    rpcURL,
		}
	})
}

func AddNetworkCmd(no *NetworkOperations, name, chainID, rpcEndpoint string) tea.Cmd {
	return tea.Cmd(func() tea.Msg {
		key, err := no.AddCustomNetwork(name, chainID, rpcEndpoint)
		if err != nil {
			return NetworkErrorMsg(err.Error())
		}

		// Get the added network for the message
		if network, exists := no.config.GetNetworkByKey(key); exists {
			return NetworkAddedMsg{
				Key:     key,
				Network: network,
			}
		}

		return NetworkErrorMsg("failed to retrieve added network")
	})
}
