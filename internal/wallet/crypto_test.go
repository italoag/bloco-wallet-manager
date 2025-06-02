package wallet

import (
	"encoding/base64"
	"strings"
	"testing"
)

func TestArgon2IDParameters(t *testing.T) {
	// Test that our Argon2ID parameters are reasonable
	if argon2IDTime < 1 {
		t.Fatal("Argon2ID time parameter should be at least 1")
	}

	if argon2IDMemory < 64*1024 {
		t.Fatal("Argon2ID memory parameter should be at least 64KB")
	}

	if argon2IDThreads < 1 {
		t.Fatal("Argon2ID threads parameter should be at least 1")
	}

	if argon2IDKeyLen < 32 {
		t.Fatal("Argon2ID key length should be at least 32 bytes")
	}
}

func TestSaltGeneration(t *testing.T) {
	// We can't test the internal generateSalt function directly since it's not exported
	// Instead, we'll test that encryption generates different results each time
	mnemonic := "test mnemonic"
	password := "password"

	encrypted1, err := EncryptMnemonic(mnemonic, password)
	if err != nil {
		t.Fatalf("Failed to encrypt first time: %v", err)
	}

	encrypted2, err := EncryptMnemonic(mnemonic, password)
	if err != nil {
		t.Fatalf("Failed to encrypt second time: %v", err)
	}

	// Different encryptions should produce different results due to random salt
	if encrypted1 == encrypted2 {
		t.Fatal("Encrypted results should be different due to random salt")
	}
}

// TestPasswordSecurity tests various password scenarios
func TestPasswordSecurity(t *testing.T) {
	testCases := []struct {
		name        string
		password    string
		shouldError bool
	}{
		{"Empty password", "", true},
		{"Single character", "a", false},
		{"Normal password", "password123", false},
		{"Long password", strings.Repeat("x", 1000), false},
		{"Unicode password", "пароль123", false},
		{"Special characters", "!@#$%^&*()", false},
		{"Spaces in password", "pass word 123", false},
		{"Only numbers", "123456789", false},
		{"Mixed case", "PassWord123", false},
	}

	mnemonic := "abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon about"

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			encrypted, err := EncryptMnemonic(mnemonic, tc.password)
			if tc.shouldError {
				if err == nil {
					t.Fatal("Expected error but got none")
				}
				return
			}

			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}

			// Test decryption
			decrypted, err := DecryptMnemonic(encrypted, tc.password)
			if err != nil {
				t.Fatalf("Failed to decrypt: %v", err)
			}

			if decrypted != mnemonic {
				t.Fatalf("Decrypted mnemonic doesn't match. Expected %q, got %q", mnemonic, decrypted)
			}
		})
	}
}

// TestMnemonicSecurity tests various mnemonic scenarios
func TestMnemonicSecurity(t *testing.T) {
	testCases := []struct {
		name     string
		mnemonic string
	}{
		{"Empty mnemonic", ""},
		{"Single word", "abandon"},
		{"Standard 12 word", "abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon about"},
		{"Standard 24 word", "abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon art"},
		{"Unicode mnemonic", "тест слово пароль"},
		{"Special characters", "test!@#$ word%^&*"},
		{"Very long mnemonic", strings.Repeat("word ", 100)},
		{"Newlines and tabs", "word1\nword2\tword3"},
		{"Mixed case", "Abandon Abandon ABANDON"},
	}

	password := "testpassword123"

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			encrypted, err := EncryptMnemonic(tc.mnemonic, password)
			if err != nil {
				t.Fatalf("Failed to encrypt: %v", err)
			}

			decrypted, err := DecryptMnemonic(encrypted, password)
			if err != nil {
				t.Fatalf("Failed to decrypt: %v", err)
			}

			if decrypted != tc.mnemonic {
				t.Fatalf("Decrypted mnemonic doesn't match. Expected %q, got %q", tc.mnemonic, decrypted)
			}
		})
	}
}

// TestWrongPasswordDetection tests that wrong passwords are properly detected
func TestWrongPasswordDetection(t *testing.T) {
	mnemonic := "abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon about"
	correctPassword := "correct_password"
	wrongPassword := "wrong_password"

	encrypted, err := EncryptMnemonic(mnemonic, correctPassword)
	if err != nil {
		t.Fatalf("Failed to encrypt: %v", err)
	}

	// Test with wrong password
	_, err = DecryptMnemonic(encrypted, wrongPassword)
	if err == nil {
		t.Fatal("Expected error with wrong password but got none")
	}

	if !strings.Contains(err.Error(), "hash verification failed") {
		t.Fatalf("Expected hash verification error, got: %v", err)
	}

	// Test with correct password should work
	decrypted, err := DecryptMnemonic(encrypted, correctPassword)
	if err != nil {
		t.Fatalf("Failed to decrypt with correct password: %v", err)
	}

	if decrypted != mnemonic {
		t.Fatalf("Decrypted mnemonic doesn't match")
	}
}

