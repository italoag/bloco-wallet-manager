package ui

import (
	"strings"
	"testing"
)

// Testes expandidos para melhorar a cobertura

func TestValidateWalletInput_EmptyName(t *testing.T) {
	err := ValidateWalletInput("", "password123")
	if err == nil {
		t.Error("Expected error for empty name")
	}
	if !strings.Contains(err.Error(), "wallet name is required") {
		t.Errorf("Expected name required error, got: %s", err.Error())
	}
}

func TestValidateWalletInput_WhitespaceName(t *testing.T) {
	err := ValidateWalletInput("   ", "password123")
	if err == nil {
		t.Error("Expected error for whitespace-only name")
	}
	if !strings.Contains(err.Error(), "wallet name is required") {
		t.Errorf("Expected name required error, got: %s", err.Error())
	}
}

func TestValidateWalletInput_EmptyPassword(t *testing.T) {
	err := ValidateWalletInput("MyWallet", "")
	if err == nil {
		t.Error("Expected error for empty password")
	}
	if !strings.Contains(err.Error(), "password is required") {
		t.Errorf("Expected password required error, got: %s", err.Error())
	}
}

func TestValidateWalletInput_WhitespacePassword(t *testing.T) {
	err := ValidateWalletInput("MyWallet", "   ")
	if err == nil {
		t.Error("Expected error for whitespace-only password")
	}
	if !strings.Contains(err.Error(), "password is required") {
		t.Errorf("Expected password required error, got: %s", err.Error())
	}
}

func TestValidateWalletInput_ShortPassword(t *testing.T) {
	err := ValidateWalletInput("MyWallet", "12345")
	if err == nil {
		t.Error("Expected error for short password")
	}
	if !strings.Contains(err.Error(), "password must be at least 6 characters long") {
		t.Errorf("Expected short password error, got: %s", err.Error())
	}
}

func TestValidateWalletInput_Valid(t *testing.T) {
	err := ValidateWalletInput("MyWallet", "password123")
	if err != nil {
		t.Errorf("Expected no error for valid input, got: %s", err.Error())
	}
}

func TestValidateWalletInput_ValidWithSpaces(t *testing.T) {
	err := ValidateWalletInput("  MyWallet  ", "  password123  ")
	if err != nil {
		t.Errorf("Expected no error for valid input with spaces, got: %s", err.Error())
	}
}

func TestValidateImportInput_InvalidWalletData(t *testing.T) {
	err := ValidateImportInput("", "pass", "abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon about")
	if err == nil {
		t.Error("Expected error for invalid wallet data")
	}
	if !strings.Contains(err.Error(), "wallet name is required") {
		t.Errorf("Expected wallet validation error, got: %s", err.Error())
	}
}

func TestValidateImportInput_EmptyMnemonic(t *testing.T) {
	err := ValidateImportInput("MyWallet", "password123", "")
	if err == nil {
		t.Error("Expected error for empty mnemonic")
	}
	if !strings.Contains(err.Error(), "mnemonic phrase is required") {
		t.Errorf("Expected mnemonic required error, got: %s", err.Error())
	}
}

func TestValidateImportInput_WhitespaceMnemonic(t *testing.T) {
	err := ValidateImportInput("MyWallet", "password123", "   ")
	if err == nil {
		t.Error("Expected error for whitespace-only mnemonic")
	}
	if !strings.Contains(err.Error(), "mnemonic phrase is required") {
		t.Errorf("Expected mnemonic required error, got: %s", err.Error())
	}
}

func TestValidateImportInput_InvalidMnemonicLength(t *testing.T) {
	err := ValidateImportInput("MyWallet", "password123", "abandon abandon abandon")
	if err == nil {
		t.Error("Expected error for invalid mnemonic length")
	}
	if !strings.Contains(err.Error(), "mnemonic phrase must be 12 or 24 words") {
		t.Errorf("Expected mnemonic length error, got: %s", err.Error())
	}
}

func TestValidateImportInput_Valid12Words(t *testing.T) {
	mnemonic := "abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon about"
	err := ValidateImportInput("MyWallet", "password123", mnemonic)
	if err != nil {
		t.Errorf("Expected no error for valid 12-word mnemonic, got: %s", err.Error())
	}
}

func TestValidateImportInput_Valid24Words(t *testing.T) {
	mnemonic := "abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon art"
	err := ValidateImportInput("MyWallet", "password123", mnemonic)
	if err != nil {
		t.Errorf("Expected no error for valid 24-word mnemonic, got: %s", err.Error())
	}
}

func TestValidatePrivateKeyInput_InvalidWalletData(t *testing.T) {
	err := ValidatePrivateKeyInput("", "pass", "abcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890")
	if err == nil {
		t.Error("Expected error for invalid wallet data")
	}
	if !strings.Contains(err.Error(), "wallet name is required") {
		t.Errorf("Expected wallet validation error, got: %s", err.Error())
	}
}

