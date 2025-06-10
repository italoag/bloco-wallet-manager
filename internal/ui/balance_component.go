package ui

import (
	"blocowallet/internal/wallet"
	"fmt"
	"math/big"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

// BalanceComponent represents the balance display component
type BalanceComponent struct {
	id             string
	selectedWallet *wallet.Wallet
	balance        *wallet.Balance
	multiBalance   *wallet.MultiNetworkBalance
	loading        bool
	err            error
	width          int
	height         int
}

// NewBalanceComponent creates a new balance component
func NewBalanceComponent() BalanceComponent {
	return BalanceComponent{
		id: "balance-display",
	}
}

// SetWallet updates the selected wallet
func (c *BalanceComponent) SetWallet(w *wallet.Wallet) {
	c.selectedWallet = w
	c.balance = nil
	c.multiBalance = nil
	c.err = nil
}

// SetBalance updates the balance information
func (c *BalanceComponent) SetBalance(balance *wallet.Balance) {
	c.balance = balance
	c.loading = false
	c.err = nil
}

// SetMultiBalance updates the multi-network balance information
func (c *BalanceComponent) SetMultiBalance(multiBalance *wallet.MultiNetworkBalance) {
	c.multiBalance = multiBalance
	c.loading = false
	c.err = nil
}

// SetLoading sets the loading state
func (c *BalanceComponent) SetLoading(loading bool) {
	c.loading = loading
	if loading {
		c.err = nil
	}
}

// SetError sets an error state
func (c *BalanceComponent) SetError(err error) {
	c.err = err
	c.loading = false
}

// SetSize updates the component size
func (c *BalanceComponent) SetSize(width, height int) {
	c.width = width
	c.height = height
}

// Update handles messages for the balance component
func (c *BalanceComponent) Update(msg tea.Msg) (*BalanceComponent, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		c.width = msg.Width
		c.height = msg.Height

	case balanceLoadedMsg:
		c.SetBalance(msg)

	case multiBalanceLoadedMsg:
		c.SetMultiBalance(msg)

	case errorMsg:
		c.SetError(fmt.Errorf("%s", string(msg)))
	}

	return c, nil
}

// View renders the balance component
func (c *BalanceComponent) View() string {
	if c.selectedWallet == nil {
		return InfoStyle.Render("Select a wallet to view balance information.")
	}

	var b strings.Builder

	// Wallet header
	b.WriteString(HeaderStyle.Render("💼 " + c.selectedWallet.Name))
	b.WriteString("\n")
	b.WriteString(AddressStyle.Render("Address: " + c.selectedWallet.Address))
	b.WriteString("\n\n")

	// Loading state
	if c.loading {
		b.WriteString(LoadingStyle.Render("⏳ Loading balance..."))
		return b.String()
	}

	// Error state
	if c.err != nil {
		b.WriteString(ErrorStyle.Render("❌ Error loading balance: " + c.err.Error()))
		return b.String()
	}

	// Multi-network balance display
	if c.multiBalance != nil {
		b.WriteString(HeaderStyle.Render("🌐 Multi-Network Balance"))
		b.WriteString("\n\n")

		for _, networkBalance := range c.multiBalance.NetworkBalances {
			if networkBalance.Error != nil {
				b.WriteString(fmt.Sprintf("%s %s: %s\n",
					ErrorStyle.Render("❌"),
					LabelStyle.Render(networkBalance.NetworkName),
					ErrorStyle.Render(networkBalance.Error.Error())))
			} else {
				// Convert big.Int to string for display
				balanceStr := "0"
				if networkBalance.Amount != nil {
					// Convert from wei to ether (divide by 10^18)
					amount := new(big.Float).SetInt(networkBalance.Amount)
					divisor := new(big.Float).SetFloat64(1e18)
					result := new(big.Float).Quo(amount, divisor)
					balanceStr = result.Text('f', 6)
				}

				b.WriteString(fmt.Sprintf("%s %s: %s %s\n",
					SuccessStyle.Render("💰"),
					LabelStyle.Render(networkBalance.NetworkName),
					BalanceStyle.Render(balanceStr),
					NetworkStyle.Render(networkBalance.Symbol)))
			}
		}

		return b.String()
	}

	// Single network balance display
	if c.balance != nil {
		// Convert big.Int to string for display
		balanceStr := "0"
		if c.balance.Amount != nil {
			// Convert from wei to ether (divide by 10^18)
			amount := new(big.Float).SetInt(c.balance.Amount)
			divisor := new(big.Float).SetFloat64(1e18)
			result := new(big.Float).Quo(amount, divisor)
			balanceStr = result.Text('f', 6)
		}

		b.WriteString(fmt.Sprintf("%s %s %s",
			SuccessStyle.Render("💰 Balance:"),
			BalanceStyle.Render(balanceStr),
			NetworkStyle.Render(c.balance.Symbol)))

		return b.String()
	}

	// No balance data
	b.WriteString(InfoStyle.Render("No balance information available."))
	return b.String()
}
