package interfaces

import (
	"blocowallet/constants"
	"blocowallet/entities"
	"blocowallet/localization"
	"blocowallet/usecases"
	"bytes"
	"fmt"
	"github.com/arsham/figurine/figurine"
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/go-errors/errors"
	"log"
	"strings"
	"time"
)

type CLIModel struct {
	Service        *usecases.WalletService
	currentView    string
	menuItems      []menuItem
	selectedMenu   int
	importWords    []string
	importStage    int
	textInputs     []textinput.Model
	wallets        []entities.Wallet
	selectedWallet *entities.Wallet
	err            error
	passwordInput  textinput.Model
	mnemonic       string
	walletTable    table.Model // Alterado para valor, não ponteiro
	width          int
	height         int
	walletDetails  *usecases.WalletDetails
	styles         Styles
}

const (
	defaultView              = "menu"
	createWalletView         = "create_wallet_password"
	importWalletView         = "import_wallet"
	importWalletPasswordView = "import_wallet_password"
	listWalletsView          = "list_wallets"
	walletPasswordView       = "wallet_password"
	walletDetailsView        = "wallet_details"
	styleWidth               = 40
	styleMargin              = 1
)

type Styles struct {
	Header        lipgloss.Style
	Content       lipgloss.Style
	Footer        lipgloss.Style
	TopStrip      lipgloss.Style
	MenuItem      lipgloss.Style
	MenuSelected  lipgloss.Style
	SelectedTitle lipgloss.Style
	MenuTitle     lipgloss.Style
	MenuDesc      lipgloss.Style
	ErrorStyle    lipgloss.Style
	WalletDetails lipgloss.Style
	StatusBar     lipgloss.Style
}

func NewCLIModel(service *usecases.WalletService) *CLIModel {
	return &CLIModel{
		Service:      service,
		currentView:  defaultView,
		menuItems:    NewMenu(),
		selectedMenu: 0,
		styles:       createStyles(),
	}
}

func createStyles() Styles {
	return Styles{
		Header: lipgloss.NewStyle().
			Align(lipgloss.Left).
			Padding(1, 2),

		Content: lipgloss.NewStyle().
			Align(lipgloss.Left).
			Padding(1, 2),

		Footer: lipgloss.NewStyle().
			Align(lipgloss.Left).
			PaddingLeft(1).
			PaddingRight(1).
			Background(lipgloss.Color("#7D56F4")),

		TopStrip: lipgloss.NewStyle().Margin(1, styleMargin).Padding(0, styleMargin),
		MenuItem: lipgloss.NewStyle().Width(styleWidth).Margin(0, styleMargin).Padding(0, styleMargin),
		MenuSelected: lipgloss.NewStyle().
			Foreground(lipgloss.Color("99")).
			Bold(true).
			Margin(0, styleMargin).
			Padding(0, styleMargin).
			Border(lipgloss.NormalBorder(), false, false, false, true).
			Width(styleWidth),
		SelectedTitle: lipgloss.NewStyle().Bold(true).Margin(0, styleMargin).Padding(0, styleMargin).Foreground(lipgloss.Color("99")),
		MenuTitle:     lipgloss.NewStyle().Margin(0, styleMargin).Padding(0, styleMargin).Bold(true),
		MenuDesc:      lipgloss.NewStyle().Margin(0, styleMargin).Padding(0, styleMargin).Width(styleWidth).Foreground(lipgloss.Color("244")),
		ErrorStyle:    lipgloss.NewStyle().Padding(1, 2).Margin(1, styleMargin),
		WalletDetails: lipgloss.NewStyle().Margin(1, styleMargin).Padding(1, 2),
		StatusBar:     lipgloss.NewStyle().Foreground(lipgloss.Color("241")).Padding(0, styleMargin),
	}
}

func (m *CLIModel) Init() tea.Cmd {
	return nil
}

func (m *CLIModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if msg == nil {
		return m, nil
	}

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height

		// Atualizar estilos com novas dimensões
		m.styles.Header = m.styles.Header.Width(m.width)
		m.styles.Content = m.styles.Content.Width(m.width)
		m.styles.Footer = m.styles.Footer.Width(m.width)

		// Atualizar dimensões da tabela
		if m.currentView == listWalletsView {
			m.updateTableDimensions()
		}
		return m, nil
	}

	if m.err != nil {
		if _, ok := msg.(tea.KeyMsg); ok {
			m.err = nil
			m.currentView = defaultView
		}
		return m, nil
	}

	switch m.currentView {
	case defaultView:
		return m.updateMenu(msg)
	case createWalletView:
		return m.updateCreateWalletPassword(msg)
	case importWalletView:
		return m.updateImportWallet(msg)
	case importWalletPasswordView:
		return m.updateImportWalletPassword(msg)
	case listWalletsView:
		return m.updateListWallets(msg)
	case walletPasswordView:
		return m.updateWalletPassword(msg)
	case walletDetailsView:
		return m.updateWalletDetails(msg)
	default:
		m.currentView = defaultView
		return m, nil
	}
}

