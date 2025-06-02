package wallet

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum/crypto"
)

// MockRepository implements Repository interface for testing
type MockRepository struct {
	wallets map[string]*Wallet
	nextID  int
}

func NewMockRepository() *MockRepository {
	return &MockRepository{
		wallets: make(map[string]*Wallet),
		nextID:  1,
	}
}

func (m *MockRepository) Create(ctx context.Context, w *Wallet) error {
	if w.ID == "" {
		w.ID = fmt.Sprintf("%d", m.nextID)
		m.nextID++
	}
	m.wallets[w.ID] = w
	return nil
}

func (m *MockRepository) GetByID(ctx context.Context, id string) (*Wallet, error) {
	if w, exists := m.wallets[id]; exists {
		return w, nil
	}
	return nil, fmt.Errorf("wallet not found")
}

func (m *MockRepository) GetByAddress(ctx context.Context, address string) (*Wallet, error) {
	for _, w := range m.wallets {
		if w.Address == address {
			return w, nil
		}
	}
	return nil, fmt.Errorf("wallet not found")
}

func (m *MockRepository) List(ctx context.Context) ([]*Wallet, error) {
	var wallets []*Wallet
	for _, w := range m.wallets {
		wallets = append(wallets, w)
	}
	return wallets, nil
}

func (m *MockRepository) Update(ctx context.Context, w *Wallet) error {
	if _, exists := m.wallets[w.ID]; !exists {
		return fmt.Errorf("wallet not found")
	}
	m.wallets[w.ID] = w
	return nil
}

func (m *MockRepository) Delete(ctx context.Context, id string) error {
	if _, exists := m.wallets[id]; !exists {
		return fmt.Errorf("wallet not found")
	}
	delete(m.wallets, id)
	return nil
}

func (m *MockRepository) Close() error {
	return nil
}

// TestCreateKeyStoreV3File tests keystore file creation
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

	// Create a service instance with mock repository
	service := NewService(NewMockRepository(), nil)

	// Generate a test private key
	privateKey, err := crypto.GenerateKey()
	if err != nil {
		t.Fatalf("Failed to generate private key: %v", err)
	}

	password := "testpassword123"

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
	if !strings.HasPrefix(keystorePath, expectedDir) {
		t.Fatalf("Keystore file not in expected directory. Got: %s, Expected prefix: %s", keystorePath, expectedDir)
	}

	// Test loading private key back from keystore
	loadedPrivateKey, err := service.LoadPrivateKeyFromKeyStoreV3(keystorePath, password)
	if err != nil {
		t.Fatalf("Failed to load private key from keystore: %v", err)
	}

	// Compare the original and loaded private keys
	originalAddress := GetAddressFromPrivateKey(privateKey)
	loadedAddress := GetAddressFromPrivateKey(loadedPrivateKey)

	if originalAddress != loadedAddress {
		t.Fatalf("Addresses don't match. Original: %s, Loaded: %s", originalAddress, loadedAddress)
	}
}

// TestCreateKeyStoreV3FilePasswordScenarios tests keystore creation with various password scenarios
func TestCreateKeyStoreV3FilePasswordScenarios(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "blocowallet_test")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	originalHome := os.Getenv("HOME")
	defer os.Setenv("HOME", originalHome)
	os.Setenv("HOME", tempDir)

	service := NewService(NewMockRepository(), nil)

	privateKey, err := crypto.GenerateKey()
	if err != nil {
		t.Fatalf("Failed to generate private key: %v", err)
	}

	testCases := []struct {
		name        string
		password    string
		shouldError bool
	}{
		{"Empty password", "", true},
		{"Short password", "a", false},
		{"Normal password", "password123", false},
		{"Long password", strings.Repeat("x", 1000), false},
		{"Unicode password", "пароль123", false},
		{"Special characters", "!@#$%^&*()", false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			keystorePath, err := service.CreateKeyStoreV3File(privateKey, tc.password)

			if tc.shouldError {
				if err == nil {
					t.Fatal("Expected error but got none")
				}
				return
			}

			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}

			// Verify we can load the key back
			_, err = service.LoadPrivateKeyFromKeyStoreV3(keystorePath, tc.password)
			if err != nil {
				t.Fatalf("Failed to load key with password %q: %v", tc.password, err)
			}
		})
	}
}

