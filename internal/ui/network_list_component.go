package ui

import (
	"blocowallet/pkg/config"
	"fmt"
	"strings"

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
	title       string
	description string
	key         string
	network     config.Network
}

func (i networkItem) Title() string       { return i.title }
func (i networkItem) Description() string { return i.description }
func (i networkItem) FilterValue() string { return i.title }

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
	delegate := newNetworkDelegate(keys)

	networkList := list.New([]list.Item{}, delegate, 0, 0)
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

// newNetworkDelegate creates a delegate for the network list
func newNetworkDelegate(keys *networkKeyMap) list.DefaultDelegate {
	d := list.NewDefaultDelegate()

	d.UpdateFunc = func(msg tea.Msg, m *list.Model) tea.Cmd {
		var item networkItem
		var ok bool

		if i := m.SelectedItem(); i != nil {
			item, ok = i.(networkItem)
			if !ok {
				return nil
			}
		}

		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch {
			case key.Matches(msg, keys.choose):
				return func() tea.Msg {
					return NetworkSelectedMsg{Key: item.key, Network: item.network}
				}
			case key.Matches(msg, keys.toggle):
				return func() tea.Msg {
					return NetworkToggleMsg{Key: item.key}
				}
			case key.Matches(msg, keys.edit):
				return func() tea.Msg {
					return NetworkEditMsg{Key: item.key, Network: item.network}
				}
			case key.Matches(msg, keys.deleteNet):
				if item.network.IsCustom {
					return func() tea.Msg {
						return NetworkDeleteMsg{Key: item.key}
					}
				}
			case key.Matches(msg, keys.addNetwork):
				return func() tea.Msg {
					return NetworkAddRequestMsg{}
				}
			}
		}

		return nil
	}

	return d
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
			customTag := ""
			if network.IsCustom {
				customTag = " [Custom]"
			}

			title := fmt.Sprintf("%s%s%s", network.Name, status, customTag)
			description := fmt.Sprintf("Chain ID: %d • %s", network.ChainID, network.RPCEndpoint)

			items = append(items, networkItem{
				title:       title,
				description: description,
				key:         key,
				network:     network,
			})
		}
	}

	// Add special items
	items = append(items, networkItem{
		title:       "➕ Add Custom Network",
		description: "Add a new custom network configuration",
		key:         "add-network",
	})

	items = append(items, networkItem{
		title:       "🔙 Back to Settings",
		description: "Return to the settings menu",
		key:         "back-to-settings",
	})

	c.list.SetItems(items)
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
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		c.width = msg.Width
		c.height = msg.Height
		c.list.SetSize(msg.Width, msg.Height)

	case tea.MouseMsg:
		// Check if any network item was clicked
		for i := 0; i < len(c.list.Items()); i++ {
			itemZoneID := fmt.Sprintf("network-item-%d", i)
			if zone.Get(itemZoneID).InBounds(msg) {
				// Select and activate the clicked item
				c.list.Select(i)
				if item, ok := c.list.SelectedItem().(networkItem); ok {
					return c, func() tea.Msg { return NetworkSelectedMsg{Key: item.key, Network: item.network} }
				}
			}
		}

	case tea.KeyMsg:
		// Handle number shortcuts for quick navigation
		switch msg.String() {
		case "esc", "q":
			return c, func() tea.Msg { return BackToSettingsMsg{} }
		}
	}

	// Update the list
	c.list, cmd = c.list.Update(msg)
	return c, cmd
}

// View renders the network list component
func (c *NetworkListComponent) View() string {
	// Render the list first
	listView := c.list.View()

	// Apply zone marking to each network item for mouse support
	lines := strings.Split(listView, "\n")
	var markedLines []string

	itemIndex := 0
	for _, line := range lines {
		// Check if this line contains a network item (has content and isn't just formatting)
		if strings.TrimSpace(line) != "" &&
			!strings.Contains(line, "Network Configuration") &&
			!strings.Contains(line, "Help") &&
			(strings.Contains(line, "►") || strings.Contains(line, "•") || strings.Contains(line, "🌐") || strings.Contains(line, "🔙") || strings.Contains(line, "➕")) {

			// Mark this line as clickable
			zoneID := fmt.Sprintf("network-item-%d", itemIndex)
			markedLine := zone.Mark(zoneID, line)
			markedLines = append(markedLines, markedLine)
			itemIndex++
		} else {
			markedLines = append(markedLines, line)
		}
	}

	return appStyle.Render(strings.Join(markedLines, "\n"))
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
