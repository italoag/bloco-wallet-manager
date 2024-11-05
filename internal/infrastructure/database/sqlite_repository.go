package database

import (
	domain2 "blocowallet/internal/domain/entities"
	"blocowallet/internal/domain/gateway"
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

type SQLiteRepository struct {
	conn *sql.DB
}

// Implement the WalletRepository interface from entities package
var _ gateway.WalletRepository = &SQLiteRepository{}

func NewSQLiteRepository(dbPath string) (*SQLiteRepository, error) {
	conn, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}

	createTableQuery := `
	CREATE TABLE IF NOT EXISTS wallets (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		address TEXT UNIQUE NOT NULL,
		keystore_path TEXT NOT NULL,
		mnemonic TEXT NOT NULL
	);
	`
	_, err = conn.Exec(createTableQuery)
	if err != nil {
		return nil, err
	}

	return &SQLiteRepository{conn: conn}, nil
}

func (repo *SQLiteRepository) CreateWallet(wallet *domain2.Wallet) error {
	insertQuery := `
	INSERT INTO wallets (address, keystore_path, mnemonic)
	VALUES (?, ?, ?);
	`
	_, err := repo.conn.Exec(insertQuery, wallet.Address, wallet.KeyStorePath, wallet.Mnemonic)
	return err
}

func (repo *SQLiteRepository) ListWallets() ([]domain2.Wallet, error) {
	selectQuery := `
	SELECT id, address, keystore_path, mnemonic FROM wallets;
	`
	rows, err := repo.conn.Query(selectQuery)
	if err != nil {
		return nil, err
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			panic(err)
		}
	}(rows)

	var wallets []domain2.Wallet
	for rows.Next() {
		var w domain2.Wallet
		err := rows.Scan(&w.ID, &w.Address, &w.KeyStorePath, &w.Mnemonic)
		if err != nil {
			return nil, err
		}
		wallets = append(wallets, w)
	}

	return wallets, nil
}

func (repo *SQLiteRepository) Close() error {
	return repo.conn.Close()
}
