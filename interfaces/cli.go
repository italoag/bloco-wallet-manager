package interfaces

import (
	"blocowallet/constants"
	"blocowallet/localization"
	"blocowallet/usecases"
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/digitallyserviced/tdfgo/tdf"
	"github.com/go-errors/errors"
)

// Função para construir a lista de fontes disponíveis tanto do diretório customizado quanto das embutidas
func buildFontsList(customFontDir string) []*tdf.FontInfo {
	var fonts []*tdf.FontInfo

	// Primeiro, tenta adicionar fontes do diretório personalizado, se existir
	if customFontDir != "" {
		if _, err := os.Stat(customFontDir); err == nil {
			// Adicionar fontes do diretório personalizado
			files, err := os.ReadDir(customFontDir)
			if err == nil {
				for _, file := range files {
					if !file.IsDir() && strings.HasSuffix(strings.ToLower(file.Name()), ".tdf") {
						fontPath := filepath.Join(customFontDir, file.Name())
						fontInfo := tdf.NewFontInfo(file.Name(), fontPath)
						fontInfo.FontDir = customFontDir
						fonts = append(fonts, fontInfo)
					}
				}
			}
		}
	}

	// Se nenhuma fonte foi encontrada no diretório personalizado ou se ele não existe,
	// usar as fontes embutidas
	if len(fonts) == 0 {
		builtinFonts := tdf.SearchBuiltinFonts("*")
		fonts = append(fonts, builtinFonts...)
	}

	return fonts
}

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
	// Obter o diretório home do usuário
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return errors.Wrap(err, 0)
	}

	// Definir o diretório da aplicação
	appDir := filepath.Join(homeDir, ".wallets")

	// Definir o diretório de fontes personalizado
	customFontDir := filepath.Join(appDir, "config", "fonts")

	// Verificar se o diretório de fontes personalizado existe
	if _, err := os.Stat(customFontDir); err != nil {
		// Se não existir, tentar criar o diretório
		if os.IsNotExist(err) {
			err = os.MkdirAll(customFontDir, os.ModePerm)
			if err != nil {
				log.Printf("Não foi possível criar o diretório de fontes personalizado: %v\n", err)
				// Continuar com as fontes embutidas
				customFontDir = ""
			}
		} else {
			log.Printf("Erro ao verificar o diretório de fontes personalizado: %v\n", err)
			customFontDir = ""
		}
	}

	// Construir a lista de fontes disponíveis (personalizadas + embutidas)
	availableFonts := buildFontsList(customFontDir)

	if len(availableFonts) == 0 {
		return errors.New("nenhuma fonte disponível, nem personalizada nem embutida")
	}

	// Carregar nomes das fontes configuradas
	configuredFontNames, err := loadFontsList(appDir)
	if err != nil {
		log.Println("Erro ao carregar a lista de fontes configuradas:", err)
		// Se houver erro, escolher qualquer fonte disponível
		rand.NewSource(time.Now().UnixNano())
		selectedFontInfo := availableFonts[rand.Intn(len(availableFonts))]
		return loadSelectedFont(model, selectedFontInfo)
	}

	// Se não houver fontes configuradas, escolher qualquer fonte disponível
	if len(configuredFontNames) == 0 {
		log.Println("A lista de fontes configuradas está vazia, selecionando aleatoriamente.")
		rand.NewSource(time.Now().UnixNano())
		selectedFontInfo := availableFonts[rand.Intn(len(availableFonts))]
		return loadSelectedFont(model, selectedFontInfo)
	}

	// Selecionar uma fonte da lista configurada
	selectedName, err := selectRandomFont(configuredFontNames)
	if err != nil {
		log.Println("Erro ao selecionar uma fonte aleatoriamente:", err)
		// Selecionar qualquer fonte disponível como fallback
		rand.NewSource(time.Now().UnixNano())
		selectedFontInfo := availableFonts[rand.Intn(len(availableFonts))]
		return loadSelectedFont(model, selectedFontInfo)
	}

	// Procurar a fonte selecionada nas fontes disponíveis
	var selectedFontInfo *tdf.FontInfo
	for _, fontInfo := range availableFonts {
		baseName := strings.TrimSuffix(fontInfo.File, ".tdf")
		if strings.EqualFold(baseName, selectedName) {
			selectedFontInfo = fontInfo
			break
		}
	}

	// Se não encontrada, usar qualquer fonte disponível como fallback
	if selectedFontInfo == nil {
		log.Printf("Fonte '%s' não encontrada. Usando uma fonte aleatória como fallback.\n", selectedName)
		rand.NewSource(time.Now().UnixNano())
		selectedFontInfo = availableFonts[rand.Intn(len(availableFonts))]
	}

	return loadSelectedFont(model, selectedFontInfo)
}

