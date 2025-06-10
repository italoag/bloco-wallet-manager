package ui

import (
	"blocowallet/internal/wallet"
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	zone "github.com/lrstanley/bubblezone"
)

func TestNewWalletListComponent(t *testing.T) {
	component := NewWalletListComponent()

	if component.id != "wallet-list" {
		t.Errorf("Expected id to be 'wallet-list', got %s", component.id)
	}

	if component.selected != 0 {
		t.Errorf("Expected initial selected to be 0, got %d", component.selected)
	}

	if len(component.wallets) != 0 {
		t.Errorf("Expected initial wallets to be empty, got %d", len(component.wallets))
	}
}

func TestWalletListComponent_SetWallets(t *testing.T) {
	component := NewWalletListComponent()

	// Create test wallets
	testWallets := []*wallet.Wallet{
		{
			Name:    "Wallet 1",
			Address: "0x1234567890123456789012345678901234567890",
		},
		{
			Name:    "Wallet 2",
			Address: "0xabcdefabcdefabcdefabcdefabcdefabcdefabcd",
		},
	}

	component.SetWallets(testWallets)

	if len(component.wallets) != 2 {
		t.Errorf("Expected 2 wallets, got %d", len(component.wallets))
	}

	if component.wallets[0].Name != "Wallet 1" {
		t.Errorf("Expected first wallet name to be 'Wallet 1', got '%s'", component.wallets[0].Name)
	}

	if component.wallets[1].Name != "Wallet 2" {
		t.Errorf("Expected second wallet name to be 'Wallet 2', got '%s'", component.wallets[1].Name)
	}
}

func TestWalletListComponent_SetWallets_SelectedOutOfBounds(t *testing.T) {
	component := NewWalletListComponent()
	component.selected = 5 // Out of bounds

	// Create test wallets
	testWallets := []*wallet.Wallet{
		{
			Name:    "Wallet 1",
			Address: "0x1234567890123456789012345678901234567890",
		},
	}

	component.SetWallets(testWallets)

	// Should adjust selected to last valid index
	if component.selected != 0 {
		t.Errorf("Expected selected to be adjusted to 0, got %d", component.selected)
	}
}

func TestWalletListComponent_SetWallets_EmptyList(t *testing.T) {
	component := NewWalletListComponent()
	component.selected = 2

	component.SetWallets([]*wallet.Wallet{})

	// Should reset selected to 0 for empty list
	if component.selected != 0 {
		t.Errorf("Expected selected to be reset to 0 for empty list, got %d", component.selected)
	}
}

func TestWalletListComponent_SetSize(t *testing.T) {
	component := NewWalletListComponent()

	component.SetSize(800, 600)

	if component.width != 800 {
		t.Errorf("Expected width to be 800, got %d", component.width)
	}

	if component.height != 600 {
		t.Errorf("Expected height to be 600, got %d", component.height)
	}
}

func TestWalletListComponent_GetSelected(t *testing.T) {
	component := NewWalletListComponent()

	// Test with no wallets
	if component.GetSelected() != nil {
		t.Error("Expected nil for no wallets")
	}

	// Create test wallets
	testWallets := []*wallet.Wallet{
		{
			Name:    "Wallet 1",
			Address: "0x1234567890123456789012345678901234567890",
		},
		{
			Name:    "Wallet 2",
			Address: "0xabcdefabcdefabcdefabcdefabcdefabcdefabcd",
		},
	}

	component.SetWallets(testWallets)

	// Test valid selection
	selectedWallet := component.GetSelected()
	if selectedWallet == nil {
		t.Error("Expected valid wallet, got nil")
	}

	if selectedWallet.Name != "Wallet 1" {
		t.Errorf("Expected selected wallet name to be 'Wallet 1', got '%s'", selectedWallet.Name)
	}

	// Test different selection
	component.selected = 1
	selectedWallet = component.GetSelected()
	if selectedWallet.Name != "Wallet 2" {
		t.Errorf("Expected selected wallet name to be 'Wallet 2', got '%s'", selectedWallet.Name)
	}

	// Test out of bounds
	component.selected = 5
	if component.GetSelected() != nil {
		t.Error("Expected nil for out of bounds selection")
	}

	// Test negative selection
	component.selected = -1
	if component.GetSelected() != nil {
		t.Error("Expected nil for negative selection")
	}
}

