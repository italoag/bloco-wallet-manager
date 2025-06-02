package wallet

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/ethereum/go-ethereum/crypto"
)

func TestCreateKeyStoreV3File(t *testing.T) {
	// Create temporary directory for test
	tempDir, err := os.MkdirTemp("", "blocowallet_test")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Override the home directory for this test
	originalHome := os.Getenv("HOME")
	defer os.Setenv("HOME", originalHome)
	os.Setenv("HOME", tempDir)

	// Create a service instance (without repository for this test)
	service := &Service{}

	// Generate a test private key
	privateKey, err := crypto.GenerateKey()
	if err != nil {
		t.Fatalf("Failed to generate private key: %v", err)
	}

	password := "testpassword123"
	address := GetAddressFromPrivateKey(privateKey)

	// Test keystore creation
	keystorePath, err := service.CreateKeyStoreV3File(privateKey, password)
	if err != nil {
		t.Fatalf("Failed to create keystore file: %v", err)
	}

	// Verify file exists
	if _, err := os.Stat(keystorePath); os.IsNotExist(err) {
		t.Fatalf("Keystore file was not created at path: %s", keystorePath)
	}

	// Verify file is in correct directory
	expectedDir := filepath.Join(tempDir, KeyStoreV3Dir)
	if !filepath.HasPrefix(keystorePath, expectedDir) {
		t.Fatalf("Keystore file not in expected directory. Got: %s, Expected prefix: %s", keystorePath, expectedDir)
	}

	// Verify filename contains address
	fileName := filepath.Base(keystorePath)
	if fileName != address+".json" {
		t.Fatalf("Keystore filename incorrect. Got: %s, Expected: %s", fileName, address+".json")
	}

	// Test loading private key back
	loadedKey, err := service.LoadPrivateKeyFromKeyStoreV3(keystorePath, password)
	if err != nil {
		t.Fatalf("Failed to load private key from keystore: %v", err)
	}

	// Verify addresses match
	originalAddress := crypto.PubkeyToAddress(privateKey.PublicKey).Hex()
	loadedAddress := crypto.PubkeyToAddress(loadedKey.PublicKey).Hex()

	if originalAddress != loadedAddress {
		t.Fatalf("Addresses don't match. Original: %s, Loaded: %s", originalAddress, loadedAddress)
	}

	// Test with wrong password
	_, err = service.LoadPrivateKeyFromKeyStoreV3(keystorePath, "wrongpassword")
	if err == nil {
		t.Fatal("Expected error with wrong password, but got none")
	}
}

func TestEncryptDecryptMnemonic(t *testing.T) {
	testMnemonic := "abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon about"
	password := "testpassword123"

	// Test encryption
	encrypted, err := EncryptMnemonic(testMnemonic, password)
	if err != nil {
		t.Fatalf("Failed to encrypt mnemonic: %v", err)
	}

	if encrypted == "" {
		t.Fatal("Encrypted mnemonic is empty")
	}

	if encrypted == testMnemonic {
		t.Fatal("Encrypted mnemonic should not be the same as original")
	}

	// Test decryption with correct password
	decrypted, err := DecryptMnemonic(encrypted, password)
	if err != nil {
		t.Fatalf("Failed to decrypt mnemonic: %v", err)
	}

	if decrypted != testMnemonic {
		t.Fatalf("Decrypted mnemonic doesn't match original. Got: %s, Expected: %s", decrypted, testMnemonic)
	}

	// Test decryption with wrong password
	_, err = DecryptMnemonic(encrypted, "wrongpassword")
	if err == nil {
		t.Fatal("Expected error with wrong password, but got none")
	}
}

func TestPasswordCache(t *testing.T) {
	service := &Service{
		passwordCache: make(map[string]string),
	}

	address := "0x1234567890123456789012345678901234567890"
	password := "testpassword123"

	// Test setting password
	service.SetWalletPassword(address, password)

	// Test getting password
	cachedPassword, exists := service.GetWalletPassword(address)
	if !exists {
		t.Fatal("Password should exist in cache")
	}

	if cachedPassword != password {
		t.Fatalf("Cached password doesn't match. Got: %s, Expected: %s", cachedPassword, password)
	}

	// Test getting non-existent password
	_, exists = service.GetWalletPassword("0xnonexistent")
	if exists {
		t.Fatal("Password should not exist for non-existent address")
	}

	// Test clearing password
	service.ClearWalletPassword(address)
	_, exists = service.GetWalletPassword(address)
	if exists {
		t.Fatal("Password should not exist after clearing")
	}
}

func TestGenerateMnemonic(t *testing.T) {
	mnemonic, err := GenerateMnemonic()
	if err != nil {
		t.Fatalf("Failed to generate mnemonic: %v", err)
	}

	if mnemonic == "" {
		t.Fatal("Generated mnemonic is empty")
	}

	// Test that we get different mnemonics on consecutive calls
	mnemonic2, err := GenerateMnemonic()
	if err != nil {
		t.Fatalf("Failed to generate second mnemonic: %v", err)
	}

	if mnemonic == mnemonic2 {
		t.Fatal("Generated mnemonics should be different")
	}
}

func TestDerivePrivateKey(t *testing.T) {
	testMnemonic := "abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon about"

	privateKeyHex, err := DerivePrivateKey(testMnemonic)
	if err != nil {
		t.Fatalf("Failed to derive private key: %v", err)
	}

	if privateKeyHex == "" {
		t.Fatal("Derived private key is empty")
	}

	// Test conversion to ECDSA
	privateKey, err := HexToECDSA(privateKeyHex)
	if err != nil {
		t.Fatalf("Failed to convert hex to ECDSA: %v", err)
	}

	// Verify we can get an address from it
	address := GetAddressFromPrivateKey(privateKey)
	if address == "" {
		t.Fatal("Generated address is empty")
	}

	// Test that the same mnemonic always generates the same key
	privateKeyHex2, err := DerivePrivateKey(testMnemonic)
	if err != nil {
		t.Fatalf("Failed to derive private key second time: %v", err)
	}

	if privateKeyHex != privateKeyHex2 {
		t.Fatal("Same mnemonic should generate same private key")
	}
}
