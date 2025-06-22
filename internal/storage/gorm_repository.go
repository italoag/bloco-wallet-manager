package storage

import (
	"blocowallet/internal/wallet"
	"blocowallet/pkg/config"
	"fmt"
	"os"
	"path/filepath"

	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// GORMRepository implements the Repository interface using GORM
type GORMRepository struct {
	db *gorm.DB
}

// Ensure that GORMRepository implements the Repository interface
var _ wallet.Repository = &GORMRepository{}

// NewWalletRepository creates a new instance of GORMRepository based on the configuration
func NewWalletRepository(cfg *config.Config) (*GORMRepository, error) {
	var dialector gorm.Dialector

	switch cfg.Database.Type {
	case "postgres":
		dsn := cfg.Database.DSN
		if dsn == "" {
			return nil, fmt.Errorf("empty DSN configuration for PostgreSQL")
		}
		dialector = postgres.Open(dsn)
	case "mysql":
		dsn := cfg.Database.DSN
		if dsn == "" {
			return nil, fmt.Errorf("empty DSN configuration for MySQL")
		}
		dialector = mysql.Open(dsn)
	case "sqlite", "":
		// Use SQLite by default
		dbPath := cfg.DatabasePath
		if cfg.Database.DSN != "" {
			dbPath = cfg.Database.DSN
		}

		// Ensure the directory exists
		dir := filepath.Dir(dbPath)
		if err := ensureDir(dir); err != nil {
			return nil, fmt.Errorf("failed to create directory for database: %w", err)
		}

		dialector = sqlite.Open(dbPath)
	default:
		return nil, fmt.Errorf("unsupported database type: %s", cfg.Database.Type)
	}

	db, err := gorm.Open(dialector, &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Auto Migrate creates the table if it doesn't exist
	err = db.AutoMigrate(&wallet.Wallet{})
	if err != nil {
		return nil, fmt.Errorf("failed to migrate wallets table: %w", err)
	}

	return &GORMRepository{db: db}, nil
}

// ensureDir ensures that the directory exists
func ensureDir(dir string) error {
	return os.MkdirAll(dir, os.ModePerm)
}

// AddWallet adds a new wallet to the database
func (repo *GORMRepository) AddWallet(wallet *wallet.Wallet) error {
	return repo.db.Create(wallet).Error
}

// GetAllWallets returns all saved wallets
func (repo *GORMRepository) GetAllWallets() ([]wallet.Wallet, error) {
	var wallets []wallet.Wallet
	result := repo.db.Find(&wallets)
	return wallets, result.Error
}

// DeleteWallet removes a wallet by ID
func (repo *GORMRepository) DeleteWallet(walletID int) error {
	return repo.db.Delete(&wallet.Wallet{}, walletID).Error
}

// Close closes the database connection
func (repo *GORMRepository) Close() error {
	sqlDB, err := repo.db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}
