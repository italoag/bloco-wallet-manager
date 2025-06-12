package ui

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"blocowallet/internal/blockchain"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// AddNetworkComponent represents the add network component
type AddNetworkComponent struct {
	id     string
	width  int
	height int
	err    error
	adding bool

	// Text input fields
	searchInput      textinput.Model
	chainIDInput     textinput.Model
	rpcEndpointInput textinput.Model
	symbolInput      textinput.Model
	nameInput        textinput.Model

	// Form state
	focusIndex         int
	inputs             []textinput.Model
	selectedSuggestion int
	isSearchFocused    bool

	// Chain service for suggestions
	chainListService *blockchain.ChainListService

	// Autocomplete data
	suggestions        []blockchain.NetworkSuggestion
	suggestionList     list.Model // Novo: lista interativa de sugestões
	loadingSuggestions bool
	lastSearchTerm     string
	typingDebounce     time.Time
}

// NewAddNetworkComponent creates a new add network component
func NewAddNetworkComponent() AddNetworkComponent {
	c := AddNetworkComponent{
		id:               "add-network",
		chainListService: blockchain.NewChainListService(),
	}
	c.initInputs()
	return c
}

// initInputs initializes the text input fields
func (c *AddNetworkComponent) initInputs() {
	// Search input for network search
	c.searchInput = textinput.New()
	c.searchInput.Placeholder = "Type to search networks (e.g., Polygon, Binance)..."
	c.searchInput.Width = 60
	c.searchInput.ShowSuggestions = true
	c.searchInput.Focus()
	c.isSearchFocused = true

	// Network name input for display
	c.nameInput = textinput.New()
	c.nameInput.Placeholder = "Network name will be filled automatically..."
	c.nameInput.Width = 60

	// Chain ID input
	c.chainIDInput = textinput.New()
	c.chainIDInput.Placeholder = "Chain ID will be filled automatically..."
	c.chainIDInput.Width = 60

	// Symbol input
	c.symbolInput = textinput.New()
	c.symbolInput.Placeholder = "Symbol will be filled automatically..."
	c.symbolInput.Width = 60

	// RPC endpoint input
	c.rpcEndpointInput = textinput.New()
	c.rpcEndpointInput.Placeholder = "RPC URL will be filled automatically..."
	c.rpcEndpointInput.Width = 60

	// Initialize inputs slice for easy navigation
	c.inputs = []textinput.Model{
		c.searchInput,
		c.nameInput,
		c.chainIDInput,
		c.symbolInput,
		c.rpcEndpointInput,
	}

	// Inicializa a lista de sugestões
	delegate := list.NewDefaultDelegate()
	delegate.Styles.NormalTitle = delegate.Styles.NormalTitle.
		Foreground(menuItemForeground)
	delegate.Styles.NormalDesc = delegate.Styles.NormalDesc.
		Foreground(menuItemDescriptionForeground)
	delegate.Styles.SelectedTitle = delegate.Styles.SelectedTitle.
		Foreground(menuItemForeground)
	delegate.Styles.SelectedDesc = delegate.Styles.SelectedDesc.
		Foreground(menuItemDescriptionForeground)
	c.suggestionList = list.New([]list.Item{}, delegate, 60, 5)
	c.suggestionList.SetShowStatusBar(false)
	c.suggestionList.SetShowHelp(false)
	c.suggestionList.SetFilteringEnabled(false)
	c.suggestionList.Title = "Suggestions"

	// Initialize other fields
	c.selectedSuggestion = -1
	c.typingDebounce = time.Time{}
}

// SetSize updates the component size
func (c *AddNetworkComponent) SetSize(width, height int) {
	c.width = width
	c.height = height
}

// SetError sets an error state
func (c *AddNetworkComponent) SetError(err error) {
	c.err = err
	c.adding = false
}

// SetAdding sets the adding state
func (c *AddNetworkComponent) SetAdding(adding bool) {
	c.adding = adding
	if adding {
		c.err = nil
	}
}

// GetNetworkName returns the entered network name
func (c *AddNetworkComponent) GetNetworkName() string {
	return c.nameInput.Value()
}

// GetChainID returns the entered chain ID as integer
func (c *AddNetworkComponent) GetChainID() (int, error) {
	chainID, err := strconv.Atoi(strings.TrimSpace(c.chainIDInput.Value()))
	if err != nil {
		return 0, fmt.Errorf("invalid chain ID: must be a number")
	}
	return chainID, nil
}