func (m *CLIModel) renderMenuItems() []string {
	var menuItems []string
	for i, item := range m.menuItems {
		style := m.styles.MenuItem
		titleStyle := m.styles.MenuTitle
		if i == m.selectedMenu {
			style = m.styles.MenuSelected
			titleStyle = m.styles.SelectedTitle
		}
		menuText := fmt.Sprintf("%s\n%s", titleStyle.Render(item.title), m.styles.MenuDesc.Render(item.description))
		menuItems = append(menuItems, style.Render(menuText))
	}

	numRows := (len(menuItems) + 1) / 2
	var menuRows []string
	for i := 0; i < numRows; i++ {
		startIndex := i * 2
		endIndex := startIndex + 2
		if endIndex > len(menuItems) {
			endIndex = len(menuItems)
		}
		row := lipgloss.JoinHorizontal(lipgloss.Top, menuItems[startIndex:endIndex]...)
		menuRows = append(menuRows, row)
	}
	return menuRows
}

func (m *CLIModel) View() string {
	if m.err != nil {
		return m.styles.ErrorStyle.Render(fmt.Sprintf(localization.Labels["error_message"], m.err))
	}

	// Preparar conteúdo do header
	//banner := figurine.Write(io.Writer(os.Stdout), "bloco", "Test1.plf")
	// Preparar conteúdo do header com o logo colorido usando figurine
	var logoBuffer bytes.Buffer
	err := figurine.Write(&logoBuffer, "bloco", "Test1.flf")
	if err != nil {
		log.Println(errors.Wrap(err, 0))
		// Fallback para logo sem estilização
		logoBuffer.WriteString("bloco")
	}
	logo := logoBuffer.String()
	//logo := figure.NewColorFigure("bloco", "speed", "pink", true)
	walletCount := len(m.wallets)
	currentTime := time.Now().Format("2006-01-02 15:04:05")

	headerLeft := lipgloss.JoinVertical(
		lipgloss.Left,
		logo,
		fmt.Sprintf("Wallets: %d", walletCount),
		fmt.Sprintf("Date: %s", currentTime),
		fmt.Sprintf("Version: %s", localization.Labels["version"]),
	)

	menuItems := m.renderMenuItems()
	menuGrid := lipgloss.JoinVertical(lipgloss.Left, menuItems...)

	// Montar header
	headerContent := lipgloss.JoinHorizontal(
		lipgloss.Top,
		headerLeft,
		lipgloss.NewStyle().Width(m.width-lipgloss.Width(headerLeft)-lipgloss.Width(menuGrid)).Render(""),
		menuGrid,
	)

	// Renderizar header com altura fixa
	renderedHeader := m.styles.Header.Render(headerContent)
	headerHeight := lipgloss.Height(renderedHeader)

	// Preparar conteúdo do footer
	statusBar := fmt.Sprintf("Current view: %s | Wallets: %d", localization.Labels[m.currentView], walletCount)
	renderedFooter := m.styles.Footer.Render(statusBar)
	footerHeight := lipgloss.Height(renderedFooter)

	// Calcular altura da área de conteúdo
	contentHeight := m.height - headerHeight - footerHeight - 2 // Subtrai 2 para evitar overflow

	if contentHeight < 0 {
		contentHeight = 0
	}

	// Obter a visualização do conteúdo
	content := m.getContentView()

	// Renderizar conteúdo com altura ajustada
	renderedContent := m.styles.Content.Height(contentHeight).Render(content)

	// Inserir espaço vazio para empurrar o footer para baixo
	remainingHeight := m.height - headerHeight - lipgloss.Height(renderedContent) - footerHeight
	if remainingHeight < 0 {
		remainingHeight = 0
	}
	emptySpace := lipgloss.NewStyle().Height(remainingHeight).Render("")

	// Montar a visualização final
	finalView := lipgloss.JoinVertical(
		lipgloss.Top,
		renderedHeader,
		renderedContent,
		emptySpace,
		renderedFooter,
	)

	return finalView
}

