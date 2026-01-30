package settings

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/lsilvatti/bakasub/internal/config"
	"github.com/lsilvatti/bakasub/internal/core/ai"
	"github.com/lsilvatti/bakasub/internal/locales"
	"github.com/lsilvatti/bakasub/internal/ui/components/langselector"
	"github.com/lsilvatti/bakasub/internal/ui/components/modelselect"
	"github.com/lsilvatti/bakasub/internal/ui/focus"
	"github.com/lsilvatti/bakasub/internal/ui/layout"
	"github.com/lsilvatti/bakasub/internal/ui/styles"
)

// Message types for settings screen communication
type (
	// SavedMsg is sent when settings are saved successfully
	SavedMsg struct {
		Config *config.Config
	}

	// CancelledMsg is sent when settings are cancelled
	CancelledMsg struct{}

	// modelsLoadedMsg is sent when models are fetched
	modelsLoadedMsg struct {
		models []modelselect.ModelInfo
		err    error
	}
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
	width  int
	height int

	// Tab navigation
	activeTab Tab

	// Config
	config      *config.Config
	originalCfg *config.Config // Original for cancel/restore

	// Focus management
	focusManager *focus.Manager

	// General tab
	targetLangSelector langselector.Model
	selectedTargetLang int // 6 = OTHER
	customISOInput     textinput.Model

	// Providers tab
	selectedProvider int // 0=openrouter, 1=gemini, 2=openai, 3=local
	apiKeyInput      textinput.Model
	apiEndpointInput textinput.Model
	showAPIKey       bool

	// Models tab
	modelSelector modelselect.Model
	modelsLoaded  bool

	// Prompts tab
	profileKeys        []string
	selectedProfile    int
	editingPrompt      bool
	promptInput        textinput.Model
	editingProfileName bool
	profileNameInput   textinput.Model
	showProfileList    bool
	temperatureInput   textinput.Model

	// Touchless configuration modal
	showTouchlessModal bool
	touchlessMultiSub  int // 0=largest, 1=smallest, 2=skip
	touchlessProfile   int // index in profileKeys
	touchlessMuxMode   int // 0=replace, 1=new file

	// Advanced tab
	selectedLogLevel int // 0=info, 1=debug

	// State
	saved    bool
	hasError bool
	errMsg   string
}

