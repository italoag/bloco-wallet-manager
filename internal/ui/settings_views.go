package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// renderSettings renders the main settings menu
func (m Model) renderSettings() string {
	var b strings.Builder

	headerStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("86")).
		MarginBottom(2)

	itemStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("240"))

	selectedStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("86")).
		Background(lipgloss.Color("235"))

	b.WriteString(headerStyle.Render("⚙️  Settings"))
	b.WriteString("\n\n")

	for i, item := range m.settingsItems {
		style := itemStyle
		if i == m.settingsSelected {
			style = selectedStyle
		}

		prefix := "  "
		if i == m.settingsSelected {
			prefix = "▶ "
		}

		b.WriteString(style.Render(prefix + item))
		b.WriteString("\n")
	}

	b.WriteString("\n")
	b.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color("240")).Render("↑/↓: Navigate • Enter: Select • Esc: Back"))

	return b.String()
}

// renderNetworkConfig renders the network configuration view
func (m Model) renderNetworkConfig() string {
	var b strings.Builder

	headerStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("86")).
		MarginBottom(2)

	itemStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("240"))

	selectedStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("86")).
		Background(lipgloss.Color("235"))

	b.WriteString(headerStyle.Render("🌐 Network Configuration"))
	b.WriteString("\n\n")

	if m.editingRPC {
		b.WriteString("Edit RPC Endpoint:\n\n")
		b.WriteString(m.rpcInput.View())
		b.WriteString("\n\n")
		b.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color("240")).Render("Enter: Save • Esc: Cancel"))
	} else {
		for i, item := range m.networkItems {
			style := itemStyle
			if i == m.networkSelected {
				style = selectedStyle
			}

			prefix := "  "
			if i == m.networkSelected {
				prefix = "▶ "
			}

			b.WriteString(style.Render(prefix + item))
			b.WriteString("\n")
		}

		b.WriteString("\n")

		// Show different help text based on selected item
		networkKeys := m.config.GetAllNetworkKeys()
		if m.networkSelected < len(networkKeys) {
			key := networkKeys[m.networkSelected]
			if network, exists := m.config.GetNetworkByKey(key); exists {
				helpText := "↑/↓: Navigate • Enter: Details • A: Activate • E: Edit RPC"
				if network.IsCustom {
					helpText += " • D: Delete"
				}
				helpText += " • Esc: Back"
				b.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color("240")).Render(helpText))
			}
		} else {
			// For "Add Custom Network" and "Back to Settings"
			b.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color("240")).Render("↑/↓: Navigate • Enter: Select • Esc: Back"))
		}
	}

	return b.String()
}

// renderLanguage renders the language selection view
func (m Model) renderLanguage() string {
	var b strings.Builder

	headerStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("86")).
		MarginBottom(2)

	itemStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("240"))

	selectedStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("86")).
		Background(lipgloss.Color("235"))

	b.WriteString(headerStyle.Render("🌍 Language Selection"))
	b.WriteString("\n\n")

	for i, item := range m.languageItems {
		style := itemStyle
		if i == m.languageSelected {
			style = selectedStyle
		}

		prefix := "  "
		if i == m.languageSelected {
			prefix = "▶ "
		}

		b.WriteString(style.Render(prefix + item))
		b.WriteString("\n")
	}

	b.WriteString("\n")
	b.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color("240")).Render("↑/↓: Navigate • Enter: Select • Esc: Back"))

	return b.String()
}

// renderAddNetwork renders the add custom network view
func (m Model) renderAddNetwork() string {
	var b strings.Builder

	headerStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("86")).
		MarginBottom(2)

	inputStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("240")).
		MarginBottom(1)

	focusedStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("86")).
		Background(lipgloss.Color("235")).
		MarginBottom(1)

	b.WriteString(headerStyle.Render("🌐 Add Custom Network"))
	b.WriteString("\n\n")

	// Network Name input
	nameLabel := "Network Name:"
	nameStyle := inputStyle
	if m.addNetworkFocus == 0 {
		nameStyle = focusedStyle
	}
	b.WriteString(nameStyle.Render(nameLabel))
	b.WriteString("\n")
	b.WriteString(m.networkNameInput.View())
	b.WriteString("\n\n")

	// Chain ID input
	chainLabel := "Chain ID:"
	chainStyle := inputStyle
	if m.addNetworkFocus == 1 {
		chainStyle = focusedStyle
	}
	b.WriteString(chainStyle.Render(chainLabel))
	b.WriteString("\n")
	b.WriteString(m.chainIDInput.View())
	b.WriteString("\n\n")

	// RPC Endpoint input
	rpcLabel := "RPC Endpoint:"
	rpcStyle := inputStyle
	if m.addNetworkFocus == 2 {
		rpcStyle = focusedStyle
	}
	b.WriteString(rpcStyle.Render(rpcLabel))
	b.WriteString("\n")
	b.WriteString(m.rpcEndpointInput.View())
	b.WriteString("\n\n")

	if m.addingNetwork {
		b.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color("202")).Render("⏳ Adding network and fetching symbol..."))
		b.WriteString("\n\n")
	}

	if m.err != nil {
		b.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color("196")).Render("❌ Error: " + m.err.Error()))
		b.WriteString("\n\n")
	}

	b.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color("240")).Render("Tab: Next Field • Enter: Add Network • Esc: Cancel"))

	return b.String()
}

// renderNetworkDetails renders detailed view of a selected network
func (m Model) renderNetworkDetails() string {
	var b strings.Builder

	headerStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("86")).
		MarginBottom(2)

	labelStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("86"))

	valueStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("240"))

	if m.selectedNetworkKey == "" {
		return "No network selected"
	}

	network, exists := m.config.GetNetworkByKey(m.selectedNetworkKey)
	if !exists {
		return "Network not found"
	}

	b.WriteString(headerStyle.Render("🌐 Network Details"))
	b.WriteString("\n\n")

	b.WriteString(labelStyle.Render("Name: "))
	b.WriteString(valueStyle.Render(network.Name))
	b.WriteString("\n\n")

	b.WriteString(labelStyle.Render("Chain ID: "))
	b.WriteString(valueStyle.Render(fmt.Sprintf("%d", network.ChainID)))
	b.WriteString("\n\n")

	b.WriteString(labelStyle.Render("Symbol: "))
	b.WriteString(valueStyle.Render(network.Symbol))
	b.WriteString("\n\n")

	b.WriteString(labelStyle.Render("RPC Endpoint: "))
	b.WriteString(valueStyle.Render(network.RPCEndpoint))
	b.WriteString("\n\n")

	b.WriteString(labelStyle.Render("Explorer: "))
	b.WriteString(valueStyle.Render(network.Explorer))
	b.WriteString("\n\n")

	b.WriteString(labelStyle.Render("Status: "))
	status := "Inactive"
	if network.IsActive {
		status = "Active"
	}
	b.WriteString(valueStyle.Render(status))
	b.WriteString("\n\n")

	if network.IsCustom {
		b.WriteString(labelStyle.Render("Type: "))
		b.WriteString(valueStyle.Render("Custom Network"))
		b.WriteString("\n\n")
	}

	b.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color("240")).Render("A: Activate • E: Edit • D: Delete (Custom only) • Esc: Back"))

	return b.String()
}
