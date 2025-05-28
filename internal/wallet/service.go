package wallet

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"os"
	"time"

	"blocowallet/internal/blockchain"

	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/google/uuid"
)

// Service provides wallet business logic
type Service struct {
	repo            Repository
	balanceProvider BalanceProvider
	multiProvider   *blockchain.MultiProvider
	keystore        *keystore.KeyStore
}

// NewService creates a new wallet service
func NewService(repo Repository, balanceProvider BalanceProvider) *Service {
	return &Service{
		repo:            repo,
		balanceProvider: balanceProvider,
	}
}

// NewServiceWithMultiProvider creates a new wallet service with multi-provider support
func NewServiceWithMultiProvider(repo Repository, multiProvider *blockchain.MultiProvider) *Service {
	return &Service{
		repo:          repo,
		multiProvider: multiProvider,
	}
}

// NewServiceWithKeystore creates a new wallet service with keystore support
func NewServiceWithKeystore(repo Repository, balanceProvider BalanceProvider, ks *keystore.KeyStore) *Service {
	return &Service{
		repo:            repo,
		balanceProvider: balanceProvider,
		keystore:        ks,
	}
}

// NewServiceWithMultiProviderAndKeystore creates a new wallet service with multi-provider and keystore support
func NewServiceWithMultiProviderAndKeystore(repo Repository, multiProvider *blockchain.MultiProvider, ks *keystore.KeyStore) *Service {
	return &Service{
		repo:          repo,
		multiProvider: multiProvider,
		keystore:      ks,
	}
}

// Create creates a new wallet
func (s *Service) Create(ctx context.Context, name, address, keystorePath string) (*Wallet, error) {
	if name == "" {
		return nil, fmt.Errorf("wallet name cannot be empty")
	}

	if address == "" {
		return nil, fmt.Errorf("wallet address cannot be empty")
	}

	// Check if wallet with this address already exists
	existing, err := s.repo.GetByAddress(ctx, address)
	if err == nil && existing != nil {
		return nil, fmt.Errorf("wallet with address %s already exists", address)
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
		return nil, fmt.Errorf("failed to create wallet: %w", err)
	}

	return wallet, nil
}

// GetByID retrieves a wallet by ID
func (s *Service) GetByID(ctx context.Context, id string) (*Wallet, error) {
	if id == "" {
		return nil, fmt.Errorf("wallet ID cannot be empty")
	}

	wallet, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get wallet: %w", err)
	}

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
	if name == "" {
		return nil, fmt.Errorf("wallet name cannot be empty")
	}

	if password == "" {
		return nil, fmt.Errorf("password cannot be empty")
	}

	// Generate mnemonic
	mnemonic, err := GenerateMnemonic()
	if err != nil {
		return nil, fmt.Errorf("failed to generate mnemonic: %w", err)
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

	// Create keystore file if keystore is available
	var keystorePath string
	if s.keystore != nil {
		account, err := s.keystore.ImportECDSA(privKey, password)
		if err != nil {
			return nil, fmt.Errorf("failed to import key to keystore: %w", err)
		}
		keystorePath = account.URL.Path
	}

	// Create wallet entity
	wallet := &Wallet{
		ID:           uuid.New().String(),
		Name:         name,
		Address:      address,
		KeyStorePath: keystorePath,
		Mnemonic:     mnemonic,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
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

	// Create keystore file if keystore is available
	var keystorePath string
	if s.keystore != nil {
		account, err := s.keystore.ImportECDSA(privKey, password)
		if err != nil {
			return nil, fmt.Errorf("failed to import key to keystore: %w", err)
		}
		keystorePath = account.URL.Path
	}

	// Create wallet entity
	wallet := &Wallet{
		ID:           uuid.New().String(),
		Name:         name,
		Address:      address,
		KeyStorePath: keystorePath,
		Mnemonic:     mnemonic,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
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
