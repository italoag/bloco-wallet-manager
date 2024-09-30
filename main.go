package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"blocowallet/config"
	"blocowallet/infrastructure"
	"blocowallet/interfaces"
	"blocowallet/localization"
	"blocowallet/usecases"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/go-errors/errors"
)

func main() {
	// Configuração de logging
	logFile, err := os.OpenFile("blocowallet.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		fmt.Printf("Erro ao abrir o arquivo de log: %v\n", err)
		os.Exit(1)
	}
	defer func(logFile *os.File) {
		err := logFile.Close()
		if err != nil {
			fmt.Printf("Erro ao fechar o arquivo de log: %v\n", err)
		}
	}(logFile)
	log.SetOutput(logFile)

	// Obter o diretório home do usuário
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Println(errors.Wrap(err, 0).ErrorStack())
		fmt.Println("Erro ao obter o diretório home do usuário.")
		os.Exit(1)
	}

	// Definir o diretório base para a aplicação
	appDir := filepath.Join(homeDir, ".wallets")

	// Garantir que o diretório base exista
	if _, err := os.Stat(appDir); os.IsNotExist(err) {
		err := os.MkdirAll(appDir, os.ModePerm)
		if err != nil {
			log.Println(errors.Wrap(err, 0).ErrorStack())
			fmt.Println("Erro ao criar o diretório da aplicação.")
			os.Exit(1)
		}
	}

	// Inicializar a configuração
	cfg, err := config.LoadConfig(appDir)
	if err != nil {
		log.Println(errors.Wrap(err, 0).ErrorStack())
		fmt.Println("Erro ao carregar a configuração.")
		os.Exit(1)
	}

	// Inicializar a localização
	err = localization.SetLanguage(cfg.Language, appDir)
	if err != nil {
		log.Println(errors.Wrap(err, 0).ErrorStack())
		fmt.Println("Erro ao carregar os arquivos de localização.")
		os.Exit(1)
	}

	// Garantir que o diretório de wallets exista
	if _, err := os.Stat(cfg.WalletsDir); os.IsNotExist(err) {
		err := os.MkdirAll(cfg.WalletsDir, os.ModePerm)
		if err != nil {
			log.Println(errors.Wrap(err, 0).ErrorStack())
			fmt.Println("Erro ao criar o diretório de wallets.")
			os.Exit(1)
		}
	}

	// Inicializar o repositório
	repo, err := infrastructure.NewSQLiteRepository(cfg.DatabasePath)
	if err != nil {
		log.Println(errors.Wrap(err, 0).ErrorStack())
		fmt.Println("Erro ao inicializar o banco de dados.")
		os.Exit(1)
	}
	defer func() {
		if err := repo.Close(); err != nil {
			log.Println(errors.Wrap(err, 0).ErrorStack())
			fmt.Println("Erro ao fechar o banco de dados.")
		}
	}()

	// Inicializar o keystore
	ks := keystore.NewKeyStore(cfg.WalletsDir, keystore.StandardScryptN, keystore.StandardScryptP)
	service := usecases.NewWalletService(repo, ks)
	model := interfaces.NewCLIModel(service)

	// Iniciar o programa Bubble Tea com tela cheia
	p := tea.NewProgram(model, tea.WithAltScreen())
	if err, _ := p.Run(); err != nil {
		log.Println(errors.Wrap(err, 0).ErrorStack())
		fmt.Printf("Erro ao executar o programa.\n")
		os.Exit(1)
	}
}
