package wallet

import (
	"crypto/ecdsa"
	"encoding/hex"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/tyler-smith/go-bip32"
	"github.com/tyler-smith/go-bip39"
	"os"
	"path/filepath"
)

type WalletDetails struct {
	Wallet     *Wallet
	Mnemonic   string
	PrivateKey *ecdsa.PrivateKey
	PublicKey  *ecdsa.PublicKey
}

type WalletService struct {
	Repo     WalletRepository
	KeyStore *keystore.KeyStore
}

func NewWalletService(repo WalletRepository, ks *keystore.KeyStore) *WalletService {
	return &WalletService{
		Repo:     repo,
		KeyStore: ks,
	}
}

func (ws *WalletService) CreateWallet(name, password string) (*WalletDetails, error) {
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

	wallet := &Wallet{
		Name:         name,
		Address:      account.Address.Hex(),
		KeyStorePath: newPath,
		Mnemonic:     mnemonic, // Store the mnemonic
	}

	err = ws.Repo.AddWallet(wallet)
	if err != nil {
		return nil, err
	}

	walletDetails := &WalletDetails{
		Wallet:     wallet,
		Mnemonic:   mnemonic,
		PrivateKey: privKey,
		PublicKey:  &privKey.PublicKey,
	}

	return walletDetails, nil
}

func (ws *WalletService) ImportWallet(name, mnemonic, password string) (*WalletDetails, error) {
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

	wallet := &Wallet{
		Name:         name,
		Address:      account.Address.Hex(),
		KeyStorePath: newPath,
		Mnemonic:     mnemonic, // Store the mnemonic
	}

	err = ws.Repo.AddWallet(wallet)
	if err != nil {
		return nil, err
	}

	walletDetails := &WalletDetails{
		Wallet:     wallet,
		Mnemonic:   mnemonic,
		PrivateKey: privKey,
		PublicKey:  &privKey.PublicKey,
	}

	return walletDetails, nil
}

func (ws *WalletService) ImportWalletFromPrivateKey(name, privateKeyHex, password string) (*WalletDetails, error) {
	// Remove "0x" prefix if present
	if len(privateKeyHex) > 2 && privateKeyHex[:2] == "0x" {
		privateKeyHex = privateKeyHex[2:]
	}

	// Validate private key format
	if len(privateKeyHex) != 64 {
		return nil, fmt.Errorf("invalid private key format")
	}

	// Convert hex to ECDSA private key
	privKey, err := HexToECDSA(privateKeyHex)
	if err != nil {
		return nil, fmt.Errorf("invalid private key: %v", err)
	}

	// Generate a mnemonic from private key
	privateKeyBytes, _ := hex.DecodeString(privateKeyHex)
	entropy := crypto.Keccak256(privateKeyBytes)[:16] // Use first 16 bytes for BIP39 entropy
	mnemonic, err := bip39.NewMnemonic(entropy)
	if err != nil {
		return nil, fmt.Errorf("error generating mnemonic: %v", err)
	}

	// Import the private key to keystore
	account, err := ws.KeyStore.ImportECDSA(privKey, password)
	if err != nil {
		return nil, err
	}

	// Rename the keystore file to match Ethereum address
	originalPath := account.URL.Path
	newFilename := fmt.Sprintf("%s.json", account.Address.Hex())
	newPath := filepath.Join(filepath.Dir(originalPath), newFilename)
	err = os.Rename(originalPath, newPath)
	if err != nil {
		return nil, fmt.Errorf("error renaming the wallet file: %v", err)
	}

	// Create the wallet entry with the generated mnemonic
	wallet := &Wallet{
		Name:         name,
		Address:      account.Address.Hex(),
		KeyStorePath: newPath,
		Mnemonic:     mnemonic, // Store the generated mnemonic
	}

	// Add wallet to repository
	err = ws.Repo.AddWallet(wallet)
	if err != nil {
		return nil, err
	}

	// Return wallet details with the generated mnemonic
	walletDetails := &WalletDetails{
		Wallet:     wallet,
		Mnemonic:   mnemonic,
		PrivateKey: privKey,
		PublicKey:  &privKey.PublicKey,
	}

	return walletDetails, nil
}

