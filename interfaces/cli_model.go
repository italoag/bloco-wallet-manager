package interfaces

import (
	"blocowallet/domain"
	"blocowallet/usecases"
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/digitallyserviced/tdfgo/tdf"
)

type CLIModel struct {
	Service        *usecases.WalletService
	currentView    string
	menuItems      []menuItem
	selectedMenu   int
	importWords    []string
	importStage    int
	textInputs     []textinput.Model
	wallets        []domain.Wallet
	selectedWallet *domain.Wallet
	err            error
	passwordInput  textinput.Model
	deleteConfirmationInput textinput.Model
	mnemonic       string
	walletTable    table.Model
	operationTable table.Model
	width          int
	height         int
	walletDetails  *usecases.WalletDetails
	styles         Styles
	fontsList      []string         // Lista de nomes de fontes carregadas do arquivo externo
	selectedFont   *tdf.TheDrawFont // Fonte selecionada aleatoriamente
	fontInfo       *tdf.FontInfo    // Informação da fonte selecionada
}
