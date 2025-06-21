package ui

import (
	"blocowallet/internal/constants"
	"blocowallet/internal/wallet"
	"blocowallet/pkg/config"
	"blocowallet/pkg/localization"
	"bytes"
	"fmt"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/arsham/figurine/figurine"
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/digitallyserviced/tdfgo/tdf"
	"github.com/ethereum/go-ethereum/accounts/keystore"
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

func NewCLIModel(service *wallet.WalletService) *CLIModel {
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

// loadFontsList returns the list of available fonts from the configuration
func loadFontsList(appDir string) ([]string, error) {
	// The fonts are now loaded from the main configuration
	// This function is kept for compatibility, but it's now a simple wrapper
	// that returns the fonts from the global configuration

	// Get the fonts from the global configuration
	cfg, err := config.LoadConfig(appDir)
	if err != nil {
		return nil, fmt.Errorf("erro ao carregar a configuração: %v", err)
	}

	return cfg.GetFontsList(), nil
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
	case constants.CreateWalletNameView:
		return m.updateCreateWalletName(msg)
	case constants.CreateWalletView:
		return m.updateCreateWalletPassword(msg)
	case constants.ImportMethodSelectionView:
		return m.updateImportMethodSelection(msg)
	case constants.ImportWalletView:
		return m.updateImportWallet(msg)
	case constants.ImportPrivateKeyView:
		return m.updateImportPrivateKey(msg)
	case constants.ImportKeystoreView:
		return m.updateImportKeystore(msg)
	case constants.ImportWalletPasswordView:
		return m.updateImportWalletPassword(msg)
	case constants.ListWalletsView:
		return m.updateListWallets(msg)
	case constants.WalletPasswordView:
		return m.updateWalletPassword(msg)
	case constants.WalletDetailsView:
		return m.updateWalletDetails(msg)
	case constants.ConfigurationView:
		return m.updateConfigMenu(msg)
	case constants.LanguageSelectionView:
		return m.updateLanguageSelection(msg)
	case constants.NetworkMenuView:
		return m.updateNetworkMenu(msg)
	case constants.NetworkListView:
		return m.updateNetworkList(msg)
	case constants.AddNetworkView:
		return m.updateAddNetwork(msg)
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

		if tableHeight > 0 && len(m.wallets) > 0 {
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

// renderConfigMenuItems renderiza os itens do menu de configuração
func (m *CLIModel) renderConfigMenuItems() string {
	// Criar o menu de configuração
	configMenu := NewConfigMenu()

	// Renderizar cada item do menu
	var menuItems []string
	for i, item := range configMenu {
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

// renderLanguageMenuItems renderiza os itens do menu de idiomas
func (m *CLIModel) renderLanguageMenuItems() string {
	// Renderizar cada item do menu
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
	case constants.CreateWalletNameView:
		return m.viewCreateWalletName()
	case constants.CreateWalletView:
		return m.viewCreateWalletPassword()
	case constants.ImportMethodSelectionView:
		return m.viewImportMethodSelection()
	case constants.ImportWalletView:
		return m.viewImportWallet()
	case constants.ImportPrivateKeyView:
		return m.viewImportPrivateKey()
	case constants.ImportKeystoreView:
		return m.viewImportKeystore()
	case constants.ImportWalletPasswordView:
		return m.viewImportWalletPassword()
	case constants.ListWalletsView:
		return m.viewListWallets()
	case constants.WalletPasswordView:
		return m.viewWalletPassword()
	case constants.WalletDetailsView:
		return m.viewWalletDetails()
	case constants.ConfigurationView:
		return m.viewConfigMenu()
	case constants.LanguageSelectionView:
		return m.viewLanguageSelection()
	case constants.NetworkMenuView:
		return m.viewNetworkMenu()
	case constants.NetworkListView:
		return m.viewNetworkList()
	case constants.AddNetworkView:
		return m.viewAddNetwork()
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
			case localization.Labels["configuration"]:
				m.initConfigMenu()
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

func (m *CLIModel) updateCreateWalletName(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			name := strings.TrimSpace(m.nameInput.Value())
			if name == "" {
				m.err = errors.Wrap(fmt.Errorf("O nome da wallet não pode estar vazio"), 0)
				if wrappedErr, ok := m.err.(*errors.Error); ok {
					log.Println(wrappedErr.ErrorStack())
				} else {
					log.Println("Error:", m.err)
				}
				m.currentView = constants.DefaultView
				return m, nil
			}
			// Proceed to password input
			m.passwordInput.Focus()
			m.currentView = constants.CreateWalletView
			return m, nil
		case "esc", "backspace":
			// Reset the name input field and go back to menu
			m.nameInput = textinput.New()
			m.currentView = constants.DefaultView
		default:
			var cmd tea.Cmd
			m.nameInput, cmd = m.nameInput.Update(msg)
			return m, cmd
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

			// Validar a complexidade da senha
			validationErr, isValid := wallet.ValidatePassword(password)
			if !isValid {
				m.err = errors.Wrap(fmt.Errorf(validationErr.GetErrorMessage()), 0)
				log.Println(m.err.(*errors.Error).ErrorStack())
				return m, nil
			}

			name := strings.TrimSpace(m.nameInput.Value())
			walletDetails, err := m.Service.CreateWallet(name, password)
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
			// Go back to name input
			m.nameInput.Focus()
			m.currentView = constants.CreateWalletNameView
			return m, nil
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
				m.passwordInput.Validate = func(s string) error {
					_, isValid := wallet.ValidatePassword(s)
					if !isValid && s != "" {
						return fmt.Errorf("")
					}
					return nil
				}
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

			// Validar a complexidade da senha
			validationErr, isValid := wallet.ValidatePassword(password)
			if !isValid {
				m.err = errors.Wrap(fmt.Errorf(validationErr.GetErrorMessage()), 0)
				log.Println(m.err.(*errors.Error).ErrorStack())
				return m, nil
			}

			var walletDetails *wallet.WalletDetails
			var err error

			// Use a default name based on the import method
			var name string
			if m.currentView == constants.ImportWalletPasswordView && len(m.privateKeyInput.Value()) > 0 {
				name = "Imported Private Key Wallet"
			} else if m.mnemonic != "" && strings.HasSuffix(m.mnemonic, ".json") {
				// If mnemonic field contains a path to a keystore file
				name = "Imported Keystore Wallet"
			} else {
				name = "Imported Mnemonic Wallet"
			}


			// Check which import method we're using
			if m.currentView == constants.ImportWalletPasswordView && (len(m.privateKeyInput.Value()) == 64 || len(m.privateKeyInput.Value()) == 66 && strings.HasPrefix(m.privateKeyInput.Value(), "0x")) {
				// Import from private key
				privateKey := strings.TrimSpace(m.privateKeyInput.Value())
				walletDetails, err = m.Service.ImportWalletFromPrivateKey(name, privateKey, password)
			} else if strings.Contains(m.mnemonic, string(filepath.Separator)) || filepath.IsAbs(m.mnemonic) {
				// Import from keystore file
				keystorePath := m.mnemonic // We stored the keystore path in the mnemonic field
				walletDetails, err = m.Service.ImportWalletFromKeystore(name, keystorePath, password)
			} else {
				// Import from mnemonic
				mnemonic := strings.Join(m.importWords, " ")
				walletDetails, err = m.Service.ImportWallet(name, mnemonic, password)
			}

			if err != nil {
				m.err = errors.Wrap(err, 0)
				log.Println(m.err.(*errors.Error).ErrorStack())

				// If it's a password error for keystore, stay on the password screen
				if errors.Is(err, keystore.ErrDecrypt) {
					return m, nil
				}

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

func (m *CLIModel) updateImportMethodSelection(msg tea.Msg) (tea.Model, tea.Cmd) {
	// Criar o menu de importação
	importMenu := NewImportMenu()

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

			case 2: // Terceira opção: Importar por arquivo keystore
				m.initImportKeystore()

			case 3: // Quarta opção: Voltar ao menu principal
				m.menuItems = NewMenu() // Recarregar o menu principal
				m.selectedMenu = 0      // Resetar a seleção
				m.currentView = constants.DefaultView
			}
		case "esc", "backspace":
			m.menuItems = NewMenu() // Recarregar o menu principal
			m.selectedMenu = 0      // Resetar a seleção
			m.currentView = constants.DefaultView
		}
	}
	return m, nil
}

func (m *CLIModel) updateConfigMenu(msg tea.Msg) (tea.Model, tea.Cmd) {
	// Criar o menu de configuração
	configMenu := NewConfigMenu()

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			if m.selectedMenu > 0 {
				m.selectedMenu--
			}
		case "down", "j":
			if m.selectedMenu < len(configMenu)-1 {
				m.selectedMenu++
			}
		case "enter":
			// Usar o menu de configuração para determinar a ação baseada na seleção
			switch m.selectedMenu {
			case 0: // Primeira opção: Redes
				// Mostrar o submenu de redes
				m.menuItems = NewNetworkMenu()
				m.selectedMenu = 0
				m.currentView = constants.NetworkMenuView
				return m, nil

			case 1: // Segunda opção: Idioma
				// Implementar a lógica para configurar idioma
				m.initLanguageSelection()
				return m, nil

			case 2: // Terceira opção: Voltar ao menu principal
				m.menuItems = NewMenu() // Recarregar o menu principal
				m.selectedMenu = 0      // Resetar a seleção
				m.currentView = constants.DefaultView
			}
		case "esc", "backspace":
			m.menuItems = NewMenu() // Recarregar o menu principal
			m.selectedMenu = 0      // Resetar a seleção
			m.currentView = constants.DefaultView
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
			m.passwordInput.Validate = func(s string) error {
				_, isValid := wallet.ValidatePassword(s)
				if !isValid && s != "" {
					return fmt.Errorf("")
				}
				return nil
			}
			m.passwordInput.Focus()
			m.currentView = constants.ImportWalletPasswordView

		case "esc", "backspace":
			m.currentView = constants.DefaultView
		default:
			var cmd tea.Cmd
			m.privateKeyInput, cmd = m.privateKeyInput.Update(msg)

			// Update suggestions as the user types
			if msg.Type == tea.KeyRunes || msg.Type == tea.KeyBackspace || msg.Type == tea.KeyDelete {
				// Get current path
				currentPath := m.privateKeyInput.Value()
				if currentPath == "" {
					currentPath = "."
				}

				// Get the directory and partial filename
				dir := filepath.Dir(currentPath)
				if dir == "." && !strings.HasPrefix(currentPath, "./") && !strings.HasPrefix(currentPath, "/") {
					dir = currentPath
				}

				// Read the directory
				files, err := os.ReadDir(dir)
				if err == nil {
					// Find matching files
					var matches []string
					partial := filepath.Base(currentPath)
					for _, file := range files {
						if strings.HasPrefix(file.Name(), partial) {
							fullPath := filepath.Join(dir, file.Name())
							if file.IsDir() {
								fullPath += "/"
							}
							matches = append(matches, fullPath)
						}
					}

					// Set all matches as suggestions
					if len(matches) > 0 {
						m.privateKeyInput.SetSuggestions(matches)
					}
				}
			}

			return m, cmd
		}
	}
	return m, nil
}

func (m *CLIModel) updateImportKeystore(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			keystorePath := strings.TrimSpace(m.privateKeyInput.Value())
			if keystorePath == "" {
				m.err = errors.Wrap(fmt.Errorf(localization.Labels["invalid_keystore"]), 0)
				log.Println(m.err.(*errors.Error).ErrorStack())
				return m, nil
			}

			// Check if file exists
			if _, err := os.Stat(keystorePath); os.IsNotExist(err) {
				m.err = errors.Wrap(fmt.Errorf(localization.Labels["invalid_keystore"]), 0)
				log.Println(m.err.(*errors.Error).ErrorStack())
				return m, nil
			}

			// Store the keystore path for later use
			m.mnemonic = keystorePath // Reusing mnemonic field to store keystore path

			// Move to password input screen
			m.passwordInput = textinput.New()
			m.passwordInput.Placeholder = localization.Labels["enter_password"]
			m.passwordInput.CharLimit = constants.PasswordCharLimit
			m.passwordInput.Width = constants.PasswordWidth
			m.passwordInput.EchoMode = textinput.EchoPassword
			m.passwordInput.EchoCharacter = '•'
			m.passwordInput.Validate = func(s string) error {
				_, isValid := wallet.ValidatePassword(s)
				if !isValid && s != "" {
					return fmt.Errorf("")
				}
				return nil
			}
			m.passwordInput.Focus()
			m.currentView = constants.ImportWalletPasswordView

		case "esc", "backspace":
			m.currentView = constants.ImportMethodSelectionView
		case "tab":
			// Implement path autocomplete
			currentPath := m.privateKeyInput.Value()
			if currentPath == "" {
				currentPath = "."
			}

			// Get the directory and partial filename
			dir := filepath.Dir(currentPath)
			if dir == "." && !strings.HasPrefix(currentPath, "./") && !strings.HasPrefix(currentPath, "/") {
				dir = currentPath
			}

			// Read the directory
			files, err := os.ReadDir(dir)
			if err != nil {
				return m, nil
			}

			// Find matching files
			var matches []string
			partial := filepath.Base(currentPath)
			for _, file := range files {
				if strings.HasPrefix(file.Name(), partial) {
					fullPath := filepath.Join(dir, file.Name())
					if file.IsDir() {
						fullPath += "/"
					}
					matches = append(matches, fullPath)
				}
			}

			// Set all matches as suggestions
			if len(matches) > 0 {
				m.privateKeyInput.SetSuggestions(matches)

				// If there's exactly one match, use it
				if len(matches) == 1 {
					m.privateKeyInput.SetValue(matches[0])
				}
			}

			// Let the textinput component handle the tab key
			var cmd tea.Cmd
			m.privateKeyInput, cmd = m.privateKeyInput.Update(msg)
			return m, cmd
		default:
			var cmd tea.Cmd
			m.privateKeyInput, cmd = m.privateKeyInput.Update(msg)

			// Update suggestions as the user types
			if msg.Type == tea.KeyRunes || msg.Type == tea.KeyBackspace || msg.Type == tea.KeyDelete {
				// Get current path
				currentPath := m.privateKeyInput.Value()
				if currentPath == "" {
					currentPath = "."
				}

				// Get the directory and partial filename
				dir := filepath.Dir(currentPath)
				if dir == "." && !strings.HasPrefix(currentPath, "./") && !strings.HasPrefix(currentPath, "/") {
					dir = currentPath
				}

				// Read the directory
				files, err := os.ReadDir(dir)
				if err == nil {
					// Find matching files
					var matches []string
					partial := filepath.Base(currentPath)
					for _, file := range files {
						if strings.HasPrefix(file.Name(), partial) {
							fullPath := filepath.Join(dir, file.Name())
							if file.IsDir() {
								fullPath += "/"
							}
							matches = append(matches, fullPath)
						}
					}

					// Set all matches as suggestions
					if len(matches) > 0 {
						m.privateKeyInput.SetSuggestions(matches)
					}
				}
			}

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
							rows[i] = table.Row{fmt.Sprintf("%d", w.ID), w.Name, w.Address}
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
			// Only try to access the table if there are wallets
			if len(m.wallets) > 0 {
				selectedRow := m.walletTable.SelectedRow()
				if len(selectedRow) > 4 {
					address := selectedRow[4]
					for i, w := range m.wallets {
						if w.Address == address {
							m.deletingWallet = &m.wallets[i]
							return m, nil
						}
					}
				}
			}
		case "enter":
			// Only try to access the table if there are wallets
			if len(m.wallets) > 0 {
				selectedRow := m.walletTable.SelectedRow()
				if len(selectedRow) > 4 {
					address := selectedRow[4]
					// Buscar wallet pela address
					for _, w := range m.wallets {
						if w.Address == address {
							m.selectedWallet = &w
							m.initWalletPassword()
							return m, nil
						}
					}
				}
			}
		case "esc", "backspace":
			m.currentView = constants.DefaultView
			return m, nil
		}
	}

	// Atualizar a tabela com a mensagem apenas se houver wallets
	if len(m.wallets) > 0 {
		var cmd tea.Cmd
		m.walletTable, cmd = m.walletTable.Update(msg)
		return m, cmd
	}
	return m, nil
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

func (m *CLIModel) updateWalletDetails(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc", "backspace":
			m.walletDetails = nil
			m.currentView = constants.ListWalletsView

			// Ensure the wallet list is properly initialized before showing it
			wallets, err := m.Service.GetAllWallets()
			if err == nil {
				m.wallets = wallets
				m.walletCount = len(wallets)

				// Always rebuild the table, even if there are no wallets
				// The rebuildWalletsTable method already has a check for empty wallets
				m.rebuildWalletsTable()
			}

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
	// Reduzir a largura da tabela para evitar quebra de linha
	m.walletTable.SetWidth(m.width - 12)
	if len(m.wallets) > 0 {
		m.walletTable.SetHeight(contentAreaHeight)
	}

	// Calcular larguras das colunas
	idColWidth := 10
	nameColWidth := 20
	typeColWidth := 20
	createdAtColWidth := 20
	// Aumentar a margem para evitar quebra de linha
	addressColWidth := m.width - idColWidth - nameColWidth - typeColWidth - createdAtColWidth - 20

	if addressColWidth < 20 {
		addressColWidth = 20
	}

	// Atualizar colunas - manter consistente com initListWallets e rebuildWalletsTable
	m.walletTable.SetColumns([]table.Column{
		{Title: localization.Labels["id"], Width: idColWidth},
		{Title: "Nome", Width: nameColWidth},
		{Title: localization.Labels["wallet_type"], Width: typeColWidth},
		{Title: localization.Labels["created_at"], Width: createdAtColWidth},
		{Title: localization.Labels["ethereum_address"], Width: addressColWidth},
	})
}

// Funções de inicialização

func (m *CLIModel) initCreateWallet() {
	m.mnemonic, _ = wallet.GenerateMnemonic()

	// Initialize name input first
	m.nameInput = textinput.New()
	m.nameInput.Placeholder = "Digite o nome da wallet"
	m.nameInput.CharLimit = 50
	m.nameInput.Width = constants.PasswordWidth
	m.nameInput.Focus()
	m.currentView = constants.CreateWalletNameView

	// Initialize password input (will be used after name is entered)
	m.passwordInput = textinput.New()
	m.passwordInput.Placeholder = localization.Labels["enter_password"]
	m.passwordInput.CharLimit = constants.PasswordCharLimit
	m.passwordInput.Width = constants.PasswordWidth
	m.passwordInput.EchoMode = textinput.EchoPassword
	m.passwordInput.EchoCharacter = '•'
	m.passwordInput.Validate = func(s string) error {
		_, isValid := wallet.ValidatePassword(s)
		if !isValid && s != "" {
			return fmt.Errorf("")
		}
		return nil
	}
}

func (m *CLIModel) initImportMethodSelection() {
	// Usar o menu de importação que inclui a opção de voltar ao menu principal
	m.menuItems = NewImportMenu()
	m.selectedMenu = 0
	m.currentView = constants.ImportMethodSelectionView
}

func (m *CLIModel) initConfigMenu() {
	// Usar o menu de configuração que inclui a opção de voltar ao menu principal
	m.menuItems = NewConfigMenu()
	m.selectedMenu = 0
	m.currentView = constants.ConfigurationView
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

// initImportKeystore initializes the keystore import view
func (m *CLIModel) initImportKeystore() {
	// Setup keystore path input with autocomplete
	m.privateKeyInput = textinput.New()
	m.privateKeyInput.Placeholder = localization.Labels["enter_keystore_path"]
	m.privateKeyInput.CharLimit = 256 // Path can be long
	m.privateKeyInput.Width = 66
	m.privateKeyInput.Focus()
	m.privateKeyInput.ShowSuggestions = true // Enable suggestions
	m.currentView = constants.ImportKeystoreView
}

func (m *CLIModel) initImportWallet() {
	// Instead of directly initializing the mnemonic import view,
	// now we show the selection screen first
	m.initImportMethodSelection()
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
	nameColWidth := 20
	typeColWidth := 20
	createdAtColWidth := 20
	addressColWidth := m.width - idColWidth - nameColWidth - typeColWidth - createdAtColWidth - 20 // Subtrai 20 para padding e margens

	if addressColWidth < 20 {
		addressColWidth = 20
	}

	columns := []table.Column{
		{Title: localization.Labels["id"], Width: idColWidth},
		{Title: "Nome", Width: nameColWidth},
		{Title: localization.Labels["wallet_type"], Width: typeColWidth},
		{Title: localization.Labels["created_at"], Width: createdAtColWidth},
		{Title: localization.Labels["ethereum_address"], Width: addressColWidth},
	}

	var rows []table.Row
	for _, w := range m.wallets {
		// Determine wallet type based on mnemonic presence
		walletType := localization.Labels["imported_mnemonic"]
		if w.Mnemonic == "" {
			walletType = localization.Labels["imported_private_key"]
		}

		// Format created at date
		createdAt := w.CreatedAt.Format("2006-01-02 15:04")

		rows = append(rows, table.Row{
			fmt.Sprintf("%d", w.ID),
			w.Name,
			walletType,
			createdAt,
			w.Address,
		})
	}

	m.walletTable = table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
	)

	// Definir largura explicitamente para evitar quebra de linha
	m.walletTable.SetWidth(m.width - 12)

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
	if len(m.wallets) > 0 {
		m.walletTable.SetHeight(contentAreaHeight)
	}

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
	m.passwordInput.Validate = func(s string) error {
		_, isValid := wallet.ValidatePassword(s)
		if !isValid && s != "" {
			return fmt.Errorf("")
		}
		return nil
	}
	m.passwordInput.Focus()
	m.currentView = constants.WalletPasswordView
}

// initLanguageSelection initializes the language selection view
func (m *CLIModel) initLanguageSelection() {
	// Load the current configuration
	appDir := filepath.Join(os.Getenv("HOME"), ".wallets")
	cfg, err := config.LoadConfig(appDir)
	if err != nil {
		m.err = errors.Wrap(err, 0)
		m.currentView = constants.DefaultView
		return
	}

	// Store the current configuration
	m.currentConfig = cfg

	// Set the menu items to the language menu items
	m.menuItems = NewLanguageMenu(cfg)

	// Reset the selected menu item
	m.selectedMenu = 0

	// Set the current view to language selection
	m.currentView = constants.LanguageSelectionView
}

// updateLanguageSelection handles user input in the language selection view
func (m *CLIModel) updateLanguageSelection(msg tea.Msg) (tea.Model, tea.Cmd) {
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
			// If the last item (Back) is selected, return to the config menu
			if m.selectedMenu == len(m.menuItems)-1 {
				m.menuItems = NewConfigMenu()
				m.selectedMenu = 0
				m.currentView = constants.ConfigurationView
				return m, nil
			}

			// Otherwise, change the language
			selectedLang := m.menuItems[m.selectedMenu].description

			// Update the configuration
			if m.currentConfig != nil && selectedLang != m.currentConfig.Language {
				// Get the config file path
				configPath := filepath.Join(m.currentConfig.AppDir, "config.toml")

				// Read the current config file
				content, err := os.ReadFile(configPath)
				if err != nil {
					m.err = errors.Wrap(err, 0)
					return m, nil
				}

				// Replace the language setting
				lines := strings.Split(string(content), "\n")
				for i, line := range lines {
					if strings.HasPrefix(strings.TrimSpace(line), "language") {
						lines[i] = fmt.Sprintf("language = \"%s\"", selectedLang)
						break
					}
				}

				// Write the updated config back to the file
				err = os.WriteFile(configPath, []byte(strings.Join(lines, "\n")), 0644)
				if err != nil {
					m.err = errors.Wrap(err, 0)
					return m, nil
				}

				// Reload the configuration
				newCfg, err := config.LoadConfig(m.currentConfig.AppDir)
				if err != nil {
					m.err = errors.Wrap(err, 0)
					return m, nil
				}

				// Update the current configuration
				m.currentConfig = newCfg

				// Reinitialize localization with the new language
				err = localization.InitLocalization(newCfg)
				if err != nil {
					m.err = errors.Wrap(err, 0)
					return m, nil
				}

				// Return to the config menu
				m.menuItems = NewConfigMenu()
				m.selectedMenu = 0
				m.currentView = constants.ConfigurationView
			} else {
				// If no change or error, just return to the config menu
				m.menuItems = NewConfigMenu()
				m.selectedMenu = 0
				m.currentView = constants.ConfigurationView
			}
		case "esc", "backspace":
			// Return to the config menu
			m.menuItems = NewConfigMenu()
			m.selectedMenu = 0
			m.currentView = constants.ConfigurationView
		}
	}
	return m, nil
}

// updateNetworkMenu handles user input in the network menu view
func (m *CLIModel) updateNetworkMenu(msg tea.Msg) (tea.Model, tea.Cmd) {
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
		case "enter":
			// If the last item (Back) is selected, return to the config menu
			if m.selectedMenu == len(m.menuItems)-1 {
				m.menuItems = NewConfigMenu()
				m.selectedMenu = 0
				m.currentView = constants.ConfigurationView
				return m, nil
			}

			// Otherwise, handle the selected option
			switch m.selectedMenu {
			case 0: // Add Network
				m.initAddNetwork()
				return m, nil
			case 1: // Network List
				m.initNetworkList()
				return m, nil
			}
		case "esc", "backspace":
			// Return to the config menu
			m.menuItems = NewConfigMenu()
			m.selectedMenu = 0
			m.currentView = constants.ConfigurationView
		}
	}
	return m, nil
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

		// Se houver wallets, reconstruir a tabela completamente
		// para garantir que ela seja inicializada corretamente
		if len(m.wallets) > 0 {
			m.rebuildWalletsTable()
		}

		// Retornar uma mensagem personalizada para indicar que a lista foi atualizada
		return walletsRefreshedMsg{}
	}
}

