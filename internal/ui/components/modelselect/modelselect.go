package modelselect

import (
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/lsilvatti/bakasub/internal/locales"
	"github.com/lsilvatti/bakasub/internal/ui/focus"
	"github.com/lsilvatti/bakasub/internal/ui/styles"
)

// ModelInfo contains detailed information about an AI model
type ModelInfo struct {
	ID          string
	Name        string
	ContextSize string
	PricePerM   string
	IsFree      bool
	// Internal fields for sorting
	priceValue   float64 // Numeric value for sorting
	contextValue int     // Numeric value for sorting (in tokens)
}

// Model represents the model selector state
type Model struct {
	// Display
	width        int
	contentWidth int
	isActive     bool // Whether this component is focused

	// Models
	availableModels []ModelInfo
	filteredModels  []ModelInfo
	selectedIndex   int
	scrollOffset    int

	// Tabs and search
	currentTab   int // 0=FREE, 1=ALL
	searchInput  textinput.Model
	focusManager *focus.Manager

	// Config
	visibleRows int // Number of rows to display
}

// New creates a new model selector
func New(focusManager *focus.Manager) Model {
	searchInput := textinput.New()
	searchInput.Placeholder = locales.T("settings.models.search_placeholder")
	searchInput.CharLimit = 50
	searchInput.Width = 54

	return Model{
		searchInput:     searchInput,
		focusManager:    focusManager,
		currentTab:      0, // Start on FREE
		selectedIndex:   0,
		scrollOffset:    0,
		visibleRows:     7, // Default
		availableModels: []ModelInfo{},
		filteredModels:  []ModelInfo{},
		contentWidth:    70, // Default
	}
}

// SetWidth sets the width for responsive rendering
func (m *Model) SetWidth(width int) {
	m.width = width
	m.contentWidth = width - 6
	if m.contentWidth < 50 {
		m.contentWidth = 50
	}
	m.searchInput.Width = m.contentWidth - 20
}

// SetActive sets whether this component is focused/active
func (m *Model) SetActive(active bool) {
	m.isActive = active
}

// IsActive returns whether this component is focused
func (m Model) IsActive() bool {
	return m.isActive
}

// SetModels updates the available models and sorts them
func (m *Model) SetModels(models []ModelInfo) {
	// Parse and enrich models with numeric values for sorting
	for i := range models {
		models[i].priceValue = m.parsePriceValue(models[i].PricePerM)
		models[i].contextValue = m.parseContextValue(models[i].ContextSize)
	}

	// Sort: Primary by price (ascending - cheaper first), Secondary by context (descending - larger first)
	sort.SliceStable(models, func(i, j int) bool {
		if models[i].priceValue != models[j].priceValue {
			return models[i].priceValue < models[j].priceValue
		}
		return models[i].contextValue > models[j].contextValue
	})

	m.availableModels = models
	m.filterModels()
}

// parsePriceValue converts price string to float for sorting
func (m Model) parsePriceValue(priceStr string) float64 {
	if priceStr == "FREE" || priceStr == "" {
		return 0.0
	}
	// Remove $ and parse
	priceStr = strings.TrimPrefix(priceStr, "$")
	price, _ := strconv.ParseFloat(priceStr, 64)
	return price
}

// parseContextValue converts context string to int for sorting
func (m Model) parseContextValue(contextStr string) int {
	contextStr = strings.ToLower(strings.TrimSpace(contextStr))

	// Handle special cases
	if contextStr == "varies" || contextStr == "" {
		return 0
	}

	// Parse: "128k" -> 128000, "1M" -> 1000000
	multiplier := 1
	if strings.HasSuffix(contextStr, "m") {
		multiplier = 1000000
		contextStr = strings.TrimSuffix(contextStr, "m")
	} else if strings.HasSuffix(contextStr, "k") {
		multiplier = 1000
		contextStr = strings.TrimSuffix(contextStr, "k")
	}

	value, _ := strconv.Atoi(contextStr)
	return value * multiplier
}

// GetSelectedModel returns the currently selected model (nil if none)
func (m Model) GetSelectedModel() *ModelInfo {
	models := m.getDisplayModels()
	if m.selectedIndex >= 0 && m.selectedIndex < len(models) {
		return &models[m.selectedIndex]
	}
	return nil
}

