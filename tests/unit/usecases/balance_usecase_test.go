// tests/unit/usecases/balance_usecase_test.go
package usecases_test

import (
	"blocowallet/internal/usecases"
	"errors"
	"testing"
)

// MockBalanceRepository implementa a interface BalanceRepository para testes.
type MockBalanceRepository struct {
	balance float64
	err     error
}

func (m *MockBalanceRepository) GetBalance(address string) (float64, error) {
	if m.err != nil {
		return 0, m.err
	}
	if address == "" {
		return 0, errors.New("endereço inválido")
	}
	return m.balance, nil
}

func TestBalanceUseCase_GetBalance_Success(t *testing.T) {
	mockRepo := &MockBalanceRepository{
		balance: 2.7182,
		err:     nil,
	}
	useCase := usecases.NewBalanceUseCase(mockRepo)

	address := "0xABCDEF1234567890"
	balance, err := useCase.GetBalance(address)
	if err != nil {
		t.Fatalf("Esperado sem erro, mas obteve: %v", err)
	}

	if balance != 2.7182 {
		t.Fatalf("Esperado saldo de 2.7182, mas obteve: %.4f", balance)
	}
}

func TestBalanceUseCase_GetBalance_InvalidAddress(t *testing.T) {
	mockRepo := &MockBalanceRepository{
		balance: 0,
		err:     errors.New("endereço inválido"),
	}
	useCase := usecases.NewBalanceUseCase(mockRepo)

	address := ""
	balance, err := useCase.GetBalance(address)
	if err == nil {
		t.Fatal("Esperado erro, mas obteve nil")
	}

	if balance != 0 {
		t.Fatalf("Esperado saldo de 0, mas obteve: %.4f", balance)
	}
}
