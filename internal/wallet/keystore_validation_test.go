package wallet

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateKeystoreV3(t *testing.T) {
	validator := &KeystoreValidator{}

	tests := []struct {
		name          string
		json          string
		expectedError KeystoreErrorType
	}{
		{
			name: "Valid keystore",
			json: `{
				"version": 3,
				"id": "f06e0f8e-7d91-4b09-8f5a-3c2c1a9b2b88",
				"address": "0x5d8c5d3a5e6f6d6c5b4a3a2b1c0d9e8f7a6b5c4d",
				"crypto": {
					"cipher": "aes-128-ctr",
					"ciphertext": "5d8c5d3a5e6f6d6c5b4a3a2b1c0d9e8f7a6b5c4d3e2f1g0h9i8j7k6l5m4n3o2p1",
					"cipherparams": {
						"iv": "5d8c5d3a5e6f6d6c5b4a3a2b1c0d9e8f"
					},
					"kdf": "scrypt",
					"kdfparams": {
						"dklen": 32,
						"n": 262144,
						"p": 1,
						"r": 8,
						"salt": "5d8c5d3a5e6f6d6c5b4a3a2b1c0d9e8f7a6b5c4d3e2f1g0h9i8j7k6l5m4n3o2p1"
					},
					"mac": "5d8c5d3a5e6f6d6c5b4a3a2b1c0d9e8f7a6b5c4d3e2f1g0h9i8j7k6l5m4n3o2p1"
				}
			}`,
			expectedError: -1, // No error expected
		}, {

			name: "Invalid JSON",
			json: `{
				"version": 3,
				"id": "f06e0f8e-7d91-4b09-8f5a-3c2c1a9b2b88",
				"address": "0x5d8c5d3a5e6f6d6c5b4a3a2b1c0d9e8f7a6b5c4d",
				"crypto": {
					"cipher": "aes-128-ctr",
					"ciphertext": "5d8c5d3a5e6f6d6c5b4a3a2b1c0d9e8f7a6b5c4d3e2f1g0h9i8j7k6l5m4n3o2p1",
					"cipherparams": {
						"iv": "5d8c5d3a5e6f6d6c5b4a3a2b1c0d9e8f"
					},
					"kdf": "scrypt",
					"kdfparams": {
						"dklen": 32,
						"n": 262144,
						"p": 1,
						"r": 8,
						"salt": "5d8c5d3a5e6f6d6c5b4a3a2b1c0d9e8f7a6b5c4d3e2f1g0h9i8j7k6l5m4n3o2p1"
					},
					"mac": "5d8c5d3a5e6f6d6c5b4a3a2b1c0d9e8f7a6b5c4d3e2f1g0h9i8j7k6l5m4n3o2p1"
				}
			`,
			expectedError: ErrorInvalidJSON,
		},
		{
			name: "Invalid version",
			json: `{
				"version": 2,
				"id": "f06e0f8e-7d91-4b09-8f5a-3c2c1a9b2b88",
				"address": "0x5d8c5d3a5e6f6d6c5b4a3a2b1c0d9e8f7a6b5c4d",
				"crypto": {
					"cipher": "aes-128-ctr",
					"ciphertext": "5d8c5d3a5e6f6d6c5b4a3a2b1c0d9e8f7a6b5c4d3e2f1g0h9i8j7k6l5m4n3o2p1",
					"cipherparams": {
						"iv": "5d8c5d3a5e6f6d6c5b4a3a2b1c0d9e8f"
					},
					"kdf": "scrypt",
					"kdfparams": {
						"dklen": 32,
						"n": 262144,
						"p": 1,
						"r": 8,
						"salt": "5d8c5d3a5e6f6d6c5b4a3a2b1c0d9e8f7a6b5c4d3e2f1g0h9i8j7k6l5m4n3o2p1"
					},
					"mac": "5d8c5d3a5e6f6d6c5b4a3a2b1c0d9e8f7a6b5c4d3e2f1g0h9i8j7k6l5m4n3o2p1"
				}
			}`,
			expectedError: ErrorInvalidVersion,
		}, {

			name: "Missing address",
			json: `{
				"version": 3,
				"id": "f06e0f8e-7d91-4b09-8f5a-3c2c1a9b2b88",
				"crypto": {
					"cipher": "aes-128-ctr",
					"ciphertext": "5d8c5d3a5e6f6d6c5b4a3a2b1c0d9e8f7a6b5c4d3e2f1g0h9i8j7k6l5m4n3o2p1",
					"cipherparams": {
						"iv": "5d8c5d3a5e6f6d6c5b4a3a2b1c0d9e8f"
					},
					"kdf": "scrypt",
					"kdfparams": {
						"dklen": 32,
						"n": 262144,
						"p": 1,
						"r": 8,
						"salt": "5d8c5d3a5e6f6d6c5b4a3a2b1c0d9e8f7a6b5c4d3e2f1g0h9i8j7k6l5m4n3o2p1"
					},
					"mac": "5d8c5d3a5e6f6d6c5b4a3a2b1c0d9e8f7a6b5c4d3e2f1g0h9i8j7k6l5m4n3o2p1"
				}
			}`,
			expectedError: ErrorMissingRequiredFields,
		},
		{
			name: "Invalid address format",
			json: `{
				"version": 3,
				"id": "f06e0f8e-7d91-4b09-8f5a-3c2c1a9b2b88",
				"address": "not-a-valid-address",
				"crypto": {
					"cipher": "aes-128-ctr",
					"ciphertext": "5d8c5d3a5e6f6d6c5b4a3a2b1c0d9e8f7a6b5c4d3e2f1g0h9i8j7k6l5m4n3o2p1",
					"cipherparams": {
						"iv": "5d8c5d3a5e6f6d6c5b4a3a2b1c0d9e8f"
					},
					"kdf": "scrypt",
					"kdfparams": {
						"dklen": 32,
						"n": 262144,
						"p": 1,
						"r": 8,
						"salt": "5d8c5d3a5e6f6d6c5b4a3a2b1c0d9e8f7a6b5c4d3e2f1g0h9i8j7k6l5m4n3o2p1"
					},
					"mac": "5d8c5d3a5e6f6d6c5b4a3a2b1c0d9e8f7a6b5c4d3e2f1g0h9i8j7k6l5m4n3o2p1"
				}
			}`,
			expectedError: ErrorInvalidAddress,
		}, {

			name: "Missing crypto.cipher field",
			json: `{
				"version": 3,
				"id": "f06e0f8e-7d91-4b09-8f5a-3c2c1a9b2b88",
				"address": "0x5d8c5d3a5e6f6d6c5b4a3a2b1c0d9e8f7a6b5c4d",
				"crypto": {
					"ciphertext": "5d8c5d3a5e6f6d6c5b4a3a2b1c0d9e8f7a6b5c4d3e2f1g0h9i8j7k6l5m4n3o2p1",
					"cipherparams": {
						"iv": "5d8c5d3a5e6f6d6c5b4a3a2b1c0d9e8f"
					},
					"kdf": "scrypt",
					"kdfparams": {
						"dklen": 32,
						"n": 262144,
						"p": 1,
						"r": 8,
						"salt": "5d8c5d3a5e6f6d6c5b4a3a2b1c0d9e8f7a6b5c4d3e2f1g0h9i8j7k6l5m4n3o2p1"
					},
					"mac": "5d8c5d3a5e6f6d6c5b4a3a2b1c0d9e8f7a6b5c4d3e2f1g0h9i8j7k6l5m4n3o2p1"
				}
			}`,
			expectedError: ErrorMissingRequiredFields,
		},
		{
			name: "Missing crypto.ciphertext field",
			json: `{
				"version": 3,
				"id": "f06e0f8e-7d91-4b09-8f5a-3c2c1a9b2b88",
				"address": "0x5d8c5d3a5e6f6d6c5b4a3a2b1c0d9e8f7a6b5c4d",
				"crypto": {
					"cipher": "aes-128-ctr",
					"cipherparams": {
						"iv": "5d8c5d3a5e6f6d6c5b4a3a2b1c0d9e8f"
					},
					"kdf": "scrypt",
					"kdfparams": {
						"dklen": 32,
						"n": 262144,
						"p": 1,
						"r": 8,
						"salt": "5d8c5d3a5e6f6d6c5b4a3a2b1c0d9e8f7a6b5c4d3e2f1g0h9i8j7k6l5m4n3o2p1"
					},
					"mac": "5d8c5d3a5e6f6d6c5b4a3a2b1c0d9e8f7a6b5c4d3e2f1g0h9i8j7k6l5m4n3o2p1"
				}
			}`,
			expectedError: ErrorMissingRequiredFields,
		},
		{
			name: "Missing crypto.cipherparams.iv field",
			json: `{
				"version": 3,
				"id": "f06e0f8e-7d91-4b09-8f5a-3c2c1a9b2b88",
				"address": "0x5d8c5d3a5e6f6d6c5b4a3a2b1c0d9e8f7a6b5c4d",
				"crypto": {
					"cipher": "aes-128-ctr",
					"ciphertext": "5d8c5d3a5e6f6d6c5b4a3a2b1c0d9e8f7a6b5c4d3e2f1g0h9i8j7k6l5m4n3o2p1",
					"cipherparams": {},
					"kdf": "scrypt",
					"kdfparams": {
						"dklen": 32,
						"n": 262144,
						"p": 1,
						"r": 8,
						"salt": "5d8c5d3a5e6f6d6c5b4a3a2b1c0d9e8f7a6b5c4d3e2f1g0h9i8j7k6l5m4n3o2p1"
					},
					"mac": "5d8c5d3a5e6f6d6c5b4a3a2b1c0d9e8f7a6b5c4d3e2f1g0h9i8j7k6l5m4n3o2p1"
				}
			}`,
			expectedError: ErrorMissingRequiredFields,
		},
		{
			name: "Missing crypto.kdf field",
			json: `{
				"version": 3,
				"id": "f06e0f8e-7d91-4b09-8f5a-3c2c1a9b2b88",
				"address": "0x5d8c5d3a5e6f6d6c5b4a3a2b1c0d9e8f7a6b5c4d",
				"crypto": {
					"cipher": "aes-128-ctr",
					"ciphertext": "5d8c5d3a5e6f6d6c5b4a3a2b1c0d9e8f7a6b5c4d3e2f1g0h9i8j7k6l5m4n3o2p1",
					"cipherparams": {
						"iv": "5d8c5d3a5e6f6d6c5b4a3a2b1c0d9e8f"
					},
					"kdfparams": {
						"dklen": 32,
						"n": 262144,
						"p": 1,
						"r": 8,
						"salt": "5d8c5d3a5e6f6d6c5b4a3a2b1c0d9e8f7a6b5c4d3e2f1g0h9i8j7k6l5m4n3o2p1"
					},
					"mac": "5d8c5d3a5e6f6d6c5b4a3a2b1c0d9e8f7a6b5c4d3e2f1g0h9i8j7k6l5m4n3o2p1"
				}
			}`,
			expectedError: ErrorMissingRequiredFields,
		},
		{
			name: "Missing crypto.kdfparams field",
			json: `{
				"version": 3,
				"id": "f06e0f8e-7d91-4b09-8f5a-3c2c1a9b2b88",
				"address": "0x5d8c5d3a5e6f6d6c5b4a3a2b1c0d9e8f7a6b5c4d",
				"crypto": {
					"cipher": "aes-128-ctr",
					"ciphertext": "5d8c5d3a5e6f6d6c5b4a3a2b1c0d9e8f7a6b5c4d3e2f1g0h9i8j7k6l5m4n3o2p1",
					"cipherparams": {
						"iv": "5d8c5d3a5e6f6d6c5b4a3a2b1c0d9e8f"
					},
					"kdf": "scrypt",
					"mac": "5d8c5d3a5e6f6d6c5b4a3a2b1c0d9e8f7a6b5c4d3e2f1g0h9i8j7k6l5m4n3o2p1"
				}
			}`,
			expectedError: ErrorMissingRequiredFields,
		},
		{
			name: "Missing crypto.mac field",
			json: `{
				"version": 3,
				"id": "f06e0f8e-7d91-4b09-8f5a-3c2c1a9b2b88",
				"address": "0x5d8c5d3a5e6f6d6c5b4a3a2b1c0d9e8f7a6b5c4d",
				"crypto": {
					"cipher": "aes-128-ctr",
					"ciphertext": "5d8c5d3a5e6f6d6c5b4a3a2b1c0d9e8f7a6b5c4d3e2f1g0h9i8j7k6l5m4n3o2p1",
					"cipherparams": {
						"iv": "5d8c5d3a5e6f6d6c5b4a3a2b1c0d9e8f"
					},
					"kdf": "scrypt",
					"kdfparams": {
						"dklen": 32,
						"n": 262144,
						"p": 1,
						"r": 8,
						"salt": "5d8c5d3a5e6f6d6c5b4a3a2b1c0d9e8f7a6b5c4d3e2f1g0h9i8j7k6l5m4n3o2p1"
					}
				}
			}`,
			expectedError: ErrorMissingRequiredFields,
		}, {

			name: "PBKDF2 KDF",
			json: `{
				"version": 3,
				"id": "f06e0f8e-7d91-4b09-8f5a-3c2c1a9b2b88",
				"address": "0x5d8c5d3a5e6f6d6c5b4a3a2b1c0d9e8f7a6b5c4d",
				"crypto": {
					"cipher": "aes-128-ctr",
					"ciphertext": "5d8c5d3a5e6f6d6c5b4a3a2b1c0d9e8f7a6b5c4d3e2f1g0h9i8j7k6l5m4n3o2p1",
					"cipherparams": {
						"iv": "5d8c5d3a5e6f6d6c5b4a3a2b1c0d9e8f"
					},
					"kdf": "pbkdf2",
					"kdfparams": {
						"dklen": 32,
						"c": 10240,
						"prf": "hmac-sha256",
						"salt": "5d8c5d3a5e6f6d6c5b4a3a2b1c0d9e8f7a6b5c4d3e2f1g0h9i8j7k6l5m4n3o2p1"
					},
					"mac": "5d8c5d3a5e6f6d6c5b4a3a2b1c0d9e8f7a6b5c4d3e2f1g0h9i8j7k6l5m4n3o2p1"
				}
			}`,
			expectedError: -1, // No error expected
		},
		{
			name: "Missing PBKDF2 required field - dklen",
			json: `{
				"version": 3,
				"id": "f06e0f8e-7d91-4b09-8f5a-3c2c1a9b2b88",
				"address": "0x5d8c5d3a5e6f6d6c5b4a3a2b1c0d9e8f7a6b5c4d",
				"crypto": {
					"cipher": "aes-128-ctr",
					"ciphertext": "5d8c5d3a5e6f6d6c5b4a3a2b1c0d9e8f7a6b5c4d3e2f1g0h9i8j7k6l5m4n3o2p1",
					"cipherparams": {
						"iv": "5d8c5d3a5e6f6d6c5b4a3a2b1c0d9e8f"
					},
					"kdf": "pbkdf2",
					"kdfparams": {
						"c": 10240,
						"prf": "hmac-sha256",
						"salt": "5d8c5d3a5e6f6d6c5b4a3a2b1c0d9e8f7a6b5c4d3e2f1g0h9i8j7k6l5m4n3o2p1"
					},
					"mac": "5d8c5d3a5e6f6d6c5b4a3a2b1c0d9e8f7a6b5c4d3e2f1g0h9i8j7k6l5m4n3o2p1"
				}
			}`,
			expectedError: ErrorMissingRequiredFields,
		}, {

			name: "Missing Scrypt required field - dklen",
			json: `{
				"version": 3,
				"id": "f06e0f8e-7d91-4b09-8f5a-3c2c1a9b2b88",
				"address": "0x5d8c5d3a5e6f6d6c5b4a3a2b1c0d9e8f7a6b5c4d",
				"crypto": {
					"cipher": "aes-128-ctr",
					"ciphertext": "5d8c5d3a5e6f6d6c5b4a3a2b1c0d9e8f7a6b5c4d3e2f1g0h9i8j7k6l5m4n3o2p1",
					"cipherparams": {
						"iv": "5d8c5d3a5e6f6d6c5b4a3a2b1c0d9e8f"
					},
					"kdf": "scrypt",
					"kdfparams": {
						"n": 262144,
						"p": 1,
						"r": 8,
						"salt": "5d8c5d3a5e6f6d6c5b4a3a2b1c0d9e8f7a6b5c4d3e2f1g0h9i8j7k6l5m4n3o2p1"
					},
					"mac": "5d8c5d3a5e6f6d6c5b4a3a2b1c0d9e8f7a6b5c4d3e2f1g0h9i8j7k6l5m4n3o2p1"
				}
			}`,
			expectedError: ErrorMissingRequiredFields,
		},
		{
			name: "Unsupported KDF",
			json: `{
				"version": 3,
				"id": "f06e0f8e-7d91-4b09-8f5a-3c2c1a9b2b88",
				"address": "0x5d8c5d3a5e6f6d6c5b4a3a2b1c0d9e8f7a6b5c4d",
				"crypto": {
					"cipher": "aes-128-ctr",
					"ciphertext": "5d8c5d3a5e6f6d6c5b4a3a2b1c0d9e8f7a6b5c4d3e2f1g0h9i8j7k6l5m4n3o2p1",
					"cipherparams": {
						"iv": "5d8c5d3a5e6f6d6c5b4a3a2b1c0d9e8f"
					},
					"kdf": "unsupported-kdf",
					"kdfparams": {
						"dklen": 32,
						"n": 262144,
						"p": 1,
						"r": 8,
						"salt": "5d8c5d3a5e6f6d6c5b4a3a2b1c0d9e8f7a6b5c4d3e2f1g0h9i8j7k6l5m4n3o2p1"
					},
					"mac": "5d8c5d3a5e6f6d6c5b4a3a2b1c0d9e8f7a6b5c4d3e2f1g0h9i8j7k6l5m4n3o2p1"
				}
			}`,
			expectedError: ErrorInvalidKeystore,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := validator.ValidateKeystoreV3([]byte(tt.json))

			if tt.expectedError == -1 {
				if err != nil {
					t.Errorf("Expected no error, got %v", err)
				}
				return
			}

			if err == nil {
				t.Errorf("Expected error of type %v, got nil", tt.expectedError)
				return
			}

			keystoreErr, ok := err.(*KeystoreImportError)
			if !ok {
				t.Errorf("Expected KeystoreImportError, got %T", err)
				return
			}

			if keystoreErr.Type != tt.expectedError {
				t.Errorf("Expected error type %v, got %v", tt.expectedError, keystoreErr.Type)
			}
		})
	}
}

