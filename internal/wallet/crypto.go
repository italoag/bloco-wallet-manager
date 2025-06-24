package wallet

import (
	"blocowallet/pkg/config"
	"blocowallet/pkg/localization"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/base64"
	"fmt"

	"golang.org/x/crypto/argon2"
)

// CryptoService provides cryptographic functionality for sensitive data
type CryptoService struct {
	config *config.Config
}

// NewCryptoService creates a new instance of the cryptography service
func NewCryptoService(cfg *config.Config) *CryptoService {
	return &CryptoService{
		config: cfg,
	}
}

// EncryptMnemonic encrypts a mnemonic phrase using Argon2ID
func (cs *CryptoService) EncryptMnemonic(mnemonic, password string) (string, error) {
	if password == "" {
		return "", fmt.Errorf(localization.Get("error_empty_password"))
	}

	// Get Argon2id configurations
	saltLength := int(cs.config.Security.SaltLength)
	argon2IDTime := cs.config.Security.Argon2Time
	argon2IDMemory := cs.config.Security.Argon2Memory
	argon2IDThreads := cs.config.Security.Argon2Threads
	argon2IDKeyLen := cs.config.Security.Argon2KeyLen

	// Generate random salt
	salt := make([]byte, saltLength)
	if _, err := rand.Read(salt); err != nil {
		return "", fmt.Errorf(localization.Get("error_generate_salt")+": %w", err)
	}

	// Derive key using Argon2ID
	key := argon2.IDKey([]byte(password), salt, argon2IDTime, argon2IDMemory, argon2IDThreads, argon2IDKeyLen)

	// Create a verification hash of the original mnemonic + password
	// This ensures that even empty mnemonics have unique hashes per password
	mnemonicBytes := []byte(mnemonic)
	hashInput := append(mnemonicBytes, []byte(password)...)
	hash := sha256.Sum256(hashInput)

	// Encrypt the mnemonic using AES-GCM
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", fmt.Errorf("failed to create cipher: %w", err)
	}

	aead, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("failed to create AEAD: %w", err)
	}

	nonce := make([]byte, aead.NonceSize())
	if _, err := rand.Read(nonce); err != nil {
		return "", fmt.Errorf("failed to generate nonce: %w", err)
	}

	encrypted := aead.Seal(nil, nonce, mnemonicBytes, nil)

	// Combine salt + hash + nonce + encrypted data and encode to base64
	combined := append(salt, hash[:]...)
	combined = append(combined, nonce...)
	combined = append(combined, encrypted...)
	return base64.StdEncoding.EncodeToString(combined), nil
}

// DecryptMnemonic decrypts a mnemonic using Argon2ID
func (cs *CryptoService) DecryptMnemonic(encryptedMnemonic, password string) (string, error) {
	if encryptedMnemonic == "" {
		return "", fmt.Errorf(localization.Get("error_empty_encrypted_mnemonic"))
	}
	if password == "" {
		return "", fmt.Errorf(localization.Get("error_empty_password"))
	}

	// Decodificar de base64
	combined, err := base64.StdEncoding.DecodeString(encryptedMnemonic)
	if err != nil {
		return "", fmt.Errorf(localization.Get("error_decode_mnemonic")+": %w", err)
	}

	// Obter configurações do Argon2id
	saltLength := int(cs.config.Security.SaltLength)
	hashLength := 32 // SHA-256 sempre tem 32 bytes

	if len(combined) < saltLength+hashLength {
		return "", fmt.Errorf(localization.Get("error_invalid_mnemonic_format"))
	}

	// Extrair salt, hash e dados criptografados
	salt := combined[:saltLength]
	expectedHash := combined[saltLength : saltLength+hashLength]
	encrypted := combined[saltLength+hashLength:]

	// Derivar chave usando Argon2ID com os mesmos parâmetros
	key := argon2.IDKey(
		[]byte(password),
		salt,
		cs.config.Security.Argon2Time,
		cs.config.Security.Argon2Memory,
		cs.config.Security.Argon2Threads,
		cs.config.Security.Argon2KeyLen,
	)

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", fmt.Errorf("failed to create cipher: %w", err)
	}
	aead, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("failed to create AEAD: %w", err)
	}

	if len(encrypted) < aead.NonceSize() {
		return "", fmt.Errorf(localization.Get("error_invalid_mnemonic_format"))
	}

	nonce := encrypted[:aead.NonceSize()]
	cipherText := encrypted[aead.NonceSize():]

	decrypted, err := aead.Open(nil, nonce, cipherText, nil)
	if err != nil {
		return "", fmt.Errorf(localization.Get("error_invalid_password"))
	}

	// Verificar o hash para garantir que a senha está correta
	hashInput := append(decrypted, []byte(password)...)
	actualHash := sha256.Sum256(hashInput)
	if subtle.ConstantTimeCompare(expectedHash, actualHash[:]) != 1 {
		return "", fmt.Errorf(localization.Get("error_invalid_password"))
	}

	mnemonic := string(decrypted)
	return mnemonic, nil
}

