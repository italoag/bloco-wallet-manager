package wallet

import (
	"context"
	"crypto/ecdsa"
	"encoding/hex"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/tyler-smith/go-bip32"
	"github.com/tyler-smith/go-bip39"
)

// Wallet represents a blockchain wallet
type Wallet struct {
	ID                string    `json:"id" db:"id"`
	Name              string    `json:"name" db:"name"`
	Address           string    `json:"address" db:"address"`
	KeyStorePath      string    `json:"keystore_path" db:"keystore_path"`
	EncryptedMnemonic string    `json:"encrypted_mnemonic,omitempty" db:"encrypted_mnemonic"`
	CreatedAt         time.Time `json:"created_at" db:"created_at"`
	UpdatedAt         time.Time `json:"updated_at" db:"updated_at"`
}

// Repository defines wallet storage operations
type Repository interface {
	Create(ctx context.Context, wallet *Wallet) error
	GetByID(ctx context.Context, id string) (*Wallet, error)
	GetByAddress(ctx context.Context, address string) (*Wallet, error)
	List(ctx context.Context) ([]*Wallet, error)
	Update(ctx context.Context, wallet *Wallet) error
	Delete(ctx context.Context, id string) error
	Close() error
}

// Balance represents wallet balance information
type Balance struct {
	Address   string    `json:"address"`
	Amount    *big.Int  `json:"amount"`
	Symbol    string    `json:"symbol"`
	Decimals  int       `json:"decimals"`
	UpdatedAt time.Time `json:"updated_at"`
}

// MultiNetworkBalance represents a wallet balance across multiple networks
type MultiNetworkBalance struct {
	Address         string            `json:"address"`
	NetworkBalances []*NetworkBalance `json:"network_balances"`
	UpdatedAt       time.Time         `json:"updated_at"`
}

// NetworkBalance represents a balance on a specific network
type NetworkBalance struct {
	NetworkKey  string   `json:"network_key"`
	NetworkName string   `json:"network_name"`
	Amount      *big.Int `json:"amount"`
	Symbol      string   `json:"symbol"`
	Decimals    int      `json:"decimals"`
	Error       error    `json:"error,omitempty"`
}

// BalanceProvider defines operations for getting balance from blockchain
type BalanceProvider interface {
	GetBalance(ctx context.Context, address string) (*big.Int, error)
	GetNetworkSymbol() string
	GetNetworkDecimals() int
}

// Cryptographic helper functions for wallet creation

// GenerateMnemonic creates a new BIP39 mnemonic phrase
func GenerateMnemonic() (string, error) {
	entropy, err := bip39.NewEntropy(128)
	if err != nil {
		return "", fmt.Errorf("failed to generate entropy: %w", err)
	}

	mnemonic, err := bip39.NewMnemonic(entropy)
	if err != nil {
		return "", fmt.Errorf("failed to generate mnemonic: %w", err)
	}

	return mnemonic, nil
}

// DerivePrivateKey derives a private key from a mnemonic phrase using BIP44 derivation path
func DerivePrivateKey(mnemonic string) (string, error) {
	if !bip39.IsMnemonicValid(mnemonic) {
		return "", fmt.Errorf("invalid mnemonic phrase")
	}

	seed := bip39.NewSeed(mnemonic, "")

	// BIP44 derivation path for Ethereum: m/44'/60'/0'/0/0
	masterKey, err := bip32.NewMasterKey(seed)
	if err != nil {
		return "", fmt.Errorf("failed to create master key: %w", err)
	}

	// Purpose (44')
	purposeKey, err := masterKey.NewChildKey(bip32.FirstHardenedChild + 44)
	if err != nil {
		return "", fmt.Errorf("failed to derive purpose key: %w", err)
	}

	// Coin type (60' for Ethereum)
	coinTypeKey, err := purposeKey.NewChildKey(bip32.FirstHardenedChild + 60)
	if err != nil {
		return "", fmt.Errorf("failed to derive coin type key: %w", err)
	}

	// Account (0')
	accountKey, err := coinTypeKey.NewChildKey(bip32.FirstHardenedChild + 0)
	if err != nil {
		return "", fmt.Errorf("failed to derive account key: %w", err)
	}

	// Change (0)
	changeKey, err := accountKey.NewChildKey(0)
	if err != nil {
		return "", fmt.Errorf("failed to derive change key: %w", err)
	}

	// Address index (0)
	addressKey, err := changeKey.NewChildKey(0)
	if err != nil {
		return "", fmt.Errorf("failed to derive address key: %w", err)
	}

	privateKeyBytes := addressKey.Key
	return hex.EncodeToString(privateKeyBytes), nil
}

// HexToECDSA converts a hex string to an ECDSA private key
func HexToECDSA(hexkey string) (*ecdsa.PrivateKey, error) {
	privateKeyBytes, err := hex.DecodeString(hexkey)
	if err != nil {
		return nil, fmt.Errorf("failed to decode hex string: %w", err)
	}

	privateKey, err := crypto.ToECDSA(privateKeyBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to convert to ECDSA: %w", err)
	}

	return privateKey, nil
}

// GetAddressFromPrivateKey derives an Ethereum address from a private key
func GetAddressFromPrivateKey(privateKey *ecdsa.PrivateKey) string {
	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return ""
	}

	address := crypto.PubkeyToAddress(*publicKeyECDSA)
	return address.Hex()
}

// WalletDetails contains full wallet information including cryptographic details
type WalletDetails struct {
	Wallet     *Wallet
	Mnemonic   string
	PrivateKey *ecdsa.PrivateKey
	PublicKey  *ecdsa.PublicKey
}

// IsValidMnemonic validates a BIP39 mnemonic phrase
func IsValidMnemonic(mnemonic string) bool {
	return bip39.IsMnemonicValid(mnemonic)
}
