package ui

import (
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	zone "github.com/lrstanley/bubblezone"
)

func TestNewDeleteWalletDialog(t *testing.T) {
	walletName := "My Test Wallet"
	address := "0x1234567890123456789012345678901234567890"

	dialog := NewDeleteWalletDialog(walletName, address)

	if dialog.id != "delete-wallet-dialog" {
		t.Errorf("Expected id to be 'delete-wallet-dialog', got %s", dialog.id)
	}

	if dialog.active != "cancel" {
		t.Errorf("Expected initial active to be 'cancel', got %s", dialog.active)
	}

	if dialog.walletName != walletName {
		t.Errorf("Expected walletName to be '%s', got '%s'", walletName, dialog.walletName)
	}

	if dialog.address != address {
		t.Errorf("Expected address to be '%s', got '%s'", address, dialog.address)
	}

	if dialog.question != "Are you sure you want to delete this wallet?" {
		t.Errorf("Expected default question, got '%s'", dialog.question)
	}
}

func TestDeleteWalletDialog_Init(t *testing.T) {
	dialog := NewDeleteWalletDialog("Test", "0x123")

	cmd := dialog.Init()

	if cmd != nil {
		t.Error("Expected Init to return nil command")
	}
}

func TestDeleteWalletDialog_Update_WindowSize(t *testing.T) {
	dialog := NewDeleteWalletDialog("Test", "0x123")

	msg := tea.WindowSizeMsg{Width: 800, Height: 600}

	updatedModel, cmd := dialog.Update(msg)

	updatedDialog, ok := updatedModel.(DeleteWalletDialog)
	if !ok {
		t.Fatal("Expected returned model to be DeleteWalletDialog")
	}

	if updatedDialog.width != 800 {
		t.Errorf("Expected width to be 800, got %d", updatedDialog.width)
	}

	if cmd != nil {
		t.Error("Expected no command for window size message")
	}
}

func TestDeleteWalletDialog_Update_LeftKey(t *testing.T) {
	dialog := NewDeleteWalletDialog("Test", "0x123")
	dialog.active = "confirm"

	leftMsg := tea.KeyMsg{Type: tea.KeyLeft}

	updatedModel, cmd := dialog.Update(leftMsg)

	updatedDialog, ok := updatedModel.(DeleteWalletDialog)
	if !ok {
		t.Fatal("Expected returned model to be DeleteWalletDialog")
	}

	if updatedDialog.active != "cancel" {
		t.Errorf("Expected active to change to 'cancel', got '%s'", updatedDialog.active)
	}

	if cmd != nil {
		t.Error("Expected no command for left key")
	}
}

func TestDeleteWalletDialog_Update_LeftKeyFromCancel(t *testing.T) {
	dialog := NewDeleteWalletDialog("Test", "0x123")
	dialog.active = "cancel"

	leftMsg := tea.KeyMsg{Type: tea.KeyLeft}

	updatedModel, cmd := dialog.Update(leftMsg)

	updatedDialog, ok := updatedModel.(DeleteWalletDialog)
	if !ok {
		t.Fatal("Expected returned model to be DeleteWalletDialog")
	}

	if updatedDialog.active != "confirm" {
		t.Errorf("Expected active to change to 'confirm', got '%s'", updatedDialog.active)
	}

	if cmd != nil {
		t.Error("Expected no command for left key")
	}
}

func TestDeleteWalletDialog_Update_RightKey(t *testing.T) {
	dialog := NewDeleteWalletDialog("Test", "0x123")
	dialog.active = "cancel"

	rightMsg := tea.KeyMsg{Type: tea.KeyRight}

	updatedModel, cmd := dialog.Update(rightMsg)

	updatedDialog, ok := updatedModel.(DeleteWalletDialog)
	if !ok {
		t.Fatal("Expected returned model to be DeleteWalletDialog")
	}

	if updatedDialog.active != "confirm" {
		t.Errorf("Expected active to change to 'confirm', got '%s'", updatedDialog.active)
	}

	if cmd != nil {
		t.Error("Expected no command for right key")
	}
}

