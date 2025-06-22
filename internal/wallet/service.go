package wallet

import (
	"blocowallet/pkg/localization"
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

type Details struct {
	Wallet     *Wallet
	Mnemonic   string
	PrivateKey *ecdsa.PrivateKey
	PublicKey  *ecdsa.PublicKey
}

type Service struct {
	Repo     Repository
	KeyStore *keystore.KeyStore
}

func NewWalletService(repo Repository, ks *keystore.KeyStore) *Service {
	return &Service{
		Repo:     repo,
		KeyStore: ks,
	}
}

func (ws *Service) CreateWallet(name, password string) (*Details, error) {
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
		return nil, fmt.Errorf(localization.T("error_renaming_wallet_file", map[string]interface{}{"v": err}))
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

	walletDetails := &Details{
		Wallet:     wallet,
		Mnemonic:   mnemonic,
		PrivateKey: privKey,
		PublicKey:  &privKey.PublicKey,
	}

	return walletDetails, nil
}

func (ws *Service) ImportWallet(name, mnemonic, password string) (*Details, error) {
	if !bip39.IsMnemonicValid(mnemonic) {
		return nil, fmt.Errorf(localization.T("error_invalid_mnemonic_phrase", nil))
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
		return nil, fmt.Errorf(localization.T("error_renaming_wallet_file", map[string]interface{}{"v": err}))
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

	walletDetails := &Details{
		Wallet:     wallet,
		Mnemonic:   mnemonic,
		PrivateKey: privKey,
		PublicKey:  &privKey.PublicKey,
	}

	return walletDetails, nil
}

func (ws *Service) ImportWalletFromPrivateKey(name, privateKeyHex, password string) (*Details, error) {
	// Remove "0x" prefix if present
	if len(privateKeyHex) > 2 && privateKeyHex[:2] == "0x" {
		privateKeyHex = privateKeyHex[2:]
	}

	// Validate a private key format
	if len(privateKeyHex) != 64 {
		return nil, fmt.Errorf(localization.T("error_invalid_private_key_format", nil))
	}

	// Convert hex to ECDSA private key
	privKey, err := HexToECDSA(privateKeyHex)
	if err != nil {
		return nil, fmt.Errorf(localization.T("error_invalid_private_key", map[string]interface{}{"v": err}))
	}

	// Generate a mnemonic from a private key
	privateKeyBytes, _ := hex.DecodeString(privateKeyHex)
	entropy := crypto.Keccak256(privateKeyBytes)[:16] // Use first 16 bytes for BIP39 entropy
	mnemonic, err := bip39.NewMnemonic(entropy)
	if err != nil {
		return nil, fmt.Errorf(localization.T("error_generating_mnemonic", map[string]interface{}{"v": err}))
	}

	// Import the private key to the keystore
	account, err := ws.KeyStore.ImportECDSA(privKey, password)
	if err != nil {
		return nil, err
	}

	// Rename the keystore file to match the Ethereum address
	originalPath := account.URL.Path
	newFilename := fmt.Sprintf("%s.json", account.Address.Hex())
	newPath := filepath.Join(filepath.Dir(originalPath), newFilename)
	err = os.Rename(originalPath, newPath)
	if err != nil {
		return nil, fmt.Errorf(localization.T("error_renaming_wallet_file", map[string]interface{}{"v": err}))
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
	walletDetails := &Details{
		Wallet:     wallet,
		Mnemonic:   mnemonic,
		PrivateKey: privKey,
		PublicKey:  &privKey.PublicKey,
	}

	return walletDetails, nil
}

func (ws *Service) ImportWalletFromKeystore(name, keystorePath, password string) (*Details, error) {
	// Read the keystore file
	keyJSON, err := os.ReadFile(keystorePath)
	if err != nil {
		return nil, fmt.Errorf(localization.T("error_reading_keystore_file", map[string]interface{}{"v": err}))
	}

	// Decrypt the keystore file to verify the password and extract the address
	key, err := keystore.DecryptKey(keyJSON, password)
	if err != nil {
		return nil, fmt.Errorf(localization.T("error_incorrect_keystore_password", nil))
	}

	// Get the address from the key
	address := key.Address.Hex()

	// Create the destination filename with 0x prefix
	destFilename := fmt.Sprintf("0x%s.json", address[2:]) // Remove the "0x" prefix if present and add it back

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
			return nil, fmt.Errorf(localization.T("error_getting_home_directory", map[string]interface{}{"v": err}))
		}
		keystoreDir = filepath.Join(homeDir, ".wallets", "keystore")
	}

	// Create the destination path
	destPath := filepath.Join(keystoreDir, destFilename)

	// Copy the keystore file to the destination
	destFile, err := os.Create(destPath)
	if err != nil {
		return nil, fmt.Errorf(localization.T("error_creating_destination_file", map[string]interface{}{"v": err}))
	}
	defer func(destFile *os.File) {
		err := destFile.Close()
		if err != nil {
			fmt.Printf(localization.T("error_writing_destination_file", map[string]interface{}{"v": err}))
		}
	}(destFile)

	// Write the keystore JSON to the destination file
	_, err = destFile.Write(keyJSON)
	if err != nil {
		return nil, fmt.Errorf(localization.T("error_writing_destination_file", map[string]interface{}{"v": err}))
	}

	// Generate a mnemonic from a private key (for compatibility with the rest of the app)
	privateKeyBytes := crypto.FromECDSA(key.PrivateKey)
	entropy := crypto.Keccak256(privateKeyBytes)[:16] // Use first 16 bytes for BIP39 entropy
	mnemonic, err := bip39.NewMnemonic(entropy)
	if err != nil {
		return nil, fmt.Errorf(localization.T("error_generating_mnemonic", map[string]interface{}{"v": err}))
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
	walletDetails := &Details{
		Wallet:     wallet,
		Mnemonic:   mnemonic,
		PrivateKey: key.PrivateKey,
		PublicKey:  &key.PrivateKey.PublicKey,
	}

	return walletDetails, nil
}

func (ws *Service) LoadWallet(wallet *Wallet, password string) (*Details, error) {
	keyJSON, err := os.ReadFile(wallet.KeyStorePath)
	if err != nil {
		return nil, fmt.Errorf(localization.T("error_reading_wallet_file", map[string]interface{}{"v": err}))
	}
	key, err := keystore.DecryptKey(keyJSON, password)
	if err != nil {
		return nil, fmt.Errorf(localization.T("error_incorrect_password", nil))
	}

	walletDetails := &Details{
		Wallet:     wallet,
		Mnemonic:   wallet.Mnemonic,
		PrivateKey: key.PrivateKey,
		PublicKey:  &key.PrivateKey.PublicKey,
	}
	return walletDetails, nil
}

func (ws *Service) GetAllWallets() ([]Wallet, error) {
	return ws.Repo.GetAllWallets()
}

func (ws *Service) DeleteWallet(wallet *Wallet) error {
	// Remove the keystore file from the system
	err := os.Remove(wallet.KeyStorePath)
	if err != nil && !os.IsNotExist(err) {
		return fmt.Errorf(localization.T("error_remove_keystore_file", map[string]interface{}{"v": err}))
	}
	// Remove from the database
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
		return "", fmt.Errorf(localization.T("error_invalid_mnemonic_phrase", nil))
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
