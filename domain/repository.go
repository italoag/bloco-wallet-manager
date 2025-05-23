package domain

type WalletRepository interface {
	AddWallet(wallet *Wallet) error
	GetAllWallets() ([]Wallet, error)
	DeleteWallet(walletID int) error
	Close() error
}