func (m *CLIModel) getContentView() string {
	switch m.currentView {
	case defaultView:
		return localization.Labels["welcome_message"]
	case createWalletView:
		return m.viewCreateWalletPassword()
	case importWalletView:
		return m.viewImportWallet()
	case importWalletPasswordView:
		return m.viewImportWalletPassword()
	case listWalletsView:
		return m.viewListWallets()
	case walletPasswordView:
		return m.viewWalletPassword()
	case walletDetailsView:
		return m.viewWalletDetails()
	default:
		return localization.Labels["unknown_state"]
	}
}

func (m *CLIModel) updateMenu(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			if m.selectedMenu > 0 {
				m.selectedMenu--
			}
		case "down", "j":
			if m.selectedMenu < len(m.menuItems)-1 {
				m.selectedMenu++
			}
		case "left", "h":
			if m.selectedMenu > 1 {
				m.selectedMenu -= 2
			}
		case "right", "l":
			if m.selectedMenu < len(m.menuItems)-2 {
				m.selectedMenu += 2
			}
		case "enter":
			switch m.menuItems[m.selectedMenu].title {
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
	return m, nil
}

func (m *CLIModel) updateCreateWalletPassword(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			password := strings.TrimSpace(m.passwordInput.Value())
			if len(password) < constants.PasswordMinLength {
				m.err = errors.Wrap(fmt.Errorf(localization.Labels["password_too_short"]), 0)
				log.Println(m.err.(*errors.Error).ErrorStack())
				m.currentView = defaultView
				return m, nil
			}
			walletDetails, err := m.Service.CreateWallet(password)
			if err != nil {
				m.err = errors.Wrap(err, 0)
				log.Println(m.err.(*errors.Error).ErrorStack())
				m.currentView = defaultView
				return m, nil
			}
			m.walletDetails = walletDetails
			m.currentView = walletDetailsView
		default:
			var cmd tea.Cmd
			m.passwordInput, cmd = m.passwordInput.Update(msg)
			return m, cmd
		}
	}
	return m, nil
}

func (m *CLIModel) updateImportWallet(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			word := strings.TrimSpace(m.textInputs[m.importStage].Value())
			if word == "" {
				m.err = errors.Wrap(fmt.Errorf(localization.Labels["all_words_required"]), 0)
				log.Println(m.err.(*errors.Error).ErrorStack())
				return m, nil
			}
			m.importWords[m.importStage] = word
			m.textInputs[m.importStage].Blur()
			m.importStage++
			if m.importStage < len(m.textInputs) {
				m.textInputs[m.importStage].Focus()
			} else {
				m.passwordInput = textinput.New()
				m.passwordInput.Placeholder = localization.Labels["enter_password"]
				m.passwordInput.CharLimit = constants.PasswordCharLimit
				m.passwordInput.Width = constants.PasswordWidth
				m.passwordInput.EchoMode = textinput.EchoPassword
				m.passwordInput.EchoCharacter = '•'
				m.passwordInput.Focus()
				m.currentView = importWalletPasswordView
			}
		case "esc":
			m.currentView = defaultView
		default:
			var cmd tea.Cmd
			m.textInputs[m.importStage], cmd = m.textInputs[m.importStage].Update(msg)
			return m, cmd
		}
	}
	return m, nil
}

func (m *CLIModel) updateImportWalletPassword(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			password := strings.TrimSpace(m.passwordInput.Value())
			if len(password) < constants.PasswordMinLength {
				m.err = errors.Wrap(fmt.Errorf(localization.Labels["password_too_short"]), 0)
				log.Println(m.err.(*errors.Error).ErrorStack())
				m.currentView = defaultView
				return m, nil
			}
			mnemonic := strings.Join(m.importWords, " ")
			walletDetails, err := m.Service.ImportWallet(mnemonic, password)
			if err != nil {
				m.err = errors.Wrap(err, 0)
				log.Println(m.err.(*errors.Error).ErrorStack())
				m.currentView = defaultView
				return m, nil
			}
			m.walletDetails = walletDetails
			m.currentView = walletDetailsView
		case "esc":
			m.currentView = defaultView
		default:
			var cmd tea.Cmd
			m.passwordInput, cmd = m.passwordInput.Update(msg)
			return m, cmd
		}
	}
	return m, nil
}

func (m *CLIModel) updateListWallets(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			selectedRow := m.walletTable.SelectedRow()
			if len(selectedRow) > 1 {
				address := selectedRow[1]
				// Buscar wallet pela address
				for _, w := range m.wallets {
					if w.Address == address {
						m.selectedWallet = &w
						m.initWalletPassword()
						return m, nil
					}
				}
			}
		case "esc":
			m.currentView = defaultView
			return m, nil
		}
	}

	// Atualizar a tabela com a mensagem
	var cmd tea.Cmd
	m.walletTable, cmd = m.walletTable.Update(msg)
	return m, cmd
}

