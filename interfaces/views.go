package interfaces

import (
	"blocowallet/constants"
	"blocowallet/localization"
	"bytes"
	"fmt"
	"github.com/arsham/figurine/figurine"
	"github.com/charmbracelet/lipgloss"
	"github.com/digitallyserviced/tdfgo/tdf"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/go-errors/errors"
	"log"
	"strings"
	"time"
)

// viewCreateWalletPassword renderiza a visualização de criação de wallet
func (m *CLIModel) viewCreateWalletPassword() string {
	if localization.Labels == nil {
		return "Localization labels not initialized."
	}

	var view strings.Builder
	view.WriteString(
		lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#00FF00")).Render(localization.Labels["mnemonic_phrase"]) + "\n\n" +
			fmt.Sprintf("%s\n\n", m.mnemonic) +
			localization.Labels["enter_password"] + "\n\n" +
			m.passwordInput.View() + "\n\n" +
			localization.Labels["press_enter"],
	)
	return view.String()
}

// renderSplash renderiza a tela de splash screen
func (m *CLIModel) renderSplash() string {
	// Verificar se a fonte selecionada está disponível
	if m.selectedFont == nil {
		log.Println("Fonte selecionada não está carregada.")
		return m.styles.ErrorStyle.Render(constants.ErrorFontNotFoundMessage)
	}

	// Inicializar o renderizador de string para a fonte selecionada
	fontString := tdf.NewTheDrawFontStringFont(m.selectedFont)

	// Renderizar o logo "bloco"
	renderedLogo := fontString.RenderString("bloco")
	renderedLogo = strings.TrimSpace(renderedLogo) // Remove any extra whitespace

	projectInfo := fmt.Sprintf("%s v%s", "BLOCO Wallet Manager", localization.Labels["version"])

	// Center the projectInfo text
	projectInfoStyled := lipgloss.NewStyle().
		Align(lipgloss.Center).
		Render(projectInfo)

	// Create the splash screen content
	splashContent := lipgloss.JoinVertical(
		lipgloss.Center,
		renderedLogo,
		projectInfoStyled,
	)

	// Usar lipgloss.Place para centralizar horizontal e verticalmente
	finalSplash := lipgloss.Place(
		m.width,
		m.height,
		lipgloss.Center,
		lipgloss.Center,
		splashContent,
	)

	return finalSplash
}

func (m *CLIModel) renderStatusBar() string {
	// Left part: Number of wallets
	leftStyle := m.styles.StatusBarLeft // Used assignment for copying.
	left := leftStyle.
		SetString(fmt.Sprintf("Wallets: %d", m.walletCount)).
		String()

	// Right part: Current date and time
	currentTime := time.Now().Format("02-01-2006 15:04:05")
	rightStyle := m.styles.StatusBarRight // Used assignment for copying.
	right := rightStyle.
		SetString(fmt.Sprintf("Date: %s", currentTime)).
		String()

	// Center part: Current view and shortcut keys
	centerContent := fmt.Sprintf("View: %s | Press 'esc' or 'backspace' to return | Press 'q' to quit", localization.Labels[m.currentView])

	centerWidth := m.width - lipgloss.Width(left) - lipgloss.Width(right)
	centerStyle := m.styles.StatusBarCenter // Used assignment for copying.
	center := centerStyle.
		SetString(centerContent).
		Width(centerWidth).
		Align(lipgloss.Center).
		String()

	// Join all parts
	statusBar := lipgloss.JoinHorizontal(
		lipgloss.Top,
		left,
		center,
		right,
	)

	return statusBar
}

