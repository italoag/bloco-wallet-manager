package main

import (
	"fmt"
	"os"
	"path/filepath"

	"blocowallet/config"
	"blocowallet/infrastructure"
	"blocowallet/interfaces"
	"blocowallet/localization"
	"blocowallet/usecases"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/ethereum/go-ethereum/accounts/keystore"
)

func main() {
	// Get the user's home directory
	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Println("Error getting the user's home directory:", err)
		os.Exit(1)
	}

	// Define the base directory for the application
	appDir := filepath.Join(homeDir, ".wallets")

	// Ensure the base directory exists
	if _, err := os.Stat(appDir); os.IsNotExist(err) {
		err := os.MkdirAll(appDir, os.ModePerm)
		if err != nil {
			fmt.Println("Error creating application directory:", err)
			os.Exit(1)
		}
	}

	// Initialize configuration
	cfg, err := config.LoadConfig(appDir)
	if err != nil {
		fmt.Println("Error loading configuration:", err)
		os.Exit(1)
	}

	// Initialize localization
	err = localization.SetLanguage(cfg.Language, appDir)
	if err != nil {
		fmt.Println("Error loading localization files:", err)
		os.Exit(1)
	}

	// Ensure the wallets directory exists
	if _, err := os.Stat(cfg.WalletsDir); os.IsNotExist(err) {
		err := os.MkdirAll(cfg.WalletsDir, os.ModePerm)
		if err != nil {
			fmt.Println("Error creating wallets directory:", err)
			os.Exit(1)
		}
	}

	// Initialize the repository
	repo, err := infrastructure.NewSQLiteRepository(cfg.DatabasePath)
	if err != nil {
		fmt.Println("Error initializing the database:", err)
		os.Exit(1)
	}
	defer func() {
		if err := repo.Close(); err != nil {
			fmt.Println("Error closing the database:", err)
		}
	}()

	// Initialize the keystore
	ks := keystore.NewKeyStore(cfg.WalletsDir, keystore.StandardScryptN, keystore.StandardScryptP)
	service := usecases.NewWalletService(repo, ks)
	model := interfaces.NewCLIModel(service)

	// Start the Bubble Tea program
	p := tea.NewProgram(&model, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error running the program: %v\n", err)
		os.Exit(1)
	}
}