func (m *CLIModel) rebuildWalletsTable() {
	// Only create a table if there are wallets
	if len(m.wallets) == 0 {
		return
	}

	// Inicialize as colunas com larguras adequadas
	idColWidth := 10
	nameColWidth := 20
	typeColWidth := 20
	createdAtColWidth := 20
	addressColWidth := m.width - idColWidth - nameColWidth - typeColWidth - createdAtColWidth - 20 // Subtrai 20 para padding e margens

	if addressColWidth < 20 {
		addressColWidth = 20
	}

	columns := []table.Column{
		{Title: localization.Labels["id"], Width: idColWidth},
		{Title: "Nome", Width: nameColWidth},
		{Title: localization.Labels["wallet_type"], Width: typeColWidth},
		{Title: localization.Labels["created_at"], Width: createdAtColWidth},
		{Title: localization.Labels["ethereum_address"], Width: addressColWidth},
	}

	var rows []table.Row
	for _, w := range m.wallets {
		// Determine wallet type based on mnemonic presence
		walletType := localization.Labels["imported_mnemonic"]
		if w.Mnemonic == "" {
			walletType = localization.Labels["imported_private_key"]
		}

		// Format created at date
		createdAt := w.CreatedAt.Format("2006-01-02 15:04")

		rows = append(rows, table.Row{
			fmt.Sprintf("%d", w.ID),
			w.Name,
			walletType,
			createdAt,
			w.Address,
		})
	}

	m.walletTable = table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
	)

	// Definir largura explicitamente para evitar quebra de linha
	m.walletTable.SetWidth(m.width - 12)

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