func (ws *WalletService) ImportWalletFromKeystore(name, keystorePath, password string) (*WalletDetails, error) {
	// Read the keystore file
	keyJSON, err := os.ReadFile(keystorePath)
	if err != nil {
		return nil, fmt.Errorf("error reading the keystore file: %v", err)
	}

	// Decrypt the keystore file to verify the password and extract the address
	key, err := keystore.DecryptKey(keyJSON, password)
	if err != nil {
		return nil, fmt.Errorf("incorrect password for keystore file")
	}

	// Get the address from the key
	address := key.Address.Hex()

	// Create the destination filename with 0x prefix
	destFilename := fmt.Sprintf("0x%s.json", address[2:]) // Remove "0x" prefix if present and add it back

	// Get the keystore directory
	var keystoreDir string

	// Check if there are any accounts in the keystore
	accounts := ws.KeyStore.Accounts()
	if len(accounts) > 0 {
		// Use the directory of the first account
		keystoreDir = filepath.Dir(accounts[0].URL.Path)
	} else {
		// If there are no accounts, use a default path
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return nil, fmt.Errorf("error getting user home directory: %v", err)
		}
		keystoreDir = filepath.Join(homeDir, ".wallets", "keystore")
	}

	// Create the destination path
	destPath := filepath.Join(keystoreDir, destFilename)

	// Copy the keystore file to the destination
	destFile, err := os.Create(destPath)
	if err != nil {
		return nil, fmt.Errorf("error creating destination file: %v", err)
	}
	defer destFile.Close()

	// Write the keystore JSON to the destination file
	_, err = destFile.Write(keyJSON)
	if err != nil {
		return nil, fmt.Errorf("error writing to destination file: %v", err)
	}

	// Generate a mnemonic from private key (for compatibility with the rest of the app)
	privateKeyBytes := crypto.FromECDSA(key.PrivateKey)
	entropy := crypto.Keccak256(privateKeyBytes)[:16] // Use first 16 bytes for BIP39 entropy
	mnemonic, err := bip39.NewMnemonic(entropy)
	if err != nil {
		return nil, fmt.Errorf("error generating mnemonic: %v", err)
	}

	// Create the wallet entry
	wallet := &Wallet{
		Name:         name,
		Address:      address,
		KeyStorePath: destPath,
		Mnemonic:     mnemonic, // Store the generated mnemonic for compatibility
	}

	// Add wallet to repository
	err = ws.Repo.AddWallet(wallet)
	if err != nil {
		return nil, err
	}

	// Return wallet details
	walletDetails := &WalletDetails{
		Wallet:     wallet,
		Mnemonic:   mnemonic,
		PrivateKey: key.PrivateKey,
		PublicKey:  &key.PrivateKey.PublicKey,
	}

	return walletDetails, nil
}

func (ws *WalletService) LoadWallet(wallet *Wallet, password string) (*WalletDetails, error) {
	keyJSON, err := os.ReadFile(wallet.KeyStorePath)
	if err != nil {
		return nil, fmt.Errorf("error reading the wallet file: %v", err)
	}
	key, err := keystore.DecryptKey(keyJSON, password)
	if err != nil {
		return nil, fmt.Errorf("incorrect password")
	}

	walletDetails := &WalletDetails{
		Wallet:     wallet,
		Mnemonic:   wallet.Mnemonic,
		PrivateKey: key.PrivateKey,
		PublicKey:  &key.PrivateKey.PublicKey,
	}
	return walletDetails, nil
}

func (ws *WalletService) GetAllWallets() ([]Wallet, error) {
	return ws.Repo.GetAllWallets()
}

func (ws *WalletService) DeleteWallet(wallet *Wallet) error {
	// Remove o arquivo keystore do sistema
	err := os.Remove(wallet.KeyStorePath)
	if err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to remove keystore file: %v", err)
	}
	// Remove do banco de dados
	return ws.Repo.DeleteWallet(wallet.ID)
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
