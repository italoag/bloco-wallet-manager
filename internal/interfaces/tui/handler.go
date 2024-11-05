package tui

import (
	"blocowallet/internal/interfaces/tui/components"
	"blocowallet/internal/usecases"
	"blocowallet/pkg/logger"

	tea "github.com/charmbracelet/bubbletea"
)

type TUI struct {
	walletView     *components.WalletView
	balanceView    *components.BalanceView
	currentView    tea.Model
	useCase        usecases.WalletUseCase
	balanceUseCase usecases.BalanceUseCase
	logr           logger.Logger
}

func NewTUI(useCase usecases.WalletUseCase, balanceUseCase usecases.BalanceUseCase, logr logger.Logger) *TUI {
	walletView := components.NewWalletView(useCase, logr)
	balanceView := components.NewBalanceView(balanceUseCase, logr)

	return &TUI{
		walletView:     walletView,
		balanceView:    balanceView,
		currentView:    walletView,
		useCase:        useCase,
		balanceUseCase: balanceUseCase,
		logr:           logr,
	}
}

func (t *TUI) Init() tea.Cmd {
	return t.currentView.Init()
}

func (t *TUI) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	// Atualizar a view atual
	newModel, cmd := t.currentView.Update(msg)

	// Verificar se há necessidade de alternar entre as views
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "b":
			if t.currentView == t.walletView {
				t.currentView = t.balanceView
				return t, nil
			}
		case "backspace":
			if t.currentView == t.balanceView {
				t.currentView = t.walletView
				return t, nil
			}
		}
	}

	t.currentView = newModel

	return t, cmd
}

func (t *TUI) View() string {
	return t.currentView.View()
}

func StartTUI(useCase usecases.WalletUseCase, balanceUseCase usecases.BalanceUseCase, logr logger.Logger) error {
	tui := NewTUI(useCase, balanceUseCase, logr)
	p := tea.NewProgram(tui)
	_, err := p.Run()
	return err
}