func TestWalletListComponent_GetSelectedIndex(t *testing.T) {
	component := NewWalletListComponent()

	if component.GetSelectedIndex() != 0 {
		t.Errorf("Expected initial selected index to be 0, got %d", component.GetSelectedIndex())
	}

	component.selected = 3
	if component.GetSelectedIndex() != 3 {
		t.Errorf("Expected selected index to be 3, got %d", component.GetSelectedIndex())
	}
}

func TestWalletListComponent_Update_WindowSize(t *testing.T) {
	component := NewWalletListComponent()

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

func TestWalletListComponent_Update_UpKey(t *testing.T) {
	component := NewWalletListComponent()

	// Create test wallets
	testWallets := []*wallet.Wallet{
		{Name: "Wallet 1", Address: "0x1234567890123456789012345678901234567890"},
		{Name: "Wallet 2", Address: "0xabcdefabcdefabcdefabcdefabcdefabcdefabcd"},
		{Name: "Wallet 3", Address: "0x9876543210987654321098765432109876543210"},
	}
	component.SetWallets(testWallets)
	component.selected = 2

	upMsg := tea.KeyMsg{Type: tea.KeyUp}

	updatedComponent, cmd := component.Update(upMsg)

	if updatedComponent.selected != 1 {
		t.Errorf("Expected selected to be 1 after up key, got %d", updatedComponent.selected)
	}

	if cmd != nil {
		t.Error("Expected no command for up key")
	}
}

func TestWalletListComponent_Update_UpKeyAtTop(t *testing.T) {
	component := NewWalletListComponent()

	// Create test wallets
	testWallets := []*wallet.Wallet{
		{Name: "Wallet 1", Address: "0x1234567890123456789012345678901234567890"},
		{Name: "Wallet 2", Address: "0xabcdefabcdefabcdefabcdefabcdefabcdefabcd"},
	}
	component.SetWallets(testWallets)
	component.selected = 0

	upMsg := tea.KeyMsg{Type: tea.KeyUp}

	updatedComponent, cmd := component.Update(upMsg)

	if updatedComponent.selected != 0 {
		t.Errorf("Expected selected to stay 0 at top, got %d", updatedComponent.selected)
	}

	if cmd != nil {
		t.Error("Expected no command for up key at top")
	}
}

func TestWalletListComponent_Update_DownKey(t *testing.T) {
	component := NewWalletListComponent()

	// Create test wallets
	testWallets := []*wallet.Wallet{
		{Name: "Wallet 1", Address: "0x1234567890123456789012345678901234567890"},
		{Name: "Wallet 2", Address: "0xabcdefabcdefabcdefabcdefabcdefabcdefabcd"},
		{Name: "Wallet 3", Address: "0x9876543210987654321098765432109876543210"},
	}
	component.SetWallets(testWallets)
	component.selected = 0

	downMsg := tea.KeyMsg{Type: tea.KeyDown}

	updatedComponent, cmd := component.Update(downMsg)

	if updatedComponent.selected != 1 {
		t.Errorf("Expected selected to be 1 after down key, got %d", updatedComponent.selected)
	}

	if cmd != nil {
		t.Error("Expected no command for down key")
	}
}

func TestWalletListComponent_Update_DownKeyAtBottom(t *testing.T) {
	component := NewWalletListComponent()

	// Create test wallets
	testWallets := []*wallet.Wallet{
		{Name: "Wallet 1", Address: "0x1234567890123456789012345678901234567890"},
		{Name: "Wallet 2", Address: "0xabcdefabcdefabcdefabcdefabcdefabcdefabcd"},
	}
	component.SetWallets(testWallets)
	component.selected = 1 // Last item

	downMsg := tea.KeyMsg{Type: tea.KeyDown}

	updatedComponent, cmd := component.Update(downMsg)

	if updatedComponent.selected != 1 {
		t.Errorf("Expected selected to stay 1 at bottom, got %d", updatedComponent.selected)
	}

	if cmd != nil {
		t.Error("Expected no command for down key at bottom")
	}
}

func TestWalletListComponent_Update_VimKeys(t *testing.T) {
	component := NewWalletListComponent()

	// Create test wallets
	testWallets := []*wallet.Wallet{
		{Name: "Wallet 1", Address: "0x1234567890123456789012345678901234567890"},
		{Name: "Wallet 2", Address: "0xabcdefabcdefabcdefabcdefabcdefabcdefabcd"},
		{Name: "Wallet 3", Address: "0x9876543210987654321098765432109876543210"},
	}
	component.SetWallets(testWallets)
	component.selected = 1

	// Test 'k' (up in vim)
	kMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'k'}}
	updatedComponent, _ := component.Update(kMsg)

	if updatedComponent.selected != 0 {
		t.Errorf("Expected selected to be 0 after 'k' key, got %d", updatedComponent.selected)
	}

	// Test 'j' (down in vim)
	jMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}}
	updatedComponent, _ = updatedComponent.Update(jMsg)

	if updatedComponent.selected != 1 {
		t.Errorf("Expected selected to be 1 after 'j' key, got %d", updatedComponent.selected)
	}
}