// New creates a new settings model
func New(cfg *config.Config) Model {
	if cfg == nil {
		cfg = config.Default()
	}

	// Clone config for original state
	origCfg := *cfg

	// Create focus manager
	fm := focus.NewManager(5)

	// Create text inputs
	apiKeyInput := textinput.New()
	apiKeyInput.Placeholder = "sk-or-v1-..."
	apiKeyInput.EchoMode = textinput.EchoPassword
	apiKeyInput.EchoCharacter = '•'
	apiKeyInput.CharLimit = 200
	apiKeyInput.Width = 54
	apiKeyInput.SetValue(cfg.APIKey)

	apiEndpointInput := textinput.New()
	apiEndpointInput.Placeholder = "http://localhost:11434"
	apiEndpointInput.CharLimit = 200
	apiEndpointInput.Width = 54
	apiEndpointInput.SetValue(cfg.LocalEndpoint)

	// Create custom ISO code input
	customISOInput := textinput.New()
	customISOInput.Placeholder = "e.g., it, ru, zh-cn"
	customISOInput.CharLimit = 10
	customISOInput.Width = 20

	// Create language selector
	targetLangSelector := langselector.NewTargetLanguageSelector(fm)
	targetLangSelector.SetSelectedByCode(cfg.TargetLang)

	// Create model selector
	modelSelector := modelselect.New(fm)

	// Build profile keys
	profileKeys := make([]string, 0)
	for key, profile := range cfg.PromptProfiles {
		if profile.IsFactory {
			profileKeys = append([]string{key}, profileKeys...)
		} else {
			profileKeys = append(profileKeys, key)
		}
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

	// Determine selected log level
	selectedLogLevel := 0
	if cfg.LogLevel == "debug" {
		selectedLogLevel = 1
	}

	// Determine target language selection (0-5 = known, 6 = OTHER)
	selectedTargetLang := 0
	customISO := ""
	normalizedTargetLang := strings.ToLower(strings.TrimSpace(cfg.TargetLang))
	switch normalizedTargetLang {
	case "pt-br", "pt_br", "ptbr":
		selectedTargetLang = 0
	case "en-us", "en_us", "enus", "en":
		selectedTargetLang = 1
	case "es", "es-la", "es_la", "esla":
		selectedTargetLang = 2
	case "ja-jp", "ja_jp", "jajp", "ja":
		selectedTargetLang = 3
	case "fr-fr", "fr_fr", "frfr", "fr":
		selectedTargetLang = 4
	case "de", "de-de", "de_de", "dede":
		selectedTargetLang = 5
	default:
		// Unknown language - set as "OTHER"
		if cfg.TargetLang != "" {
			selectedTargetLang = 6
			customISO = cfg.TargetLang
		}
	}
	customISOInput.SetValue(customISO)

	// Create prompt editor input
	promptInput := textinput.New()
	promptInput.Placeholder = locales.T("settings.prompts.edit_placeholder")
	promptInput.CharLimit = 5000
	promptInput.Width = 60

	// Create profile name input
	profileNameInput := textinput.New()
	profileNameInput.Placeholder = locales.T("settings.prompts.name_placeholder")
	profileNameInput.CharLimit = 50
	profileNameInput.Width = 40

	// Create temperature input
	temperatureInput := textinput.New()
	temperatureInput.Placeholder = "0.0 - 1.0"
	temperatureInput.CharLimit = 5
	temperatureInput.Width = 10

	return Model{
		width:              0, // Will be set by WindowSizeMsg
		height:             0, // Will be set by WindowSizeMsg
		activeTab:          TabGeneral,
		config:             cfg,
		originalCfg:        &origCfg,
		focusManager:       fm,
		targetLangSelector: targetLangSelector,
		selectedTargetLang: selectedTargetLang,
		customISOInput:     customISOInput,
		selectedProvider:   selectedProvider,
		apiKeyInput:        apiKeyInput,
		apiEndpointInput:   apiEndpointInput,
		modelSelector:      modelSelector,
		profileKeys:        profileKeys,
		selectedLogLevel:   selectedLogLevel,
		promptInput:        promptInput,
		profileNameInput:   profileNameInput,
		temperatureInput:   temperatureInput,
	}
}

// SetSize updates the model dimensions
func (m *Model) SetSize(width, height int) {
	m.width = width
	m.height = height
	// Update component widths
	contentWidth := width - 10
	m.targetLangSelector.SetWidth(contentWidth)
	m.modelSelector.SetWidth(contentWidth)
	m.apiKeyInput.Width = min(contentWidth-20, 54)
	m.apiEndpointInput.Width = min(contentWidth-20, 54)
}

func (m Model) Init() tea.Cmd {
	// Request initial window size and start loading models
	return tea.Batch(tea.WindowSize(), m.fetchModelsCmd())
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		// Update component widths
		contentWidth := m.width - 10
		m.targetLangSelector.SetWidth(contentWidth)
		m.modelSelector.SetWidth(contentWidth)
		m.apiKeyInput.Width = min(contentWidth-20, 54)
		m.apiEndpointInput.Width = min(contentWidth-20, 54)
		m.customISOInput.Width = min(contentWidth-40, 20)
		return m, nil

	case modelsLoadedMsg:
		// Models were loaded - update model selector
		if msg.err == nil && len(msg.models) > 0 {
			m.modelSelector.SetModels(msg.models)
			m.modelsLoaded = true
		}
		return m, nil

	case tea.KeyMsg:
		key := msg.String()

		// Handle Touchless modal if open
		if m.showTouchlessModal {
			return m.handleTouchlessModal(msg)
		}

		// Handle input mode for text fields
		if m.focusManager.Mode() == focus.ModeInput {
			return m.handleInputMode(msg)
		}

		// Navigation mode - handle global keys first (these work in all tabs)
		switch key {
		case "ctrl+c":
			return m, func() tea.Msg { return CancelledMsg{} }

		case "q":
			// Q cancels without saving
			return m, func() tea.Msg { return CancelledMsg{} }

		case "esc":
			// ESC saves and exits
			if err := m.saveConfig(); err == nil {
				return m, func() tea.Msg { return SavedMsg{Config: m.config} }
			}
			return m, func() tea.Msg { return CancelledMsg{} }

		case "ctrl+s":
			// Save config in place
			if err := m.saveConfig(); err != nil {
				m.hasError = true
				m.errMsg = err.Error()
			} else {
				m.saved = true
				m.hasError = false
			}
			return m, nil

		case "enter":
			// Enter saves and exits
			if err := m.saveConfig(); err != nil {
				m.hasError = true
				m.errMsg = err.Error()
				return m, nil
			}
			m.saved = true
			return m, func() tea.Msg { return SavedMsg{Config: m.config} }

		case "tab":
			// Tab cycles through tabs to the right
			if m.activeTab < TabAdvanced {
				m.activeTab++
			} else {
				m.activeTab = TabGeneral
			}
			return m, nil

		case "shift+tab":
			// Shift+Tab cycles through tabs to the left
			if m.activeTab > TabGeneral {
				m.activeTab--
			} else {
				m.activeTab = TabAdvanced
			}
			return m, nil

		case "1":
			m.activeTab = TabGeneral
			return m, nil
		case "2":
			m.activeTab = TabProviders
			return m, nil
		case "3":
			m.activeTab = TabModels
			return m, nil
		case "4":
			m.activeTab = TabPrompts
			return m, nil
		case "5":
			m.activeTab = TabAdvanced
			return m, nil

		case "/":
			// "/" triggers search in models tab
			if m.activeTab == TabModels {
				m.focusManager.EnterInput(0)
				m.modelSelector.SetActive(true)
				return m, nil
			}

		default:
			// Tab-specific navigation
			return m.handleTabNavigation(msg)
		}
	}

	return m, cmd
}

// handleInputMode handles key events when a text input is active
func (m Model) handleInputMode(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	key := msg.String()

	switch key {
	case "esc":
		// Exit input mode without saving
		m.focusManager.ExitInput()
		m.apiKeyInput.Blur()
		m.apiEndpointInput.Blur()
		m.customISOInput.Blur()
		m.promptInput.Blur()
		m.profileNameInput.Blur()
		m.temperatureInput.Blur()
		m.editingPrompt = false
		m.editingProfileName = false
		m.targetLangSelector.BlurCustomInput()
		return m, nil
	case "enter":
		// Exit input mode and save
		m.focusManager.ExitInput()
		m.apiKeyInput.Blur()
		m.apiEndpointInput.Blur()
		m.customISOInput.Blur()
		// Save prompt if editing
		if m.editingPrompt && len(m.profileKeys) > m.selectedProfile {
			currentKey := m.profileKeys[m.selectedProfile]
			if profile, exists := m.config.PromptProfiles[currentKey]; exists {
				profile.SystemPrompt = m.promptInput.Value()
				m.config.PromptProfiles[currentKey] = profile
			}
		}
		// Save profile name if editing
		if m.editingProfileName && len(m.profileKeys) > m.selectedProfile {
			currentKey := m.profileKeys[m.selectedProfile]
			if profile, exists := m.config.PromptProfiles[currentKey]; exists {
				newName := strings.TrimSpace(m.profileNameInput.Value())
				if newName != "" {
					profile.Name = newName
					m.config.PromptProfiles[currentKey] = profile
				}
			}
		}
		m.promptInput.Blur()
		m.profileNameInput.Blur()
		m.temperatureInput.Blur()
		m.editingPrompt = false
		m.editingProfileName = false
		m.targetLangSelector.BlurCustomInput()
		return m, nil
	}

	// Forward to the appropriate input
	switch m.activeTab {
	case TabGeneral:
		if m.selectedTargetLang == 6 {
			// Custom ISO input
			m.customISOInput, cmd = m.customISOInput.Update(msg)
			return m, cmd
		} else if m.targetLangSelector.IsCustomSelected() {
			m.targetLangSelector, cmd = m.targetLangSelector.Update(msg)
			return m, cmd
		}
	case TabProviders:
		if m.selectedProvider == 3 {
			m.apiEndpointInput, cmd = m.apiEndpointInput.Update(msg)
			return m, cmd
		} else {
			m.apiKeyInput, cmd = m.apiKeyInput.Update(msg)
			return m, cmd
		}
	case TabPrompts:
		if m.editingPrompt {
			m.promptInput, cmd = m.promptInput.Update(msg)
			return m, cmd
		} else if m.editingProfileName {
			m.profileNameInput, cmd = m.profileNameInput.Update(msg)
			return m, cmd
		} else if m.temperatureInput.Focused() {
			// Handle temperature input
			m.temperatureInput, cmd = m.temperatureInput.Update(msg)
			// Save on enter
			if msg.String() == "enter" {
				if tempVal, err := strconv.ParseFloat(m.temperatureInput.Value(), 64); err == nil {
					if tempVal >= 0 && tempVal <= 1.0 {
						if len(m.profileKeys) > m.selectedProfile {
							currentKey := m.profileKeys[m.selectedProfile]
							if profile, exists := m.config.PromptProfiles[currentKey]; exists {
								profile.Temperature = tempVal
								m.config.PromptProfiles[currentKey] = profile
							}
						}
					}
				}
				m.temperatureInput.Blur()
				m.focusManager.ExitInput()
			}
			return m, cmd
		}
	}

	return m, cmd
}