// TestLoadPrivateKeyFromKeyStoreV3WithWrongPassword tests loading with wrong password
func TestLoadPrivateKeyFromKeyStoreV3WithWrongPassword(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "blocowallet_test")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	originalHome := os.Getenv("HOME")
	defer os.Setenv("HOME", originalHome)
	os.Setenv("HOME", tempDir)

	service := NewService(NewMockRepository(), nil)

	privateKey, err := crypto.GenerateKey()
	if err != nil {
		t.Fatalf("Failed to generate private key: %v", err)
	}

	correctPassword := "correct_password"
	wrongPassword := "wrong_password"

	keystorePath, err := service.CreateKeyStoreV3File(privateKey, correctPassword)
	if err != nil {
		t.Fatalf("Failed to create keystore: %v", err)
	}

	// Test with wrong password
	_, err = service.LoadPrivateKeyFromKeyStoreV3(keystorePath, wrongPassword)
	if err == nil {
		t.Fatal("Expected error with wrong password but got none")
	}

	// Test with correct password should work
	_, err = service.LoadPrivateKeyFromKeyStoreV3(keystorePath, correctPassword)
	if err != nil {
		t.Fatalf("Failed to load with correct password: %v", err)
	}
}

// TestPasswordCache tests the password caching functionality
func TestPasswordCache(t *testing.T) {
	service := NewService(NewMockRepository(), nil)

	address := "0x742d35Cc6634C0532925a3b8D86CAE2e1aD84c3e"
	password := "test_password"

	// Initially should have no password
	cachedPassword, exists := service.GetWalletPassword(address)
	if exists || cachedPassword != "" {
		t.Fatalf("Expected empty password and false, got %q and %v", cachedPassword, exists)
	}

	// Set password
	service.SetWalletPassword(address, password)

	// Should now return the cached password
	cachedPassword, exists = service.GetWalletPassword(address)
	if !exists || cachedPassword != password {
		t.Fatalf("Expected %q and true, got %q and %v", password, cachedPassword, exists)
	}

	// Clear password
	service.ClearWalletPassword(address)

	// Should be empty again
	cachedPassword, exists = service.GetWalletPassword(address)
	if exists || cachedPassword != "" {
		t.Fatalf("Expected empty password and false after clear, got %q and %v", cachedPassword, exists)
	}
}

// TestConcurrentPasswordCache tests thread safety of password cache
func TestConcurrentPasswordCache(t *testing.T) {
	service := NewService(NewMockRepository(), nil)

	// Run concurrent operations
	done := make(chan bool)

	// Goroutine 1: Set passwords
	go func() {
		for i := 0; i < 100; i++ {
			address := fmt.Sprintf("0x%040d", i)
			password := fmt.Sprintf("password_%d", i)
			service.SetWalletPassword(address, password)
		}
		done <- true
	}()

	// Goroutine 2: Get passwords
	go func() {
		for i := 0; i < 100; i++ {
			address := fmt.Sprintf("0x%040d", i)
			service.GetWalletPassword(address)
		}
		done <- true
	}()

	// Goroutine 3: Clear passwords
	go func() {
		for i := 0; i < 100; i++ {
			address := fmt.Sprintf("0x%040d", i)
			service.ClearWalletPassword(address)
		}
		done <- true
	}()

	// Wait for all goroutines to complete
	for i := 0; i < 3; i++ {
		<-done
	}
}

