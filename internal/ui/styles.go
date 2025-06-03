package ui

import "github.com/charmbracelet/lipgloss"

// Style definitions for consistent UI appearance
var (
	// Header styles
	HeaderStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("86")).
			MarginBottom(2)

	// Menu and selection styles
	SelectedStyle = lipgloss.NewStyle().
			Background(lipgloss.Color("86")).
			Foreground(lipgloss.Color("232")).
			Padding(0, 1)

	MenuSelectedStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(lipgloss.Color("86")).
				Background(lipgloss.Color("235"))

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
			Foreground(lipgloss.Color("86")).
			Align(lipgloss.Center)

	SubtitleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("241")).
			Align(lipgloss.Center)

	InstructionStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("86")).
				Align(lipgloss.Center)

	// Balance styles
	BalanceStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("86"))

	NetworkStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("240"))
)
