package ui

import (
	"blocowallet/internal/blockchain"
	"blocowallet/internal/wallet"
	"blocowallet/pkg/config"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/digitallyserviced/tdfgo/tdf"
)

// View constants
const (
	SplashView           = "splash"
	MenuView             = "menu"
	WalletListView       = "wallet_list"
	WalletAuthView       = "wallet_auth"
	WalletDetailsView    = "wallet_details"
	CreateWalletView     = "create_wallet"
	ImportWalletView     = "import_wallet"
	ImportPrivateKeyView = "import_private_key"
	SettingsView         = "settings"
	NetworkConfigView    = "network_config"
	LanguageView         = "language"
	AddNetworkView       = "add_network"
	NetworkDetailsView   = "network_details"
)

// Model represents the TUI application state
type Model struct {
	walletService       *wallet.Service
	wallets             []*wallet.Wallet
	selected            int
	loading             bool
	err                 error
	width               int
	height              int
	currentView         string
	selectedWallet      *wallet.Wallet
	currentBalance      *wallet.Balance
	currentMultiBalance *wallet.MultiNetworkBalance

	// Configuration
	config *config.Config

	// Components
	splashComponent           SplashComponent
	mainMenuComponent         MainMenuComponent
	settingsMenuComponent     SettingsMenuComponent
	networkListComponent      NetworkListComponent
	addNetworkComponent       AddNetworkComponent
	languageMenuComponent     LanguageMenuComponent
	walletListComponent       WalletListComponent
	balanceComponent          BalanceComponent
	createWalletComponent     CreateWalletComponent
	importMnemonicComponent   ImportMnemonicComponent
	importPrivateKeyComponent ImportPrivateKeyComponent

	// Loading and table components
	loadingSpinner spinner.Model
	isLoading      bool
	loadingText    string
	walletTable    table.Model

	// Legacy input fields (to be removed gradually)
	nameInput       textinput.Model
	passwordInput   textinput.Model
	mnemonicInput   textinput.Model
	privateKeyInput textinput.Model
	inputFocus      int

	// Private key viewing for imported wallets
	privateKeyPassword  textinput.Model
	extractedPrivateKey string
	privateKeyError     string

	// Wallet deletion dialog
	deleteDialog *DeleteWalletDialog

	// Wallet authentication for keystore access
	needsWalletAuth    bool
	walletAuthPassword textinput.Model
	walletAuthError    string

	// Settings fields
	settingsSelected int
	networkSelected  int
	languageSelected int
	editingRPC       bool
	rpcInput         textinput.Model

	// Add network fields
	networkNameInput   textinput.Model
	chainIDInput       textinput.Model
	rpcEndpointInput   textinput.Model
	addNetworkFocus    int
	addingNetwork      bool
	selectedNetworkKey string
	networkSuggestions []blockchain.NetworkSuggestion
	showingSuggestions bool
	selectedSuggestion int
	chainListService   *blockchain.ChainListService

	// Settings menu items
	settingsItems []string
	networkItems  []string
	languageItems []string

	// TDF font support
	selectedFont *tdf.TheDrawFont
	fontName     string

	// Sensitive information visibility
	showSensitiveInfo bool
}

// Message types
type walletsLoadedMsg []*wallet.Wallet
type balanceLoadedMsg *wallet.Balance
type multiBalanceLoadedMsg *wallet.MultiNetworkBalance
type walletCreatedMsg struct{}
type errorMsg string
type networkSuggestionsMsg []blockchain.NetworkSuggestion
type chainInfoLoadedMsg struct {
	chainInfo *blockchain.ChainInfo
	rpcURL    string
}
