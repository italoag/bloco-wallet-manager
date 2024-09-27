package interfaces

import (
	"fmt"
	"strings"

	"blocowallet/constants"
	"blocowallet/entities"
	"blocowallet/localization"
	"blocowallet/usecases"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	contentStyle = lipgloss.NewStyle().Padding(1).Align(lipgloss.Left)
)

type menuItem struct {
	title       string
	description string
}

func (i menuItem) Title() string       { return i.title }
func (i menuItem) Description() string { return i.description }
func (i menuItem) FilterValue() string { return i.title }

type CLIModel struct {
	Service        *usecases.WalletService
	currentView    string
	menuList       list.Model
	importWords    []string
	importStage    int
	textInputs     []textinput.Model
	wallets        []entities.Wallet
	selectedWallet *entities.Wallet
	err            error
	passwordInput  textinput.Model
	mnemonic       string
	walletTable    table.Model
	width          int
	height         int
}

func NewCLIModel(service *usecases.WalletService) CLIModel {
	menuList := NewMenu()
	return CLIModel{
		Service:     service,
		currentView: "menu",
		menuList:    menuList,
	}
}

func (m *CLIModel) Init() tea.Cmd {
	return nil
}

func (m *CLIModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if m.err != nil {
		switch msg.(type) {
		case tea.KeyMsg:
			m.err = nil
			m.currentView = "menu"
			return m, nil
		}
		return m, nil
	}

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

		// Set menu height
		m.menuList.SetHeight(m.height - 2)

		// Calculate content width
		contentWidth := m.width - constants.MenuWidth - 2
		if contentWidth < 20 {
			contentWidth = 20 // Minimum content width
		}

		// Adjust table size if it exists
		if m.walletTable.Columns() != nil {
			m.walletTable.SetWidth(contentWidth - 4)
			m.walletTable.SetHeight(m.height - 4)
		}
		return m, nil
	}

	switch m.currentView {
	case "menu":
		var cmd tea.Cmd
		m.menuList, cmd = m.menuList.Update(msg)
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case "enter":
				i, ok := m.menuList.SelectedItem().(menuItem)
				if ok {
					switch i.title {
					case localization.Labels["create_new_wallet"]:
						m.initCreateWallet()
					case localization.Labels["import_wallet"]:
						m.initImportWallet()
					case localization.Labels["list_wallets"]:
						m.initListWallets()
					case localization.Labels["exit"]:
						return m, tea.Quit
					}
				}
			}
		}
		return m, cmd
	case "create_wallet_password":
		switch msg := msg.(type) {
		case tea.KeyMsg:
			if msg.String() == "enter" {
				password := strings.TrimSpace(m.passwordInput.Value())
				if len(password) < constants.PasswordMinLength {
					m.err = fmt.Errorf(localization.Labels["password_too_short"])
					m.currentView = "menu"
					return m, nil
				}
				wallet, err := m.Service.CreateWallet(password)
				if err != nil {
					m.err = err
					m.currentView = "menu"
					return m, nil
				}
				m.selectedWallet = wallet
				m.currentView = "wallet_details"
			} else {
				var cmd tea.Cmd
				m.passwordInput, cmd = m.passwordInput.Update(msg)
				return m, cmd
			}
		}
	case "import_wallet":
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case "enter":
				word := strings.TrimSpace(m.textInputs[m.importStage].Value())
				if word == "" {
					m.err = fmt.Errorf(localization.Labels["all_words_required"])
					return m, nil
				}
				m.importWords[m.importStage] = word
				m.textInputs[m.importStage].Blur()
				m.importStage++
				if m.importStage < 12 {
					m.textInputs[m.importStage].Focus()
				} else {
					m.passwordInput = textinput.New()
					m.passwordInput.Placeholder = localization.Labels["enter_password"]
					m.passwordInput.CharLimit = constants.PasswordCharLimit
					m.passwordInput.Width = constants.PasswordWidth
					m.passwordInput.EchoMode = textinput.EchoPassword
					m.passwordInput.EchoCharacter = '•'
					m.passwordInput.Focus()
					m.currentView = "import_wallet_password"
				}
			default:
				var cmd tea.Cmd
				m.textInputs[m.importStage], cmd = m.textInputs[m.importStage].Update(msg)
				return m, cmd
			}
		}
	case "import_wallet_password":
		switch msg := msg.(type) {
		case tea.KeyMsg:
			if msg.String() == "enter" {
				password := strings.TrimSpace(m.passwordInput.Value())
				if len(password) < constants.PasswordMinLength {
					m.err = fmt.Errorf(localization.Labels["password_too_short"])
					m.currentView = "menu"
					return m, nil
				}
				mnemonic := strings.Join(m.importWords, " ")
				wallet, err := m.Service.ImportWallet(mnemonic, password)
				if err != nil {
					m.err = err
					m.currentView = "menu"
					return m, nil
				}
				m.selectedWallet = wallet
				m.currentView = "wallet_details"
			} else {
				var cmd tea.Cmd
				m.passwordInput, cmd = m.passwordInput.Update(msg)
				return m, cmd
			}
		}
	case "list_wallets":
		var cmd tea.Cmd
		m.walletTable, cmd = m.walletTable.Update(msg)
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case "up", "k":
				m.walletTable.MoveUp(1)
			case "down", "j":
				m.walletTable.MoveDown(1)
			case "enter":
				selectedRow := m.walletTable.SelectedRow()
				if len(selectedRow) > 1 {
					for _, w := range m.wallets {
						if w.Address == selectedRow[1] {
							m.selectedWallet = &w
							break
						}
					}
					m.initWalletPassword()
				}
			case "esc":
				m.currentView = "menu"
			}
		}
		return m, cmd
	case "wallet_password":
		switch msg := msg.(type) {
		case tea.KeyMsg:
			if msg.String() == "enter" {
				password := strings.TrimSpace(m.passwordInput.Value())
				if password == "" {
					m.err = fmt.Errorf(localization.Labels["password_cannot_be_empty"])
					m.currentView = "menu"
					return m, nil
				}
				err := m.Service.LoadWallet(m.selectedWallet, password)
				if err != nil {
					m.err = err
					m.currentView = "menu"
					return m, nil
				}
				m.currentView = "wallet_details"
			} else {
				var cmd tea.Cmd
				m.passwordInput, cmd = m.passwordInput.Update(msg)
				return m, cmd
			}
		}
	case "wallet_details":
		switch msg := msg.(type) {
		case tea.KeyMsg:
			if msg.String() == "esc" {
				m.currentView = "list_wallets"
			}
		}
	}
	return m, nil
}