func TestDeleteWalletDialog_Update_RightKeyFromConfirm(t *testing.T) {
	dialog := NewDeleteWalletDialog("Test", "0x123")
	dialog.active = "confirm"

	rightMsg := tea.KeyMsg{Type: tea.KeyRight}

	updatedModel, cmd := dialog.Update(rightMsg)

	updatedDialog, ok := updatedModel.(DeleteWalletDialog)
	if !ok {
		t.Fatal("Expected returned model to be DeleteWalletDialog")
	}

	if updatedDialog.active != "cancel" {
		t.Errorf("Expected active to change to 'cancel', got '%s'", updatedDialog.active)
	}

	if cmd != nil {
		t.Error("Expected no command for right key")
	}
}

func TestDeleteWalletDialog_Update_VimKeys(t *testing.T) {
	dialog := NewDeleteWalletDialog("Test", "0x123")
	dialog.active = "confirm"

	// Test 'h' (left in vim)
	hMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'h'}}
	updatedModel, _ := dialog.Update(hMsg)
	updatedDialog, _ := updatedModel.(DeleteWalletDialog)

	if updatedDialog.active != "cancel" {
		t.Errorf("Expected active to change to 'cancel' with 'h' key, got '%s'", updatedDialog.active)
	}

	// Test 'l' (right in vim)
	lMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'l'}}
	updatedModel, _ = updatedDialog.Update(lMsg)
	updatedDialog, _ = updatedModel.(DeleteWalletDialog)

	if updatedDialog.active != "confirm" {
		t.Errorf("Expected active to change to 'confirm' with 'l' key, got '%s'", updatedDialog.active)
	}
}

func TestDeleteWalletDialog_Update_EnterConfirm(t *testing.T) {
	dialog := NewDeleteWalletDialog("Test", "0x123")
	dialog.active = "confirm"

	enterMsg := tea.KeyMsg{Type: tea.KeyEnter}

	updatedModel, cmd := dialog.Update(enterMsg)

	_, ok := updatedModel.(DeleteWalletDialog)
	if !ok {
		t.Fatal("Expected returned model to be DeleteWalletDialog")
	}

	if cmd == nil {
		t.Error("Expected command to be returned for enter on confirm")
	}

	// Execute the command to test the message
	msg := cmd()
	if _, ok := msg.(ConfirmDeleteMsg); !ok {
		t.Error("Expected ConfirmDeleteMsg")
	}
}

func TestDeleteWalletDialog_Update_EnterCancel(t *testing.T) {
	dialog := NewDeleteWalletDialog("Test", "0x123")
	dialog.active = "cancel"

	enterMsg := tea.KeyMsg{Type: tea.KeyEnter}

	updatedModel, cmd := dialog.Update(enterMsg)

	_, ok := updatedModel.(DeleteWalletDialog)
	if !ok {
		t.Fatal("Expected returned model to be DeleteWalletDialog")
	}

	if cmd == nil {
		t.Error("Expected command to be returned for enter on cancel")
	}

	// Execute the command to test the message
	msg := cmd()
	if _, ok := msg.(CancelDeleteMsg); !ok {
		t.Error("Expected CancelDeleteMsg")
	}
}

func TestDeleteWalletDialog_Update_EscapeKey(t *testing.T) {
	dialog := NewDeleteWalletDialog("Test", "0x123")

	escMsg := tea.KeyMsg{Type: tea.KeyEsc}

	updatedModel, cmd := dialog.Update(escMsg)

	_, ok := updatedModel.(DeleteWalletDialog)
	if !ok {
		t.Fatal("Expected returned model to be DeleteWalletDialog")
	}

	if cmd == nil {
		t.Error("Expected command to be returned for escape key")
	}

	// Execute the command to test the message
	msg := cmd()
	if _, ok := msg.(CancelDeleteMsg); !ok {
		t.Error("Expected CancelDeleteMsg for escape key")
	}
}

func TestDeleteWalletDialog_Update_MouseClickConfirm(t *testing.T) {
	dialog := NewDeleteWalletDialog("Test", "0x123")

	// Initialize bubblezone
	zone.NewGlobal()

	// Create a mouse release message
	mouseMsg := tea.MouseMsg{
		X:      10,
		Y:      5,
		Action: tea.MouseActionRelease,
		Button: tea.MouseButtonLeft,
	}

	// Test mouse click (we can't easily test the zone boundaries without proper setup)
	updatedModel, cmd := dialog.Update(mouseMsg)

	_, ok := updatedModel.(DeleteWalletDialog)
	if !ok {
		t.Fatal("Expected returned model to be DeleteWalletDialog")
	}

	// Should return nil command if no zone is hit
	if cmd != nil {
		t.Error("Expected no command for mouse click outside zones")
	}
}