// Update handles messages
func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd

	// Only process input if this component is active
	if !m.isActive {
		return m, nil
	}

	if keyMsg, ok := msg.(tea.KeyMsg); ok {
		// If in input mode (search active), forward to search input
		if m.focusManager != nil && m.focusManager.Mode() == focus.ModeInput {
			// Handle ESC to exit search mode
			if keyMsg.String() == "esc" {
				m.focusManager.ExitInput()
				m.searchInput.Blur()
				// Clear search and re-filter
				m.searchInput.SetValue("")
				m.filterModels()
				return m, nil
			}

			m.searchInput, cmd = m.searchInput.Update(msg)

			// Re-filter on every keystroke
			if keyMsg.Type == tea.KeyRunes || keyMsg.Type == tea.KeyBackspace || keyMsg.Type == tea.KeyDelete {
				m.filterModels()
			}

			return m, cmd
		}

		// Navigation mode - handle tab switching and model selection
		switch keyMsg.String() {
		case "left", "h":
			if m.currentTab > 0 {
				m.currentTab--
				m.selectedIndex = 0
				m.scrollOffset = 0
				m.filterModels()
			}
		case "right", "l":
			if m.currentTab < 1 {
				m.currentTab++
				m.selectedIndex = 0
				m.scrollOffset = 0
				m.filterModels()
			}
		case "up", "k":
			if m.selectedIndex > 0 {
				m.selectedIndex--
				if m.selectedIndex < m.scrollOffset {
					m.scrollOffset = m.selectedIndex
				}
			}
		case "down", "j":
			models := m.getDisplayModels()
			if m.selectedIndex < len(models)-1 {
				m.selectedIndex++
				if m.selectedIndex-m.scrollOffset >= m.visibleRows {
					m.scrollOffset = m.selectedIndex - m.visibleRows + 1
				}
			}
		case "/":
			// Activate search mode (only in ALL tab)
			if m.currentTab == 1 && m.focusManager != nil {
				m.focusManager.EnterInput(0)
				m.searchInput.Focus()
				return m, textinput.Blink
			}
		default:
			// If in ALL tab and user types a letter, activate search automatically
			if m.currentTab == 1 && m.focusManager != nil && keyMsg.Type == tea.KeyRunes {
				m.focusManager.EnterInput(0)
				m.searchInput.Focus()
				// Forward the key to the input
				m.searchInput, cmd = m.searchInput.Update(msg)
				m.filterModels()
				return m, textinput.Blink
			}
		}
	}

	return m, cmd
}

// filterModels filters available models based on current tab and search query
func (m *Model) filterModels() {
	m.filteredModels = []ModelInfo{}
	searchQuery := strings.ToLower(strings.TrimSpace(m.searchInput.Value()))

	for _, model := range m.availableModels {
		// Filter by tab (FREE vs ALL)
		if m.currentTab == 0 && !model.IsFree {
			continue
		}

		// Filter by search (only in ALL tab)
		if m.currentTab == 1 && searchQuery != "" {
			nameMatch := strings.Contains(strings.ToLower(model.Name), searchQuery)
			idMatch := strings.Contains(strings.ToLower(model.ID), searchQuery)
			if !nameMatch && !idMatch {
				continue
			}
		}

		m.filteredModels = append(m.filteredModels, model)
	}

	// Reset selection if out of bounds
	if m.selectedIndex >= len(m.filteredModels) {
		m.selectedIndex = 0
		m.scrollOffset = 0
	}
}

// getDisplayModels returns the models to display based on current filters
func (m Model) getDisplayModels() []ModelInfo {
	if len(m.filteredModels) > 0 {
		return m.filteredModels
	}

	// If no filters applied, show based on tab
	if m.currentTab == 0 {
		// FREE tab - only free models
		free := []ModelInfo{}
		for _, model := range m.availableModels {
			if model.IsFree {
				free = append(free, model)
			}
		}
		return free
	}

	// ALL tab - all models
	return m.availableModels
}

