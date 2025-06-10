package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// renderSettings renders the main settings menu
func (m Model) renderSettings() string {
	var contentParts []string

	// Header
	contentParts = append(contentParts, HeaderStyle.Render("⚙️  Settings"))
	contentParts = append(contentParts, "")

	// Menu items
	for i, item := range m.settingsItems {
		prefix := "  "
		if i == m.settingsSelected {
			prefix = "▶ "
			contentParts = append(contentParts, MenuSelectedStyle.Render(prefix+item))
		} else {
			contentParts = append(contentParts, ItemStyle.Render(prefix+item))
		}
	}

	// Create main content
	mainContent := lipgloss.JoinVertical(lipgloss.Left, contentParts...)

	// Footer
	footer := FooterStyle.Render("↑/↓: Navigate • Enter: Select • Esc: Back")

	// Combine all parts with proper spacing
	var finalParts []string
	finalParts = append(finalParts, mainContent)

	// Add spacing before footer
	availableHeight := m.height - lipgloss.Height(mainContent) - 4
	if availableHeight > 0 {
		padding := strings.Repeat("\n", availableHeight)
		finalParts = append(finalParts, padding)
	}

	finalParts = append(finalParts, footer)

	return lipgloss.JoinVertical(lipgloss.Left, finalParts...)
}

