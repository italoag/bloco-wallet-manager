package ui

import (
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	zone "github.com/lrstanley/bubblezone"
)

func TestNewMainMenuComponent(t *testing.T) {
	component := NewMainMenuComponent()

	if component.id != "main-menu" {
		t.Errorf("Expected id to be 'main-menu', got %s", component.id)
	}

	if component.selected != 0 {
		t.Errorf("Expected initial selected to be 0, got %d", component.selected)
	}

	expectedItemCount := 6
	if len(component.items) != expectedItemCount {
		t.Errorf("Expected %d menu items, got %d", expectedItemCount, len(component.items))
	}

	// Test specific menu items
	expectedItems := []string{
		"📋 View Wallets",
		"➕ Create New Wallet",
		"📥 Import Wallet from Mnemonic",
		"🔑 Import Wallet from Private Key",
		"⚙️  Settings",
		"❌ Exit",
	}

	for i, expected := range expectedItems {
		if component.items[i] != expected {
			t.Errorf("Expected item %d to be '%s', got '%s'", i, expected, component.items[i])
		}
	}
}

func TestMainMenuComponent_SetSize(t *testing.T) {
	component := NewMainMenuComponent()

	component.SetSize(800, 600)

	if component.width != 800 {
		t.Errorf("Expected width to be 800, got %d", component.width)
	}

	if component.height != 600 {
		t.Errorf("Expected height to be 600, got %d", component.height)
	}
}

func TestMainMenuComponent_GetSelected(t *testing.T) {
	component := NewMainMenuComponent()

	if component.GetSelected() != 0 {
		t.Errorf("Expected initial selected to be 0, got %d", component.GetSelected())
	}

	component.selected = 3
	if component.GetSelected() != 3 {
		t.Errorf("Expected selected to be 3, got %d", component.GetSelected())
	}
}

func TestMainMenuComponent_GetSelectedItem(t *testing.T) {
	component := NewMainMenuComponent()

	// Test default selection
	expected := "📋 View Wallets"
	if component.GetSelectedItem() != expected {
		t.Errorf("Expected selected item to be '%s', got '%s'", expected, component.GetSelectedItem())
	}

	// Test different selection
	component.selected = 2
	expected = "📥 Import Wallet from Mnemonic"
	if component.GetSelectedItem() != expected {
		t.Errorf("Expected selected item to be '%s', got '%s'", expected, component.GetSelectedItem())
	}

	// Test out of bounds (negative)
	component.selected = -1
	if component.GetSelectedItem() != "" {
		t.Errorf("Expected empty string for out of bounds selection, got '%s'", component.GetSelectedItem())
	}

	// Test out of bounds (positive)
	component.selected = 10
	if component.GetSelectedItem() != "" {
		t.Errorf("Expected empty string for out of bounds selection, got '%s'", component.GetSelectedItem())
	}
}

func TestMainMenuComponent_Update_WindowSize(t *testing.T) {
	component := NewMainMenuComponent()

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

func TestMainMenuComponent_Update_UpKey(t *testing.T) {
	component := NewMainMenuComponent()
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

func TestMainMenuComponent_Update_UpKeyAtTop(t *testing.T) {
	component := NewMainMenuComponent()
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

func TestMainMenuComponent_Update_DownKey(t *testing.T) {
	component := NewMainMenuComponent()
	component.selected = 2

	downMsg := tea.KeyMsg{Type: tea.KeyDown}

	updatedComponent, cmd := component.Update(downMsg)

	if updatedComponent.selected != 3 {
		t.Errorf("Expected selected to be 3 after down key, got %d", updatedComponent.selected)
	}

	if cmd != nil {
		t.Error("Expected no command for down key")
	}
}

func TestMainMenuComponent_Update_DownKeyAtBottom(t *testing.T) {
	component := NewMainMenuComponent()
	component.selected = 5 // Last item

	downMsg := tea.KeyMsg{Type: tea.KeyDown}

	updatedComponent, cmd := component.Update(downMsg)

	if updatedComponent.selected != 5 {
		t.Errorf("Expected selected to stay 5 at bottom, got %d", updatedComponent.selected)
	}

	if cmd != nil {
		t.Error("Expected no command for down key at bottom")
	}
}

func TestMainMenuComponent_Update_VimKeys(t *testing.T) {
	component := NewMainMenuComponent()
	component.selected = 2

	// Test 'k' (up in vim)
	kMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'k'}}
	updatedComponent, _ := component.Update(kMsg)

	if updatedComponent.selected != 1 {
		t.Errorf("Expected selected to be 1 after 'k' key, got %d", updatedComponent.selected)
	}

	// Test 'j' (down in vim)
	jMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}}
	updatedComponent, _ = updatedComponent.Update(jMsg)

	if updatedComponent.selected != 2 {
		t.Errorf("Expected selected to be 2 after 'j' key, got %d", updatedComponent.selected)
	}
}