// TestWalletCRUDOperations tests basic CRUD operations
func TestWalletCRUDOperations(t *testing.T) {
	service := NewService(NewMockRepository(), nil)
	ctx := context.Background()

	// Create wallet
	name := "Test Wallet"
	address := "0x742d35Cc6634C0532925a3b8D86CAE2e1aD84c3e"
	keystorePath := "/path/to/keystore"

	wallet, err := service.Create(ctx, name, address, keystorePath)
	if err != nil {
		t.Fatalf("Failed to create wallet: %v", err)
	}

	if wallet.Name != name {
		t.Fatalf("Expected name %q, got %q", name, wallet.Name)
	}

	if wallet.Address != address {
		t.Fatalf("Expected address %q, got %q", address, wallet.Address)
	}

	// Get wallet by ID
	retrievedWallet, err := service.GetByID(ctx, wallet.ID)
	if err != nil {
		t.Fatalf("Failed to get wallet by ID: %v", err)
	}

	if retrievedWallet.ID != wallet.ID {
		t.Fatalf("Expected ID %q, got %q", wallet.ID, retrievedWallet.ID)
	}

	// Get wallet by address
	retrievedWallet, err = service.GetByAddress(ctx, address)
	if err != nil {
		t.Fatalf("Failed to get wallet by address: %v", err)
	}

	if retrievedWallet.Address != address {
		t.Fatalf("Expected address %q, got %q", address, retrievedWallet.Address)
	}

	// List wallets
	wallets, err := service.List(ctx)
	if err != nil {
		t.Fatalf("Failed to list wallets: %v", err)
	}

	if len(wallets) != 1 {
		t.Fatalf("Expected 1 wallet, got %d", len(wallets))
	}

	// Update wallet
	wallet.Name = "Updated Wallet"
	err = service.Update(ctx, wallet)
	if err != nil {
		t.Fatalf("Failed to update wallet: %v", err)
	}

	// Verify update
	updatedWallet, err := service.GetByID(ctx, wallet.ID)
	if err != nil {
		t.Fatalf("Failed to get updated wallet: %v", err)
	}

	if updatedWallet.Name != "Updated Wallet" {
		t.Fatalf("Expected name 'Updated Wallet', got %q", updatedWallet.Name)
	}

	// Delete wallet
	err = service.Delete(ctx, wallet.ID)
	if err != nil {
		t.Fatalf("Failed to delete wallet: %v", err)
	}

	// Verify deletion
	_, err = service.GetByID(ctx, wallet.ID)
	if err == nil {
		t.Fatal("Expected error when getting deleted wallet")
	}
}

// TestCreateWalletWithMnemonic tests wallet creation with mnemonic
func TestCreateWalletWithMnemonic(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "blocowallet_test")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	originalHome := os.Getenv("HOME")
	defer os.Setenv("HOME", originalHome)
	os.Setenv("HOME", tempDir)

	service := NewService(NewMockRepository(), nil)
	ctx := context.Background()

	name := "Mnemonic Wallet"
	password := "password123"

	walletDetails, err := service.CreateWalletWithMnemonic(ctx, name, password)
	if err != nil {
		t.Fatalf("Failed to create wallet with mnemonic: %v", err)
	}

	if walletDetails.Wallet.Name != name {
		t.Fatalf("Expected name %q, got %q", name, walletDetails.Wallet.Name)
	}

	if walletDetails.Wallet.Address == "" {
		t.Fatal("Address should not be empty")
	}

	if walletDetails.Mnemonic == "" {
		t.Fatal("Mnemonic should not be empty")
	}

	// Verify mnemonic has 12 words
	words := strings.Fields(walletDetails.Mnemonic)
	if len(words) != 12 {
		t.Fatalf("Expected 12 words in mnemonic, got %d", len(words))
	}
}

