package wallet

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
)

// KeystoreErrorType represents different types of keystore import errors
type KeystoreErrorType int

const (
	ErrorFileNotFound KeystoreErrorType = iota
	ErrorInvalidJSON
	ErrorInvalidKeystore
	ErrorInvalidVersion
	ErrorIncorrectPassword
	ErrorCorruptedFile
	ErrorAddressMismatch
	ErrorMissingRequiredFields
	ErrorInvalidAddress
)

// KeystoreImportError represents a specific error that occurred during keystore import
type KeystoreImportError struct {
	Type    KeystoreErrorType
	Message string
	Cause   error
	Field   string // Optional field name that caused the error
}

// Error implements the error interface
func (e *KeystoreImportError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Cause)
	}
	return e.Message
}

// Unwrap returns the underlying error for error unwrapping
func (e *KeystoreImportError) Unwrap() error {
	return e.Cause
}

// KeystoreV3 represents the structure of a keystore v3 file
type KeystoreV3 struct {
	Version int              `json:"version"`
	ID      string           `json:"id"`
	Address string           `json:"address"`
	Crypto  KeystoreV3Crypto `json:"crypto"`
}

// KeystoreV3Crypto represents the crypto section of a keystore v3 file
type KeystoreV3Crypto struct {
	Cipher       string                 `json:"cipher"`
	CipherText   string                 `json:"ciphertext"`
	CipherParams KeystoreV3CipherParams `json:"cipherparams"`
	KDF          string                 `json:"kdf"`
	KDFParams    any                    `json:"kdfparams"`
	MAC          string                 `json:"mac"`
}

// KeystoreV3CipherParams represents the cipher parameters in a keystore v3 file
type KeystoreV3CipherParams struct {
	IV string `json:"iv"`
}

// KeystoreV3ScryptParams represents scrypt KDF parameters
type KeystoreV3ScryptParams struct {
	DKLen int    `json:"dklen"`
	N     int    `json:"n"`
	P     int    `json:"p"`
	R     int    `json:"r"`
	Salt  string `json:"salt"`
}

// KeystoreV3PBKDF2Params represents PBKDF2 KDF parameters
type KeystoreV3PBKDF2Params struct {
	DKLen int    `json:"dklen"`
	C     int    `json:"c"`
	PRF   string `json:"prf"`
	Salt  string `json:"salt"`
}

// KeystoreValidator provides methods to validate keystore files
type KeystoreValidator struct{}

// ValidateKeystoreV3 parses JSON data and validates the keystore structure
func (kv *KeystoreValidator) ValidateKeystoreV3(data []byte) (*KeystoreV3, error) {
	var keystore KeystoreV3

	// Parse JSON
	if err := json.Unmarshal(data, &keystore); err != nil {
		return nil, NewKeystoreImportError(ErrorInvalidJSON, "Invalid JSON format in keystore file", err)
	}

	// Validate structure
	if err := kv.ValidateStructure(&keystore); err != nil {
		return nil, err
	}

	return &keystore, nil
}

// ValidateStructure checks if the keystore has all required fields
func (kv *KeystoreValidator) ValidateStructure(keystore *KeystoreV3) error {
	// Check version
	if err := kv.ValidateVersion(keystore.Version); err != nil {
		return err
	}

	// Check address
	if err := kv.ValidateAddress(keystore.Address); err != nil {
		return err
	}

	// Check crypto section
	if err := kv.ValidateCrypto(&keystore.Crypto); err != nil {
		return err
	}

	return nil
}

// ValidateVersion ensures the keystore version is exactly 3
func (kv *KeystoreValidator) ValidateVersion(version int) error {
	if version != 3 {
		return NewKeystoreImportErrorWithField(
			ErrorInvalidVersion,
			fmt.Sprintf("Invalid keystore version: %d, expected version 3", version),
			"version",
			nil,
		)
	}
	return nil
}

