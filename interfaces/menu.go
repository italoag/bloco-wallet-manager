package interfaces

import (
	"blocowallet/constants"
	"blocowallet/localization"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/lipgloss"
)

func NewMenu() list.Model {
	menuItems := []list.Item{
		menuItem{title: localization.Labels["create_new_wallet"], description: localization.Labels["create_new_wallet_desc"]},
		menuItem{title: localization.Labels["import_wallet"], description: localization.Labels["import_wallet_desc"]},
		menuItem{title: localization.Labels["list_wallets"], description: localization.Labels["list_wallets_desc"]},
		menuItem{title: localization.Labels["exit"], description: localization.Labels["exit_desc"]},
	}

	delegate := list.NewDefaultDelegate()
	delegate.ShowDescription = true
	delegate.Styles.SelectedTitle = delegate.Styles.SelectedTitle.Foreground(lipgloss.Color("#00FF00")).Bold(true)
	delegate.Styles.SelectedDesc = delegate.Styles.SelectedDesc.Foreground(lipgloss.Color("#00FF00"))
	delegate.Styles.NormalTitle = delegate.Styles.NormalTitle.Foreground(lipgloss.Color("#FFFFFF"))
	delegate.Styles.NormalDesc = delegate.Styles.NormalDesc.Foreground(lipgloss.Color("#888888"))

	menuList := list.New(menuItems, delegate, constants.MenuWidth, 0)
	menuList.Title = localization.Labels["main_menu_title"]
	menuList.SetShowStatusBar(false)
	menuList.SetFilteringEnabled(false)
	menuList.SetShowHelp(false)
	menuList.Styles.Title = lipgloss.NewStyle().
		Background(lipgloss.Color("#25A065")).
		Foreground(lipgloss.Color("#FFFDF5")).
		Padding(0, 1).
		Bold(true)

	return menuList
}
