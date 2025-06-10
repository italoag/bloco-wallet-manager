package ui

import (
	"fmt"
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
)

func TestNewImportPrivateKeyComponent(t *testing.T) {
	component := NewImportPrivateKeyComponent()

	// Test that component is properly initialized
	if component.id != "import-private-key" {
		t.Errorf("Expected id to be 'import-private-key', got %s", component.id)
	}

	if component.form == nil {
		t.Error("Expected form to be initialized")
	}

	// Test initial values
	if component.walletName != "" {
		t.Errorf("Expected empty wallet name, got %s", component.walletName)
	}

	if component.privateKey != "" {
		t.Errorf("Expected empty private key, got %s", component.privateKey)
	}

	if component.password != "" {
		t.Errorf("Expected empty password, got %s", component.password)
	}
}

func TestImportPrivateKeyComponent_SetSize(t *testing.T) {
	component := NewImportPrivateKeyComponent()
	
	component.SetSize(100, 50)
	
	if component.width != 100 {
		t.Errorf("Expected width 100, got %d", component.width)
	}
	
	if component.height != 50 {
		t.Errorf("Expected height 50, got %d", component.height)
	}
}

func TestImportPrivateKeyComponent_SetError(t *testing.T) {
	component := NewImportPrivateKeyComponent()
	testError := fmt.Errorf("test error")
	
	component.SetError(testError)
	
	if component.err == nil {
		t.Error("Expected error to be set")
	}
	
	if component.err.Error() != "test error" {
		t.Errorf("Expected 'test error', got %s", component.err.Error())
	}
	
	if component.importing {
		t.Error("Expected importing to be false when error is set")
	}
}

func TestImportPrivateKeyComponent_SetImporting(t *testing.T) {
	component := NewImportPrivateKeyComponent()
	
	component.SetImporting(true)
	
	if !component.importing {
		t.Error("Expected importing to be true")
	}
	
	if component.err != nil {
		t.Error("Expected error to be nil when importing is set to true")
	}
}

func TestImportPrivateKeyComponent_GetWalletName(t *testing.T) {
	component := NewImportPrivateKeyComponent()
	
	// Set wallet name
	component.walletName = "test wallet"
	
	if component.GetWalletName() != "test wallet" {
		t.Errorf("Expected 'test wallet', got %s", component.GetWalletName())
	}
}

func TestImportPrivateKeyComponent_GetPrivateKey(t *testing.T) {
	component := NewImportPrivateKeyComponent()
	
	// Set private key
	component.privateKey = "abcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890"
	
	if component.GetPrivateKey() != "abcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890" {
		t.Errorf("Expected private key, got %s", component.GetPrivateKey())
	}
}

func TestImportPrivateKeyComponent_GetPassword(t *testing.T) {
	component := NewImportPrivateKeyComponent()
	
	// Set password
	component.password = "password123"
	
	if component.GetPassword() != "password123" {
		t.Errorf("Expected 'password123', got %s", component.GetPassword())
	}
}

func TestImportPrivateKeyComponent_Reset(t *testing.T) {
	component := NewImportPrivateKeyComponent()
	
	// Set some values
	component.walletName = "test"
	component.privateKey = "abcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890"
	component.password = "password"
	component.err = fmt.Errorf("error")
	component.importing = true
	
	component.Reset()
	
	// All values should be reset
	if component.walletName != "" {
		t.Errorf("Expected empty wallet name after reset, got %s", component.walletName)
	}
	
	if component.privateKey != "" {
		t.Errorf("Expected empty private key after reset, got %s", component.privateKey)
	}
	
	if component.password != "" {
		t.Errorf("Expected empty password after reset, got %s", component.password)
	}
	
	if component.err != nil {
		t.Error("Expected nil error after reset")
	}
	
	if component.importing {
		t.Error("Expected importing to be false after reset")
	}
}

func TestImportPrivateKeyComponent_ValidateInputs(t *testing.T) {
	component := NewImportPrivateKeyComponent()
	
	// Test empty wallet name
	if component.validateInputs() {
		t.Error("Expected validation to fail with empty wallet name")
	}
	
	// Set wallet name
	component.walletName = "test wallet"
	
	// Test empty private key
	if component.validateInputs() {
		t.Error("Expected validation to fail with empty private key")
	}
	
	// Set private key
	component.privateKey = "abcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890"
	
	// Test empty password
	if component.validateInputs() {
		t.Error("Expected validation to fail with empty password")
	}
	
	// Set password
	component.password = "password123"
	
	// Now validation should pass
	if !component.validateInputs() {
		t.Error("Expected validation to pass with all fields filled")
	}
}

func TestImportPrivateKeyComponent_ValidateInputs_InvalidPrivateKey(t *testing.T) {
	component := NewImportPrivateKeyComponent()
	
	// Set valid wallet name and password
	component.walletName = "test wallet"
	component.password = "password123"
	
	// Test short private key
	component.privateKey = "abcdef"
	if component.validateInputs() {
		t.Error("Expected validation to fail with short private key")
	}
	
	// Test private key with invalid characters
	component.privateKey = "abcdef1234567890abcdef1234567890abcdef1234567890abcdef12345678gg"
	if component.validateInputs() {
		t.Error("Expected validation to fail with invalid hex characters")
	}
}

func TestImportPrivateKeyComponent_Update(t *testing.T) {
	component := NewImportPrivateKeyComponent()
	
	// Test window size message
	windowMsg := tea.WindowSizeMsg{Width: 120, Height: 60}
	updatedComponent, _ := component.Update(windowMsg)
	
	if updatedComponent.width != 120 {
		t.Errorf("Expected width 120, got %d", updatedComponent.width)
	}
	
	if updatedComponent.height != 60 {
		t.Errorf("Expected height 60, got %d", updatedComponent.height)
	}
}

