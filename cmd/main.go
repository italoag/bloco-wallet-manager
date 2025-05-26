package main

import (
	"blocowallet/internal/infrastructure/blockchain"
	"blocowallet/internal/infrastructure/database"
	"blocowallet/internal/interfaces/cli"
	"blocowallet/internal/interfaces/tui"
	"blocowallet/internal/usecases"
	"blocowallet/pkg/config"
	"blocowallet/pkg/logger"
	"fmt"
	"go.uber.org/zap"
	"os"
	"path/filepath"
)

func main() {
	// Obter o diretório home do usuário
	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Printf("Erro ao obter o diretório home do usuário: %v\n", err)
		os.Exit(1)
	}

	appDir := filepath.Join(homeDir, ".wallets")

	// Garantir que o diretório base exista
	if _, err := os.Stat(appDir); os.IsNotExist(err) {
		err := os.MkdirAll(appDir, os.ModePerm)
		if err != nil {
			fmt.Printf("Erro ao criar o diretório da aplicação: %v\n", err)
			os.Exit(1)
		}
	}

	// Inicializar configuração
	cfg, err := config.LoadConfig(appDir)
	if err != nil {
		fmt.Printf("Erro ao carregar configuração: %v\n", err)
		os.Exit(1)
	}

	// Inicializar logger
	logr := logger.NewLogger(cfg.LogLevel)

	walletRepo, err := database.NewPostgresWalletRepository(cfg.DatabaseURL)
	if err != nil {
		handleError(logr, "Erro ao inicializar o banco de dados", err)
	}
	defer func() {
		if err := walletRepo.Close(); err != nil {
			logr.Error("Erro ao fechar o banco de dados", zap.Error(err))
		}
	}()

	balanceRepo := blockchain.NewBalanceRepository( /* parâmetros, como endpoint da blockchain */ )

	// Inicializar casos de uso
	walletUseCase := usecases.NewWalletUseCase(walletRepo)
	balanceUseCase := usecases.NewBalanceUseCase(balanceRepo)

	// Decidir entre CLI ou TUI
	if len(os.Args) > 1 {
		// Executa a CLI com os argumentos fornecidos
		if err := cli.Execute(walletUseCase, balanceUseCase, logr); err != nil {
			logr.Error("Erro na execução da CLI", zap.Error(err))
			os.Exit(1)
		}
	} else {
		// Executa a TUI se nenhum argumento for fornecido
		if err := tui.StartTUI(walletUseCase, balanceUseCase, logr); err != nil {
			logr.Error("Erro na execução da TUI", zap.Error(err))
			os.Exit(1)
		}
	}
}

// Função auxiliar para tratar erros usando apenas zap
func handleError(logr logger.Logger, message string, err error) {
	logr.Error(message, zap.Error(err))
	fmt.Println(message)
	os.Exit(1)
}
