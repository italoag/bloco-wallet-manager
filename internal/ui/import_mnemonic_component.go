package ui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
)

// WizardStep represents the current step of the wizard
type WizardStep int

const (
	StepBasicInfo WizardStep = iota + 1
	StepMnemonicWords
)

// ImportMnemonicComponent represents the mnemonic import component
type ImportMnemonicComponent struct {
	id        string
	form      *huh.Form
	width     int
	height    int
	err       error
	importing bool

	// Wizard state
	currentStep WizardStep

	// Form values
	walletName                                  string
	password                                    string
	word1, word2, word3, word4, word5, word6    string
	word7, word8, word9, word10, word11, word12 string
}

// NewImportMnemonicComponent creates a new mnemonic import component
func NewImportMnemonicComponent() ImportMnemonicComponent {
	c := ImportMnemonicComponent{
		id:          "import-mnemonic",
		currentStep: StepBasicInfo,
	}
	c.initBasicInfoForm()
	return c
}

// initBasicInfoForm initializes the first step form (wallet name and password)
func (c *ImportMnemonicComponent) initBasicInfoForm() {
	// Reset values
	c.walletName = ""
	c.password = ""

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
				EchoMode(huh.EchoModePassword).
				Value(&c.password).
				Validate(func(s string) error {
					if strings.TrimSpace(s) == "" {
						return fmt.Errorf("password cannot be empty")
					}
					return nil
				}),
		),
	).WithWidth(60).WithShowHelp(false).WithShowErrors(false)
}

// initMnemonicForm initializes the second step form (12 mnemonic words)
func (c *ImportMnemonicComponent) initMnemonicForm() {
	// Reset mnemonic words
	c.word1, c.word2, c.word3, c.word4, c.word5, c.word6 = "", "", "", "", "", ""
	c.word7, c.word8, c.word9, c.word10, c.word11, c.word12 = "", "", "", "", "", ""

	c.form = huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Key("word1").
				Title("Word 1").
				Placeholder("1st word").
				Value(&c.word1),
			huh.NewInput().
				Key("word2").
				Title("Word 2").
				Placeholder("2nd word").
				Value(&c.word2),
			huh.NewInput().
				Key("word3").
				Title("Word 3").
				Placeholder("3rd word").
				Value(&c.word3),
			huh.NewInput().
				Key("word4").
				Title("Word 4").
				Placeholder("4th word").
				Value(&c.word4),
			huh.NewInput().
				Key("word5").
				Title("Word 5").
				Placeholder("5th word").
				Value(&c.word5),
			huh.NewInput().
				Key("word6").
				Title("Word 6").
				Placeholder("6th word").
				Value(&c.word6),
		),
		huh.NewGroup(
			huh.NewInput().
				Key("word7").
				Title("Word 7").
				Placeholder("7th word").
				Value(&c.word7),
			huh.NewInput().
				Key("word8").
				Title("Word 8").
				Placeholder("8th word").
				Value(&c.word8),
			huh.NewInput().
				Key("word9").
				Title("Word 9").
				Placeholder("9th word").
				Value(&c.word9),
			huh.NewInput().
				Key("word10").
				Title("Word 10").
				Placeholder("10th word").
				Value(&c.word10),
			huh.NewInput().
				Key("word11").
				Title("Word 11").
				Placeholder("11th word").
				Value(&c.word11),
			huh.NewInput().
				Key("word12").
				Title("Word 12").
				Placeholder("12th word").
				Value(&c.word12),
		),
	).WithWidth(80).WithShowHelp(false).WithShowErrors(false).WithLayout(huh.LayoutColumns(2))
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
	if c.form != nil {
		return strings.TrimSpace(c.form.GetString("walletName"))
	}
	return strings.TrimSpace(c.walletName)
}

