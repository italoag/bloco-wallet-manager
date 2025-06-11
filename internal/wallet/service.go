package wallet

import (
	"context"
	"crypto/ecdsa"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"blocowallet/internal/blockchain"
	"blocowallet/pkg/config"
	"blocowallet/pkg/logger"

	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/google/uuid"
)

// KeystoreFile represents the structure of a keystore JSON file
type KeystoreFile struct {
	Address string `json:"address"`
	Crypto  struct {
		Cipher       string `json:"cipher"`
		CipherText   string `json:"ciphertext"`
		CipherParams struct {
			IV string `json:"iv"`
		} `json:"cipherparams"`
		KDF       string `json:"kdf"`
		KDFParams struct {
			DKLen int    `json:"dklen"`
			Salt  string `json:"salt"`
			N     int    `json:"n"`
			R     int    `json:"r"`
			P     int    `json:"p"`
		} `json:"kdfparams"`
		MAC string `json:"mac"`
	} `json:"crypto"`
	ID      string `json:"id"`
	Version int    `json:"version"`
}

// KeyStore V3 constants
const (
	KeyStoreV3Dir   = ".blocowallet/keystore"
	KeyStoreVersion = 3
)

// Service provides wallet business logic
type Service struct {
	repo            Repository
	balanceProvider BalanceProvider
	multiProvider   *blockchain.MultiProvider
	keystore        *keystore.KeyStore
	passwordCache   map[string]string // Cache for wallet passwords
	passwordMutex   sync.RWMutex      // Mutex for thread-safe password cache access
	logger          logger.Logger     // Logger for structured logging
}

// NewService creates a new wallet service
func NewService(repo Repository, balanceProvider BalanceProvider, log logger.Logger) *Service {
	return &Service{
		repo:            repo,
		balanceProvider: balanceProvider,
		passwordCache:   make(map[string]string),
		logger:          log,
	}
}

// NewServiceWithMultiProvider creates a new wallet service with multi-provider support
func NewServiceWithMultiProvider(repo Repository, multiProvider *blockchain.MultiProvider, log logger.Logger) *Service {
	return &Service{
		repo:          repo,
		multiProvider: multiProvider,
		passwordCache: make(map[string]string),
		logger:        log,
	}
}

// NewServiceWithKeystore creates a new wallet service with keystore support
func NewServiceWithKeystore(repo Repository, balanceProvider BalanceProvider, ks *keystore.KeyStore, log logger.Logger) *Service {
	return &Service{
		repo:            repo,
		balanceProvider: balanceProvider,
		keystore:        ks,
		passwordCache:   make(map[string]string),
		logger:          log,
	}
}

// NewServiceWithMultiProviderAndKeystore creates a new wallet service with multi-provider and keystore support
func NewServiceWithMultiProviderAndKeystore(repo Repository, multiProvider *blockchain.MultiProvider, ks *keystore.KeyStore, log logger.Logger) *Service {
	return &Service{
		repo:          repo,
		multiProvider: multiProvider,
		keystore:      ks,
		passwordCache: make(map[string]string),
		logger:        log,
	}
}