// GetSymbol returns the entered symbol
func (c *AddNetworkComponent) GetSymbol() string {
	return c.symbolInput.Value()
}

// GetRPCEndpoint returns the entered RPC endpoint
func (c *AddNetworkComponent) GetRPCEndpoint() string {
	return c.rpcEndpointInput.Value()
}

// Reset clears all inputs
func (c *AddNetworkComponent) Reset() {
	c.searchInput.SetValue("")
	c.nameInput.SetValue("")
	c.chainIDInput.SetValue("")
	c.symbolInput.SetValue("")
	c.rpcEndpointInput.SetValue("")
	c.err = nil
	c.adding = false
	c.suggestions = []blockchain.NetworkSuggestion{}
	c.loadingSuggestions = false
	c.focusIndex = 0
	c.selectedSuggestion = -1
	c.isSearchFocused = true
	c.initInputs()
}

// searchNetworks searches for networks based on the query
func (c *AddNetworkComponent) searchNetworks(query string) tea.Cmd {
	return func() tea.Msg {
		query = strings.TrimSpace(query)

		// If empty query, return popular networks
		if query == "" {
			popular := []blockchain.NetworkSuggestion{
				{ChainID: 137, Name: "Polygon Mainnet", Symbol: "POL"},
				{ChainID: 42161, Name: "Arbitrum One", Symbol: "ETH"},
				{ChainID: 10, Name: "Optimism", Symbol: "ETH"},
				{ChainID: 8453, Name: "Base", Symbol: "ETH"},
			}
			return networkSuggestionsMsg(popular)
		}

		suggestions, err := c.chainListService.SearchNetworksByName(query)
		if err != nil {
			return errorMsg(err.Error())
		}

		return networkSuggestionsMsg(suggestions)
	}
}

// fillNetworkData fills the form with network data when a suggestion is selected
func (c *AddNetworkComponent) fillNetworkData(suggestion blockchain.NetworkSuggestion) {
	// Find the full chain info for this suggestion
	_, rpcURL, err := c.chainListService.GetChainInfoWithRetry(suggestion.ChainID)
	if err != nil {
		c.err = fmt.Errorf("failed to get network details: %v", err)
		return
	}

	// Update input values directly
	c.nameInput.SetValue(suggestion.Name)
	c.chainIDInput.SetValue(strconv.Itoa(suggestion.ChainID))
	c.symbolInput.SetValue(suggestion.Symbol)
	c.rpcEndpointInput.SetValue(rpcURL)

	// Update search input with the selected name
	c.searchInput.SetValue(suggestion.Name)

	// Clear selection highlighting
	c.selectedSuggestion = -1

	// Move focus to the network name field for possible editing
	c.focusIndex = 1
	c.updateFocus()
}

// Init initializes the component
func (c *AddNetworkComponent) Init() tea.Cmd {
	// Initialize the search input to be focused
	c.focusIndex = 0
	c.searchInput.Focus()
	c.isSearchFocused = true
	c.selectedSuggestion = -1

	// Start with some popular networks
	return c.searchNetworks("")
}

