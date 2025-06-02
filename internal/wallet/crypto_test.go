package wallet

import (
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
