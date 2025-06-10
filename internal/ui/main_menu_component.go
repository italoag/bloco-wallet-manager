package ui

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
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
	title       string
	description string
	index       int
}

func (i menuItem) Title() string       { return i.title }
func (i menuItem) Description() string { return i.description }
func (i menuItem) FilterValue() string { return i.title }

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
			Foreground(lipgloss.Color("#FFFDF5")).
			Background(lipgloss.Color("#25A065")).
			Padding(0, 1)

	statusMessageStyle = lipgloss.NewStyle().
				Foreground(lipgloss.AdaptiveColor{Light: "#04B575", Dark: "#04B575"}).
				Render
)

// NewMainMenuComponent creates a new main menu component
func NewMainMenuComponent() MainMenuComponent {
	items := []list.Item{
		menuItem{title: "📋 View Wallets", description: "View and manage your wallets", index: 0},
		menuItem{title: "➕ Create New Wallet", description: "Create a new wallet with a generated mnemonic", index: 1},
		menuItem{title: "📥 Import Wallet from Mnemonic", description: "Import an existing wallet using mnemonic phrase", index: 2},
		menuItem{title: "🔑 Import Wallet from Private Key", description: "Import an existing wallet using private key", index: 3},
		menuItem{title: "⚙️  Settings", description: "Configure networks and preferences", index: 4},
		menuItem{title: "❌ Exit", description: "Exit the application", index: 5},
	}

	keys := newMenuKeyMap()
	delegate := newMenuDelegate(keys)
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

// newMenuDelegate creates a delegate for the menu list
func newMenuDelegate(keys *menuKeyMap) list.DefaultDelegate {
	d := list.NewDefaultDelegate()

	d.UpdateFunc = func(msg tea.Msg, m *list.Model) tea.Cmd {
		var item menuItem
		var ok bool

		if i := m.SelectedItem(); i != nil {
			item, ok = i.(menuItem)
			if !ok {
				return nil
			}
		}

		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch {
			case key.Matches(msg, keys.choose):
				return func() tea.Msg {
					return MenuItemSelectedMsg{Index: item.index, Item: item.title}
				}
			}
		}

		return nil
	}

	return d
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
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		c.width = msg.Width
		c.height = msg.Height
		c.list.SetSize(msg.Width, msg.Height)

	case tea.KeyMsg:
		// Handle number shortcuts
		switch msg.String() {
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
		// Adicionando tratamento específico para teclas de navegação
		case "up", "k":
			var cmd tea.Cmd
			c.list, cmd = c.list.Update(msg)
			// Garantir um comando mesmo que cmd seja nil
			if cmd == nil {
				cmd = func() tea.Msg { return nil }
			}
			return c, cmd
		case "down", "j":
			var cmd tea.Cmd
			c.list, cmd = c.list.Update(msg)
			// Garantir um comando mesmo que cmd seja nil
			if cmd == nil {
				cmd = func() tea.Msg { return nil }
			}
			return c, cmd
		}
	}

	// Update the list
	c.list, cmd = c.list.Update(msg)
	// Garantir um comando mesmo que cmd seja nil
	if cmd == nil {
		cmd = func() tea.Msg { return nil }
	}
	return c, cmd
}

// View renders the main menu component
func (c *MainMenuComponent) View() string {
	return appStyle.Render(c.list.View())
}

// MenuItemSelectedMsg is sent when a menu item is selected
type MenuItemSelectedMsg struct {
	Index int
	Item  string
}