func TestImportPrivateKeyComponent_EscapeKey(t *testing.T) {
	component := NewImportPrivateKeyComponent()
	
	// Test escape key
	keyMsg := tea.KeyMsg{Type: tea.KeyEsc}
	_, cmd := component.Update(keyMsg)
	
	// Should return BackToMenuMsg
	if cmd == nil {
		t.Error("Expected command from escape key")
	}
	
	msg := cmd()
	if _, ok := msg.(BackToMenuMsg); !ok {
		t.Error("Expected BackToMenuMsg from escape key")
	}
}

func TestImportPrivateKeyComponent_FormCompletion(t *testing.T) {
	component := NewImportPrivateKeyComponent()
	
	// Fill valid data
	component.walletName = "test wallet"
	component.privateKey = "abcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890"
	component.password = "password123"
	
	// Set form state to completed
	component.form.State = huh.StateCompleted
	
	// Update should trigger import
	updatedComponent, cmd := component.Update(tea.KeyMsg{})
	
	if !updatedComponent.importing {
		t.Error("Expected importing to be true after form completion")
	}
	
	if cmd == nil {
		t.Error("Expected command from form completion")
	}
	
	// Execute command to check message type
	msg := cmd()
	if importMsg, ok := msg.(ImportPrivateKeyRequestMsg); ok {
		if importMsg.Name != "test wallet" {
			t.Errorf("Expected name 'test wallet', got %s", importMsg.Name)
		}
		if importMsg.PrivateKey != "abcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890" {
			t.Errorf("Expected private key, got %s", importMsg.PrivateKey)
		}
		if importMsg.Password != "password123" {
			t.Errorf("Expected password 'password123', got %s", importMsg.Password)
		}
	} else {
		t.Error("Expected ImportPrivateKeyRequestMsg from form completion")
	}
}

func TestImportPrivateKeyComponent_FormCompletionValidationFails(t *testing.T) {
	component := NewImportPrivateKeyComponent()
	
	// Don't fill data (validation should fail)
	component.walletName = ""
	component.privateKey = ""
	component.password = ""
	
	// Set form state to completed
	component.form.State = huh.StateCompleted
	
	// Update should not trigger import due to validation failure
	updatedComponent, cmd := component.Update(tea.KeyMsg{})
	
	if updatedComponent.importing {
		t.Error("Expected importing to be false when validation fails")
	}
	
	// Form state should be reset to normal
	if updatedComponent.form.State != huh.StateNormal {
		t.Error("Expected form state to be reset to normal after validation failure")
	}
	
	// Should not return an import command
	if cmd != nil {
		msg := cmd()
		if _, ok := msg.(ImportPrivateKeyRequestMsg); ok {
			t.Error("Expected no ImportPrivateKeyRequestMsg when validation fails")
		}
	}
}

func TestImportPrivateKeyComponent_WalletCreatedMessage(t *testing.T) {
	component := NewImportPrivateKeyComponent()
	
	// Set some state
	component.importing = true
	component.walletName = "test"
	component.privateKey = "abcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890"
	component.password = "password"
	
	// Simulate wallet created message
	updatedComponent, cmd := component.Update(walletCreatedMsg{})
	
	// Component should be reset
	if updatedComponent.walletName != "" {
		t.Error("Expected wallet name to be reset after wallet created")
	}
	
	if updatedComponent.privateKey != "" {
		t.Error("Expected private key to be reset after wallet created")
	}
	
	if updatedComponent.password != "" {
		t.Error("Expected password to be reset after wallet created")
	}
	
	if updatedComponent.importing {
		t.Error("Expected importing to be false after wallet created")
	}
	
	// Should return BackToMenuMsg
	if cmd == nil {
		t.Error("Expected command from wallet created message")
	}
	
	msg := cmd()
	if _, ok := msg.(BackToMenuMsg); !ok {
		t.Error("Expected BackToMenuMsg from wallet created")
	}
}

func TestImportPrivateKeyComponent_ErrorMessage(t *testing.T) {
	component := NewImportPrivateKeyComponent()
	
	// Simulate error message
	updatedComponent, _ := component.Update(errorMsg("test error"))
	
	if updatedComponent.err == nil {
		t.Error("Expected error to be set")
	}
	
	if updatedComponent.err.Error() != "test error" {
		t.Errorf("Expected 'test error', got %s", updatedComponent.err.Error())
	}
}

func TestImportPrivateKeyComponent_View(t *testing.T) {
	component := NewImportPrivateKeyComponent()
	component.SetSize(80, 24)
	
	view := component.View()
	
	// Should not be empty
	if view == "" {
		t.Error("Expected non-empty view")
	}
	
	// Should contain header text
	if !strings.Contains(view, "Import Wallet from Private Key") {
		t.Error("Expected view to contain header text")
	}
}

func TestImportPrivateKeyComponent_ViewWithError(t *testing.T) {
	component := NewImportPrivateKeyComponent()
	component.SetError(fmt.Errorf("test error"))
	
	view := component.View()
	
	// Should contain error message
	if !strings.Contains(view, "test error") {
		t.Error("Expected view to contain error message")
	}
}

func TestImportPrivateKeyComponent_ViewWithImporting(t *testing.T) {
	component := NewImportPrivateKeyComponent()
	component.SetImporting(true)
	
	view := component.View()
	
	// Should contain importing message
	if !strings.Contains(view, "Importing wallet") {
		t.Error("Expected view to contain importing message")
	}
}