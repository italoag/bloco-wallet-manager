package ui

import (
	"blocowallet/pkg/config"
	"fmt"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	zone "github.com/lrstanley/bubblezone"
)

// NetworkListComponent represents the network list component
type NetworkListComponent struct {
	id     string
	list   list.Model
	width  int
	height int
	keys   *networkKeyMap
	config *config.Config
}

// networkItem represents a network item
type networkItem struct {
	id          string
	title       string
	description string
	key         string
	network     config.Network
}

func (i networkItem) Title() string       { return zone.Mark(i.id, i.title) }
func (i networkItem) Description() string { return i.description }
func (i networkItem) FilterValue() string { return zone.Mark(i.id, i.title) }

// networkKeyMap defines key bindings for the network menu
type networkKeyMap struct {
	choose     key.Binding
	toggle     key.Binding
	edit       key.Binding
	deleteNet  key.Binding
	addNetwork key.Binding
}

func (k networkKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.choose, k.toggle, k.edit, k.addNetwork}
}

func (k networkKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{{k.choose, k.toggle, k.edit, k.deleteNet, k.addNetwork}}
}

func newNetworkKeyMap() *networkKeyMap {
	return &networkKeyMap{
		choose: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter", "details"),
		),
		toggle: key.NewBinding(
			key.WithKeys("a"),
			key.WithHelp("a", "toggle active"),
		),
		edit: key.NewBinding(
			key.WithKeys("e"),
			key.WithHelp("e", "edit RPC"),
		),
		deleteNet: key.NewBinding(
			key.WithKeys("d"),
			key.WithHelp("d", "delete"),
		),
		addNetwork: key.NewBinding(
			key.WithKeys("n"),
			key.WithHelp("n", "add network"),
		),
	}
}

// NewNetworkListComponent creates a new network list component
func NewNetworkListComponent(cfg *config.Config) NetworkListComponent {
	keys := newNetworkKeyMap()
	// Use default delegate instead of custom delegate to avoid conflicts
	delegate := list.NewDefaultDelegate()
	delegate.Styles.NormalTitle = delegate.Styles.NormalTitle.
		Foreground(menuItemForeground)
	delegate.Styles.NormalDesc = delegate.Styles.NormalDesc.
		Foreground(menuItemDescriptionForeground)
	delegate.Styles.SelectedTitle = delegate.Styles.SelectedTitle.
		Foreground(menuItemForeground)
	delegate.Styles.SelectedDesc = delegate.Styles.SelectedDesc.
		Foreground(menuItemDescriptionForeground)

	networkList := list.New([]list.Item{}, delegate, 80, 20)
	networkList.Title = "🌐 Network Configuration"
	networkList.Styles.Title = titleStyle
	networkList.SetShowStatusBar(false)
	networkList.SetFilteringEnabled(false)

	c := NetworkListComponent{
		id:     "network-list",
		list:   networkList,
		keys:   keys,
		config: cfg,
	}

	c.RefreshNetworks()
	return c
}

// RefreshNetworks updates the network list with current networks
func (c *NetworkListComponent) RefreshNetworks() {
	var items []list.Item

	networkKeys := c.config.GetAllNetworkKeys()
	for _, key := range networkKeys {
		if network, exists := c.config.GetNetworkByKey(key); exists {
			status := ""
			if network.IsActive {
				status = " ✓ Active"
			}

			title := fmt.Sprintf("%s%s", network.Name, status)
			description := fmt.Sprintf("Chain ID: %d • %s", network.ChainID, network.RPCEndpoint)

			items = append(items, networkItem{
				id:          "network_" + key,
				title:       title,
				description: description,
				key:         key,
				network:     network,
			})
		}
	}

	// Add special items
	items = append(items, networkItem{
		id:          "network_add",
		title:       "➕ Add Network",
		description: "Add a new network configuration",
		key:         "add-network",
	})

	items = append(items, networkItem{
		id:          "network_back",
		title:       "🔙 Back to Settings",
		description: "Return to the settings menu",
		key:         "back-to-settings",
	})

	c.list.SetItems(items)
	// Force a refresh of the list view to ensure it displays correctly
	if len(items) > 0 {
		c.list.Select(0)
	}
}

// SetSize updates the component size
func (c *NetworkListComponent) SetSize(width, height int) {
	c.width = width
	c.height = height
	c.list.SetSize(width, height)
}

// GetSelected returns the currently selected network key
func (c *NetworkListComponent) GetSelected() string {
	if item, ok := c.list.SelectedItem().(networkItem); ok {
		return item.key
	}
	return ""
}

// Update handles messages for the network list component
func (c *NetworkListComponent) Update(msg tea.Msg) (*NetworkListComponent, tea.Cmd) {
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
				v, _ := listItem.(networkItem)
				// Check each item to see if it's in bounds.
				if zone.Get(v.id).InBounds(msg) {
					// If so, select it in the list.
					c.list.Select(i)
					// Trigger selection action
					return c, func() tea.Msg { return NetworkSelectedMsg{Key: v.key, Network: v.network} }
				}
			}
		}

		return c, nil

	case tea.KeyMsg:
		// Handle number shortcuts for quick navigation
		switch msg.String() {
		case "a":
			if item, ok := c.list.SelectedItem().(networkItem); ok && item.key != "add-network" && item.key != "back-to-settings" {
				return c, func() tea.Msg { return NetworkToggleMsg{Key: item.key} }
			}
		case "enter":
			if item, ok := c.list.SelectedItem().(networkItem); ok {
				return c, func() tea.Msg { return NetworkSelectedMsg{Key: item.key, Network: item.network} }
			}
		case "esc", "q":
			return c, func() tea.Msg { return BackToSettingsMsg{} }
		}
	}

	var cmd tea.Cmd
	c.list, cmd = c.list.Update(msg)
	return c, cmd
}

// View renders the network list component
func (c *NetworkListComponent) View() string {
	// Use zone.Scan to wrap the list view for mouse support
	return zone.Scan(appStyle.Render(c.list.View()))
}

// Network-related messages
type NetworkSelectedMsg struct {
	Key     string
	Network config.Network
}

type NetworkToggleMsg struct {
	Key string
}

type NetworkEditMsg struct {
	Key     string
	Network config.Network
}

type NetworkDeleteMsg struct {
	Key string
}

type NetworkAddRequestMsg struct{}

type BackToSettingsMsg struct{}