func TestDeleteWalletDialog_Update_MouseClickWrongButton(t *testing.T) {
	dialog := NewDeleteWalletDialog("Test", "0x123")

	// Test right mouse button (should be ignored)
	mouseMsg := tea.MouseMsg{
		X:      10,
		Y:      5,
		Action: tea.MouseActionRelease,
		Button: tea.MouseButtonRight,
	}

	updatedModel, cmd := dialog.Update(mouseMsg)

	_, ok := updatedModel.(DeleteWalletDialog)
	if !ok {
		t.Fatal("Expected returned model to be DeleteWalletDialog")
	}

	if cmd != nil {
		t.Error("Expected no command for right mouse button")
	}
}

func TestDeleteWalletDialog_Update_MouseClickWrongAction(t *testing.T) {
	dialog := NewDeleteWalletDialog("Test", "0x123")

	// Test mouse press (not release)
	mouseMsg := tea.MouseMsg{
		X:      10,
		Y:      5,
		Action: tea.MouseActionPress,
		Button: tea.MouseButtonLeft,
	}

	updatedModel, cmd := dialog.Update(mouseMsg)

	_, ok := updatedModel.(DeleteWalletDialog)
	if !ok {
		t.Fatal("Expected returned model to be DeleteWalletDialog")
	}

	if cmd != nil {
		t.Error("Expected no command for mouse press")
	}
}

func TestDeleteWalletDialog_View_CancelActive(t *testing.T) {
	dialog := NewDeleteWalletDialog("My Wallet", "0x1234567890123456789012345678901234567890")
	dialog.active = "cancel"

	view := dialog.View()

	if !strings.Contains(view, "Are you sure you want to delete this wallet?") {
		t.Error("Expected view to contain question")
	}

	if !strings.Contains(view, "My Wallet") {
		t.Error("Expected view to contain wallet name")
	}

	if !strings.Contains(view, "0x12345678...1234567890") {
		t.Error("Expected view to contain abbreviated address")
	}

	if !strings.Contains(view, "This action cannot be undone.") {
		t.Error("Expected view to contain warning")
	}

	if !strings.Contains(view, "Yes, Delete") {
		t.Error("Expected view to contain confirm button")
	}

	if !strings.Contains(view, "Cancel") {
		t.Error("Expected view to contain cancel button")
	}
}

func TestDeleteWalletDialog_View_ConfirmActive(t *testing.T) {
	dialog := NewDeleteWalletDialog("My Wallet", "0x1234567890123456789012345678901234567890")
	dialog.active = "confirm"

	view := dialog.View()

	if !strings.Contains(view, "Are you sure you want to delete this wallet?") {
		t.Error("Expected view to contain question")
	}

	if !strings.Contains(view, "My Wallet") {
		t.Error("Expected view to contain wallet name")
	}

	if !strings.Contains(view, "Yes, Delete") {
		t.Error("Expected view to contain confirm button")
	}

	if !strings.Contains(view, "Cancel") {
		t.Error("Expected view to contain cancel button")
	}
}

func TestDeleteWalletDialog_AddressAbbreviation(t *testing.T) {
	// Test with standard Ethereum address
	dialog := NewDeleteWalletDialog("Test", "0x1234567890abcdef1234567890abcdef12345678")

	view := dialog.View()

	// Expected format: first 10 chars + "..." + last 10 chars
	if !strings.Contains(view, "0x12345678...ef12345678") {
		t.Error("Expected view to contain correctly abbreviated address")
	}
}

func TestConfirmDeleteMsg(t *testing.T) {
	msg := ConfirmDeleteMsg{}
	// Just test that the struct can be instantiated
	_ = msg
}

func TestCancelDeleteMsg(t *testing.T) {
	msg := CancelDeleteMsg{}
	// Just test that the struct can be instantiated
	_ = msg
}