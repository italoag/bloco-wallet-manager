package wallet

// Wallet representa uma carteira de criptomoeda
type Wallet struct {
	ID           int    `gorm:"primaryKey"`
	Address      string `gorm:"uniqueIndex;not null"`
	KeyStorePath string `gorm:"not null"`
	Mnemonic     string `gorm:"not null"`
}

// TableName define o nome da tabela no banco de dados
func (Wallet) TableName() string {
	return "wallets"
}
