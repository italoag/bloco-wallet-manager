package storage

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"blocowallet/internal/wallet"
	"blocowallet/pkg/logger"

	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
)

// SQLite implements wallet.Repository using SQLite
type SQLite struct {
	db     *sql.DB
	logger logger.Logger
}

// NewSQLite creates a new SQLite repository
func NewSQLite(databasePath string, log logger.Logger) (*SQLite, error) {
	db, err := sql.Open("sqlite3", databasePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	sqlite := &SQLite{
		db:     db,
		logger: log,
	}

	log.Debug("Initializing SQLite database",
		logger.String("database_path", databasePath),
		logger.String("operation", "initialize_database"))

	if err := sqlite.createTables(); err != nil {
		log.Error("Failed to create database tables",
			logger.Error(err),
			logger.String("database_path", databasePath),
			logger.String("operation", "initialize_database"))
		return nil, fmt.Errorf("failed to create tables: %w", err)
	}

	log.Debug("SQLite database initialized successfully",
		logger.String("database_path", databasePath),
		logger.String("operation", "initialize_database"))

	return sqlite, nil
}

// createTables creates the necessary database tables
func (s *SQLite) createTables() error {
	query := `
	CREATE TABLE IF NOT EXISTS wallets (
		id TEXT PRIMARY KEY,
		name TEXT NOT NULL,
		address TEXT UNIQUE NOT NULL,
		keystore_path TEXT NOT NULL,
		encrypted_mnemonic TEXT,
		created_at DATETIME NOT NULL,
		updated_at DATETIME NOT NULL
	);
	`

	_, err := s.db.Exec(query)
	if err != nil {
		return fmt.Errorf("failed to create wallets table: %w", err)
	}

	return nil
}

// Create creates a new wallet in the database
func (s *SQLite) Create(ctx context.Context, w *wallet.Wallet) error {
	correlationID, _ := ctx.Value("correlation_id").(string)
	if correlationID == "" {
		correlationID = uuid.New().String()
	}

	s.logger.Debug("Creating wallet in database",
		logger.String("correlation_id", correlationID),
		logger.String("wallet_id", w.ID),
		logger.String("address", w.Address),
		logger.String("operation", "db_create_wallet"))

	query := `
	INSERT INTO wallets (id, name, address, keystore_path, encrypted_mnemonic, created_at, updated_at)
	VALUES (?, ?, ?, ?, ?, ?, ?)
	`

	_, err := s.db.ExecContext(ctx, query,
		w.ID, w.Name, w.Address, w.KeyStorePath, w.EncryptedMnemonic, w.CreatedAt, w.UpdatedAt)
	if err != nil {
		s.logger.Error("Database wallet creation failed",
			logger.String("correlation_id", correlationID),
			logger.String("wallet_id", w.ID),
			logger.Error(err),
			logger.String("operation", "db_create_wallet"))
		return fmt.Errorf("failed to create wallet: %w", err)
	}

	s.logger.Debug("Wallet saved to database successfully",
		logger.String("correlation_id", correlationID),
		logger.String("wallet_id", w.ID),
		logger.String("operation", "db_create_wallet"))

	return nil
}

// GetByID retrieves a wallet by ID
func (s *SQLite) GetByID(ctx context.Context, id string) (*wallet.Wallet, error) {
	correlationID, _ := ctx.Value("correlation_id").(string)
	if correlationID == "" {
		correlationID = uuid.New().String()
	}

	s.logger.Debug("Retrieving wallet by ID from database",
		logger.String("correlation_id", correlationID),
		logger.String("wallet_id", id),
		logger.String("operation", "db_get_wallet_by_id"))

	query := `
	SELECT id, name, address, keystore_path, encrypted_mnemonic, created_at, updated_at
	FROM wallets WHERE id = ?
	`

	row := s.db.QueryRowContext(ctx, query, id)

	var w wallet.Wallet
	var createdAt, updatedAt string

	err := row.Scan(&w.ID, &w.Name, &w.Address, &w.KeyStorePath, &w.EncryptedMnemonic, &createdAt, &updatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			s.logger.Debug("Wallet not found in database",
				logger.String("correlation_id", correlationID),
				logger.String("wallet_id", id),
				logger.String("operation", "db_get_wallet_by_id"))
			return nil, fmt.Errorf("wallet not found")
		}
		s.logger.Error("Database scan failed",
			logger.String("correlation_id", correlationID),
			logger.String("wallet_id", id),
			logger.Error(err),
			logger.String("operation", "db_get_wallet_by_id"))
		return nil, fmt.Errorf("failed to scan wallet: %w", err)
	}

	if w.CreatedAt, err = time.Parse(time.RFC3339, createdAt); err != nil {
		s.logger.Error("Failed to parse created_at timestamp",
			logger.String("correlation_id", correlationID),
			logger.String("wallet_id", id),
			logger.Error(err),
			logger.String("operation", "db_get_wallet_by_id"))
		return nil, fmt.Errorf("failed to parse created_at: %w", err)
	}

	if w.UpdatedAt, err = time.Parse(time.RFC3339, updatedAt); err != nil {
		s.logger.Error("Failed to parse updated_at timestamp",
			logger.String("correlation_id", correlationID),
			logger.String("wallet_id", id),
			logger.Error(err),
			logger.String("operation", "db_get_wallet_by_id"))
		return nil, fmt.Errorf("failed to parse updated_at: %w", err)
	}

	s.logger.Debug("Wallet retrieved from database successfully",
		logger.String("correlation_id", correlationID),
		logger.String("wallet_id", id),
		logger.String("operation", "db_get_wallet_by_id"))

	return &w, nil
}

