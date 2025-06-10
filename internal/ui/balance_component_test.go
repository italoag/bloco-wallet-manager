package ui

import (
	"blocowallet/internal/wallet"
	"fmt"
	"math/big"
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

func TestNewBalanceComponent(t *testing.T) {
	component := NewBalanceComponent()

	if component.id != "balance-display" {
		t.Errorf("Expected id to be 'balance-display', got %s", component.id)
	}

	if component.selectedWallet != nil {
		t.Error("Expected selectedWallet to be nil initially")
	}

	if component.balance != nil {
		t.Error("Expected balance to be nil initially")
	}

	if component.multiBalance != nil {
		t.Error("Expected multiBalance to be nil initially")
	}

	if component.loading {
		t.Error("Expected loading to be false initially")
	}

	if component.err != nil {
		t.Error("Expected err to be nil initially")
	}
}

func TestBalanceComponent_SetWallet(t *testing.T) {
	component := NewBalanceComponent()

	// Set some initial state
	component.balance = &wallet.Balance{Amount: big.NewInt(1000), Symbol: "ETH"}
	component.multiBalance = &wallet.MultiNetworkBalance{}
	component.err = fmt.Errorf("test error")

	testWallet := &wallet.Wallet{
		Name:    "Test Wallet",
		Address: "0x1234567890123456789012345678901234567890",
	}

	component.SetWallet(testWallet)

	if component.selectedWallet != testWallet {
		t.Error("Expected selectedWallet to be set")
	}

	if component.balance != nil {
		t.Error("Expected balance to be reset to nil")
	}

	if component.multiBalance != nil {
		t.Error("Expected multiBalance to be reset to nil")
	}

	if component.err != nil {
		t.Error("Expected err to be reset to nil")
	}
}

func TestBalanceComponent_SetBalance(t *testing.T) {
	component := NewBalanceComponent()
	component.loading = true
	component.err = fmt.Errorf("test error")

	testBalance := &wallet.Balance{
		Amount: big.NewInt(1000000000000000000), // 1 ETH in wei
		Symbol: "ETH",
	}

	component.SetBalance(testBalance)

	if component.balance != testBalance {
		t.Error("Expected balance to be set")
	}

	if component.loading {
		t.Error("Expected loading to be false after setting balance")
	}

	if component.err != nil {
		t.Error("Expected error to be cleared after setting balance")
	}
}

func TestBalanceComponent_SetMultiBalance(t *testing.T) {
	component := NewBalanceComponent()
	component.loading = true
	component.err = fmt.Errorf("test error")

	testMultiBalance := &wallet.MultiNetworkBalance{
		NetworkBalances: []*wallet.NetworkBalance{
			{
				NetworkName: "Ethereum",
				Symbol:      "ETH",
				Amount:      big.NewInt(1000000000000000000),
			},
		},
	}

	component.SetMultiBalance(testMultiBalance)

	if component.multiBalance != testMultiBalance {
		t.Error("Expected multiBalance to be set")
	}

	if component.loading {
		t.Error("Expected loading to be false after setting multiBalance")
	}

	if component.err != nil {
		t.Error("Expected error to be cleared after setting multiBalance")
	}
}

func TestBalanceComponent_SetLoading(t *testing.T) {
	component := NewBalanceComponent()
	component.err = fmt.Errorf("test error")

	component.SetLoading(true)

	if !component.loading {
		t.Error("Expected loading to be true")
	}

	if component.err != nil {
		t.Error("Expected error to be cleared when setting loading to true")
	}

	component.SetLoading(false)

	if component.loading {
		t.Error("Expected loading to be false")
	}
}

func TestBalanceComponent_SetError(t *testing.T) {
	component := NewBalanceComponent()
	component.loading = true

	testError := fmt.Errorf("test error")
	component.SetError(testError)

	if component.err != testError {
		t.Error("Expected error to be set")
	}

	if component.loading {
		t.Error("Expected loading to be false after setting error")
	}
}

func TestBalanceComponent_SetSize(t *testing.T) {
	component := NewBalanceComponent()

	component.SetSize(800, 600)

	if component.width != 800 {
		t.Errorf("Expected width to be 800, got %d", component.width)
	}

	if component.height != 600 {
		t.Errorf("Expected height to be 600, got %d", component.height)
	}
}

func TestBalanceComponent_Update_WindowSize(t *testing.T) {
	component := NewBalanceComponent()

	msg := tea.WindowSizeMsg{Width: 1024, Height: 768}

	updatedComponent, cmd := component.Update(msg)

	if updatedComponent.width != 1024 {
		t.Errorf("Expected width to be 1024, got %d", updatedComponent.width)
	}

	if updatedComponent.height != 768 {
		t.Errorf("Expected height to be 768, got %d", updatedComponent.height)
	}

	if cmd != nil {
		t.Error("Expected no command for window size message")
	}
}

func TestBalanceComponent_Update_BalanceLoadedMsg(t *testing.T) {
	component := NewBalanceComponent()
	component.loading = true
	component.err = fmt.Errorf("test error")

	testBalance := &wallet.Balance{
		Amount: big.NewInt(1000000000000000000),
		Symbol: "ETH",
	}

	balanceMsg := balanceLoadedMsg(testBalance)

	updatedComponent, cmd := component.Update(balanceMsg)

	if updatedComponent.balance != testBalance {
		t.Error("Expected balance to be set from message")
	}

	if updatedComponent.loading {
		t.Error("Expected loading to be false after balance loaded")
	}

	if updatedComponent.err != nil {
		t.Error("Expected error to be cleared after balance loaded")
	}

	if cmd != nil {
		t.Error("Expected no command for balance loaded message")
	}
}

func TestBalanceComponent_Update_MultiBalanceLoadedMsg(t *testing.T) {
	component := NewBalanceComponent()
	component.loading = true
	component.err = fmt.Errorf("test error")

	testMultiBalance := &wallet.MultiNetworkBalance{
		NetworkBalances: []*wallet.NetworkBalance{
			{
				NetworkName: "Ethereum",
				Symbol:      "ETH",
				Amount:      big.NewInt(1000000000000000000),
			},
		},
	}

	multiBalanceMsg := multiBalanceLoadedMsg(testMultiBalance)

	updatedComponent, cmd := component.Update(multiBalanceMsg)

	if updatedComponent.multiBalance != testMultiBalance {
		t.Error("Expected multiBalance to be set from message")
	}

	if updatedComponent.loading {
		t.Error("Expected loading to be false after multiBalance loaded")
	}

	if updatedComponent.err != nil {
		t.Error("Expected error to be cleared after multiBalance loaded")
	}

	if cmd != nil {
		t.Error("Expected no command for multiBalance loaded message")
	}
}

func TestBalanceComponent_Update_ErrorMsg(t *testing.T) {
	component := NewBalanceComponent()
	component.loading = true

	errorMessage := errorMsg("Test error message")

	updatedComponent, cmd := component.Update(errorMessage)

	if updatedComponent.err == nil {
		t.Error("Expected error to be set from message")
	}

	if updatedComponent.err.Error() != "Test error message" {
		t.Errorf("Expected error message to be 'Test error message', got '%s'", updatedComponent.err.Error())
	}

	if updatedComponent.loading {
		t.Error("Expected loading to be false after error message")
	}

	if cmd != nil {
		t.Error("Expected no command for error message")
	}
}

func TestBalanceComponent_View_NoWallet(t *testing.T) {
	component := NewBalanceComponent()

	view := component.View()

	expectedMessage := "Select a wallet to view balance information."
	if !strings.Contains(view, expectedMessage) {
		t.Errorf("Expected view to contain no wallet message, got: %s", view)
	}
}

func TestBalanceComponent_View_WithWallet(t *testing.T) {
	component := NewBalanceComponent()

	testWallet := &wallet.Wallet{
		Name:    "My Test Wallet",
		Address: "0x1234567890123456789012345678901234567890",
	}

	component.SetWallet(testWallet)

	view := component.View()

	if !strings.Contains(view, "💼 My Test Wallet") {
		t.Error("Expected view to contain wallet name with icon")
	}

	if !strings.Contains(view, "Address: 0x1234567890123456789012345678901234567890") {
		t.Error("Expected view to contain wallet address")
	}

	// Should show no balance information message by default
	if !strings.Contains(view, "No balance information available.") {
		t.Error("Expected view to contain no balance information message")
	}
}

func TestBalanceComponent_View_Loading(t *testing.T) {
	component := NewBalanceComponent()

	testWallet := &wallet.Wallet{
		Name:    "My Test Wallet",
		Address: "0x1234567890123456789012345678901234567890",
	}

	component.SetWallet(testWallet)
	component.SetLoading(true)

	view := component.View()

	if !strings.Contains(view, "⏳ Loading balance...") {
		t.Error("Expected view to contain loading message")
	}

	// Should not contain other content when loading
	if strings.Contains(view, "No balance information available.") {
		t.Error("Expected view to not contain no balance message when loading")
	}
}

func TestBalanceComponent_View_Error(t *testing.T) {
	component := NewBalanceComponent()

	testWallet := &wallet.Wallet{
		Name:    "My Test Wallet",
		Address: "0x1234567890123456789012345678901234567890",
	}

	component.SetWallet(testWallet)
	component.SetError(fmt.Errorf("Network connection failed"))

	view := component.View()

	if !strings.Contains(view, "❌ Error loading balance: Network connection failed") {
		t.Error("Expected view to contain error message")
	}

	// Should not contain other content when showing error
	if strings.Contains(view, "No balance information available.") {
		t.Error("Expected view to not contain no balance message when showing error")
	}
}

func TestBalanceComponent_View_SingleBalance(t *testing.T) {
	component := NewBalanceComponent()

	testWallet := &wallet.Wallet{
		Name:    "My Test Wallet",
		Address: "0x1234567890123456789012345678901234567890",
	}

	testBalance := &wallet.Balance{
		Amount: big.NewInt(1500000000000000000), // 1.5 ETH in wei
		Symbol: "ETH",
	}

	component.SetWallet(testWallet)
	component.SetBalance(testBalance)

	view := component.View()

	if !strings.Contains(view, "💰 Balance:") {
		t.Error("Expected view to contain balance label")
	}

	if !strings.Contains(view, "1.500000") {
		t.Error("Expected view to contain formatted balance amount")
	}

	if !strings.Contains(view, "ETH") {
		t.Error("Expected view to contain currency symbol")
	}
}

func TestBalanceComponent_View_SingleBalanceNil(t *testing.T) {
	component := NewBalanceComponent()

	testWallet := &wallet.Wallet{
		Name:    "My Test Wallet",
		Address: "0x1234567890123456789012345678901234567890",
	}

	testBalance := &wallet.Balance{
		Amount: nil, // Nil amount
		Symbol: "ETH",
	}

	component.SetWallet(testWallet)
	component.SetBalance(testBalance)

	view := component.View()

	if !strings.Contains(view, "💰 Balance:") {
		t.Error("Expected view to contain balance label")
	}

	if !strings.Contains(view, "0") {
		t.Error("Expected view to show 0 for nil balance amount")
	}

	if !strings.Contains(view, "ETH") {
		t.Error("Expected view to contain currency symbol")
	}
}

func TestBalanceComponent_View_MultiBalance(t *testing.T) {
	component := NewBalanceComponent()

	testWallet := &wallet.Wallet{
		Name:    "My Test Wallet",
		Address: "0x1234567890123456789012345678901234567890",
	}

	testMultiBalance := &wallet.MultiNetworkBalance{
		NetworkBalances: []*wallet.NetworkBalance{
			{
				NetworkName: "Ethereum",
				Symbol:      "ETH",
				Amount:      big.NewInt(1000000000000000000), // 1 ETH
			},
			{
				NetworkName: "Polygon",
				Symbol:      "MATIC",
				Amount:      big.NewInt(2000000000000000000), // 2 MATIC
			},
		},
	}

	component.SetWallet(testWallet)
	component.SetMultiBalance(testMultiBalance)

	view := component.View()

	if !strings.Contains(view, "🌐 Multi-Network Balance") {
		t.Error("Expected view to contain multi-network balance header")
	}

	if !strings.Contains(view, "💰") {
		t.Error("Expected view to contain balance icons")
	}

	if !strings.Contains(view, "Ethereum") {
		t.Error("Expected view to contain Ethereum network name")
	}

	if !strings.Contains(view, "Polygon") {
		t.Error("Expected view to contain Polygon network name")
	}

	if !strings.Contains(view, "1.000000") {
		t.Error("Expected view to contain ETH balance")
	}

	if !strings.Contains(view, "2.000000") {
		t.Error("Expected view to contain MATIC balance")
	}

	if !strings.Contains(view, "ETH") {
		t.Error("Expected view to contain ETH symbol")
	}

	if !strings.Contains(view, "MATIC") {
		t.Error("Expected view to contain MATIC symbol")
	}
}

func TestBalanceComponent_View_MultiBalanceWithError(t *testing.T) {
	component := NewBalanceComponent()

	testWallet := &wallet.Wallet{
		Name:    "My Test Wallet",
		Address: "0x1234567890123456789012345678901234567890",
	}

	testMultiBalance := &wallet.MultiNetworkBalance{
		NetworkBalances: []*wallet.NetworkBalance{
			{
				NetworkName: "Ethereum",
				Symbol:      "ETH",
				Amount:      big.NewInt(1000000000000000000),
			},
			{
				NetworkName: "Polygon",
				Symbol:      "MATIC",
				Error:       fmt.Errorf("Network unavailable"),
			},
		},
	}

	component.SetWallet(testWallet)
	component.SetMultiBalance(testMultiBalance)

	view := component.View()

	if !strings.Contains(view, "Ethereum") {
		t.Error("Expected view to contain Ethereum network name")
	}

	if !strings.Contains(view, "Polygon") {
		t.Error("Expected view to contain Polygon network name")
	}

	if !strings.Contains(view, "1.000000") {
		t.Error("Expected view to contain successful ETH balance")
	}

	if !strings.Contains(view, "❌") {
		t.Error("Expected view to contain error icon for failed network")
	}

	if !strings.Contains(view, "Network unavailable") {
		t.Error("Expected view to contain error message")
	}
}

func TestBalanceComponent_View_MultiBalanceNilAmount(t *testing.T) {
	component := NewBalanceComponent()

	testWallet := &wallet.Wallet{
		Name:    "My Test Wallet",
		Address: "0x1234567890123456789012345678901234567890",
	}

	testMultiBalance := &wallet.MultiNetworkBalance{
		NetworkBalances: []*wallet.NetworkBalance{
			{
				NetworkName: "Ethereum",
				Symbol:      "ETH",
				Amount:      nil, // Nil amount
			},
		},
	}

	component.SetWallet(testWallet)
	component.SetMultiBalance(testMultiBalance)

	view := component.View()

	if !strings.Contains(view, "Ethereum") {
		t.Error("Expected view to contain Ethereum network name")
	}

	if !strings.Contains(view, "0") {
		t.Error("Expected view to show 0 for nil amount")
	}

	if !strings.Contains(view, "ETH") {
		t.Error("Expected view to contain ETH symbol")
	}
}
