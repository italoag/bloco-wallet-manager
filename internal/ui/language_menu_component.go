package ui

import (
	"blocowallet/pkg/config"
	"fmt"

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
	id          string
	title       string
	description string
	code        string
	isCurrent   bool
}

func (i languageItem) Title() string       { return zone.Mark(i.id, i.title) }
func (i languageItem) Description() string { return i.description }
func (i languageItem) FilterValue() string { return zone.Mark(i.id, i.title) }

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
	// Use default delegate instead of custom delegate to avoid conflicts
	delegate := list.NewDefaultDelegate()

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

// RefreshLanguages updates the language list with current languages
func (c *LanguageMenuComponent) RefreshLanguages() {
	var items []list.Item

	langCodes := c.config.GetLanguageCodes()
	for i, code := range langCodes {
		name := config.SupportedLanguages[code]
		isCurrent := c.config.Language == code

		title := name
		description := fmt.Sprintf("Language: %s", code)
		if isCurrent {
			title += " ✓ Current"
			description += " • Currently selected"
		}

		items = append(items, languageItem{
			id:          fmt.Sprintf("lang_%s_%d", code, i),
			title:       title,
			description: description,
			code:        code,
			isCurrent:   isCurrent,
		})
	}

	// Add back button
	items = append(items, languageItem{
		id:          "lang_back_to_settings",
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
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		c.width = msg.Width
		c.height = msg.Height
		c.list.SetSize(msg.Width, msg.Height)

	case tea.MouseMsg:
		if msg.Button == tea.MouseButtonWheelUp {
			c.list.CursorUp()
			return c, nil
		}

		if msg.Button == tea.MouseButtonWheelDown {
			c.list.CursorDown()
			return c, nil
		}

		if msg.Action == tea.MouseActionRelease && msg.Button == tea.MouseButtonLeft {
			for i, listItem := range c.list.VisibleItems() {
				v, _ := listItem.(languageItem)
				// Check each item to see if it's in bounds.
				if zone.Get(v.id).InBounds(msg) {
					// If so, select it in the list.
					c.list.Select(i)
					// Trigger selection action
					return c, func() tea.Msg { return LanguageSelectedMsg{Code: v.code} }
				}
			}
		}

		return c, nil

	case tea.KeyMsg:
		// Handle escape key
		switch msg.String() {
		case "enter":
			if item, ok := c.list.SelectedItem().(languageItem); ok {
				return c, func() tea.Msg { return LanguageSelectedMsg{Code: item.code} }
			}
		case "esc", "q":
			return c, func() tea.Msg { return BackToSettingsMsg{} }
		}
	}

	var cmd tea.Cmd
	c.list, cmd = c.list.Update(msg)
	return c, cmd
}

// View renders the language menu component
func (c *LanguageMenuComponent) View() string {
	// Use zone.Scan to wrap the list view for mouse support
	return zone.Scan(appStyle.Render(c.list.View()))
}

// Language-related messages
type LanguageSelectedMsg struct {
	Code string
}