// Função auxiliar para carregar a fonte selecionada
func loadSelectedFont(model *CLIModel, fontInfo *tdf.FontInfo) error {
	// Carregar a fonte selecionada
	fontFile, err := tdf.LoadFont(fontInfo)
	if err != nil {
		log.Println("Erro ao carregar a fonte:", err)
		return errors.Wrap(err, 0)
	}

	if len(fontFile.Fonts) == 0 {
		log.Printf("Nenhuma fonte carregada de '%s'\n", fontInfo.File)
		return errors.New("nenhuma fonte carregada")
	}

	// Armazenar a informação da fonte selecionada no modelo
	model.selectedFont = &fontFile.Fonts[0]
	model.fontInfo = fontInfo

	log.Printf("Fonte carregada com sucesso: %s\n", fontInfo.File)
	return nil
}

type FontsConfig struct {
	Fonts []string `json:"fonts"`
}

func loadFontsList(appDir string) ([]string, error) {
	// Construir o caminho correto para o arquivo de configuração
	configPath := filepath.Join(appDir, "config", "fonts.json")

	// Verificar se o diretório config existe, se não, criar
	configDir := filepath.Dir(configPath)
	if _, err := os.Stat(configDir); os.IsNotExist(err) {
		err := os.MkdirAll(configDir, os.ModePerm)
		if err != nil {
			return nil, fmt.Errorf("erro ao criar o diretório de configuração: %v", err)
		}
	}

	// Verificar se o arquivo de configuração existe, se não, criar com valores padrão
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		// Tentar ler do caminho relativo (para o caso de desenvolvimento)
		relativeConfigPath := filepath.Join("config", "fonts.json")
		data, err := os.ReadFile(relativeConfigPath)
		if err == nil {
			// Escrever para o caminho correto
			err = os.WriteFile(configPath, data, 0644)
			if err != nil {
				return nil, fmt.Errorf("erro ao criar o arquivo de configuração das fontes: %v", err)
			}
		} else {
			// Se não conseguir ler do relativo, criar um config padrão
			defaultFonts := FontsConfig{
				Fonts: []string{
					"1911", "dynasty", "etherx", "commx", "intensex",
					"icex", "royfour", "phudge", "portal", "wild",
				},
			}
			data, err := json.Marshal(defaultFonts)
			if err != nil {
				return nil, fmt.Errorf("erro ao criar configuração padrão: %v", err)
			}
			err = os.WriteFile(configPath, data, 0644)
			if err != nil {
				return nil, fmt.Errorf("erro ao criar o arquivo de configuração das fontes: %v", err)
			}
		}
	}

	// Agora ler o arquivo
	data, err := os.ReadFile(configPath)
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

	// Tratar as teclas de navegação global (esc/backspace) antes de qualquer outro processamento
	if keyMsg, ok := msg.(tea.KeyMsg); ok {
		switch keyMsg.String() {
		case "esc", "backspace":
			// Se estiver na tela de lista de wallets e tiver um diálogo de exclusão aberto,
			// não faça nada aqui e deixe o handler específico da view tratar
			if m.currentView == constants.ListWalletsView && m.deletingWallet != nil {
				// Não faz nada, deixa o handler específico tratar
			} else if m.currentView != constants.DefaultView && m.currentView != constants.SplashView {
				// Para a maioria das telas, voltar para o menu principal
				if m.currentView == constants.WalletDetailsView {
					// Comportamento específico para tela de detalhes: voltar para lista de wallets
					m.walletDetails = nil
					m.currentView = constants.ListWalletsView
				} else {
					// Comportamento padrão: voltar ao menu principal
					m.menuItems = NewMenu()
					m.selectedMenu = 0
					m.currentView = constants.DefaultView
				}
				// Sempre retorne imediatamente após processar a tecla de navegação
				return m, nil
			}
		case "q":
			if m.currentView != constants.SplashView {
				return m, tea.Quit
			}
		}
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

	case walletsRefreshedMsg:
		// Apenas retornar o modelo sem fazer nada, pois a atualização já foi feita
		// Isso evita que a tela inteira seja redesenhada
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

	// Processamento específico para cada tela
	switch m.currentView {
	case constants.SplashView:
		// Nenhuma atualização adicional necessária durante a splash screen
		return m, nil
	case constants.DefaultView:
		return m.updateMenu(msg)
	case constants.CreateWalletView:
		return m.updateCreateWalletPassword(msg)
	case constants.ImportMethodSelectionView:
		return m.updateImportMethodSelection(msg)
	case constants.ImportWalletView:
		return m.updateImportWallet(msg)
	case constants.ImportPrivateKeyView:
		return m.updateImportPrivateKey(msg)
	case constants.ImportWalletPasswordView:
		return m.updateImportWalletPassword(msg)
	case constants.ListWalletsView:
		return m.updateListWallets(msg)
	case constants.WalletPasswordView:
		return m.updateWalletPassword(msg)
	case constants.SelectWalletOperationView:
		return m.updateSelectWalletOperation(msg)
	case constants.WalletDetailsView:
		return m.updateWalletDetails(msg)
	case constants.DeleteWalletConfirmationView:
		return m.updateDeleteWalletConfirmation(msg)
	default:
		m.currentView = constants.DefaultView
		return m, nil
	}
}

