package domain

type Wallet struct {
	ID           int
	Address      string
	KeyStorePath string
	Mnemonic     string
}
