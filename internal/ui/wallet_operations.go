package ui

import (
	"context"
	"fmt"
	"strings"

	"blocowallet/internal/wallet"

	tea "github.com/charmbracelet/bubbletea"
)

// WalletOperations contém operações específicas de wallet que podem ser reutilizadas
type WalletOperations struct {
	service *wallet.Service
}

// NewWalletOperations cria uma nova instância de WalletOperations
func NewWalletOperations(service *wallet.Service) *WalletOperations {
	return &WalletOperations{
		service: service,
	}
}

// LoadWallets carrega todas as wallets
func (wo *WalletOperations) LoadWallets() ([]*wallet.Wallet, error) {
	return wo.service.List(context.Background())
}

// CreateWallet cria uma nova wallet com mnemônico
func (wo *WalletOperations) CreateWallet(name, password string) (*wallet.WalletDetails, error) {
	return wo.service.CreateWalletWithMnemonic(context.Background(), name, password)
}

// ImportWallet importa uma wallet usando mnemônico
func (wo *WalletOperations) ImportWallet(name, password, mnemonic string) (*wallet.WalletDetails, error) {
	return wo.service.ImportWalletFromMnemonic(context.Background(), name, mnemonic, password)
}

// ImportWalletFromPrivateKey importa uma wallet usando chave privada
func (wo *WalletOperations) ImportWalletFromPrivateKey(name, password, privateKey string) (*wallet.WalletDetails, error) {
	return wo.service.ImportWalletFromPrivateKey(context.Background(), name, privateKey, password)
}

// GetBalance obtém o saldo de uma wallet
func (wo *WalletOperations) GetBalance(address string) (*wallet.Balance, error) {
	return wo.service.GetBalance(context.Background(), address)
}

// GetMultiNetworkBalance obtém o saldo em múltiplas redes
func (wo *WalletOperations) GetMultiNetworkBalance(address string) (*wallet.MultiNetworkBalance, error) {
	return wo.service.GetMultiNetworkBalance(context.Background(), address)
}

// DeleteWallet exclui uma wallet pelo endereço
func (wo *WalletOperations) DeleteWallet(address string) error {
	return wo.service.DeleteWalletByAddress(context.Background(), address)
}

// ValidateWalletInput valida os dados de entrada para criação de wallet
func ValidateWalletInput(name, password string) error {
	name = strings.TrimSpace(name)
	password = strings.TrimSpace(password)

	if name == "" {
		return fmt.Errorf("wallet name is required")
	}
	if password == "" {
		return fmt.Errorf("password is required")
	}
	if len(password) < 6 {
		return fmt.Errorf("password must be at least 6 characters long")
	}
	return nil
}

// ValidateImportInput valida os dados de entrada para importação de wallet
func ValidateImportInput(name, password, mnemonic string) error {
	if err := ValidateWalletInput(name, password); err != nil {
		return err
	}

	mnemonic = strings.TrimSpace(mnemonic)
	if mnemonic == "" {
		return fmt.Errorf("mnemonic phrase is required")
	}

	words := strings.Fields(mnemonic)
	if len(words) != 12 && len(words) != 24 {
		return fmt.Errorf("mnemonic phrase must be 12 or 24 words")
	}

	return nil
}

// ValidatePrivateKeyInput valida os dados de entrada para importação via chave privada
func ValidatePrivateKeyInput(name, password, privateKey string) error {
	if err := ValidateWalletInput(name, password); err != nil {
		return err
	}

	privateKey = strings.TrimSpace(privateKey)
	if privateKey == "" {
		return fmt.Errorf("private key is required")
	}

	// Remove 0x prefix if present
	if strings.HasPrefix(privateKey, "0x") {
		privateKey = privateKey[2:]
	}

	if len(privateKey) != 64 {
		return fmt.Errorf("private key must be 64 hexadecimal characters")
	}

	// Validate hex characters
	for _, c := range privateKey {
		if !((c >= '0' && c <= '9') || (c >= 'a' && c <= 'f') || (c >= 'A' && c <= 'F')) {
			return fmt.Errorf("private key must contain only hexadecimal characters")
		}
	}

	return nil
}

// FormatWalletAddress formata um endereço de wallet para exibição
func FormatWalletAddress(address string) string {
	if len(address) <= 13 {
		return address
	}
	return address[:10] + "..."
}

// WalletCommand types para comandos específicos de wallet
type WalletLoadedMsg []*wallet.Wallet
type WalletBalanceLoadedMsg *wallet.Balance
type WalletMultiBalanceLoadedMsg *wallet.MultiNetworkBalance
type WalletCreatedMsg struct{}
type WalletDeletedMsg struct{}
type WalletAuthenticatedMsg struct{}

// Wallet command creators
func CreateLoadWalletsCmd(wo *WalletOperations) tea.Cmd {
	return tea.Cmd(func() tea.Msg {
		wallets, err := wo.LoadWallets()
		if err != nil {
			return errorMsg(err.Error())
		}
		return WalletLoadedMsg(wallets)
	})
}

func CreateGetBalanceCmd(wo *WalletOperations, address string) tea.Cmd {
	return tea.Cmd(func() tea.Msg {
		balance, err := wo.GetBalance(address)
		if err != nil {
			return errorMsg(err.Error())
		}
		return WalletBalanceLoadedMsg(balance)
	})
}

func CreateGetMultiBalanceCmd(wo *WalletOperations, address string) tea.Cmd {
	return tea.Cmd(func() tea.Msg {
		multiBalance, err := wo.GetMultiNetworkBalance(address)
		if err != nil {
			return errorMsg(err.Error())
		}
		return WalletMultiBalanceLoadedMsg(multiBalance)
	})
}

func CreateWalletCmd(wo *WalletOperations, name, password string) tea.Cmd {
	return tea.Cmd(func() tea.Msg {
		if err := ValidateWalletInput(name, password); err != nil {
			return errorMsg(err.Error())
		}

		_, err := wo.CreateWallet(name, password)
		if err != nil {
			return errorMsg(err.Error())
		}
		return WalletCreatedMsg{}
	})
}

func ImportWalletCmd(wo *WalletOperations, name, password, mnemonic string) tea.Cmd {
	return tea.Cmd(func() tea.Msg {
		if err := ValidateImportInput(name, password, mnemonic); err != nil {
			return errorMsg(err.Error())
		}

		_, err := wo.ImportWallet(name, password, mnemonic)
		if err != nil {
			return errorMsg(err.Error())
		}
		return WalletCreatedMsg{}
	})
}

func ImportPrivateKeyCmd(wo *WalletOperations, name, password, privateKey string) tea.Cmd {
	return tea.Cmd(func() tea.Msg {
		if err := ValidatePrivateKeyInput(name, password, privateKey); err != nil {
			return errorMsg(err.Error())
		}

		_, err := wo.ImportWalletFromPrivateKey(name, password, privateKey)
		if err != nil {
			return errorMsg(err.Error())
		}
		return WalletCreatedMsg{}
	})
}

func DeleteWalletCmd(wo *WalletOperations, address string) tea.Cmd {
	return tea.Cmd(func() tea.Msg {
		if err := wo.DeleteWallet(address); err != nil {
			return errorMsg(err.Error())
		}
		return WalletDeletedMsg{}
	})
}
