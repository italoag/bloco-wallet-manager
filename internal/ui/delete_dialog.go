package ui

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	zone "github.com/lrstanley/bubblezone"
)

type DeleteWalletDialog struct {
	id         string
	width      int
	active     string
	question   string
	walletName string
	address    string
}

func NewDeleteWalletDialog(walletName, address string) DeleteWalletDialog {
	return DeleteWalletDialog{
		id:         "delete-wallet-dialog",
		active:     "cancel", // Começar com "cancel" ativo por segurança
		walletName: walletName,
		address:    address,
		question:   "Are you sure you want to delete this wallet?",
	}
}

func (m DeleteWalletDialog) Init() tea.Cmd {
	return nil
}

func (m DeleteWalletDialog) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
	case tea.MouseMsg:
		if msg.Action != tea.MouseActionRelease || msg.Button != tea.MouseButtonLeft {
			return m, nil
		}

		if zone.Get(m.id + "confirm").InBounds(msg) {
			m.active = "confirm"
			return m, func() tea.Msg { return ConfirmDeleteMsg{} }
		} else if zone.Get(m.id + "cancel").InBounds(msg) {
			m.active = "cancel"
			return m, func() tea.Msg { return CancelDeleteMsg{} }
		}

		return m, nil
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, key.NewBinding(key.WithKeys("left", "h"))):
			if m.active == "confirm" {
				m.active = "cancel"
			} else {
				m.active = "confirm"
			}
		case key.Matches(msg, key.NewBinding(key.WithKeys("right", "l"))):
			if m.active == "cancel" {
				m.active = "confirm"
			} else {
				m.active = "cancel"
			}
		case key.Matches(msg, key.NewBinding(key.WithKeys("enter"))):
			if m.active == "confirm" {
				return m, func() tea.Msg { return ConfirmDeleteMsg{} }
			} else {
				return m, func() tea.Msg { return CancelDeleteMsg{} }
			}
		case key.Matches(msg, key.NewBinding(key.WithKeys("esc"))):
			return m, func() tea.Msg { return CancelDeleteMsg{} }
		}
	}
	return m, nil
}

func (m DeleteWalletDialog) View() string {
	var confirmButton, cancelButton string

	if m.active == "confirm" {
		confirmButton = ActiveButtonStyle.Render("Yes, Delete")
		cancelButton = ButtonStyle.Render("Cancel")
	} else {
		confirmButton = ButtonStyle.Render("Yes, Delete")
		cancelButton = ActiveButtonStyle.Render("Cancel")
	}

	// Criar a pergunta com o nome da wallet
	questionText := lipgloss.JoinVertical(
		lipgloss.Center,
		m.question,
		"",
		WalletNameStyle.Render("\""+m.walletName+"\""),
		"",
		AddressStyle.Render(m.address[:10]+"..."+m.address[len(m.address)-10:]),
		"",
		WarningStyle.Render("This action cannot be undone."),
	)

	question := DialogQuestionStyle.Render(questionText)

	buttons := lipgloss.JoinHorizontal(
		lipgloss.Top,
		zone.Mark(m.id+"cancel", cancelButton),
		"  ",
		zone.Mark(m.id+"confirm", confirmButton),
	)

	return DialogBoxStyle.Render(
		lipgloss.JoinVertical(
			lipgloss.Center,
			question,
			"",
			buttons,
		),
	)
}

// Mensagens para comunicação
type ConfirmDeleteMsg struct{}
type CancelDeleteMsg struct{}
