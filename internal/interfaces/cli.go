package interfaces

import (
	"blocowallet/internal/constants"
	"blocowallet/internal/usecases"
	"blocowallet/localization"
	"encoding/json"
	"fmt"
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/digitallyserviced/tdfgo/tdf"
	"github.com/go-errors/errors"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type splashMsg struct{}

func NewCLIModel(service *usecases.WalletService) *CLIModel {
	model := &CLIModel{
		Service:      service,
		currentView:  constants.SplashView,
		menuItems:    NewMenu(),
		selectedMenu: 0,
		styles:       createStyles(),
	}

	if err := initializeFont(model); err != nil {
		model.err = err
		return model
	}

	return model
}

func initializeFont(model *CLIModel) error {
	// Carregar a lista de fontes do arquivo JSON
	fonts, err := loadFontsList(constants.ConfigFontsPath)
	if err != nil {
		log.Println("Erro ao carregar a lista de fontes:", err)
		return err
	}

	if len(fonts) == 0 {
		log.Println("A lista de fontes está vazia.")
		return errors.New("a lista de fontes está vazia")
	}

	// Selecionar uma fonte aleatoriamente
	selectedFontName, err := selectRandomFont(fonts)
	if err != nil {
		log.Println("Erro ao selecionar uma fonte aleatoriamente:", err)
		return err
	}

	log.Printf("Fonte selecionada aleatoriamente: %s\n", selectedFontName)

	// Encontrar a fonte utilizando o módulo tdfgo
	fontInfo := tdf.FindFont(selectedFontName)
	if fontInfo == nil {
		log.Printf("Fonte '%s' não encontrada nos diretórios especificados.\n", selectedFontName)
		return errors.New(constants.ErrorFontNotFoundMessage)
	}

	// Carregar a fonte selecionada
	fontFile, err := tdf.LoadFont(fontInfo)
	if err != nil {
		log.Println("Erro ao carregar a fonte:", err)
		return errors.Wrap(err, 0)
	}

	if len(fontFile.Fonts) == 0 {
		log.Printf("Nenhuma fonte carregada de '%s.tdf'\n", selectedFontName)
		return errors.New("nenhuma fonte carregada")
	}

	// Armazenar a informação da fonte selecionada no modelo
	model.selectedFont = &fontFile.Fonts[0]
	model.fontInfo = fontInfo

	return nil
}

type FontsConfig struct {
	Fonts []string `json:"fonts"`
}

func loadFontsList(configPath string) ([]string, error) {
	absPath, err := filepath.Abs(configPath)
	if err != nil {
		return nil, fmt.Errorf("erro ao obter caminho absoluto do config: %v", err)
	}

	data, err := os.ReadFile(absPath)
	if err != nil {
		return nil, fmt.Errorf("erro ao ler o arquivo de configuração das fontes: %v", err)
	}

	var config FontsConfig
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("erro ao deserializar o JSON de fontes: %v", err)
	}

	return config.Fonts, nil
}

func selectRandomFont(fonts []string) (string, error) {
	if len(fonts) == 0 {
		return "", fmt.Errorf("lista de fontes está vazia")
	}

	rand.NewSource(time.Now().UnixNano())
	index := rand.Intn(len(fonts))
	return fonts[index], nil
}

func (m *CLIModel) Init() tea.Cmd {
	return tea.Batch(
		splashCmd(),
		walletCountCmd(m.Service),
	)
}

func splashCmd() tea.Cmd {
	return tea.Tick(constants.SplashDuration, func(t time.Time) tea.Msg {
		return splashMsg{}
	})
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
		if m.currentView == constants.ListWalletsView {
			m.updateTableDimensions()
		}
		return m, nil

	case splashMsg:
		// Transitar para o menu principal após a splash screen
		m.currentView = constants.DefaultView
		// Iniciar o comando para buscar a quantidade de wallets
		return m, walletCountCmd(m.Service)
	case walletCountMsg:
		if msg.err != nil {
			m.err = msg.err
			log.Println("Erro ao buscar a quantidade de wallets:", msg.err)
		} else {
			m.walletCount = msg.count
		}
		return m, nil
	}

	if m.err != nil {
		if _, ok := msg.(tea.KeyMsg); ok {
			m.err = nil
			m.currentView = constants.DefaultView
		}
		return m, nil
	}

	switch m.currentView {
	case constants.SplashView:
		// Nenhuma atualização adicional necessária durante a splash screen
		return m, nil
	case constants.DefaultView:
		return m.updateMenu(msg)
	case constants.CreateWalletView:
		return m.updateCreateWalletPassword(msg)
	case constants.ImportWalletView:
		return m.updateImportWallet(msg)
	case constants.ImportWalletPasswordView:
		return m.updateImportWalletPassword(msg)
	case constants.ListWalletsView:
		return m.updateListWallets(msg)
	case constants.WalletPasswordView:
		return m.updateWalletPassword(msg)
	case constants.WalletDetailsView:
		return m.updateWalletDetails(msg)
	default:
		m.currentView = constants.DefaultView
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

	switch m.currentView {
	case constants.SplashView:
		return m.renderSplash()
	default:
		return m.renderMainView()
	}
}

