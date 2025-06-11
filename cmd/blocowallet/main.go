package main

import (
	"log"
	"os"
	"path/filepath"

	"blocowallet/internal/blockchain"
	"blocowallet/internal/storage"
	"blocowallet/internal/ui"
	"blocowallet/internal/wallet"
	"blocowallet/pkg/config"
	"blocowallet/pkg/logger"

	tea "github.com/charmbracelet/bubbletea"
	zone "github.com/lrstanley/bubblezone"
)

func main() {
	// Initialize global zone manager for mouse interactions
	zone.NewGlobal()

	// Initialize logger first - use error level to avoid interfering with TUI
	appLogger, err := logger.NewLogger("error")
	if err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}
	defer appLogger.Sync()

	appLogger.Debug("Starting BlockoWallet application",
		logger.String("operation", "application_startup"))

	// Setup application directory
	homeDir, err := os.UserHomeDir()
	if err != nil {
		appLogger.Error("Failed to get user home directory",
			logger.Error(err),
			logger.String("operation", "application_startup"))
		log.Fatalf("Failed to get user home directory: %v", err)
	}

	appDir := filepath.Join(homeDir, ".blocowallet")
	if err := os.MkdirAll(appDir, 0755); err != nil {
		appLogger.Error("Failed to create app directory",
			logger.Error(err),
			logger.String("app_dir", appDir),
			logger.String("operation", "application_startup"))
		log.Fatalf("Failed to create app directory: %v", err)
	}

	appLogger.Debug("Application directory created",
		logger.String("app_dir", appDir),
		logger.String("operation", "application_startup"))

	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		appLogger.Error("Failed to load configuration",
			logger.Error(err),
			logger.String("operation", "application_startup"))
		log.Fatalf("Failed to load config: %v", err)
	}

	appLogger.Debug("Configuration loaded successfully",
		logger.String("operation", "application_startup"))

	// Initialize storage
	dbPath := cfg.Database.Path
	if dbPath == "" || dbPath == "wallets.db" || !filepath.IsAbs(dbPath) {
		dbPath = filepath.Join(appDir, "wallets.db")
	}

	// Ensure the directory for the database exists
	if err := os.MkdirAll(filepath.Dir(dbPath), 0755); err != nil {
		appLogger.Error("Failed to create database directory",
			logger.Error(err),
			logger.String("db_path", dbPath),
			logger.String("operation", "application_startup"))
		log.Fatalf("Failed to create database directory: %v", err)
	}

	repo, err := storage.NewSQLite(dbPath, appLogger)
	if err != nil {
		appLogger.Error("Failed to initialize storage",
			logger.Error(err),
			logger.String("db_path", dbPath),
			logger.String("operation", "application_startup"))
		log.Fatalf("Failed to initialize storage: %v", err)
	}
	defer repo.Close()

	// Initialize multi provider for all active networks
	multiProvider := blockchain.NewMultiProvider()
	defer multiProvider.Close()

	// Setup providers for all networks
	multiProvider.RefreshProviders(cfg)

	appLogger.Debug("Blockchain providers initialized",
		logger.String("operation", "application_startup"))

	// Initialize wallet service with multi-provider
	walletService := wallet.NewServiceWithMultiProvider(repo, multiProvider, appLogger)

	appLogger.Debug("Wallet service initialized",
		logger.String("operation", "application_startup"))

	// Initialize and run TUI
	model := ui.NewModel(walletService, cfg)
	p := tea.NewProgram(model, tea.WithAltScreen(), tea.WithMouseCellMotion())

	appLogger.Debug("Starting TUI interface",
		logger.String("operation", "application_startup"))

	if _, err := p.Run(); err != nil {
		appLogger.Error("TUI execution failed",
			logger.Error(err),
			logger.String("operation", "application_startup"))
		log.Fatalf("Failed to run TUI: %v", err)
	}

	appLogger.Debug("Application shutdown completed",
		logger.String("operation", "application_shutdown"))
}