// Create creates a new wallet
func (s *Service) Create(ctx context.Context, name, address, keystorePath string) (*Wallet, error) {
	correlationID := uuid.New().String()
	ctx = context.WithValue(ctx, "correlation_id", correlationID)

	// s.logger.Info("Creating new wallet",
	// 	logger.String("correlation_id", correlationID),
	// 	logger.String("wallet_name", name),
	// 	logger.String("address", address),
	// 	logger.String("operation", "create_wallet"))

	if name == "" {
		err := NewValidationError("wallet name cannot be empty")
		s.logger.Error("Wallet creation failed",
			logger.String("correlation_id", correlationID),
			logger.String("error", err.Error()),
			logger.String("operation", "create_wallet"))
		return nil, AddCorrelationID(ctx, err)
	}

	if address == "" {
		err := NewValidationError("wallet address cannot be empty")
		s.logger.Error("Wallet creation failed",
			logger.String("correlation_id", correlationID),
			logger.String("error", err.Error()),
			logger.String("operation", "create_wallet"))
		return nil, AddCorrelationID(ctx, err)
	}

	// Check if wallet with this address already exists
	existing, err := s.repo.GetByAddress(ctx, address)
	if err == nil && existing != nil {
		err := NewAlreadyExistsError(fmt.Sprintf("wallet with address %s already exists", address))
		s.logger.Error("Wallet creation failed - already exists",
			logger.String("correlation_id", correlationID),
			logger.String("address", address),
			logger.String("error", err.Error()),
			logger.String("operation", "create_wallet"))
		return nil, AddCorrelationID(ctx, err)
	}

	wallet := &Wallet{
		ID:           uuid.New().String(),
		Name:         name,
		Address:      address,
		KeyStorePath: keystorePath,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	if err := s.repo.Create(ctx, wallet); err != nil {
		s.logger.Error("Failed to save wallet to database",
			logger.String("correlation_id", correlationID),
			logger.String("wallet_id", wallet.ID),
			logger.Error(err),
			logger.String("operation", "create_wallet"))
		return nil, AddCorrelationID(ctx, NewStorageError("failed to create wallet", err))
	}

	// s.logger.Info("Wallet created successfully",
	// 	logger.String("correlation_id", correlationID),
	// 	logger.String("wallet_id", wallet.ID),
	// 	logger.String("address", wallet.Address),
	// 	logger.String("operation", "create_wallet"))

	return wallet, nil
}

// GetByID retrieves a wallet by ID
func (s *Service) GetByID(ctx context.Context, id string) (*Wallet, error) {
	correlationID := uuid.New().String()
	ctx = context.WithValue(ctx, "correlation_id", correlationID)

	s.logger.Debug("Getting wallet by ID",
		logger.String("correlation_id", correlationID),
		logger.String("wallet_id", id),
		logger.String("operation", "get_wallet_by_id"))

	if id == "" {
		err := NewValidationError("wallet ID cannot be empty")
		s.logger.Error("Get wallet failed",
			logger.String("correlation_id", correlationID),
			logger.String("error", err.Error()),
			logger.String("operation", "get_wallet_by_id"))
		return nil, AddCorrelationID(ctx, err)
	}

	wallet, err := s.repo.GetByID(ctx, id)
	if err != nil {
		s.logger.Error("Failed to retrieve wallet from database",
			logger.String("correlation_id", correlationID),
			logger.String("wallet_id", id),
			logger.Error(err),
			logger.String("operation", "get_wallet_by_id"))
		return nil, AddCorrelationID(ctx, NewStorageError("failed to get wallet", err))
	}

	s.logger.Debug("Wallet retrieved successfully",
		logger.String("correlation_id", correlationID),
		logger.String("wallet_id", id),
		logger.String("operation", "get_wallet_by_id"))

	return wallet, nil
}

// GetByAddress retrieves a wallet by address
func (s *Service) GetByAddress(ctx context.Context, address string) (*Wallet, error) {
	if address == "" {
		return nil, fmt.Errorf("wallet address cannot be empty")
	}

	wallet, err := s.repo.GetByAddress(ctx, address)
	if err != nil {
		return nil, fmt.Errorf("failed to get wallet: %w", err)
	}

	return wallet, nil
}

// List retrieves all wallets
func (s *Service) List(ctx context.Context) ([]*Wallet, error) {
	wallets, err := s.repo.List(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list wallets: %w", err)
	}

	return wallets, nil
}

// Update updates a wallet
func (s *Service) Update(ctx context.Context, wallet *Wallet) error {
	if wallet == nil {
		return fmt.Errorf("wallet cannot be nil")
	}

	if wallet.ID == "" {
		return fmt.Errorf("wallet ID cannot be empty")
	}

	wallet.UpdatedAt = time.Now()

	if err := s.repo.Update(ctx, wallet); err != nil {
		return fmt.Errorf("failed to update wallet: %w", err)
	}

	return nil
}

// Delete deletes a wallet
func (s *Service) Delete(ctx context.Context, id string) error {
	if id == "" {
		return fmt.Errorf("wallet ID cannot be empty")
	}

	if err := s.repo.Delete(ctx, id); err != nil {
		return fmt.Errorf("failed to delete wallet: %w", err)
	}

	return nil
}

// GetBalance gets the balance for a wallet
func (s *Service) GetBalance(ctx context.Context, address string) (*Balance, error) {
	if address == "" {
		return nil, fmt.Errorf("wallet address cannot be empty")
	}

	amount, err := s.balanceProvider.GetBalance(ctx, address)
	if err != nil {
		return nil, fmt.Errorf("failed to get balance: %w", err)
	}

	return &Balance{
		Address:   address,
		Amount:    amount,
		Symbol:    "ETH",
		Decimals:  18,
		UpdatedAt: time.Now(),
	}, nil
}

// CreateWalletWithMnemonic creates a new wallet with mnemonic and keystore
func (s *Service) CreateWalletWithMnemonic(ctx context.Context, name, password string) (*WalletDetails, error) {
	correlationID := uuid.New().String()
	ctx = context.WithValue(ctx, "correlation_id", correlationID)

	// s.logger.Info("Creating wallet with mnemonic",
	// 	logger.String("correlation_id", correlationID),
	// 	logger.String("wallet_name", name),
	// 	logger.String("operation", "create_wallet_with_mnemonic"))

	if name == "" {
		err := NewValidationError("wallet name cannot be empty")
		s.logger.Error("Wallet creation failed - empty name",
			logger.String("correlation_id", correlationID),
			logger.String("error", err.Error()),
			logger.String("operation", "create_wallet_with_mnemonic"))
		return nil, AddCorrelationID(ctx, err)
	}

	if password == "" {
		err := NewValidationError("password cannot be empty")
		s.logger.Error("Wallet creation failed - empty password",
			logger.String("correlation_id", correlationID),
			logger.String("error", err.Error()),
			logger.String("operation", "create_wallet_with_mnemonic"))
		return nil, AddCorrelationID(ctx, err)
	}

	// Generate mnemonic
	s.logger.Debug("Generating mnemonic",
		logger.String("correlation_id", correlationID),
		logger.String("operation", "create_wallet_with_mnemonic"))

	mnemonic, err := GenerateMnemonic()
	if err != nil {
		s.logger.Error("Failed to generate mnemonic",
			logger.String("correlation_id", correlationID),
			logger.Error(err),
			logger.String("operation", "create_wallet_with_mnemonic"))
		return nil, AddCorrelationID(ctx, NewCryptoError("failed to generate mnemonic", err))
	}

	// Derive private key
	s.logger.Debug("Deriving private key from mnemonic",
		logger.String("correlation_id", correlationID),
		logger.String("operation", "create_wallet_with_mnemonic"))

	privateKeyHex, err := DerivePrivateKey(mnemonic)
	if err != nil {
		s.logger.Error("Failed to derive private key",
			logger.String("correlation_id", correlationID),
			logger.Error(err),
			logger.String("operation", "create_wallet_with_mnemonic"))
		return nil, AddCorrelationID(ctx, NewCryptoError("failed to derive private key", err))
	}

	privKey, err := HexToECDSA(privateKeyHex)
	if err != nil {
		s.logger.Error("Failed to convert private key",
			logger.String("correlation_id", correlationID),
			logger.Error(err),
			logger.String("operation", "create_wallet_with_mnemonic"))
		return nil, AddCorrelationID(ctx, NewCryptoError("failed to convert private key", err))
	}

	// Get address
	address := GetAddressFromPrivateKey(privKey)
	s.logger.Debug("Generated wallet address",
		logger.String("correlation_id", correlationID),
		logger.String("address", address),
		logger.String("operation", "create_wallet_with_mnemonic"))

	// Create KeyStore V3 file in proper directory structure
	s.logger.Debug("Creating keystore file",
		logger.String("correlation_id", correlationID),
		logger.String("operation", "create_wallet_with_mnemonic"))

	keystorePath, err := s.CreateKeyStoreV3File(privKey, password)
	if err != nil {
		s.logger.Error("Failed to create keystore file",
			logger.String("correlation_id", correlationID),
			logger.Error(err),
			logger.String("operation", "create_wallet_with_mnemonic"))
		return nil, AddCorrelationID(ctx, NewCryptoError("failed to create keystore file", err))
	}

	// Encrypt mnemonic for secure storage in database
	s.logger.Debug("Encrypting mnemonic for storage",
		logger.String("correlation_id", correlationID),
		logger.String("operation", "create_wallet_with_mnemonic"))

	encryptedMnemonic, err := EncryptMnemonic(mnemonic, password)
	if err != nil {
		s.logger.Error("Failed to encrypt mnemonic",
			logger.String("correlation_id", correlationID),
			logger.Error(err),
			logger.String("operation", "create_wallet_with_mnemonic"))
		// Cleanup keystore file on encryption failure
		if keystorePath != "" {
			os.Remove(keystorePath)
		}
		return nil, AddCorrelationID(ctx, NewCryptoError("failed to encrypt mnemonic", err))
	}

	// Create wallet entity - store encrypted mnemonic
	wallet := &Wallet{
		ID:                uuid.New().String(),
		Name:              name,
		Address:           address,
		KeyStorePath:      keystorePath,
		EncryptedMnemonic: encryptedMnemonic,
		CreatedAt:         time.Now(),
		UpdatedAt:         time.Now(),
	}

	// Save to repository
	s.logger.Debug("Saving wallet to database",
		logger.String("correlation_id", correlationID),
		logger.String("wallet_id", wallet.ID),
		logger.String("operation", "create_wallet_with_mnemonic"))

	if err := s.repo.Create(ctx, wallet); err != nil {
		s.logger.Error("Failed to save wallet to database",
			logger.String("correlation_id", correlationID),
			logger.String("wallet_id", wallet.ID),
			logger.Error(err),
			logger.String("operation", "create_wallet_with_mnemonic"))

		// If database save fails and keystore was created, try to clean up
		if s.keystore != nil && keystorePath != "" {
			os.Remove(keystorePath)
			s.logger.Debug("Cleaned up keystore file after database failure",
				logger.String("correlation_id", correlationID),
				logger.String("keystore_path", keystorePath),
				logger.String("operation", "create_wallet_with_mnemonic"))
		}
		return nil, AddCorrelationID(ctx, NewStorageError("failed to save wallet", err))
	}

	// Get public key
	publicKey := privKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		err := NewCryptoError("failed to get public key", fmt.Errorf("invalid public key type"))
		s.logger.Error("Failed to extract public key",
			logger.String("correlation_id", correlationID),
			logger.String("wallet_id", wallet.ID),
			logger.String("error", err.Error()),
			logger.String("operation", "create_wallet_with_mnemonic"))
		return nil, AddCorrelationID(ctx, err)
	}

	// s.logger.Info("Wallet created successfully with mnemonic",
	// 	logger.String("correlation_id", correlationID),
	// 	logger.String("wallet_id", wallet.ID),
	// 	logger.String("address", wallet.Address),
	// 	logger.String("keystore_path", keystorePath),
	// 	logger.String("operation", "create_wallet_with_mnemonic"))

	return &WalletDetails{
		Wallet:     wallet,
		Mnemonic:   mnemonic,
		PrivateKey: privKey,
		PublicKey:  publicKeyECDSA,
	}, nil
}