// GetMnemonic returns the entered mnemonic
func (c *ImportMnemonicComponent) GetMnemonic() string {
	if c.form != nil {
		// Try to get values from form first
		words := []string{
			c.form.GetString("word1"), c.form.GetString("word2"), c.form.GetString("word3"),
			c.form.GetString("word4"), c.form.GetString("word5"), c.form.GetString("word6"),
			c.form.GetString("word7"), c.form.GetString("word8"), c.form.GetString("word9"),
			c.form.GetString("word10"), c.form.GetString("word11"), c.form.GetString("word12"),
		}

		// Check if any form values are present
		hasFormValues := false
		for _, word := range words {
			if strings.TrimSpace(word) != "" {
				hasFormValues = true
				break
			}
		}

		// If form has values, use them; otherwise fall back to component variables
		if hasFormValues {
			return strings.Join(words, " ")
		}
	}

	// Fallback to component variables (for tests and when form is empty)
	words := []string{
		c.word1, c.word2, c.word3, c.word4, c.word5, c.word6,
		c.word7, c.word8, c.word9, c.word10, c.word11, c.word12,
	}
	return strings.Join(words, " ")
}

// GetPassword returns the entered password
func (c *ImportMnemonicComponent) GetPassword() string {
	if c.form != nil {
		return strings.TrimSpace(c.form.GetString("password"))
	}
	return strings.TrimSpace(c.password)
}

// Reset clears all inputs and resets to first step
func (c *ImportMnemonicComponent) Reset() {
	c.walletName = ""
	c.word1, c.word2, c.word3, c.word4, c.word5, c.word6 = "", "", "", "", "", ""
	c.word7, c.word8, c.word9, c.word10, c.word11, c.word12 = "", "", "", "", "", ""
	c.password = ""
	c.err = nil
	c.importing = false
	c.currentStep = StepBasicInfo
	c.initBasicInfoForm()
}

// Init initializes the component
func (c *ImportMnemonicComponent) Init() tea.Cmd {
	// Initialize the form so input fields are ready
	return c.form.Init()
}

// Update handles messages for the import mnemonic component
func (c *ImportMnemonicComponent) Update(msg tea.Msg) (*ImportMnemonicComponent, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		c.width = msg.Width
		c.height = msg.Height

	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			if c.currentStep == StepMnemonicWords {
				// Go back to previous step
				c.currentStep = StepBasicInfo
				c.initBasicInfoForm()
				return c, c.form.Init()
			} else {
				// Exit to menu
				return c, func() tea.Msg { return BackToMenuMsg{} }
			}
		}

	case walletCreatedMsg:
		c.Reset()
		return c, func() tea.Msg { return BackToMenuMsg{} }

	case errorMsg:
		c.SetError(fmt.Errorf("%s", string(msg)))
	}

	// Process the form
	form, cmd := c.form.Update(msg)
	if f, ok := form.(*huh.Form); ok {
		c.form = f
		cmds = append(cmds, cmd)
	}

	// Check if current step is completed
	if c.form.State == huh.StateCompleted && !c.importing {
		switch c.currentStep {
		case StepBasicInfo:
			// Store values from basic info step (try form first, fallback to variables)
			formWalletName := strings.TrimSpace(c.form.GetString("walletName"))
			formPassword := strings.TrimSpace(c.form.GetString("password"))

			if formWalletName != "" {
				c.walletName = formWalletName
			}
			if formPassword != "" {
				c.password = formPassword
			}

			// Validate basic info
			if strings.TrimSpace(c.walletName) == "" || strings.TrimSpace(c.password) == "" {
				c.SetError(fmt.Errorf("wallet name and password are required"))
				return c, nil
			}

			// First step completed, move to mnemonic step
			c.currentStep = StepMnemonicWords
			c.initMnemonicForm()
			return c, c.form.Init()

		case StepMnemonicWords:
			// All steps completed, proceed with import
			// Get mnemonic words from form (try form first, fallback to variables)
			formWords := []string{
				strings.TrimSpace(c.form.GetString("word1")), strings.TrimSpace(c.form.GetString("word2")), strings.TrimSpace(c.form.GetString("word3")),
				strings.TrimSpace(c.form.GetString("word4")), strings.TrimSpace(c.form.GetString("word5")), strings.TrimSpace(c.form.GetString("word6")),
				strings.TrimSpace(c.form.GetString("word7")), strings.TrimSpace(c.form.GetString("word8")), strings.TrimSpace(c.form.GetString("word9")),
				strings.TrimSpace(c.form.GetString("word10")), strings.TrimSpace(c.form.GetString("word11")), strings.TrimSpace(c.form.GetString("word12")),
			}

			// Check if form has values, otherwise use component variables
			hasFormWords := false
			for _, word := range formWords {
				if word != "" {
					hasFormWords = true
					break
				}
			}

			var words []string
			if hasFormWords {
				words = formWords
			} else {
				// Fallback to component variables (for tests)
				words = []string{
					strings.TrimSpace(c.word1), strings.TrimSpace(c.word2), strings.TrimSpace(c.word3),
					strings.TrimSpace(c.word4), strings.TrimSpace(c.word5), strings.TrimSpace(c.word6),
					strings.TrimSpace(c.word7), strings.TrimSpace(c.word8), strings.TrimSpace(c.word9),
					strings.TrimSpace(c.word10), strings.TrimSpace(c.word11), strings.TrimSpace(c.word12),
				}
			}

			// Filter out empty words and join
			var validWords []string
			for _, word := range words {
				if word != "" {
					validWords = append(validWords, word)
				}
			}

			if len(validWords) != 12 {
				c.SetError(fmt.Errorf("mnemonic must contain exactly 12 words"))
				return c, nil
			}

			mnemonic := strings.Join(validWords, " ")

			c.importing = true
			return c, func() tea.Msg {
				return ImportMnemonicRequestMsg{
					Name:     c.walletName,
					Mnemonic: mnemonic,
					Password: c.password,
				}
			}
		}
	}

	return c, tea.Batch(cmds...)
}

