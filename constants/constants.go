package constants

import "time"

const (
	PasswordCharLimit        = 64
	PasswordWidth            = 30
	PasswordMinLength        = 8
	DefaultView              = "menu"
	SplashView               = "splash"
	CreateWalletView         = "create_wallet_password"
	ImportWalletView         = "import_wallet"
	ImportWalletPasswordView = "import_wallet_password"
	ListWalletsView          = "list_wallets"
	WalletPasswordView       = "wallet_password"
	WalletDetailsView        = "wallet_details"
	StyleWidth               = 40
	StyleMargin              = 1
	ConfigFontsPath          = "config/fonts.json"
	SplashDuration           = 2 * time.Second
	ErrorFontNotFoundMessage = "Fonte não encontrada nos diretórios especificados."
)