// GetByAddress retrieves a wallet by address
func (s *SQLite) GetByAddress(ctx context.Context, address string) (*wallet.Wallet, error) {
	query := `
	SELECT id, name, address, keystore_path, encrypted_mnemonic, created_at, updated_at
	FROM wallets WHERE address = ?
	`

	row := s.db.QueryRowContext(ctx, query, address)

	var w wallet.Wallet
	var createdAt, updatedAt string

	err := row.Scan(&w.ID, &w.Name, &w.Address, &w.KeyStorePath, &w.EncryptedMnemonic, &createdAt, &updatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("wallet not found")
		}
		return nil, fmt.Errorf("failed to scan wallet: %w", err)
	}

	if w.CreatedAt, err = time.Parse(time.RFC3339, createdAt); err != nil {
		return nil, fmt.Errorf("failed to parse created_at: %w", err)
	}

	if w.UpdatedAt, err = time.Parse(time.RFC3339, updatedAt); err != nil {
		return nil, fmt.Errorf("failed to parse updated_at: %w", err)
	}

	return &w, nil
}

// List retrieves all wallets
func (s *SQLite) List(ctx context.Context) ([]*wallet.Wallet, error) {
	query := `
	SELECT id, name, address, keystore_path, encrypted_mnemonic, created_at, updated_at
	FROM wallets ORDER BY created_at DESC
	`

	rows, err := s.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to query wallets: %w", err)
	}
	defer rows.Close()

	var wallets []*wallet.Wallet

	for rows.Next() {
		var w wallet.Wallet
		var createdAt, updatedAt string

		err := rows.Scan(&w.ID, &w.Name, &w.Address, &w.KeyStorePath, &w.EncryptedMnemonic, &createdAt, &updatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan wallet: %w", err)
		}

		if w.CreatedAt, err = time.Parse(time.RFC3339, createdAt); err != nil {
			return nil, fmt.Errorf("failed to parse created_at: %w", err)
		}

		if w.UpdatedAt, err = time.Parse(time.RFC3339, updatedAt); err != nil {
			return nil, fmt.Errorf("failed to parse updated_at: %w", err)
		}

		wallets = append(wallets, &w)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate over rows: %w", err)
	}

	return wallets, nil
}

// Update updates a wallet
func (s *SQLite) Update(ctx context.Context, w *wallet.Wallet) error {
	query := `
	UPDATE wallets 
	SET name = ?, address = ?, keystore_path = ?, encrypted_mnemonic = ?, updated_at = ?
	WHERE id = ?
	`

	_, err := s.db.ExecContext(ctx, query,
		w.Name, w.Address, w.KeyStorePath, w.EncryptedMnemonic, w.UpdatedAt, w.ID)
	if err != nil {
		return fmt.Errorf("failed to update wallet: %w", err)
	}

	return nil
}

// Delete deletes a wallet
func (s *SQLite) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM wallets WHERE id = ?`

	_, err := s.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete wallet: %w", err)
	}

	return nil
}

// Close closes the database connection
func (s *SQLite) Close() error {
	return s.db.Close()
}