func (m *CLIModel) View() string {
	if m.err != nil {
		return m.styles.ErrorStyle.Render(fmt.Sprintf(localization.Labels["error_message"], m.err))
	}

	switch m.currentView {
	case constants.SplashView:
		return m.renderSplash()
	case constants.ListWalletsView:
		// Tratamento especial para a visualização de listagem de carteiras
		// para garantir que ela se encaixe corretamente no layout
		return m.renderListWalletsWithLayout()
	default:
		return m.renderMainView()
	}
}

// renderListWalletsWithLayout renderiza a tela de listagem de carteiras com o layout completo
func (m *CLIModel) renderListWalletsWithLayout() string {
	// Renderizar o cabeçalho da mesma forma que renderMainView
	var logoBuffer bytes.Buffer
	err := figurine.Write(&logoBuffer, "bloco", "Test1.flf")
	if err != nil {
		log.Println(errors.Wrap(err, 0))
		logoBuffer.WriteString("bloco")
	}
	renderedLogo := logoBuffer.String()

	walletCount := m.walletCount
	currentTime := time.Now().Format("02-01-2006 15:04:05")

	headerLeft := lipgloss.JoinVertical(
		lipgloss.Left,
		renderedLogo,
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
	renderedFooter := m.renderStatusBar()
	footerHeight := lipgloss.Height(renderedFooter)

	// Calcular altura disponível para o conteúdo
	contentHeight := m.height - headerHeight - footerHeight - 2

	// Ajustar o tamanho da tabela para caber na área de conteúdo
	if contentHeight > 0 {
		// Reservar espaço para título e instruções
		titleAndInstructionsHeight := 4
		tableHeight := contentHeight - titleAndInstructionsHeight

		if tableHeight > 0 {
			m.walletTable.SetHeight(tableHeight)
		}
	}

	// Obter conteúdo da visualização de carteiras
	content := m.viewListWallets()

	// Renderizar o conteúdo na área apropriada
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

// renderImportMenuItems renderiza os itens do menu de importação
func (m *CLIModel) renderImportMenuItems() string {
	// Criar o menu de importação
	importMenu := NewImportMenu()

	// Renderizar cada item do menu
	var menuItems []string
	for i, item := range importMenu {
		style := m.styles.MenuItem
		titleStyle := m.styles.MenuTitle
		if i == m.selectedMenu {
			style = m.styles.MenuSelected
			titleStyle = m.styles.SelectedTitle
		}
		menuText := fmt.Sprintf("%s\n%s", titleStyle.Render(item.title), m.styles.MenuDesc.Render(item.description))
		menuItems = append(menuItems, style.Render(menuText))
	}

	// Organizar itens em linhas
	numRows := (len(menuItems) + 1) / 2
	var menuRows []string
	for i := 0; i < numRows; i++ {
		startIndex := i * 2
		endIndex := startIndex + 2
		if endIndex > len(menuItems) {
			endIndex = len(menuItems)
		}

		// Se temos dois itens na linha
		if endIndex-startIndex == 2 {
			// Unir horizontalmente com espaçamento
			menuRows = append(menuRows, lipgloss.JoinHorizontal(lipgloss.Top, menuItems[startIndex], "  ", menuItems[startIndex+1]))
		} else {
			// Apenas um item na linha
			menuRows = append(menuRows, menuItems[startIndex])
		}
	}

	// Unir todas as linhas verticalmente
	return lipgloss.JoinVertical(lipgloss.Left, menuRows...)
}

func (m *CLIModel) getContentView() string {
	switch m.currentView {
	case constants.DefaultView:
		return localization.Labels["welcome_message"]
	case constants.CreateWalletView:
		return m.viewCreateWalletPassword()
	case constants.ImportMethodSelectionView:
		return m.viewImportMethodSelection()
	case constants.ImportWalletView:
		return m.viewImportWallet()
	case constants.ImportPrivateKeyView:
		return m.viewImportPrivateKey()
	case constants.ImportWalletPasswordView:
		return m.viewImportWalletPassword()
	case constants.ListWalletsView:
		return m.viewListWallets()
	case constants.WalletPasswordView:
		return m.viewWalletPassword()
	case constants.SelectWalletOperationView:
		return m.viewSelectWalletOperation()
	case constants.DeleteWalletConfirmationView:
		return m.viewDeleteWallet()
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
		case "esc", "backspace":
			// Voltar para o menu principal
			m.menuItems = NewMenu() // Recarregar o menu principal
			m.selectedMenu = 0      // Resetar a seleção
			m.currentView = constants.DefaultView
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

			// Atualizar a contagem de wallets
			return m, m.refreshWalletsTable()
		case "esc", "backspace":
			// Reset the password input field
			m.passwordInput = textinput.New()
			m.currentView = constants.DefaultView
		default:
			var cmd tea.Cmd
			m.passwordInput, cmd = m.passwordInput.Update(msg)
			return m, cmd
		}
	}
	return m, nil
}

func (m *CLIModel) updateDeleteWalletConfirmation(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			deleteConfirmation := strings.TrimSpace(m.deleteConfirmationInput.Value())
			if strings.ToUpper(deleteConfirmation) == "DELETE" {
				err := m.Service.DeleteWallet(m.selectedWallet.ID)

				if err != nil {
					m.err = errors.Wrap(err, 0)
				}

				m.currentView = constants.DefaultView
			}

		case "esc":
			m.currentView = constants.DefaultView

		default:
			var cmd tea.Cmd
			m.deleteConfirmationInput, cmd = m.deleteConfirmationInput.Update(msg)
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
		case "esc", "backspace":
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

			var walletDetails *usecases.WalletDetails
			var err error

			// Check if we're coming from private key import or mnemonic import
			if m.currentView == constants.ImportWalletPasswordView && len(m.privateKeyInput.Value()) > 0 {
				// Import from private key
				privateKey := strings.TrimSpace(m.privateKeyInput.Value())
				walletDetails, err = m.Service.ImportWalletFromPrivateKey(privateKey, password)
			} else {
				// Import from mnemonic
				mnemonic := strings.Join(m.importWords, " ")
				walletDetails, err = m.Service.ImportWallet(mnemonic, password)
			}

			if err != nil {
				m.err = errors.Wrap(err, 0)
				log.Println(m.err.(*errors.Error).ErrorStack())
				m.currentView = constants.DefaultView
				return m, nil
			}

			m.walletDetails = walletDetails
			m.currentView = constants.WalletDetailsView

			// Atualizar a contagem de wallets
			return m, m.refreshWalletsTable()
		case "esc", "backspace":
			m.currentView = constants.DefaultView
		default:
			var cmd tea.Cmd
			m.passwordInput, cmd = m.passwordInput.Update(msg)
			return m, cmd
		}
	}
	return m, nil
}

