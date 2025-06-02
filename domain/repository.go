package domain

type WalletRepository interface {
	AddWallet(wallet *Wallet) error
	DeleteWallet(id int) error
	GetAllWallets() ([]Wallet, error)
	DeleteWallet(walletID int) error
	Close() error
}
