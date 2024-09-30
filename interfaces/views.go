package interfaces

import (
	"blocowallet/localization"
	"fmt"
	"github.com/charmbracelet/lipgloss"
	"github.com/ethereum/go-ethereum/crypto"
	"strings"
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

// viewImportWallet renderiza a visualização de importação de wallet
func (m *CLIModel) viewImportWallet() string {
	if localization.Labels == nil {
		return "Localization labels not initialized."
	}

	var view strings.Builder
	view.WriteString(
		lipgloss.NewStyle().Bold(true).Render(localization.Labels["import_wallet_title"] + "\n\n"),
	)
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

// viewListWallets renderiza a visualização de listagem de wallets
func (m *CLIModel) viewListWallets() string {
	if localization.Labels == nil {
		return "Localization labels not initialized."
	}

	return m.walletTable.View()
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
