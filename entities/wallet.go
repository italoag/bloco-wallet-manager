package entities

type Wallet struct {
	ID           int
	Address      string
	KeyStorePath string
	Mnemonic     string
}
