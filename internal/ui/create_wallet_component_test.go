package ui

import (
	"fmt"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
)

func TestNewCreateWalletComponent(t *testing.T) {
	component := NewCreateWalletComponent()

	if component.id != "create-wallet" {
		t.Errorf("Expected id to be 'create-wallet', got %s", component.id)
	}

	if component.form == nil {
		t.Error("Expected form to be initialized")
	}

	if component.GetWalletName() != "" {
		t.Errorf("Expected empty wallet name, got %s", component.GetWalletName())
	}

	if component.GetPassword() != "" {
		t.Errorf("Expected empty password, got %s", component.GetPassword())
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

	if component.err != testError {
		t.Error("Expected error to be set")
	}

	if component.creating {
		t.Error("Expected creating to be false after setting error")
	}
}

func TestCreateWalletComponent_SetCreating(t *testing.T) {
	component := NewCreateWalletComponent()

	component.SetCreating(true)

	if !component.creating {
		t.Error("Expected creating to be true")
	}

	if component.err != nil {
		t.Error("Expected error to be nil when creating")
	}

	component.SetCreating(false)

	if component.creating {
		t.Error("Expected creating to be false")
	}
}

func TestCreateWalletComponent_Reset(t *testing.T) {
	component := NewCreateWalletComponent()

	component.err = fmt.Errorf("error")
	component.creating = true

	component.Reset()

	if component.err != nil {
		t.Error("Expected nil error after reset")
	}

	if component.creating {
		t.Error("Expected creating to be false after reset")
	}

	if component.form == nil {
		t.Error("Expected form to be reinitialized after reset")
	}
}

func TestCreateWalletComponent_FormCompletion(t *testing.T) {
	component := NewCreateWalletComponent()

	// Set values directly to the component variables
	component.walletName = "test wallet"
	component.password = "password123"

	// Set form state to completed to simulate completion
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

func TestCreateWalletComponent_EscapeKey(t *testing.T) {
	component := NewCreateWalletComponent()

	// Simulate escape key press
	keyMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("esc")}
	_, cmd := component.Update(keyMsg)

	if cmd == nil {
		t.Error("Expected command from escape key")
	}

	// Execute command to check message type
	msg := cmd()
	if _, ok := msg.(BackToMenuMsg); !ok {
		t.Error("Expected BackToMenuMsg from escape key")
	}
}

func TestCreateWalletComponent_WalletCreatedMessage(t *testing.T) {
	component := NewCreateWalletComponent()
	component.creating = true

	// Simulate wallet created message
	_, cmd := component.Update(walletCreatedMsg{})

	if component.creating {
		t.Error("Expected creating to be false after wallet created")
	}

	if cmd == nil {
		t.Error("Expected command from wallet created message")
	}

	// Execute command to check message type
	msg := cmd()
	if _, ok := msg.(BackToMenuMsg); !ok {
		t.Error("Expected BackToMenuMsg from wallet created")
	}
}

func TestCreateWalletComponent_ErrorMessage(t *testing.T) {
	component := NewCreateWalletComponent()

	testErrorMsg := errorMsg("test error")
	component.Update(testErrorMsg)

	if component.err == nil {
		t.Error("Expected error to be set")
	}

	if component.err.Error() != "test error" {
		t.Errorf("Expected error message 'test error', got %s", component.err.Error())
	}
}

func TestCreateWalletComponent_View(t *testing.T) {
	component := NewCreateWalletComponent()
	component.SetSize(80, 24)

	view := component.View()

	if view == "" {
		t.Error("Expected non-empty view")
	}

	// Check for basic elements
	if !contains(view, "Create New Wallet") {
		t.Error("Expected view to contain title")
	}
}

// Helper function to check if string contains substring
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) &&
		(s[:len(substr)] == substr || s[len(s)-len(substr):] == substr ||
			containsSubstring(s, substr)))
}

func containsSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
