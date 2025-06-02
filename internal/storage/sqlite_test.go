package storage

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	"blocowallet/internal/wallet"
)

func TestSQLiteCreateAndGet(t *testing.T) {
	// Create temporary database file
	tempDir, err := os.MkdirTemp("", "blocowallet_test")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	dbPath := filepath.Join(tempDir, "test.db")
	db, err := NewSQLite(dbPath)
	if err != nil {
		t.Fatalf("Failed to create SQLite instance: %v", err)
	}
	defer db.Close()

	ctx := context.Background()

	// Create test wallet
	testWallet := &wallet.Wallet{
		ID:                "test-id-1",
		Name:              "Test Wallet",
		Address:           "0x1234567890123456789012345678901234567890",
		KeyStorePath:      "/path/to/keystore.json",
		EncryptedMnemonic: "encrypted_mnemonic_data",
		CreatedAt:         time.Now().UTC(),
		UpdatedAt:         time.Now().UTC(),
	}

	// Test Create
	err = db.Create(ctx, testWallet)
	if err != nil {
		t.Fatalf("Failed to create wallet: %v", err)
	}

	// Test GetByID
	retrieved, err := db.GetByID(ctx, testWallet.ID)
	if err != nil {
		t.Fatalf("Failed to get wallet by ID: %v", err)
	}

	// Verify all fields
	if retrieved.ID != testWallet.ID {
		t.Errorf("ID mismatch. Got: %s, Expected: %s", retrieved.ID, testWallet.ID)
	}
	if retrieved.Name != testWallet.Name {
		t.Errorf("Name mismatch. Got: %s, Expected: %s", retrieved.Name, testWallet.Name)
	}
	if retrieved.Address != testWallet.Address {
		t.Errorf("Address mismatch. Got: %s, Expected: %s", retrieved.Address, testWallet.Address)
	}
	if retrieved.KeyStorePath != testWallet.KeyStorePath {
		t.Errorf("KeyStorePath mismatch. Got: %s, Expected: %s", retrieved.KeyStorePath, testWallet.KeyStorePath)
	}
	if retrieved.EncryptedMnemonic != testWallet.EncryptedMnemonic {
		t.Errorf("EncryptedMnemonic mismatch. Got: %s, Expected: %s", retrieved.EncryptedMnemonic, testWallet.EncryptedMnemonic)
	}

	// Test GetByAddress
	retrievedByAddr, err := db.GetByAddress(ctx, testWallet.Address)
	if err != nil {
		t.Fatalf("Failed to get wallet by address: %v", err)
	}

	if retrievedByAddr.ID != testWallet.ID {
		t.Errorf("GetByAddress returned wrong wallet. Got ID: %s, Expected: %s", retrievedByAddr.ID, testWallet.ID)
	}
}

func TestSQLiteList(t *testing.T) {
	// Create temporary database file
	tempDir, err := os.MkdirTemp("", "blocowallet_test")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	dbPath := filepath.Join(tempDir, "test.db")
	db, err := NewSQLite(dbPath)
	if err != nil {
		t.Fatalf("Failed to create SQLite instance: %v", err)
	}
	defer db.Close()

	ctx := context.Background()

	// Create multiple test wallets
	wallets := []*wallet.Wallet{
		{
			ID:                "test-id-1",
			Name:              "Wallet 1",
			Address:           "0x1111111111111111111111111111111111111111",
			KeyStorePath:      "/path/to/keystore1.json",
			EncryptedMnemonic: "encrypted_mnemonic_1",
			CreatedAt:         time.Now().UTC().Add(-2 * time.Hour),
			UpdatedAt:         time.Now().UTC(),
		},
		{
			ID:                "test-id-2",
			Name:              "Wallet 2",
			Address:           "0x2222222222222222222222222222222222222222",
			KeyStorePath:      "/path/to/keystore2.json",
			EncryptedMnemonic: "encrypted_mnemonic_2",
			CreatedAt:         time.Now().UTC().Add(-1 * time.Hour),
			UpdatedAt:         time.Now().UTC(),
		},
		{
			ID:                "test-id-3",
			Name:              "Wallet 3",
			Address:           "0x3333333333333333333333333333333333333333",
			KeyStorePath:      "/path/to/keystore3.json",
			EncryptedMnemonic: "",
			CreatedAt:         time.Now().UTC(),
			UpdatedAt:         time.Now().UTC(),
		},
	}

	// Create all wallets
	for _, w := range wallets {
		err = db.Create(ctx, w)
		if err != nil {
			t.Fatalf("Failed to create wallet %s: %v", w.Name, err)
		}
	}

	// Test List
	retrievedWallets, err := db.List(ctx)
	if err != nil {
		t.Fatalf("Failed to list wallets: %v", err)
	}

	if len(retrievedWallets) != len(wallets) {
		t.Fatalf("List returned wrong number of wallets. Got: %d, Expected: %d", len(retrievedWallets), len(wallets))
	}

	// Verify order (should be DESC by created_at, so newest first)
	if retrievedWallets[0].ID != "test-id-3" {
		t.Errorf("First wallet should be newest (test-id-3), got: %s", retrievedWallets[0].ID)
	}
	if retrievedWallets[2].ID != "test-id-1" {
		t.Errorf("Last wallet should be oldest (test-id-1), got: %s", retrievedWallets[2].ID)
	}
}