func (m *CLIModel) updateWalletPassword(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			password := strings.TrimSpace(m.passwordInput.Value())
			if password == "" {
				m.err = errors.Wrap(fmt.Errorf(localization.Labels["password_cannot_be_empty"]), 0)
				log.Println(m.err.(*errors.Error).ErrorStack())
				m.currentView = defaultView
				return m, nil
			}
			walletDetails, err := m.Service.LoadWallet(m.selectedWallet, password)
			if err != nil {
				m.err = errors.Wrap(err, 0)
				log.Println(m.err.(*errors.Error).ErrorStack())
				m.currentView = defaultView
				return m, nil
			}
			m.walletDetails = walletDetails
			m.currentView = walletDetailsView
		case "esc":
			m.currentView = defaultView
		default:
			var cmd tea.Cmd
			m.passwordInput, cmd = m.passwordInput.Update(msg)
			return m, cmd
		}
	}
	return m, nil
}

func (m *CLIModel) updateWalletDetails(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			m.walletDetails = nil
			m.currentView = listWalletsView
		}
	}
	return m, nil
}

func (m *CLIModel) updateTableDimensions() {
	if m.currentView != listWalletsView || len(m.wallets) == 0 {
		return
	}

	contentAreaHeight := m.height - lipgloss.Height(m.styles.Header.Render("")) - lipgloss.Height(m.styles.Footer.Render("")) - 4
	if contentAreaHeight < 0 {
		contentAreaHeight = 0
	}

	m.walletTable.SetWidth(m.width - 4)
	m.walletTable.SetHeight(contentAreaHeight - 2) // Subtrai 2 para evitar overflow

	// Calcular larguras das colunas
	idColWidth := 10
	addressColWidth := m.width - idColWidth - 8 // Subtrai 8 para padding e margens

	if addressColWidth < 20 {
		addressColWidth = 20
	}

	// Atualizar colunas
	m.walletTable.SetColumns([]table.Column{
		{Title: localization.Labels["id"], Width: idColWidth},
		{Title: localization.Labels["ethereum_address"], Width: addressColWidth},
	})
}

// Funções de inicialização

func (m *CLIModel) initCreateWallet() {
	m.mnemonic, _ = usecases.GenerateMnemonic()
	m.passwordInput = textinput.New()
	m.passwordInput.Placeholder = localization.Labels["enter_password"]
	m.passwordInput.CharLimit = constants.PasswordCharLimit // Corrigido para PasswordCharLimit
	m.passwordInput.Width = constants.PasswordWidth
	m.passwordInput.EchoMode = textinput.EchoPassword
	m.passwordInput.EchoCharacter = '•'
	m.passwordInput.Focus()
	m.currentView = createWalletView
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
	m.currentView = importWalletView
}

func (m *CLIModel) initListWallets() {
	wallets, err := m.Service.GetAllWallets()
	if err != nil {
		m.err = errors.Wrap(fmt.Errorf(localization.Labels["error_loading_wallets"], err), 0)
		log.Println(m.err.(*errors.Error).ErrorStack())
		m.currentView = defaultView
		return
	}
	m.wallets = wallets

	// Inicialize as colunas com larguras adequadas
	idColWidth := 10
	addressColWidth := m.width - idColWidth - 8 // Subtrai 8 para padding e margens

	if addressColWidth < 20 {
		addressColWidth = 20
	}

	columns := []table.Column{
		{Title: localization.Labels["id"], Width: idColWidth},
		{Title: localization.Labels["ethereum_address"], Width: addressColWidth},
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

	// Ajustar os estilos da tabela
	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true).
		Bold(true)
	s.Selected = s.Selected.
		Foreground(lipgloss.Color("229")).
		Background(lipgloss.Color("57")).
		Bold(false)
	s.Cell = s.Cell.Align(lipgloss.Left)
	m.walletTable.SetStyles(s)

	// Atualizar dimensões da tabela
	m.updateTableDimensions()

	m.currentView = listWalletsView
}

func (m *CLIModel) initWalletPassword() {
	m.passwordInput = textinput.New()
	m.passwordInput.Placeholder = localization.Labels["enter_wallet_password"]
	m.passwordInput.CharLimit = constants.PasswordCharLimit
	m.passwordInput.Width = constants.PasswordWidth
	m.passwordInput.EchoMode = textinput.EchoPassword
	m.passwordInput.EchoCharacter = '•'
	m.passwordInput.Focus()
	m.currentView = walletPasswordView
}