// handleTouchlessModal handles keyboard input for the touchless configuration modal
func (m Model) handleTouchlessModal(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	key := msg.String()

	switch key {
	case "esc":
		m.showTouchlessModal = false
	case "enter":
		// Save touchless settings to config
		switch m.touchlessMultiSub {
		case 0:
			m.config.TouchlessRules.MultipleSubtitles = "largest"
		case 1:
			m.config.TouchlessRules.MultipleSubtitles = "smallest"
		case 2:
			m.config.TouchlessRules.MultipleSubtitles = "skip"
		}
		switch m.touchlessMuxMode {
		case 0:
			m.config.TouchlessRules.MuxingStrategy = "replace"
		case 1:
			m.config.TouchlessRules.MuxingStrategy = "create_new"
		}
		// Set default profile if available
		if len(m.profileKeys) > m.touchlessProfile {
			m.config.TouchlessRules.DefaultProfile = m.profileKeys[m.touchlessProfile]
		}
		m.showTouchlessModal = false
	case "up", "k":
		// Cycle through options - multi-sub selection
		if m.touchlessMultiSub > 0 {
			m.touchlessMultiSub--
		}
	case "down", "j":
		if m.touchlessMultiSub < 2 {
			m.touchlessMultiSub++
		}
	case "1":
		m.touchlessMultiSub = 0
	case "2":
		m.touchlessMultiSub = 1
	case "3":
		m.touchlessMultiSub = 2
	case "m":
		// Toggle mux mode
		m.touchlessMuxMode = (m.touchlessMuxMode + 1) % 2
	case "p":
		// Cycle profile
		if len(m.profileKeys) > 0 {
			m.touchlessProfile = (m.touchlessProfile + 1) % len(m.profileKeys)
		}
	}

	return m, nil
}