func (m *CLIModel) updateSelectWalletOperation(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			selectedRow := m.operationTable.Cursor()
			if selectedRow == 0 {
				m.currentView = constants.WalletDetailsView
			} else if selectedRow == 1 {
				m.initDeleteWalletConfirmation()
			}
		case "esc":
			m.currentView = constants.DefaultView
			return m, nil
		}

	}

	var cmd tea.Cmd
	m.operationTable, cmd = m.operationTable.Update(msg)
	return m, cmd
}

func (m *CLIModel) updateListWallets(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			if m.selectedMenu > 0 {
				m.selectedMenu--
			}
		case "down", "j":
			if m.selectedMenu < len(importMenu)-1 {
				m.selectedMenu++
			}
		case "enter":
			// Usar o menu de importação para determinar a ação baseada na seleção
			switch m.selectedMenu {
			case 0: // Primeira opção: Importar por frase mnemônica
				// Preparar campos de entrada para as 12 palavras
				m.textInputs = make([]textinput.Model, constants.MnemonicWordCount)
				m.importWords = make([]string, constants.MnemonicWordCount)
				for i := 0; i < constants.MnemonicWordCount; i++ {
					ti := textinput.New()
					ti.Placeholder = fmt.Sprintf("%s %d", localization.Labels["word"], i+1)
					ti.CharLimit = 50
					ti.Width = 30
					if i == 0 {
						ti.Focus()
					}
					m.textInputs[i] = ti
				}
				m.importStage = 0
				m.currentView = constants.ImportWalletView

			case 1: // Segunda opção: Importar por chave privada
				m.privateKeyInput = textinput.New()
				m.privateKeyInput.Placeholder = localization.Labels["enter_private_key"]
				m.privateKeyInput.CharLimit = 66 // 0x + 64 caracteres hexadecimais
				m.privateKeyInput.Width = 66
				m.privateKeyInput.Focus()
				m.currentView = constants.ImportPrivateKeyView

			case 2: // Terceira opção: Voltar ao menu principal
				m.currentView = constants.DefaultView
				m.selectedMenu = 0
			}
		case "esc", "backspace":
			m.currentView = constants.DefaultView
			m.selectedMenu = 0
		}
	}
	return m, nil
}