func TestWalletListComponent_Update_EnterKey(t *testing.T) {
	component := NewWalletListComponent()

	// Create test wallets
	testWallets := []*wallet.Wallet{
		{Name: "Wallet 1", Address: "0x1234567890123456789012345678901234567890"},
		{Name: "Wallet 2", Address: "0xabcdefabcdefabcdefabcdefabcdefabcdefabcd"},
	}
	component.SetWallets(testWallets)
	component.selected = 1

	enterMsg := tea.KeyMsg{Type: tea.KeyEnter}

	updatedComponent, cmd := component.Update(enterMsg)

	if updatedComponent.selected != 1 {
		t.Errorf("Expected selected to remain 1, got %d", updatedComponent.selected)
	}

	if cmd == nil {
		t.Error("Expected command to be returned for enter key")
	}

	// Execute the command to test the message
	msg := cmd()
	if walletMsg, ok := msg.(WalletSelectedMsg); ok {
		if walletMsg.Wallet.Name != "Wallet 2" {
			t.Errorf("Expected selected wallet name to be 'Wallet 2', got '%s'", walletMsg.Wallet.Name)
		}
	} else {
		t.Error("Expected WalletSelectedMsg")
	}
}

func TestWalletListComponent_Update_EnterKeyOutOfBounds(t *testing.T) {
	component := NewWalletListComponent()

	// Create test wallets
	testWallets := []*wallet.Wallet{
		{Name: "Wallet 1", Address: "0x1234567890123456789012345678901234567890"},
	}
	component.SetWallets(testWallets)
	component.selected = 5 // Out of bounds

	enterMsg := tea.KeyMsg{Type: tea.KeyEnter}

	updatedComponent, cmd := component.Update(enterMsg)

	if updatedComponent.selected != 5 {
		t.Errorf("Expected selected to remain 5, got %d", updatedComponent.selected)
	}

	if cmd != nil {
		t.Error("Expected no command for out of bounds enter key")
	}
}

func TestWalletListComponent_Update_KeysWithNoWallets(t *testing.T) {
	component := NewWalletListComponent()

	testKeys := []tea.KeyMsg{
		{Type: tea.KeyUp},
		{Type: tea.KeyDown},
		{Type: tea.KeyEnter},
		{Type: tea.KeyRunes, Runes: []rune{'k'}},
		{Type: tea.KeyRunes, Runes: []rune{'j'}},
	}

	for _, keyMsg := range testKeys {
		updatedComponent, cmd := component.Update(keyMsg)

		if updatedComponent.selected != 0 {
			t.Errorf("Expected selected to remain 0 with no wallets for key, got %d", updatedComponent.selected)
		}

		if cmd != nil {
			t.Error("Expected no command with no wallets")
		}
	}
}