// TestImportWalletFromMnemonic tests wallet import from mnemonic
func TestImportWalletFromMnemonic(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "blocowallet_test")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	originalHome := os.Getenv("HOME")
	defer os.Setenv("HOME", originalHome)
	os.Setenv("HOME", tempDir)

	service := NewService(NewMockRepository(), nil)
	ctx := context.Background()

	name := "Imported Wallet"
	mnemonic := "abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon about"
	password := "password123"

	walletDetails, err := service.ImportWalletFromMnemonic(ctx, name, mnemonic, password)
	if err != nil {
		t.Fatalf("Failed to import wallet from mnemonic: %v", err)
	}

	if walletDetails.Wallet.Name != name {
		t.Fatalf("Expected name %q, got %q", name, walletDetails.Wallet.Name)
	}

	if walletDetails.Mnemonic != mnemonic {
		t.Fatalf("Expected mnemonic %q, got %q", mnemonic, walletDetails.Mnemonic)
	}

	// Test getting mnemonic from wallet
	retrievedMnemonic, err := service.GetMnemonicFromWallet(walletDetails.Wallet, password)
	if err != nil {
		t.Fatalf("Failed to get mnemonic from wallet: %v", err)
	}

	if retrievedMnemonic != mnemonic {
		t.Fatalf("Retrieved mnemonic doesn't match. Expected %q, got %q", mnemonic, retrievedMnemonic)
	}
}

// TestImportWalletFromPrivateKey tests wallet import from private key
func TestImportWalletFromPrivateKey(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "blocowallet_test")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	originalHome := os.Getenv("HOME")
	defer os.Setenv("HOME", originalHome)
	os.Setenv("HOME", tempDir)

	service := NewService(NewMockRepository(), nil)
	ctx := context.Background()

	// Generate a test private key
	privateKey, err := crypto.GenerateKey()
	if err != nil {
		t.Fatalf("Failed to generate private key: %v", err)
	}

	privateKeyHex := fmt.Sprintf("%x", crypto.FromECDSA(privateKey)) // Convert to hex without 0x prefix
	name := "Private Key Wallet"
	password := "password123"

	walletDetails, err := service.ImportWalletFromPrivateKey(ctx, name, privateKeyHex, password)
	if err != nil {
		t.Fatalf("Failed to import wallet from private key: %v", err)
	}

	if walletDetails.Wallet.Name != name {
		t.Fatalf("Expected name %q, got %q", name, walletDetails.Wallet.Name)
	}

	expectedAddress := GetAddressFromPrivateKey(privateKey)
	if walletDetails.Wallet.Address != expectedAddress {
		t.Fatalf("Expected address %q, got %q", expectedAddress, walletDetails.Wallet.Address)
	}

	// Private key wallets should have empty mnemonic
	if walletDetails.Wallet.EncryptedMnemonic != "" {
		t.Fatal("Private key wallet should have empty encrypted mnemonic")
	}
}

// TestErrorHandling tests various error scenarios
func TestErrorHandling(t *testing.T) {
	service := NewService(NewMockRepository(), nil)
	ctx := context.Background()

	t.Run("Create wallet with empty name", func(t *testing.T) {
		_, err := service.Create(ctx, "", "0x123", "/path")
		if err == nil {
			t.Fatal("Expected error for empty name")
		}
	})

	t.Run("Get non-existent wallet", func(t *testing.T) {
		_, err := service.GetByID(ctx, "non-existent")
		if err == nil {
			t.Fatal("Expected error for non-existent wallet")
		}
	})

	t.Run("Update non-existent wallet", func(t *testing.T) {
		wallet := &Wallet{ID: "non-existent", Name: "Test", Address: "0x123"}
		err := service.Update(ctx, wallet)
		if err == nil {
			t.Fatal("Expected error for updating non-existent wallet")
		}
	})

	t.Run("Delete non-existent wallet", func(t *testing.T) {
		err := service.Delete(ctx, "non-existent")
		if err == nil {
			t.Fatal("Expected error for deleting non-existent wallet")
		}
	})
}

