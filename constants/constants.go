package constants

import "time"

const (
	PasswordCharLimit         = 64
	PasswordWidth             = 30
	PasswordMinLength         = 8
	DefaultView               = "menu"
	SplashView                = "splash"
	CreateWalletView          = "create_wallet_password"
	ImportWalletView          = "import_wallet"
	ImportWalletPasswordView  = "import_wallet_password"
	ImportMethodSelectionView = "import_method_selection"
	ImportPrivateKeyView      = "import_private_key"
	ListWalletsView           = "list_wallets"
	SelectWalletOperationView = "wallet_select_operation_view"
	DeleteWalletConfirmationView = "delete_wallet_confirmation"
	WalletPasswordView        = "wallet_password"
	WalletDetailsView         = "wallet_details"
	StyleWidth                = 40
	StyleMargin               = 1
	ConfigFontsPath           = "config/fonts.json"
	SplashDuration            = 2 * time.Second
	ErrorFontNotFoundMessage  = "Fonte não encontrada nos diretórios especificados."
	MnemonicWordCount         = 12
)