func (m *CLIModel) renderMainView() string {
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
	//statusBar := fmt.Sprintf("Current view: %s | Wallets: %d", localization.Labels[m.currentView], walletCount)
	//renderedFooter := m.styles.Footer.Render(statusBar)
	//footerHeight := lipgloss.Height(renderedFooter)
	renderedFooter := m.renderStatusBar()
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

// viewImportWallet renderiza a visualização de importação de wallet
func (m *CLIModel) viewImportWallet() string {
	if localization.Labels == nil {
		return "Localization labels not initialized."
	}

	var view strings.Builder

	// Renderizando o título com destaque
	title := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#7D56F4")).
		MarginBottom(1).
		Render(localization.Labels["import_wallet_title"])

	view.WriteString(title + "\n\n")

	// Estilo para o campo ativo
	activeStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#00FF00")).
		Bold(true)

	// Estilo para campos inativos
	inactiveStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#AAAAAA"))

	// Renderizar cada campo de entrada
	for i, ti := range m.textInputs {
		wordLabel := fmt.Sprintf("%s %d:", localization.Labels["word"], i+1)
		paddedLabel := fmt.Sprintf("%-10s", wordLabel) // Padding para alinhamento

		if i == m.importStage {
			// Campo ativo com destaque
			view.WriteString(activeStyle.Render(paddedLabel) + " " + ti.View() + "\n\n")
		} else {
			// Campos inativos
			view.WriteString(inactiveStyle.Render(paddedLabel) + " " + ti.Value() + "\n")
		}
	}

	// Instruções para o usuário
	instructions := lipgloss.NewStyle().
		MarginTop(1).
		Italic(true).
		Render(localization.Labels["press_enter"])

	view.WriteString("\n" + instructions)

	// Adicionar uma borda ao redor de tudo
	content := view.String()
	return lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#7D56F4")).
		Padding(1, 2).
		Render(content)
}

// viewImportWalletPassword renderiza a visualização de senha após importação
func (m *CLIModel) viewImportWalletPassword() string {
	if localization.Labels == nil {
		return "Localization labels not initialized."
	}

	var view strings.Builder
	view.WriteString(
		lipgloss.NewStyle().Bold(true).Render(localization.Labels["enter_password"]+"\n\n") +
			m.passwordInput.View() + "\n\n" +
			localization.Labels["press_enter"],
	)
	return view.String()
}

// viewImportMethodSelection renderiza a visualização de seleção de methods de importação
func (m *CLIModel) viewImportMethodSelection() string {
	if localization.Labels == nil {
		return "Localization labels not initialized."
	}

	// Em vez de renderizar o menu de importação novamente, exibir apenas uma mensagem informativa
	// já que o menu já é exibido na área padrão de menu
	return localization.Labels["welcome_message"]
}

// viewImportPrivateKey renderiza a visualização de importação de chave privada
func (m *CLIModel) viewImportPrivateKey() string {
	// Use MenuTitle style for the header instead of non-existent Title style
	title := m.styles.MenuTitle.Render(localization.Labels["private_key_title"])
	input := m.privateKeyInput.View()
	// Use MenuDesc instead of non-existent Instructions style
	instructions := m.styles.MenuDesc.Render(localization.Labels["press_enter"])

	return lipgloss.JoinVertical(
		lipgloss.Left,
		title,
		"",
		input,
		"",
		instructions,
	)
}

// viewListWallets renderiza a visualização de listagem de wallets
func (m *CLIModel) viewListWallets() string {
	if localization.Labels == nil {
		return "Localization labels not initialized."
	}

	// Se não há diálogo de exclusão, retornar apenas a tabela
	if m.deletingWallet == nil {
		var view strings.Builder
		view.WriteString(m.walletTable.View())
		return view.String()
	}

	// Se há um diálogo de confirmação de exclusão, renderizar o diálogo
	return m.renderDeleteConfirmationDialog()
}

// renderDeleteConfirmationDialog renderiza o diálogo de confirmação de exclusão
func (m *CLIModel) renderDeleteConfirmationDialog() string {
	// Primeiro, renderizar a tabela de wallets
	tableView := m.walletTable.View()

	// Caixa de diálogo centralizada com botões estilizados e seleção
	question := localization.Labels["confirm_delete_wallet"]
	address := fmt.Sprintf("%s: %s", localization.Labels["ethereum_address"], m.deletingWallet.Address)

	// Botões com seleção (garante espaçamento entre os textos)
	var confirmBtn, cancelBtn string
	if m.dialogButtonIndex == 0 {
		confirmBtn = m.styles.DialogButtonActive.Render("[ " + localization.Labels["confirm"] + " ]")
		cancelBtn = m.styles.DialogButton.Render("[ " + localization.Labels["cancel"] + " ]")
	} else {
		confirmBtn = m.styles.DialogButton.Render("[ " + localization.Labels["confirm"] + " ]")
		cancelBtn = m.styles.DialogButtonActive.Render("[ " + localization.Labels["cancel"] + " ]")
	}

	buttons := lipgloss.JoinHorizontal(lipgloss.Center, confirmBtn, "   ", cancelBtn)
	content := lipgloss.JoinVertical(lipgloss.Center, question, address, "", buttons)
	dialog := m.styles.Dialog.Render(content)

	// Calcular a posição do diálogo para centralizá-lo na área da tabela
	tableWidth := lipgloss.Width(tableView)
	tableHeight := lipgloss.Height(tableView)
	dialogWidth := lipgloss.Width(dialog)
	dialogHeight := lipgloss.Height(dialog)

	// Calcular posições para centralizar o diálogo na tabela
	leftPadding := (tableWidth - dialogWidth) / 2
	if leftPadding < 0 {
		leftPadding = 0
	}

	// Dividir a tabela em linhas
	tableLines := strings.Split(tableView, "\n")

	// Calcular a linha inicial para o diálogo
	startLine := (tableHeight - dialogHeight) / 2
	if startLine < 0 {
		startLine = 0
	}

	// Dividir o diálogo em linhas
	dialogLines := strings.Split(dialog, "\n")

	// Inserir o diálogo nas linhas da tabela
	for i := 0; i < dialogHeight && i+startLine < len(tableLines); i++ {
		// Garantir que a linha da tabela é longa o suficiente
		for len(tableLines[i+startLine]) < leftPadding {
			tableLines[i+startLine] += " "
		}

		// Inserir a linha do diálogo na posição correta
		if leftPadding < len(tableLines[i+startLine]) {
			prefix := tableLines[i+startLine][:leftPadding]
			suffix := ""
			if leftPadding+dialogWidth < len(tableLines[i+startLine]) {
				suffix = tableLines[i+startLine][leftPadding+dialogWidth:]
			}
			tableLines[i+startLine] = prefix + dialogLines[i] + suffix
		} else {
			padding := strings.Repeat(" ", leftPadding-len(tableLines[i+startLine]))
			tableLines[i+startLine] += padding + dialogLines[i]
		}
	}

	// Reconstruir a visualização da tabela com o diálogo
	return strings.Join(tableLines, "\n")
}

// viewWalletPassword renderiza a visualização de entrada de senha para wallet selecionada
func (m *CLIModel) viewWalletPassword() string {
	if localization.Labels == nil {
		return "Localization labels not initialized."
	}

	var view strings.Builder
	view.WriteString(
		lipgloss.NewStyle().Bold(true).Render(localization.Labels["enter_wallet_password"]+"\n\n") +
			m.passwordInput.View() + "\n\n" +
			localization.Labels["press_enter"],
	)
	return view.String()
}

// viewWalletDetails renderiza a visualização de detalhes da wallet
func (m *CLIModel) viewWalletDetails() string {
	if localization.Labels == nil {
		return "Localization labels not initialized."
	}

	if m.walletDetails != nil {
		var view strings.Builder
		view.WriteString(
			lipgloss.NewStyle().Bold(true).Render(localization.Labels["wallet_details_title"]+"\n\n") +
				fmt.Sprintf("%-*s %s\n", 20, localization.Labels["ethereum_address"], m.walletDetails.Wallet.Address) +
				fmt.Sprintf("%-*s 0x%x\n", 20, localization.Labels["private_key"], crypto.FromECDSA(m.walletDetails.PrivateKey)) +
				fmt.Sprintf("%-*s %x\n", 20, localization.Labels["public_key"], crypto.FromECDSAPub(m.walletDetails.PublicKey)) +
				fmt.Sprintf("%-*s %s\n\n", 20, localization.Labels["mnemonic_phrase_label"], m.walletDetails.Mnemonic) +
				localization.Labels["press_esc"],
		)
		return view.String()
	}
	return localization.Labels["select_wallet_prompt"]
}
