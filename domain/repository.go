package domain

type WalletRepository interface {
	AddWallet(wallet *Wallet) error
	DeleteWallet(id int) error
	GetAllWallets() ([]Wallet, error)
	Close() error
}
