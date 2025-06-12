package wallet

import (
	"blocowallet/pkg/config"
	"blocowallet/pkg/localization"
	"crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/base64"
	"fmt"

	"golang.org/x/crypto/argon2"
)

// CryptoService fornece funcionalidades de criptografia para dados sensíveis
type CryptoService struct {
	config *config.Config
}

// NewCryptoService cria uma nova instância do serviço de criptografia
func NewCryptoService(cfg *config.Config) *CryptoService {
	return &CryptoService{
		config: cfg,
	}
}

// EncryptMnemonic criptografa uma frase mnemônica usando Argon2ID
func (cs *CryptoService) EncryptMnemonic(mnemonic, password string) (string, error) {
	if password == "" {
		return "", fmt.Errorf(localization.Get("error_empty_password"))
	}

	// Obter configurações do Argon2id
	saltLength := int(cs.config.Security.SaltLength)
	argon2IDTime := cs.config.Security.Argon2Time
	argon2IDMemory := cs.config.Security.Argon2Memory
	argon2IDThreads := cs.config.Security.Argon2Threads
	argon2IDKeyLen := cs.config.Security.Argon2KeyLen

	// Gerar salt aleatório
	salt := make([]byte, saltLength)
	if _, err := rand.Read(salt); err != nil {
		return "", fmt.Errorf(localization.Get("error_generate_salt")+": %w", err)
	}

	// Derivar chave usando Argon2ID
	key := argon2.IDKey([]byte(password), salt, argon2IDTime, argon2IDMemory, argon2IDThreads, argon2IDKeyLen)

	// Criar hash de verificação da mnemônica original + senha
	// Isso garante que mesmo mnemônicas vazias tenham hashes únicos por senha
	mnemonicBytes := []byte(mnemonic)
	hashInput := append(mnemonicBytes, []byte(password)...)
	hash := sha256.Sum256(hashInput)

	// Criptografar a mnemônica com XOR
	encrypted := make([]byte, len(mnemonicBytes))

	// Repetir a chave se a mnemônica for maior que a chave
	for i := 0; i < len(mnemonicBytes); i++ {
		encrypted[i] = mnemonicBytes[i] ^ key[i%len(key)]
	}

	// Combinar salt + hash + dados criptografados e codificar para base64
	combined := append(salt, hash[:]...)
	combined = append(combined, encrypted...)
	return base64.StdEncoding.EncodeToString(combined), nil
}

// DecryptMnemonic descriptografa uma mnemônica usando Argon2ID
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

	// Descriptografar a mnemônica com XOR
	decrypted := make([]byte, len(encrypted))
	for i := 0; i < len(encrypted); i++ {
		decrypted[i] = encrypted[i] ^ key[i%len(key)]
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