func TestValidateAddress(t *testing.T) {
	validator := &KeystoreValidator{}

	tests := []struct {
		name          string
		address       string
		expectedError KeystoreErrorType
	}{
		{
			name:          "Valid address with 0x prefix",
			address:       "0x5d8c5d3a5e6f6d6c5b4a3a2b1c0d9e8f7a6b5c4d",
			expectedError: -1, // No error expected
		},
		{
			name:          "Valid address without 0x prefix",
			address:       "5d8c5d3a5e6f6d6c5b4a3a2b1c0d9e8f7a6b5c4d",
			expectedError: -1, // No error expected
		},
		{
			name:          "Empty address",
			address:       "",
			expectedError: ErrorMissingRequiredFields,
		},
		{
			name:          "Invalid address - too short",
			address:       "0x5d8c5d3a5e6f6d6c5b4a3a2b1c0d9e8f",
			expectedError: ErrorInvalidAddress,
		},
		{
			name:          "Invalid address - too long",
			address:       "0x5d8c5d3a5e6f6d6c5b4a3a2b1c0d9e8f7a6b5c4d3e2f",
			expectedError: ErrorInvalidAddress,
		},
		{
			name:          "Invalid address - non-hex characters",
			address:       "0x5d8c5d3a5e6f6d6c5b4a3a2b1c0d9e8f7a6b5c4z",
			expectedError: ErrorInvalidAddress,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator.ValidateAddress(tt.address)

			if tt.expectedError == -1 {
				if err != nil {
					t.Errorf("Expected no error, got %v", err)
				}
				return
			}

			if err == nil {
				t.Errorf("Expected error of type %v, got nil", tt.expectedError)
				return
			}

			keystoreErr, ok := err.(*KeystoreImportError)
			if !ok {
				t.Errorf("Expected KeystoreImportError, got %T", err)
				return
			}

			if keystoreErr.Type != tt.expectedError {
				t.Errorf("Expected error type %v, got %v", tt.expectedError, keystoreErr.Type)
			}
		})
	}
}

