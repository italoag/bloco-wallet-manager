package interfaces

import (
	"blocowallet/internal/constants"
	"github.com/charmbracelet/lipgloss"
)

type Styles struct {
	Header          lipgloss.Style
	Content         lipgloss.Style
	Footer          lipgloss.Style
	TopStrip        lipgloss.Style
	MenuItem        lipgloss.Style
	MenuSelected    lipgloss.Style
	SelectedTitle   lipgloss.Style
	MenuTitle       lipgloss.Style
	MenuDesc        lipgloss.Style
	ErrorStyle      lipgloss.Style
	WalletDetails   lipgloss.Style
	StatusBar       lipgloss.Style
	Splash          lipgloss.Style
	StatusBarLeft   lipgloss.Style
	StatusBarCenter lipgloss.Style
	StatusBarRight  lipgloss.Style
}

func createStyles() Styles {
	return Styles{
		Header: lipgloss.NewStyle().
			Align(lipgloss.Left).
			Padding(1, 2),

		Content: lipgloss.NewStyle().
			Align(lipgloss.Left).
			Padding(1, 2),

		Footer: lipgloss.NewStyle().
			Align(lipgloss.Left).
			PaddingLeft(1).
			PaddingRight(1).
			Background(lipgloss.Color("#7D56F4")),

		TopStrip: lipgloss.NewStyle().Margin(1, constants.StyleMargin).Padding(0, constants.StyleMargin),
		MenuItem: lipgloss.NewStyle().Width(constants.StyleWidth).Margin(0, constants.StyleMargin).Padding(0, constants.StyleMargin),
		MenuSelected: lipgloss.NewStyle().
			Foreground(lipgloss.Color("99")).
			Bold(true).
			Margin(0, constants.StyleMargin).
			Padding(0, constants.StyleMargin).
			Border(lipgloss.NormalBorder(), false, false, false, true).
			Width(constants.StyleWidth),
		SelectedTitle: lipgloss.NewStyle().Bold(true).
			Margin(0, constants.StyleMargin).
			Padding(0, constants.StyleMargin).
			Foreground(lipgloss.Color("99")),
		MenuTitle: lipgloss.NewStyle().
			Margin(0, constants.StyleMargin).
			Padding(0, constants.StyleMargin).
			Bold(true),
		MenuDesc: lipgloss.NewStyle().
			Margin(0, constants.StyleMargin).
			Padding(0, constants.StyleMargin).
			Width(constants.StyleWidth).
			Foreground(lipgloss.Color("244")),
		ErrorStyle: lipgloss.NewStyle().
			Padding(1, 2).
			Margin(1, constants.StyleMargin),
		WalletDetails: lipgloss.NewStyle().
			Margin(1, constants.StyleMargin).
			Padding(1, 2),
		StatusBar: lipgloss.NewStyle().
			Foreground(lipgloss.Color("241")).
			Padding(0, constants.StyleMargin),
		Splash: lipgloss.NewStyle().
			Align(lipgloss.Center).Padding(1, 2),
		StatusBarLeft: lipgloss.NewStyle().
			Background(lipgloss.Color("#7D56F4")).
			PaddingLeft(1).
			PaddingRight(1),
		StatusBarCenter: lipgloss.NewStyle().
			Background(lipgloss.Color("#454544")).
			PaddingLeft(1).
			PaddingRight(1),
		StatusBarRight: lipgloss.NewStyle().
			Background(lipgloss.Color("#CC5C87")).
			PaddingLeft(1).
			PaddingRight(1),
	}
}
