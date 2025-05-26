package cli

import (
	"blocowallet/internal/usecases"
	"blocowallet/pkg/logger"
	"fmt"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

func NewCreateCmd(useCase usecases.WalletUseCase, logr logger.Logger) *cobra.Command {
	return &cobra.Command{
		Use:   "create",
		Short: "Cria uma nova carteira Ethereum",
		Long:  `Inicializa uma nova carteira compatível com Ethereum e a adiciona ao gerenciador.`,
		Run: func(cmd *cobra.Command, args []string) {
			wallet, err := useCase.CreateWallet()
			if err != nil {
				logr.Error("Erro ao criar carteira", zap.Error(err))
				fmt.Println("Erro ao criar carteira:", err)
				return
			}
			logr.Info("Carteira criada com sucesso", zap.Any("address", wallet.Address))
			fmt.Println("Carteira criada com sucesso:", wallet.Address)
		},
	}
}