// TestKeystoreImportErrorMethods tests the methods of KeystoreImportError
func TestKeystoreImportErrorMethods(t *testing.T) {
	// Test Error() method
	err1 := NewKeystoreImportError(ErrorInvalidJSON, "Invalid JSON format", nil)
	assert.Equal(t, "Invalid JSON format", err1.Error())

	err2 := NewKeystoreImportError(ErrorInvalidJSON, "Invalid JSON format", assert.AnError)
	assert.Equal(t, "Invalid JSON format: assert.AnError general error for testing", err2.Error())

	// Test Unwrap() method
	assert.Nil(t, err1.Unwrap())
	assert.Equal(t, assert.AnError, err2.Unwrap())

	// Test GetLocalizedMessage() method
	assert.Equal(t, "keystore_invalid_json", err1.GetLocalizedMessage())

	// Test GetLocalizedMessageWithField() method
	err3 := NewKeystoreImportErrorWithField(ErrorMissingRequiredFields, "Missing field", "address", nil)
	assert.Equal(t, "keystore_missing_fields:address", err3.GetLocalizedMessageWithField())
}

// TestKeystoreErrorTypeString tests the String() method of KeystoreErrorType
func TestKeystoreErrorTypeString(t *testing.T) {
	tests := []struct {
		errorType KeystoreErrorType
		expected  string
	}{
		{ErrorFileNotFound, "FILE_NOT_FOUND"},
		{ErrorInvalidJSON, "INVALID_JSON"},
		{ErrorInvalidKeystore, "INVALID_KEYSTORE"},
		{ErrorInvalidVersion, "INVALID_VERSION"},
		{ErrorIncorrectPassword, "INCORRECT_PASSWORD"},
		{ErrorCorruptedFile, "CORRUPTED_FILE"},
		{ErrorAddressMismatch, "ADDRESS_MISMATCH"},
		{ErrorMissingRequiredFields, "MISSING_REQUIRED_FIELDS"},
		{ErrorInvalidAddress, "INVALID_ADDRESS"},
		{KeystoreErrorType(999), "UNKNOWN_ERROR"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.errorType.String())
		})
	}
}

// TestKeystoreErrorTypeGetLocalizationKey tests the GetLocalizationKey() method of KeystoreErrorType
func TestKeystoreErrorTypeGetLocalizationKey(t *testing.T) {
	tests := []struct {
		errorType KeystoreErrorType
		expected  string
	}{
		{ErrorFileNotFound, "keystore_file_not_found"},
		{ErrorInvalidJSON, "keystore_invalid_json"},
		{ErrorInvalidKeystore, "keystore_invalid_structure"},
		{ErrorInvalidVersion, "keystore_invalid_version"},
		{ErrorIncorrectPassword, "keystore_incorrect_password"},
		{ErrorCorruptedFile, "keystore_corrupted_file"},
		{ErrorAddressMismatch, "keystore_address_mismatch"},
		{ErrorMissingRequiredFields, "keystore_missing_fields"},
		{ErrorInvalidAddress, "keystore_invalid_address"},
		{KeystoreErrorType(999), "unknown_error"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.errorType.GetLocalizationKey())
		})
	}
}
