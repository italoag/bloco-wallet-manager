package wallet

import "time"

// Wallet represents a cryptocurrency wallet
type Wallet struct {
	ID           int       `gorm:"primaryKey"`
	Name         string    `gorm:"not null"`
	Address      string    `gorm:"uniqueIndex;not null"`
	KeyStorePath string    `gorm:"not null"`
	Mnemonic     string    `gorm:"not null"`
	CreatedAt    time.Time `gorm:"not null;autoCreateTime"`
}

// TableName defines the table name in the database
func (Wallet) TableName() string {
	return "wallets"
}