// TestEncryptionFormat tests the format of encrypted data
func TestEncryptionFormat(t *testing.T) {
	mnemonic := "test mnemonic"
	password := "password"

	encrypted, err := EncryptMnemonic(mnemonic, password)
	if err != nil {
		t.Fatalf("Failed to encrypt: %v", err)
	}

	// Should be valid base64
	decoded, err := base64.StdEncoding.DecodeString(encrypted)
	if err != nil {
		t.Fatalf("Encrypted data is not valid base64: %v", err)
	}

	// Should have minimum length: salt(16) + hash(32) + at least some encrypted data
	expectedMinLength := saltLength + hashLength
	if len(decoded) < expectedMinLength {
		t.Fatalf("Encrypted data too short. Expected at least %d bytes, got %d", expectedMinLength, len(decoded))
	}

	// Verify structure: first 16 bytes should be salt, next 32 should be hash
	if len(decoded) < saltLength+hashLength {
		t.Fatal("Decoded data doesn't contain salt and hash")
	}
}

// TestConcurrentEncryption tests thread safety
func TestConcurrentEncryption(t *testing.T) {
	mnemonic := "abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon about"
	password := "password123"

	// Run multiple encryptions concurrently
	results := make(chan string, 10)
	errors := make(chan error, 10)

	for i := 0; i < 10; i++ {
		go func() {
			encrypted, err := EncryptMnemonic(mnemonic, password)
			if err != nil {
				errors <- err
				return
			}
			results <- encrypted
		}()
	}

	// Collect results
	var encryptedResults []string
	for i := 0; i < 10; i++ {
		select {
		case encrypted := <-results:
			encryptedResults = append(encryptedResults, encrypted)
		case err := <-errors:
			t.Fatalf("Concurrent encryption failed: %v", err)
		}
	}

	// All results should be different (due to random salt)
	seen := make(map[string]bool)
	for _, encrypted := range encryptedResults {
		if seen[encrypted] {
			t.Fatal("Found duplicate encrypted result - salt randomization failed")
		}
		seen[encrypted] = true

		// Each should decrypt correctly
		decrypted, err := DecryptMnemonic(encrypted, password)
		if err != nil {
			t.Fatalf("Failed to decrypt concurrent result: %v", err)
		}
		if decrypted != mnemonic {
			t.Fatal("Concurrent decryption produced wrong result")
		}
	}
}

// TestLargeDataHandling tests encryption with very large mnemonics
func TestLargeDataHandling(t *testing.T) {
	// Create a very large mnemonic
	largeMnemonic := strings.Repeat("word ", 10000) // ~50KB of data
	password := "password"

	encrypted, err := EncryptMnemonic(largeMnemonic, password)
	if err != nil {
		t.Fatalf("Failed to encrypt large mnemonic: %v", err)
	}

	decrypted, err := DecryptMnemonic(encrypted, password)
	if err != nil {
		t.Fatalf("Failed to decrypt large mnemonic: %v", err)
	}

	if decrypted != largeMnemonic {
		t.Fatal("Large mnemonic decryption failed")
	}
}

// TestSecureCompare tests the constant-time comparison function
func TestSecureCompare(t *testing.T) {
	testCases := []struct {
		name     string
		a, b     string
		expected bool
	}{
		{"Equal strings", "test", "test", true},
		{"Different strings", "test", "best", false},
		{"Empty strings", "", "", true},
		{"One empty", "test", "", false},
		{"Different lengths", "short", "verylongstring", false},
		{"Unicode equal", "тест", "тест", true},
		{"Unicode different", "тест", "лест", false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := SecureCompare(tc.a, tc.b)
			if result != tc.expected {
				t.Fatalf("Expected %v, got %v for comparing %q and %q", tc.expected, result, tc.a, tc.b)
			}
		})
	}
}

