package ui

import (
	"blocowallet/internal/wallet"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	zone "github.com/lrstanley/bubblezone"
)

// WalletListComponent represents the wallet list component
type WalletListComponent struct {
	id       string
	wallets  []*wallet.Wallet
	selected int
	width    int
	height   int
}

// NewWalletListComponent creates a new wallet list component
func NewWalletListComponent() WalletListComponent {
	return WalletListComponent{
		id:       "wallet-list",
		selected: 0,
	}
}

// SetWallets updates the wallet list
func (c *WalletListComponent) SetWallets(wallets []*wallet.Wallet) {
	c.wallets = wallets
	if c.selected >= len(c.wallets) && len(c.wallets) > 0 {
		c.selected = len(c.wallets) - 1
	}
	if len(c.wallets) == 0 {
		c.selected = 0
	}
}

// SetSize updates the component size
func (c *WalletListComponent) SetSize(width, height int) {
	c.width = width
	c.height = height
}

// GetSelected returns the currently selected wallet
func (c *WalletListComponent) GetSelected() *wallet.Wallet {
	if len(c.wallets) == 0 || c.selected < 0 || c.selected >= len(c.wallets) {
		return nil
	}
	return c.wallets[c.selected]
}

// GetSelectedIndex returns the selected index
func (c *WalletListComponent) GetSelectedIndex() int {
	return c.selected
}

// Update handles messages for the wallet list component
func (c *WalletListComponent) Update(msg tea.Msg) (*WalletListComponent, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		c.width = msg.Width
		c.height = msg.Height

	case tea.MouseMsg:
		if msg.Action != tea.MouseActionRelease || msg.Button != tea.MouseButtonLeft {
			return c, nil
		}

		// Check clicks on wallet items
		for i := range c.wallets {
			if zone.Get(c.id + "-wallet-" + string(rune(i))).InBounds(msg) {
				c.selected = i
				return c, func() tea.Msg { return WalletSelectedMsg{Wallet: c.wallets[i]} }
			}
		}

	case tea.KeyMsg:
		if len(c.wallets) == 0 {
			return c, nil
		}

		switch {
		case key.Matches(msg, key.NewBinding(key.WithKeys("up", "k"))):
			if c.selected > 0 {
				c.selected--
			}
		case key.Matches(msg, key.NewBinding(key.WithKeys("down", "j"))):
			if c.selected < len(c.wallets)-1 {
				c.selected++
			}
		case key.Matches(msg, key.NewBinding(key.WithKeys("enter"))):
			if c.selected >= 0 && c.selected < len(c.wallets) {
				return c, func() tea.Msg { return WalletSelectedMsg{Wallet: c.wallets[c.selected]} }
			}
		}
	}

	return c, nil
}

// View renders the wallet list component
func (c *WalletListComponent) View() string {
	if len(c.wallets) == 0 {
		return InfoStyle.Render("No wallets found. Create a new wallet to get started.")
	}

	var b strings.Builder

	for i, w := range c.wallets {
		var itemText string
		if i == c.selected {
			itemText = SelectedStyle.Render("► " + w.Name + " - " + w.Address[:10] + "..." + w.Address[len(w.Address)-6:])
		} else {
			itemText = ItemStyle.Render("  " + w.Name + " - " + w.Address[:10] + "..." + w.Address[len(w.Address)-6:])
		}

		// Mark zone for mouse interaction
		b.WriteString(zone.Mark(c.id+"-wallet-"+string(rune(i)), itemText))
		if i < len(c.wallets)-1 {
			b.WriteString("\n")
		}
	}

	return b.String()
}

// WalletSelectedMsg is sent when a wallet is selected
type WalletSelectedMsg struct {
	Wallet *wallet.Wallet
}