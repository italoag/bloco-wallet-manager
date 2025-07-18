package wallet

import "time"

// Wallet representa uma carteira de criptomoeda
type Wallet struct {
	ID           int       `gorm:"primaryKey"`
	Name         string    `gorm:"not null"`
	Address      string    `gorm:"uniqueIndex;not null"`
	KeyStorePath string    `gorm:"not null"`
	Mnemonic     string    `gorm:"type:text"` // Removed NOT NULL constraint for migration compatibility
	CreatedAt    time.Time `gorm:"not null;autoCreateTime"`
}

// TableName define o nome da tabela no banco de dados
func (Wallet) TableName() string {
	return "wallets"
}
