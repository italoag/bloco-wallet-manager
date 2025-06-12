package ui

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	zone "github.com/lrstanley/bubblezone"
)

// SettingsMenuComponent represents the settings menu component
type SettingsMenuComponent struct {
	id     string
	list   list.Model
	width  int
	height int
	keys   *settingsKeyMap
}

// settingsItem represents a settings menu item
type settingsItem struct {
	id          string
	title       string
	description string
	index       int
}

func (i settingsItem) Title() string       { return zone.Mark(i.id, i.title) }
func (i settingsItem) Description() string { return i.description }
func (i settingsItem) FilterValue() string { return zone.Mark(i.id, i.title) }

// settingsKeyMap defines key bindings for the settings menu
type settingsKeyMap struct {
	choose key.Binding
}

func (k settingsKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.choose}
}

func (k settingsKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{{k.choose}}
}

func newSettingsKeyMap() *settingsKeyMap {
	return &settingsKeyMap{
		choose: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter", "select"),
		),
	}
}

// NewSettingsMenuComponent creates a new settings menu component
func NewSettingsMenuComponent() SettingsMenuComponent {
	items := []list.Item{
		settingsItem{id: "settings_network", title: "🌐 Network Configuration", description: "Configure networks and RPC endpoints", index: 0},
		settingsItem{id: "settings_language", title: "🌍 Language", description: "Change interface language", index: 1},
		settingsItem{id: "settings_back", title: "🔙 Back to Main Menu", description: "Return to the main menu", index: 2},
	}

	keys := newSettingsKeyMap()
	// Use default delegate instead of custom delegate to avoid conflicts
	delegate := list.NewDefaultDelegate()
	delegate.Styles.NormalTitle = delegate.Styles.NormalTitle.
		Foreground(menuItemForeground)
	delegate.Styles.NormalDesc = delegate.Styles.NormalDesc.
		Foreground(menuItemForeground)
	delegate.Styles.SelectedTitle = delegate.Styles.SelectedTitle.
		Foreground(menuItemForeground)
	delegate.Styles.SelectedDesc = delegate.Styles.SelectedDesc.
		Foreground(menuItemForeground)
	settingsList := list.New(items, delegate, 0, 0)
	settingsList.Title = "⚙️  Settings"
	settingsList.Styles.Title = titleStyle
	settingsList.SetShowStatusBar(false)
	settingsList.SetFilteringEnabled(false)

	return SettingsMenuComponent{
		id:   "settings-menu",
		list: settingsList,
		keys: keys,
	}
}

// SetSize updates the component size
func (c *SettingsMenuComponent) SetSize(width, height int) {
	c.width = width
	c.height = height
	c.list.SetSize(width, height)
}

// GetSelected returns the currently selected settings index
func (c *SettingsMenuComponent) GetSelected() int {
	if item, ok := c.list.SelectedItem().(settingsItem); ok {
		return item.index
	}
	return 0
}

// GetSelectedItem returns the currently selected settings item
func (c *SettingsMenuComponent) GetSelectedItem() string {
	if item, ok := c.list.SelectedItem().(settingsItem); ok {
		return item.title
	}
	return ""
}

// Update handles messages for the settings menu component
func (c *SettingsMenuComponent) Update(msg tea.Msg) (*SettingsMenuComponent, tea.Cmd) {
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
				v, _ := listItem.(settingsItem)
				// Check each item to see if it's in bounds.
				if zone.Get(v.id).InBounds(msg) {
					// If so, select it in the list.
					c.list.Select(i)
					// Trigger selection action
					return c, func() tea.Msg { return SettingsItemSelectedMsg{Index: v.index, Item: v.title} }
				}
			}
		}

		return c, nil

	case tea.KeyMsg:
		// Handle number shortcuts
		switch msg.String() {
		case "enter":
			if item, ok := c.list.SelectedItem().(settingsItem); ok {
				return c, func() tea.Msg { return SettingsItemSelectedMsg{Index: item.index, Item: item.title} }
			}
		case "1":
			c.list.Select(0)
			if item, ok := c.list.SelectedItem().(settingsItem); ok {
				return c, func() tea.Msg { return SettingsItemSelectedMsg{Index: item.index, Item: item.title} }
			}
		case "2":
			c.list.Select(1)
			if item, ok := c.list.SelectedItem().(settingsItem); ok {
				return c, func() tea.Msg { return SettingsItemSelectedMsg{Index: item.index, Item: item.title} }
			}
		case "3":
			c.list.Select(2)
			if item, ok := c.list.SelectedItem().(settingsItem); ok {
				return c, func() tea.Msg { return SettingsItemSelectedMsg{Index: item.index, Item: item.title} }
			}
		case "q", "esc":
			c.list.Select(2) // Back to Main Menu
			if item, ok := c.list.SelectedItem().(settingsItem); ok {
				return c, func() tea.Msg { return SettingsItemSelectedMsg{Index: item.index, Item: item.title} }
			}
		}
	}

	var cmd tea.Cmd
	c.list, cmd = c.list.Update(msg)
	return c, cmd
}

// View renders the settings menu component
func (c *SettingsMenuComponent) View() string {
	// Use zone.Scan to wrap the list view for mouse support
	return zone.Scan(appStyle.Render(c.list.View()))
}

// SettingsItemSelectedMsg is sent when a settings item is selected
type SettingsItemSelectedMsg struct {
	Index int
	Item  string
}