// handleTabNavigation handles navigation within each tab
func (m Model) handleTabNavigation(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	key := msg.String()

	switch m.activeTab {
	case TabGeneral:
		switch key {
		case "up", "k":
			if m.selectedTargetLang > 0 {
				m.selectedTargetLang--
			}
		case "down", "j":
			if m.selectedTargetLang < 6 { // Now includes OTHER option (6)
				m.selectedTargetLang++
			}
		case "h":
			// Toggle HI tags
			m.config.RemoveHITags = !m.config.RemoveHITags
		case "t":
			// Toggle Touchless Mode
			m.config.TouchlessMode = !m.config.TouchlessMode
		case "c":
			// Open Touchless configuration modal
			m.showTouchlessModal = true
			// Initialize touchless settings from config
			switch m.config.TouchlessRules.MultipleSubtitles {
			case "largest":
				m.touchlessMultiSub = 0
			case "smallest":
				m.touchlessMultiSub = 1
			case "skip":
				m.touchlessMultiSub = 2
			default:
				m.touchlessMultiSub = 0
			}
			switch m.config.TouchlessRules.MuxingStrategy {
			case "replace":
				m.touchlessMuxMode = 0
			case "create_new":
				m.touchlessMuxMode = 1
			default:
				m.touchlessMuxMode = 0
			}
		case "u":
			// Toggle auto-update
			m.config.AutoCheckUpdates = !m.config.AutoCheckUpdates
		case "e":
			// Enter edit mode for custom ISO code when OTHER is selected
			if m.selectedTargetLang == 6 {
				m.focusManager.EnterInput(0)
				m.customISOInput.Focus()
				return m, textinput.Blink
			}
		}

	case TabProviders:
		switch key {
		case "up", "k":
			if m.selectedProvider > 0 {
				m.selectedProvider--
			}
		case "down", "j":
			if m.selectedProvider < 3 {
				m.selectedProvider++
			}
		case " ":
			// Toggle API key visibility
			m.showAPIKey = !m.showAPIKey
			if m.showAPIKey {
				m.apiKeyInput.EchoMode = textinput.EchoNormal
			} else {
				m.apiKeyInput.EchoMode = textinput.EchoPassword
			}
		case "e":
			// Enter edit mode for credentials
			m.focusManager.EnterInput(0)
			if m.selectedProvider == 3 {
				m.apiEndpointInput.Focus()
			} else {
				m.apiKeyInput.Focus()
			}
			return m, textinput.Blink
		}

	case TabModels:
		// Model selector handles its own navigation
		// Always activate and forward keys
		m.modelSelector.SetActive(true)
		m.modelSelector, cmd = m.modelSelector.Update(msg)
		return m, cmd

	case TabPrompts:
		switch key {
		case "up", "k":
			if m.showProfileList {
				// Navigate in list
				if m.selectedProfile > 0 {
					m.selectedProfile--
				}
			}
		case "down", "j":
			if m.showProfileList {
				// Navigate in list
				if m.selectedProfile < len(m.profileKeys)-1 {
					m.selectedProfile++
				}
			}
		case " ":
			// Toggle profile list dropdown
			m.showProfileList = !m.showProfileList
		case "n":
			// Edit profile name (only if not locked)
			if len(m.profileKeys) > m.selectedProfile {
				currentKey := m.profileKeys[m.selectedProfile]
				currentProfile := m.config.PromptProfiles[currentKey]
				if !currentProfile.IsFactory && !currentProfile.IsLocked {
					m.editingProfileName = true
					m.profileNameInput.SetValue(currentProfile.Name)
					m.profileNameInput.Focus()
					m.focusManager.EnterInput(0)
					return m, textinput.Blink
				}
			}
		case "t":
			// Edit temperature (only for custom profiles)
			if len(m.profileKeys) > m.selectedProfile {
				currentKey := m.profileKeys[m.selectedProfile]
				currentProfile := m.config.PromptProfiles[currentKey]
				if !currentProfile.IsFactory && !currentProfile.IsLocked {
					m.temperatureInput.SetValue(fmt.Sprintf("%.2f", currentProfile.Temperature))
					m.temperatureInput.Focus()
					m.focusManager.EnterInput(0)
					return m, textinput.Blink
				}
			}
		case "d":
			// Delete profile (only custom profiles)
			if len(m.profileKeys) > m.selectedProfile {
				currentKey := m.profileKeys[m.selectedProfile]
				currentProfile := m.config.PromptProfiles[currentKey]
				if !currentProfile.IsFactory && !currentProfile.IsLocked {
					delete(m.config.PromptProfiles, currentKey)
					// Remove from keys list
					newKeys := []string{}
					for _, k := range m.profileKeys {
						if k != currentKey {
							newKeys = append(newKeys, k)
						}
					}
					m.profileKeys = newKeys
					// Adjust selection
					if m.selectedProfile >= len(m.profileKeys) {
						m.selectedProfile = len(m.profileKeys) - 1
					}
					if m.selectedProfile < 0 {
						m.selectedProfile = 0
					}
				}
			}
		case "e":
			// Edit current profile (only if not locked)
			if len(m.profileKeys) > m.selectedProfile {
				currentKey := m.profileKeys[m.selectedProfile]
				currentProfile := m.config.PromptProfiles[currentKey]
				if !currentProfile.IsFactory && !currentProfile.IsLocked {
					m.editingPrompt = true
					m.promptInput.SetValue(currentProfile.SystemPrompt)
					m.promptInput.Focus()
					m.focusManager.EnterInput(0)
					return m, textinput.Blink
				}
			}
		case "c":
			// Clone current profile
			if len(m.profileKeys) > m.selectedProfile {
				currentKey := m.profileKeys[m.selectedProfile]
				currentProfile := m.config.PromptProfiles[currentKey]
				// Create new profile as clone
				newKey := currentKey + "_custom"
				counter := 1
				for {
					if _, exists := m.config.PromptProfiles[newKey]; !exists {
						break
					}
					counter++
					newKey = fmt.Sprintf("%s_custom_%d", currentKey, counter)
				}
				newProfile := config.PromptProfile{
					Name:         currentProfile.Name + " (Clone)",
					SystemPrompt: currentProfile.SystemPrompt,
					Temperature:  currentProfile.Temperature,
					IsFactory:    false,
					IsLocked:     false,
				}
				m.config.PromptProfiles[newKey] = newProfile
				m.profileKeys = append(m.profileKeys, newKey)
				m.selectedProfile = len(m.profileKeys) - 1
			}
		}

	case TabAdvanced:
		switch key {
		case "up", "k":
			if m.selectedLogLevel > 0 {
				m.selectedLogLevel--
			}
		case "down", "j":
			if m.selectedLogLevel < 1 {
				m.selectedLogLevel++
			}
		case " ":
			// Toggle auto update check
			m.config.AutoCheckUpdates = !m.config.AutoCheckUpdates
		}
	}

	return m, cmd
}

// saveConfig saves the current configuration
func (m *Model) saveConfig() error {
	// Update config from UI state

	// General tab - Target Language
	targetLangs := []string{"PT-BR", "EN-US", "ES", "JA-JP", "FR-FR", "DE"}
	if m.selectedTargetLang >= 0 && m.selectedTargetLang < len(targetLangs) {
		m.config.TargetLang = targetLangs[m.selectedTargetLang]
	} else if m.selectedTargetLang == 6 {
		// OTHER - use custom ISO code
		customISO := strings.TrimSpace(m.customISOInput.Value())
		if customISO != "" {
			m.config.TargetLang = strings.ToLower(customISO)
		}
	}

	// Providers tab
	providers := []string{"openrouter", "gemini", "openai", "local"}
	if m.selectedProvider >= 0 && m.selectedProvider < len(providers) {
		m.config.AIProvider = providers[m.selectedProvider]
	}
	m.config.APIKey = m.apiKeyInput.Value()
	m.config.LocalEndpoint = m.apiEndpointInput.Value()

	// Models tab
	if selectedModel := m.modelSelector.GetSelectedModel(); selectedModel != nil {
		m.config.Model = selectedModel.ID
	}

	// Advanced tab
	if m.selectedLogLevel == 1 {
		m.config.LogLevel = "debug"
	} else {
		m.config.LogLevel = "info"
	}

	return m.config.Save()
}

func (m Model) View() string {
	// Wait for terminal size
	if layout.IsWaitingForSize(m.width, m.height) {
		return locales.T("common.loading")
	}

	// Use actual terminal dimensions
	termWidth := m.width
	termHeight := m.height
	if termWidth < 60 {
		termWidth = 60
	}
	if termHeight < 20 {
		termHeight = 20
	}

	// Check if terminal is too small
	if layout.IsTooSmall(m.width, m.height) {
		return layout.RenderTooSmallWarning(m.width, m.height)
	}

	// Content width is terminal width minus margins
	contentWidth := termWidth - 4

	header := m.renderHeader(contentWidth)
	tabs := m.renderTabs(contentWidth)
	content := m.renderTabContent(contentWidth)
	footer := m.renderFooter(contentWidth)

	view := lipgloss.JoinVertical(
		lipgloss.Left,
		header,
		tabs,
		"",
		content,
		"",
		footer,
	)

	// Wrap in main window using full terminal width
	baseView := styles.MainWindow.Width(termWidth - 2).Render(view)

	// Overlay touchless modal if open
	if m.showTouchlessModal {
		return m.renderWithOverlay(baseView, m.renderTouchlessModal())
	}

	return baseView
}