// View renders the import mnemonic component
func (c *ImportMnemonicComponent) View() string {
	var b strings.Builder

	// Header with step indicator
	headerStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("86")).
		MarginBottom(1)

	stepIndicator := ""
	switch c.currentStep {
	case StepBasicInfo:
		stepIndicator = " (Step 1/2: Basic Information)"
	case StepMnemonicWords:
		stepIndicator = " (Step 2/2: Mnemonic Phrase)"
	}

	b.WriteString(headerStyle.Render("📥 Import Wallet from Mnemonic" + stepIndicator))
	b.WriteString("\n\n")

	// Show wallet name and password from previous step if on step 2
	if c.currentStep == StepMnemonicWords && c.walletName != "" {
		infoStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("245")).
			MarginBottom(1)
		b.WriteString(infoStyle.Render(fmt.Sprintf("Wallet Name: %s", c.walletName)))
		b.WriteString("\n")
	}

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

	// Step-specific instructions
	b.WriteString("\n")
	switch c.currentStep {
	case StepBasicInfo:
		b.WriteString(InfoStyle.Render("💡 First, enter your wallet name and password."))
		b.WriteString("\n")
		b.WriteString(InfoStyle.Render("   These will be used to secure your imported wallet."))
	case StepMnemonicWords:
		b.WriteString(WarningStyle.Render("⚠️  Important: Enter your 12-word mnemonic phrase!"))
		b.WriteString("\n")
		b.WriteString(InfoStyle.Render("   Each word should be entered in the correct order (1-12)."))
		b.WriteString("\n")
		b.WriteString(InfoStyle.Render("   The words are arranged in two columns for easier input."))
	}

	b.WriteString("\n\n")

	// Footer with navigation instructions
	var footerText string
	switch c.currentStep {
	case StepBasicInfo:
		footerText = "Tab/Arrow Keys: Navigate • Enter: Next Step • Esc: Back to Menu"
	case StepMnemonicWords:
		footerText = "Tab/Arrow Keys: Navigate • Enter: Import Wallet • Esc: Previous Step"
	}
	b.WriteString(FooterStyle.Render(footerText))

	return b.String()
}

// ImportMnemonicRequestMsg is sent when the user wants to import a wallet from mnemonic
type ImportMnemonicRequestMsg struct {
	Name     string
	Mnemonic string
	Password string
}
