package entities

type Wallet struct {
	ID           string
	Address      string
	KeyStorePath string
	Mnemonic     string
	Balance      float64 // Saldo em Ether
}