// Initialization functions

func (m *CLIModel) initCreateWallet() {
	mnemonic, err := usecases.GenerateMnemonic()
	if err != nil {
		m.err = err
		m.currentView = "menu"
		return
	}
	m.mnemonic = mnemonic
	m.passwordInput = textinput.New()
	m.passwordInput.Placeholder = localization.Labels["enter_password"]
	m.passwordInput.CharLimit = constants.PasswordCharLimit
	m.passwordInput.Width = constants.PasswordWidth
	m.passwordInput.EchoMode = textinput.EchoPassword
	m.passwordInput.EchoCharacter = '•'
	m.passwordInput.Focus()
	m.currentView = "create_wallet_password"
}

func (m *CLIModel) initImportWallet() {
	m.importWords = make([]string, 12)
	m.importStage = 0
	m.textInputs = make([]textinput.Model, 12)
	for i := 0; i < 12; i++ {
		ti := textinput.New()
		ti.Placeholder = fmt.Sprintf("%s %d", localization.Labels["word"], i+1)
		ti.CharLimit = 32
		ti.Width = 20
		m.textInputs[i] = ti
	}
	m.textInputs[0].Focus()
	m.currentView = "import_wallet"
}

func (m *CLIModel) initListWallets() {
	wallets, err := m.Service.GetAllWallets()
	if err != nil {
		m.err = fmt.Errorf(localization.Labels["error_loading_wallets"], err)
		m.currentView = "menu"
		return
	}
	m.wallets = wallets
	columns := []table.Column{
		{Title: localization.Labels["id"], Width: 5},
		{Title: localization.Labels["ethereum_address"], Width: constants.MenuWidth - 10},
	}
	var rows []table.Row
	for _, w := range m.wallets {
		rows = append(rows, table.Row{fmt.Sprintf("%d", w.ID), w.Address})
	}
	m.walletTable = table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
	)
	m.walletTable.SetWidth(constants.MenuWidth - 2)
	m.walletTable.SetHeight(m.height - 4)

	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true).
		Bold(false)
	s.Selected = s.Selected.
		Foreground(lipgloss.Color("229")).
		Background(lipgloss.Color("57")).
		Bold(false)
	s.Cell = s.Cell.Align(lipgloss.Left)
	m.walletTable.SetStyles(s)
	m.currentView = "list_wallets"
}

func (m *CLIModel) initWalletPassword() {
	m.passwordInput = textinput.New()
	m.passwordInput.Placeholder = localization.Labels["enter_wallet_password"]
	m.passwordInput.CharLimit = constants.PasswordCharLimit
	m.passwordInput.Width = constants.PasswordWidth
	m.passwordInput.EchoMode = textinput.EchoPassword
	m.passwordInput.EchoCharacter = '•'
	m.passwordInput.Focus()
	m.currentView = "wallet_password"
}
