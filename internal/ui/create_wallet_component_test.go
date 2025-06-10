package ui

import (
	"fmt"
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
)

func TestNewCreateWalletComponent(t *testing.T) {
	component := NewCreateWalletComponent()

	// Test that component is properly initialized
	if component.id != "create-wallet" {
		t.Errorf("Expected id to be 'create-wallet', got %s", component.id)
	}

	if component.form == nil {
		t.Error("Expected form to be initialized")
	}

	// Test initial values
	if component.walletName != "" {
		t.Errorf("Expected empty wallet name, got %s", component.walletName)
	}

	if component.password != "" {
		t.Errorf("Expected empty password, got %s", component.password)
	}
}

func TestCreateWalletComponent_SetSize(t *testing.T) {
	component := NewCreateWalletComponent()
	
	component.SetSize(100, 50)
	
	if component.width != 100 {
		t.Errorf("Expected width 100, got %d", component.width)
	}
	
	if component.height != 50 {
		t.Errorf("Expected height 50, got %d", component.height)
	}
}

func TestCreateWalletComponent_SetError(t *testing.T) {
	component := NewCreateWalletComponent()
	testError := fmt.Errorf("test error")
	
	component.SetError(testError)
	
	if component.err == nil {
		t.Error("Expected error to be set")
	}
	
	if component.err.Error() != "test error" {
		t.Errorf("Expected 'test error', got %s", component.err.Error())
	}
	
	if component.creating {
		t.Error("Expected creating to be false when error is set")
	}
}

func TestCreateWalletComponent_SetCreating(t *testing.T) {
	component := NewCreateWalletComponent()
	
	component.SetCreating(true)
	
	if !component.creating {
		t.Error("Expected creating to be true")
	}
	
	if component.err != nil {
		t.Error("Expected error to be nil when creating is set to true")
	}
}

func TestCreateWalletComponent_GetWalletName(t *testing.T) {
	component := NewCreateWalletComponent()
	
	// Set wallet name
	component.walletName = "test wallet"
	
	if component.GetWalletName() != "test wallet" {
		t.Errorf("Expected 'test wallet', got %s", component.GetWalletName())
	}
}

func TestCreateWalletComponent_GetPassword(t *testing.T) {
	component := NewCreateWalletComponent()
	
	// Set password
	component.password = "password123"
	
	if component.GetPassword() != "password123" {
		t.Errorf("Expected 'password123', got %s", component.GetPassword())
	}
}

func TestCreateWalletComponent_Reset(t *testing.T) {
	component := NewCreateWalletComponent()
	
	// Set some values
	component.walletName = "test"
	component.password = "password"
	component.err = fmt.Errorf("error")
	component.creating = true
	
	component.Reset()
	
	// All values should be reset
	if component.walletName != "" {
		t.Errorf("Expected empty wallet name after reset, got %s", component.walletName)
	}
	
	if component.password != "" {
		t.Errorf("Expected empty password after reset, got %s", component.password)
	}
	
	if component.err != nil {
		t.Error("Expected nil error after reset")
	}
	
	if component.creating {
		t.Error("Expected creating to be false after reset")
	}
}

