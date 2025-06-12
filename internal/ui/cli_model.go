package ui

import (
	"blocowallet/internal/wallet"
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/digitallyserviced/tdfgo/tdf"
)

type CLIModel struct {
	Service           *wallet.WalletService
	currentView       string
	menuItems         []menuItem
	selectedMenu      int
	importWords       []string
	importStage       int
	textInputs        []textinput.Model
	wallets           []wallet.Wallet
	walletCount       int
	selectedWallet    *wallet.Wallet
	deletingWallet    *wallet.Wallet
	err               error
	nameInput         textinput.Model
	passwordInput     textinput.Model
	privateKeyInput   textinput.Model
	mnemonic          string
	walletTable       table.Model
	width             int
	height            int
	walletDetails     *wallet.WalletDetails
	styles            Styles
	fontsList         []string         // Lista de nomes de fontes carregadas do arquivo externo
	selectedFont      *tdf.TheDrawFont // Fonte selecionada aleatoriamente
	fontInfo          *tdf.FontInfo    // Informação da fonte selecionada
	dialogButtonIndex int              // 0 = Confirmar, 1 = Cancelar
}
