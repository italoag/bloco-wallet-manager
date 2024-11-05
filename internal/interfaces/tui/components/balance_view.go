package components

import (
	"blocowallet/internal/usecases"
	"blocowallet/pkg/logger"
	"fmt"
	"go.uber.org/zap"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type BalanceView struct {
	address string
	balance float64
	input   string
	ready   bool
	useCase usecases.BalanceUseCase
	logr    logger.Logger
}

func NewBalanceView(useCase usecases.BalanceUseCase, logr logger.Logger) *BalanceView {
	return &BalanceView{
		address: "",
		balance: 0,
		input:   "",
		ready:   false,
		useCase: useCase,
		logr:    logr,
	}
}

func (b *BalanceView) Init() tea.Cmd {
	return nil
}

func (b *BalanceView) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "b":
			// Iniciar consulta de saldo
			b.ready = true
			b.input = ""
		case "enter":
			if b.ready {
				b.address = b.input
				balance, err := b.useCase.GetBalance(b.address)
				if err != nil {
					b.logr.Error("Erro ao consultar saldo via TUI", zap.Error(err))
					b.balance = 0
					b.ready = false
					return b, nil
				}
				b.balance = balance
				b.ready = false
			}
		case "esc":
			if b.ready {
				b.ready = false
				b.input = ""
			}
		case "ctrl+c", "q":
			return b, tea.Quit
		default:
			if b.ready {
				b.input += msg.String()
			}
		}
	}
	return b, nil
}

func (b *BalanceView) View() string {
	s := "Pressione 'b' para consultar o saldo de uma carteira.\n\n"

	if b.ready {
		s += "Digite o endereço da carteira:\n" + b.input + "\n"
	} else if b.address != "" {
		s += fmt.Sprintf("Saldo da carteira %s: %.4f ETH\n", b.address, b.balance)
	}

	s += "\nPressione 'q' para sair.\n"

	// Estilização opcional com Lip Gloss
	titleStyle := lipgloss.NewStyle().Bold(true).Underline(true)
	return titleStyle.Render("BLOCO Wallet Manager") + "\n\n" + s
}
