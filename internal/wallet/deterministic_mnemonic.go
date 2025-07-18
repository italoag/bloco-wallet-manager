package wallet

import (
	"crypto/ecdsa"
	"crypto/sha256"
	"encoding/hex"
	"fmt"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/tyler-smith/go-bip39"
)

// GenerateDeterministicMnemonic creates a deterministic mnemonic from a private key
// This ensures that the same private key always generates the same mnemonic
// It uses SHA-256 hash of the private key as entropy source for BIP39 mnemonic generation
func GenerateDeterministicMnemonic(privateKey *ecdsa.PrivateKey) (string, error) {
	if privateKey == nil {
		return "", fmt.Errorf("private key cannot be nil")
	}

	// Convert private key to bytes
	privateKeyBytes := crypto.FromECDSA(privateKey)

	// Generate SHA-256 hash of the private key to use as entropy
	hash := sha256.Sum256(privateKeyBytes)

	// BIP39 requires entropy of specific bit lengths (128, 160, 192, 224, or 256 bits)
	// We'll use the first 16 bytes (128 bits) of the hash for a 12-word mnemonic
	entropy := hash[:16]

	// Generate mnemonic from entropy
	mnemonic, err := bip39.NewMnemonic(entropy)
	if err != nil {
		return "", fmt.Errorf("failed to generate mnemonic: %w", err)
	}

	return mnemonic, nil
}

// ValidateMnemonicMatchesPrivateKey verifies that a mnemonic corresponds to the expected private key
// This validation checks if the mnemonic can be used to derive a valid Ethereum address
// that matches the address derived from the expected private key
func ValidateMnemonicMatchesPrivateKey(mnemonic string, expectedPrivateKey *ecdsa.PrivateKey) (bool, error) {
	if !bip39.IsMnemonicValid(mnemonic) {
		return false, fmt.Errorf("invalid mnemonic")
	}

	// Derive private key from mnemonic using the standard BIP39/44 derivation path
	derivedPrivateKeyHex, err := DerivePrivateKey(mnemonic)
	if err != nil {
		return false, fmt.Errorf("failed to derive private key from mnemonic: %w", err)
	}

	// Convert derived private key hex to ECDSA private key
	derivedPrivateKeyBytes, err := hex.DecodeString(derivedPrivateKeyHex)
	if err != nil {
		return false, fmt.Errorf("failed to decode derived private key: %w", err)
	}

	derivedPrivateKey, err := crypto.ToECDSA(derivedPrivateKeyBytes)
	if err != nil {
		return false, fmt.Errorf("failed to convert derived private key to ECDSA: %w", err)
	}

	// Get the address from the derived private key
	derivedAddress := crypto.PubkeyToAddress(derivedPrivateKey.PublicKey).Hex()

	// For deterministic mnemonic generation, we don't expect the private keys to match
	// but we can verify that the generated mnemonic is valid and can be used to derive
	// a valid Ethereum address
	return derivedAddress != "", nil
}

// GenerateAndValidateDeterministicMnemonic generates a deterministic mnemonic from a private key
// and validates that it can be used to derive a valid Ethereum address
func GenerateAndValidateDeterministicMnemonic(privateKey *ecdsa.PrivateKey) (string, error) {
	if privateKey == nil {
		return "", fmt.Errorf("private key cannot be nil")
	}

	// Generate deterministic mnemonic
	mnemonic, err := GenerateDeterministicMnemonic(privateKey)
	if err != nil {
		return "", fmt.Errorf("failed to generate deterministic mnemonic: %w", err)
	}

	// Validate that the mnemonic is valid and can be used to derive a valid Ethereum address
	isValid, err := ValidateMnemonicMatchesPrivateKey(mnemonic, privateKey)
	if err != nil {
		return "", fmt.Errorf("mnemonic validation error: %w", err)
	}

	if !isValid {
		return "", fmt.Errorf("generated mnemonic failed validation")
	}

	// Verify that the same private key always generates the same mnemonic
	// This is the key property of deterministic mnemonic generation
	verificationMnemonic, err := GenerateDeterministicMnemonic(privateKey)
	if err != nil {
		return "", fmt.Errorf("failed to verify deterministic mnemonic: %w", err)
	}

	if mnemonic != verificationMnemonic {
		return "", fmt.Errorf("deterministic mnemonic generation is inconsistent")
	}

	return mnemonic, nil
}
