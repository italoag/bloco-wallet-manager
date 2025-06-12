package config

import (
	"gopkg.in/yaml.v2"
	"os"
	"path/filepath"
)

type Config struct {
	AppDir       string             `yaml:"app_dir"`
	Language     string             `yaml:"language"`
	WalletsDir   string             `yaml:"wallets_dir"`
	Database     DatabaseConfig     `yaml:"database"`
	Networks     map[string]Network `yaml:"networks"`
	DatabasePath string             `yaml:"database_path"`
}

type Network struct {
	Name        string `yaml:"name"`
	RPCEndpoint string `yaml:"rpc_endpoint"`
	ChainID     int64  `yaml:"chain_id"`
	Symbol      string `yaml:"symbol"`
	Explorer    string `yaml:"explorer"`
	IsActive    bool   `yaml:"is_active"`
	IsCustom    bool   `yaml:"is_custom"`
}

type DatabaseConfig struct {
	Type string `yaml:"type"` // sqlite, postgres
	Path string `yaml:"path"` // for sqlite
	URL  string `yaml:"url"`  // for postgres
}

func LoadConfig(appDir string) (*Config, error) {
	configPath := filepath.Join(appDir, "config.yaml")

	// If a config file doesn't exist, create it with default values
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		defaultConfig := &Config{
			AppDir:       appDir,
			Language:     "en",
			WalletsDir:   filepath.Join(appDir, "keystore"),
			DatabasePath: filepath.Join(appDir, "wallets.db"),
		}

		configData, err := yaml.Marshal(defaultConfig)
		if err != nil {
			return nil, err
		}

		err = os.WriteFile(configPath, configData, 0644)
		if err != nil {
			return nil, err
		}
	}

	// Load the configuration file
	configFile, err := os.Open(configPath)
	if err != nil {
		return nil, err
	}
	defer func(configFile *os.File) {
		err := configFile.Close()
		if err != nil {
			panic(err)
		}
	}(configFile)

	cfg := &Config{}
	decoder := yaml.NewDecoder(configFile)
	err = decoder.Decode(cfg)
	if err != nil {
		return nil, err
	}

	// Expand ~ to the home directory
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	if cfg.WalletsDir != "" {
		cfg.WalletsDir = expandPath(cfg.WalletsDir, homeDir)
	} else {
		cfg.WalletsDir = filepath.Join(appDir, "keystore")
	}

	if cfg.DatabasePath != "" {
		cfg.DatabasePath = expandPath(cfg.DatabasePath, homeDir)
	} else {
		cfg.DatabasePath = filepath.Join(appDir, "wallets.db")
	}

	return cfg, nil
}

func expandPath(path, homeDir string) string {
	if len(path) > 1 && path[:2] == "~/" {
		return filepath.Join(homeDir, path[2:])
	}
	return path
}
