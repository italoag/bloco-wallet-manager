package ui

import (
	"fmt"
	"strings"

	"blocowallet/pkg/config"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
)

// NetworkSelectionComponent represents the network selection component
type NetworkSelectionComponent struct {
	id               string
	form             *huh.Form
	width            int
	height           int
	err              error
	config           *config.Config
	selectedNetworks []string
	state            NetworkSelectionState

	// Form values
	networkChoice string
}

type NetworkSelectionState int

const (
	NetworkSelectionForm NetworkSelectionState = iota
	NetworkSelectionConfirmation
	NetworkSelectionCompleted
)

// NewNetworkSelectionComponent creates a new network selection component
func NewNetworkSelectionComponent(cfg *config.Config) NetworkSelectionComponent {
	c := NetworkSelectionComponent{
		id:     "network-selection",
		config: cfg,
		state:  NetworkSelectionForm,
	}
	c.initForm()
	return c
}

// initForm initializes the huh form for network selection
func (c *NetworkSelectionComponent) initForm() {
	// Get available networks
	allNetworks := c.config.GetAllNetworks()
	var options []huh.Option[string]

	for key, network := range allNetworks {
		label := fmt.Sprintf("%s (Chain ID: %d)", network.Name, network.ChainID)
		if network.IsActive {
			label += " ✓"
		}
		options = append(options, huh.NewOption(label, key))
	}

	c.form = huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Key("network").
				Title("Select Network").
				Description("Choose the blockchain network to use").
				Options(options...).
				Value(&c.networkChoice),
		),
	).WithWidth(80).WithShowHelp(false).WithShowErrors(false).WithTheme(huh.ThemeCharm())
}

// Init initializes the component
func (c *NetworkSelectionComponent) Init() tea.Cmd {
	// Initialize the form and focus the first field
	return c.form.Init()
}

// SetSize updates the component size
func (c *NetworkSelectionComponent) SetSize(width, height int) {
	c.width = width
	c.height = height
}

// SetError sets an error state
func (c *NetworkSelectionComponent) SetError(err error) {
	c.err = err
}

// GetSelectedNetwork returns the selected network key
func (c *NetworkSelectionComponent) GetSelectedNetwork() string {
	return c.networkChoice
}

// Reset clears the selection
func (c *NetworkSelectionComponent) Reset() {
	c.networkChoice = ""
	c.err = nil
	c.state = NetworkSelectionForm
	c.initForm()
}

// Update handles messages for the network selection component
func (c *NetworkSelectionComponent) Update(msg tea.Msg) (*NetworkSelectionComponent, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		c.width = msg.Width
		c.height = msg.Height

	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			if c.state == NetworkSelectionConfirmation {
				c.state = NetworkSelectionForm
				c.initForm()
				return c, nil
			}
			return c, func() tea.Msg { return BackToMenuMsg{} }
		case "enter":
			if c.state == NetworkSelectionConfirmation {
				// Apply network selection
				if err := c.config.SetActiveNetwork(c.networkChoice); err != nil {
					c.SetError(err)
					c.state = NetworkSelectionForm
					c.initForm()
					return c, nil
				}
				if err := c.config.Save(); err != nil {
					c.SetError(err)
					c.state = NetworkSelectionForm
					c.initForm()
					return c, nil
				}
				c.state = NetworkSelectionCompleted
				return c, func() tea.Msg { return NetworkSelectionCompleteMsg{Network: c.networkChoice} }
			}
		}

	case NetworkSelectionCompleteMsg:
		c.Reset()
		return c, func() tea.Msg { return BackToMenuMsg{} }
	}

	// Update the form only in form state
	if c.state == NetworkSelectionForm {
		form, cmd := c.form.Update(msg)
		if f, ok := form.(*huh.Form); ok {
			c.form = f
			cmds = append(cmds, cmd)
		}

		// Check if form is completed
		if c.form.State == huh.StateCompleted && c.networkChoice != "" {
			c.state = NetworkSelectionConfirmation
		}
	}

	return c, tea.Batch(cmds...)
}

// View renders the network selection component
func (c *NetworkSelectionComponent) View() string {
	var b strings.Builder

	// Header
	headerStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("86")).
		MarginBottom(1)
	b.WriteString(headerStyle.Render("🌐 Network Selection"))
	b.WriteString("\n\n")

	switch c.state {
	case NetworkSelectionForm:
		// Form
		b.WriteString(c.form.View())

		// Instructions
		b.WriteString("\n\n")
		b.WriteString(InfoStyle.Render("💡 Select the blockchain network you want to use for your wallet operations."))
		b.WriteString("\n")
		b.WriteString(InfoStyle.Render("   Only user-added networks are available."))
		b.WriteString("\n\n")

		// Footer
		b.WriteString(FooterStyle.Render("↑/↓: Navigate • Enter: Select • Esc: Back"))

	case NetworkSelectionConfirmation:
		// Show confirmation
		if network, exists := c.config.GetNetworkByKey(c.networkChoice); exists {
			confirmStyle := lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(lipgloss.Color("86")).
				Padding(1).
				MarginBottom(1)

			confirmContent := fmt.Sprintf("Network: %s\nChain ID: %d\nSymbol: %s\nRPC: %s\n\nActivate this network?",
				network.Name, network.ChainID, network.Symbol, network.RPCEndpoint)

			b.WriteString(confirmStyle.Render(confirmContent))
			b.WriteString("\n\n")
		}

		// Footer for confirmation
		b.WriteString(FooterStyle.Render("Enter: Confirm • Esc: Back"))

	case NetworkSelectionCompleted:
		// Success message
		successStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("10")).
			Bold(true)
		b.WriteString(successStyle.Render("✅ Network activated successfully!"))
		b.WriteString("\n\n")
		b.WriteString(InfoStyle.Render("Returning to main menu..."))
	}

	// Error display
	if c.err != nil {
		b.WriteString("\n\n")
		b.WriteString(ErrorStyle.Render("❌ Error: " + c.err.Error()))
	}

	return b.String()
}

// NetworkSelectionCompleteMsg is sent when network selection is completed
type NetworkSelectionCompleteMsg struct {
	Network string
}
