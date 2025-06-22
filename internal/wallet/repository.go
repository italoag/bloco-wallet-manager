package wallet

type Repository interface {
	AddWallet(wallet *Wallet) error
	GetAllWallets() ([]Wallet, error)
	DeleteWallet(walletID int) error
	Close() error
}
