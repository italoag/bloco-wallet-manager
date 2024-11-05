package gateway

import (
	"blocowallet/internal/domain/entities"
)

type WalletRepository interface {
	CreateWallet(wallet *entities.Wallet) error
	ListWallets() ([]entities.Wallet, error)
	Close() error
}