// Update handles messages for the add network component
func (c *AddNetworkComponent) Update(msg tea.Msg) (*AddNetworkComponent, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		c.width = msg.Width
		c.height = msg.Height
		c.suggestionList.SetSize(60, 5)

	case networkAddedMsg:
		c.Reset()
		return c, func() tea.Msg { return BackToNetworkListMsg{} }

	case networkSuggestionsMsg:
		c.suggestions = []blockchain.NetworkSuggestion(msg)
		c.loadingSuggestions = false
		c.selectedSuggestion = -1
		// Atualiza a lista de sugestões
		items := make([]list.Item, 0, len(c.suggestions))
		for _, s := range c.suggestions {
			items = append(items, networkSuggestionItem{suggestion: s})
		}
		c.suggestionList.SetItems(items)
		c.suggestionList.Select(0)

	case errorMsg:
		c.SetError(fmt.Errorf("%s", string(msg)))
		c.loadingSuggestions = false

	case tea.KeyMsg:
		if c.isSearchFocused && len(c.suggestions) > 0 {
			switch msg.String() {
			case "up", "down":
				var cmd tea.Cmd
				c.suggestionList, cmd = c.suggestionList.Update(msg)
				cmds = append(cmds, cmd)
				c.selectedSuggestion = c.suggestionList.Index()
				return c, tea.Batch(cmds...)
			case "enter", "tab":
				if len(c.suggestionList.Items()) > 0 && c.selectedSuggestion >= 0 {
					item := c.suggestionList.SelectedItem().(networkSuggestionItem)
					c.fillNetworkData(item.suggestion)
					return c, nil
				}
			}
		}

		// Handle global special keys
		switch msg.String() {
		case "esc":
			return c, func() tea.Msg { return BackToNetworkListMsg{} }

		case "enter":
			// Submit form if not in search mode
			if !c.isSearchFocused && c.validateInputs() {
				c.adding = true
				return c, func() tea.Msg {
					return AddNetworkRequestMsg{
						Name:        c.GetNetworkName(),
						ChainID:     c.chainIDInput.Value(),
						Symbol:      c.GetSymbol(),
						RPCEndpoint: c.GetRPCEndpoint(),
					}
				}
			}

		case "tab":
			// Move to next input (handled separately if search is focused)
			if !c.isSearchFocused {
				c.nextInput()
				return c, nil
			}

		case "shift+tab":
			// Move to previous input
			c.prevInput()
			return c, nil
		}

		// Handle number key selection for suggestions
		if msg.Type == tea.KeyRunes && len(msg.Runes) == 1 {
			key := string(msg.Runes[0])
			if num, err := strconv.Atoi(key); err == nil && num >= 1 && num <= len(c.suggestions) {
				c.fillNetworkData(c.suggestions[num-1])
				return c, nil
			}
		}

		// Update the currently focused input
		var cmd tea.Cmd
		switch c.focusIndex {
		case 0: // Search input
			oldValue := c.searchInput.Value()
			c.searchInput, cmd = c.searchInput.Update(msg)
			newValue := c.searchInput.Value()

			// Trigger search if value changed
			if oldValue != newValue {
				// Auto-search after a short delay
				c.loadingSuggestions = true
				c.selectedSuggestion = -1
				cmds = append(cmds, c.searchNetworks(newValue))
			}

		case 1: // Name input
			c.nameInput, cmd = c.nameInput.Update(msg)
			cmds = append(cmds, cmd)

		case 2: // Chain ID input
			c.chainIDInput, cmd = c.chainIDInput.Update(msg)
			cmds = append(cmds, cmd)

		case 3: // Symbol input
			c.symbolInput, cmd = c.symbolInput.Update(msg)
			cmds = append(cmds, cmd)

		case 4: // RPC endpoint input
			c.rpcEndpointInput, cmd = c.rpcEndpointInput.Update(msg)
			cmds = append(cmds, cmd)
		}
	}

	return c, tea.Batch(cmds...)
}

// nextInput focuses the next input field
func (c *AddNetworkComponent) nextInput() {
	c.focusIndex = (c.focusIndex + 1) % len(c.inputs)
	c.updateFocus()
}

// prevInput focuses the previous input field
func (c *AddNetworkComponent) prevInput() {
	c.focusIndex--
	if c.focusIndex < 0 {
		c.focusIndex = len(c.inputs) - 1
	}
	c.updateFocus()
}

// updateFocus updates the focus state of all inputs
func (c *AddNetworkComponent) updateFocus() {
	// Blur all inputs
	c.searchInput.Blur()
	c.nameInput.Blur()
	c.chainIDInput.Blur()
	c.symbolInput.Blur()
	c.rpcEndpointInput.Blur()

	// Track if search is focused
	c.isSearchFocused = (c.focusIndex == 0)

	// Focus the selected input
	switch c.focusIndex {
	case 0:
		c.searchInput.Focus()
	case 1:
		c.nameInput.Focus()
	case 2:
		c.chainIDInput.Focus()
	case 3:
		c.symbolInput.Focus()
	case 4:
		c.rpcEndpointInput.Focus()
	}
}