// TestVerifyMnemonicPassword tests password verification
func TestVerifyMnemonicPassword(t *testing.T) {
	mnemonic := "abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon about"
	correctPassword := "correct123"
	wrongPassword := "wrong123"

	encrypted, err := EncryptMnemonic(mnemonic, correctPassword)
	if err != nil {
		t.Fatalf("Failed to encrypt: %v", err)
	}

	// Test correct password
	if !VerifyMnemonicPassword(encrypted, correctPassword) {
		t.Fatal("VerifyMnemonicPassword should return true for correct password")
	}

	// Test wrong password
	if VerifyMnemonicPassword(encrypted, wrongPassword) {
		t.Fatal("VerifyMnemonicPassword should return false for wrong password")
	}

	// Test with invalid encrypted data
	if VerifyMnemonicPassword("invalid", correctPassword) {
		t.Fatal("VerifyMnemonicPassword should return false for invalid encrypted data")
	}
}

func TestEncryptDecryptCycle(t *testing.T) {
	testCases := []struct {
		name     string
		mnemonic string
		password string
	}{
		{
			name:     "Standard mnemonic",
			mnemonic: "abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon about",
			password: "testpassword123",
		},
		{
			name:     "Short password",
			mnemonic: "abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon about",
			password: "123",
		},
		{
			name:     "Long password",
			mnemonic: "abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon about",
			password: "this_is_a_very_long_password_with_special_characters_!@#$%^&*()_+-=[]{}|;:,.<>?",
		},
		{
			name:     "Unicode password",
			mnemonic: "abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon about",
			password: "пароль_with_unicode_密码_🔐",
		},
		{
			name:     "Empty mnemonic",
			mnemonic: "",
			password: "testpassword",
		},
		{
			name:     "Special characters mnemonic",
			mnemonic: "test mnemonic with special chars: !@#$%^&*()",
			password: "password",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Encrypt
			encrypted, err := EncryptMnemonic(tc.mnemonic, tc.password)
			if err != nil {
				t.Fatalf("Failed to encrypt mnemonic: %v", err)
			}

			// Verify encrypted data is not empty and different from original
			if encrypted == "" {
				t.Fatal("Encrypted mnemonic is empty")
			}

			if encrypted == tc.mnemonic && tc.mnemonic != "" {
				t.Fatal("Encrypted mnemonic should not be the same as original (unless original is empty)")
			}

			// Decrypt with correct password
			decrypted, err := DecryptMnemonic(encrypted, tc.password)
			if err != nil {
				t.Fatalf("Failed to decrypt mnemonic: %v", err)
			}

			if decrypted != tc.mnemonic {
				t.Fatalf("Decrypted mnemonic doesn't match original. Got: %s, Expected: %s", decrypted, tc.mnemonic)
			}

			// Try to decrypt with wrong password
			_, err = DecryptMnemonic(encrypted, tc.password+"wrong")
			if err == nil {
				t.Fatal("Expected error when decrypting with wrong password")
			}
		})
	}
}

func TestDecryptInvalidData(t *testing.T) {
	testCases := []struct {
		name string
		data string
	}{
		{
			name: "Empty string",
			data: "",
		},
		{
			name: "Invalid base64",
			data: "invalid_base64_data!!!",
		},
		{
			name: "Valid base64 but invalid format",
			data: "dGVzdA==", // "test" in base64
		},
		{
			name: "Too short data",
			data: "dGVzdA==",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := DecryptMnemonic(tc.data, "anypassword")
			if err == nil {
				t.Fatal("Expected error when decrypting invalid data")
			}
		})
	}
}

func TestEncryptionConsistency(t *testing.T) {
	mnemonic := "abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon about"
	password := "testpassword123"

	// Encrypt the same data multiple times
	encrypted1, err := EncryptMnemonic(mnemonic, password)
	if err != nil {
		t.Fatalf("Failed to encrypt mnemonic first time: %v", err)
	}

	encrypted2, err := EncryptMnemonic(mnemonic, password)
	if err != nil {
		t.Fatalf("Failed to encrypt mnemonic second time: %v", err)
	}

	// Encrypted data should be different each time due to random salt
	if encrypted1 == encrypted2 {
		t.Fatal("Encrypted data should be different each time due to random salt")
	}

	// But both should decrypt to the same original data
	decrypted1, err := DecryptMnemonic(encrypted1, password)
	if err != nil {
		t.Fatalf("Failed to decrypt first encrypted data: %v", err)
	}

	decrypted2, err := DecryptMnemonic(encrypted2, password)
	if err != nil {
		t.Fatalf("Failed to decrypt second encrypted data: %v", err)
	}

	if decrypted1 != mnemonic || decrypted2 != mnemonic {
		t.Fatal("Both decrypted values should match the original mnemonic")
	}

	if decrypted1 != decrypted2 {
		t.Fatal("Both decrypted values should be identical")
	}
}
