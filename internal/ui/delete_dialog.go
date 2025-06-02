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
		confirmButton = activeButtonStyle.Render("Yes, Delete")
		cancelButton = buttonStyle.Render("Cancel")
	} else {
		confirmButton = buttonStyle.Render("Yes, Delete")
		cancelButton = activeButtonStyle.Render("Cancel")
	}

	// Criar a pergunta com o nome da wallet
	questionText := lipgloss.JoinVertical(
		lipgloss.Center,
		m.question,
		"",
		lipgloss.NewStyle().
			Foreground(lipgloss.Color("#ff6b6b")).
			Bold(true).
			Render("\""+m.walletName+"\""),
		"",
		lipgloss.NewStyle().
			Foreground(lipgloss.Color("#8a8a8a")).
			Render(m.address[:10]+"..."+m.address[len(m.address)-10:]),
		"",
		lipgloss.NewStyle().
			Foreground(lipgloss.Color("#666666")).
			Render("This action cannot be undone."),
	)

	question := lipgloss.NewStyle().
		Width(45).
		Align(lipgloss.Center).
		Render(questionText)

	buttons := lipgloss.JoinHorizontal(
		lipgloss.Top,
		zone.Mark(m.id+"cancel", cancelButton),
		"  ",
		zone.Mark(m.id+"confirm", confirmButton),
	)

	return dialogBoxStyle.Render(
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

// Estilos
var (
	dialogBoxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#874BFD")).
			Padding(2, 3).
			Background(lipgloss.Color("#1a1a1a")).
			Foreground(lipgloss.Color("#ffffff")).
			MarginTop(2).
			MarginBottom(2)

	buttonStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFF")).
			Background(lipgloss.Color("#666")).
			Padding(0, 3).
			Margin(0, 1).
			Border(lipgloss.RoundedBorder())

	activeButtonStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#FFF")).
				Background(lipgloss.Color("#874BFD")).
				Padding(0, 3).
				Margin(0, 1).
				Bold(true).
				Border(lipgloss.RoundedBorder())
)