func TestSQLiteUpdate(t *testing.T) {
	// Create temporary database file
	tempDir, err := os.MkdirTemp("", "blocowallet_test")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	dbPath := filepath.Join(tempDir, "test.db")
	db, err := NewSQLite(dbPath)
	if err != nil {
		t.Fatalf("Failed to create SQLite instance: %v", err)
	}
	defer db.Close()

	ctx := context.Background()

	// Create test wallet
	testWallet := &wallet.Wallet{
		ID:                "test-id-1",
		Name:              "Original Name",
		Address:           "0x1234567890123456789012345678901234567890",
		KeyStorePath:      "/path/to/keystore.json",
		EncryptedMnemonic: "original_encrypted_mnemonic",
		CreatedAt:         time.Now().UTC(),
		UpdatedAt:         time.Now().UTC(),
	}

	// Create wallet
	err = db.Create(ctx, testWallet)
	if err != nil {
		t.Fatalf("Failed to create wallet: %v", err)
	}

	// Update wallet
	testWallet.Name = "Updated Name"
	testWallet.EncryptedMnemonic = "updated_encrypted_mnemonic"
	testWallet.UpdatedAt = time.Now().UTC()

	err = db.Update(ctx, testWallet)
	if err != nil {
		t.Fatalf("Failed to update wallet: %v", err)
	}

	// Retrieve and verify update
	retrieved, err := db.GetByID(ctx, testWallet.ID)
	if err != nil {
		t.Fatalf("Failed to get updated wallet: %v", err)
	}

	if retrieved.Name != "Updated Name" {
		t.Errorf("Name not updated. Got: %s, Expected: Updated Name", retrieved.Name)
	}
	if retrieved.EncryptedMnemonic != "updated_encrypted_mnemonic" {
		t.Errorf("EncryptedMnemonic not updated. Got: %s, Expected: updated_encrypted_mnemonic", retrieved.EncryptedMnemonic)
	}
}

func TestSQLiteDelete(t *testing.T) {
	// Create temporary database file
	tempDir, err := os.MkdirTemp("", "blocowallet_test")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	dbPath := filepath.Join(tempDir, "test.db")
	db, err := NewSQLite(dbPath)
	if err != nil {
		t.Fatalf("Failed to create SQLite instance: %v", err)
	}
	defer db.Close()

	ctx := context.Background()

	// Create test wallet
	testWallet := &wallet.Wallet{
		ID:                "test-id-1",
		Name:              "Test Wallet",
		Address:           "0x1234567890123456789012345678901234567890",
		KeyStorePath:      "/path/to/keystore.json",
		EncryptedMnemonic: "encrypted_mnemonic",
		CreatedAt:         time.Now().UTC(),
		UpdatedAt:         time.Now().UTC(),
	}

	// Create wallet
	err = db.Create(ctx, testWallet)
	if err != nil {
		t.Fatalf("Failed to create wallet: %v", err)
	}

	// Verify it exists
	_, err = db.GetByID(ctx, testWallet.ID)
	if err != nil {
		t.Fatalf("Wallet should exist before deletion: %v", err)
	}

	// Delete wallet
	err = db.Delete(ctx, testWallet.ID)
	if err != nil {
		t.Fatalf("Failed to delete wallet: %v", err)
	}

	// Verify it's gone
	_, err = db.GetByID(ctx, testWallet.ID)
	if err == nil {
		t.Fatal("Wallet should not exist after deletion")
	}
}

func TestSQLiteWithEmptyEncryptedMnemonic(t *testing.T) {
	// Create temporary database file
	tempDir, err := os.MkdirTemp("", "blocowallet_test")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	dbPath := filepath.Join(tempDir, "test.db")
	db, err := NewSQLite(dbPath)
	if err != nil {
		t.Fatalf("Failed to create SQLite instance: %v", err)
	}
	defer db.Close()

	ctx := context.Background()

	// Create wallet without encrypted mnemonic (imported from private key)
	testWallet := &wallet.Wallet{
		ID:                "test-id-1",
		Name:              "Imported Wallet",
		Address:           "0x1234567890123456789012345678901234567890",
		KeyStorePath:      "/path/to/keystore.json",
		EncryptedMnemonic: "", // Empty for imported wallets
		CreatedAt:         time.Now().UTC(),
		UpdatedAt:         time.Now().UTC(),
	}

	// Create wallet
	err = db.Create(ctx, testWallet)
	if err != nil {
		t.Fatalf("Failed to create wallet with empty mnemonic: %v", err)
	}

	// Retrieve and verify
	retrieved, err := db.GetByID(ctx, testWallet.ID)
	if err != nil {
		t.Fatalf("Failed to get wallet with empty mnemonic: %v", err)
	}

	if retrieved.EncryptedMnemonic != "" {
		t.Errorf("EncryptedMnemonic should be empty. Got: %s", retrieved.EncryptedMnemonic)
	}
}
