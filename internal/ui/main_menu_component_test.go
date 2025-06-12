package ui

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	zone "github.com/lrstanley/bubblezone"
)

func TestNewMainMenuComponent(t *testing.T) {
	component := NewMainMenuComponent()

	// Test that component is properly initialized
	if component.id != "main-menu" {
		t.Errorf("Expected id to be 'main-menu', got %s", component.id)
	}

	if component.keys == nil {
		t.Error("Expected keys to be initialized")
	}

	// Test that list has proper title
	if component.list.Title != "🏦 BlockoWallet - Main Menu" {
		t.Errorf("Expected title '🏦 BlockoWallet - Main Menu', got %s", component.list.Title)
	}

	// Test that list has 6 items (0-5)
	if len(component.list.Items()) != 6 {
		t.Errorf("Expected 6 menu items, got %d", len(component.list.Items()))
	}
}

func TestMainMenuComponent_SetSize(t *testing.T) {
	component := NewMainMenuComponent()

	component.SetSize(100, 50)

	if component.width != 100 {
		t.Errorf("Expected width 100, got %d", component.width)
	}

	if component.height != 50 {
		t.Errorf("Expected height 50, got %d", component.height)
	}
}

func TestMainMenuComponent_GetSelected(t *testing.T) {
	component := NewMainMenuComponent()

	// Default selection should be 0
	if component.GetSelected() != 0 {
		t.Errorf("Expected default selection 0, got %d", component.GetSelected())
	}
}

func TestMainMenuComponent_GetSelectedItem(t *testing.T) {
	component := NewMainMenuComponent()

	// Should return first menu item title
	selected := component.GetSelectedItem()
	if selected != "📋 View Wallets" {
		t.Errorf("Expected '📋 View Wallets', got %s", selected)
	}
}

func TestMainMenuComponent_Update(t *testing.T) {
	component := NewMainMenuComponent()

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

func TestMainMenuComponent_KeyNavigation(t *testing.T) {
	component := NewMainMenuComponent()

	// Test navigation keys that should be handled by the list
	keyMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("j")} // 'j' is mapped to "down" in vim style
	updatedComponent, _ := component.Update(keyMsg)

	// Component should be updated (cmd can be nil for navigation)
	if updatedComponent == nil {
		t.Error("Expected updated component")
	}

	// Test that we can still navigate
	if updatedComponent.GetSelected() < 0 || updatedComponent.GetSelected() >= 6 {
		t.Error("Selection should be within valid range")
	}
}

func TestMainMenuComponent_NumberShortcuts(t *testing.T) {
	component := NewMainMenuComponent()

	// Test number shortcut "1"
	keyMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'1'}}
	_, cmd := component.Update(keyMsg)

	// Should return a menu selection command
	if cmd == nil {
		t.Error("Expected menu selection command")
	}

	// Execute the command to get the message
	msg := cmd()
	if menuMsg, ok := msg.(MenuItemSelectedMsg); ok {
		if menuMsg.Index != 0 {
			t.Errorf("Expected index 0, got %d", menuMsg.Index)
		}
		if menuMsg.Item != "📋 View Wallets" {
			t.Errorf("Expected '📋 View Wallets', got %s", menuMsg.Item)
		}
	} else {
		t.Error("Expected MenuItemSelectedMsg")
	}
}

func TestMainMenuComponent_View(t *testing.T) {
	// Initialize bubblezone for the test
	zone.NewGlobal()

	component := NewMainMenuComponent()
	component.SetSize(80, 24)

	view := component.View()

	// Should not be empty
	if view == "" {
		t.Error("Expected non-empty view")
	}

	// Should contain the title
	if !contains(view, "BlockoWallet") {
		t.Error("Expected view to contain 'BlockoWallet'")
	}
}

// Helper function to check if string contains substring
func contains(s, substr string) bool {
	return len(s) >= len(substr) && s[:len(substr)] == substr ||
		(len(s) > len(substr) && contains(s[1:], substr))
}