func TestValidatePrivateKeyInput_EmptyPrivateKey(t *testing.T) {
	err := ValidatePrivateKeyInput("MyWallet", "password123", "")
	if err == nil {
		t.Error("Expected error for empty private key")
	}
	if !strings.Contains(err.Error(), "private key is required") {
		t.Errorf("Expected private key required error, got: %s", err.Error())
	}
}

func TestValidatePrivateKeyInput_WhitespacePrivateKey(t *testing.T) {
	err := ValidatePrivateKeyInput("MyWallet", "password123", "   ")
	if err == nil {
		t.Error("Expected error for whitespace-only private key")
	}
	if !strings.Contains(err.Error(), "private key is required") {
		t.Errorf("Expected private key required error, got: %s", err.Error())
	}
}

func TestValidatePrivateKeyInput_With0xPrefix(t *testing.T) {
	err := ValidatePrivateKeyInput("MyWallet", "password123", "0xabcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890")
	if err != nil {
		t.Errorf("Expected no error for valid private key with 0x prefix, got: %s", err.Error())
	}
}

func TestValidatePrivateKeyInput_InvalidLength(t *testing.T) {
	err := ValidatePrivateKeyInput("MyWallet", "password123", "abcdef123456")
	if err == nil {
		t.Error("Expected error for invalid private key length")
	}
	if !strings.Contains(err.Error(), "private key must be 64 hexadecimal characters") {
		t.Errorf("Expected length error, got: %s", err.Error())
	}
}

func TestValidatePrivateKeyInput_InvalidCharacters(t *testing.T) {
	err := ValidatePrivateKeyInput("MyWallet", "password123", "gggggg1234567890abcdef1234567890abcdef1234567890abcdef1234567890")
	if err == nil {
		t.Error("Expected error for invalid hex characters")
	}
	if !strings.Contains(err.Error(), "private key must contain only hexadecimal characters") {
		t.Errorf("Expected hex character error, got: %s", err.Error())
	}
}

func TestValidatePrivateKeyInput_ValidUppercase(t *testing.T) {
	err := ValidatePrivateKeyInput("MyWallet", "password123", "ABCDEF1234567890ABCDEF1234567890ABCDEF1234567890ABCDEF1234567890")
	if err != nil {
		t.Errorf("Expected no error for valid uppercase private key, got: %s", err.Error())
	}
}

func TestValidatePrivateKeyInput_ValidLowercase(t *testing.T) {
	err := ValidatePrivateKeyInput("MyWallet", "password123", "abcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890")
	if err != nil {
		t.Errorf("Expected no error for valid lowercase private key, got: %s", err.Error())
	}
}

func TestValidatePrivateKeyInput_ValidMixedCase(t *testing.T) {
	err := ValidatePrivateKeyInput("MyWallet", "password123", "AbCdEf1234567890aBcDeF1234567890AbCdEf1234567890aBcDeF1234567890")
	if err != nil {
		t.Errorf("Expected no error for valid mixed case private key, got: %s", err.Error())
	}
}

func TestFormatWalletAddress_ShortAddress(t *testing.T) {
	address := "0x12345"
	formatted := FormatWalletAddress(address)
	if formatted != address {
		t.Errorf("Expected short address to remain unchanged, got: %s", formatted)
	}
}

func TestFormatWalletAddress_ExactLengthLimit(t *testing.T) {
	address := "0x12345678901" // Exactly 13 characters
	formatted := FormatWalletAddress(address)
	if formatted != address {
		t.Errorf("Expected address at length limit to remain unchanged, got: %s", formatted)
	}
}

func TestFormatWalletAddress_JustOverLimit(t *testing.T) {
	address := "0x12345678901234" // 14 characters, just over limit
	formatted := FormatWalletAddress(address)
	expected := "0x12345678..."
	if formatted != expected {
		t.Errorf("Expected formatted address to be '%s', got: %s", expected, formatted)
	}
}

func TestFormatWalletAddress_LongAddress(t *testing.T) {
	address := "0x1234567890123456789012345678901234567890"
	formatted := FormatWalletAddress(address)
	expected := "0x12345678..."
	if formatted != expected {
		t.Errorf("Expected formatted address to be '%s', got: %s", expected, formatted)
	}
}

func TestFormatWalletAddress_EmptyAddress(t *testing.T) {
	address := ""
	formatted := FormatWalletAddress(address)
	if formatted != address {
		t.Errorf("Expected empty address to remain unchanged, got: %s", formatted)
	}
}

// Test message types
func TestWalletMessageTypes(t *testing.T) {
	// Test that message types can be instantiated
	_ = WalletLoadedMsg{}
	_ = WalletBalanceLoadedMsg(nil)
	_ = WalletMultiBalanceLoadedMsg(nil)
	_ = WalletCreatedMsg{}
	_ = WalletDeletedMsg{}
	_ = WalletAuthenticatedMsg{}
}
