package cli

import (
	"blocowallet/internal/usecases"
	"blocowallet/pkg/logger"
	"fmt"

	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

func NewBalanceCmd(useCase usecases.BalanceUseCase, logr logger.Logger) *cobra.Command {
	var address string

	balanceCmd := &cobra.Command{
		Use:   "balance",
		Short: "Consulta o saldo de uma carteira Ethereum",
		Long:  `Consulta e exibe o saldo de uma carteira Ethereum específica.`,
		Run: func(cmd *cobra.Command, args []string) {
			if address == "" {
				fmt.Println("O endereço da carteira é obrigatório. Use --address para especificá-lo.")
				return
			}

			balance, err := useCase.GetBalance(address)
			if err != nil {
				logr.Error("Erro ao consultar saldo", zap.String("address", address), zap.Error(err))
				fmt.Println("Erro ao consultar saldo:", err)
				return
			}

			logr.Info("Saldo consultado com sucesso", zap.String("address", address), zap.Float64("balance", balance))
			fmt.Printf("Saldo da carteira %s: %.4f ETH\n", address, balance)
		},
	}

	balanceCmd.Flags().StringVarP(&address, "address", "a", "", "Endereço da carteira Ethereum")

	return balanceCmd
}
