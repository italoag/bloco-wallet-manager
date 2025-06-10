package ui

import (
	"strconv"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	zone "github.com/lrstanley/bubblezone"
)

// MainMenuComponent represents the main menu component
type MainMenuComponent struct {
	id       string
	items    []string
	selected int
	width    int
	height   int
}

// NewMainMenuComponent creates a new main menu component
func NewMainMenuComponent() MainMenuComponent {
	return MainMenuComponent{
		id: "main-menu",
		items: []string{
			"📋 View Wallets",
			"➕ Create New Wallet",
			"📥 Import Wallet from Mnemonic",
			"🔑 Import Wallet from Private Key",
			"⚙️  Settings",
			"❌ Exit",
		},
		selected: 0,
	}
}

// SetSize updates the component size
func (c *MainMenuComponent) SetSize(width, height int) {
	c.width = width
	c.height = height
}

// GetSelected returns the currently selected menu index
func (c *MainMenuComponent) GetSelected() int {
	return c.selected
}

// GetSelectedItem returns the currently selected menu item
func (c *MainMenuComponent) GetSelectedItem() string {
	if c.selected >= 0 && c.selected < len(c.items) {
		return c.items[c.selected]
	}
	return ""
}

// Update handles messages for the main menu component
func (c *MainMenuComponent) Update(msg tea.Msg) (*MainMenuComponent, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		c.width = msg.Width
		c.height = msg.Height

	case tea.MouseMsg:
		if msg.Action != tea.MouseActionRelease || msg.Button != tea.MouseButtonLeft {
			return c, nil
		}

		// Check clicks on menu items
		for i := range c.items {
			if zone.Get(c.id + "-item-" + strconv.Itoa(i)).InBounds(msg) {
				c.selected = i
				return c, func() tea.Msg { return MenuItemSelectedMsg{Index: i, Item: c.items[i]} }
			}
		}

	case tea.KeyMsg:
		switch {
		case key.Matches(msg, key.NewBinding(key.WithKeys("up", "k"))):
			if c.selected > 0 {
				c.selected--
			}
		case key.Matches(msg, key.NewBinding(key.WithKeys("down", "j"))):
			if c.selected < len(c.items)-1 {
				c.selected++
			}
		case key.Matches(msg, key.NewBinding(key.WithKeys("enter"))):
			return c, func() tea.Msg { return MenuItemSelectedMsg{Index: c.selected, Item: c.items[c.selected]} }
		case key.Matches(msg, key.NewBinding(key.WithKeys("1"))):
			c.selected = 0
			return c, func() tea.Msg { return MenuItemSelectedMsg{Index: 0, Item: c.items[0]} }
		case key.Matches(msg, key.NewBinding(key.WithKeys("2"))):
			c.selected = 1
			return c, func() tea.Msg { return MenuItemSelectedMsg{Index: 1, Item: c.items[1]} }
		case key.Matches(msg, key.NewBinding(key.WithKeys("3"))):
			c.selected = 2
			return c, func() tea.Msg { return MenuItemSelectedMsg{Index: 2, Item: c.items[2]} }
		case key.Matches(msg, key.NewBinding(key.WithKeys("4"))):
			c.selected = 3
			return c, func() tea.Msg { return MenuItemSelectedMsg{Index: 3, Item: c.items[3]} }
		case key.Matches(msg, key.NewBinding(key.WithKeys("5"))):
			c.selected = 4
			return c, func() tea.Msg { return MenuItemSelectedMsg{Index: 4, Item: c.items[4]} }
		case key.Matches(msg, key.NewBinding(key.WithKeys("q", "esc"))):
			c.selected = 5 // Exit
			return c, func() tea.Msg { return MenuItemSelectedMsg{Index: 5, Item: c.items[5]} }
		}
	}

	return c, nil
}

// View renders the main menu component
func (c *MainMenuComponent) View() string {
	var b strings.Builder

	b.WriteString(HeaderStyle.Render("🏦 BlockoWallet - Main Menu"))
	b.WriteString("\n\n")

	for i, item := range c.items {
		var itemText string
		if i == c.selected {
			itemText = MenuSelectedStyle.Render("► " + item)
		} else {
			itemText = ItemStyle.Render("  " + item)
		}

		// Mark zone for mouse interaction
		b.WriteString(zone.Mark(c.id+"-item-"+strconv.Itoa(i), itemText))
		if i < len(c.items)-1 {
			b.WriteString("\n")
		}
	}

	return b.String()
}

// MenuItemSelectedMsg is sent when a menu item is selected
type MenuItemSelectedMsg struct {
	Index int
	Item  string
}
