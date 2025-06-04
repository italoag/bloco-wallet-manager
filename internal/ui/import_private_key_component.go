package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

// ImportPrivateKeyComponent represents the private key import component
type ImportPrivateKeyComponent struct {
	id               string
	nameInput        textinput.Model
	privateKeyInput  textinput.Model
	passwordInput    textinput.Model
	inputFocus       int
	width            int
	height           int
	err              error
	importing        bool
}

// NewImportPrivateKeyComponent creates a new private key import component
func NewImportPrivateKeyComponent() ImportPrivateKeyComponent {
	nameInput := textinput.New()
	nameInput.Placeholder = "Enter wallet name..."
	nameInput.Width = 40
	nameInput.Focus()

	privateKeyInput := textinput.New()
	privateKeyInput.Placeholder = "Enter private key (without 0x prefix)..."
	privateKeyInput.EchoMode = textinput.EchoPassword
	privateKeyInput.Width = 60

	passwordInput := textinput.New()
	passwordInput.Placeholder = "Enter password..."
	passwordInput.EchoMode = textinput.EchoPassword
	passwordInput.Width = 40

	return ImportPrivateKeyComponent{
		id:              "import-private-key",
		nameInput:       nameInput,
		privateKeyInput: privateKeyInput,
		passwordInput:   passwordInput,
		inputFocus:      0,
	}
}

// SetSize updates the component size
func (c *ImportPrivateKeyComponent) SetSize(width, height int) {
	c.width = width
	c.height = height
}

// SetError sets an error state
func (c *ImportPrivateKeyComponent) SetError(err error) {
	c.err = err
	c.importing = false
}

// SetImporting sets the importing state
func (c *ImportPrivateKeyComponent) SetImporting(importing bool) {
	c.importing = importing
	if importing {
		c.err = nil
	}
}

// GetWalletName returns the entered wallet name
func (c *ImportPrivateKeyComponent) GetWalletName() string {
	return c.nameInput.Value()
}

// GetPrivateKey returns the entered private key
func (c *ImportPrivateKeyComponent) GetPrivateKey() string {
	return c.privateKeyInput.Value()
}

// GetPassword returns the entered password
func (c *ImportPrivateKeyComponent) GetPassword() string {
	return c.passwordInput.Value()
}

// Reset clears all inputs
func (c *ImportPrivateKeyComponent) Reset() {
	c.nameInput.SetValue("")
	c.privateKeyInput.SetValue("")
	c.passwordInput.SetValue("")
	c.inputFocus = 0
	c.err = nil
	c.importing = false
	c.nameInput.Focus()
	c.privateKeyInput.Blur()
	c.passwordInput.Blur()
}

