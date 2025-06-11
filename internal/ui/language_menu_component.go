package ui

import (
	"blocowallet/pkg/config"
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	zone "github.com/lrstanley/bubblezone"
)

// LanguageMenuComponent represents the language menu component
type LanguageMenuComponent struct {
	id     string
	list   list.Model
	width  int
	height int
	keys   *languageKeyMap
	config *config.Config
}

// languageItem represents a language item
type languageItem struct {
	title       string
	description string
	code        string
	isCurrent   bool
}

func (i languageItem) Title() string       { return i.title }
func (i languageItem) Description() string { return i.description }
func (i languageItem) FilterValue() string { return i.title }

// languageKeyMap defines key bindings for the language menu
type languageKeyMap struct {
	choose key.Binding
}

func (k languageKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.choose}
}

func (k languageKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{{k.choose}}
}

func newLanguageKeyMap() *languageKeyMap {
	return &languageKeyMap{
		choose: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter", "select"),
		),
	}
}

// NewLanguageMenuComponent creates a new language menu component
func NewLanguageMenuComponent(cfg *config.Config) LanguageMenuComponent {
	keys := newLanguageKeyMap()
	delegate := newLanguageDelegate(keys)

	languageList := list.New([]list.Item{}, delegate, 0, 0)
	languageList.Title = "🌍 Language Selection"
	languageList.Styles.Title = titleStyle
	languageList.SetShowStatusBar(false)
	languageList.SetFilteringEnabled(false)

	c := LanguageMenuComponent{
		id:     "language-menu",
		list:   languageList,
		keys:   keys,
		config: cfg,
	}

	c.RefreshLanguages()
	return c
}

// newLanguageDelegate creates a delegate for the language list
func newLanguageDelegate(keys *languageKeyMap) list.DefaultDelegate {
	d := list.NewDefaultDelegate()

	d.UpdateFunc = func(msg tea.Msg, m *list.Model) tea.Cmd {
		var item languageItem
		var ok bool

		if i := m.SelectedItem(); i != nil {
			item, ok = i.(languageItem)
			if !ok {
				return nil
			}
		}

		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch {
			case key.Matches(msg, keys.choose):
				return func() tea.Msg {
					return LanguageSelectedMsg{Code: item.code}
				}
			}
		}

		return nil
	}

	return d
}

// RefreshLanguages updates the language list with current languages
func (c *LanguageMenuComponent) RefreshLanguages() {
	var items []list.Item

	langCodes := c.config.GetLanguageCodes()
	for _, code := range langCodes {
		name := config.SupportedLanguages[code]
		isCurrent := c.config.Language == code

		title := name
		description := fmt.Sprintf("Language: %s", code)
		if isCurrent {
			title += " ✓ Current"
			description += " • Currently selected"
		}

		items = append(items, languageItem{
			title:       title,
			description: description,
			code:        code,
			isCurrent:   isCurrent,
		})
	}

	// Add back button
	items = append(items, languageItem{
		title:       "🔙 Back to Settings",
		description: "Return to the settings menu",
		code:        "back-to-settings",
	})

	c.list.SetItems(items)
}

// SetSize updates the component size
func (c *LanguageMenuComponent) SetSize(width, height int) {
	c.width = width
	c.height = height
	c.list.SetSize(width, height)
}

// GetSelected returns the currently selected language code
func (c *LanguageMenuComponent) GetSelected() string {
	if item, ok := c.list.SelectedItem().(languageItem); ok {
		return item.code
	}
	return ""
}

// Update handles messages for the language menu component
func (c *LanguageMenuComponent) Update(msg tea.Msg) (*LanguageMenuComponent, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		c.width = msg.Width
		c.height = msg.Height
		c.list.SetSize(msg.Width, msg.Height)

	case tea.MouseMsg:
		// Check if any language item was clicked
		for i := 0; i < len(c.list.Items()); i++ {
			itemZoneID := fmt.Sprintf("language-item-%d", i)
			if zone.Get(itemZoneID).InBounds(msg) {
				// Select and activate the clicked item
				c.list.Select(i)
				if item, ok := c.list.SelectedItem().(languageItem); ok {
					return c, func() tea.Msg { return LanguageSelectedMsg{Code: item.code} }
				}
			}
		}

	case tea.KeyMsg:
		// Handle escape key
		switch msg.String() {
		case "esc", "q":
			return c, func() tea.Msg { return BackToSettingsMsg{} }
		}
	}

	// Update the list
	c.list, cmd = c.list.Update(msg)
	return c, cmd
}

// View renders the language menu component
func (c *LanguageMenuComponent) View() string {
	// Render the list first
	listView := c.list.View()

	// Apply zone marking to each language item for mouse support
	lines := strings.Split(listView, "\n")
	var markedLines []string

	itemIndex := 0
	for _, line := range lines {
		// Check if this line contains a language item (has content and isn't just formatting)
		if strings.TrimSpace(line) != "" &&
			!strings.Contains(line, "Language") &&
			!strings.Contains(line, "Help") &&
			(strings.Contains(line, "►") || strings.Contains(line, "•") || strings.Contains(line, "🗣️") || strings.Contains(line, "🔙")) {

			// Mark this line as clickable
			zoneID := fmt.Sprintf("language-item-%d", itemIndex)
			markedLine := zone.Mark(zoneID, line)
			markedLines = append(markedLines, markedLine)
			itemIndex++
		} else {
			markedLines = append(markedLines, line)
		}
	}

	return appStyle.Render(strings.Join(markedLines, "\n"))
}

// Language-related messages
type LanguageSelectedMsg struct {
	Code string
}
