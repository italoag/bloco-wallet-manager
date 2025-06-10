package ui

import (
	"fmt"
	"strconv"
	"strings"

	"blocowallet/internal/blockchain"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
)

// AddNetworkComponent represents the add network component
type AddNetworkComponent struct {
	id     string
	form   *huh.Form
	width  int
	height int
	err    error
	adding bool

	// Form values
	networkName string
	chainID     string
	rpcEndpoint string
	symbol      string

	// Chain service for suggestions
	chainListService *blockchain.ChainListService
}

// NewAddNetworkComponent creates a new add network component
func NewAddNetworkComponent() AddNetworkComponent {
	c := AddNetworkComponent{
		id:               "add-network",
		chainListService: blockchain.NewChainListService(),
	}
	c.initForm()
	return c
}

// initForm initializes the huh form
func (c *AddNetworkComponent) initForm() {
	c.form = huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Key("networkName").
				Title("Network Name").
				Placeholder("Type to search networks (e.g., Polygon, Binance)").
				Suggestions([]string{
					"Polygon", "Binance Smart Chain", "Ethereum", "Avalanche",
					"Fantom", "Arbitrum", "Optimism", "Base", "Linea",
					"Polygon zkEVM", "zkSync Era", "Mantle", "Scroll",
				}).
				Value(&c.networkName),
			huh.NewInput().
				Key("chainID").
				Title("Chain ID").
				Placeholder("Enter chain ID (e.g., 137 for Polygon)").
				Value(&c.chainID),
			huh.NewInput().
				Key("symbol").
				Title("Native Currency Symbol").
				Placeholder("Enter currency symbol (e.g., MATIC, BNB)").
				Value(&c.symbol),
			huh.NewInput().
				Key("rpcEndpoint").
				Title("RPC Endpoint").
				Placeholder("Enter RPC URL (e.g., https://polygon-rpc.com)").
				Value(&c.rpcEndpoint),
		),
	).WithWidth(80).WithShowHelp(false).WithShowErrors(false).WithTheme(huh.ThemeCharm())
}

// SetSize updates the component size
func (c *AddNetworkComponent) SetSize(width, height int) {
	c.width = width
	c.height = height
}

// SetError sets an error state
func (c *AddNetworkComponent) SetError(err error) {
	c.err = err
	c.adding = false
}

// SetAdding sets the adding state
func (c *AddNetworkComponent) SetAdding(adding bool) {
	c.adding = adding
	if adding {
		c.err = nil
	}
}

// GetNetworkName returns the entered network name
func (c *AddNetworkComponent) GetNetworkName() string {
	return c.networkName
}

// GetChainID returns the entered chain ID as integer
func (c *AddNetworkComponent) GetChainID() (int, error) {
	chainID, err := strconv.Atoi(strings.TrimSpace(c.chainID))
	if err != nil {
		return 0, fmt.Errorf("invalid chain ID: must be a number")
	}
	return chainID, nil
}

// GetSymbol returns the entered symbol
func (c *AddNetworkComponent) GetSymbol() string {
	return c.symbol
}

// GetRPCEndpoint returns the entered RPC endpoint
func (c *AddNetworkComponent) GetRPCEndpoint() string {
	return c.rpcEndpoint
}

// Reset clears all inputs
func (c *AddNetworkComponent) Reset() {
	c.networkName = ""
	c.chainID = ""
	c.symbol = ""
	c.rpcEndpoint = ""
	c.err = nil
	c.adding = false
	c.initForm()
}

// Init initializes the component
func (c *AddNetworkComponent) Init() tea.Cmd {
	// Initialize the form so the first field is focused
	return c.form.Init()
}