func TestWalletListComponent_Update_MouseClick(t *testing.T) {
	component := NewWalletListComponent()

	// Initialize bubblezone
	zone.NewGlobal()

	// Create test wallets
	testWallets := []*wallet.Wallet{
		{Name: "Wallet 1", Address: "0x1234567890123456789012345678901234567890"},
		{Name: "Wallet 2", Address: "0xabcdefabcdefabcdefabcdefabcdefabcdefabcd"},
	}
	component.SetWallets(testWallets)

	// Create a mouse release message
	mouseMsg := tea.MouseMsg{
		X:      10,
		Y:      5,
		Action: tea.MouseActionRelease,
		Button: tea.MouseButtonLeft,
	}

	// Test mouse click (we can't easily test the zone boundaries without proper setup)
	updatedComponent, cmd := component.Update(mouseMsg)

	// Should return the component unchanged if no zone is hit
	if updatedComponent.selected != component.selected {
		t.Error("Expected selection to remain unchanged for mouse click outside zones")
	}

	// Command should be nil if no zone is hit
	if cmd != nil {
		t.Error("Expected no command for mouse click outside zones")
	}
}

func TestWalletListComponent_Update_MouseClickWrongButton(t *testing.T) {
	component := NewWalletListComponent()

	// Test right mouse button (should be ignored)
	mouseMsg := tea.MouseMsg{
		X:      10,
		Y:      5,
		Action: tea.MouseActionRelease,
		Button: tea.MouseButtonRight,
	}

	updatedComponent, cmd := component.Update(mouseMsg)

	if updatedComponent.selected != component.selected {
		t.Error("Expected selection to remain unchanged for right mouse button")
	}

	if cmd != nil {
		t.Error("Expected no command for right mouse button")
	}
}

func TestWalletListComponent_View_EmptyList(t *testing.T) {
	component := NewWalletListComponent()

	view := component.View()

	expectedMessage := "No wallets found. Create a new wallet to get started."
	if !strings.Contains(view, expectedMessage) {
		t.Errorf("Expected view to contain empty message, got: %s", view)
	}
}

func TestWalletListComponent_View_WithWallets(t *testing.T) {
	component := NewWalletListComponent()
	component.SetSize(80, 24)

	// Create test wallets
	testWallets := []*wallet.Wallet{
		{Name: "My Wallet", Address: "0x1234567890123456789012345678901234567890"},
		{Name: "Second Wallet", Address: "0xabcdefabcdefabcdefabcdefabcdefabcdefabcd"},
	}
	component.SetWallets(testWallets)

	view := component.View()

	// Test wallet names are present
	if !strings.Contains(view, "My Wallet") {
		t.Error("Expected view to contain 'My Wallet'")
	}

	if !strings.Contains(view, "Second Wallet") {
		t.Error("Expected view to contain 'Second Wallet'")
	}

	// Test address abbreviation
	if !strings.Contains(view, "0x12345678...567890") {
		t.Error("Expected view to contain abbreviated address for first wallet")
	}

	if !strings.Contains(view, "0xabcdefab...efabcd") {
		t.Error("Expected view to contain abbreviated address for second wallet")
	}

	// Test selection indicator
	if !strings.Contains(view, "► My Wallet") {
		t.Error("Expected first wallet to be selected with indicator")
	}

	// Test non-selected formatting
	if strings.Contains(view, "► Second Wallet") {
		t.Error("Expected second wallet to not have selection indicator")
	}
}

func TestWalletListComponent_View_DifferentSelection(t *testing.T) {
	component := NewWalletListComponent()
	component.SetSize(80, 24)

	// Create test wallets
	testWallets := []*wallet.Wallet{
		{Name: "My Wallet", Address: "0x1234567890123456789012345678901234567890"},
		{Name: "Second Wallet", Address: "0xabcdefabcdefabcdefabcdefabcdefabcdefabcd"},
	}
	component.SetWallets(testWallets)
	component.selected = 1

	view := component.View()

	// Test that the selected item has the indicator
	if !strings.Contains(view, "► Second Wallet") {
		t.Error("Expected second wallet to be selected with indicator")
	}

	// Test that non-selected items don't have the indicator
	if strings.Contains(view, "► My Wallet") {
		t.Error("Expected first wallet to not have indicator when not selected")
	}
}

