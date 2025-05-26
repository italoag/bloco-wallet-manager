package database

import (
	"blocowallet/internal/domain/entities"
	"blocowallet/internal/usecases"
	"database/sql"
	"errors"
)

type PostgresWalletRepository struct {
	//TODO: implementar conexão com o banco de dados
	db *sql.DB
}

func (r *PostgresWalletRepository) Close() error {
	//TODO implement me
	panic("implement me")
}

func NewPostgresWalletRepository(url string) (usecases.WalletRepository, error) {
	return &PostgresWalletRepository{
		db: nil,
	}, nil
}

func (r *PostgresWalletRepository) CreateWallet(wallet entities.Wallet) error {
	// TODO: Implementação para criar a wallet no banco de dados
	return nil
}

func (r *PostgresWalletRepository) GetWallet(address string) (entities.Wallet, error) {
	// TODO: Implementação para obter a wallet do banco de dados
	return entities.Wallet{}, nil
}

func (r *PostgresWalletRepository) ListWallets() ([]entities.Wallet, error) {
	// TODO: Implementação para listar wallets do banco de dados
	return nil, errors.New("not implemented")
}

func (r *PostgresWalletRepository) DeleteWallet(address string) error {
	// TODO: Implementação para deletar a wallet do banco de dados
	return nil
}