func (m *CLIModel) updateImportPrivateKey(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			privateKey := strings.TrimSpace(m.privateKeyInput.Value())
			if privateKey == "" {
				m.err = errors.Wrap(fmt.Errorf(localization.Labels["invalid_private_key"]), 0)
				log.Println(m.err.(*errors.Error).ErrorStack())
				return m, nil
			}

			// Move to password input screen
			m.passwordInput = textinput.New()
			m.passwordInput.Placeholder = localization.Labels["enter_password"]
			m.passwordInput.CharLimit = constants.PasswordCharLimit
			m.passwordInput.Width = constants.PasswordWidth
			m.passwordInput.EchoMode = textinput.EchoPassword
			m.passwordInput.EchoCharacter = '•'
			m.passwordInput.Focus()
			m.currentView = constants.ImportWalletPasswordView

		case "esc", "backspace":
			m.currentView = constants.DefaultView
		default:
			var cmd tea.Cmd
			m.privateKeyInput, cmd = m.privateKeyInput.Update(msg)
			return m, cmd
		}
	}
	return m, nil
}

func (m *CLIModel) updateListWallets(msg tea.Msg) (tea.Model, tea.Cmd) {
	// Diálogo de confirmação de exclusão
	if m.deletingWallet != nil {
		if keyMsg, ok := msg.(tea.KeyMsg); ok {
			switch keyMsg.String() {
			case "left", "h":
				if m.dialogButtonIndex > 0 {
					m.dialogButtonIndex = 0
				}
				return m, nil
			case "right", "l":
				if m.dialogButtonIndex < 1 {
					m.dialogButtonIndex = 1
				}
				return m, nil
			case "enter":
				walletToDelete := m.deletingWallet
				shouldDelete := m.dialogButtonIndex == 0

				// Limpar a referência do diálogo antes de qualquer outra operação
				m.deletingWallet = nil
				m.dialogButtonIndex = 0

				if shouldDelete {
					// Executar a exclusão
					err := m.Service.DeleteWallet(walletToDelete)
					if err != nil {
						m.err = errors.Wrap(err, 0)
					}

					// Recarregar a lista de wallets
					wallets, err := m.Service.GetAllWallets()
					if err == nil {
						m.wallets = wallets
						m.walletCount = len(wallets)

						// Reconstruir linhas da tabela
						rows := make([]table.Row, len(wallets))
						for i, w := range wallets {
							rows[i] = table.Row{fmt.Sprintf("%d", w.ID), w.Address}
						}
						m.walletTable.SetRows(rows)
					}
				}

				// Forçar uma atualização da tela
				return m, m.refreshWalletsTable()
			case "esc", "backspace":
				// Limpar a referência do diálogo e forçar atualização
				m.deletingWallet = nil
				m.dialogButtonIndex = 0
				// Forçar uma atualização da tela
				return m, m.refreshWalletsTable()
			}
		}
		return m, nil
	}

	// Continuar com o código existente para quando não houver diálogo
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "d", "delete":
			selectedRow := m.walletTable.SelectedRow()
			if len(selectedRow) > 1 {
				address := selectedRow[1]
				for i, w := range m.wallets {
					if w.Address == address {
						m.deletingWallet = &m.wallets[i]
						return m, nil
					}
				}
			}
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
		case "esc", "backspace":
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

			m.initSelectWalletOperation()
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
		case "esc", "backspace":
			m.walletDetails = nil
			m.currentView = constants.ListWalletsView
			return m, nil // Return explícito para consumir o evento de teclado
		}
	}
	return m, nil
}

