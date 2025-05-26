package usecases_test

import (
	"blocowallet/internal/domain/entities"
	"blocowallet/internal/usecases"
	"errors"
	"testing"
)

type MockWalletRepository struct {
	wallets []entities.Wallet
	err     error
}

func (m *MockWalletRepository) CreateWallet(wallet entities.Wallet) error {
	if m.err != nil {
		return m.err
	}
	m.wallets = append(m.wallets, wallet)
	return nil
}

func (m *MockWalletRepository) GetWallet(address string) (entities.Wallet, error) {
	for _, w := range m.wallets {
		if w.Address == address {
			return w, nil
		}
	}
	return entities.Wallet{}, errors.New("wallet not found")
}

func (m *MockWalletRepository) ListWallets() ([]entities.Wallet, error) {
	return m.wallets, m.err
}

func (m *MockWalletRepository) DeleteWallet(address string) error {
	for i, w := range m.wallets {
		if w.Address == address {
			m.wallets = append(m.wallets[:i], m.wallets[i+1:]...)
			return nil
		}
	}
	return errors.New("wallet not found")
}

func TestCreateWallet(t *testing.T) {
	mockRepo := &MockWalletRepository{}
	useCase := usecases.NewWalletUseCase(mockRepo)

	wallet, err := useCase.CreateWallet()
	if err != nil {
		t.Fatalf("Erro ao criar carteira: %v", err)
	}

	if len(mockRepo.wallets) != 1 {
		t.Fatalf("Esperado 1 carteira, got %d", len(mockRepo.wallets))
	}

	if mockRepo.wallets[0].Address != wallet.Address {
		t.Fatalf("Endereço da carteira não corresponde")
	}
}