// Update handles messages for the add network component
func (c *AddNetworkComponent) Update(msg tea.Msg) (*AddNetworkComponent, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		c.width = msg.Width
		c.height = msg.Height

	case networkAddedMsg:
		c.Reset()
		return c, func() tea.Msg { return BackToNetworkListMsg{} }

	case errorMsg:
		c.SetError(fmt.Errorf("%s", string(msg)))

	case tea.KeyMsg:
		// Verificar se é uma tecla que devemos tratar especificamente antes do form
		switch msg.String() {
		case "tab":
			// Se temos uma mensagem tab, vamos garantir que será processada pelo formulário
			break
		}
	}

	// Update the form first (allows typing and internal navigation)
	form, cmd := c.form.Update(msg)
	if f, ok := form.(*huh.Form); ok {
		c.form = f
		cmds = append(cmds, cmd)
	}

	// Only handle escape if form didn't handle it (when form is not focused on input)
	if keyMsg, ok := msg.(tea.KeyMsg); ok && keyMsg.String() == "esc" && c.form.State == huh.StateNormal {
		return c, func() tea.Msg { return BackToNetworkListMsg{} }
	}

	// Check if form is completed
	if c.form.State == huh.StateCompleted {
		if c.validateInputs() {
			c.adding = true
			return c, func() tea.Msg {
				return AddNetworkRequestMsg{
					Name:        c.GetNetworkName(),
					ChainID:     c.chainID, // Pass as string for validation in handler
					Symbol:      c.GetSymbol(),
					RPCEndpoint: c.GetRPCEndpoint(),
				}
			}
		}
		// Reset form state if validation failed
		c.form.State = huh.StateNormal
	}

	return c, tea.Batch(cmds...)
}

// validateInputs checks if the inputs are valid
func (c *AddNetworkComponent) validateInputs() bool {
	if strings.TrimSpace(c.networkName) == "" {
		c.err = fmt.Errorf("Network name cannot be empty")
		return false
	}

	if strings.TrimSpace(c.chainID) == "" {
		c.err = fmt.Errorf("Chain ID cannot be empty")
		return false
	}

	// Validate chain ID is a number
	if _, err := c.GetChainID(); err != nil {
		c.err = err
		return false
	}

	if strings.TrimSpace(c.symbol) == "" {
		c.err = fmt.Errorf("Currency symbol cannot be empty")
		return false
	}

	if strings.TrimSpace(c.rpcEndpoint) == "" {
		c.err = fmt.Errorf("RPC endpoint cannot be empty")
		return false
	}

	// Basic URL validation
	rpc := strings.TrimSpace(c.rpcEndpoint)
	if !strings.HasPrefix(rpc, "http://") && !strings.HasPrefix(rpc, "https://") {
		c.err = fmt.Errorf("RPC endpoint must start with http:// or https://")
		return false
	}

	c.err = nil
	return true
}

// View renders the add network component
func (c *AddNetworkComponent) View() string {
	var b strings.Builder

	// Header
	headerStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("86")).
		MarginBottom(1)
	b.WriteString(headerStyle.Render("🌐 Add Custom Network"))
	b.WriteString("\n\n")

	// Form
	b.WriteString(c.form.View())

	// Status messages
	if c.adding {
		b.WriteString("\n")
		b.WriteString(LoadingStyle.Render("⏳ Adding network..."))
	} else if c.err != nil {
		b.WriteString("\n")
		b.WriteString(ErrorStyle.Render("❌ Error: " + c.err.Error()))
	}

	// Instructions
	b.WriteString("\n\n")
	b.WriteString(WarningStyle.Render("💡 Tips:"))
	b.WriteString("\n")
	b.WriteString(InfoStyle.Render("   • Chain ID must be unique (check chainlist.org for reference)"))
	b.WriteString("\n")
	b.WriteString(InfoStyle.Render("   • Use reliable RPC endpoints for better performance"))
	b.WriteString("\n")
	b.WriteString(InfoStyle.Render("   • Test the network before adding important transactions"))
	b.WriteString("\n\n")

	// Footer
	b.WriteString(FooterStyle.Render("Tab/Arrow Keys: Navigate • Enter: Add Network • Esc: Back"))

	return b.String()
}

// AddNetworkRequestMsg is sent when the user wants to add a network
type AddNetworkRequestMsg struct {
	Name        string
	ChainID     string
	Symbol      string
	RPCEndpoint string
}

type BackToNetworkListMsg struct{}

// networkAddedMsg is sent when a network is successfully added
type networkAddedMsg struct {
	key     string
	network interface{} // Will be config.Network
}
