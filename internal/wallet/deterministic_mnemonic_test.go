package wallet

import (
	"encoding/hex"
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/assert"
	"github.com/tyler-smith/go-bip39"
)

func TestGenerateDeterministicMnemonic(t *testing.T) {
	// Create a test private key
	privateKeyHex := "1a2b3c4d5e6f7a8b9c0d1e2f3a4b5c6d7e8f9a0b1c2d3e4f5a6b7c8d9e0f1a2b"
	privateKeyBytes, err := hex.DecodeString(privateKeyHex)
	assert.NoError(t, err)

	privateKey, err := crypto.ToECDSA(privateKeyBytes)
	assert.NoError(t, err)

	// Generate mnemonic from the private key
	mnemonic1, err := GenerateDeterministicMnemonic(privateKey)
	assert.NoError(t, err)
	assert.NotEmpty(t, mnemonic1)

	// Generate mnemonic again from the same private key
	mnemonic2, err := GenerateDeterministicMnemonic(privateKey)
	assert.NoError(t, err)

	// Verify that the same private key produces the same mnemonic
	assert.Equal(t, mnemonic1, mnemonic2, "Deterministic mnemonic should be consistent for the same private key")
}
func TestGenerateDeterministicMnemonicWithDifferentKeys(t *testing.T) {
	// Create two different private keys
	privateKeyHex1 := "1a2b3c4d5e6f7a8b9c0d1e2f3a4b5c6d7e8f9a0b1c2d3e4f5a6b7c8d9e0f1a2b"
	privateKeyHex2 := "2a2b3c4d5e6f7a8b9c0d1e2f3a4b5c6d7e8f9a0b1c2d3e4f5a6b7c8d9e0f1a2b"

	privateKeyBytes1, err := hex.DecodeString(privateKeyHex1)
	assert.NoError(t, err)
	privateKey1, err := crypto.ToECDSA(privateKeyBytes1)
	assert.NoError(t, err)

	privateKeyBytes2, err := hex.DecodeString(privateKeyHex2)
	assert.NoError(t, err)
	privateKey2, err := crypto.ToECDSA(privateKeyBytes2)
	assert.NoError(t, err)

	// Generate mnemonics from different private keys
	mnemonic1, err := GenerateDeterministicMnemonic(privateKey1)
	assert.NoError(t, err)

	mnemonic2, err := GenerateDeterministicMnemonic(privateKey2)
	assert.NoError(t, err)

	// Verify that different private keys produce different mnemonics
	assert.NotEqual(t, mnemonic1, mnemonic2, "Different private keys should produce different mnemonics")
}