// ImportWalletFromMnemonic imports a wallet from an existing mnemonic
func (s *Service) ImportWalletFromMnemonic(ctx context.Context, name, mnemonic, password string) (*WalletDetails, error) {
	if name == "" {
		return nil, fmt.Errorf("wallet name cannot be empty")
	}

	if mnemonic == "" {
		return nil, fmt.Errorf("mnemonic cannot be empty")
	}

	if password == "" {
		return nil, fmt.Errorf("password cannot be empty")
	}

	// Validate mnemonic
	if !IsValidMnemonic(mnemonic) {
		return nil, fmt.Errorf("invalid mnemonic phrase")
	}

	// Derive private key
	privateKeyHex, err := DerivePrivateKey(mnemonic)
	if err != nil {
		return nil, fmt.Errorf("failed to derive private key: %w", err)
	}

	privKey, err := HexToECDSA(privateKeyHex)
	if err != nil {
		return nil, fmt.Errorf("failed to convert private key: %w", err)
	}

	// Get address
	address := GetAddressFromPrivateKey(privKey)

	// Check if wallet already exists
	existing, err := s.repo.GetByAddress(ctx, address)
	if err == nil && existing != nil {
		return nil, fmt.Errorf("wallet with address %s already exists", address)
	}

	// Create KeyStore V3 file in proper directory structure
	keystorePath, err := s.CreateKeyStoreV3File(privKey, password)
	if err != nil {
		return nil, fmt.Errorf("failed to create keystore file: %w", err)
	}

	// Encrypt mnemonic for secure storage in database
	encryptedMnemonic, err := EncryptMnemonic(mnemonic, password)
	if err != nil {
		return nil, fmt.Errorf("failed to encrypt mnemonic: %w", err)
	}

	// Create wallet entity - store encrypted mnemonic
	wallet := &Wallet{
		ID:                uuid.New().String(),
		Name:              name,
		Address:           address,
		KeyStorePath:      keystorePath,
		EncryptedMnemonic: encryptedMnemonic,
		CreatedAt:         time.Now(),
		UpdatedAt:         time.Now(),
	}

	// Save to repository
	if err := s.repo.Create(ctx, wallet); err != nil {
		// If database save fails and keystore was created, try to clean up
		if s.keystore != nil && keystorePath != "" {
			os.Remove(keystorePath)
		}
		return nil, fmt.Errorf("failed to save wallet: %w", err)
	}

	// Get public key
	publicKey := privKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return nil, fmt.Errorf("failed to get public key")
	}

	return &WalletDetails{
		Wallet:     wallet,
		Mnemonic:   mnemonic,
		PrivateKey: privKey,
		PublicKey:  publicKeyECDSA,
	}, nil
}