// ValidateAddress checks if the address is a valid Ethereum address
func (kv *KeystoreValidator) ValidateAddress(address string) error {
	if address == "" {
		return NewKeystoreImportErrorWithField(
			ErrorMissingRequiredFields,
			"Missing required field: address",
			"address",
			nil,
		)
	}

	// Remove 0x prefix if present for validation
	cleanAddress := address
	if strings.HasPrefix(strings.ToLower(address), "0x") {
		cleanAddress = address[2:]
	}

	// Ethereum addresses are 40 hex characters
	matched, err := regexp.MatchString("^[0-9a-fA-F]{40}$", cleanAddress)
	if err != nil {
		return NewKeystoreImportErrorWithField(
			ErrorInvalidAddress,
			"Error validating address format",
			"address",
			err,
		)
	}

	if !matched {
		return NewKeystoreImportErrorWithField(
			ErrorInvalidAddress,
			fmt.Sprintf("Invalid Ethereum address format: %s", address),
			"address",
			nil,
		)
	}

	return nil
}

// ValidateCrypto checks if the crypto section has all required fields
func (kv *KeystoreValidator) ValidateCrypto(crypto *KeystoreV3Crypto) error {
	if crypto.Cipher == "" {
		return NewKeystoreImportErrorWithField(
			ErrorMissingRequiredFields,
			"Missing required field: crypto.cipher",
			"crypto.cipher",
			nil,
		)
	}

	if crypto.CipherText == "" {
		return NewKeystoreImportErrorWithField(
			ErrorMissingRequiredFields,
			"Missing required field: crypto.ciphertext",
			"crypto.ciphertext",
			nil,
		)
	}

	if crypto.CipherParams.IV == "" {
		return NewKeystoreImportErrorWithField(
			ErrorMissingRequiredFields,
			"Missing required field: crypto.cipherparams.iv",
			"crypto.cipherparams.iv",
			nil,
		)
	}

	if crypto.KDF == "" {
		return NewKeystoreImportErrorWithField(
			ErrorMissingRequiredFields,
			"Missing required field: crypto.kdf",
			"crypto.kdf",
			nil,
		)
	}

	if crypto.KDFParams == nil {
		return NewKeystoreImportErrorWithField(
			ErrorMissingRequiredFields,
			"Missing required field: crypto.kdfparams",
			"crypto.kdfparams",
			nil,
		)
	}

	if crypto.MAC == "" {
		return NewKeystoreImportErrorWithField(
			ErrorMissingRequiredFields,
			"Missing required field: crypto.mac",
			"crypto.mac",
			nil,
		)
	}

	// Validate KDF parameters based on KDF type
	switch strings.ToLower(crypto.KDF) {
	case "scrypt":
		return kv.validateScryptParams(crypto.KDFParams)
	case "pbkdf2":
		return kv.validatePBKDF2Params(crypto.KDFParams)
	default:
		return NewKeystoreImportErrorWithField(
			ErrorInvalidKeystore,
			fmt.Sprintf("Unsupported KDF algorithm: %s", crypto.KDF),
			"crypto.kdf",
			nil,
		)
	}
}

// validateScryptParams validates scrypt KDF parameters
func (kv *KeystoreValidator) validateScryptParams(params any) error {
	// Convert interface{} to map for validation
	paramsMap, ok := params.(map[string]any)
	if !ok {
		return NewKeystoreImportErrorWithField(
			ErrorInvalidKeystore,
			"Invalid scrypt parameters format",
			"crypto.kdfparams",
			nil,
		)
	}

	// Check required fields
	requiredFields := []string{"dklen", "n", "r", "p", "salt"}
	for _, field := range requiredFields {
		if _, exists := paramsMap[field]; !exists {
			return NewKeystoreImportErrorWithField(
				ErrorMissingRequiredFields,
				fmt.Sprintf("Missing required field: crypto.kdfparams.%s", field),
				fmt.Sprintf("crypto.kdfparams.%s", field),
				nil,
			)
		}
	}

	return nil
}

