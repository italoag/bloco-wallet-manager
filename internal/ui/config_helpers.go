package ui

import (
	"blocowallet/pkg/config"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

// loadOrCreateConfig loads the configuration from the specified directory
// or creates a new one if it doesn't exist
func loadOrCreateConfig(appDir string) (*config.Config, error) {
	// Try to load the existing configuration
	cfg, err := config.LoadConfig(appDir)
	if err != nil {
		// If the configuration doesn't exist, create a new one
		if os.IsNotExist(err) {
			// Create the directory if it doesn't exist
			if err := os.MkdirAll(appDir, os.ModePerm); err != nil {
				return nil, fmt.Errorf("failed to create config directory: %w", err)
			}

			// Create a new configuration
			cfg = &config.Config{
				AppDir:     appDir,
				Language:   "en",
				WalletsDir: filepath.Join(appDir, "wallets"),
				Networks:   make(map[string]config.Network),
			}

			// Save the configuration
			configPath := filepath.Join(appDir, "config.toml")
			v := viper.New()
			v.SetConfigFile(configPath)

			// Set the configuration values
			v.Set("app.language", cfg.Language)
			v.Set("app.app_dir", cfg.AppDir)
			v.Set("app.wallets_dir", cfg.WalletsDir)

			// Write the configuration to file
			if err := v.WriteConfigAs(configPath); err != nil {
				return nil, fmt.Errorf("failed to write config file: %w", err)
			}

			return cfg, nil
		}

		return nil, fmt.Errorf("failed to load configuration: %w", err)
	}

	return cfg, nil
}

// updateLanguageInConfig updates the language in the configuration file
func updateLanguageInConfig(configPath string, language string) error {
	// Create a new Viper instance
	v := viper.New()
	v.SetConfigFile(configPath)

	// Try to read the existing config
	err := v.ReadInConfig()
	if err != nil {
		return fmt.Errorf("failed to read config file: %w", err)
	}

	// Update the language
	v.Set("app.language", language)

	// Write the updated config back to the file
	if err := v.WriteConfig(); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}