// ImportWalletFromPrivateKey imports a wallet from a private key
func (s *Service) ImportWalletFromPrivateKey(ctx context.Context, name, privateKeyHex, password string) (*WalletDetails, error) {
	if name == "" {
		return nil, fmt.Errorf("wallet name cannot be empty")
	}

	if privateKeyHex == "" {
		return nil, fmt.Errorf("private key cannot be empty")
	}

	if password == "" {
		return nil, fmt.Errorf("password cannot be empty")
	}

	// Remove 0x prefix if present
	if len(privateKeyHex) > 2 && privateKeyHex[:2] == "0x" {
		privateKeyHex = privateKeyHex[2:]
	}

	// Convert hex string to ECDSA private key
	privKey, err := HexToECDSA(privateKeyHex)
	if err != nil {
		return nil, fmt.Errorf("invalid private key: %w", err)
	}

	// Derive address from private key
	address := GetAddressFromPrivateKey(privKey)

	// Check if wallet already exists
	existingWallets, err := s.repo.List(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to check existing wallets: %w", err)
	}

	for _, existingWallet := range existingWallets {
		if existingWallet.Address == address {
			return nil, fmt.Errorf("wallet with address %s already exists", address)
		}
		if existingWallet.Name == name {
			return nil, fmt.Errorf("wallet with name %s already exists", name)
		}
	}

	// Create KeyStore V3 file to store the private key securely
	keystorePath, err := s.CreateKeyStoreV3File(privKey, password)
	if err != nil {
		return nil, fmt.Errorf("failed to create keystore file: %w", err)
	}

	// Create wallet object - NO private key stored in database, only keystore path
	// Private key imports don't have mnemonic, so EncryptedMnemonic is empty
	wallet := &Wallet{
		ID:                uuid.New().String(),
		Name:              name,
		Address:           address,
		KeyStorePath:      keystorePath,
		EncryptedMnemonic: "", // Private key imports don't have mnemonic
		CreatedAt:         time.Now(),
		UpdatedAt:         time.Now(),
	}

	// Save to repository
	err = s.repo.Create(ctx, wallet)
	if err != nil {
		// If database save fails and keystore was created, try to clean up
		if s.keystore != nil && keystorePath != "" {
			os.Remove(keystorePath)
		}
		return nil, fmt.Errorf("failed to save wallet: %w", err)
	}

	// Get public key
	publicKey := privKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return nil, fmt.Errorf("failed to get public key")
	}

	return &WalletDetails{
		Wallet:     wallet,
		PrivateKey: privKey,
		PublicKey:  publicKeyECDSA,
	}, nil
}