func (m *CLIModel) updateTableDimensions() {
	if m.currentView != constants.ListWalletsView || len(m.wallets) == 0 {
		return
	}

	// Calcular a altura disponível para a área de conteúdo
	headerHeight := lipgloss.Height(m.styles.Header.Render(""))
	footerHeight := lipgloss.Height(m.styles.Footer.Render(""))

	// Reserva de espaço para o título e instruções dentro da área de conteúdo
	titleAndInstructionsHeight := 6 // Espaço estimado para o título e as instruções

	// Calcular altura final da tabela (altura total - cabeçalho - rodapé - título/instruções - margem)
	contentAreaHeight := m.height - headerHeight - footerHeight - titleAndInstructionsHeight - 2

	// Garantir que a tabela tenha pelo menos uma altura mínima
	if contentAreaHeight < 5 {
		contentAreaHeight = 5
	}

	// Definir largura e altura da tabela
	m.walletTable.SetWidth(m.width - 4)
	m.walletTable.SetHeight(contentAreaHeight)

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

func (m *CLIModel) initImportMethodSelection() {
	// Usar o menu de importação que inclui a opção de voltar ao menu principal
	m.menuItems = NewImportMenu()
	m.selectedMenu = 0
	m.currentView = constants.ImportMethodSelectionView
}

func (m *CLIModel) initImportPrivateKey() {
	// Setup private key input
	m.privateKeyInput = textinput.New()
	m.privateKeyInput.Placeholder = localization.Labels["enter_private_key"]
	m.privateKeyInput.CharLimit = 66 // 0x + 64 hex characters
	m.privateKeyInput.Width = 66
	m.privateKeyInput.Focus()
	m.currentView = constants.ImportPrivateKeyView
}

func (m *CLIModel) initImportWallet() {
	// Instead of directly initializing the mnemonic import view,
	// now we show the selection screen first
	m.initImportMethodSelection()
}

func (m *CLIModel) initDeleteWalletConfirmation() {
	m.deleteConfirmationInput = textinput.New()
	m.deleteConfirmationInput.Focus()
	m.deleteConfirmationInput.CharLimit = len("DELETE")

	m.currentView = constants.DeleteWalletConfirmationView
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

	// Definir altura da tabela para usar totalmente o espaço disponível
	contentAreaHeight := m.height - lipgloss.Height(m.styles.Header.Render("")) - lipgloss.Height(m.styles.Footer.Render("")) - 2
	if contentAreaHeight < 0 {
		contentAreaHeight = 0
	}
	m.walletTable.SetHeight(contentAreaHeight)

	// Atualizar dimensões da tabela
	m.updateTableDimensions()

	m.currentView = constants.ListWalletsView
}

func (m *CLIModel) initSelectWalletOperation() {
	columns := []table.Column{
		{Title: localization.Labels["operation"], Width: 30},
	}
	rows := []table.Row{{localization.Labels["wallet_details_option"]}, {localization.Labels["wallet_delete_option"]}}

	m.operationTable = table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
	)

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

	m.operationTable.SetStyles(s)

	m.currentView = constants.SelectWalletOperationView
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

// walletsRefreshedMsg é uma mensagem personalizada para indicar que a lista de wallets foi atualizada
type walletsRefreshedMsg struct{}

func (m *CLIModel) refreshWalletsTable() tea.Cmd {
	return func() tea.Msg {
		// Recarregar a lista de wallets do serviço
		wallets, err := m.Service.GetAllWallets()
		if err != nil {
			m.err = errors.Wrap(err, 0)
			return nil
		}

		// Atualizar a lista de wallets no modelo
		m.wallets = wallets

		// Atualizar a contagem de wallets
		m.walletCount = len(wallets)

		// Reconstruir as linhas da tabela
		var rows []table.Row
		for _, w := range m.wallets {
			rows = append(rows, table.Row{fmt.Sprintf("%d", w.ID), w.Address})
		}

		// Atualizar a tabela com as novas linhas
		m.walletTable.SetRows(rows)

		// Retornar uma mensagem personalizada para indicar que a lista foi atualizada
		return walletsRefreshedMsg{}
	}
}

func (m *CLIModel) rebuildWalletsTable() {
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

	// Definir altura da tabela para usar totalmente o espaço disponível
	contentAreaHeight := m.height - lipgloss.Height(m.styles.Header.Render("")) - lipgloss.Height(m.styles.Footer.Render("")) - 2
	if contentAreaHeight < 0 {
		contentAreaHeight = 0
	}
	m.walletTable.SetHeight(contentAreaHeight)

	// Atualizar dimensões da tabela
	m.updateTableDimensions()
}
