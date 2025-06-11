package ui

import (
	"fmt"
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
)

func TestNewImportMnemonicComponent(t *testing.T) {
	component := NewImportMnemonicComponent()

	// Test that component is properly initialized
	if component.id != "import-mnemonic" {
		t.Errorf("Expected id to be 'import-mnemonic', got %s", component.id)
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

func TestImportMnemonicComponent_SetSize(t *testing.T) {
	component := NewImportMnemonicComponent()

	component.SetSize(100, 50)

	if component.width != 100 {
		t.Errorf("Expected width 100, got %d", component.width)
	}

	if component.height != 50 {
		t.Errorf("Expected height 50, got %d", component.height)
	}
}

func TestImportMnemonicComponent_SetError(t *testing.T) {
	component := NewImportMnemonicComponent()
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

func TestImportMnemonicComponent_SetImporting(t *testing.T) {
	component := NewImportMnemonicComponent()

	component.SetImporting(true)

	if !component.importing {
		t.Error("Expected importing to be true")
	}

	if component.err != nil {
		t.Error("Expected error to be nil when importing is set to true")
	}
}

func TestImportMnemonicComponent_GetMnemonic(t *testing.T) {
	component := NewImportMnemonicComponent()

	// Set individual words
	component.word1 = "abandon"
	component.word2 = "ability"
	component.word3 = "able"
	component.word4 = "about"
	component.word5 = "above"
	component.word6 = "absent"
	component.word7 = "absorb"
	component.word8 = "abstract"
	component.word9 = "absurd"
	component.word10 = "abuse"
	component.word11 = "access"
	component.word12 = "accident"

	expected := "abandon ability able about above absent absorb abstract absurd abuse access accident"
	mnemonic := component.GetMnemonic()

	if mnemonic != expected {
		t.Errorf("Expected '%s', got '%s'", expected, mnemonic)
	}
}

func TestImportMnemonicComponent_Reset(t *testing.T) {
	component := NewImportMnemonicComponent()

	// Set some values
	component.walletName = "test"
	component.password = "password"
	component.word1 = "word1"
	component.err = fmt.Errorf("error")
	component.importing = true

	component.Reset()

	// All values should be reset
	if component.walletName != "" {
		t.Errorf("Expected empty wallet name after reset, got %s", component.walletName)
	}

	if component.password != "" {
		t.Errorf("Expected empty password after reset, got %s", component.password)
	}

	if component.word1 != "" {
		t.Errorf("Expected empty word1 after reset, got %s", component.word1)
	}

	if component.err != nil {
		t.Error("Expected nil error after reset")
	}

	if component.importing {
		t.Error("Expected importing to be false after reset")
	}

	// Test that wizard step is reset to first step
	if component.currentStep != 1 {
		t.Errorf("Expected currentStep to be 1 after reset, got %d", component.currentStep)
	}
}

func TestImportMnemonicComponent_Update(t *testing.T) {
	component := NewImportMnemonicComponent()

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

func TestImportMnemonicComponent_EscapeKey(t *testing.T) {
	component := NewImportMnemonicComponent()

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

func TestImportMnemonicComponent_FormCompletion(t *testing.T) {
	component := NewImportMnemonicComponent()

	// Test that the component starts at step 1
	if component.currentStep != 1 {
		t.Errorf("Expected currentStep to be 1, got %d", component.currentStep)
	}

	// Fill valid data for step 1 (wallet name and password)
	component.walletName = "test wallet"
	component.password = "password123"

	// Complete step 1
	component.form.State = huh.StateCompleted

	// Update should advance to step 2
	updatedComponent, _ := component.Update(tea.KeyMsg{})

	if updatedComponent.currentStep != 2 {
		t.Errorf("Expected currentStep to be 2 after completing step 1, got %d", updatedComponent.currentStep)
	}

	// It's okay to have a command when advancing to step 2 (form initialization)

	// Now simulate completing step 2 with mnemonic words
	// Set the mnemonic words in component variables (for fallback)
	updatedComponent.word1 = "abandon"
	updatedComponent.word2 = "ability"
	updatedComponent.word3 = "able"
	updatedComponent.word4 = "about"
	updatedComponent.word5 = "above"
	updatedComponent.word6 = "absent"
	updatedComponent.word7 = "absorb"
	updatedComponent.word8 = "abstract"
	updatedComponent.word9 = "absurd"
	updatedComponent.word10 = "abuse"
	updatedComponent.word11 = "access"
	updatedComponent.word12 = "accident"

	// Complete step 2
	updatedComponent.form.State = huh.StateCompleted

	// Test that the component is not already importing (prevents duplicate processing)
	if updatedComponent.importing {
		t.Error("Component should not be importing initially")
	}

	// Update should trigger import
	finalComponent, cmd2 := updatedComponent.Update(tea.KeyMsg{})

	if !finalComponent.importing {
		t.Error("Expected importing to be true after completing step 2")
	}

	if cmd2 == nil {
		t.Error("Expected command from form completion")
	}

	// Execute command to check message type
	msg := cmd2()
	if importMsg, ok := msg.(ImportMnemonicRequestMsg); ok {
		// Test that the message was created with some values
		// In test context, values may come from component variables (fallback)
		if importMsg.Name == "" && importMsg.Password == "" && importMsg.Mnemonic == "" {
			t.Error("Expected non-empty values in ImportMnemonicRequestMsg")
		}
	} else {
		t.Error("Expected ImportMnemonicRequestMsg from form completion")
	}

	// Test that subsequent updates don't trigger import again (preventing duplicates)
	_, cmd3 := finalComponent.Update(tea.KeyMsg{})
	if cmd3 != nil {
		msg3 := cmd3()
		if _, ok := msg3.(ImportMnemonicRequestMsg); ok {
			t.Error("Should not create duplicate ImportMnemonicRequestMsg when already importing")
		}
	}
}

func TestImportMnemonicComponent_WizardSteps(t *testing.T) {
	component := NewImportMnemonicComponent()

	// Test initial step
	if component.currentStep != 1 {
		t.Errorf("Expected currentStep to be 1 initially, got %d", component.currentStep)
	}

	// Test that step 1 form is initialized
	if component.form == nil {
		t.Error("Expected form to be initialized")
	}

	// Simulate completing step 1
	component.walletName = "test"
	component.password = "test123"
	component.form.State = huh.StateCompleted

	// Update should advance to step 2
	updatedComponent, _ := component.Update(tea.KeyMsg{})

	if updatedComponent.currentStep != 2 {
		t.Errorf("Expected currentStep to be 2 after completing step 1, got %d", updatedComponent.currentStep)
	}

	// Test that stored values from step 1 are preserved
	if updatedComponent.walletName != "test" {
		t.Errorf("Expected walletName to be preserved, got %s", updatedComponent.walletName)
	}

	if updatedComponent.password != "test123" {
		t.Errorf("Expected password to be preserved, got %s", updatedComponent.password)
	}
}

func TestImportMnemonicComponent_View(t *testing.T) {
	component := NewImportMnemonicComponent()
	component.SetSize(80, 24)

	view := component.View()

	// Should not be empty
	if view == "" {
		t.Error("Expected non-empty view")
	}

	// Should contain header text
	if !strings.Contains(view, "Import Wallet from Mnemonic") {
		t.Error("Expected view to contain header text")
	}
}

func TestImportMnemonicComponent_ViewWithError(t *testing.T) {
	component := NewImportMnemonicComponent()
	component.SetError(fmt.Errorf("test error"))

	view := component.View()

	// Should contain error message
	if !strings.Contains(view, "test error") {
		t.Error("Expected view to contain error message")
	}
}

func TestImportMnemonicComponent_ViewWithImporting(t *testing.T) {
	component := NewImportMnemonicComponent()
	component.SetImporting(true)

	view := component.View()

	// Should contain importing message
	if !strings.Contains(view, "Importing wallet") {
		t.Error("Expected view to contain importing message")
	}
}
