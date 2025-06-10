package config

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestDefaultConfig(t *testing.T) {
	config := DefaultConfig()

	// Verify default values
	if config.Language != "en" {
		t.Errorf("Default language should be 'en', got: %s", config.Language)
	}

	if len(config.Networks) != 0 {
		t.Fatalf("Expected no default networks, got %d", len(config.Networks))
	}

	if config.Database.Type != "sqlite" {
		t.Errorf("Default database type should be sqlite, got: %s", config.Database.Type)
	}

	if config.UIConfig.Theme == "" {
		t.Fatal("Default config should have a UI theme")
	}

	if config.CustomNetworks == nil {
		t.Fatal("Custom networks should be initialized")
	}
}

func TestConfigSaveAndLoad(t *testing.T) {
	// Create temporary directory for test
	tempDir, err := os.MkdirTemp("", "blocowallet_config_test")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Override the config path for this test
	originalConfigDir := os.Getenv("XDG_CONFIG_HOME")
	defer func() {
		if originalConfigDir != "" {
			os.Setenv("XDG_CONFIG_HOME", originalConfigDir)
		} else {
			os.Unsetenv("XDG_CONFIG_HOME")
		}
	}()
	os.Setenv("XDG_CONFIG_HOME", tempDir)

	// Create test config
	testConfig := &Config{
		Language: "pt",
		Networks: map[string]Network{
			"ethereum": {
				Name:        "Ethereum Mainnet",
				RPCEndpoint: "https://test-rpc.com",
				ChainID:     1,
				Symbol:      "ETH",
				Explorer:    "https://etherscan.io",
				IsActive:    true,
				IsCustom:    false,
			},
		},
		CustomNetworks: map[string]Network{
			"custom": {
				Name:        "Custom Network",
				RPCEndpoint: "https://custom-rpc.com",
				ChainID:     999,
				Symbol:      "CUSTOM",
				Explorer:    "https://custom-explorer.com",
				IsActive:    true,
				IsCustom:    true,
			},
		},
		Database: DatabaseConfig{
			Type: "sqlite",
			Path: filepath.Join(tempDir, "test.db"),
		},
		UIConfig: UIConfig{
			Theme:      "dark",
			ShowSplash: false,
		},
	}

	// Test saving config
	err = testConfig.Save()
	if err != nil {
		t.Fatalf("Failed to save config: %v", err)
	}

	// Test loading config
	loadedConfig, err := LoadConfig()
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	// Verify loaded config
	if loadedConfig.Language != testConfig.Language {
		t.Errorf("Language mismatch. Got: %s, Expected: %s", loadedConfig.Language, testConfig.Language)
	}

	if loadedConfig.UIConfig.Theme != testConfig.UIConfig.Theme {
		t.Errorf("UI Theme mismatch. Got: %s, Expected: %s", loadedConfig.UIConfig.Theme, testConfig.UIConfig.Theme)
	}

	if loadedConfig.UIConfig.ShowSplash != testConfig.UIConfig.ShowSplash {
		t.Errorf("UI ShowSplash mismatch. Got: %t, Expected: %t", loadedConfig.UIConfig.ShowSplash, testConfig.UIConfig.ShowSplash)
	}

	if loadedConfig.Database.Path != testConfig.Database.Path {
		t.Errorf("Database path mismatch. Got: %s, Expected: %s", loadedConfig.Database.Path, testConfig.Database.Path)
	}

	// Verify networks
	if len(loadedConfig.Networks) != len(testConfig.Networks) {
		t.Errorf("Networks count mismatch. Got: %d, Expected: %d", len(loadedConfig.Networks), len(testConfig.Networks))
	}

	ethereumNetwork, exists := loadedConfig.Networks["ethereum"]
	if !exists {
		t.Fatal("Ethereum network should exist")
	}

	if ethereumNetwork.RPCEndpoint != "https://test-rpc.com" {
		t.Errorf("Ethereum RPC endpoint mismatch. Got: %s, Expected: https://test-rpc.com", ethereumNetwork.RPCEndpoint)
	}

	// Verify custom networks
	customNetwork, exists := loadedConfig.CustomNetworks["custom"]
	if !exists {
		t.Fatal("Custom network should exist")
	}

	if customNetwork.Name != "Custom Network" {
		t.Errorf("Custom network name mismatch. Got: %s, Expected: Custom Network", customNetwork.Name)
	}
}

func TestNetworkStructure(t *testing.T) {
	testCases := []struct {
		name    string
		network Network
		valid   bool
	}{
		{
			name: "Valid network",
			network: Network{
				Name:        "Test Network",
				RPCEndpoint: "https://test-rpc.com",
				ChainID:     1,
				Symbol:      "TEST",
				Explorer:    "https://test-explorer.com",
				IsActive:    true,
				IsCustom:    false,
			},
			valid: true,
		},
		{
			name: "Network with all fields",
			network: Network{
				Name:        "Custom Network",
				RPCEndpoint: "https://custom-rpc.com",
				ChainID:     999,
				Symbol:      "CUSTOM",
				Explorer:    "https://custom-explorer.com",
				IsActive:    false,
				IsCustom:    true,
			},
			valid: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Just verify the network structure is correct
			if tc.network.Name == "" && tc.valid {
				t.Error("Valid network should have a name")
			}
			if tc.network.RPCEndpoint == "" && tc.valid {
				t.Error("Valid network should have an RPC endpoint")
			}
			if tc.network.ChainID <= 0 && tc.valid {
				t.Error("Valid network should have a positive chain ID")
			}
		})
	}
}

func TestJSONMarshaling(t *testing.T) {
	config := DefaultConfig()
	config.Language = "pt"

	// Test marshaling
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		t.Fatalf("Failed to marshal config: %v", err)
	}

	// Test unmarshaling
	var unmarshaledConfig Config
	err = json.Unmarshal(data, &unmarshaledConfig)
	if err != nil {
		t.Fatalf("Failed to unmarshal config: %v", err)
	}

	// Verify
	if unmarshaledConfig.Language != config.Language {
		t.Errorf("Language mismatch after JSON roundtrip. Got: %s, Expected: %s", unmarshaledConfig.Language, config.Language)
	}

	if len(unmarshaledConfig.Networks) != len(config.Networks) {
		t.Errorf("Networks count mismatch after JSON roundtrip. Got: %d, Expected: %d", len(unmarshaledConfig.Networks), len(config.Networks))
	}
}