// validatePBKDF2Params validates PBKDF2 KDF parameters
func (kv *KeystoreValidator) validatePBKDF2Params(params any) error {
	// Convert interface{} to map for validation
	paramsMap, ok := params.(map[string]any)
	if !ok {
		return NewKeystoreImportErrorWithField(
			ErrorInvalidKeystore,
			"Invalid PBKDF2 parameters format",
			"crypto.kdfparams",
			nil,
		)
	}

	// Check required fields
	requiredFields := []string{"dklen", "c", "prf", "salt"}
	for _, field := range requiredFields {
		if _, exists := paramsMap[field]; !exists {
			return NewKeystoreImportErrorWithField(
				ErrorMissingRequiredFields,
				fmt.Sprintf("Missing required field: crypto.kdfparams.%s", field),
				fmt.Sprintf("crypto.kdfparams.%s", field),
				nil,
			)
		}
	}

	return nil
}

// NewKeystoreImportError creates a new KeystoreImportError
func NewKeystoreImportError(errorType KeystoreErrorType, message string, cause error) *KeystoreImportError {
	return &KeystoreImportError{
		Type:    errorType,
		Message: message,
		Cause:   cause,
	}
}

// NewKeystoreImportErrorWithField creates a new KeystoreImportError with a specific field
func NewKeystoreImportErrorWithField(errorType KeystoreErrorType, message string, field string, cause error) *KeystoreImportError {
	return &KeystoreImportError{
		Type:    errorType,
		Message: message,
		Field:   field,
		Cause:   cause,
	}
}

// GetErrorTypeString returns a string representation of the error type
func (e KeystoreErrorType) String() string {
	switch e {
	case ErrorFileNotFound:
		return "FILE_NOT_FOUND"
	case ErrorInvalidJSON:
		return "INVALID_JSON"
	case ErrorInvalidKeystore:
		return "INVALID_KEYSTORE"
	case ErrorInvalidVersion:
		return "INVALID_VERSION"
	case ErrorIncorrectPassword:
		return "INCORRECT_PASSWORD"
	case ErrorCorruptedFile:
		return "CORRUPTED_FILE"
	case ErrorAddressMismatch:
		return "ADDRESS_MISMATCH"
	case ErrorMissingRequiredFields:
		return "MISSING_REQUIRED_FIELDS"
	case ErrorInvalidAddress:
		return "INVALID_ADDRESS"
	default:
		return "UNKNOWN_ERROR"
	}
}

// GetLocalizationKey returns the localization key for this error type
func (e KeystoreErrorType) GetLocalizationKey() string {
	switch e {
	case ErrorFileNotFound:
		return "keystore_file_not_found"
	case ErrorInvalidJSON:
		return "keystore_invalid_json"
	case ErrorInvalidKeystore:
		return "keystore_invalid_structure"
	case ErrorInvalidVersion:
		return "keystore_invalid_version"
	case ErrorIncorrectPassword:
		return "keystore_incorrect_password"
	case ErrorCorruptedFile:
		return "keystore_corrupted_file"
	case ErrorAddressMismatch:
		return "keystore_address_mismatch"
	case ErrorMissingRequiredFields:
		return "keystore_missing_fields"
	case ErrorInvalidAddress:
		return "keystore_invalid_address"
	default:
		return "unknown_error"
	}
}

// GetLocalizedMessage returns a localized error message for this error
func (e *KeystoreImportError) GetLocalizedMessage() string {
	// This will be used by UI code to get localized messages
	key := e.Type.GetLocalizationKey()

	// The actual localization will be done by the UI layer
	// to avoid import cycles between packages
	return key
}

// GetLocalizedMessageWithField returns a localized error message with field information
func (e *KeystoreImportError) GetLocalizedMessageWithField() string {
	key := e.Type.GetLocalizationKey()

	// The actual localization will be done by the UI layer
	// but we return the field information here
	if e.Field != "" {
		return key + ":" + e.Field
	}
	return key
}