// VerifyMnemonicPassword verifica se a senha pode descriptografar a mnemônica
func (cs *CryptoService) VerifyMnemonicPassword(encryptedMnemonic, password string) bool {
	_, err := cs.DecryptMnemonic(encryptedMnemonic, password)
	return err == nil
}

// SecureCompare realiza comparação em tempo constante de duas strings
func SecureCompare(a, b string) bool {
	return subtle.ConstantTimeCompare([]byte(a), []byte(b)) == 1
}

// PasswordValidationError representa os erros de validação de senha
type PasswordValidationError struct {
	TooShort         bool
	NoLowercase      bool
	NoUppercase      bool
	NoDigitOrSpecial bool
}

// Error implementa a interface error
func (e PasswordValidationError) Error() string {
	if !e.HasErrors() {
		return ""
	}

	return localization.Get("password_validation_failed")
}

// HasErrors verifica se há algum erro de validação
func (e PasswordValidationError) HasErrors() bool {
	return e.TooShort || e.NoLowercase || e.NoUppercase || e.NoDigitOrSpecial
}

// GetErrorMessage retorna uma mensagem de erro específica para o primeiro erro encontrado
func (e PasswordValidationError) GetErrorMessage() string {
	if e.TooShort {
		return localization.Get("password_too_short")
	}
	if e.NoLowercase {
		return localization.Get("password_no_lowercase")
	}
	if e.NoUppercase {
		return localization.Get("password_no_uppercase")
	}
	if e.NoDigitOrSpecial {
		return localization.Get("password_no_digit_or_special")
	}
	return ""
}

// ValidatePassword valida a complexidade da senha
func ValidatePassword(password string) (PasswordValidationError, bool) {
	var err PasswordValidationError

	// Verificar tamanho mínimo
	if len(password) < 8 {
		err.TooShort = true
		return err, false
	}

	// Verificar presença de letra minúscula
	hasLower := false
	// Verificar presença de letra maiúscula
	hasUpper := false
	// Verificar presença de dígito ou caractere especial
	hasDigitOrSpecial := false

	for _, c := range password {
		switch {
		case c >= 'a' && c <= 'z':
			hasLower = true
		case c >= 'A' && c <= 'Z':
			hasUpper = true
		default:
			hasDigitOrSpecial = true
		}
	}

	err.NoLowercase = !hasLower
	err.NoUppercase = !hasUpper
	err.NoDigitOrSpecial = !hasDigitOrSpecial

	return err, !err.HasErrors()
}

// Para compatibilidade com código existente, fornecendo funções estáticas
var defaultCryptoService *CryptoService

// InitCryptoService inicializa o serviço de criptografia padrão
func InitCryptoService(cfg *config.Config) {
	defaultCryptoService = NewCryptoService(cfg)
}

// Funções auxiliares para compatibilidade com código existente
func EncryptMnemonic(mnemonic, password string) (string, error) {
	if defaultCryptoService == nil {
		return "", fmt.Errorf(localization.Get("error_crypto_service_not_initialized"))
	}
	return defaultCryptoService.EncryptMnemonic(mnemonic, password)
}

func DecryptMnemonic(encryptedMnemonic, password string) (string, error) {
	if defaultCryptoService == nil {
		return "", fmt.Errorf(localization.Get("error_crypto_service_not_initialized"))
	}
	return defaultCryptoService.DecryptMnemonic(encryptedMnemonic, password)
}

func VerifyMnemonicPassword(encryptedMnemonic, password string) bool {
	if defaultCryptoService == nil {
		return false
	}
	return defaultCryptoService.VerifyMnemonicPassword(encryptedMnemonic, password)
}
