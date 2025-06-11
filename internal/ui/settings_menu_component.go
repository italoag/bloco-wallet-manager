package ui

import (
	"fmt"
	"strings"

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
	title       string
	description string
	index       int
}

func (i settingsItem) Title() string       { return i.title }
func (i settingsItem) Description() string { return i.description }
func (i settingsItem) FilterValue() string { return i.title }

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
		settingsItem{title: "🌐 Network Configuration", description: "Configure networks and RPC endpoints", index: 0},
		settingsItem{title: "🌍 Language", description: "Change interface language", index: 1},
		settingsItem{title: "🔙 Back to Main Menu", description: "Return to the main menu", index: 2},
	}

	keys := newSettingsKeyMap()
	delegate := newSettingsDelegate(keys)
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

// newSettingsDelegate creates a delegate for the settings list
func newSettingsDelegate(keys *settingsKeyMap) list.DefaultDelegate {
	d := list.NewDefaultDelegate()

	d.UpdateFunc = func(msg tea.Msg, m *list.Model) tea.Cmd {
		var item settingsItem
		var ok bool

		if i := m.SelectedItem(); i != nil {
			item, ok = i.(settingsItem)
			if !ok {
				return nil
			}
		}

		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch {
			case key.Matches(msg, keys.choose):
				return func() tea.Msg {
					return SettingsItemSelectedMsg{Index: item.index, Item: item.title}
				}
			}
		}

		return nil
	}

	return d
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
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		c.width = msg.Width
		c.height = msg.Height
		c.list.SetSize(msg.Width, msg.Height)

	case tea.MouseMsg:
		// Check if any settings item was clicked
		for i := 0; i < len(c.list.Items()); i++ {
			itemZoneID := fmt.Sprintf("settings-item-%d", i)
			if zone.Get(itemZoneID).InBounds(msg) {
				// Select and activate the clicked item
				c.list.Select(i)
				if item, ok := c.list.SelectedItem().(settingsItem); ok {
					return c, func() tea.Msg { return SettingsItemSelectedMsg{Index: item.index, Item: item.title} }
				}
			}
		}

	case tea.KeyMsg:
		// Handle number shortcuts
		switch msg.String() {
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

	// Update the list
	c.list, cmd = c.list.Update(msg)
	return c, cmd
}

// View renders the settings menu component
func (c *SettingsMenuComponent) View() string {
	// Render the list first
	listView := c.list.View()

	// Apply zone marking to each settings item for mouse support
	lines := strings.Split(listView, "\n")
	var markedLines []string

	itemIndex := 0
	for _, line := range lines {
		// Check if this line contains a settings item (has content and isn't just formatting)
		if strings.TrimSpace(line) != "" &&
			!strings.Contains(line, "Settings") &&
			!strings.Contains(line, "Help") &&
			(strings.Contains(line, "►") || strings.Contains(line, "•") || strings.Contains(line, "🌐") || strings.Contains(line, "🗣️") || strings.Contains(line, "🔙")) {

			// Mark this line as clickable
			zoneID := fmt.Sprintf("settings-item-%d", itemIndex)
			markedLine := zone.Mark(zoneID, line)
			markedLines = append(markedLines, markedLine)
			itemIndex++
		} else {
			markedLines = append(markedLines, line)
		}
	}

	return appStyle.Render(strings.Join(markedLines, "\n"))
}

// SettingsItemSelectedMsg is sent when a settings item is selected
type SettingsItemSelectedMsg struct {
	Index int
	Item  string
}
