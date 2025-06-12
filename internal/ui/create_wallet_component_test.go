package ui

import (
	"fmt"
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

func TestNewCreateWalletComponent(t *testing.T) {
	component := NewCreateWalletComponent()

	if component.id != "create-wallet" {
		t.Errorf("Expected id to be 'create-wallet', got %s", component.id)
	}

	if component.inputFocus != 0 {
		t.Errorf("Expected initial inputFocus to be 0, got %d", component.inputFocus)
	}

	if component.creating {
		t.Error("Expected component to not be in creating state initially")
	}

	if component.err != nil {
		t.Errorf("Expected no initial error, got %v", component.err)
	}

	// Test input configurations
	if component.nameInput.Placeholder != "Enter wallet name..." {
		t.Errorf("Expected nameInput placeholder to be 'Enter wallet name...', got %s", component.nameInput.Placeholder)
	}

	if component.passwordInput.Placeholder != "Enter password..." {
		t.Errorf("Expected passwordInput placeholder to be 'Enter password...', got %s", component.passwordInput.Placeholder)
	}
}

func TestCreateWalletComponent_SetSize(t *testing.T) {
	component := NewCreateWalletComponent()

	component.SetSize(800, 600)

	if component.width != 800 {
		t.Errorf("Expected width to be 800, got %d", component.width)
	}

	if component.height != 600 {
		t.Errorf("Expected height to be 600, got %d", component.height)
	}
}

func TestCreateWalletComponent_SetError(t *testing.T) {
	component := NewCreateWalletComponent()
	component.creating = true

	testError := fmt.Errorf("test error")
	component.SetError(testError)

	if component.err == nil {
		t.Error("Expected error to be set")
	}

	if component.creating {
		t.Error("Expected creating to be false after setting error")
	}
}

func TestCreateWalletComponent_SetCreating(t *testing.T) {
	component := NewCreateWalletComponent()
	component.err = fmt.Errorf("test error")

	component.SetCreating(true)

	if !component.creating {
		t.Error("Expected creating to be true")
	}

	if component.err != nil {
		t.Error("Expected error to be cleared when setting creating to true")
	}

	component.SetCreating(false)

	if component.creating {
		t.Error("Expected creating to be false")
	}
}

func TestCreateWalletComponent_GetMethods(t *testing.T) {
	component := NewCreateWalletComponent()

	// Set test values
	component.nameInput.SetValue("test-wallet")
	component.passwordInput.SetValue("test-password")

	if component.GetWalletName() != "test-wallet" {
		t.Errorf("Expected wallet name to be 'test-wallet', got %s", component.GetWalletName())
	}

	if component.GetPassword() != "test-password" {
		t.Errorf("Expected password to be 'test-password', got %s", component.GetPassword())
	}
}

func TestCreateWalletComponent_Reset(t *testing.T) {
	component := NewCreateWalletComponent()

	// Set some values and state
	component.nameInput.SetValue("test-wallet")
	component.passwordInput.SetValue("test-password")
	component.inputFocus = 1
	component.err = fmt.Errorf("test error")
	component.creating = true

	component.Reset()

	if component.GetWalletName() != "" {
		t.Error("Expected wallet name to be empty after reset")
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

	if component.creating {
		t.Error("Expected creating to be false after reset")
	}
}

func TestCreateWalletComponent_Update_WindowSize(t *testing.T) {
	component := NewCreateWalletComponent()

	msg := tea.WindowSizeMsg{Width: 1024, Height: 768}

	updatedComponent, _ := component.Update(msg)

	if updatedComponent.width != 1024 {
		t.Errorf("Expected width to be 1024, got %d", updatedComponent.width)
	}

	if updatedComponent.height != 768 {
		t.Errorf("Expected height to be 768, got %d", updatedComponent.height)
	}
}

func TestCreateWalletComponent_Update_TabNavigation(t *testing.T) {
	component := NewCreateWalletComponent()

	// Test forward tab navigation
	tabMsg := tea.KeyMsg{Type: tea.KeyTab}

	updatedComponent, _ := component.Update(tabMsg)
	if updatedComponent.inputFocus != 1 {
		t.Errorf("Expected inputFocus to be 1 after tab, got %d", updatedComponent.inputFocus)
	}

	updatedComponent, _ = updatedComponent.Update(tabMsg)
	if updatedComponent.inputFocus != 0 {
		t.Errorf("Expected inputFocus to wrap to 0 after second tab, got %d", updatedComponent.inputFocus)
	}
}

func TestCreateWalletComponent_Update_ShiftTabNavigation(t *testing.T) {
	component := NewCreateWalletComponent()

	// Test reverse tab navigation
	shiftTabMsg := tea.KeyMsg{Type: tea.KeyShiftTab}

	updatedComponent, _ := component.Update(shiftTabMsg)
	if updatedComponent.inputFocus != 1 {
		t.Errorf("Expected inputFocus to wrap to 1 after shift+tab, got %d", updatedComponent.inputFocus)
	}

	updatedComponent, _ = updatedComponent.Update(shiftTabMsg)
	if updatedComponent.inputFocus != 0 {
		t.Errorf("Expected inputFocus to be 0 after second shift+tab, got %d", updatedComponent.inputFocus)
	}
}

func TestCreateWalletComponent_Update_EnterOnPasswordField(t *testing.T) {
	component := NewCreateWalletComponent()

	// Set valid inputs
	component.nameInput.SetValue("test-wallet")
	component.passwordInput.SetValue("test-password")
	component.inputFocus = 1 // Password field

	enterMsg := tea.KeyMsg{Type: tea.KeyEnter}

	updatedComponent, cmd := component.Update(enterMsg)

	if !updatedComponent.creating {
		t.Error("Expected component to be in creating state after valid enter")
	}

	if cmd == nil {
		t.Error("Expected command to be returned after valid enter")
	}
}

func TestCreateWalletComponent_Update_EnterOnNameField(t *testing.T) {
	component := NewCreateWalletComponent()
	component.inputFocus = 0 // Name field

	enterMsg := tea.KeyMsg{Type: tea.KeyEnter}

	updatedComponent, _ := component.Update(enterMsg)

	if updatedComponent.inputFocus != 1 {
		t.Errorf("Expected inputFocus to move to 1 after enter on name field, got %d", updatedComponent.inputFocus)
	}
}

func TestCreateWalletComponent_Update_EscapeKey(t *testing.T) {
	component := NewCreateWalletComponent()

	escMsg := tea.KeyMsg{Type: tea.KeyEsc}

	_, cmd := component.Update(escMsg)

	if cmd == nil {
		t.Error("Expected command to be returned after escape key")
	}
}

func TestCreateWalletComponent_Update_WalletCreatedMsg(t *testing.T) {
	component := NewCreateWalletComponent()

	// Set some initial state
	component.nameInput.SetValue("test-wallet")
	component.passwordInput.SetValue("test-password")
	component.inputFocus = 1
	component.creating = true

	walletCreatedMsg := walletCreatedMsg{}

	updatedComponent, cmd := component.Update(walletCreatedMsg)

	// Should reset the component
	if updatedComponent.GetWalletName() != "" {
		t.Error("Expected wallet name to be reset after walletCreatedMsg")
	}

	if updatedComponent.GetPassword() != "" {
		t.Error("Expected password to be reset after walletCreatedMsg")
	}

	if updatedComponent.inputFocus != 0 {
		t.Error("Expected inputFocus to be reset to 0 after walletCreatedMsg")
	}

	if updatedComponent.creating {
		t.Error("Expected creating to be false after walletCreatedMsg")
	}

	if cmd == nil {
		t.Error("Expected BackToMenuMsg command to be returned after walletCreatedMsg")
	}
}

func TestCreateWalletComponent_Update_ErrorMsg(t *testing.T) {
	component := NewCreateWalletComponent()
	component.creating = true

	errorMsg := errorMsg("Test error message")

	updatedComponent, _ := component.Update(errorMsg)

	if updatedComponent.err == nil {
		t.Error("Expected error to be set after errorMsg")
	}

	if updatedComponent.creating {
		t.Error("Expected creating to be false after errorMsg")
	}
}

func TestCreateWalletComponent_ValidateInputs_EmptyName(t *testing.T) {
	component := NewCreateWalletComponent()

	component.nameInput.SetValue("")
	component.passwordInput.SetValue("test-password")

	isValid := component.validateInputs()

	if isValid {
		t.Error("Expected validation to fail for empty name")
	}

	if component.err == nil {
		t.Error("Expected error to be set for empty name")
	}
}

func TestCreateWalletComponent_ValidateInputs_EmptyPassword(t *testing.T) {
	component := NewCreateWalletComponent()

	component.nameInput.SetValue("test-wallet")
	component.passwordInput.SetValue("")

	isValid := component.validateInputs()

	if isValid {
		t.Error("Expected validation to fail for empty password")
	}

	if component.err == nil {
		t.Error("Expected error to be set for empty password")
	}
}

func TestCreateWalletComponent_ValidateInputs_WhitespaceOnly(t *testing.T) {
	component := NewCreateWalletComponent()

	component.nameInput.SetValue("   ")
	component.passwordInput.SetValue("   ")

	isValid := component.validateInputs()

	if isValid {
		t.Error("Expected validation to fail for whitespace-only inputs")
	}

	if component.err == nil {
		t.Error("Expected error to be set for whitespace-only inputs")
	}
}

func TestCreateWalletComponent_ValidateInputs_Valid(t *testing.T) {
	component := NewCreateWalletComponent()

	component.nameInput.SetValue("test-wallet")
	component.passwordInput.SetValue("test-password")

	isValid := component.validateInputs()

	if !isValid {
		t.Error("Expected validation to pass for valid inputs")
	}

	if component.err != nil {
		t.Errorf("Expected no error for valid inputs, got %v", component.err)
	}
}

func TestCreateWalletComponent_View_Basic(t *testing.T) {
	component := NewCreateWalletComponent()
	component.SetSize(80, 24)

	view := component.View()

	if !strings.Contains(view, "➕ Create New Wallet") {
		t.Error("Expected view to contain create wallet header")
	}

	if !strings.Contains(view, "Wallet Name:") {
		t.Error("Expected view to contain wallet name label")
	}

	if !strings.Contains(view, "Password:") {
		t.Error("Expected view to contain password label")
	}

	if !strings.Contains(view, "Your wallet will be secured with a mnemonic phrase.") {
		t.Error("Expected view to contain info message")
	}

	if !strings.Contains(view, "Tab: Next Field • Enter: Create • Esc: Back") {
		t.Error("Expected view to contain footer instructions")
	}
}

func TestCreateWalletComponent_View_CreatingState(t *testing.T) {
	component := NewCreateWalletComponent()
	component.SetSize(80, 24)
	component.SetCreating(true)

	view := component.View()

	if !strings.Contains(view, "⏳ Creating wallet...") {
		t.Error("Expected view to contain creating message")
	}
}

func TestCreateWalletComponent_View_ErrorState(t *testing.T) {
	component := NewCreateWalletComponent()
	component.SetSize(80, 24)
	component.err = fmt.Errorf("test error")

	view := component.View()

	if !strings.Contains(view, "❌ Error:") {
		t.Error("Expected view to contain error message")
	}
}

func TestCreateWalletComponent_View_FocusedFields(t *testing.T) {
	component := NewCreateWalletComponent()
	component.SetSize(80, 24)
	component.nameInput.SetValue("test-wallet")
	component.passwordInput.SetValue("test-password")

	// Test name field focus
	component.inputFocus = 0
	view := component.View()
	if !strings.Contains(view, "test-wallet") {
		t.Error("Expected view to contain wallet name when name field is focused")
	}

	// Test password field focus
	component.inputFocus = 1
	view = component.View()
	// Password field should show as masked, but we can still check basic structure
	if !strings.Contains(view, "Password:") {
		t.Error("Expected view to contain password label when password field is focused")
	}
}