func (m Model) renderHeader(contentWidth int) string {
	title := locales.T("settings.title")
	barWidth := contentWidth - len(title) - 4
	if barWidth < 10 {
		barWidth = 10
	}
	if barWidth > 70 {
		barWidth = 70
	}

	headerBar := strings.Repeat("▒", barWidth)
	headerLine := " " + title + " " + headerBar

	return styles.HeaderBorder.Render(headerLine)
}

func (m Model) renderTabs(contentWidth int) string {
	// Define short tab names for compact display
	tabNames := []string{"GERAL", "PROV.", "MODELOS", "PROMPTS", "AVANÇ."}
	currentLang := locales.GetCurrentLocale()
	if currentLang == "en" {
		tabNames = []string{"GENERAL", "PROVID.", "MODELS", "PROMPTS", "ADVANC."}
	} else if currentLang == "es" {
		tabNames = []string{"GENERAL", "PROV.", "MODELOS", "PROMPTS", "AVANZ."}
	}

	var tabs strings.Builder
	for i, name := range tabNames {
		shortcut := fmt.Sprintf("[%d]", i+1)
		label := shortcut + " " + name

		if Tab(i) == m.activeTab {
			tabs.WriteString(styles.Highlight.Render("[ " + label + " ]"))
		} else {
			tabs.WriteString(styles.Dimmed.Render("[ " + label + " ]"))
		}
	}

	return tabs.String()
}

func (m Model) renderTabContent(contentWidth int) string {
	// Panel takes full available width
	panelWidth := contentWidth - 2
	if panelWidth < 50 {
		panelWidth = 50
	}

	switch m.activeTab {
	case TabGeneral:
		return m.renderGeneralTab(panelWidth)
	case TabProviders:
		return m.renderProvidersTab(panelWidth)
	case TabModels:
		return m.renderModelsTab(panelWidth)
	case TabPrompts:
		return m.renderPromptsTab(panelWidth)
	case TabAdvanced:
		return m.renderAdvancedTab(panelWidth)
	}
	return ""
}

func (m Model) renderGeneralTab(panelWidth int) string {
	// Target language section
	var langList strings.Builder
	targetLangs := []string{
		"PT-BR (Português)",
		"EN-US (English)",
		"ES (Español)",
		"JA-JP (Japanese)",
		"FR-FR (Français)",
		"DE (Deutsch)",
	}

	for i, lang := range targetLangs {
		icon := "( )"
		if i == m.selectedTargetLang {
			icon = "(o)"
		}
		langList.WriteString(fmt.Sprintf("   %s %s\n", styles.Highlight.Render(icon), lang))
	}

	// OTHER option with custom ISO input
	otherIcon := "( )"
	if m.selectedTargetLang == 6 {
		otherIcon = "(o)"
	}
	customValue := m.customISOInput.Value()
	if customValue == "" {
		customValue = "_______"
	}
	otherLabel := locales.T("wizard.step3.other")
	langList.WriteString(fmt.Sprintf("   %s %s: [%s]", styles.Highlight.Render(otherIcon), otherLabel, customValue))
	if m.selectedTargetLang == 6 {
		langList.WriteString("  " + styles.KeyHintStyle.Render("[E] "+locales.T("common.edit")))
	}
	langList.WriteString("\n")

	// HI tags checkbox
	hiCheckbox := "[ ]"
	if m.config.RemoveHITags {
		hiCheckbox = "[X]"
	}

	// Touchless mode checkbox
	touchlessCheckbox := "[ ]"
	if m.config.TouchlessMode {
		touchlessCheckbox = "[X]"
	}

	// Auto update checkbox
	autoUpdateCheckbox := "[ ]"
	if m.config.AutoCheckUpdates {
		autoUpdateCheckbox = "[X]"
	}

	// Temperature - use GlobalTemperature if set, otherwise Temperature
	temp := m.config.GlobalTemperature
	if temp == 0 {
		temp = m.config.Temperature
	}
	if temp == 0 {
		temp = 0.3 // default
	}

	content := lipgloss.JoinVertical(
		lipgloss.Left,
		styles.PanelTitle.Render(locales.T("settings.general.target_lang")),
		langList.String(),
		"",
		styles.PanelTitle.Render(locales.T("settings.general.preferences")),
		fmt.Sprintf("   %s %s  %s", hiCheckbox, locales.T("settings.general.remove_hi"), styles.KeyHintStyle.Render("[H]")),
		fmt.Sprintf("   %s %s  %s", touchlessCheckbox, locales.T("settings.general.touchless_mode"), styles.KeyHintStyle.Render("[T]")),
		"      "+styles.Dimmed.Render(locales.T("settings.general.touchless_warning")),
		"      "+styles.KeyHintStyle.Render("[C] "+locales.T("settings.general.configure_rules")),
		fmt.Sprintf("   %s %s  %s", autoUpdateCheckbox, locales.T("settings.general.auto_update"), styles.KeyHintStyle.Render("[U]")),
		fmt.Sprintf("   %s: [ %.1f ]", locales.T("settings.general.temperature"), temp),
	)

	return styles.Panel.Width(panelWidth).BorderForeground(styles.Yellow).Render(content)
}

