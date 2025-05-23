package infrastructure

import (
	"blocowallet/domain"
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

type SQLiteRepository struct {
	conn *sql.DB
}

// Implement the WalletRepository interface from entities package
var _ domain.WalletRepository = &SQLiteRepository{}

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

func (repo *SQLiteRepository) AddWallet(wallet *domain.Wallet) error {
	insertQuery := `
	INSERT INTO wallets (address, keystore_path, mnemonic)
	VALUES (?, ?, ?);
	`
	_, err := repo.conn.Exec(insertQuery, wallet.Address, wallet.KeyStorePath, wallet.Mnemonic)
	return err
}

func (repo *SQLiteRepository) GetAllWallets() ([]domain.Wallet, error) {
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

	var wallets []domain.Wallet
	for rows.Next() {
		var w domain.Wallet
		err := rows.Scan(&w.ID, &w.Address, &w.KeyStorePath, &w.Mnemonic)
		if err != nil {
			return nil, err
		}
		wallets = append(wallets, w)
	}

	return wallets, nil
}

func (repo *SQLiteRepository) DeleteWallet(walletID int) error {
	deleteQuery := `DELETE FROM wallets WHERE id = ?;`
	_, err := repo.conn.Exec(deleteQuery, walletID)
	return err
}

func (repo *SQLiteRepository) Close() error {
	return repo.conn.Close()
}