// ExtractPrivateKeyFromKeystore extracts the private key from a keystore file
func (s *Service) ExtractPrivateKeyFromKeystore(keystorePath, password string) (string, error) {
	if keystorePath == "" {
		return "", fmt.Errorf("keystore path is empty")
	}

	if password == "" {
		return "", fmt.Errorf("password is required to decrypt keystore")
	}

	// Use the new LoadPrivateKeyFromKeyStoreV3 method
	privateKey, err := s.LoadPrivateKeyFromKeyStoreV3(keystorePath, password)
	if err != nil {
		return "", fmt.Errorf("failed to load private key from keystore: %w", err)
	}

	// Convert private key to hex string
	privateKeyBytes := crypto.FromECDSA(privateKey)
	return hex.EncodeToString(privateKeyBytes), nil
}

// CreateKeyStoreV3File creates a KeyStore V3 file in the proper directory structure
func (s *Service) CreateKeyStoreV3File(privateKey *ecdsa.PrivateKey, password string) (string, error) {
	if password == "" {
		return "", fmt.Errorf("password cannot be empty")
	}

	address := GetAddressFromPrivateKey(privateKey)

	// Get home directory
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get home directory: %w", err)
	}

	// Create keystore directory if it doesn't exist
	keystoreDir := filepath.Join(homeDir, KeyStoreV3Dir)
	if err := os.MkdirAll(keystoreDir, 0700); err != nil {
		return "", fmt.Errorf("failed to create keystore directory: %w", err)
	}

	// Create temporary keystore to generate the encrypted key
	tempDir, err := os.MkdirTemp("", "blocowallet_temp")
	if err != nil {
		return "", fmt.Errorf("failed to create temp directory: %w", err)
	}
	defer os.RemoveAll(tempDir)

	// Create temporary keystore
	tempKeystore := keystore.NewKeyStore(tempDir, keystore.StandardScryptN, keystore.StandardScryptP)

	// Import the key to generate encrypted keystore file
	account, err := tempKeystore.ImportECDSA(privateKey, password)
	if err != nil {
		return "", fmt.Errorf("failed to encrypt private key: %w", err)
	}

	// Read the generated keystore file
	keystoreData, err := os.ReadFile(account.URL.Path)
	if err != nil {
		return "", fmt.Errorf("failed to read keystore file: %w", err)
	}

	// Define the final keystore file path
	fileName := fmt.Sprintf("%s.json", address)
	finalPath := filepath.Join(keystoreDir, fileName)

	// Write to final location
	if err := os.WriteFile(finalPath, keystoreData, 0600); err != nil {
		return "", fmt.Errorf("failed to write keystore file: %w", err)
	}

	return finalPath, nil
}

