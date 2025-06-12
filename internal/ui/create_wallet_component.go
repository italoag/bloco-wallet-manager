package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

// CreateWalletComponent represents the wallet creation component
type CreateWalletComponent struct {
	id           string
	nameInput    textinput.Model
	passwordInput textinput.Model
	inputFocus   int
	width        int
	height       int
	err          error
	creating     bool
}

// NewCreateWalletComponent creates a new wallet creation component
func NewCreateWalletComponent() CreateWalletComponent {
	nameInput := textinput.New()
	nameInput.Placeholder = "Enter wallet name..."
	nameInput.Width = 40
	nameInput.Focus()

	passwordInput := textinput.New()
	passwordInput.Placeholder = "Enter password..."
	passwordInput.EchoMode = textinput.EchoPassword
	passwordInput.Width = 40

	return CreateWalletComponent{
		id:            "create-wallet",
		nameInput:     nameInput,
		passwordInput: passwordInput,
		inputFocus:    0,
	}
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
	return c.nameInput.Value()
}

// GetPassword returns the entered password
func (c *CreateWalletComponent) GetPassword() string {
	return c.passwordInput.Value()
}

// Reset clears all inputs
func (c *CreateWalletComponent) Reset() {
	c.nameInput.SetValue("")
	c.passwordInput.SetValue("")
	c.inputFocus = 0
	c.err = nil
	c.creating = false
	c.nameInput.Focus()
	c.passwordInput.Blur()
}

// Update handles messages for the create wallet component
func (c *CreateWalletComponent) Update(msg tea.Msg) (*CreateWalletComponent, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		c.width = msg.Width
		c.height = msg.Height

	case tea.KeyMsg:
		switch {
		case key.Matches(msg, key.NewBinding(key.WithKeys("tab", "shift+tab", "enter", "up", "down"))):
			if key.Matches(msg, key.NewBinding(key.WithKeys("enter"))) {
				if c.inputFocus == 1 { // Password field, attempt to create wallet
					if c.validateInputs() {
						c.creating = true
						return c, func() tea.Msg {
							return CreateWalletRequestMsg{
								Name:     c.nameInput.Value(),
								Password: c.passwordInput.Value(),
							}
						}
					}
				} else {
					// Move to next field
					c.inputFocus++
					if c.inputFocus > 1 {
						c.inputFocus = 0
					}
				}
			} else {
				// Handle tab navigation
				if key.Matches(msg, key.NewBinding(key.WithKeys("shift+tab", "up"))) {
					c.inputFocus--
					if c.inputFocus < 0 {
						c.inputFocus = 1
					}
				} else {
					c.inputFocus++
					if c.inputFocus > 1 {
						c.inputFocus = 0
					}
				}
			}

			// Update focus
			if c.inputFocus == 0 {
				c.nameInput.Focus()
				c.passwordInput.Blur()
			} else {
				c.nameInput.Blur()
				c.passwordInput.Focus()
			}

		case key.Matches(msg, key.NewBinding(key.WithKeys("esc"))):
			return c, func() tea.Msg { return BackToMenuMsg{} }
		}

	case walletCreatedMsg:
		c.Reset()
		return c, func() tea.Msg { return BackToMenuMsg{} }

	case errorMsg:
		c.SetError(fmt.Errorf("%s", string(msg)))
	}

	// Update the focused input
	var cmd tea.Cmd
	if c.inputFocus == 0 {
		c.nameInput, cmd = c.nameInput.Update(msg)
		cmds = append(cmds, cmd)
	} else {
		c.passwordInput, cmd = c.passwordInput.Update(msg)
		cmds = append(cmds, cmd)
	}

	return c, tea.Batch(cmds...)
}

// validateInputs checks if the inputs are valid
func (c *CreateWalletComponent) validateInputs() bool {
	if strings.TrimSpace(c.nameInput.Value()) == "" {
		c.err = fmt.Errorf("Wallet name cannot be empty")
		return false
	}
	if strings.TrimSpace(c.passwordInput.Value()) == "" {
		c.err = fmt.Errorf("Password cannot be empty")
		return false
	}
	c.err = nil
	return true
}

// View renders the create wallet component
func (c *CreateWalletComponent) View() string {
	var b strings.Builder

	b.WriteString(HeaderStyle.Render("➕ Create New Wallet"))
	b.WriteString("\n\n")

	// Name input
	b.WriteString(LabelStyle.Render("Wallet Name:"))
	b.WriteString("\n")
	if c.inputFocus == 0 {
		b.WriteString(FocusedInputStyle.Render(c.nameInput.View()))
	} else {
		b.WriteString(InputStyle.Render(c.nameInput.View()))
	}
	b.WriteString("\n\n")

	// Password input
	b.WriteString(LabelStyle.Render("Password:"))
	b.WriteString("\n")
	if c.inputFocus == 1 {
		b.WriteString(FocusedInputStyle.Render(c.passwordInput.View()))
	} else {
		b.WriteString(InputStyle.Render(c.passwordInput.View()))
	}
	b.WriteString("\n\n")

	// Status messages
	if c.creating {
		b.WriteString(LoadingStyle.Render("⏳ Creating wallet..."))
		b.WriteString("\n")
	} else if c.err != nil {
		b.WriteString(ErrorStyle.Render("❌ Error: " + c.err.Error()))
		b.WriteString("\n")
	}

	// Instructions
	b.WriteString(InfoStyle.Render("💡 Your wallet will be secured with a mnemonic phrase."))
	b.WriteString("\n")
	b.WriteString(InfoStyle.Render("   Make sure to save it in a secure location!"))
	b.WriteString("\n\n")

	// Footer
	b.WriteString(FooterStyle.Render("Tab: Next Field • Enter: Create • Esc: Back"))

	return b.String()
}

// CreateWalletRequestMsg is sent when the user wants to create a wallet
type CreateWalletRequestMsg struct {
	Name     string
	Password string
}

// BackToMenuMsg is sent when the user wants to go back to the menu
type BackToMenuMsg struct{}