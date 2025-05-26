package cli

import (
	"blocowallet/internal/usecases"
	"blocowallet/pkg/logger"
	"github.com/spf13/cobra"
)

// Execute configura e executa a CLI.
func Execute(useCase usecases.WalletUseCase, balanceUseCase usecases.BalanceUseCase, logr logger.Logger) error {
	var rootCmd = &cobra.Command{
		Use:   "bloco",
		Short: "BLOCO Wallet Manager CLI",
		Long:  `Uma interface de linha de comando para gerenciar carteiras de criptomoedas compatíveis com Ethereum.`,
		Run: func(cmd *cobra.Command, args []string) {
			err := cmd.Help()
			if err != nil {
				return
			}
		},
	}
	rootCmd.AddCommand(NewCreateCmd(useCase, logr))
	rootCmd.AddCommand(NewBalanceCmd(balanceUseCase, logr))

	return rootCmd.Execute()
}