// validateInputs checks if the inputs are valid
func (c *AddNetworkComponent) validateInputs() bool {
	if strings.TrimSpace(c.nameInput.Value()) == "" {
		c.err = fmt.Errorf("network name cannot be empty")
		return false
	}

	if strings.TrimSpace(c.chainIDInput.Value()) == "" {
		c.err = fmt.Errorf("chain ID cannot be empty")
		return false
	}

	// Validate chain ID is a number
	if _, err := c.GetChainID(); err != nil {
		c.err = err
		return false
	}

	if strings.TrimSpace(c.symbolInput.Value()) == "" {
		c.err = fmt.Errorf("currency symbol cannot be empty")
		return false
	}

	if strings.TrimSpace(c.rpcEndpointInput.Value()) == "" {
		c.err = fmt.Errorf("RPC endpoint cannot be empty")
		return false
	}

	// Basic URL validation
	rpc := strings.TrimSpace(c.rpcEndpointInput.Value())
	if !strings.HasPrefix(rpc, "http://") && !strings.HasPrefix(rpc, "https://") {
		c.err = fmt.Errorf("RPC endpoint must start with http:// or https://")
		return false
	}

	c.err = nil
	return true
}

// View renders the add network component
func (c *AddNetworkComponent) View() string {
	var b strings.Builder

	// Header
	headerStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#FFF")).
		Background(lipgloss.Color("#874BFD")).
		MarginLeft(2).
		MarginBottom(1)
	b.WriteString(headerStyle.Render("🌐 Add Custom Network"))
	b.WriteString("\n\n")

	// Styles (InfoStyle, ErrorStyle, LoadingStyle are from styles.go)
	fieldStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("250")).
		MarginLeft(2).
		MarginBottom(1)

	labelStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#874BFD"))

	searchLabelStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("13"))

	// Search field
	b.WriteString(searchLabelStyle.Render("🔍 Search Networks:"))
	b.WriteString("\n")
	b.WriteString(fieldStyle.Render(c.searchInput.View()))

	// Sugestões interativas
	if c.loadingSuggestions {
		b.WriteString("\n")
		b.WriteString(InfoStyle.Render("🔍 Searching networks..."))
	} else if len(c.suggestions) > 0 {
		b.WriteString("\n")
		b.WriteString(c.suggestionList.View())
	}

	b.WriteString("\n\n")
	detailHeaderStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#FFF")).
		Background(lipgloss.Color("#874BFD")).
		MarginLeft(2).
		MarginBottom(1)
	b.WriteString(detailHeaderStyle.Render("Network Details:"))
	b.WriteString("\n\n")

	// Network Name field
	b.WriteString(labelStyle.Render("Network Name:"))
	b.WriteString("\n")
	b.WriteString(fieldStyle.Render(c.nameInput.View()))
	b.WriteString("\n")

	// Chain ID field
	b.WriteString(labelStyle.Render("Chain ID:"))
	b.WriteString("\n")
	b.WriteString(fieldStyle.Render(c.chainIDInput.View()))
	b.WriteString("\n")

	// Symbol field
	b.WriteString(labelStyle.Render("Native Currency Symbol:"))
	b.WriteString("\n")
	b.WriteString(fieldStyle.Render(c.symbolInput.View()))
	b.WriteString("\n")

	// RPC Endpoint field
	b.WriteString(labelStyle.Render("RPC Endpoint:"))
	b.WriteString("\n")
	b.WriteString(fieldStyle.Render(c.rpcEndpointInput.View()))
	b.WriteString("\n")

	// Status messages
	if c.adding {
		b.WriteString("\n")
		b.WriteString(LoadingStyle.Render("⏳ Adding network...")) // Uses LoadingStyle from styles.go
	} else if c.err != nil {
		b.WriteString("\n")
		b.WriteString(ErrorStyle.Render("❌ Error: " + c.err.Error())) // Uses ErrorStyle from styles.go
	}

	// Instructions
	b.WriteString("\n\n")
	b.WriteString(WarningStyle.Render("💡 Tips:"))
	b.WriteString("\n")
	b.WriteString(InfoStyle.Render("   • Search for networks by name and select from suggestions"))
	b.WriteString("\n")
	b.WriteString(InfoStyle.Render("   • Chain ID must be unique (check chainlist.org for reference)"))
	b.WriteString("\n")
	b.WriteString(InfoStyle.Render("   • Use reliable RPC endpoints for better performance"))
	b.WriteString("\n\n")

	// Footer
	b.WriteString(FooterStyle.Render("↑/↓: Navigate Suggestions • Tab: Next Field • Enter: Select/Submit • Esc: Back"))

	return b.String()
}

// AddNetworkRequestMsg is sent when the user wants to add a network
type AddNetworkRequestMsg struct {
	Name        string
	ChainID     string
	Symbol      string
	RPCEndpoint string
}
