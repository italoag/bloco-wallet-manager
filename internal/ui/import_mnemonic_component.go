package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

// ImportMnemonicComponent represents the mnemonic import component
type ImportMnemonicComponent struct {
	id            string
	nameInput     textinput.Model
	mnemonicInput textinput.Model
	passwordInput textinput.Model
	inputFocus    int
	width         int
	height        int
	err           error
	importing     bool
}

// NewImportMnemonicComponent creates a new mnemonic import component
func NewImportMnemonicComponent() ImportMnemonicComponent {
	nameInput := textinput.New()
	nameInput.Placeholder = "Enter wallet name..."
	nameInput.Width = 40
	nameInput.Focus()

	mnemonicInput := textinput.New()
	mnemonicInput.Placeholder = "Enter 12-word mnemonic phrase..."
	mnemonicInput.Width = 60

	passwordInput := textinput.New()
	passwordInput.Placeholder = "Enter password..."
	passwordInput.EchoMode = textinput.EchoPassword
	passwordInput.Width = 40

	return ImportMnemonicComponent{
		id:            "import-mnemonic",
		nameInput:     nameInput,
		mnemonicInput: mnemonicInput,
		passwordInput: passwordInput,
		inputFocus:    0,
	}
}

// SetSize updates the component size
func (c *ImportMnemonicComponent) SetSize(width, height int) {
	c.width = width
	c.height = height
}

// SetError sets an error state
func (c *ImportMnemonicComponent) SetError(err error) {
	c.err = err
	c.importing = false
}

// SetImporting sets the importing state
func (c *ImportMnemonicComponent) SetImporting(importing bool) {
	c.importing = importing
	if importing {
		c.err = nil
	}
}

// GetWalletName returns the entered wallet name
func (c *ImportMnemonicComponent) GetWalletName() string {
	return c.nameInput.Value()
}

// GetMnemonic returns the entered mnemonic
func (c *ImportMnemonicComponent) GetMnemonic() string {
	return c.mnemonicInput.Value()
}

// GetPassword returns the entered password
func (c *ImportMnemonicComponent) GetPassword() string {
	return c.passwordInput.Value()
}

// Reset clears all inputs
func (c *ImportMnemonicComponent) Reset() {
	c.nameInput.SetValue("")
	c.mnemonicInput.SetValue("")
	c.passwordInput.SetValue("")
	c.inputFocus = 0
	c.err = nil
	c.importing = false
	c.nameInput.Focus()
	c.mnemonicInput.Blur()
	c.passwordInput.Blur()
}

// Update handles messages for the import mnemonic component
func (c *ImportMnemonicComponent) Update(msg tea.Msg) (*ImportMnemonicComponent, tea.Cmd) {
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
							return ImportMnemonicRequestMsg{
								Name:     c.nameInput.Value(),
								Mnemonic: c.mnemonicInput.Value(),
								Password: c.passwordInput.Value(),
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
		c.mnemonicInput, cmd = c.mnemonicInput.Update(msg)
		cmds = append(cmds, cmd)
	case 2:
		c.passwordInput, cmd = c.passwordInput.Update(msg)
		cmds = append(cmds, cmd)
	}

	return c, tea.Batch(cmds...)
}

// updateFocus updates which input field has focus
func (c *ImportMnemonicComponent) updateFocus() {
	c.nameInput.Blur()
	c.mnemonicInput.Blur()
	c.passwordInput.Blur()

	switch c.inputFocus {
	case 0:
		c.nameInput.Focus()
	case 1:
		c.mnemonicInput.Focus()
	case 2:
		c.passwordInput.Focus()
	}
}

// validateInputs checks if the inputs are valid
func (c *ImportMnemonicComponent) validateInputs() bool {
	if strings.TrimSpace(c.nameInput.Value()) == "" {
		c.err = fmt.Errorf("Wallet name cannot be empty")
		return false
	}
	if strings.TrimSpace(c.mnemonicInput.Value()) == "" {
		c.err = fmt.Errorf("Mnemonic phrase cannot be empty")
		return false
	}
	if strings.TrimSpace(c.passwordInput.Value()) == "" {
		c.err = fmt.Errorf("Password cannot be empty")
		return false
	}
	
	// Basic mnemonic validation
	words := strings.Fields(strings.TrimSpace(c.mnemonicInput.Value()))
	if len(words) != 12 && len(words) != 24 {
		c.err = fmt.Errorf("Mnemonic must be 12 or 24 words")
		return false
	}
	
	c.err = nil
	return true
}

// View renders the import mnemonic component
func (c *ImportMnemonicComponent) View() string {
	var b strings.Builder

	b.WriteString(HeaderStyle.Render("📥 Import Wallet from Mnemonic"))
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

	// Mnemonic input
	b.WriteString(LabelStyle.Render("Mnemonic Phrase (12 or 24 words):"))
	b.WriteString("\n")
	if c.inputFocus == 1 {
		b.WriteString(FocusedInputStyle.Render(c.mnemonicInput.View()))
	} else {
		b.WriteString(InputStyle.Render(c.mnemonicInput.View()))
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
	b.WriteString(WarningStyle.Render("⚠️  Important: Make sure your mnemonic phrase is correct!"))
	b.WriteString("\n")
	b.WriteString(InfoStyle.Render("   Each word should be separated by a space."))
	b.WriteString("\n\n")

	// Footer
	b.WriteString(FooterStyle.Render("Tab: Next Field • Enter: Import • Esc: Back"))

	return b.String()
}

// ImportMnemonicRequestMsg is sent when the user wants to import a wallet from mnemonic
type ImportMnemonicRequestMsg struct {
	Name     string
	Mnemonic string
	Password string
}