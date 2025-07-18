package wallet

import (
	"blocowallet/pkg/config"
)

// CreateMockConfig creates a mock configuration for testing
func CreateMockConfig() *config.Config {
	return &config.Config{
		AppDir:       "/tmp/blocowallet-test",
		Language:     "en",
		WalletsDir:   "/tmp/blocowallet-test/keystore",
		DatabasePath: "/tmp/blocowallet-test/wallets.db",
		LocaleDir:    "/tmp/blocowallet-test/locale",
		Fonts:        []string{"test-font"},
		Database: config.DatabaseConfig{
			Type: "sqlite",
			DSN:  ":memory:",
		},
		Security: config.SecurityConfig{
			Argon2Time:    1,
			Argon2Memory:  64 * 1024, // 64MB
			Argon2Threads: 4,
			Argon2KeyLen:  32,
			SaltLength:    16,
		},
		Networks: map[string]config.Network{
			"ethereum": {
				Name:        "Ethereum",
				RPCEndpoint: "https://mainnet.infura.io/v3/your-api-key",
				ChainID:     1,
				Symbol:      "ETH",
				Explorer:    "https://etherscan.io",
				IsActive:    true,
			},
		},
	}
}
