package ui

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	zone "github.com/lrstanley/bubblezone"
)

// MainMenuComponent represents the main menu component
type MainMenuComponent struct {
	id     string
	list   list.Model
	width  int
	height int
	keys   *menuKeyMap
}

// menuItem represents a menu item
type menuItem struct {
	id          string
	title       string
	description string
	index       int
}

func (i menuItem) Title() string       { return zone.Mark(i.id, i.title) }
func (i menuItem) Description() string { return i.description }
func (i menuItem) FilterValue() string { return zone.Mark(i.id, i.title) }

// menuKeyMap defines key bindings for the menu
type menuKeyMap struct {
	choose key.Binding
}

func (k menuKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.choose}
}

func (k menuKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{{k.choose}}
}

func newMenuKeyMap() *menuKeyMap {
	return &menuKeyMap{
		choose: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter", "select"),
		),
	}
}

var (
	appStyle = lipgloss.NewStyle().Padding(1, 2)

	titleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#874BFD")).
			Background(lipgloss.Color("#FFFFFF")).
			Padding(0, 1)

	statusMessageStyle = lipgloss.NewStyle().
				Foreground(lipgloss.AdaptiveColor{Light: "#04B575", Dark: "#04B575"}).
				Render
)

// NewMainMenuComponent creates a new main menu component
func NewMainMenuComponent() MainMenuComponent {
	items := []list.Item{
		menuItem{id: "menu_view_wallets", title: "📋 View Wallets", description: "View and manage your wallets", index: 0},
		menuItem{id: "menu_create_wallet", title: "➕ Create New Wallet", description: "Create a new wallet with a generated mnemonic", index: 1},
		menuItem{id: "menu_import_mnemonic", title: "📥 Import Wallet from Mnemonic", description: "Import an existing wallet using mnemonic phrase", index: 2},
		menuItem{id: "menu_import_private_key", title: "🔑 Import Wallet from Private Key", description: "Import an existing wallet using private key", index: 3},
		menuItem{id: "menu_settings", title: "⚙️  Settings", description: "Configure networks and preferences", index: 4},
		menuItem{id: "menu_exit", title: "❌ Exit", description: "Exit the application", index: 5},
	}

	keys := newMenuKeyMap()
	// Use default delegate instead of custom delegate to avoid conflicts
	delegate := list.NewDefaultDelegate()
	menuList := list.New(items, delegate, 0, 0)
	menuList.Title = "🏦 BlockoWallet - Main Menu"
	menuList.Styles.Title = titleStyle
	menuList.SetShowStatusBar(false)
	menuList.SetFilteringEnabled(false)

	return MainMenuComponent{
		id:   "main-menu",
		list: menuList,
		keys: keys,
	}
}

// SetSize updates the component size
func (c *MainMenuComponent) SetSize(width, height int) {
	c.width = width
	c.height = height
	c.list.SetSize(width, height)
}

// GetSelected returns the currently selected menu index
func (c *MainMenuComponent) GetSelected() int {
	if item, ok := c.list.SelectedItem().(menuItem); ok {
		return item.index
	}
	return 0
}

// GetSelectedItem returns the currently selected menu item
func (c *MainMenuComponent) GetSelectedItem() string {
	if item, ok := c.list.SelectedItem().(menuItem); ok {
		return item.title
	}
	return ""
}

// Update handles messages for the main menu component
func (c *MainMenuComponent) Update(msg tea.Msg) (*MainMenuComponent, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		c.width = msg.Width
		c.height = msg.Height
		c.list.SetSize(msg.Width, msg.Height)

	case tea.MouseMsg:
		if msg.Button == tea.MouseButtonWheelUp {
			c.list.CursorUp()
			return c, nil
		}

		if msg.Button == tea.MouseButtonWheelDown {
			c.list.CursorDown()
			return c, nil
		}

		if msg.Action == tea.MouseActionRelease && msg.Button == tea.MouseButtonLeft {
			for i, listItem := range c.list.VisibleItems() {
				v, _ := listItem.(menuItem)
				// Check each item to see if it's in bounds.
				if zone.Get(v.id).InBounds(msg) {
					// If so, select it in the list.
					c.list.Select(i)
					// Trigger selection action
					return c, func() tea.Msg { return MenuItemSelectedMsg{Index: v.index, Item: v.title} }
				}
			}
		}

		return c, nil

	case tea.KeyMsg:
		// Handle number shortcuts
		switch msg.String() {
		case "enter":
			if item, ok := c.list.SelectedItem().(menuItem); ok {
				return c, func() tea.Msg { return MenuItemSelectedMsg{Index: item.index, Item: item.title} }
			}
		case "1":
			c.list.Select(0)
			if item, ok := c.list.SelectedItem().(menuItem); ok {
				return c, func() tea.Msg { return MenuItemSelectedMsg{Index: item.index, Item: item.title} }
			}
		case "2":
			c.list.Select(1)
			if item, ok := c.list.SelectedItem().(menuItem); ok {
				return c, func() tea.Msg { return MenuItemSelectedMsg{Index: item.index, Item: item.title} }
			}
		case "3":
			c.list.Select(2)
			if item, ok := c.list.SelectedItem().(menuItem); ok {
				return c, func() tea.Msg { return MenuItemSelectedMsg{Index: item.index, Item: item.title} }
			}
		case "4":
			c.list.Select(3)
			if item, ok := c.list.SelectedItem().(menuItem); ok {
				return c, func() tea.Msg { return MenuItemSelectedMsg{Index: item.index, Item: item.title} }
			}
		case "5":
			c.list.Select(4)
			if item, ok := c.list.SelectedItem().(menuItem); ok {
				return c, func() tea.Msg { return MenuItemSelectedMsg{Index: item.index, Item: item.title} }
			}
		case "q", "esc":
			c.list.Select(5) // Exit
			if item, ok := c.list.SelectedItem().(menuItem); ok {
				return c, func() tea.Msg { return MenuItemSelectedMsg{Index: item.index, Item: item.title} }
			}
		}
	}

	var cmd tea.Cmd
	c.list, cmd = c.list.Update(msg)
	return c, cmd
}

// View renders the main menu component
func (c *MainMenuComponent) View() string {
	// Use zone.Scan to wrap the list view for mouse support
	return zone.Scan(appStyle.Render(c.list.View()))
}

// MenuItemSelectedMsg is sent when a menu item is selected
type MenuItemSelectedMsg struct {
	Index int
	Item  string
}
