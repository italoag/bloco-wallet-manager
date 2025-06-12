package ui

import "github.com/charmbracelet/lipgloss"

// Style definitions for consistent UI appearance
var (
	// Menu item colors
	menuItemForeground = lipgloss.Color("#874BFD")
	menuItemBackground = lipgloss.Color("#FFFFFF")

	// Header styles
	HeaderStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(menuItemForeground).
			MarginBottom(2)

		// Menu and selection styles
	SelectedStyle = lipgloss.NewStyle().
			Background(menuItemBackground).
			Foreground(menuItemForeground).
			Padding(0, 1)

	MenuSelectedStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(menuItemForeground).
				Background(menuItemBackground)

		// Item styles for lists
	ItemStyle = lipgloss.NewStyle().
			Foreground(menuItemForeground).
			Background(menuItemBackground)

	// Text styles
	InfoStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("240"))

	ErrorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("196"))

	FooterStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("241"))

	SensitiveStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("208")).
			Bold(true)

	// Splash screen styles
	LogoStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#874BFD")).
			Align(lipgloss.Center)

	SubtitleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("241")).
			Align(lipgloss.Center)

	InstructionStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#874BFD")).
				Align(lipgloss.Center)

	ProjectInfoStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("241")).
				Align(lipgloss.Center)

	// Balance styles
	BalanceStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#874BFD"))

	NetworkStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("240"))

	// Dialog styles
	DialogBoxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder(), true).
			BorderForeground(lipgloss.Color("#874BFD")).
			Padding(3, 3)

	DialogQuestionStyle = lipgloss.NewStyle().
				Width(45).
				Align(lipgloss.Center)

	WalletNameStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#ff6b6b")).
			Bold(true)

	AddressStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#8a8a8a"))

	WarningStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#666666"))

	// Button styles
	ButtonStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFF")).
			Background(lipgloss.Color("#666")).
			Padding(0, 3).
			Margin(0, 1)

	ActiveButtonStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#FFF")).
				Background(lipgloss.Color("#874BFD")).
				Padding(0, 3).
				Margin(0, 1).
				Bold(true)

	// Input styles
	InputStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#874BFD")).
			Background(lipgloss.Color("235")).
			Padding(0, 1).
			Margin(1, 0)

	FocusedInputStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#874BFD")).
				Background(lipgloss.Color("235")).
				Padding(0, 1).
				Margin(1, 0).
				Border(lipgloss.RoundedBorder()).
				BorderForeground(lipgloss.Color("#874BFD"))

	// Suggestion styles
	SuggestionStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("250")).
			Background(lipgloss.Color("236")).
			Padding(0, 1).
			Margin(0, 1)

	SelectedSuggestionStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#874BFD")).
				Background(lipgloss.Color("240")).
				Padding(0, 1).
				Margin(0, 1).
				Bold(true)

	// Status styles
	LoadingStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("202"))

	SuccessStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("2"))

	// Network configuration styles
	LabelStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#874BFD")).
			Bold(true)

	ValueStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("250"))

	// Active/Inactive status styles
	ActiveStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("2")).
			Bold(true)

	InactiveStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("240"))

	CustomTagStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("208")).
			Bold(true)
)