func (m Model) renderProvidersTab(panelWidth int) string {
	providers := []string{
		"♾️  OpenRouter (Recommended)",
		"💎 Google Gemini API",
		"🤖 OpenAI API",
		"🏠 Local LLM (Ollama/LMStudio)",
	}

	var providerList strings.Builder
	for i, prov := range providers {
		icon := "( )"
		if i == m.selectedProvider {
			icon = "(o)"
		}
		providerList.WriteString(fmt.Sprintf("   %s %s\n", styles.Highlight.Render(icon), prov))
	}

	// Credentials section
	var credContent strings.Builder
	if m.selectedProvider == 3 {
		// Local LLM - endpoint
		credContent.WriteString(locales.T("settings.providers.endpoint") + ":\n")
		credContent.WriteString("   " + m.apiEndpointInput.View() + "\n")
	} else {
		// API Key
		credContent.WriteString(locales.T("settings.providers.api_key") + ":\n")
		credContent.WriteString("   " + m.apiKeyInput.View() + "\n")
		showText := locales.T("settings.providers.show")
		if m.showAPIKey {
			showText = locales.T("settings.providers.hide")
		}
		credContent.WriteString("   " + styles.KeyHintStyle.Render("[SPACE] "+showText+" | [E] "+locales.T("common.edit")) + "\n")
	}

	content := lipgloss.JoinVertical(
		lipgloss.Left,
		styles.PanelTitle.Render(locales.T("settings.providers.select")),
		providerList.String(),
		"",
		styles.PanelTitle.Render(locales.T("settings.providers.credentials")),
		credContent.String(),
	)

	return styles.Panel.Width(panelWidth).BorderForeground(styles.Yellow).Render(content)
}

func (m Model) renderModelsTab(panelWidth int) string {
	m.modelSelector.SetWidth(panelWidth)
	m.modelSelector.SetActive(m.activeTab == TabModels)

	content := m.modelSelector.RenderWithInfo()

	return styles.Panel.Width(panelWidth).BorderForeground(styles.Yellow).Render(content)
}

func (m Model) renderPromptsTab(panelWidth int) string {
	// Profile selector/dropdown
	var dropdown strings.Builder
	dropdown.WriteString(locales.T("settings.prompts.select_profile") + ":\n")

	if m.showProfileList {
		// Show all profiles in list
		dropdown.WriteString("\n")
		for i, key := range m.profileKeys {
			profile := m.config.PromptProfiles[key]
			icon := "🔒"
			if !profile.IsFactory {
				icon = "👤"
			}
			if i == m.selectedProfile {
				dropdown.WriteString("   " + styles.Highlight.Render("→ "+icon+" "+profile.Name) + "\n")
			} else {
				dropdown.WriteString("     " + icon + " " + profile.Name + "\n")
			}
		}
		dropdown.WriteString("\n   " + styles.KeyHintStyle.Render("[SPACE] "+locales.T("settings.prompts.hide_list")) + "\n")
	} else {
		// Show selected profile only
		if len(m.profileKeys) > m.selectedProfile {
			currentKey := m.profileKeys[m.selectedProfile]
			currentProfile := m.config.PromptProfiles[currentKey]
			icon := "🔒"
			if !currentProfile.IsFactory {
				icon = "👤"
			}
			dropdown.WriteString(fmt.Sprintf("   [ %s %s           ▼ ]\n", icon, currentProfile.Name))
		}
		dropdown.WriteString("   " + styles.KeyHintStyle.Render("[SPACE] "+locales.T("settings.prompts.show_list")) + "\n")
	}

	status := ""
	if len(m.profileKeys) > m.selectedProfile {
		currentProfile := m.config.PromptProfiles[m.profileKeys[m.selectedProfile]]
		if currentProfile.IsLocked {
			status = styles.StatusWarning.Render("["+locales.T("settings.prompts.locked")+"]") + " " + locales.T("settings.prompts.locked_hint")
		} else {
			status = styles.StatusOK.Render("["+locales.T("settings.prompts.editable")+"]") + " " + locales.T("settings.prompts.editable_hint")
		}
	}

	// Show name edit input if editing name
	var nameEditSection string
	if m.editingProfileName {
		nameEditSection = lipgloss.JoinVertical(
			lipgloss.Left,
			"",
			styles.PanelTitle.Render(locales.T("settings.prompts.edit_name")),
			"   "+m.profileNameInput.View(),
			"   "+styles.KeyHintStyle.Render("[ENTER] "+locales.T("common.save")+"  [ESC] "+locales.T("common.cancel")),
			"",
		)
	}

	// Show edit input or preview
	var previewSection string
	if m.editingPrompt {
		// Show editable input
		previewSection = lipgloss.JoinVertical(
			lipgloss.Left,
			styles.PanelTitle.Render(locales.T("settings.prompts.title")+" - "+locales.T("common.edit")),
			"   "+m.promptInput.View(),
			"   "+styles.KeyHintStyle.Render("[ENTER] "+locales.T("common.save")+"  [ESC] "+locales.T("common.cancel")),
		)
	} else {
		preview := ""
		if len(m.profileKeys) > m.selectedProfile {
			currentProfile := m.config.PromptProfiles[m.profileKeys[m.selectedProfile]]
			lines := strings.Split(currentProfile.SystemPrompt, "\n")
			if len(lines) > 4 {
				lines = lines[:4]
			}
			preview = strings.Join(lines, "\n") + "\n" + styles.Dimmed.Render("...")
		}
		previewSection = lipgloss.JoinVertical(
			lipgloss.Left,
			styles.PanelTitle.Render(locales.T("settings.prompts.preview")),
			styles.CodeBlock.Render(preview),
		)
	}

	// Actions based on profile type
	actions := ""
	if len(m.profileKeys) > m.selectedProfile {
		currentProfile := m.config.PromptProfiles[m.profileKeys[m.selectedProfile]]
		if currentProfile.IsFactory || currentProfile.IsLocked {
			actions = styles.KeyHintStyle.Render("[C] " + locales.T("settings.prompts.clone_to_edit"))
		} else {
			actions = styles.KeyHintStyle.Render("[N] " + locales.T("settings.prompts.edit_name") + "  [T] " + locales.T("settings.prompts.edit_temp") + "  [E] " + locales.T("common.edit") + "  [C] " + locales.T("settings.prompts.clone") + "  [D] " + locales.T("settings.prompts.delete"))
		}
	}

	// Show temperature info
	var tempSection string
	if len(m.profileKeys) > m.selectedProfile {
		currentProfile := m.config.PromptProfiles[m.profileKeys[m.selectedProfile]]
		if currentProfile.IsFactory || currentProfile.IsLocked {
			// Factory profiles use global temperature
			globalTemp := m.config.GlobalTemperature
			if globalTemp == 0 {
				globalTemp = m.config.Temperature
			}
			if globalTemp == 0 {
				globalTemp = 0.3
			}
			tempSection = fmt.Sprintf("\n   %s: %.2f (%s)", locales.T("settings.prompts.temperature"), globalTemp, locales.T("settings.prompts.uses_global"))
		} else {
			// Custom profiles have their own temperature
			temp := currentProfile.Temperature
			if temp == 0 {
				temp = 0.3
			}
			tempSection = fmt.Sprintf("\n   %s: %.2f", locales.T("settings.prompts.temperature"), temp)
		}
	}

	content := lipgloss.JoinVertical(
		lipgloss.Left,
		styles.PanelTitle.Render(locales.T("settings.prompts.title")),
		dropdown.String(),
		status,
		tempSection,
		actions,
		nameEditSection,
		previewSection,
	)

	return styles.Panel.Width(panelWidth).BorderForeground(styles.Yellow).Render(content)
}