// View renders the component
func (m Model) View() string {
	var content strings.Builder

	content.WriteString(fmt.Sprintf("   %s\n\n", locales.T("wizard.step2.model_selection_title")))

	// Render tabs
	freeTab := " " + locales.T("settings.models.tab_free") + " "
	allTab := " " + locales.T("settings.models.tab_all") + " "

	if m.currentTab == 0 {
		freeTab = styles.StatusOK.Render(" [ " + locales.T("settings.models.tab_free") + " ] ")
		allTab = styles.Dimmed.Render(" " + locales.T("settings.models.tab_all") + " ")
	} else {
		freeTab = styles.Dimmed.Render(" " + locales.T("settings.models.tab_free") + " ")
		allTab = styles.StatusOK.Render(" [ " + locales.T("settings.models.tab_all") + " ] ")
	}

	switchHint := styles.KeyHintStyle.Render("[← →] " + locales.T("settings.models.tab_switch"))
	content.WriteString("   " + freeTab + allTab + "  " + switchHint + "\n\n")

	// Search input (only in ALL MODELS tab)
	if m.currentTab == 1 {
		searchView := m.searchInput.View()
		searchHint := ""

		if m.isActive && m.focusManager != nil {
			if m.focusManager.Mode() == focus.ModeInput {
				searchHint = " " + styles.StatusOK.Render("["+locales.T("common.edit")+"]") + " " + styles.KeyHintStyle.Render("[ESC] "+locales.T("settings.models.stop_search"))
			} else {
				// Make it very clear how to start searching
				searchHint = " " + styles.KeyHintStyle.Render("[ / ] "+locales.T("settings.models.start_search")+" | "+locales.T("wizard.step2.search_hint"))
			}
		}

		content.WriteString("   " + locales.T("settings.models.search_label") + " > " + searchView + searchHint + "\n")
		// Add extra hint if not searching
		if m.isActive && m.focusManager != nil && m.focusManager.Mode() != focus.ModeInput {
			content.WriteString("   " + styles.Dimmed.Render(locales.T("wizard.step2.search_instruction")) + "\n")
		}
		content.WriteString("\n")
	}

	// Calculate column widths
	nameColWidth := m.contentWidth - 30
	if nameColWidth < 20 {
		nameColWidth = 20
	}

	// Table header
	content.WriteString("   ┌" + strings.Repeat("─", nameColWidth) + "┬───────────┬──────────────┐\n")
	content.WriteString(fmt.Sprintf("   │ %-*s│ %-9s │ %-12s │\n",
		nameColWidth-1,
		locales.T("settings.models.columns.name"),
		locales.T("settings.models.columns.context"),
		locales.T("settings.models.columns.cost")))
	content.WriteString("   ├" + strings.Repeat("─", nameColWidth) + "┼───────────┼──────────────┤\n")

	// Table rows
	models := m.getDisplayModels()
	if len(models) == 0 {
		noModelsMsg := locales.T("settings.models.no_models")
		padding := m.contentWidth - len(noModelsMsg) - 10
		if padding < 0 {
			padding = 0
		}
		content.WriteString("   │ " + styles.Dimmed.Render(noModelsMsg) + strings.Repeat(" ", padding) + "│\n")
	} else {
		for i := 0; i < m.visibleRows && i+m.scrollOffset < len(models); i++ {
			idx := i + m.scrollOffset
			model := models[idx]

			// Truncate name based on available width
			name := model.Name
			maxNameLen := nameColWidth - 3
			if len(name) > maxNameLen {
				name = name[:maxNameLen-3] + "..."
			}

			row := fmt.Sprintf(" %-*s│ %-9s │ %-12s ",
				nameColWidth-1, name, model.ContextSize, model.PricePerM)

			if idx == m.selectedIndex {
				content.WriteString("   │" + styles.Highlight.Render(row) + "│\n")
			} else {
				content.WriteString("   │" + row + "│\n")
			}
		}
	}

	content.WriteString("   └" + strings.Repeat("─", nameColWidth) + "┴───────────┴──────────────┘\n")

	// Scroll indicator
	if len(models) > m.visibleRows {
		content.WriteString(fmt.Sprintf("   %s (%d/%d)\n",
			styles.Dimmed.Render("↑↓ "+locales.T("settings.models.scroll_indicator")),
			m.selectedIndex+1,
			len(models)))
	}

	return content.String()
}

// RenderWithInfo renders the component with info box
func (m Model) RenderWithInfo() string {
	mainView := m.View()

	// Info box
	var infoBox strings.Builder
	infoBox.WriteString("\n")
	infoBox.WriteString("   " + locales.T("settings.models.info_title") + "\n")
	infoBox.WriteString("   " + styles.Dimmed.Render(locales.T("settings.models.info_free")) + "\n")
	infoBox.WriteString("   " + styles.Dimmed.Render(locales.T("settings.models.info_all")) + "\n")
	infoBox.WriteString("   " + styles.Dimmed.Render(locales.T("settings.models.info_search")) + "\n")
	infoBox.WriteString("   " + styles.Dimmed.Render(locales.T("settings.models.info_controls")) + "\n")

	return mainView + infoBox.String()
}