// renderNetworkConfig renders the network configuration view
func (m Model) renderNetworkConfig() string {
	var contentParts []string

	// Header
	contentParts = append(contentParts, HeaderStyle.Render("🌐 Network Configuration"))
	contentParts = append(contentParts, "")

	if m.editingRPC {
		contentParts = append(contentParts, LabelStyle.Render("Edit RPC Endpoint:"))
		contentParts = append(contentParts, "")
		contentParts = append(contentParts, m.rpcInput.View())

		// Create main content
		mainContent := lipgloss.JoinVertical(lipgloss.Left, contentParts...)

		// Footer
		footer := FooterStyle.Render("Enter: Save • Esc: Cancel")

		// Combine all parts
		var finalParts []string
		finalParts = append(finalParts, mainContent)

		// Add spacing before footer
		availableHeight := m.height - lipgloss.Height(mainContent) - 4
		if availableHeight > 0 {
			padding := strings.Repeat("\n", availableHeight)
			finalParts = append(finalParts, padding)
		}

		finalParts = append(finalParts, footer)

		return lipgloss.JoinVertical(lipgloss.Left, finalParts...)
	} else {
		// Network items
		for i, item := range m.networkItems {
			prefix := "  "
			if i == m.networkSelected {
				prefix = "▶ "
				contentParts = append(contentParts, MenuSelectedStyle.Render(prefix+item))
			} else {
				contentParts = append(contentParts, ItemStyle.Render(prefix+item))
			}
		}

		// Create main content
		mainContent := lipgloss.JoinVertical(lipgloss.Left, contentParts...)

		// Footer - different help text based on selected item
		var footerText string
		networkKeys := m.config.GetAllNetworkKeys()
		if m.networkSelected < len(networkKeys) {
			key := networkKeys[m.networkSelected]
			if network, exists := m.config.GetNetworkByKey(key); exists {
				footerText = "↑/↓: Navigate • Enter: Details • A: Toggle Active • E: Edit RPC"
				if network.IsCustom {
					footerText += " • D: Delete"
				}
				footerText += " • Esc: Back"
			}
		} else {
			footerText = "↑/↓: Navigate • Enter: Select • Esc: Back"
		}
		footer := FooterStyle.Render(footerText)

		// Combine all parts with proper spacing
		var finalParts []string
		finalParts = append(finalParts, mainContent)

		// Add spacing before footer
		availableHeight := m.height - lipgloss.Height(mainContent) - 4
		if availableHeight > 0 {
			padding := strings.Repeat("\n", availableHeight)
			finalParts = append(finalParts, padding)
		}

		finalParts = append(finalParts, footer)

		return lipgloss.JoinVertical(lipgloss.Left, finalParts...)
	}
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
	var contentParts []string

	// Header
	contentParts = append(contentParts, HeaderStyle.Render("🌐 Add Custom Network"))
	contentParts = append(contentParts, "")

	// Network Name input
	nameLabel := "Network Name:"
	if m.addNetworkFocus == 0 {
		nameLabel = LabelStyle.Render(nameLabel)
	} else {
		nameLabel = InfoStyle.Render(nameLabel)
	}
	contentParts = append(contentParts, nameLabel)
	contentParts = append(contentParts, m.networkNameInput.View())

	// Show network suggestions if available
	if m.showingSuggestions && len(m.networkSuggestions) > 0 && m.addNetworkFocus == 0 {
		contentParts = append(contentParts, "")
		contentParts = append(contentParts, LabelStyle.Render("📡 Suggestions (Press Enter to select):"))

		for i, suggestion := range m.networkSuggestions {
			var style lipgloss.Style
			prefix := "  "
			if i == m.selectedSuggestion {
				style = SelectedSuggestionStyle
				prefix = "▶ "
			} else {
				style = SuggestionStyle
			}

			suggestionText := fmt.Sprintf("%s%s (Chain ID: %d, Symbol: %s)",
				prefix, suggestion.Name, suggestion.ChainID, suggestion.Symbol)
			contentParts = append(contentParts, style.Render(suggestionText))
		}

		contentParts = append(contentParts, "")
		contentParts = append(contentParts, FooterStyle.Render("↑/↓: Navigate suggestions • Enter: Select • Esc: Close"))
	}

	contentParts = append(contentParts, "")

	// Chain ID input
	chainLabel := "Chain ID:"
	if m.addNetworkFocus == 1 {
		chainLabel = LabelStyle.Render(chainLabel)
	} else {
		chainLabel = InfoStyle.Render(chainLabel)
	}
	contentParts = append(contentParts, chainLabel)
	contentParts = append(contentParts, m.chainIDInput.View())
	contentParts = append(contentParts, "")

	// RPC Endpoint input
	rpcLabel := "RPC Endpoint (optional - auto-filled from ChainList):"
	if m.addNetworkFocus == 2 {
		rpcLabel = LabelStyle.Render(rpcLabel)
	} else {
		rpcLabel = InfoStyle.Render(rpcLabel)
	}
	contentParts = append(contentParts, rpcLabel)
	contentParts = append(contentParts, m.rpcEndpointInput.View())

	// Adding network status
	if m.addingNetwork {
		contentParts = append(contentParts, "")
		contentParts = append(contentParts, LoadingStyle.Render("⏳ Adding network and finding best RPC endpoint..."))
	}

	// Error display
	if m.err != nil {
		contentParts = append(contentParts, "")
		contentParts = append(contentParts, ErrorStyle.Render("❌ Error: "+m.err.Error()))
	}

	// Create main content
	mainContent := lipgloss.JoinVertical(lipgloss.Left, contentParts...)

	// Footer instructions
	var footerText string
	if m.showingSuggestions && m.addNetworkFocus == 0 {
		footerText = "↑/↓: Navigate suggestions • Enter: Select • Esc: Close suggestions"
	} else {
		footerText = "Tab: Next Field • Enter: Add Network • Esc: Cancel\n💡 Tip: Type network name for suggestions or enter Chain ID for auto-completion"
	}
	footer := FooterStyle.Render(footerText)

	// Loading status (appears at bottom)
	var loadingStatus string
	if m.isLoading {
		loadingText := m.loadingText
		if loadingText == "" {
			loadingText = "Loading..."
		}
		loadingStatus = LoadingStyle.Render(fmt.Sprintf("%s %s", m.loadingSpinner.View(), loadingText))
	}

	// Combine all parts with proper spacing
	var finalParts []string
	finalParts = append(finalParts, mainContent)

	// Add spacing before footer/loading
	availableHeight := m.height - lipgloss.Height(mainContent) - 4
	if availableHeight > 0 {
		padding := strings.Repeat("\n", availableHeight)
		finalParts = append(finalParts, padding)
	}

	// Add loading status if present
	if loadingStatus != "" {
		finalParts = append(finalParts, loadingStatus)
	}

	// Add footer
	finalParts = append(finalParts, footer)

	return lipgloss.JoinVertical(lipgloss.Left, finalParts...)
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

	b.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color("240")).Render("A: Toggle Active • E: Edit • D: Delete (Custom only) • Esc: Back"))

	return b.String()
}
