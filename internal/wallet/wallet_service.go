package wallet

import (
	"crypto/ecdsa"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"

	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/tyler-smith/go-bip32"
	"github.com/tyler-smith/go-bip39"
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

	// Encrypt the mnemonic before storing
	encryptedMnemonic, err := EncryptMnemonic(mnemonic, password)
	if err != nil {
		return nil, fmt.Errorf("failed to encrypt mnemonic: %v", err)
	}

	wallet := &Wallet{
		Name:         name,
		Address:      account.Address.Hex(),
		KeyStorePath: newPath,
		Mnemonic:     encryptedMnemonic, // Store the encrypted mnemonic
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

	// Encrypt the mnemonic before storing
	encryptedMnemonic, err := EncryptMnemonic(mnemonic, password)
	if err != nil {
		return nil, fmt.Errorf("failed to encrypt mnemonic: %v", err)
	}

	wallet := &Wallet{
		Name:         name,
		Address:      account.Address.Hex(),
		KeyStorePath: newPath,
		Mnemonic:     encryptedMnemonic, // Store the encrypted mnemonic
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

	// Generate a deterministic mnemonic from private key
	mnemonic, err := GenerateDeterministicMnemonic(privKey)
	if err != nil {
		return nil, fmt.Errorf("error generating deterministic mnemonic: %v", err)
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

	// Encrypt the mnemonic before storing
	encryptedMnemonic, err := EncryptMnemonic(mnemonic, password)
	if err != nil {
		return nil, fmt.Errorf("failed to encrypt mnemonic: %v", err)
	}

	// Create the wallet entry with the encrypted mnemonic
	wallet := &Wallet{
		Name:         name,
		Address:      account.Address.Hex(),
		KeyStorePath: newPath,
		Mnemonic:     encryptedMnemonic, // Store the encrypted mnemonic
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

// ImportWalletFromKeystoreV3 imports a wallet from a keystore v3 file with enhanced validation
func (ws *WalletService) ImportWalletFromKeystoreV3(name, keystorePath, password string) (*WalletDetails, error) {
	// Step 1: Validate file existence
	if _, err := os.Stat(keystorePath); os.IsNotExist(err) {
		return nil, NewKeystoreImportError(
			ErrorFileNotFound,
			"Keystore file not found at specified path",
			err,
		)
	}

	// Step 2: Read the keystore file
	keyJSON, err := os.ReadFile(keystorePath)
	if err != nil {
		return nil, NewKeystoreImportError(
			ErrorFileNotFound,
			"Error reading the keystore file",
			err,
		)
	}

	// Step 3: Validate keystore structure
	validator := &KeystoreValidator{}
	keystoreData, err := validator.ValidateKeystoreV3(keyJSON)
	if err != nil {
		// The validator already returns a KeystoreImportError
		return nil, err
	}

	// Step 4: Decrypt the keystore file to verify the password and extract the private key
	key, err := keystore.DecryptKey(keyJSON, password)
	if err != nil {
		return nil, NewKeystoreImportError(
			ErrorIncorrectPassword,
			"Incorrect password for keystore file",
			err,
		)
	}

	// Step 5: Verify that the decrypted private key matches the address in the keystore
	derivedAddress := crypto.PubkeyToAddress(key.PrivateKey.PublicKey).Hex()

	// Normalize addresses for comparison (ensure both have 0x prefix and are lowercase)
	normalizedKeystoreAddress := common.HexToAddress(keystoreData.Address).Hex()
	normalizedDerivedAddress := common.HexToAddress(derivedAddress).Hex()

	if normalizedKeystoreAddress != normalizedDerivedAddress {
		return nil, NewKeystoreImportError(
			ErrorAddressMismatch,
			fmt.Sprintf("Address mismatch: keystore address %s does not match derived address %s",
				normalizedKeystoreAddress, normalizedDerivedAddress),
			nil,
		)
	}

	// Step 6: Generate a deterministic mnemonic from private key
	mnemonic, err := GenerateAndValidateDeterministicMnemonic(key.PrivateKey)
	if err != nil {
		return nil, NewKeystoreImportError(
			ErrorCorruptedFile,
			"Error generating deterministic mnemonic",
			err,
		)
	}

	// Step 7: Create the destination filename with proper address format
	address := key.Address.Hex()
	destFilename := fmt.Sprintf("%s.json", address)

	// Step 8: Get the keystore directory
	var keystoreDir string
	accounts := ws.KeyStore.Accounts()
	if len(accounts) > 0 {
		keystoreDir = filepath.Dir(accounts[0].URL.Path)
	} else {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return nil, NewKeystoreImportError(
				ErrorFileNotFound,
				"Error getting user home directory",
				err,
			)
		}
		keystoreDir = filepath.Join(homeDir, ".wallets", "keystore")

		// Ensure the directory exists
		if err := os.MkdirAll(keystoreDir, 0700); err != nil {
			return nil, NewKeystoreImportError(
				ErrorFileNotFound,
				"Error creating keystore directory",
				err,
			)
		}
	}

	// Step 9: Create the destination path
	destPath := filepath.Join(keystoreDir, destFilename)

	// Step 10: Copy the keystore file to the destination
	destFile, err := os.Create(destPath)
	if err != nil {
		return nil, NewKeystoreImportError(
			ErrorFileNotFound,
			"Error creating destination file",
			err,
		)
	}
	defer destFile.Close()

	// Write the keystore JSON to the destination file
	_, err = destFile.Write(keyJSON)
	if err != nil {
		return nil, NewKeystoreImportError(
			ErrorFileNotFound,
			"Error writing to destination file",
			err,
		)
	}

	// Step 11: Encrypt the mnemonic before storing
	encryptedMnemonic, err := EncryptMnemonic(mnemonic, password)
	if err != nil {
		return nil, NewKeystoreImportError(
			ErrorCorruptedFile,
			"Failed to encrypt mnemonic",
			err,
		)
	}

	// Step 12: Create the wallet entry
	wallet := &Wallet{
		Name:         name,
		Address:      address,
		KeyStorePath: destPath,
		Mnemonic:     encryptedMnemonic,
	}

	// Step 13: Add wallet to repository
	err = ws.Repo.AddWallet(wallet)
	if err != nil {
		return nil, NewKeystoreImportError(
			ErrorCorruptedFile,
			"Failed to add wallet to repository",
			err,
		)
	}

	// Step 14: Return wallet details
	walletDetails := &WalletDetails{
		Wallet:     wallet,
		Mnemonic:   mnemonic,
		PrivateKey: key.PrivateKey,
		PublicKey:  &key.PrivateKey.PublicKey,
	}

	return walletDetails, nil
}

// ImportWalletFromKeystore is kept for backward compatibility
// It calls the new ImportWalletFromKeystoreV3 function
func (ws *WalletService) ImportWalletFromKeystore(name, keystorePath, password string) (*WalletDetails, error) {
	return ws.ImportWalletFromKeystoreV3(name, keystorePath, password)
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

	// Decrypt the mnemonic
	decryptedMnemonic, err := DecryptMnemonic(wallet.Mnemonic, password)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt mnemonic: %v", err)
	}

	walletDetails := &WalletDetails{
		Wallet:     wallet,
		Mnemonic:   decryptedMnemonic,
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
