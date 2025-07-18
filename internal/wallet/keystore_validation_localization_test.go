package wallet

import (
	"testing"
)

func TestKeystoreErrorTypeLocalizationKey(t *testing.T) {
	tests := []struct {
		name      string
		errorType KeystoreErrorType
		expected  string
	}{
		{
			name:      "File not found error",
			errorType: ErrorFileNotFound,
			expected:  "keystore_file_not_found",
		},
		{
			name:      "Invalid JSON error",
			errorType: ErrorInvalidJSON,
			expected:  "keystore_invalid_json",
		},
		{
			name:      "Invalid keystore structure error",
			errorType: ErrorInvalidKeystore,
			expected:  "keystore_invalid_structure",
		},
		{
			name:      "Invalid version error",
			errorType: ErrorInvalidVersion,
			expected:  "keystore_invalid_version",
		},
		{
			name:      "Incorrect password error",
			errorType: ErrorIncorrectPassword,
			expected:  "keystore_incorrect_password",
		},
		{
			name:      "Corrupted file error",
			errorType: ErrorCorruptedFile,
			expected:  "keystore_corrupted_file",
		},
		{
			name:      "Address mismatch error",
			errorType: ErrorAddressMismatch,
			expected:  "keystore_address_mismatch",
		},
		{
			name:      "Missing required fields error",
			errorType: ErrorMissingRequiredFields,
			expected:  "keystore_missing_fields",
		},
		{
			name:      "Invalid address error",
			errorType: ErrorInvalidAddress,
			expected:  "keystore_invalid_address",
		},
		{
			name:      "Unknown error type",
			errorType: 999, // Invalid error type
			expected:  "unknown_error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.errorType.GetLocalizationKey()
			if result != tt.expected {
				t.Errorf("GetLocalizationKey() = %q, want %q", result, tt.expected)
			}
		})
	}
}

func TestKeystoreImportErrorGetLocalizedMessage(t *testing.T) {
	tests := []struct {
		name     string
		err      *KeystoreImportError
		expected string
	}{
		{
			name: "File not found error",
			err: &KeystoreImportError{
				Type:    ErrorFileNotFound,
				Message: "File not found",
			},
			expected: "keystore_file_not_found",
		},
		{
			name: "Invalid JSON error",
			err: &KeystoreImportError{
				Type:    ErrorInvalidJSON,
				Message: "Invalid JSON",
			},
			expected: "keystore_invalid_json",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.err.GetLocalizedMessage()
			if result != tt.expected {
				t.Errorf("GetLocalizedMessage() = %q, want %q", result, tt.expected)
			}
		})
	}
}

func TestKeystoreImportErrorGetLocalizedMessageWithField(t *testing.T) {
	tests := []struct {
		name     string
		err      *KeystoreImportError
		expected string
	}{
		{
			name: "Error without field",
			err: &KeystoreImportError{
				Type:    ErrorFileNotFound,
				Message: "File not found",
			},
			expected: "keystore_file_not_found",
		},
		{
			name: "Error with field",
			err: &KeystoreImportError{
				Type:    ErrorMissingRequiredFields,
				Message: "Missing required field",
				Field:   "crypto.cipher",
			},
			expected: "keystore_missing_fields:crypto.cipher",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.err.GetLocalizedMessageWithField()
			if result != tt.expected {
				t.Errorf("GetLocalizedMessageWithField() = %q, want %q", result, tt.expected)
			}
		})
	}
}