func TestMainMenuComponent_Update_EnterKey(t *testing.T) {
	component := NewMainMenuComponent()
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
	if menuMsg, ok := msg.(MenuItemSelectedMsg); ok {
		if menuMsg.Index != 1 {
			t.Errorf("Expected message index to be 1, got %d", menuMsg.Index)
		}
		if menuMsg.Item != "➕ Create New Wallet" {
			t.Errorf("Expected message item to be 'Create New Wallet', got '%s'", menuMsg.Item)
		}
	} else {
		t.Error("Expected MenuItemSelectedMsg")
	}
}

func TestMainMenuComponent_Update_NumberKeys(t *testing.T) {
	component := NewMainMenuComponent()

	testCases := []struct {
		key           rune
		expectedIndex int
		expectedItem  string
	}{
		{'1', 0, "📋 View Wallets"},
		{'2', 1, "➕ Create New Wallet"},
		{'3', 2, "📥 Import Wallet from Mnemonic"},
		{'4', 3, "🔑 Import Wallet from Private Key"},
		{'5', 4, "⚙️  Settings"},
	}

	for _, tc := range testCases {
		keyMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{tc.key}}

		updatedComponent, cmd := component.Update(keyMsg)

		if updatedComponent.selected != tc.expectedIndex {
			t.Errorf("Expected selected to be %d for key '%c', got %d", tc.expectedIndex, tc.key, updatedComponent.selected)
		}

		if cmd == nil {
			t.Errorf("Expected command to be returned for key '%c'", tc.key)
			continue
		}

		// Execute the command to test the message
		msg := cmd()
		if menuMsg, ok := msg.(MenuItemSelectedMsg); ok {
			if menuMsg.Index != tc.expectedIndex {
				t.Errorf("Expected message index to be %d for key '%c', got %d", tc.expectedIndex, tc.key, menuMsg.Index)
			}
			if menuMsg.Item != tc.expectedItem {
				t.Errorf("Expected message item to be '%s' for key '%c', got '%s'", tc.expectedItem, tc.key, menuMsg.Item)
			}
		} else {
			t.Errorf("Expected MenuItemSelectedMsg for key '%c'", tc.key)
		}
	}
}

func TestMainMenuComponent_Update_QuitKeys(t *testing.T) {
	component := NewMainMenuComponent()

	testKeys := []tea.KeyMsg{
		{Type: tea.KeyRunes, Runes: []rune{'q'}},
		{Type: tea.KeyEsc},
	}

	for _, keyMsg := range testKeys {
		updatedComponent, cmd := component.Update(keyMsg)

		if updatedComponent.selected != 5 {
			t.Errorf("Expected selected to be 5 (exit) for quit key, got %d", updatedComponent.selected)
		}

		if cmd == nil {
			t.Error("Expected command to be returned for quit key")
			continue
		}

		// Execute the command to test the message
		msg := cmd()
		if menuMsg, ok := msg.(MenuItemSelectedMsg); ok {
			if menuMsg.Index != 5 {
				t.Errorf("Expected message index to be 5 for quit key, got %d", menuMsg.Index)
			}
			if menuMsg.Item != "❌ Exit" {
				t.Errorf("Expected message item to be 'Exit' for quit key, got '%s'", menuMsg.Item)
			}
		} else {
			t.Error("Expected MenuItemSelectedMsg for quit key")
		}
	}
}

