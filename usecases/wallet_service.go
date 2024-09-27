package usecases

import (
	"crypto/ecdsa"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"

	"blocowallet/entities"

	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/tyler-smith/go-bip32"
	"github.com/tyler-smith/go-bip39"
)

type WalletService struct {
	Repo     entities.WalletRepository
	KeyStore *keystore.KeyStore
}

func NewWalletService(repo entities.WalletRepository, ks *keystore.KeyStore) *WalletService {
	return &WalletService{
		Repo:     repo,
		KeyStore: ks,
	}
}

func (ws *WalletService) CreateWallet(password string) (*entities.Wallet, error) {
	mnemonic, err := GenerateMnemonic()
	if err != nil {
		return nil, err
	}

	privateKeyHex, err := DerivePrivateKey(mnemonic)
	if err != nil {
		return nil, err
	}

	privKey, err := HexToECDSA(privateKeyHex)
	if err != nil {
		return nil, err
	}

	account, err := ws.KeyStore.ImportECDSA(privKey, password)
	if err != nil {
		return nil, err
	}

	originalPath := account.URL.Path
	newFilename := fmt.Sprintf("%s.json", account.Address.Hex())
	newPath := filepath.Join(filepath.Dir(originalPath), newFilename)
	err = os.Rename(originalPath, newPath)
	if err != nil {
		return nil, fmt.Errorf("error renaming the wallet file: %v", err)
	}

	wallet := &entities.Wallet{
		Address:      account.Address.Hex(),
		KeyStorePath: newPath,
		Mnemonic:     "", // Do not store the mnemonic
	}

	err = ws.Repo.AddWallet(wallet)
	if err != nil {
		return nil, err
	}

	return wallet, nil
}

func (ws *WalletService) ImportWallet(mnemonic, password string) (*entities.Wallet, error) {
	if !bip39.IsMnemonicValid(mnemonic) {
		return nil, fmt.Errorf("invalid mnemonic phrase")
	}

	privateKeyHex, err := DerivePrivateKey(mnemonic)
	if err != nil {
		return nil, err
	}

	privKey, err := HexToECDSA(privateKeyHex)
	if err != nil {
		return nil, err
	}

	account, err := ws.KeyStore.ImportECDSA(privKey, password)
	if err != nil {
		return nil, err
	}

	originalPath := account.URL.Path
	newFilename := fmt.Sprintf("%s.json", account.Address.Hex())
	newPath := filepath.Join(filepath.Dir(originalPath), newFilename)
	err = os.Rename(originalPath, newPath)
	if err != nil {
		return nil, fmt.Errorf("error renaming the wallet file: %v", err)
	}

	wallet := &entities.Wallet{
		Address:      account.Address.Hex(),
		KeyStorePath: newPath,
		Mnemonic:     "", // Do not store the mnemonic
	}

	err = ws.Repo.AddWallet(wallet)
	if err != nil {
		return nil, err
	}

	return wallet, nil
}

func (ws *WalletService) GetAllWallets() ([]entities.Wallet, error) {
	return ws.Repo.GetAllWallets()
}

func (ws *WalletService) LoadWallet(wallet *entities.Wallet, password string) error {
	keyJSON, err := os.ReadFile(wallet.KeyStorePath)
	if err != nil {
		return fmt.Errorf("error reading the wallet file: %v", err)
	}
	_, err = keystore.DecryptKey(keyJSON, password)
	if err != nil {
		return fmt.Errorf("incorrect password")
	}
	return nil
}

// Helper functions

func GenerateMnemonic() (string, error) {
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

func DerivePrivateKey(mnemonic string) (string, error) {
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

func HexToECDSA(hexkey string) (*ecdsa.PrivateKey, error) {
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
