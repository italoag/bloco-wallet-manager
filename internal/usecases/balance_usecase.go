package usecases

// BalanceRepository define as operações necessárias para consultar o saldo.
type BalanceRepository interface {
	GetBalance(address string) (float64, error)
}

// BalanceUseCase define os métodos de negócios para consultar saldo.
type BalanceUseCase interface {
	GetBalance(address string) (float64, error)
}

type balanceUseCase struct {
	balanceRepo BalanceRepository
}

func NewBalanceUseCase(repo BalanceRepository) BalanceUseCase {
	return &balanceUseCase{
		balanceRepo: repo,
	}
}

func (u *balanceUseCase) GetBalance(address string) (float64, error) {
	balance, err := u.balanceRepo.GetBalance(address)
	if err != nil {
		return 0, err
	}
	return balance, nil
}