func (m Model) renderAdvancedTab(panelWidth int) string {
	logLevels := []string{"INFO (Default)", "DEBUG (Verbose)"}
	var logList strings.Builder
	for i, level := range logLevels {
		icon := "( )"
		if i == m.selectedLogLevel {
			icon = "(o)"
		}
		logList.WriteString(fmt.Sprintf("   %s %s\n", icon, level))
	}

	updatesCheck := "[ ]"
	if m.config.AutoCheckUpdates {
		updatesCheck = "[X]"
	}

	content := lipgloss.JoinVertical(
		lipgloss.Left,
		styles.PanelTitle.Render(locales.T("settings.advanced.log_level")),
		logList.String(),
		"",
		styles.PanelTitle.Render(locales.T("settings.advanced.updates")),
		fmt.Sprintf("   %s %s  %s", updatesCheck, locales.T("settings.advanced.auto_check"), styles.KeyHintStyle.Render("[SPACE]")),
		"",
		styles.PanelTitle.Render(locales.T("settings.advanced.system_info")),
		"   VERSION: v1.0.0",
		"   GO VERSION: 1.24.0",
	)

	return styles.Panel.Width(panelWidth).BorderForeground(styles.Yellow).Render(content)
}

func (m Model) renderFooter(contentWidth int) string {
	saveStatus := ""
	if m.saved {
		saveStatus = styles.StatusOK.Render("[" + locales.T("common.saved") + "] ")
	}
	if m.hasError {
		saveStatus = styles.StatusError.Render("[ERROR: " + m.errMsg + "] ")
	}

	footerText := saveStatus +
		styles.RenderHotkey("ESC/ENTER", locales.T("common.save_exit")) + "  " +
		styles.RenderHotkey("TAB", locales.T("settings.switch_tab")) + "  " +
		styles.RenderHotkey("↑/↓", locales.T("common.navigate"))

	return styles.Footer.Render(footerText)
}

// fetchModelsCmd creates a command to fetch models from the provider
func (m Model) fetchModelsCmd() tea.Cmd {
	return func() tea.Msg {
		ctx := context.Background()

		var provider ai.LLMProvider
		var err error

		switch m.selectedProvider {
		case 0: // OpenRouter
			apiKey := m.apiKeyInput.Value()
			if apiKey == "" {
				apiKey = m.config.APIKey
			}
			if apiKey == "" {
				// Return sample data if no API key
				return modelsLoadedMsg{models: getSampleModels(), err: nil}
			}
			provider = ai.NewOpenRouterAdapter(apiKey, "", 0.3)
			modelStrings, err := provider.ListModels(ctx)
			if err != nil {
				return modelsLoadedMsg{models: getSampleModels(), err: nil}
			}
			models := []modelselect.ModelInfo{}
			for _, model := range modelStrings {
				isFree := strings.Contains(strings.ToLower(model), ":free")
				models = append(models, parseModelToModelInfo(model, isFree))
			}
			return modelsLoadedMsg{models: models, err: nil}

		case 1: // Gemini
			apiKey := m.apiKeyInput.Value()
			if apiKey == "" {
				apiKey = m.config.APIKey
			}
			provider, err = ai.NewGeminiAdapter(ctx, apiKey, "", 0.3)
			if err != nil {
				return modelsLoadedMsg{models: nil, err: err}
			}
			modelStrings, err := provider.ListModels(ctx)
			if err != nil {
				return modelsLoadedMsg{models: nil, err: err}
			}
			models := []modelselect.ModelInfo{}
			for _, model := range modelStrings {
				models = append(models, parseModelToModelInfo(model, false))
			}
			return modelsLoadedMsg{models: models, err: nil}

		case 2: // OpenAI
			modelStrings := []string{"gpt-4o", "gpt-4o-mini", "gpt-4-turbo", "gpt-3.5-turbo"}
			models := []modelselect.ModelInfo{}
			for _, model := range modelStrings {
				models = append(models, parseModelToModelInfo(model, false))
			}
			return modelsLoadedMsg{models: models, err: nil}

		case 3: // Local
			endpoint := m.apiEndpointInput.Value()
			if endpoint == "" {
				endpoint = m.config.LocalEndpoint
			}
			if endpoint == "" {
				endpoint = "http://localhost:11434"
			}
			provider = ai.NewLocalLLMAdapter(endpoint, "", 0.3)
			modelStrings, err := provider.ListModels(ctx)
			if err != nil {
				return modelsLoadedMsg{models: nil, err: err}
			}
			models := []modelselect.ModelInfo{}
			for _, model := range modelStrings {
				models = append(models, parseModelToModelInfo(model, true))
			}
			return modelsLoadedMsg{models: models, err: nil}
		}

		return modelsLoadedMsg{models: nil, err: fmt.Errorf("provider not configured")}
	}
}

