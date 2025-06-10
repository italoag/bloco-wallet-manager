package ui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
)

// ImportPrivateKeyComponent represents the private key import component
type ImportPrivateKeyComponent struct {
	id        string
	form      *huh.Form
	width     int
	height    int
	err       error
	importing bool

	// Form values
	walletName string
	privateKey string
	password   string
}

// NewImportPrivateKeyComponent creates a new private key import component
func NewImportPrivateKeyComponent() ImportPrivateKeyComponent {
	c := ImportPrivateKeyComponent{
		id: "import-private-key",
	}
	c.initForm()
	return c
}

// initForm initializes the huh form
func (c *ImportPrivateKeyComponent) initForm() {
	c.form = huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Key("walletName").
				Title("Wallet Name").
				Placeholder("Enter wallet name...").
				Value(&c.walletName),
			huh.NewInput().
				Key("privateKey").
				Title("Private Key").
				Placeholder("Enter private key (64 hex characters, with or without 0x prefix)").
				EchoMode(huh.EchoModePassword).
				Value(&c.privateKey),
			huh.NewInput().
				Key("password").
				Title("Password").
				Placeholder("Enter password...").
				EchoMode(huh.EchoModePassword).
				Value(&c.password),
		),
	).WithWidth(80).WithShowHelp(false).WithShowErrors(false)
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
	return c.walletName
}

// GetPrivateKey returns the entered private key
func (c *ImportPrivateKeyComponent) GetPrivateKey() string {
	return c.privateKey
}

// GetPassword returns the entered password
func (c *ImportPrivateKeyComponent) GetPassword() string {
	return c.password
}

// Reset clears all inputs
func (c *ImportPrivateKeyComponent) Reset() {
	c.walletName = ""
	c.privateKey = ""
	c.password = ""
	c.err = nil
	c.importing = false
	c.initForm()
}

// Init initializes the component
func (c *ImportPrivateKeyComponent) Init() tea.Cmd {
	// Initialize the form so the first input is focused
	return c.form.Init()
}

// Update handles messages for the import private key component
func (c *ImportPrivateKeyComponent) Update(msg tea.Msg) (*ImportPrivateKeyComponent, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		c.width = msg.Width
		c.height = msg.Height

	case walletCreatedMsg:
		c.Reset()
		return c, func() tea.Msg { return BackToMenuMsg{} }

	case errorMsg:
		c.SetError(fmt.Errorf("%s", string(msg)))
	}

	// Update the form first (allows typing and internal navigation)
	form, cmd := c.form.Update(msg)
	if f, ok := form.(*huh.Form); ok {
		c.form = f
		cmds = append(cmds, cmd)
	}

	// Only handle escape if form didn't handle it (when form is not focused on input)
	if keyMsg, ok := msg.(tea.KeyMsg); ok && keyMsg.String() == "esc" && c.form.State == huh.StateNormal {
		return c, func() tea.Msg { return BackToMenuMsg{} }
	}

	// Check if form is completed
	if c.form.State == huh.StateCompleted {
		if c.validateInputs() {
			c.importing = true
			return c, func() tea.Msg {
				return ImportPrivateKeyRequestMsg{
					Name:       c.GetWalletName(),
					PrivateKey: c.GetPrivateKey(),
					Password:   c.GetPassword(),
				}
			}
		}
		// Reset form state if validation failed
		c.form.State = huh.StateNormal
	}

	return c, tea.Batch(cmds...)
}

// validateInputs checks if the inputs are valid
func (c *ImportPrivateKeyComponent) validateInputs() bool {
	if strings.TrimSpace(c.walletName) == "" {
		c.err = fmt.Errorf("Wallet name cannot be empty")
		return false
	}
	if strings.TrimSpace(c.privateKey) == "" {
		c.err = fmt.Errorf("Private key cannot be empty")
		return false
	}
	if strings.TrimSpace(c.password) == "" {
		c.err = fmt.Errorf("Password cannot be empty")
		return false
	}

	// Basic private key validation
	privateKey := strings.TrimSpace(c.privateKey)
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

	// Header
	headerStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("86")).
		MarginBottom(1)
	b.WriteString(headerStyle.Render("🔑 Import Wallet from Private Key"))
	b.WriteString("\n\n")

	// Form
	b.WriteString(c.form.View())

	// Status messages
	if c.importing {
		b.WriteString("\n")
		b.WriteString(LoadingStyle.Render("⏳ Importing wallet..."))
	} else if c.err != nil {
		b.WriteString("\n")
		b.WriteString(ErrorStyle.Render("❌ Error: " + c.err.Error()))
	}

	// Instructions
	b.WriteString("\n\n")
	b.WriteString(WarningStyle.Render("⚠️  Important: Never share your private key with anyone!"))
	b.WriteString("\n")
	b.WriteString(InfoStyle.Render("   Private keys give full access to your wallet."))
	b.WriteString("\n")
	b.WriteString(InfoStyle.Render("   You can enter it with or without the '0x' prefix."))
	b.WriteString("\n\n")

	// Footer
	b.WriteString(FooterStyle.Render("Tab/Arrow Keys: Navigate • Enter: Import • Esc: Back"))

	return b.String()
}

// ImportPrivateKeyRequestMsg is sent when the user wants to import a wallet from private key
type ImportPrivateKeyRequestMsg struct {
	Name       string
	PrivateKey string
	Password   string
}