func TestMainMenuComponent_Update_MouseClick(t *testing.T) {
	component := NewMainMenuComponent()

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

func TestMainMenuComponent_Update_MouseClickWrongButton(t *testing.T) {
	component := NewMainMenuComponent()

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

func TestMainMenuComponent_Update_MouseClickWrongAction(t *testing.T) {
	component := NewMainMenuComponent()

	// Test mouse press (not release)
	mouseMsg := tea.MouseMsg{
		X:      10,
		Y:      5,
		Action: tea.MouseActionPress,
		Button: tea.MouseButtonLeft,
	}

	updatedComponent, cmd := component.Update(mouseMsg)

	if updatedComponent.selected != component.selected {
		t.Error("Expected selection to remain unchanged for mouse press")
	}

	if cmd != nil {
		t.Error("Expected no command for mouse press")
	}
}

func TestMainMenuComponent_View(t *testing.T) {
	component := NewMainMenuComponent()
	component.SetSize(80, 24)

	view := component.View()

	// Test header
	if !strings.Contains(view, "🏦 BlockoWallet - Main Menu") {
		t.Error("Expected view to contain main menu header")
	}

	// Test all menu items are present
	expectedItems := []string{
		"📋 View Wallets",
		"➕ Create New Wallet",
		"📥 Import Wallet from Mnemonic",
		"🔑 Import Wallet from Private Key",
		"⚙️  Settings",
		"❌ Exit",
	}

	for _, item := range expectedItems {
		if !strings.Contains(view, item) {
			t.Errorf("Expected view to contain menu item '%s'", item)
		}
	}

	// Test selected item indicator
	if !strings.Contains(view, "► 📋 View Wallets") {
		t.Error("Expected first item to be selected with indicator")
	}
}

func TestMainMenuComponent_View_DifferentSelection(t *testing.T) {
	component := NewMainMenuComponent()
	component.SetSize(80, 24)
	component.selected = 2

	view := component.View()

	// Test that the selected item has the indicator
	if !strings.Contains(view, "► 📥 Import Wallet from Mnemonic") {
		t.Error("Expected third item to be selected with indicator")
	}

	// Test that non-selected items don't have the indicator for first item
	if strings.Contains(view, "► 📋 View Wallets") {
		t.Error("Expected first item to not have indicator when not selected")
	}
}

func TestMainMenuComponent_Update_UpKeyTooHigh(t *testing.T) {
	component := NewMainMenuComponent()
	component.selected = 10 // Set to invalid high value

	upMsg := tea.KeyMsg{Type: tea.KeyUp}

	updatedComponent, _ := component.Update(upMsg)

	// Should decrease regardless of being out of bounds
	if updatedComponent.selected != 9 {
		t.Errorf("Expected selected to decrease to 9, got %d", updatedComponent.selected)
	}
}

func TestMainMenuComponent_Update_DownKeyTooLow(t *testing.T) {
	component := NewMainMenuComponent()
	component.selected = -1 // Set to invalid low value

	downMsg := tea.KeyMsg{Type: tea.KeyDown}

	updatedComponent, _ := component.Update(downMsg)

	// Should increase regardless of being out of bounds
	if updatedComponent.selected != 0 {
		t.Errorf("Expected selected to increase to 0, got %d", updatedComponent.selected)
	}
}

func TestMainMenuComponent_Update_UnknownKey(t *testing.T) {
	component := NewMainMenuComponent()
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

func TestMainMenuComponent_Update_OtherMessageType(t *testing.T) {
	component := NewMainMenuComponent()
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

func TestMenuItemSelectedMsg(t *testing.T) {
	msg := MenuItemSelectedMsg{
		Index: 2,
		Item:  "Test Item",
	}

	if msg.Index != 2 {
		t.Errorf("Expected index to be 2, got %d", msg.Index)
	}

	if msg.Item != "Test Item" {
		t.Errorf("Expected item to be 'Test Item', got '%s'", msg.Item)
	}
}