// Update handles messages for the import private key component
func (c *ImportPrivateKeyComponent) Update(msg tea.Msg) (*ImportPrivateKeyComponent, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		c.width = msg.Width
		c.height = msg.Height

	case tea.KeyMsg:
		switch {
		case key.Matches(msg, key.NewBinding(key.WithKeys("tab", "shift+tab", "enter", "up", "down"))):
			if key.Matches(msg, key.NewBinding(key.WithKeys("enter"))) {
				if c.inputFocus == 2 { // Password field, attempt to import wallet
					if c.validateInputs() {
						c.importing = true
						return c, func() tea.Msg {
							return ImportPrivateKeyRequestMsg{
								Name:       c.nameInput.Value(),
								PrivateKey: c.privateKeyInput.Value(),
								Password:   c.passwordInput.Value(),
							}
						}
					}
				} else {
					// Move to next field
					c.inputFocus++
					if c.inputFocus > 2 {
						c.inputFocus = 0
					}
				}
			} else {
				// Handle tab navigation
				if key.Matches(msg, key.NewBinding(key.WithKeys("shift+tab", "up"))) {
					c.inputFocus--
					if c.inputFocus < 0 {
						c.inputFocus = 2
					}
				} else {
					c.inputFocus++
					if c.inputFocus > 2 {
						c.inputFocus = 0
					}
				}
			}

			// Update focus
			c.updateFocus()

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
	switch c.inputFocus {
	case 0:
		c.nameInput, cmd = c.nameInput.Update(msg)
		cmds = append(cmds, cmd)
	case 1:
		c.privateKeyInput, cmd = c.privateKeyInput.Update(msg)
		cmds = append(cmds, cmd)
	case 2:
		c.passwordInput, cmd = c.passwordInput.Update(msg)
		cmds = append(cmds, cmd)
	}

	return c, tea.Batch(cmds...)
}

// updateFocus updates which input field has focus
func (c *ImportPrivateKeyComponent) updateFocus() {
	c.nameInput.Blur()
	c.privateKeyInput.Blur()
	c.passwordInput.Blur()

	switch c.inputFocus {
	case 0:
		c.nameInput.Focus()
	case 1:
		c.privateKeyInput.Focus()
	case 2:
		c.passwordInput.Focus()
	}
}

// validateInputs checks if the inputs are valid
func (c *ImportPrivateKeyComponent) validateInputs() bool {
	if strings.TrimSpace(c.nameInput.Value()) == "" {
		c.err = fmt.Errorf("Wallet name cannot be empty")
		return false
	}
	if strings.TrimSpace(c.privateKeyInput.Value()) == "" {
		c.err = fmt.Errorf("Private key cannot be empty")
		return false
	}
	if strings.TrimSpace(c.passwordInput.Value()) == "" {
		c.err = fmt.Errorf("Password cannot be empty")
		return false
	}
	
	// Basic private key validation
	privateKey := strings.TrimSpace(c.privateKeyInput.Value())
	privateKey = strings.TrimPrefix(privateKey, "0x")
	privateKey = strings.TrimPrefix(privateKey, "0X")
	
	if len(privateKey) != 64 {
		c.err = fmt.Errorf("Private key must be 64 characters long (32 bytes)")
		return false
	}
	
	// Check if it's a valid hex string
	for _, char := range privateKey {
		if !((char >= '0' && char <= '9') || (char >= 'a' && char <= 'f') || (char >= 'A' && char <= 'F')) {
			c.err = fmt.Errorf("Private key must contain only hexadecimal characters")
			return false
		}
	}
	
	c.err = nil
	return true
}

// View renders the import private key component
func (c *ImportPrivateKeyComponent) View() string {
	var b strings.Builder

	b.WriteString(HeaderStyle.Render("🔑 Import Wallet from Private Key"))
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

	// Private key input
	b.WriteString(LabelStyle.Render("Private Key (32 bytes / 64 hex characters):"))
	b.WriteString("\n")
	if c.inputFocus == 1 {
		b.WriteString(FocusedInputStyle.Render(c.privateKeyInput.View()))
	} else {
		b.WriteString(InputStyle.Render(c.privateKeyInput.View()))
	}
	b.WriteString("\n\n")

	// Password input
	b.WriteString(LabelStyle.Render("Password:"))
	b.WriteString("\n")
	if c.inputFocus == 2 {
		b.WriteString(FocusedInputStyle.Render(c.passwordInput.View()))
	} else {
		b.WriteString(InputStyle.Render(c.passwordInput.View()))
	}
	b.WriteString("\n\n")

	// Status messages
	if c.importing {
		b.WriteString(LoadingStyle.Render("⏳ Importing wallet..."))
		b.WriteString("\n")
	} else if c.err != nil {
		b.WriteString(ErrorStyle.Render("❌ Error: " + c.err.Error()))
		b.WriteString("\n")
	}

	// Instructions
	b.WriteString(WarningStyle.Render("⚠️  Important: Never share your private key with anyone!"))
	b.WriteString("\n")
	b.WriteString(InfoStyle.Render("   Private keys give full access to your wallet."))
	b.WriteString("\n")
	b.WriteString(InfoStyle.Render("   You can enter it with or without the '0x' prefix."))
	b.WriteString("\n\n")

	// Footer
	b.WriteString(FooterStyle.Render("Tab: Next Field • Enter: Import • Esc: Back"))

	return b.String()
}

// ImportPrivateKeyRequestMsg is sent when the user wants to import a wallet from private key
type ImportPrivateKeyRequestMsg struct {
	Name       string
	PrivateKey string
	Password   string
}