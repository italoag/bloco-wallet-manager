package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
)

type Config struct {
	Language       string             `json:"language"`
	Networks       map[string]Network `json:"networks"`
	CustomNetworks map[string]Network `json:"custom_networks"`
	Database       DatabaseConfig     `json:"database"`
	UIConfig       UIConfig           `json:"ui"`
}

type Network struct {
	Name        string `json:"name"`
	RPCEndpoint string `json:"rpc_endpoint"`
	ChainID     int64  `json:"chain_id"`
	Symbol      string `json:"symbol"`
	Explorer    string `json:"explorer"`
	IsActive    bool   `json:"is_active"`
	IsCustom    bool   `json:"is_custom"`
}

type DatabaseConfig struct {
	Type string `json:"type"` // sqlite, postgres
	Path string `json:"path"` // for sqlite
	URL  string `json:"url"`  // for postgres
}

type UIConfig struct {
	Theme      string `json:"theme"`
	ShowSplash bool   `json:"show_splash"`
}

// DefaultNetworks is now empty - users must add their own networks
var DefaultNetworks = map[string]Network{}

var SupportedLanguages = map[string]string{
	"en": "English",
	"pt": "Português",
	"es": "Español",
	"fr": "Français",
}

func DefaultConfig() *Config {
	homeDir, _ := os.UserHomeDir()
	return &Config{
		Language:       "en",
		Networks:       DefaultNetworks,
		CustomNetworks: make(map[string]Network),
		Database: DatabaseConfig{
			Type: "sqlite",
			Path: filepath.Join(homeDir, ".blocowallet", "wallets.db"),
		},
		UIConfig: UIConfig{
			Theme:      "default",
			ShowSplash: true,
		},
	}
}

func LoadConfig() (*Config, error) {
	configPath := getConfigPath()

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		cfg := DefaultConfig()
		if err := cfg.Save(); err != nil {
			return nil, fmt.Errorf("failed to create default config: %w", err)
		}
		return cfg, nil
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	// Merge with default networks if missing
	if cfg.Networks == nil {
		cfg.Networks = DefaultNetworks
	}

	// Initialize custom networks if missing
	if cfg.CustomNetworks == nil {
		cfg.CustomNetworks = make(map[string]Network)
	}

	return &cfg, nil
}

func (c *Config) Save() error {
	configPath := getConfigPath()

	// Criar diretório se não existir
	if err := os.MkdirAll(filepath.Dir(configPath), 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	if err := os.WriteFile(configPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}

func getConfigPath() string {
	homeDir, _ := os.UserHomeDir()
	return filepath.Join(homeDir, ".blocowallet", "config.json")
}

func (c *Config) GetActiveNetwork() *Network {
	// Check regular networks
	for _, network := range c.Networks {
		if network.IsActive {
			return &network
		}
	}
	// Check custom networks
	for _, network := range c.CustomNetworks {
		if network.IsActive {
			return &network
		}
	}
	return nil
}

func (c *Config) SetActiveNetwork(networkKey string) error {
	// Deactivate all networks first
	for k, v := range c.Networks {
		v.IsActive = false
		c.Networks[k] = v
	}
	for k, v := range c.CustomNetworks {
		v.IsActive = false
		c.CustomNetworks[k] = v
	}

	// Activate the selected network
	if network, exists := c.Networks[networkKey]; exists {
		network.IsActive = true
		c.Networks[networkKey] = network
		return nil
	}

	if network, exists := c.CustomNetworks[networkKey]; exists {
		network.IsActive = true
		c.CustomNetworks[networkKey] = network
		return nil
	}

	return fmt.Errorf("network %s not found", networkKey)
}

func (c *Config) UpdateNetworkRPC(networkKey, rpcEndpoint string) error {
	if network, exists := c.Networks[networkKey]; exists {
		network.RPCEndpoint = rpcEndpoint
		c.Networks[networkKey] = network
		return nil
	}
	return fmt.Errorf("network %s not found", networkKey)
}

func (c *Config) GetNetworkKeys() []string {
	var keys []string
	for key := range c.Networks {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return keys
}

func (c *Config) GetLanguageCodes() []string {
	var codes []string
	for code := range SupportedLanguages {
		codes = append(codes, code)
	}
	sort.Strings(codes)
	return codes
}

// Legacy support methods
func (c *Config) GetDatabasePath() string {
	return c.Database.Path
}

func (c *Config) GetRPCEndpoint() string {
	if activeNetwork := c.GetActiveNetwork(); activeNetwork != nil {
		return activeNetwork.RPCEndpoint
	}
	return "" // No default RPC - user must add networks
}

func (c *Config) AddCustomNetwork(key string, network Network) {
	if c.CustomNetworks == nil {
		c.CustomNetworks = make(map[string]Network)
	}
	network.IsCustom = true
	c.CustomNetworks[key] = network
}

func (c *Config) RemoveCustomNetwork(key string) {
	delete(c.CustomNetworks, key)
}

func (c *Config) GetCustomNetwork(key string) (Network, bool) {
	network, exists := c.CustomNetworks[key]
	return network, exists
}

func (c *Config) GetAllNetworks() map[string]Network {
	allNetworks := make(map[string]Network)
	for k, v := range c.Networks {
		allNetworks[k] = v
	}
	for k, v := range c.CustomNetworks {
		allNetworks[k] = v
	}
	return allNetworks
}

func (c *Config) GetAllNetworkKeys() []string {
	var keys []string

	// Add default networks
	for key := range c.Networks {
		keys = append(keys, key)
	}

	// Add custom networks
	for key := range c.CustomNetworks {
		keys = append(keys, key)
	}

	sort.Strings(keys)
	return keys
}

func (c *Config) GetNetworkByKey(key string) (Network, bool) {
	// Check default networks first
	if network, exists := c.Networks[key]; exists {
		return network, true
	}

	// Check custom networks
	if network, exists := c.CustomNetworks[key]; exists {
		return network, true
	}

	return Network{}, false
}

func (c *Config) UpdateNetwork(key string, network Network) error {
	// Check if it's a default network
	if _, exists := c.Networks[key]; exists {
		c.Networks[key] = network
		return nil
	}

	// Check if it's a custom network
	if _, exists := c.CustomNetworks[key]; exists {
		network.IsCustom = true
		c.CustomNetworks[key] = network
		return nil
	}

	return fmt.Errorf("network with key %s not found", key)
}

// ToggleNetworkActive toggles the active status of a network without affecting others
func (c *Config) ToggleNetworkActive(networkKey string) error {
	// Check in regular networks
	if network, exists := c.Networks[networkKey]; exists {
		network.IsActive = !network.IsActive
		c.Networks[networkKey] = network
		return nil
	}

	// Check in custom networks
	if network, exists := c.CustomNetworks[networkKey]; exists {
		network.IsActive = !network.IsActive
		c.CustomNetworks[networkKey] = network
		return nil
	}

	return fmt.Errorf("network %s not found", networkKey)
}

// GetActiveNetworks returns all active networks
func (c *Config) GetActiveNetworks() map[string]Network {
	activeNetworks := make(map[string]Network)

	// Add active regular networks
	for key, network := range c.Networks {
		if network.IsActive {
			activeNetworks[key] = network
		}
	}

	// Add active custom networks
	for key, network := range c.CustomNetworks {
		if network.IsActive {
			activeNetworks[key] = network
		}
	}

	return activeNetworks
}