// LoadPrivateKeyFromKeyStoreV3 loads a private key from a KeyStore V3 file
func (s *Service) LoadPrivateKeyFromKeyStoreV3(keystorePath, password string) (*ecdsa.PrivateKey, error) {
	// Read keystore file
	keystoreData, err := os.ReadFile(keystorePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read keystore file: %w", err)
	}

	// Use keystore.DecryptKey directly to decrypt the JSON keystore
	key, err := keystore.DecryptKey(keystoreData, password)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt keystore with password: %w", err)
	}

	return key.PrivateKey, nil
}

// GetMnemonicFromWallet retrieves and decrypts the mnemonic for a wallet
func (s *Service) GetMnemonicFromWallet(wallet *Wallet, password string) (string, error) {
	if wallet.EncryptedMnemonic == "" {
		return "", fmt.Errorf("wallet has no mnemonic (imported from private key)")
	}

	// Decrypt the mnemonic
	mnemonic, err := DecryptMnemonic(wallet.EncryptedMnemonic, password)
	if err != nil {
		return "", fmt.Errorf("failed to decrypt mnemonic: %w", err)
	}

	return mnemonic, nil
}

// GetAllWallets retrieves all wallets (alias for List for backward compatibility)
func (s *Service) GetAllWallets(ctx context.Context) ([]*Wallet, error) {
	return s.List(ctx)
}

