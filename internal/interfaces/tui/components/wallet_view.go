package components

import (
	"blocowallet/internal/domain/entities"
	"blocowallet/internal/usecases"
	"blocowallet/pkg/logger"
	"fmt"
	"go.uber.org/zap"

	tea "github.com/charmbracelet/bubbletea"
)

type WalletView struct {
	wallets  []entities.Wallet
	selected int
	useCase  usecases.WalletUseCase
	logr     logger.Logger
}

func NewWalletView(useCase usecases.WalletUseCase, logr logger.Logger) *WalletView {
	return &WalletView{
		wallets:  []entities.Wallet{},
		selected: 0,
		useCase:  useCase,
		logr:     logr,
	}
}

func (w *WalletView) Init() tea.Cmd {
	// Carregar wallets se necessário
	return nil
}

func (w *WalletView) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			if w.selected > 0 {
				w.selected--
			}
		case "down", "j":
			if w.selected < len(w.wallets)-1 {
				w.selected++
			}
		case "c":
			// Criar nova wallet via TUI
			wallet, err := w.useCase.CreateWallet()
			if err != nil {
				w.logr.Error("Erro ao criar carteira via TUI", zap.Error(err))
				return w, nil
			}
			w.wallets = append(w.wallets, wallet)
			w.logr.Info("Carteira criada via TUI", zap.Any("address", wallet.Address))
		case "q", "ctrl+c":
			return w, tea.Quit
		}
	}
	return w, nil
}

func (w *WalletView) View() string {
	s := "Carteiras:\n\n"
	for i, wallet := range w.wallets {
		cursor := " "
		if i == w.selected {
			cursor = ">"
		}
		s += fmt.Sprintf("%s %s (%s)\n", cursor, wallet.Address, wallet.ID)
	}
	s += "\nUse as teclas ↑/↓ para navegar, 'c' para criar, 'q' para sair.\n"
	return s
}
