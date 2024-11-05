//go:build wireinject
// +build wireinject

package main

import (
	"blocowallet/internal/infrastructure/database"
	"blocowallet/internal/interfaces/cli"
	"blocowallet/internal/interfaces/tui"
	"blocowallet/internal/usecases"
	"blocowallet/pkg/config"
	"blocowallet/pkg/logger"
	"github.com/google/wire"
)

func InitializeApp() func() error {
	wire.Build(
		config.LoadConfig,
		logger.NewLogger,
		database.NewPostgresWalletRepository,
		usecases.NewWalletUseCase,
		wire.Struct(new(cli.CLI), "UseCase", "Logger"),
		wire.Struct(new(tui.TUI), "UseCase", "Logger"),
		wire.Bind(new(usecases.WalletUseCase), new(*usecases.walletUseCase)),
	)
	return nil
}