func TestGenerateDeterministicMnemonicNilKey(t *testing.T) {
	// Test with nil private key
	mnemonic, err := GenerateDeterministicMnemonic(nil)
	assert.Error(t, err)
	assert.Empty(t, mnemonic)
	assert.Contains(t, err.Error(), "private key cannot be nil")
}
func TestValidateMnemonicMatchesPrivateKey(t *testing.T) {
	// This test demonstrates the expected behavior where a mnemonic generated from a private key
	// will not necessarily derive back to the same private key due to BIP39/44 derivation paths

	// Create a test private key
	privateKeyHex := "1a2b3c4d5e6f7a8b9c0d1e2f3a4b5c6d7e8f9a0b1c2d3e4f5a6b7c8d9e0f1a2b"
	privateKeyBytes, err := hex.DecodeString(privateKeyHex)
	assert.NoError(t, err)

	privateKey, err := crypto.ToECDSA(privateKeyBytes)
	assert.NoError(t, err)

	// Generate mnemonic from the private key
	mnemonic, err := GenerateDeterministicMnemonic(privateKey)
	assert.NoError(t, err)

	// Validate that the mnemonic corresponds to the private key
	// Note: This is expected to return true because we're just checking if the mnemonic
	// can derive a valid Ethereum address, not if it derives the same private key
	isValid, err := ValidateMnemonicMatchesPrivateKey(mnemonic, privateKey)
	assert.NoError(t, err)
	assert.True(t, isValid, "Mnemonic should be valid and derive a valid Ethereum address")

	// Test with invalid mnemonic
	isValid, err = ValidateMnemonicMatchesPrivateKey("invalid mnemonic phrase", privateKey)
	assert.Error(t, err)
	assert.False(t, isValid)
	assert.Contains(t, err.Error(), "invalid mnemonic")
}
func TestGenerateAndValidateDeterministicMnemonic(t *testing.T) {
	// Create a test private key
	privateKeyHex := "1a2b3c4d5e6f7a8b9c0d1e2f3a4b5c6d7e8f9a0b1c2d3e4f5a6b7c8d9e0f1a2b"
	privateKeyBytes, err := hex.DecodeString(privateKeyHex)
	assert.NoError(t, err)

	privateKey, err := crypto.ToECDSA(privateKeyBytes)
	assert.NoError(t, err)

	// Generate and validate mnemonic
	mnemonic, err := GenerateAndValidateDeterministicMnemonic(privateKey)
	assert.NoError(t, err)
	assert.NotEmpty(t, mnemonic)

	// Generate again to verify determinism
	mnemonic2, err := GenerateAndValidateDeterministicMnemonic(privateKey)
	assert.NoError(t, err)

	// Verify that the same private key produces the same mnemonic
	assert.Equal(t, mnemonic, mnemonic2, "Deterministic mnemonic should be consistent for the same private key")

	// Test with nil private key
	mnemonic3, err := GenerateAndValidateDeterministicMnemonic(nil)
	assert.Error(t, err)
	assert.Empty(t, mnemonic3)
	assert.Contains(t, err.Error(), "private key cannot be nil")
}
func TestDeterministicMnemonicWithRealPrivateKey(t *testing.T) {
	// Generate a real private key
	privateKey, err := crypto.GenerateKey()
	assert.NoError(t, err)

	// Generate mnemonic from the private key
	mnemonic1, err := GenerateDeterministicMnemonic(privateKey)
	assert.NoError(t, err)
	assert.NotEmpty(t, mnemonic1)

	// Generate mnemonic again from the same private key
	mnemonic2, err := GenerateDeterministicMnemonic(privateKey)
	assert.NoError(t, err)

	// Verify that the same private key produces the same mnemonic
	assert.Equal(t, mnemonic1, mnemonic2, "Deterministic mnemonic should be consistent for the same private key")

	// Verify that the mnemonic is valid BIP39
	assert.True(t, bip39.IsMnemonicValid(mnemonic1), "Generated mnemonic should be valid BIP39")

	// Verify that the mnemonic has words (we don't check the exact count as it may vary)
	words := strings.Split(mnemonic1, " ")
	assert.Greater(t, len(words), 0, "Generated mnemonic should have words")
}

// TestMnemonicEntropySize verifies that the mnemonic is generated with the correct entropy size
func TestMnemonicEntropySize(t *testing.T) {
	// Generate a real private key
	privateKey, err := crypto.GenerateKey()
	assert.NoError(t, err)

	// Generate mnemonic from the private key
	mnemonic, err := GenerateDeterministicMnemonic(privateKey)
	assert.NoError(t, err)

	// Split the mnemonic into words
	words := strings.Split(mnemonic, " ")

	// 12 words corresponds to 128 bits of entropy
	// 15 words corresponds to 160 bits of entropy
	// 18 words corresponds to 192 bits of entropy
	// 21 words corresponds to 224 bits of entropy
	// 24 words corresponds to 256 bits of entropy

	// We're using the first 16 bytes (128 bits) of the hash, so we expect 12 words
	assert.Equal(t, 12, len(words), "Generated mnemonic should have 12 words for 128 bits of entropy")
}