func TestCreateWalletComponent_ValidateInputs(t *testing.T) {
	component := NewCreateWalletComponent()
	
	// Test empty wallet name
	if component.validateInputs() {
		t.Error("Expected validation to fail with empty wallet name")
	}
	
	// Set wallet name
	component.walletName = "test wallet"
	
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

func TestCreateWalletComponent_Update(t *testing.T) {
	component := NewCreateWalletComponent()
	
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

func TestCreateWalletComponent_EscapeKey(t *testing.T) {
	component := NewCreateWalletComponent()
	
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

func TestCreateWalletComponent_FormCompletion(t *testing.T) {
	component := NewCreateWalletComponent()
	
	// Fill valid data
	component.walletName = "test wallet"
	component.password = "password123"
	
	// Set form state to completed
	component.form.State = huh.StateCompleted
	
	// Update should trigger wallet creation
	updatedComponent, cmd := component.Update(tea.KeyMsg{})
	
	if !updatedComponent.creating {
		t.Error("Expected creating to be true after form completion")
	}
	
	if cmd == nil {
		t.Error("Expected command from form completion")
	}
	
	// Execute command to check message type
	msg := cmd()
	if createMsg, ok := msg.(CreateWalletRequestMsg); ok {
		if createMsg.Name != "test wallet" {
			t.Errorf("Expected name 'test wallet', got %s", createMsg.Name)
		}
		if createMsg.Password != "password123" {
			t.Errorf("Expected password 'password123', got %s", createMsg.Password)
		}
	} else {
		t.Error("Expected CreateWalletRequestMsg from form completion")
	}
}

func TestCreateWalletComponent_FormCompletionValidationFails(t *testing.T) {
	component := NewCreateWalletComponent()
	
	// Don't fill data (validation should fail)
	component.walletName = ""
	component.password = ""
	
	// Set form state to completed
	component.form.State = huh.StateCompleted
	
	// Update should not trigger wallet creation due to validation failure
	updatedComponent, cmd := component.Update(tea.KeyMsg{})
	
	if updatedComponent.creating {
		t.Error("Expected creating to be false when validation fails")
	}
	
	// Form state should be reset to normal
	if updatedComponent.form.State != huh.StateNormal {
		t.Error("Expected form state to be reset to normal after validation failure")
	}
	
	// Should not return a creation command
	if cmd != nil {
		msg := cmd()
		if _, ok := msg.(CreateWalletRequestMsg); ok {
			t.Error("Expected no CreateWalletRequestMsg when validation fails")
		}
	}
}

func TestCreateWalletComponent_WalletCreatedMessage(t *testing.T) {
	component := NewCreateWalletComponent()
	
	// Set some state
	component.creating = true
	component.walletName = "test"
	component.password = "password"
	
	// Simulate wallet created message
	updatedComponent, cmd := component.Update(walletCreatedMsg{})
	
	// Component should be reset
	if updatedComponent.walletName != "" {
		t.Error("Expected wallet name to be reset after wallet created")
	}
	
	if updatedComponent.password != "" {
		t.Error("Expected password to be reset after wallet created")
	}
	
	if updatedComponent.creating {
		t.Error("Expected creating to be false after wallet created")
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

func TestCreateWalletComponent_ErrorMessage(t *testing.T) {
	component := NewCreateWalletComponent()
	
	// Simulate error message
	updatedComponent, _ := component.Update(errorMsg("test error"))
	
	if updatedComponent.err == nil {
		t.Error("Expected error to be set")
	}
	
	if updatedComponent.err.Error() != "test error" {
		t.Errorf("Expected 'test error', got %s", updatedComponent.err.Error())
	}
}

func TestCreateWalletComponent_View(t *testing.T) {
	component := NewCreateWalletComponent()
	component.SetSize(80, 24)
	
	view := component.View()
	
	// Should not be empty
	if view == "" {
		t.Error("Expected non-empty view")
	}
	
	// Should contain header text
	if !strings.Contains(view, "Create New Wallet") {
		t.Error("Expected view to contain header text")
	}
}

func TestCreateWalletComponent_ViewWithError(t *testing.T) {
	component := NewCreateWalletComponent()
	component.SetError(fmt.Errorf("test error"))
	
	view := component.View()
	
	// Should contain error message
	if !strings.Contains(view, "test error") {
		t.Error("Expected view to contain error message")
	}
}

func TestCreateWalletComponent_ViewWithCreating(t *testing.T) {
	component := NewCreateWalletComponent()
	component.SetCreating(true)
	
	view := component.View()
	
	// Should contain creating message
	if !strings.Contains(view, "Creating wallet") {
		t.Error("Expected view to contain creating message")
	}
}