package ui

import (
	"fmt"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/digitallyserviced/tdfgo/tdf"
)

// SplashComponent represents the splash screen component
type SplashComponent struct {
	id           string
	width        int
	height       int
	selectedFont *tdf.TheDrawFont
	fontName     string
	progress     int
	maxProgress  int
	loading      bool
	showWelcome  bool
}

// NewSplashComponent creates a new splash component
func NewSplashComponent() SplashComponent {
	return SplashComponent{
		id:          "splash-screen",
		maxProgress: 100,
		loading:     true,
	}
}

// SetSize updates the component size
func (c *SplashComponent) SetSize(width, height int) {
	c.width = width
	c.height = height
}

// SetFont sets the ASCII art font
func (c *SplashComponent) SetFont(font *tdf.TheDrawFont, fontName string) {
	c.selectedFont = font
	c.fontName = fontName
}

// StartLoading starts the loading animation
func (c *SplashComponent) StartLoading() tea.Cmd {
	c.loading = true
	c.progress = 0
	return c.tickCmd()
}

// StopLoading stops the loading and shows welcome message
func (c *SplashComponent) StopLoading() {
	c.loading = false
	c.progress = c.maxProgress
	c.showWelcome = true
}

// Update handles messages for the splash component
func (c *SplashComponent) Update(msg tea.Msg) (*SplashComponent, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		c.width = msg.Width
		c.height = msg.Height

	case tea.KeyMsg:
		if !c.loading {
			// Any key press after loading completes will exit splash
			return c, func() tea.Msg { return SplashCompletedMsg{} }
		}

	case splashTickMsg:
		if c.loading && c.progress < c.maxProgress {
			c.progress += 2
			if c.progress >= c.maxProgress {
				c.StopLoading()
				return c, tea.Tick(time.Second*2, func(t time.Time) tea.Msg {
					return SplashCompletedMsg{}
				})
			}
			return c, c.tickCmd()
		}
	}

	return c, nil
}

// tickCmd creates a tick command for the loading animation
func (c *SplashComponent) tickCmd() tea.Cmd {
	return tea.Tick(time.Millisecond*50, func(t time.Time) tea.Msg {
		return splashTickMsg(t)
	})
}

// View renders the splash component
func (c *SplashComponent) View() string {
	var b strings.Builder

	// Calculate vertical centering
	logoHeight := 8                // Approximate height of ASCII logo
	totalHeight := logoHeight + 10 // Logo + spacing + other elements
	topPadding := (c.height - totalHeight) / 2
	if topPadding < 0 {
		topPadding = 0
	}

	// Add top padding
	for i := 0; i < topPadding; i++ {
		b.WriteString("\n")
	}

	// ASCII Art Logo
	var logo string
	if c.selectedFont != nil {
		// Initialize string renderer for the selected font
		fontString := tdf.NewTheDrawFontStringFont(c.selectedFont)

		// Render the logo using TDF font
		renderedLogo := fontString.RenderString("BlockoWallet")
		if renderedLogo != "" {
			logo = strings.TrimSpace(renderedLogo)
		} else {
			logo = "BlockoWallet"
		}
	} else {
		logo = `
 ██████╗ ██╗      ██████╗  ██████╗██╗  ██╗ ██████╗ ██╗    ██╗ █████╗ ██╗     ██╗     ███████╗████████╗
 ██╔══██╗██║     ██╔═══██╗██╔════╝██║ ██╔╝██╔═══██╗██║    ██║██╔══██╗██║     ██║     ██╔════╝╚══██╔══╝
 ██████╔╝██║     ██║   ██║██║     █████╔╝ ██║   ██║██║ █╗ ██║███████║██║     ██║     █████╗     ██║   
 ██╔══██╗██║     ██║   ██║██║     ██╔═██╗ ██║   ██║██║███╗██║██╔══██║██║     ██║     ██╔══╝     ██║   
 ██████╔╝███████╗╚██████╔╝╚██████╗██║  ██╗╚██████╔╝╚███╔███╔╝██║  ██║███████╗███████╗███████╗   ██║   
 ╚═════╝ ╚══════╝ ╚═════╝  ╚═════╝╚═╝  ╚═╝ ╚═════╝  ╚══╝╚══╝ ╚═╝  ╚═╝╚══════╝╚══════╝╚══════╝   ╚═╝   
		`
	}

	b.WriteString(LogoStyle.Render(logo))
	b.WriteString("\n\n")

	// Subtitle
	b.WriteString(SubtitleStyle.Render("Secure Multi-Network Cryptocurrency Wallet"))
	b.WriteString("\n\n")

	// Loading bar or welcome message
	if c.loading {
		progressBar := c.createProgressBar()
		b.WriteString(lipgloss.NewStyle().Align(lipgloss.Center).Render(progressBar))
		b.WriteString("\n")
		b.WriteString(InfoStyle.Render(fmt.Sprintf("Loading... %d%%", c.progress)))
	} else if c.showWelcome {
		b.WriteString(InstructionStyle.Render("🎉 Welcome to BlockoWallet!"))
		b.WriteString("\n")
		b.WriteString(InfoStyle.Render("Press any key to continue..."))
	}

	b.WriteString("\n\n")

	// Project info
	projectInfo := lipgloss.JoinVertical(
		lipgloss.Center,
		"Version 1.0.0",
		"Built with ❤️  by the BlockoWallet Team",
		"Secure • Fast • Multi-Chain",
	)
	b.WriteString(ProjectInfoStyle.Render(projectInfo))

	return b.String()
}

// createProgressBar creates a visual progress bar
func (c *SplashComponent) createProgressBar() string {
	width := 50
	filled := int(float64(width) * float64(c.progress) / float64(c.maxProgress))
	empty := width - filled

	bar := strings.Repeat("█", filled) + strings.Repeat("░", empty)

	return lipgloss.NewStyle().
		Foreground(lipgloss.Color("86")).
		Render("▕" + bar + "▏")
}

// splashTickMsg represents a tick in the loading animation
type splashTickMsg time.Time

// SplashCompletedMsg is sent when the splash screen should close
type SplashCompletedMsg struct{}