func (m *CLIModel) getContentView() string {
	switch m.currentView {
	case constants.DefaultView:
		return localization.Labels["welcome_message"]
	case constants.CreateWalletView:
		return m.viewCreateWalletPassword()
	case constants.ImportWalletView:
		return m.viewImportWallet()
	case constants.ImportWalletPasswordView:
		return m.viewImportWalletPassword()
	case constants.ListWalletsView:
		return m.viewListWallets()
	case constants.WalletPasswordView:
		return m.viewWalletPassword()
	case constants.WalletDetailsView:
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
			case tea.KeyCtrlX.String(), "q", localization.Labels["exit"]:
				return m, tea.Quit
			}
		case tea.KeyCtrlX.String(), "q":
			return m, tea.Quit
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
				m.currentView = constants.DefaultView
				return m, nil
			}
			walletDetails, err := m.Service.CreateWallet(password)
			if err != nil {
				m.err = errors.Wrap(err, 0)
				log.Println(m.err.(*errors.Error).ErrorStack())
				m.currentView = constants.DefaultView
				return m, nil
			}
			m.walletDetails = walletDetails
			m.currentView = constants.WalletDetailsView
		case "esc":
			m.currentView = constants.DefaultView
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
				m.currentView = constants.ImportWalletPasswordView
			}
		case "esc":
			m.currentView = constants.DefaultView
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
				m.currentView = constants.DefaultView
				return m, nil
			}
			mnemonic := strings.Join(m.importWords, " ")
			walletDetails, err := m.Service.ImportWallet(mnemonic, password)
			if err != nil {
				m.err = errors.Wrap(err, 0)
				log.Println(m.err.(*errors.Error).ErrorStack())
				m.currentView = constants.DefaultView
				return m, nil
			}
			m.walletDetails = walletDetails
			m.currentView = constants.WalletDetailsView
		case "esc":
			m.currentView = constants.DefaultView
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
			m.currentView = constants.DefaultView
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
				m.currentView = constants.DefaultView
				return m, nil
			}
			walletDetails, err := m.Service.LoadWallet(m.selectedWallet, password)
			if err != nil {
				m.err = errors.Wrap(err, 0)
				log.Println(m.err.(*errors.Error).ErrorStack())
				m.currentView = constants.DefaultView
				return m, nil
			}
			m.walletDetails = walletDetails
			m.currentView = constants.WalletDetailsView
		case "esc":
			m.currentView = constants.DefaultView
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
			m.currentView = constants.ListWalletsView
		}
	}
	return m, nil
}

func (m *CLIModel) updateTableDimensions() {
	if m.currentView != constants.ListWalletsView || len(m.wallets) == 0 {
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
	m.currentView = constants.CreateWalletView
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
	m.currentView = constants.ImportWalletView
}

func (m *CLIModel) initListWallets() {
	wallets, err := m.Service.GetAllWallets()
	if err != nil {
		m.err = errors.Wrap(fmt.Errorf(localization.Labels["error_loading_wallets"], err), 0)
		log.Println(m.err.(*errors.Error).ErrorStack())
		m.currentView = constants.DefaultView
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

	m.currentView = constants.ListWalletsView
}

func (m *CLIModel) initWalletPassword() {
	m.passwordInput = textinput.New()
	m.passwordInput.Placeholder = localization.Labels["enter_wallet_password"]
	m.passwordInput.CharLimit = constants.PasswordCharLimit
	m.passwordInput.Width = constants.PasswordWidth
	m.passwordInput.EchoMode = textinput.EchoPassword
	m.passwordInput.EchoCharacter = '•'
	m.passwordInput.Focus()
	m.currentView = constants.WalletPasswordView
}
