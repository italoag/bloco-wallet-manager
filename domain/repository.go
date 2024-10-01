package domain

type WalletRepository interface {
	AddWallet(wallet *Wallet) error
	GetAllWallets() ([]Wallet, error)
	Close() error
}
