package usecases

import (
	"blocowallet/internal/domain/entities"
	"crypto/ecdsa"
	"encoding/hex"
	"fmt"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/tyler-smith/go-bip32"
	"github.com/tyler-smith/go-bip39"
)

type WalletRepository interface {
	CreateWallet(wallet entities.Wallet) error
	GetWallet(address string) (entities.Wallet, error)
	ListWallets() ([]entities.Wallet, error)
	DeleteWallet(address string) error
	Close() error
}

type WalletUseCase interface {
	CreateWallet() (entities.Wallet, error)
	GetWallet(address string) (entities.Wallet, error)
	ListWallets() ([]entities.Wallet, error)
	DeleteWallet(address string) error
}

type walletUseCase struct {
	repo WalletRepository
}

func NewWalletUseCase(repo WalletRepository) WalletUseCase {
	return &walletUseCase{
		repo: repo,
	}
}

func (u *walletUseCase) CreateWallet() (entities.Wallet, error) {
	wallet := entities.Wallet{
		Address: "0x1234567890abcdef",                  // Lógica para gerar o endereço
		ID:      "0a21b3c4-5d6e-7f8-9g0h-1i2j3k4l5m6n", // Lógica para gerar o ID
	}
	err := u.repo.CreateWallet(wallet)
	if err != nil {
		return entities.Wallet{}, err
	}
	return wallet, nil
}

func (u *walletUseCase) GetWallet(address string) (entities.Wallet, error) {
	return u.repo.GetWallet(address)
}

func (u *walletUseCase) ListWallets() ([]entities.Wallet, error) {
	return u.repo.ListWallets()
}

func (u *walletUseCase) DeleteWallet(address string) error {
	return u.repo.DeleteWallet(address)
}

// Helper functions

func CreateMnemonic() (string, error) {
	entropy, err := bip39.NewEntropy(128)
	if err != nil {
		return "", err
	}
	mnemonic, err := bip39.NewMnemonic(entropy)
	if err != nil {
		return "", err
	}
	return mnemonic, nil
}

func DerivePk(mnemonic string) (string, error) {
	if !bip39.IsMnemonicValid(mnemonic) {
		return "", fmt.Errorf("invalid mnemonic phrase")
	}
	seed := bip39.NewSeed(mnemonic, "")
	masterKey, err := bip32.NewMasterKey(seed)
	if err != nil {
		return "", err
	}
	purposeKey, err := masterKey.NewChildKey(bip32.FirstHardenedChild + 44)
	if err != nil {
		return "", err
	}
	coinTypeKey, err := purposeKey.NewChildKey(bip32.FirstHardenedChild + 60)
	if err != nil {
		return "", err
	}
	accountKey, err := coinTypeKey.NewChildKey(bip32.FirstHardenedChild + 0)
	if err != nil {
		return "", err
	}
	changeKey, err := accountKey.NewChildKey(0)
	if err != nil {
		return "", err
	}
	addressKey, err := changeKey.NewChildKey(0)
	if err != nil {
		return "", err
	}
	privateKeyBytes := addressKey.Key
	return hex.EncodeToString(privateKeyBytes), nil
}

func ConvertHexToECDSA(hexkey string) (*ecdsa.PrivateKey, error) {
	privateKeyBytes, err := hex.DecodeString(hexkey)
	if err != nil {
		return nil, err
	}
	privateKey, err := crypto.ToECDSA(privateKeyBytes)
	if err != nil {
		return nil, err
	}
	return privateKey, nil
}
