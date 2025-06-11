package ui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
)

// CreateWalletComponent represents the wallet creation component
type CreateWalletComponent struct {
	id         string
	form       *huh.Form
	width      int
	height     int
	err        error
	creating   bool
	walletName string
	password   string
}

// NewCreateWalletComponent creates a new wallet creation component
func NewCreateWalletComponent() CreateWalletComponent {
	c := CreateWalletComponent{
		id: "create-wallet",
	}
	c.initForm()
	return c
}

// initForm initializes the huh form
func (c *CreateWalletComponent) initForm() {
	// Reset the variables
	c.walletName = ""
	c.password = ""

	// Create form with Value() pointing to our variables
	c.form = huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Key("walletName").
				Title("Wallet Name").
				Placeholder("Enter wallet name...").
				Value(&c.walletName).
				Validate(func(s string) error {
					if strings.TrimSpace(s) == "" {
						return fmt.Errorf("wallet name cannot be empty")
					}
					return nil
				}),
			huh.NewInput().
				Key("password").
				Title("Password").
				Placeholder("Enter password...").
				Value(&c.password).
				EchoMode(huh.EchoModePassword).
				Validate(func(s string) error {
					if strings.TrimSpace(s) == "" {
						return fmt.Errorf("password cannot be empty")
					}
					return nil
				}),
		),
	).WithWidth(60).WithShowHelp(false).WithShowErrors(false)
}

// SetSize updates the component size
func (c *CreateWalletComponent) SetSize(width, height int) {
	c.width = width
	c.height = height
}

// SetError sets an error state
func (c *CreateWalletComponent) SetError(err error) {
	c.err = err
	c.creating = false
}

// SetCreating sets the creating state
func (c *CreateWalletComponent) SetCreating(creating bool) {
	c.creating = creating
	if creating {
		c.err = nil
	}
}

// GetWalletName returns the entered wallet name
func (c *CreateWalletComponent) GetWalletName() string {
	if c.form != nil {
		return strings.TrimSpace(c.form.GetString("walletName"))
	}
	return strings.TrimSpace(c.walletName)
}

// GetPassword returns the entered password
func (c *CreateWalletComponent) GetPassword() string {
	if c.form != nil {
		return strings.TrimSpace(c.form.GetString("password"))
	}
	return strings.TrimSpace(c.password)
}

// Reset clears all inputs
func (c *CreateWalletComponent) Reset() {
	c.err = nil
	c.creating = false
	c.initForm()
}

// Init initializes the component
func (c *CreateWalletComponent) Init() tea.Cmd {
	// Initialize the form so the first field receives focus
	return c.form.Init()
}

// Update handles messages for the create wallet component
func (c *CreateWalletComponent) Update(msg tea.Msg) (*CreateWalletComponent, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		c.width = msg.Width
		c.height = msg.Height

	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			return c, func() tea.Msg { return BackToMenuMsg{} }
		}

	case walletCreatedMsg:
		c.Reset()
		return c, func() tea.Msg { return BackToMenuMsg{} }

	case errorMsg:
		c.SetError(fmt.Errorf("%s", string(msg)))
	}

	// Process the form - exactly like the example
	form, cmd := c.form.Update(msg)
	if f, ok := form.(*huh.Form); ok {
		c.form = f
		cmds = append(cmds, cmd)
	}
	// Check if form is completed - exactly like the example
	if c.form.State == huh.StateCompleted && !c.creating {
		// Get values directly from form instead of variables
		walletName := strings.TrimSpace(c.form.GetString("walletName"))
		password := strings.TrimSpace(c.form.GetString("password"))

		// Create wallet if form is valid (form validation handles empty values)
		c.creating = true
		return c, func() tea.Msg {
			return CreateWalletRequestMsg{
				Name:     walletName,
				Password: password,
			}
		}
	}

	return c, tea.Batch(cmds...)
}

// View renders the create wallet component
func (c *CreateWalletComponent) View() string {
	var b strings.Builder

	// Header
	headerStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("86")).
		MarginBottom(1)
	b.WriteString(headerStyle.Render("➕ Create New Wallet"))
	b.WriteString("\n\n")

	// Form
	b.WriteString(c.form.View())

	// Status messages
	if c.creating {
		b.WriteString("\n")
		b.WriteString(LoadingStyle.Render("⏳ Creating wallet..."))
	} else if c.err != nil {
		b.WriteString("\n")
		b.WriteString(ErrorStyle.Render("❌ Error: " + c.err.Error()))
	}

	// Show form errors if any
	errors := c.form.Errors()
	if len(errors) > 0 {
		b.WriteString("\n")
		for _, err := range errors {
			b.WriteString(ErrorStyle.Render("❌ " + err.Error()))
			b.WriteString("\n")
		}
	}

	// Instructions
	b.WriteString("\n")
	b.WriteString(InfoStyle.Render("💡 Your wallet will be secured with a mnemonic phrase."))
	b.WriteString("\n")
	b.WriteString(InfoStyle.Render("   Make sure to save it in a secure location!"))
	b.WriteString("\n\n")

	// Footer
	b.WriteString(FooterStyle.Render("Tab/Arrow Keys: Navigate • Enter: Create • Esc: Back"))

	return b.String()
}

// CreateWalletRequestMsg is sent when the user wants to create a wallet
type CreateWalletRequestMsg struct {
	Name     string
	Password string
}

// BackToMenuMsg is sent when the user wants to go back to the menu
type BackToMenuMsg struct{}
