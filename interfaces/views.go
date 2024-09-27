package interfaces

import (
	"fmt"
	"strings"

	"blocowallet/constants"
	"blocowallet/localization"

	"github.com/charmbracelet/lipgloss"
)

func (m *CLIModel) View() string {
	if m.err != nil {
		return fmt.Sprintf(localization.Labels["error_message"], m.err)
	}

	var menuView string
	var contentView string

	menuStyle := lipgloss.NewStyle().Padding(1).Align(lipgloss.Left).Width(constants.MenuWidth)

	switch m.currentView {
	case "menu":
		menuView = m.menuList.View()
		contentView = m.viewWelcome()
	case "create_wallet_password":
		menuView = m.menuList.View()
		contentView = m.viewCreateWalletPassword()
	case "import_wallet":
		menuView = m.menuList.View()
		contentView = m.viewImportWallet()
	case "import_wallet_password":
		menuView = m.menuList.View()
		contentView = m.viewImportWalletPassword()
	case "list_wallets":
		menuView = m.menuList.View()
		contentView = m.viewWalletList()
	case "wallet_password":
		menuView = m.menuList.View()
		contentView = m.viewWalletPassword()
	case "wallet_details":
		menuView = m.menuList.View()
		contentView = m.viewWalletDetails()
	default:
		menuView = m.menuList.View()
		contentView = localization.Labels["unknown_state"]
	}

	return lipgloss.JoinHorizontal(
		lipgloss.Top,
		menuStyle.Render(menuView),
		contentStyle.Width(m.width-constants.MenuWidth-2).Height(m.height).Render(contentView),
	)
}

// View functions

func (m *CLIModel) viewWelcome() string {
	return localization.Labels["welcome_message"]
}

func (m *CLIModel) viewCreateWalletPassword() string {
	var view strings.Builder
	view.WriteString(lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#00FF00")).Render(localization.Labels["mnemonic_phrase"]) + "\n\n")
	view.WriteString(fmt.Sprintf("%s\n\n", m.mnemonic))
	view.WriteString(localization.Labels["enter_password"] + "\n\n")
	passwordStyle := lipgloss.NewStyle().Align(lipgloss.Left)
	view.WriteString(passwordStyle.Render(m.passwordInput.View()))
	view.WriteString("\n\n" + localization.Labels["press_enter"])
	return view.String()
}

func (m *CLIModel) viewImportWallet() string {
	var view strings.Builder
	view.WriteString(lipgloss.NewStyle().Bold(true).Render(localization.Labels["import_wallet_title"] + "\n\n"))
	for i, ti := range m.textInputs {
		if i == m.importStage {
			view.WriteString(fmt.Sprintf("%s %d: %s\n", localization.Labels["word"], i+1, ti.View()))
		} else {
			view.WriteString(fmt.Sprintf("%s %d: %s\n", localization.Labels["word"], i+1, ti.Value()))
		}
	}
	view.WriteString("\n" + localization.Labels["press_enter"])
	return view.String()
}

func (m *CLIModel) viewImportWalletPassword() string {
	var view strings.Builder
	view.WriteString(lipgloss.NewStyle().Bold(true).Render(localization.Labels["enter_password"] + "\n\n"))
	passwordStyle := lipgloss.NewStyle().Align(lipgloss.Left)
	view.WriteString(passwordStyle.Render(m.passwordInput.View()))
	view.WriteString("\n\n" + localization.Labels["press_enter"])
	return view.String()
}

func (m *CLIModel) viewWalletList() string {
	var view strings.Builder
	view.WriteString(m.walletTable.View())
	view.WriteString("\n" + localization.Labels["wallet_list_instructions"])
	return view.String()
}

func (m *CLIModel) viewWalletPassword() string {
	var view strings.Builder
	view.WriteString(lipgloss.NewStyle().Bold(true).Render(localization.Labels["enter_wallet_password"] + "\n\n"))
	view.WriteString(m.passwordInput.View())
	view.WriteString("\n\n" + localization.Labels["press_enter"])
	return view.String()
}

func (m *CLIModel) viewWalletDetails() string {
	if m.selectedWallet == nil {
		return localization.Labels["select_wallet_prompt"]
	}

	// Omit displaying private key and mnemonic for security

	detailStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#00FF00")).Align(lipgloss.Left)
	labelStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#FFFFFF")).Align(lipgloss.Left)
	valueStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#AAAAAA")).Align(lipgloss.Left)

	labelWidth := 20

	var view strings.Builder
	view.WriteString(detailStyle.Render(localization.Labels["wallet_details_title"] + "\n\n"))

	// Ethereum Address
	ethLabel := fmt.Sprintf("%-*s", labelWidth, localization.Labels["ethereum_address"])
	view.WriteString(labelStyle.Render(ethLabel) + valueStyle.Render(fmt.Sprintf(" %s\n", m.selectedWallet.Address)))

	// Additional details can be added here

	view.WriteString("\n" + localization.Labels["press_esc"])
	return view.String()
}