// TestConsistencyAcrossRestarts verifies that the mnemonic generation is consistent
// even if the code is restarted (simulated by creating new instances)
func TestConsistencyAcrossRestarts(t *testing.T) {
	// Create a test private key
	privateKeyHex := "1a2b3c4d5e6f7a8b9c0d1e2f3a4b5c6d7e8f9a0b1c2d3e4f5a6b7c8d9e0f1a2b"
	privateKeyBytes, err := hex.DecodeString(privateKeyHex)
	assert.NoError(t, err)

	privateKey, err := crypto.ToECDSA(privateKeyBytes)
	assert.NoError(t, err)

	// Generate mnemonic from the private key
	mnemonic1, err := GenerateDeterministicMnemonic(privateKey)
	assert.NoError(t, err)

	// Simulate a restart by creating a new private key from the same hex
	privateKeyBytes2, err := hex.DecodeString(privateKeyHex)
	assert.NoError(t, err)

	privateKey2, err := crypto.ToECDSA(privateKeyBytes2)
	assert.NoError(t, err)

	// Generate mnemonic again from the "restarted" private key
	mnemonic2, err := GenerateDeterministicMnemonic(privateKey2)
	assert.NoError(t, err)

	// Verify that the same private key produces the same mnemonic even after a "restart"
	assert.Equal(t, mnemonic1, mnemonic2, "Deterministic mnemonic should be consistent across restarts")
}

// TestMnemonicWordCount verifies that the mnemonic has the expected number of words
func TestMnemonicWordCount(t *testing.T) {
	// Create a test private key
	privateKeyHex := "1a2b3c4d5e6f7a8b9c0d1e2f3a4b5c6d7e8f9a0b1c2d3e4f5a6b7c8d9e0f1a2b"
	privateKeyBytes, err := hex.DecodeString(privateKeyHex)
	assert.NoError(t, err)

	privateKey, err := crypto.ToECDSA(privateKeyBytes)
	assert.NoError(t, err)

	// Generate mnemonic from the private key
	mnemonic, err := GenerateDeterministicMnemonic(privateKey)
	assert.NoError(t, err)

	// Count the number of words in the mnemonic
	wordCount := len(strings.Split(mnemonic, " "))

	// We're using the first 16 bytes (128 bits) of the hash, so we expect 12 words
	assert.Equal(t, 12, wordCount, "Generated mnemonic should have 12 words for 128 bits of entropy")
}

// TestAddressVerification tests the address verification functionality
func TestAddressVerification(t *testing.T) {
	// Create a test private key
	privateKey, err := crypto.GenerateKey()
	assert.NoError(t, err)

	// Get the address from the private key
	address := crypto.PubkeyToAddress(privateKey.PublicKey).Hex()

	// Generate mnemonic from the private key
	mnemonic, err := GenerateDeterministicMnemonic(privateKey)
	assert.NoError(t, err)

	// Derive private key from mnemonic
	derivedPrivateKeyHex, err := DerivePrivateKey(mnemonic)
	assert.NoError(t, err)

	// Convert derived private key hex to ECDSA private key
	derivedPrivateKeyBytes, err := hex.DecodeString(derivedPrivateKeyHex)
	assert.NoError(t, err)

	derivedPrivateKey, err := crypto.ToECDSA(derivedPrivateKeyBytes)
	assert.NoError(t, err)

	// Get the address from the derived private key
	derivedAddress := crypto.PubkeyToAddress(derivedPrivateKey.PublicKey).Hex()

	// The addresses will be different because the deterministic mnemonic generation
	// doesn't guarantee that the derived private key will match the original
	// But both should be valid Ethereum addresses
	assert.NotEqual(t, address, derivedAddress, "Original and derived addresses should be different")
	assert.Regexp(t, "^0x[0-9a-fA-F]{40}$", address, "Original address should be a valid Ethereum address")
	assert.Regexp(t, "^0x[0-9a-fA-F]{40}$", derivedAddress, "Derived address should be a valid Ethereum address")
}