// TestGetMnemonicFromWalletErrors tests error scenarios for mnemonic retrieval
func TestGetMnemonicFromWalletErrors(t *testing.T) {
	service := NewService(NewMockRepository(), nil)

	t.Run("Wallet with no encrypted mnemonic", func(t *testing.T) {
		wallet := &Wallet{
			ID:                "1",
			Name:              "Test",
			Address:           "0x123",
			EncryptedMnemonic: "",
		}

		_, err := service.GetMnemonicFromWallet(wallet, "password")
		if err == nil {
			t.Fatal("Expected error for wallet with no encrypted mnemonic")
		}
	})

	t.Run("Wrong password", func(t *testing.T) {
		// First create a wallet with encrypted mnemonic
		tempDir, err := os.MkdirTemp("", "blocowallet_test")
		if err != nil {
			t.Fatalf("Failed to create temp directory: %v", err)
		}
		defer os.RemoveAll(tempDir)

		originalHome := os.Getenv("HOME")
		defer os.Setenv("HOME", originalHome)
		os.Setenv("HOME", tempDir)

		ctx := context.Background()
		mnemonic := "abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon about"
		correctPassword := "correct_password"
		wrongPassword := "wrong_password"

		walletDetails, err := service.ImportWalletFromMnemonic(ctx, "Test", mnemonic, correctPassword)
		if err != nil {
			t.Fatalf("Failed to import wallet: %v", err)
		}

		// Try with wrong password
		_, err = service.GetMnemonicFromWallet(walletDetails.Wallet, wrongPassword)
		if err == nil {
			t.Fatal("Expected error with wrong password")
		}
	})
}

func TestEncryptDecryptMnemonic(t *testing.T) {
	mnemonic := "abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon about"
	password := "password123"

	// Encrypt
	encrypted, err := EncryptMnemonic(mnemonic, password)
	if err != nil {
		t.Fatalf("Failed to encrypt mnemonic: %v", err)
	}

	// Decrypt
	decrypted, err := DecryptMnemonic(encrypted, password)
	if err != nil {
		t.Fatalf("Failed to decrypt mnemonic: %v", err)
	}

	if decrypted != mnemonic {
		t.Fatalf("Decrypted mnemonic doesn't match. Expected %q, got %q", mnemonic, decrypted)
	}
}

func TestGenerateMnemonic(t *testing.T) {
	mnemonic, err := GenerateMnemonic()
	if err != nil {
		t.Fatalf("Failed to generate mnemonic: %v", err)
	}

	// Basic validation
	words := strings.Fields(mnemonic)
	if len(words) != 12 {
		t.Fatalf("Expected 12 words, got %d", len(words))
	}

	// Generate multiple mnemonics and ensure they're different
	mnemonics := make(map[string]bool)
	for i := 0; i < 10; i++ {
		mnemonic, err := GenerateMnemonic()
		if err != nil {
			t.Fatalf("Failed to generate mnemonic %d: %v", i, err)
		}

		if mnemonics[mnemonic] {
			t.Fatal("Generated duplicate mnemonic")
		}
		mnemonics[mnemonic] = true
	}
}

func TestDerivePrivateKey(t *testing.T) {
	mnemonic := "abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon about"

	privateKeyHex, err := DerivePrivateKey(mnemonic)
	if err != nil {
		t.Fatalf("Failed to derive private key: %v", err)
	}

	if privateKeyHex == "" {
		t.Fatal("Private key should not be empty")
	}

	// Test that the same mnemonic always produces the same private key
	privateKeyHex2, err := DerivePrivateKey(mnemonic)
	if err != nil {
		t.Fatalf("Failed to derive private key second time: %v", err)
	}

	if privateKeyHex != privateKeyHex2 {
		t.Fatalf("Same mnemonic produced different private keys: %s vs %s", privateKeyHex, privateKeyHex2)
	}
}