// parseModelToModelInfo converts a model string to ModelInfo
func parseModelToModelInfo(model string, isFree bool) modelselect.ModelInfo {
	// Check if model is in format: id|price|context (from OpenRouter API)
	pipeParts := strings.Split(model, "|")
	modelID := model
	contextSize := "32k"
	price := "$0.10"

	if len(pipeParts) >= 3 {
		// Format: id|price|context
		modelID = pipeParts[0]
		priceStr := pipeParts[1]
		contextStr := pipeParts[2]

		// Parse context length
		if ctx, err := strconv.Atoi(contextStr); err == nil {
			if ctx >= 1000000 {
				contextSize = fmt.Sprintf("%dM", ctx/1000000)
			} else if ctx >= 1000 {
				contextSize = fmt.Sprintf("%dk", ctx/1000)
			} else {
				contextSize = fmt.Sprintf("%d", ctx)
			}
		}

		// Parse price (format: "0.00000150" per token)
		if priceVal, err := strconv.ParseFloat(priceStr, 64); err == nil {
			if priceVal == 0 {
				price = "FREE"
				isFree = true
			} else {
				// Convert per-token price to per-million tokens
				pricePerM := priceVal * 1000000
				if pricePerM < 0.01 {
					price = "FREE"
					isFree = true
				} else {
					price = fmt.Sprintf("$%.2f", pricePerM)
				}
			}
		}
	}

	// Extract name from model ID (format: provider/model-name)
	slashParts := strings.Split(modelID, "/")
	name := modelID
	provider := ""
	if len(slashParts) > 1 {
		provider = slashParts[0]
		name = slashParts[len(slashParts)-1]
	}

	// Remove :free suffix from name
	name = strings.TrimSuffix(name, ":free")
	if strings.Contains(strings.ToLower(modelID), ":free") {
		isFree = true
		price = "FREE"
	}

	// Clean up the name for display
	name = strings.ReplaceAll(name, "-", " ")
	name = strings.Title(name)

	// Format as provider/name if provider exists
	if provider != "" {
		name = provider + "/" + name
	}

	return modelselect.ModelInfo{
		ID:          modelID,
		Name:        name,
		ContextSize: contextSize,
		PricePerM:   price,
		IsFree:      isFree,
	}
}

// getSampleModels returns sample models for testing
func getSampleModels() []modelselect.ModelInfo {
	return []modelselect.ModelInfo{
		{ID: "meta-llama/llama-3.3-70b-instruct:free", Name: "Llama 3.3 70B", ContextSize: "128k", PricePerM: "FREE", IsFree: true},
		{ID: "qwen/qwen-2.5-72b-instruct:free", Name: "Qwen 2.5 72B", ContextSize: "32k", PricePerM: "FREE", IsFree: true},
		{ID: "mistralai/mistral-7b-instruct:free", Name: "Mistral 7B", ContextSize: "32k", PricePerM: "FREE", IsFree: true},
		{ID: "google/gemma-2-9b-it:free", Name: "Gemma 2 9B", ContextSize: "8k", PricePerM: "FREE", IsFree: true},
		{ID: "anthropic/claude-3-opus", Name: "Claude 3 Opus", ContextSize: "200k", PricePerM: "$15.00", IsFree: false},
		{ID: "anthropic/claude-3-sonnet", Name: "Claude 3 Sonnet", ContextSize: "200k", PricePerM: "$3.00", IsFree: false},
		{ID: "openai/gpt-4o", Name: "GPT-4o", ContextSize: "128k", PricePerM: "$2.50", IsFree: false},
		{ID: "openai/gpt-4o-mini", Name: "GPT-4o Mini", ContextSize: "128k", PricePerM: "$0.15", IsFree: false},
		{ID: "google/gemini-1.5-pro", Name: "Gemini 1.5 Pro", ContextSize: "1M", PricePerM: "$1.25", IsFree: false},
		{ID: "google/gemini-2.0-flash-exp:free", Name: "Gemini 2.0 Flash", ContextSize: "1M", PricePerM: "FREE", IsFree: true},
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func (m Model) renderWithOverlay(base, modal string) string {
	return base + "\n\n" + modal
}

func (m Model) renderTouchlessModal() string {
	var s strings.Builder

	title := locales.T("settings.touchless.title")
	s.WriteString(styles.TitleStyle.Render(title))
	s.WriteString("\n\n")

	// Multiple subtitles strategy
	s.WriteString(styles.SectionStyle.Render(locales.T("settings.touchless.multi_sub")))
	s.WriteString("\n")

	strategies := []string{
		locales.T("settings.touchless.largest"),
		locales.T("settings.touchless.smallest"),
		locales.T("settings.touchless.skip"),
	}

	for i, strat := range strategies {
		icon := "( )"
		if i == m.touchlessMultiSub {
			icon = "(o)"
		}
		s.WriteString(fmt.Sprintf("  %s %s\n", styles.Highlight.Render(icon), strat))
	}

	s.WriteString("\n")

	// Default context profile
	s.WriteString(styles.SectionStyle.Render(locales.T("settings.touchless.default_profile")))
	s.WriteString("\n")
	profileName := "Anime (Factory)"
	if len(m.profileKeys) > m.touchlessProfile {
		profileName = m.profileKeys[m.touchlessProfile]
	}
	s.WriteString(fmt.Sprintf("  [ %s ]  %s\n", profileName, styles.KeyHintStyle.Render("[P] cycle")))

	s.WriteString("\n")

	// Muxing strategy
	s.WriteString(styles.SectionStyle.Render(locales.T("settings.touchless.mux_strategy")))
	s.WriteString("\n")

	muxStrategies := []string{
		locales.T("settings.touchless.replace"),
		locales.T("settings.touchless.new_file"),
	}

	for i, strat := range muxStrategies {
		icon := "( )"
		if i == m.touchlessMuxMode {
			icon = "(o)"
		}
		s.WriteString(fmt.Sprintf("  %s %s\n", styles.Highlight.Render(icon), strat))
	}

	footer := "\n" + styles.KeyHintStyle.Render("[ESC]") + " Cancel  " + styles.KeyHintStyle.Render("[ENTER]") + " Save"
	s.WriteString(footer)

	return styles.ModalStyle.Width(60).Render(s.String())
}