func TestWalletListComponent_Update_UnknownKey(t *testing.T) {
	component := NewWalletListComponent()

	// Create test wallets
	testWallets := []*wallet.Wallet{
		{Name: "Wallet 1", Address: "0x1234567890123456789012345678901234567890"},
	}
	component.SetWallets(testWallets)
	initialSelected := component.selected

	unknownMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'x'}}

	updatedComponent, cmd := component.Update(unknownMsg)

	if updatedComponent.selected != initialSelected {
		t.Errorf("Expected selection to remain unchanged for unknown key, got %d", updatedComponent.selected)
	}

	if cmd != nil {
		t.Error("Expected no command for unknown key")
	}
}

func TestWalletListComponent_Update_OtherMessageType(t *testing.T) {
	component := NewWalletListComponent()

	// Create test wallets
	testWallets := []*wallet.Wallet{
		{Name: "Wallet 1", Address: "0x1234567890123456789012345678901234567890"},
	}
	component.SetWallets(testWallets)
	initialSelected := component.selected

	// Test with some other message type
	otherMsg := "some random message"

	updatedComponent, cmd := component.Update(otherMsg)

	if updatedComponent.selected != initialSelected {
		t.Errorf("Expected selection to remain unchanged for other message type, got %d", updatedComponent.selected)
	}

	if cmd != nil {
		t.Error("Expected no command for other message type")
	}
}

func TestWalletListComponent_SetWallets_BoundaryConditions(t *testing.T) {
	component := NewWalletListComponent()

	// Test setting wallets multiple times
	testWallets1 := []*wallet.Wallet{
		{Name: "Wallet 1", Address: "0x1111111111111111111111111111111111111111"},
		{Name: "Wallet 2", Address: "0x2222222222222222222222222222222222222222"},
		{Name: "Wallet 3", Address: "0x3333333333333333333333333333333333333333"},
	}
	component.SetWallets(testWallets1)
	component.selected = 2 // Select last wallet

	// Now set to fewer wallets
	testWallets2 := []*wallet.Wallet{
		{Name: "Wallet A", Address: "0xaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"},
	}
	component.SetWallets(testWallets2)

	// Should adjust selected to last valid index
	if component.selected != 0 {
		t.Errorf("Expected selected to be adjusted to 0, got %d", component.selected)
	}

	if len(component.wallets) != 1 {
		t.Errorf("Expected 1 wallet, got %d", len(component.wallets))
	}

	if component.wallets[0].Name != "Wallet A" {
		t.Errorf("Expected wallet name to be 'Wallet A', got '%s'", component.wallets[0].Name)
	}
}

func TestWalletListComponent_GetSelected_EdgeCases(t *testing.T) {
	component := NewWalletListComponent()

	// Test with single wallet
	testWallets := []*wallet.Wallet{
		{Name: "Only Wallet", Address: "0x1111111111111111111111111111111111111111"},
	}
	component.SetWallets(testWallets)

	selected := component.GetSelected()
	if selected == nil {
		t.Error("Expected to get the only wallet")
	}
	if selected.Name != "Only Wallet" {
		t.Errorf("Expected wallet name to be 'Only Wallet', got '%s'", selected.Name)
	}

	// Test with large selection index
	component.selected = 1000
	selected = component.GetSelected()
	if selected != nil {
		t.Error("Expected nil for large selection index")
	}
}

func TestWalletSelectedMsg(t *testing.T) {
	testWallet := &wallet.Wallet{
		Name:    "Test Wallet",
		Address: "0x1234567890123456789012345678901234567890",
	}

	msg := WalletSelectedMsg{
		Wallet: testWallet,
	}

	if msg.Wallet.Name != "Test Wallet" {
		t.Errorf("Expected wallet name to be 'Test Wallet', got '%s'", msg.Wallet.Name)
	}

	if msg.Wallet.Address != "0x1234567890123456789012345678901234567890" {
		t.Errorf("Expected wallet address to match, got '%s'", msg.Wallet.Address)
	}
}
