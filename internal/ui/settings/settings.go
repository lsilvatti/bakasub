package settings

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/lsilvatti/bakasub/internal/config"
	"github.com/lsilvatti/bakasub/internal/ui/focus"
	"github.com/lsilvatti/bakasub/internal/ui/layout"
	"github.com/lsilvatti/bakasub/internal/ui/styles"
)

// Tab represents a settings tab
type Tab int

const (
	TabGeneral Tab = iota
	TabProviders
	TabModels
	TabPrompts
	TabAdvanced
)

// Model represents the settings screen state
type Model struct {
	width            int
	height           int
	activeTab        Tab
	config           *config.Config
	selectedLang     int
	customLangInput  textinput.Model
	selectedProvider int
	apiKeyInput      textinput.Model
	apiEndpointInput textinput.Model
	showAPIKey       bool
	profileKeys      []string
	selectedProfile  int
	selectedLogLevel int
	focusedInput     int // Deprecated - use focusManager
	focusManager     *focus.Manager
	saved            bool
}

// New creates a new settings model
func New(cfg *config.Config) Model {
	if cfg == nil {
		cfg = config.Default()
	}

	// Create text inputs
	apiKeyInput := textinput.New()
	apiKeyInput.Placeholder = "sk-or-v1-..."
	apiKeyInput.EchoMode = textinput.EchoPassword
	apiKeyInput.EchoCharacter = '•'
	apiKeyInput.CharLimit = 200
	apiKeyInput.SetValue(cfg.APIKey)

	apiEndpointInput := textinput.New()
	apiEndpointInput.Placeholder = "http://localhost:11434"
	apiEndpointInput.CharLimit = 200
	apiEndpointInput.SetValue(cfg.LocalEndpoint)

	customLangInput := textinput.New()
	customLangInput.Placeholder = "e.g. 'it', 'ru', 'zh-cn'"
	customLangInput.CharLimit = 10

	// Build profile keys
	profileKeys := make([]string, 0)
	for key, profile := range cfg.PromptProfiles {
		if profile.IsFactory {
			profileKeys = append([]string{key}, profileKeys...)
		} else {
			profileKeys = append(profileKeys, key)
		}
	}

	// Determine selected language
	selectedLang := 0
	switch cfg.TargetLang {
	case "PT-BR":
		selectedLang = 0
	case "EN-US":
		selectedLang = 1
	case "ES":
		selectedLang = 2
	case "JA-JP":
		selectedLang = 3
	case "FR-FR":
		selectedLang = 4
	default:
		selectedLang = 5
		customLangInput.SetValue(cfg.TargetLang)
	}

	// Determine selected provider
	selectedProvider := 0
	switch cfg.AIProvider {
	case "openrouter":
		selectedProvider = 0
	case "gemini":
		selectedProvider = 1
	case "openai":
		selectedProvider = 2
	case "local":
		selectedProvider = 3
	}

	selectedLogLevel := 0
	if cfg.LogLevel == "debug" {
		selectedLogLevel = 1
	}

	return Model{
		width:            80,
		height:           24,
		activeTab:        TabGeneral,
		config:           cfg,
		selectedLang:     selectedLang,
		customLangInput:  customLangInput,
		selectedProvider: selectedProvider,
		apiKeyInput:      apiKeyInput,
		apiEndpointInput: apiEndpointInput,
		showAPIKey:       false,
		profileKeys:      profileKeys,
		selectedProfile:  0,
		selectedLogLevel: selectedLogLevel,
		focusManager:     focus.NewManager(3), // 3 text inputs: customLang, apiKey, apiEndpoint
		focusedInput:     -1,
		saved:            false,
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		key := msg.String()

		// GATEKEEPER LOGIC: Route keys based on focus mode
		if m.focusManager.Mode() == focus.ModeInput {
			// In input mode - only ESC exits, all other keys go to active input
			if key == "esc" {
				m.focusManager.ExitInput()
				return m, nil
			}
			// Forward to active input
			return m.updateActiveInput(msg)
		}

		// In nav mode - handle Tab cycling first
		if m.focusManager.HandleTabCycle(key) {
			return m, nil
		}

		// Handle global hotkeys and navigation
		switch key {
		case "q", "esc":
			return m, tea.Quit
		case "ctrl+s":
			if err := m.saveConfig(); err == nil {
				m.saved = true
			}
			return m, nil
		case "enter":
			// Try to enter input mode for selected field
			if m.focusManager.HandleEnter() {
				return m, nil
			}
			// Otherwise save config
			if err := m.saveConfig(); err == nil {
				m.saved = true
			}
			return m, nil
		case "left", "h":
			if m.activeTab > TabGeneral {
				m.activeTab--
			}
		case "right", "l":
			if m.activeTab < TabAdvanced {
				m.activeTab++
			}
		case "up", "k":
			m.handleUp()
		case "down", "j":
			m.handleDown()
		case " ", "space":
			m.handleSpace()
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	}

	return m, cmd
}

// updateActiveInput forwards messages to the currently active text input
func (m Model) updateActiveInput(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	activeField := m.focusManager.ActiveField()

	switch m.activeTab {
	case TabGeneral:
		if m.selectedLang == 5 && activeField == 0 {
			m.customLangInput, cmd = m.customLangInput.Update(msg)
		}
	case TabProviders:
		if activeField == 1 {
			m.apiKeyInput, cmd = m.apiKeyInput.Update(msg)
		} else if activeField == 2 {
			m.apiEndpointInput, cmd = m.apiEndpointInput.Update(msg)
		}
	}

	return m, cmd
}

func (m *Model) handleUp() {
	switch m.activeTab {
	case TabGeneral:
		if m.selectedLang > 0 {
			m.selectedLang--
		}
	case TabProviders:
		if m.selectedProvider > 0 {
			m.selectedProvider--
		}
	case TabPrompts:
		if m.selectedProfile > 0 {
			m.selectedProfile--
		}
	case TabAdvanced:
		if m.selectedLogLevel > 0 {
			m.selectedLogLevel--
		}
	}
}

func (m *Model) handleDown() {
	switch m.activeTab {
	case TabGeneral:
		if m.selectedLang < 5 {
			m.selectedLang++
		}
	case TabProviders:
		if m.selectedProvider < 3 {
			m.selectedProvider++
		}
	case TabPrompts:
		if m.selectedProfile < len(m.profileKeys)-1 {
			m.selectedProfile++
		}
	case TabAdvanced:
		if m.selectedLogLevel < 1 {
			m.selectedLogLevel++
		}
	}
}

func (m *Model) handleSpace() {
	switch m.activeTab {
	case TabGeneral:
		m.config.RemoveHITags = !m.config.RemoveHITags
	case TabProviders:
		m.showAPIKey = !m.showAPIKey
		if m.showAPIKey {
			m.apiKeyInput.EchoMode = textinput.EchoNormal
		} else {
			m.apiKeyInput.EchoMode = textinput.EchoPassword
		}
	case TabAdvanced:
		m.config.AutoCheckUpdates = !m.config.AutoCheckUpdates
	}
}

func (m *Model) saveConfig() error {
	// Update config from UI state
	switch m.selectedLang {
	case 0:
		m.config.TargetLang = "PT-BR"
	case 1:
		m.config.TargetLang = "EN-US"
	case 2:
		m.config.TargetLang = "ES"
	case 3:
		m.config.TargetLang = "JA-JP"
	case 4:
		m.config.TargetLang = "FR-FR"
	case 5:
		m.config.TargetLang = m.customLangInput.Value()
	}

	switch m.selectedProvider {
	case 0:
		m.config.AIProvider = "openrouter"
	case 1:
		m.config.AIProvider = "gemini"
	case 2:
		m.config.AIProvider = "openai"
	case 3:
		m.config.AIProvider = "local"
	}

	m.config.APIKey = m.apiKeyInput.Value()
	m.config.LocalEndpoint = m.apiEndpointInput.Value()

	if m.selectedLogLevel == 1 {
		m.config.LogLevel = "debug"
	} else {
		m.config.LogLevel = "info"
	}

	return m.config.Save()
}

func (m Model) View() string {
	// Check if terminal is too small
	if layout.IsTooSmall(m.width, m.height) {
		return layout.RenderTooSmallWarning(m.width, m.height)
	}

	header := m.renderHeader()
	tabs := m.renderTabs()
	content := m.renderTabContent()
	footer := m.renderFooter()

	view := lipgloss.JoinVertical(
		lipgloss.Left,
		header,
		tabs,
		"",
		content,
		"",
		footer,
	)

	contentWidth := layout.SafeWidth(m.width-4, 76)
	return styles.MainWindow.Width(contentWidth).Render(view)
}

func (m Model) renderHeader() string {
	title := styles.PanelTitle.Render("CONFIGURATION")
	headerBar := strings.Repeat("▒", m.width-4)

	return lipgloss.JoinVertical(
		lipgloss.Left,
		headerBar,
		title,
		headerBar,
	)
}

func (m Model) renderTabs() string {
	tabs := []string{"GENERAL", "AI PROVIDERS", "AI MODELS", "PROMPTS", "ADVANCED"}

	rendered := make([]string, len(tabs))
	for i, tab := range tabs {
		if Tab(i) == m.activeTab {
			rendered[i] = styles.Highlight.Render("[ " + tab + " ]")
		} else {
			rendered[i] = styles.Dimmed.Render("[ " + tab + " ]")
		}
	}

	return lipgloss.JoinHorizontal(lipgloss.Top, rendered...)
}

func (m Model) renderTabContent() string {
	switch m.activeTab {
	case TabGeneral:
		return m.renderGeneralTab()
	case TabProviders:
		return m.renderProvidersTab()
	case TabModels:
		return m.renderModelsTab()
	case TabPrompts:
		return m.renderPromptsTab()
	case TabAdvanced:
		return m.renderAdvancedTab()
	}
	return ""
}

func (m Model) renderGeneralTab() string {
	targetLangs := []string{
		"PT-BR (Português)",
		"EN-US (English)",
		"ES (Español)",
		"JA-JP (Japanese)",
		"FR-FR (Français)",
		"OTHER ISO CODE",
	}

	targetList := ""
	for i, lang := range targetLangs {
		icon := "( )"
		if i == m.selectedLang {
			icon = "(o)"
		}
		targetList += fmt.Sprintf("%s %s\n", styles.Highlight.Render(icon), lang)

		if i == 5 && m.selectedLang == 5 {
			// Configure and style the custom lang input with focus manager
			m.focusManager.ConfigureInput(&m.customLangInput, 0)
			inputStyle := m.focusManager.FieldStyle(0, m.focusManager.ActiveField() == 0)
			inputView := inputStyle.Render(m.customLangInput.View())
			targetList += "    " + inputView + "\n"
			if m.focusManager.Mode() == focus.ModeNav {
				targetList += "    " + styles.KeyHintStyle.Render("[ENTER] to edit, [ESC] to exit") + "\n"
			}
		}
	}

	hiCheckbox := "[ ]"
	if m.config.RemoveHITags {
		hiCheckbox = "[X]"
	}

	content := lipgloss.JoinVertical(
		lipgloss.Left,
		styles.PanelTitle.Render("DEFAULT TARGET LANG (Output Language)"),
		targetList,
		"",
		styles.PanelTitle.Render("PREFERENCES"),
		fmt.Sprintf("%s REMOVE HEARING IMPAIRED TAGS", hiCheckbox),
		fmt.Sprintf("GLOBAL TEMPERATURE: [ %.1f   ]", m.config.GlobalTemperature),
	)

	panelStyle := styles.Panel.Width(m.width - 10)
	if m.activeTab == TabGeneral {
		panelStyle = panelStyle.BorderForeground(styles.Yellow)
	}
	return panelStyle.Render(content)
}

func (m Model) renderProvidersTab() string {
	providers := []string{
		"♾️  OpenRouter (Recommended)",
		"💎 Google Gemini API",
		"🤖 OpenAI API",
		"🏠 Local LLM (Ollama/LMStudio)",
	}

	providerList := ""
	for i, prov := range providers {
		icon := "( )"
		if i == m.selectedProvider {
			icon = "(o)"
		}
		providerList += fmt.Sprintf("%s %s\n", styles.Highlight.Render(icon), prov)
	}

	configContent := ""
	if m.selectedProvider == 3 {
		// Local LLM - endpoint input (field index 2)
		m.focusManager.ConfigureInput(&m.apiEndpointInput, 2)
		inputStyle := m.focusManager.FieldStyle(2, m.focusManager.ActiveField() == 2)
		inputView := inputStyle.Render(m.apiEndpointInput.View())
		configContent = "ENDPOINT:\n" + inputView + "\n"
		if m.focusManager.Mode() == focus.ModeNav {
			configContent += styles.KeyHintStyle.Render("[ENTER] to edit, [ESC] to exit")
		}
	} else {
		// API Key input (field index 1)
		m.focusManager.ConfigureInput(&m.apiKeyInput, 1)
		inputStyle := m.focusManager.FieldStyle(1, m.focusManager.ActiveField() == 1)
		inputView := inputStyle.Render(m.apiKeyInput.View())

		showText := "SHOW"
		if m.showAPIKey {
			showText = "HIDE"
		}
		configContent = "API KEY:\n" + inputView + "\n" +
			"[ " + styles.KeyHintStyle.Render(showText+" (SPACE)") + " ]"
		if m.focusManager.Mode() == focus.ModeNav {
			configContent += "\n" + styles.KeyHintStyle.Render("[ENTER] to edit, [ESC] to exit")
		}
	}

	content := lipgloss.JoinVertical(
		lipgloss.Left,
		styles.PanelTitle.Render("ACTIVE PROVIDER"),
		"",
		providerList,
		"",
		styles.PanelTitle.Render("CONFIGURATION"),
		configContent,
	)

	panelStyle := styles.Panel.Width(m.width - 10)
	if m.activeTab == TabProviders {
		panelStyle = panelStyle.BorderForeground(styles.Yellow)
	}
	return panelStyle.Render(content)
}

func (m Model) renderModelsTab() string {
	// Model selection with FREE and ALL MODELS tabs only
	subTabs := styles.Dimmed.Render("SUB-TABS:  ") +
		styles.Highlight.Render("< FREE >") +
		styles.Dimmed.Render("   < ALL MODELS >") +
		styles.Dimmed.Render("   < SEARCH >")

	searchBox := "SEARCH > " + styles.Dimmed.Render("[Type to filter models...]")

	// Example model list (FREE models)
	modelList := lipgloss.JoinVertical(
		lipgloss.Left,
		"",
		styles.Dimmed.Render("NAME                          COST(1M)   CTX      TAGS"),
		styles.Dimmed.Render("────────────────────────────────────────────────────────────────"),
		"( ) Llama 3.3 70B              FREE       128k     [OpenSource]",
		"( ) Qwen 2.5 72B               FREE       32k      [Multilingual]",
		"( ) Mistral 7B Instruct        FREE       32k      [Fast]",
		"( ) Gemma 2 9B                 FREE       8k       [Google]",
		"( ) Phi-3 Medium               FREE       128k     [Microsoft]",
		"",
		styles.Dimmed.Render("[← PREV]   Page 1/5   [NEXT →]"),
	)

	helperBox := lipgloss.JoinVertical(
		lipgloss.Left,
		"",
		styles.PanelTitle.Render("💡 INFO"),
		styles.Dimmed.Render("• FREE: Zero-cost models (community/open-source)"),
		styles.Dimmed.Render("• ALL MODELS: Complete catalog with pricing"),
		styles.Dimmed.Render("• SEARCH: Filter by name, provider, or capability"),
		styles.Dimmed.Render("• Use ↑↓ to select, [ENTER] to confirm"),
	)

	content := lipgloss.JoinVertical(
		lipgloss.Left,
		styles.PanelTitle.Render("MODEL SELECTION"),
		"",
		subTabs,
		"",
		searchBox,
		modelList,
		helperBox,
	)

	panelStyle := styles.Panel.Width(m.width - 10)
	if m.activeTab == TabModels {
		panelStyle = panelStyle.BorderForeground(styles.Yellow)
	}
	return panelStyle.Render(content)
}

func (m Model) renderPromptsTab() string {
	dropdown := "SELECT PROFILE:\n"
	if len(m.profileKeys) > m.selectedProfile {
		currentKey := m.profileKeys[m.selectedProfile]
		currentProfile := m.config.PromptProfiles[currentKey]
		icon := "🔒"
		if !currentProfile.IsFactory {
			icon = "👤"
		}
		dropdown += fmt.Sprintf("[ %s %s           ▼ ]\n", icon, currentProfile.Name)
	}

	status := ""
	if len(m.profileKeys) > m.selectedProfile {
		currentProfile := m.config.PromptProfiles[m.profileKeys[m.selectedProfile]]
		if currentProfile.IsLocked {
			status = styles.StatusWarning.Render("[LOCKED]") + " Cannot be modified. Clone to customize."
		} else {
			status = styles.StatusOK.Render("[EDITABLE]") + " Can be modified or deleted."
		}
	}

	preview := ""
	if len(m.profileKeys) > m.selectedProfile {
		currentProfile := m.config.PromptProfiles[m.profileKeys[m.selectedProfile]]
		lines := strings.Split(currentProfile.SystemPrompt, "\n")
		if len(lines) > 5 {
			lines = lines[:5]
		}
		preview = strings.Join(lines, "\n") + "\n" + styles.Dimmed.Render("... (content dimmed)")
	}

	content := lipgloss.JoinVertical(
		lipgloss.Left,
		styles.PanelTitle.Render("PROFILE MANAGER"),
		"",
		dropdown,
		status,
		"",
		styles.PanelTitle.Render("SYSTEM PROMPT PREVIEW (READ ONLY)"),
		styles.CodeBlock.Render(preview),
	)

	panel := styles.Panel.Width(m.width - 10)
	if m.activeTab == 3 {
		panel = panel.BorderForeground(styles.Yellow)
	}
	return panel.Render(content)
}

func (m Model) renderAdvancedTab() string {
	logLevels := []string{"INFO (Default)", "DEBUG (Verbose)"}
	logList := ""
	for i, level := range logLevels {
		icon := "( )"
		if i == m.selectedLogLevel {
			icon = "(o)"
		}
		logList += fmt.Sprintf("%s %s\n", icon, level)
	}

	content := lipgloss.JoinVertical(
		lipgloss.Left,
		styles.PanelTitle.Render("DEBUGGING"),
		logList,
		"",
		styles.PanelTitle.Render("SYSTEM INFO"),
		"VERSION: v1.0.0 (Build 20240520)",
		"GO VERSION: 1.24.0",
	)

	panel := styles.Panel.Width(m.width - 10)
	if m.activeTab == 4 {
		panel = panel.BorderForeground(styles.Yellow)
	}
	return panel.Render(content)
}

func (m Model) renderFooter() string {
	saveStatus := ""
	if m.saved {
		saveStatus = styles.StatusOK.Render("[SAVED] ")
	}

	return styles.Footer.Render(
		saveStatus +
			styles.RenderHotkey("ESC", "CANCEL") + "  " +
			styles.RenderHotkey("ENTER", "SAVE") + "  " +
			styles.RenderHotkey("←/→", "SWITCH TAB") + "  " +
			styles.RenderHotkey("↑/↓", "NAVIGATE"),
	)
}