// DeleteWalletByAddress deletes a wallet by address and cleans up keystore
func (s *Service) DeleteWalletByAddress(ctx context.Context, address string) error {
	if address == "" {
		return fmt.Errorf("address cannot be empty")
	}

	// Get wallet first to get keystore path
	wallet, err := s.repo.GetByAddress(ctx, address)
	if err != nil {
		return fmt.Errorf("failed to get wallet: %w", err)
	}

	// Delete from database
	if err := s.repo.Delete(ctx, wallet.ID); err != nil {
		return fmt.Errorf("failed to delete wallet from database: %w", err)
	}

	// Clean up keystore file if it exists
	if wallet.KeyStorePath != "" {
		if err := os.Remove(wallet.KeyStorePath); err != nil {
			// Log the error but don't fail the operation
			// The wallet is already deleted from the database
			fmt.Printf("Warning: failed to delete keystore file %s: %v\n", wallet.KeyStorePath, err)
		}
	}

	return nil
}

// RefreshMultiProvider updates the multi-provider with current network configuration
func (s *Service) RefreshMultiProvider(cfg *config.Config) {
	if s.multiProvider != nil {
		s.multiProvider.RefreshProviders(cfg)
	}
}

// GetMultiNetworkBalance gets the balance for a wallet across all active networks
func (s *Service) GetMultiNetworkBalance(ctx context.Context, address string) (*MultiNetworkBalance, error) {
	if address == "" {
		return nil, fmt.Errorf("wallet address cannot be empty")
	}

	if s.multiProvider == nil {
		// Fallback to single provider if multiProvider is not available
		if s.balanceProvider != nil {
			amount, err := s.balanceProvider.GetBalance(ctx, address)
			if err != nil {
				return nil, fmt.Errorf("failed to get balance: %w", err)
			}

			return &MultiNetworkBalance{
				Address: address,
				NetworkBalances: []*NetworkBalance{
					{
						NetworkKey:  "ethereum",
						NetworkName: "Ethereum",
						Amount:      amount,
						Symbol:      s.balanceProvider.GetNetworkSymbol(),
						Decimals:    s.balanceProvider.GetNetworkDecimals(),
					},
				},
				UpdatedAt: time.Now(),
			}, nil
		}
		return nil, fmt.Errorf("no balance provider available")
	}

	// Use multi-provider to get balances from all active networks
	networkBalances := s.multiProvider.GetAllBalances(ctx, address)

	result := &MultiNetworkBalance{
		Address:         address,
		NetworkBalances: make([]*NetworkBalance, 0, len(networkBalances)),
		UpdatedAt:       time.Now(),
	}

	// Convert blockchain.NetworkBalance to wallet.NetworkBalance
	for _, nb := range networkBalances {
		balance := &NetworkBalance{
			NetworkKey:  nb.NetworkKey,
			NetworkName: nb.NetworkName,
			Amount:      nb.Amount,
			Symbol:      nb.Symbol,
			Decimals:    nb.Decimals,
		}

		if nb.Error != nil {
			balance.Error = nb.Error
		}

		result.NetworkBalances = append(result.NetworkBalances, balance)
	}

	return result, nil
}

// SetWalletPassword stores the password for a wallet address in memory
func (s *Service) SetWalletPassword(address, password string) {
	s.passwordMutex.Lock()
	defer s.passwordMutex.Unlock()
	s.passwordCache[address] = password
}

// GetWalletPassword retrieves the cached password for a wallet address
func (s *Service) GetWalletPassword(address string) (string, bool) {
	s.passwordMutex.RLock()
	defer s.passwordMutex.RUnlock()
	password, exists := s.passwordCache[address]
	return password, exists
}

// ClearWalletPassword removes the cached password for a wallet address
func (s *Service) ClearWalletPassword(address string) {
	s.passwordMutex.Lock()
	defer s.passwordMutex.Unlock()
	delete(s.passwordCache, address)
}
