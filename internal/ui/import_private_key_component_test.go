package ui

import (
	"fmt"
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

func TestNewImportPrivateKeyComponent(t *testing.T) {
	component := NewImportPrivateKeyComponent()

	if component.id != "import-private-key" {
		t.Errorf("Expected id to be 'import-private-key', got %s", component.id)
	}

	if component.inputFocus != 0 {
		t.Errorf("Expected initial inputFocus to be 0, got %d", component.inputFocus)
	}

	if component.importing {
		t.Error("Expected component to not be in importing state initially")
	}

	if component.err != nil {
		t.Errorf("Expected no initial error, got %v", component.err)
	}

	// Test input configurations
	if component.nameInput.Placeholder != "Enter wallet name..." {
		t.Errorf("Expected nameInput placeholder to be 'Enter wallet name...', got %s", component.nameInput.Placeholder)
	}

	if component.privateKeyInput.Placeholder != "Enter private key (without 0x prefix)..." {
		t.Errorf("Expected privateKeyInput placeholder to be 'Enter private key (without 0x prefix)...', got %s", component.privateKeyInput.Placeholder)
	}

	if component.passwordInput.Placeholder != "Enter password..." {
		t.Errorf("Expected passwordInput placeholder to be 'Enter password...', got %s", component.passwordInput.Placeholder)
	}
}

func TestImportPrivateKeyComponent_SetSize(t *testing.T) {
	component := NewImportPrivateKeyComponent()

	component.SetSize(800, 600)

	if component.width != 800 {
		t.Errorf("Expected width to be 800, got %d", component.width)
	}

	if component.height != 600 {
		t.Errorf("Expected height to be 600, got %d", component.height)
	}
}

func TestImportPrivateKeyComponent_SetError(t *testing.T) {
	component := NewImportPrivateKeyComponent()
	component.importing = true

	testError := fmt.Errorf("test error")
	component.SetError(testError)

	if component.err == nil {
		t.Error("Expected error to be set")
	}

	if component.importing {
		t.Error("Expected importing to be false after setting error")
	}
}

func TestImportPrivateKeyComponent_SetImporting(t *testing.T) {
	component := NewImportPrivateKeyComponent()
	component.err = fmt.Errorf("test error")

	component.SetImporting(true)

	if !component.importing {
		t.Error("Expected importing to be true")
	}

	if component.err != nil {
		t.Error("Expected error to be cleared when setting importing to true")
	}

	component.SetImporting(false)

	if component.importing {
		t.Error("Expected importing to be false")
	}
}

func TestImportPrivateKeyComponent_GetMethods(t *testing.T) {
	component := NewImportPrivateKeyComponent()

	// Set test values
	component.nameInput.SetValue("test-wallet")
	component.privateKeyInput.SetValue("abcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890")
	component.passwordInput.SetValue("test-password")

	if component.GetWalletName() != "test-wallet" {
		t.Errorf("Expected wallet name to be 'test-wallet', got %s", component.GetWalletName())
	}

	if component.GetPrivateKey() != "abcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890" {
		t.Errorf("Expected private key to be set correctly, got %s", component.GetPrivateKey())
	}

	if component.GetPassword() != "test-password" {
		t.Errorf("Expected password to be 'test-password', got %s", component.GetPassword())
	}
}

func TestImportPrivateKeyComponent_Reset(t *testing.T) {
	component := NewImportPrivateKeyComponent()

	// Set some values and state
	component.nameInput.SetValue("test-wallet")
	component.privateKeyInput.SetValue("test-key")
	component.passwordInput.SetValue("test-password")
	component.inputFocus = 2
	component.err = fmt.Errorf("test error")
	component.importing = true

	component.Reset()

	if component.GetWalletName() != "" {
		t.Error("Expected wallet name to be empty after reset")
	}

	if component.GetPrivateKey() != "" {
		t.Error("Expected private key to be empty after reset")
	}

	if component.GetPassword() != "" {
		t.Error("Expected password to be empty after reset")
	}

	if component.inputFocus != 0 {
		t.Errorf("Expected inputFocus to be 0 after reset, got %d", component.inputFocus)
	}

	if component.err != nil {
		t.Error("Expected error to be nil after reset")
	}

	if component.importing {
		t.Error("Expected importing to be false after reset")
	}
}

func TestImportPrivateKeyComponent_Update_WindowSize(t *testing.T) {
	component := NewImportPrivateKeyComponent()

	msg := tea.WindowSizeMsg{Width: 1024, Height: 768}

	updatedComponent, _ := component.Update(msg)

	if updatedComponent.width != 1024 {
		t.Errorf("Expected width to be 1024, got %d", updatedComponent.width)
	}

	if updatedComponent.height != 768 {
		t.Errorf("Expected height to be 768, got %d", updatedComponent.height)
	}
}

func TestImportPrivateKeyComponent_Update_TabNavigation(t *testing.T) {
	component := NewImportPrivateKeyComponent()

	// Test forward tab navigation
	tabMsg := tea.KeyMsg{Type: tea.KeyTab}

	updatedComponent, _ := component.Update(tabMsg)
	if updatedComponent.inputFocus != 1 {
		t.Errorf("Expected inputFocus to be 1 after tab, got %d", updatedComponent.inputFocus)
	}

	updatedComponent, _ = updatedComponent.Update(tabMsg)
	if updatedComponent.inputFocus != 2 {
		t.Errorf("Expected inputFocus to be 2 after second tab, got %d", updatedComponent.inputFocus)
	}

	updatedComponent, _ = updatedComponent.Update(tabMsg)
	if updatedComponent.inputFocus != 0 {
		t.Errorf("Expected inputFocus to wrap to 0 after third tab, got %d", updatedComponent.inputFocus)
	}
}

func TestImportPrivateKeyComponent_Update_ShiftTabNavigation(t *testing.T) {
	component := NewImportPrivateKeyComponent()

	// Test reverse tab navigation
	shiftTabMsg := tea.KeyMsg{Type: tea.KeyShiftTab}

	updatedComponent, _ := component.Update(shiftTabMsg)
	if updatedComponent.inputFocus != 2 {
		t.Errorf("Expected inputFocus to wrap to 2 after shift+tab, got %d", updatedComponent.inputFocus)
	}

	updatedComponent, _ = updatedComponent.Update(shiftTabMsg)
	if updatedComponent.inputFocus != 1 {
		t.Errorf("Expected inputFocus to be 1 after second shift+tab, got %d", updatedComponent.inputFocus)
	}

	updatedComponent, _ = updatedComponent.Update(shiftTabMsg)
	if updatedComponent.inputFocus != 0 {
		t.Errorf("Expected inputFocus to be 0 after third shift+tab, got %d", updatedComponent.inputFocus)
	}
}

func TestImportPrivateKeyComponent_Update_EnterOnPasswordField(t *testing.T) {
	component := NewImportPrivateKeyComponent()

	// Set valid inputs
	component.nameInput.SetValue("test-wallet")
	component.privateKeyInput.SetValue("abcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890")
	component.passwordInput.SetValue("test-password")
	component.inputFocus = 2 // Password field

	enterMsg := tea.KeyMsg{Type: tea.KeyEnter}

	updatedComponent, cmd := component.Update(enterMsg)

	if !updatedComponent.importing {
		t.Error("Expected component to be in importing state after valid enter")
	}

	if cmd == nil {
		t.Error("Expected command to be returned after valid enter")
	}
}

func TestImportPrivateKeyComponent_Update_EnterOnOtherFields(t *testing.T) {
	component := NewImportPrivateKeyComponent()
	component.inputFocus = 0 // Name field

	enterMsg := tea.KeyMsg{Type: tea.KeyEnter}

	updatedComponent, _ := component.Update(enterMsg)

	if updatedComponent.inputFocus != 1 {
		t.Errorf("Expected inputFocus to move to 1 after enter on name field, got %d", updatedComponent.inputFocus)
	}
}

func TestImportPrivateKeyComponent_Update_EscapeKey(t *testing.T) {
	component := NewImportPrivateKeyComponent()

	escMsg := tea.KeyMsg{Type: tea.KeyEsc}

	_, cmd := component.Update(escMsg)

	if cmd == nil {
		t.Error("Expected command to be returned after escape key")
	}
}

func TestImportPrivateKeyComponent_Update_WalletCreatedMsg(t *testing.T) {
	component := NewImportPrivateKeyComponent()

	// Set some initial state
	component.nameInput.SetValue("test-wallet")
	component.privateKeyInput.SetValue("test-key")
	component.passwordInput.SetValue("test-password")
	component.inputFocus = 2
	component.importing = true

	walletCreatedMsg := walletCreatedMsg{}

	updatedComponent, cmd := component.Update(walletCreatedMsg)

	// Should reset the component
	if updatedComponent.GetWalletName() != "" {
		t.Error("Expected wallet name to be reset after walletCreatedMsg")
	}

	if updatedComponent.GetPrivateKey() != "" {
		t.Error("Expected private key to be reset after walletCreatedMsg")
	}

	if updatedComponent.GetPassword() != "" {
		t.Error("Expected password to be reset after walletCreatedMsg")
	}

	if updatedComponent.inputFocus != 0 {
		t.Error("Expected inputFocus to be reset to 0 after walletCreatedMsg")
	}

	if updatedComponent.importing {
		t.Error("Expected importing to be false after walletCreatedMsg")
	}

	if cmd == nil {
		t.Error("Expected BackToMenuMsg command to be returned after walletCreatedMsg")
	}
}

func TestImportPrivateKeyComponent_Update_ErrorMsg(t *testing.T) {
	component := NewImportPrivateKeyComponent()
	component.importing = true

	errorMsg := errorMsg("Test error message")

	updatedComponent, _ := component.Update(errorMsg)

	if updatedComponent.err == nil {
		t.Error("Expected error to be set after errorMsg")
	}

	if updatedComponent.importing {
		t.Error("Expected importing to be false after errorMsg")
	}
}

func TestImportPrivateKeyComponent_ValidateInputs_EmptyName(t *testing.T) {
	component := NewImportPrivateKeyComponent()

	component.nameInput.SetValue("")
	component.privateKeyInput.SetValue("abcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890")
	component.passwordInput.SetValue("test-password")

	isValid := component.validateInputs()

	if isValid {
		t.Error("Expected validation to fail for empty name")
	}

	if component.err == nil {
		t.Error("Expected error to be set for empty name")
	}
}

func TestImportPrivateKeyComponent_ValidateInputs_EmptyPrivateKey(t *testing.T) {
	component := NewImportPrivateKeyComponent()

	component.nameInput.SetValue("test-wallet")
	component.privateKeyInput.SetValue("")
	component.passwordInput.SetValue("test-password")

	isValid := component.validateInputs()

	if isValid {
		t.Error("Expected validation to fail for empty private key")
	}

	if component.err == nil {
		t.Error("Expected error to be set for empty private key")
	}
}

func TestImportPrivateKeyComponent_ValidateInputs_EmptyPassword(t *testing.T) {
	component := NewImportPrivateKeyComponent()

	component.nameInput.SetValue("test-wallet")
	component.privateKeyInput.SetValue("abcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890")
	component.passwordInput.SetValue("")

	isValid := component.validateInputs()

	if isValid {
		t.Error("Expected validation to fail for empty password")
	}

	if component.err == nil {
		t.Error("Expected error to be set for empty password")
	}
}

func TestImportPrivateKeyComponent_ValidateInputs_InvalidPrivateKeyLength(t *testing.T) {
	component := NewImportPrivateKeyComponent()

	component.nameInput.SetValue("test-wallet")
	component.privateKeyInput.SetValue("abcdef") // Too short
	component.passwordInput.SetValue("test-password")

	isValid := component.validateInputs()

	if isValid {
		t.Error("Expected validation to fail for invalid private key length")
	}

	if component.err == nil {
		t.Error("Expected error to be set for invalid private key length")
	}
}

func TestImportPrivateKeyComponent_ValidateInputs_InvalidPrivateKeyChars(t *testing.T) {
	component := NewImportPrivateKeyComponent()

	component.nameInput.SetValue("test-wallet")
	component.privateKeyInput.SetValue("gggggg1234567890abcdef1234567890abcdef1234567890abcdef1234567890") // Invalid hex
	component.passwordInput.SetValue("test-password")

	isValid := component.validateInputs()

	if isValid {
		t.Error("Expected validation to fail for invalid private key characters")
	}

	if component.err == nil {
		t.Error("Expected error to be set for invalid private key characters")
	}
}

func TestImportPrivateKeyComponent_ValidateInputs_WithPrefix(t *testing.T) {
	component := NewImportPrivateKeyComponent()

	component.nameInput.SetValue("test-wallet")
	component.privateKeyInput.SetValue("0xabcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890")
	component.passwordInput.SetValue("test-password")

	isValid := component.validateInputs()

	if !isValid {
		t.Error("Expected validation to pass for valid private key with 0x prefix")
	}

	if component.err != nil {
		t.Errorf("Expected no error for valid input, got %v", component.err)
	}
}

func TestImportPrivateKeyComponent_ValidateInputs_Valid(t *testing.T) {
	component := NewImportPrivateKeyComponent()

	component.nameInput.SetValue("test-wallet")
	component.privateKeyInput.SetValue("abcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890")
	component.passwordInput.SetValue("test-password")

	isValid := component.validateInputs()

	if !isValid {
		t.Error("Expected validation to pass for valid inputs")
	}

	if component.err != nil {
		t.Errorf("Expected no error for valid inputs, got %v", component.err)
	}
}

func TestImportPrivateKeyComponent_View_Basic(t *testing.T) {
	component := NewImportPrivateKeyComponent()
	component.SetSize(80, 24)

	view := component.View()

	if !strings.Contains(view, "🔑 Import Wallet from Private Key") {
		t.Error("Expected view to contain import header")
	}

	if !strings.Contains(view, "Wallet Name:") {
		t.Error("Expected view to contain wallet name label")
	}

	if !strings.Contains(view, "Private Key (32 bytes / 64 hex characters):") {
		t.Error("Expected view to contain private key label")
	}

	if !strings.Contains(view, "Password:") {
		t.Error("Expected view to contain password label")
	}

	if !strings.Contains(view, "Never share your private key with anyone!") {
		t.Error("Expected view to contain security warning")
	}

	if !strings.Contains(view, "Tab: Next Field • Enter: Import • Esc: Back") {
		t.Error("Expected view to contain footer instructions")
	}
}

func TestImportPrivateKeyComponent_View_ImportingState(t *testing.T) {
	component := NewImportPrivateKeyComponent()
	component.SetSize(80, 24)
	component.SetImporting(true)

	view := component.View()

	if !strings.Contains(view, "⏳ Importing wallet...") {
		t.Error("Expected view to contain importing message")
	}
}

func TestImportPrivateKeyComponent_View_ErrorState(t *testing.T) {
	component := NewImportPrivateKeyComponent()
	component.SetSize(80, 24)
	component.err = fmt.Errorf("test error")

	view := component.View()

	if !strings.Contains(view, "❌ Error:") {
		t.Error("Expected view to contain error message")
	}
}