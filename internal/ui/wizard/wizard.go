package wizard

import (
	"context"
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/lsilvatti/bakasub/internal/config"
	"github.com/lsilvatti/bakasub/internal/core/ai"
	"github.com/lsilvatti/bakasub/internal/core/dependencies"
	"github.com/lsilvatti/bakasub/internal/locales"
	"github.com/lsilvatti/bakasub/internal/ui/components"
	"github.com/lsilvatti/bakasub/internal/ui/focus"
	"github.com/lsilvatti/bakasub/internal/ui/layout"
	"github.com/lsilvatti/bakasub/internal/ui/styles"
)

// isoLanguageNames maps ISO language codes to human-readable names
var isoLanguageNames = map[string]string{
	"pt-br": "Português (Brasil)",
	"pt":    "Português",
	"en":    "English",
	"en-us": "English (US)",
	"en-gb": "English (UK)",
	"es":    "Español",
	"es-la": "Español (Latinoamérica)",
	"es-es": "Español (España)",
	"fr":    "Français",
	"fr-fr": "Français (France)",
	"de":    "Deutsch",
	"it":    "Italiano",
	"ja":    "日本語 (Japanese)",
	"ja-jp": "日本語 (Japanese)",
	"ko":    "한국어 (Korean)",
	"ko-kr": "한국어 (Korean)",
	"zh":    "中文 (Chinese)",
	"zh-cn": "简体中文 (Simplified Chinese)",
	"zh-tw": "繁體中文 (Traditional Chinese)",
	"ru":    "Русский (Russian)",
	"ar":    "العربية (Arabic)",
	"hi":    "हिन्दी (Hindi)",
	"th":    "ไทย (Thai)",
	"vi":    "Tiếng Việt (Vietnamese)",
	"nl":    "Nederlands (Dutch)",
	"pl":    "Polski (Polish)",
	"tr":    "Türkçe (Turkish)",
	"sv":    "Svenska (Swedish)",
	"da":    "Dansk (Danish)",
	"no":    "Norsk (Norwegian)",
	"fi":    "Suomi (Finnish)",
	"cs":    "Čeština (Czech)",
	"hu":    "Magyar (Hungarian)",
	"ro":    "Română (Romanian)",
	"el":    "Ελληνικά (Greek)",
	"he":    "עברית (Hebrew)",
	"id":    "Bahasa Indonesia",
	"ms":    "Bahasa Melayu (Malay)",
	"uk":    "Українська (Ukrainian)",
}

// ModelInfo contains detailed information about an AI model
type ModelInfo struct {
	ID          string
	Name        string
	ContextSize string
	PricePerM   string
	IsFree      bool
}

// Step represents the current wizard step
type Step int

const (
	StepLanguageDeps Step = iota // Step 1: UI Language + Dependencies
	StepProvider                 // Step 2: Provider + API Key + Model Selection
	StepDefaults                 // Step 3: Target Language + Preferences
)

// Model represents the wizard state
type Model struct {
	step     Step
	width    int
	height   int
	config   *config.Config
	quitting bool
	finished bool

	// Focus manager
	focusManager *focus.Manager

	// Step 1: UI Language + Dependencies
	uiLanguageSelection int // 0=English, 1=PT-BR, 2=Español
	depStatus           map[string]bool
	depDownloading      bool
	currentDep          string
	downloadProgress    float64
	depError            string
	checkComplete       bool
	depSpinner          components.NeonSpinner

	// Step 2: Provider + API Key + Model Selection
	activeSection     int // 0=Provider, 1=Credentials, 2=ModelSelect
	providerSelection int
	apiKeyInput       textinput.Model
	apiEndpointInput  textinput.Model
	showAPIKey        bool
	keyValidated      bool
	modelSelection    int
	modelTab          int // 0=FREE, 1=ALL MODELS
	modelSearchInput  textinput.Model
	availableModels   []ModelInfo
	filteredModels    []ModelInfo
	loadingModels     bool
	modelScrollOffset int // For scrolling through long lists

	// Step 3: Target Language + Preferences
	activeStep3Section int // 0=Language, 1=Preferences
	languageSelection  int
	targetLangOther    bool            // true if "OTHER" selected
	customLangCode     textinput.Model // for custom ISO code
	hiTagsRemoval      bool
	tempValue          float64
}

// New creates a new wizard
func New(cfg *config.Config) Model {
	// Don't reload locales - use whatever was initialized by the main app
	// This preserves the language that was loaded by locales.Init()

	apiKey := textinput.New()
	apiKey.Placeholder = "sk-or-v1-... or AIza..."
	apiKey.CharLimit = 200
	apiKey.Width = 54 // Consistent width
	apiKey.EchoMode = textinput.EchoPassword
	apiKey.EchoCharacter = '•'
	apiKey.Focus() // Start focused

	apiEndpoint := textinput.New()
	apiEndpoint.Placeholder = "http://localhost:11434"
	apiEndpoint.CharLimit = 100
	apiEndpoint.Width = 54 // Consistent width
	apiEndpoint.Focus()    // Start focused

	customLang := textinput.New()
	customLang.Placeholder = "pt-br, en-us, ja-jp, etc."
	customLang.CharLimit = 10
	customLang.Width = 30

	modelSearch := textinput.New()
	modelSearch.Placeholder = "Search models..."
	modelSearch.CharLimit = 50
	modelSearch.Width = 54

	fm := focus.NewManager(1) // 1 text input field at a time

	// Detect current language and set the correct selection index
	currentLang := locales.GetCurrentLocale()
	uiLangIndex := 0 // Default to English
	switch currentLang {
	case "en":
		uiLangIndex = 0
	case "pt-br":
		uiLangIndex = 1
	case "es":
		uiLangIndex = 2
	}

	return Model{
		step:                StepLanguageDeps,
		config:              cfg,
		focusManager:        fm,
		uiLanguageSelection: uiLangIndex, // Set based on current locale
		depSpinner:          components.NewNeonSpinner(),
		depStatus:           make(map[string]bool),
		activeSection:       0, // Start at Provider section (for Step 2)
		providerSelection:   0,
		apiKeyInput:         apiKey,
		apiEndpointInput:    apiEndpoint,
		modelTab:            0, // Start on FREE tab
		modelSearchInput:    modelSearch,
		availableModels:     []ModelInfo{},
		filteredModels:      []ModelInfo{},
		modelScrollOffset:   0,
		activeStep3Section:  0,
		languageSelection:   0, // PT-BR default
		customLangCode:      customLang,
		hiTagsRemoval:       true,
		tempValue:           0.3,
	}
}

// Init initializes the wizard
func (m Model) Init() tea.Cmd {
	// Start dependency check immediately on Step 1
	spinnerCmd := m.depSpinner.Start()
	return tea.Batch(textinput.Blink, checkDependencies, spinnerCmd)
}

// Update handles messages
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	// Update spinner if active
	var spinnerCmd tea.Cmd
	m.depSpinner, spinnerCmd = m.depSpinner.Update(msg)
	if spinnerCmd != nil {
		cmds = append(cmds, spinnerCmd)
	}

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, tea.Batch(cmds...)

	case tea.KeyMsg:
		// GATEKEEPER LOGIC: Route keys based on focus mode
		if m.focusManager.Mode() == focus.ModeInput {
			// In input mode - ESC exits (except for Step 2 where ENTER toggles)
			if msg.String() == "esc" {
				// Only allow ESC to exit if NOT in Step 2 (Provider)
				if m.step != StepProvider {
					m.focusManager.ExitInput()
					// Blur all inputs when exiting
					m.apiKeyInput.Blur()
					m.apiEndpointInput.Blur()
					m.modelSearchInput.Blur()
					m.customLangCode.Blur()
					return m, tea.Batch(cmds...)
				}
			}
			// In Step 2, ENTER toggles edit mode (handled in handleKeyPress)
			if msg.String() == "enter" && m.step == StepProvider {
				// Exit input mode
				m.focusManager.ExitInput()
				m.apiKeyInput.Blur()
				m.apiEndpointInput.Blur()
				m.modelSearchInput.Blur()
				// Fetch models after entering API key
				if m.apiKeyInput.Value() != "" || m.apiEndpointInput.Value() != "" {
					m.loadingModels = true
					cmds = append(cmds, m.fetchModelsCmd())
				}
				// Filter models after search
				m.filterModels()
				return m, tea.Batch(cmds...)
			}
			// Forward to active input
			model, cmd := m.updateActiveInput(msg)
			m = model.(Model)
			if cmd != nil {
				cmds = append(cmds, cmd)
			}
			return m, tea.Batch(cmds...)
		}

		// In nav mode - handle navigation and global hotkeys
		model, cmd := m.handleKeyPress(msg)
		m = model.(Model)
		if cmd != nil {
			cmds = append(cmds, cmd)
		}
		return m, tea.Batch(cmds...)

	case checkDepsMsg:
		m.depStatus = msg.status
		m.checkComplete = true
		m.depSpinner.Stop() // Stop spinner when check completes
		return m, tea.Batch(cmds...)

	case downloadProgressMsg:
		m.downloadProgress = msg.progress
		m.currentDep = msg.depName
		return m, tea.Batch(cmds...)

	case downloadCompleteMsg:
		m.depDownloading = false
		m.depSpinner.Stop() // Stop spinner when download completes
		if msg.err != nil {
			m.depError = msg.err.Error()
		} else {
			// Recheck dependencies
			cmds = append(cmds, checkDependencies)
		}
		return m, tea.Batch(cmds...)

	case finishMsg:
		// Wizard completed - quit the program
		if msg.err != nil {
			// Error saving config - show error and quit
			m.depError = msg.err.Error()
		}
		return m, tea.Quit

	case modelsLoadedMsg:
		m.loadingModels = false
		if msg.err != nil {
			// Failed to load models - use empty list
			m.availableModels = []ModelInfo{}
			m.filteredModels = []ModelInfo{}
		} else {
			m.availableModels = msg.models
			m.modelSelection = 0 // Reset selection
			m.modelScrollOffset = 0
			m.filterModels() // Apply initial filter
		}
		return m, tea.Batch(cmds...)
	}

	return m, tea.Batch(cmds...)
}

// updateActiveInput forwards messages to the currently active text input
func (m Model) updateActiveInput(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	if m.step == StepProvider {
		if m.activeSection == 1 { // Credentials section
			if m.providerSelection == 3 { // Local LLM - endpoint input
				m.apiEndpointInput, cmd = m.apiEndpointInput.Update(msg)
			} else { // API Key input
				m.apiKeyInput, cmd = m.apiKeyInput.Update(msg)
			}
		} else if m.activeSection == 2 && m.modelTab == 1 { // Model search in ALL tab
			m.modelSearchInput, cmd = m.modelSearchInput.Update(msg)
			// Real-time filtering as user types
			m.filterModels()
		}
	} else if m.step == StepDefaults && m.targetLangOther {
		m.customLangCode, cmd = m.customLangCode.Update(msg)
	}

	return m, cmd
}

// handleKeyPress handles keyboard input in navigation mode
func (m Model) handleKeyPress(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	key := msg.String()

	switch key {
	case "q", "ctrl+c":
		// Allow quitting from Step 1 (Language/Deps) or Step 2 (Provider)
		if m.step == StepLanguageDeps || m.step == StepProvider {
			m.quitting = true
			return m, tea.Quit
		}
		// Step 3 requires going back first
		return m, nil

	case "esc":
		// Step 1: Quit the wizard
		if m.step == StepLanguageDeps {
			m.quitting = true
			return m, tea.Quit
		}
		// Step 2+: Go back one step
		if m.step > StepLanguageDeps {
			m.step--
			// Reset navigation state when returning to Step 2 (Provider)
			if m.step == StepProvider {
				m.activeSection = 0
				m.focusManager.ExitInput()
				m.apiKeyInput.Blur()
				m.apiEndpointInput.Blur()
			}
		}
		return m, nil

	case "tab":
		// Tab cycles through sections
		if m.step == StepProvider {
			// Only cycle if in navigation mode (not editing)
			if m.focusManager.Mode() == focus.ModeNav {
				m.activeSection = (m.activeSection + 1) % 3
			}
		} else if m.step == StepDefaults {
			// Cycle between Language (0) and Preferences (1) sections
			if m.focusManager.Mode() == focus.ModeNav {
				m.activeStep3Section = (m.activeStep3Section + 1) % 2
			}
		}
		return m, nil

	case "enter":
		// In Step 2, ENTER toggles edit mode in credentials section
		if m.step == StepProvider {
			if m.activeSection == 1 && m.focusManager.Mode() == focus.ModeNav {
				// Credentials section - enter input mode
				m.focusManager.EnterInput(0)
				if m.providerSelection == 3 {
					m.apiEndpointInput.Focus()
					m.apiKeyInput.Blur()
				} else {
					m.apiKeyInput.Focus()
					m.apiEndpointInput.Blur()
				}
				return m, nil
			} else if m.activeSection == 2 && m.modelTab == 1 && m.focusManager.Mode() == focus.ModeNav {
				// Model section - enter search mode (ALL tab only)
				m.focusManager.EnterInput(0)
				m.modelSearchInput.Focus()
				return m, nil
			}
			// If in input mode, ENTER exits (handled in gatekeeper above)
			// For Provider and Model sections, ENTER goes to next step
		}
		// In Step 3, if "OTHER" is selected and nav mode, enter custom ISO input
		if m.step == StepDefaults && m.targetLangOther && m.focusManager.Mode() == focus.ModeNav {
			m.focusManager.EnterInput(0)
			m.customLangCode.Focus()
			return m, nil
		}
		// Handle as step navigation
		return m.handleEnter()

	case "up", "k":
		if m.step == StepLanguageDeps {
			// Navigate UI language selection
			if m.uiLanguageSelection > 0 {
				m.uiLanguageSelection--
				m.applyUILanguage()
			}
		} else if m.step == StepProvider && m.focusManager.Mode() == focus.ModeNav {
			// Navigate within the active section
			switch m.activeSection {
			case 0: // Provider section
				if m.providerSelection > 0 {
					m.providerSelection--
				}
			case 2: // Model selection
				if m.modelSelection > 0 {
					m.modelSelection--
					// Adjust scroll offset if needed
					if m.modelSelection < m.modelScrollOffset {
						m.modelScrollOffset = m.modelSelection
					}
				}
			}
		} else if m.step == StepDefaults {
			if m.activeStep3Section == 0 && m.languageSelection > 0 {
				m.languageSelection--
				// Reset custom lang if moving away from "OTHER"
				if m.languageSelection != 6 {
					m.targetLangOther = false
				}
			}
		}
		return m, nil

	case "down", "j":
		if m.step == StepLanguageDeps {
			// Navigate UI language selection
			if m.uiLanguageSelection < 2 {
				m.uiLanguageSelection++
				m.applyUILanguage()
			}
		} else if m.step == StepProvider && m.focusManager.Mode() == focus.ModeNav {
			// Navigate within the active section
			switch m.activeSection {
			case 0: // Provider section
				if m.providerSelection < 3 {
					m.providerSelection++
				}
			case 2: // Model selection
				models := m.getDisplayModels()
				if m.modelSelection < len(models)-1 {
					m.modelSelection++
					// Adjust scroll offset if needed (show 7 rows at a time)
					if m.modelSelection-m.modelScrollOffset >= 7 {
						m.modelScrollOffset = m.modelSelection - 6
					}
				}
			}
		} else if m.step == StepDefaults {
			if m.activeStep3Section == 0 && m.languageSelection < 6 {
				m.languageSelection++
				// Check if "OTHER" selected
				if m.languageSelection == 6 {
					m.targetLangOther = true
				}
			}
		}
		return m, nil

	case " ":
		if m.step == StepProvider {
			// In input mode, toggle API key visibility
			if m.focusManager.ActiveField() != -1 {
				if m.apiKeyInput.EchoMode == textinput.EchoPassword {
					m.apiKeyInput.EchoMode = textinput.EchoNormal
					m.showAPIKey = true
				} else {
					m.apiKeyInput.EchoMode = textinput.EchoPassword
					m.showAPIKey = false
				}
			}
		} else if m.step == StepDefaults {
			// Toggle HI tags removal
			m.hiTagsRemoval = !m.hiTagsRemoval
		}
		return m, nil

	case "left", "h":
		if m.step == StepProvider && m.activeSection == 2 && m.focusManager.Mode() == focus.ModeNav {
			// Switch model tabs
			if m.modelTab > 0 {
				m.modelTab--
				m.modelSelection = 0
				m.modelScrollOffset = 0
				m.filterModels()
			}
		} else if m.step == StepDefaults && m.activeStep3Section == 1 && m.tempValue > 0.1 {
			m.tempValue -= 0.1
			if m.tempValue < 0.0 {
				m.tempValue = 0.0
			}
		}
		return m, nil

	case "right", "l":
		if m.step == StepProvider && m.activeSection == 2 && m.focusManager.Mode() == focus.ModeNav {
			// Switch model tabs
			if m.modelTab < 1 {
				m.modelTab++
				m.modelSelection = 0
				m.modelScrollOffset = 0
				m.filterModels()
			}
		} else if m.step == StepDefaults && m.activeStep3Section == 1 && m.tempValue < 1.0 {
			m.tempValue += 0.1
			if m.tempValue > 1.0 {
				m.tempValue = 1.0
			}
		}
		return m, nil

	case "r":
		// Recheck dependencies
		if m.step == StepLanguageDeps && m.checkComplete {
			m.checkComplete = false
			spinnerCmd := m.depSpinner.Start()
			return m, tea.Batch(checkDependencies, spinnerCmd)
		}
		return m, nil
	}

	return m, nil
}

// handleEnter handles the enter key based on current step
func (m Model) handleEnter() (tea.Model, tea.Cmd) {
	switch m.step {
	case StepLanguageDeps:
		// Step 1: UI Language + Dependencies check complete
		// Allow user to proceed even if dependencies are missing
		if !m.checkComplete {
			return m, nil // Wait for check to complete
		}
		m.step = StepProvider
		return m, nil

	case StepProvider:
		// Step 2: Validate provider and API key
		if m.providerSelection == 3 {
			// Local LLM - check endpoint
			if m.apiEndpointInput.Value() == "" {
				return m, nil
			}
			m.config.LocalEndpoint = m.apiEndpointInput.Value()
			m.config.AIProvider = "local"
		} else {
			// Cloud provider - check API key
			if m.apiKeyInput.Value() == "" {
				return m, nil
			}
			m.config.APIKey = m.apiKeyInput.Value()
			switch m.providerSelection {
			case 0:
				m.config.AIProvider = "openrouter"
			case 1:
				m.config.AIProvider = "gemini"
			case 2:
				m.config.AIProvider = "openai"
			}
		}
		// Save selected model if available
		models := m.getDisplayModels()
		if m.modelSelection < len(models) {
			m.config.Model = models[m.modelSelection].ID
		}
		m.step = StepDefaults
		return m, nil

	case StepDefaults:
		// Step 3: Save config and finish
		m.finished = true
		return m, m.saveAndFinish()
	}

	return m, nil
}

// View renders the wizard
func (m Model) View() string {
	if m.finished {
		return ""
	}

	// Check if terminal is too small
	if layout.IsTooSmall(m.width, m.height) {
		return layout.RenderTooSmallWarning(m.width, m.height)
	}

	var content string
	switch m.step {
	case StepLanguageDeps:
		content = m.renderLanguageDepsStep()
	case StepProvider:
		content = m.renderProviderStep()
	case StepDefaults:
		content = m.renderDefaultsStep()
	}

	safeWidth := layout.SafeWidth(m.width, layout.MinWidth)
	safeHeight := layout.SafeHeight(m.height, layout.MinHeight)
	return lipgloss.Place(safeWidth, safeHeight,
		lipgloss.Center, lipgloss.Top, // Changed from Center to Top for vertical alignment
		content,
	)
}

// renderLanguageDepsStep renders step 1: UI Language + Dependencies
func (m Model) renderLanguageDepsStep() string {
	const contentWidth = 88

	header := styles.HeaderBorder.Width(contentWidth).Render(
		fmt.Sprintf(" %s %s [STEP 1/3] ", locales.T("wizard.title"), strings.Repeat("▒", 38)),
	)

	// Interface Language panel
	uiLangKeys := []string{"wizard.step1.lang_english", "wizard.step1.lang_portuguese", "wizard.step1.lang_spanish"}

	var uiLangList strings.Builder
	for i, key := range uiLangKeys {
		lang := locales.T(key) // Get translation at render time
		if i == m.uiLanguageSelection {
			uiLangList.WriteString("   (o) " + styles.Highlight.Render(lang) + "\n")
		} else {
			uiLangList.WriteString("   ( ) " + lang + "\n")
		}
	}

	langPanel := styles.Panel.Width(contentWidth).Render(
		fmt.Sprintf("1. %s\n", locales.T("wizard.step1.interface_language")) +
			uiLangList.String(),
	)

	// Dependencies check panel
	var depsContent strings.Builder
	depsContent.WriteString(styles.SectionStyle.Render(locales.T("wizard.step1.deps_title")) + "\n")
	depsContent.WriteString(locales.T("wizard.step1.deps_subtitle") + "\n\n")

	if !m.checkComplete {
		// Show spinner while checking
		if m.depSpinner.IsActive() {
			depsContent.WriteString("   " + m.depSpinner.View() + " " + locales.T("wizard.step1.scanning") + "\n")
		} else {
			depsContent.WriteString("   " + locales.T("wizard.step1.scanning") + "\n")
		}
	} else {
		// Show status for each dependency
		for _, dep := range dependencies.Dependencies {
			depsContent.WriteString(fmt.Sprintf("   ┌── %s ──────────────────────────────────────\n", dep.Name))

			allOK := true
			var foundBinary string
			for _, binary := range dep.TargetBinaries {
				status, exists := m.depStatus[binary]
				if !exists || !status {
					allOK = false
					break
				}
				foundBinary = binary
			}

			if allOK {
				depsContent.WriteString("   │   " + locales.T("wizard.step1.status") + ": " + styles.StatusOK.Render("✓ "+locales.T("wizard.step1.found")) + "\n")
				if dependencies.CheckSystemPath(foundBinary) {
					depsContent.WriteString("   │   " + locales.T("wizard.step1.location") + ": " + styles.Dimmed.Render(locales.T("wizard.step1.system_path")) + "\n")
				} else {
					depsContent.WriteString(fmt.Sprintf("   │   %s: %s\n", locales.T("wizard.step1.location"), styles.Dimmed.Render(dependencies.GetBinaryPath(foundBinary))))
				}
			} else {
				depsContent.WriteString("   │   " + locales.T("wizard.step1.status") + ": " + styles.StatusError.Render("✗ "+locales.T("wizard.step1.not_found")) + "\n")
				depsContent.WriteString("   │   \n")
				depsContent.WriteString("   │   " + styles.StatusWarning.Render("⚠ "+locales.T("wizard.step1.action_required")) + "\n")
				depsContent.WriteString("   │   " + locales.Tf("wizard.step1.install_manually", dep.Name) + "\n")
				depsContent.WriteString("   │   \n")
				depsContent.WriteString("   │   ➤ " + locales.T("wizard.step1.download") + ": " + styles.CodeBlock.Render(dep.GetDownloadURL()) + "\n")
			}

			depsContent.WriteString("   └──────────────────────────────────────────────\n\n")
		}
	}

	depsPanel := styles.Panel.Width(contentWidth).Render(depsContent.String())

	var footerText string
	if !m.checkComplete {
		footerText = styles.RenderHotkey("Q", locales.T("common.quit")) + strings.Repeat(" ", 35) +
			styles.Dimmed.Render("["+locales.T("wizard.step1.checking")+"...]")
	} else {
		allOK := true
		for _, ok := range m.depStatus {
			if !ok {
				allOK = false
				break
			}
		}

		if allOK {
			footerText = styles.RenderHotkey("Q", locales.T("common.quit")) + strings.Repeat(" ", 30) +
				styles.RenderHotkey("ENTER", locales.T("common.next")+" >")
		} else {
			footerText = styles.RenderHotkey("Q", locales.T("common.quit")) + " | " +
				styles.RenderHotkey("R", locales.T("wizard.step1.recheck")) + strings.Repeat(" ", 20) +
				styles.RenderHotkey("ENTER", locales.T("wizard.step1.skip")+" >")
		}
	}

	footer := styles.FooterBorder.Width(contentWidth).Render(footerText)

	return lipgloss.JoinVertical(lipgloss.Center,
		header,
		"",
		langPanel,
		"",
		depsPanel,
		footer,
	)
}

// renderProviderStep renders step 2: Provider + API Key + Model Selection
func (m Model) renderProviderStep() string {
	const contentWidth = 88

	header := styles.HeaderBorder.Width(contentWidth).Render(
		fmt.Sprintf(" %s %s [STEP 2/3] ", locales.T("wizard.title"), strings.Repeat("▒", 38)),
	)

	providerKeys := []string{
		"wizard.step1.provider_openrouter",
		"wizard.step1.provider_gemini",
		"wizard.step1.provider_openai",
		"wizard.step1.provider_local",
	}

	var providerList strings.Builder
	for i, key := range providerKeys {
		p := locales.T(key) // Get translation at render time
		if i == m.providerSelection {
			providerList.WriteString("   (o) " + styles.Highlight.Render(p) + "\n")
		} else {
			providerList.WriteString("   ( ) " + p + "\n")
		}
	}

	providerPanelStyle := styles.Panel.Width(contentWidth)
	if m.activeSection == 0 {
		providerPanelStyle = providerPanelStyle.BorderForeground(styles.Yellow)
	}
	providerPanel := providerPanelStyle.Render(
		fmt.Sprintf("1. %s (%s)\n", locales.T("wizard.step1.title"), locales.T("wizard.step1.subtitle")) +
			providerList.String(),
	)

	// Credentials panel
	var credContent string
	if m.providerSelection == 3 {
		inputView := m.apiEndpointInput.View()
		editingStatus := ""
		// Only show editing status if this section is active
		if m.activeSection == 1 {
			if m.focusManager.Mode() == focus.ModeInput {
				editingStatus = "   " + styles.StatusOK.Render("["+locales.T("wizard.step2.editing")+"]") + " " + styles.KeyHintStyle.Render("[ENTER] "+locales.T("wizard.step2.stop_editing"))
			} else {
				editingStatus = "   " + styles.KeyHintStyle.Render("[ENTER] "+locales.T("wizard.step2.start_editing"))
			}
		}
		credContent = fmt.Sprintf("2. %s\n", locales.T("wizard.step2.endpoint_title")) +
			fmt.Sprintf("   %s > ", locales.T("wizard.step2.endpoint_label")) + inputView + "\n" +
			editingStatus
	} else {
		inputView := m.apiKeyInput.View()
		showLabel := locales.T("wizard.step2.show")
		if m.showAPIKey {
			showLabel = locales.T("wizard.step2.hide")
		}
		editingStatus := ""
		// Only show editing status if this section is active
		if m.activeSection == 1 {
			if m.focusManager.Mode() == focus.ModeInput {
				editingStatus = "   " + styles.StatusOK.Render("["+locales.T("wizard.step2.editing")+"]") + " " + styles.KeyHintStyle.Render("[ENTER] "+locales.T("wizard.step2.stop_editing")+" | [SPACE] "+showLabel)
			} else {
				editingStatus = "   " + styles.KeyHintStyle.Render("[ENTER] "+locales.T("wizard.step2.start_editing")+" | [SPACE] "+showLabel)
			}
		}
		credContent = fmt.Sprintf("2. %s\n", locales.T("wizard.step1.credentials_title")) +
			fmt.Sprintf("   %s > ", locales.T("wizard.step1.api_key_label")) + inputView + "\n" +
			editingStatus
	}

	credPanelStyle := styles.Panel.Width(contentWidth)
	if m.activeSection == 1 {
		credPanelStyle = credPanelStyle.BorderForeground(styles.Yellow)
	}
	credPanel := credPanelStyle.Render(credContent)

	// Model selection panel with tabs and table
	var modelContent strings.Builder
	modelContent.WriteString(fmt.Sprintf("3. %s\n\n", locales.T("wizard.step2.model_selection_title")))

	if m.loadingModels {
		modelContent.WriteString("   " + styles.StatusWarning.Render("⏳ "+locales.T("wizard.step2.loading_models")+"...") + "\n")
	} else if len(m.availableModels) == 0 {
		modelContent.WriteString("   " + styles.Dimmed.Render("("+locales.T("wizard.step2.model_selection_hint")+")") + "\n")
	} else {
		// Render tabs
		freeTab := " FREE "
		allTab := " ALL MODELS "
		if m.modelTab == 0 {
			freeTab = styles.StatusOK.Render(" [ FREE ] ")
			allTab = styles.Dimmed.Render(" ALL MODELS ")
		} else {
			freeTab = styles.Dimmed.Render(" FREE ")
			allTab = styles.StatusOK.Render(" [ ALL MODELS ] ")
		}
		modelContent.WriteString("   " + freeTab + allTab + "  " + styles.KeyHintStyle.Render("[← →] Switch") + "\n\n")

		// Search input (only in ALL MODELS tab)
		if m.modelTab == 1 {
			searchView := m.modelSearchInput.View()
			searchHint := ""
			// Only show editing status if model section is active
			if m.activeSection == 2 {
				if m.focusManager.Mode() == focus.ModeInput {
					searchHint = " " + styles.StatusOK.Render("[EDITING]") + " [ENTER] Stop"
				} else {
					searchHint = " [ENTER] Search"
				}
			}
			modelContent.WriteString("   Search > " + searchView + searchHint + "\n\n")
		}

		// Table header
		modelContent.WriteString("   ┌────────────────────────────┬───────────┬──────────────┐\n")
		modelContent.WriteString("   │ MODEL                      │ CONTEXT   │ PRICE/1M     │\n")
		modelContent.WriteString("   ├────────────────────────────┼───────────┼──────────────┤\n")

		// Table rows (show 7 rows)
		models := m.getDisplayModels()
		if len(models) == 0 {
			modelContent.WriteString("   │ " + styles.Dimmed.Render("No models found") + strings.Repeat(" ", 55) + "│\n")
		} else {
			for i := 0; i < 7 && i+m.modelScrollOffset < len(models); i++ {
				idx := i + m.modelScrollOffset
				model := models[idx]

				// Truncate name if too long
				name := model.Name
				if len(name) > 26 {
					name = name[:23] + "..."
				}

				row := fmt.Sprintf(" %-26s │ %-9s │ %-12s ", name, model.ContextSize, model.PricePerM)

				if idx == m.modelSelection {
					modelContent.WriteString("   │" + styles.Highlight.Render(row) + "│\n")
				} else {
					modelContent.WriteString("   │" + row + "│\n")
				}
			}
		}

		modelContent.WriteString("   └────────────────────────────┴───────────┴──────────────┘\n")

		// Scroll indicator
		if len(models) > 7 {
			modelContent.WriteString(fmt.Sprintf("   %s (%d/%d)\n",
				styles.Dimmed.Render("↑↓ Scroll"),
				m.modelSelection+1,
				len(models)))
		}
	}

	modelPanelStyle := styles.Panel.Width(contentWidth)
	if m.activeSection == 2 {
		modelPanelStyle = modelPanelStyle.BorderForeground(styles.Yellow)
	}
	modelPanel := modelPanelStyle.Render(modelContent.String())

	footer := styles.FooterBorder.Width(contentWidth).Render(
		styles.RenderHotkey("ESC", "< "+locales.T("common.back")) + strings.Repeat(" ", 30) +
			styles.RenderHotkey("ENTER", locales.T("common.next_step")+" >"),
	)

	return lipgloss.JoinVertical(lipgloss.Center,
		header,
		"",
		providerPanel,
		"",
		credPanel,
		"",
		modelPanel,
		footer,
	)
}

// renderDefaultsStep renders step 3: Target Language + Preferences
func (m Model) renderDefaultsStep() string {
	const contentWidth = 88

	header := styles.HeaderBorder.Width(contentWidth).Render(
		fmt.Sprintf(" %s %s [STEP 3/3] ", locales.T("wizard.title"), strings.Repeat("▒", 38)),
	)

	languages := []string{
		"PT-BR (Português)",
		"EN-US (English)",
		"ES (Español)",
		"JA-JP (Japanese)",
		"FR-FR (Français)",
		"DE (Deutsch)",
		locales.T("wizard.step3.other") + " (" + locales.T("wizard.step3.custom_iso") + ")",
	}

	var langList strings.Builder
	for i, lang := range languages {
		if i == m.languageSelection {
			langList.WriteString("   (o) " + styles.Highlight.Render(lang) + "\n")
		} else {
			langList.WriteString("   ( ) " + lang + "\n")
		}
	}

	// Show custom input if "OTHER" selected
	if m.targetLangOther {
		inputView := m.customLangCode.View()
		editingHint := ""
		if m.focusManager.Mode() == focus.ModeInput {
			editingHint = " " + styles.StatusOK.Render("["+locales.T("wizard.step2.editing")+"]") + " [ENTER] " + locales.T("wizard.step2.stop_editing")
		} else {
			editingHint = " [ENTER] " + locales.T("wizard.step2.start_editing")
		}
		// Show language name if ISO code is recognized
		langName := getLanguageName(m.customLangCode.Value())
		langHint := ""
		if langName != "" {
			langHint = " → " + styles.StatusOK.Render(langName)
		} else if m.customLangCode.Value() != "" {
			langHint = " → " + styles.StatusWarning.Render("?")
		}
		langList.WriteString("\n   " + locales.T("wizard.step3.iso_code") + " > " + inputView + langHint + "\n")
		langList.WriteString("   " + editingHint + "\n")
	}

	// Add section indicator
	langPanelStyle := styles.Panel.Width(contentWidth)
	langIndicator := "  "
	if m.activeStep3Section == 0 {
		langPanelStyle = langPanelStyle.BorderForeground(styles.Yellow)
		langIndicator = styles.StatusOK.Render("▸ ")
	}
	langPanel := langPanelStyle.Render(
		langIndicator + styles.SectionStyle.Render(locales.T("wizard.step3.target_title")) + "\n\n" +
			langList.String(),
	)

	hiTag := " "
	if m.hiTagsRemoval {
		hiTag = "X"
	}

	// Add section indicator
	prefPanelStyle := styles.Panel.Width(contentWidth)
	prefIndicator := "  "
	if m.activeStep3Section == 1 {
		prefPanelStyle = prefPanelStyle.BorderForeground(styles.Yellow)
		prefIndicator = styles.StatusOK.Render("▸ ")
	}
	prefPanel := prefPanelStyle.Render(
		fmt.Sprintf("%s%s\n\n"+
			"   [%s] %s %s\n"+
			"   %s: [ %.1f ] %s\n\n"+
			"   ┌── [?] %s ──────────────────\n"+
			"   │  %s\n"+
			"   │  %s\n"+
			"   │  %s\n"+
			"   └───────────────────────────────────────────────────────\n",
			prefIndicator, styles.SectionStyle.Render(locales.T("wizard.step3.preferences_title")),
			hiTag, locales.T("wizard.step3.remove_hi_tags"), styles.KeyHintStyle.Render("[SPACE] "+locales.T("wizard.step3.toggle")),
			locales.T("wizard.step3.temperature"), m.tempValue, styles.KeyHintStyle.Render("[← →] "+locales.T("wizard.step3.adjust")),
			locales.T("wizard.step3.temp_helper"),
			locales.T("wizard.step3.temp_help_low"),
			locales.T("wizard.step3.temp_help_mid"),
			locales.T("wizard.step3.temp_help_high")),
	)

	footer := styles.FooterBorder.Width(contentWidth).Render(
		styles.RenderHotkey("ESC", "< "+locales.T("common.back")) + " | " +
			styles.RenderHotkey("TAB", locales.T("wizard.step3.next_section")) + strings.Repeat(" ", 15) +
			styles.RenderHotkey("ENTER", locales.T("common.finish")+" >"),
	)

	return lipgloss.JoinVertical(lipgloss.Center,
		header,
		"",
		langPanel,
		"",
		prefPanel,
		"",
		footer,
	)
}

// saveAndFinish saves the configuration and finishes the wizard
func (m Model) saveAndFinish() tea.Cmd {
	return func() tea.Msg {
		// Apply UI language from step 1
		uiLangs := []string{"en", "pt-br", "es"}
		if m.uiLanguageSelection < len(uiLangs) {
			m.config.InterfaceLang = uiLangs[m.uiLanguageSelection]
		}

		// Apply defaults from step 3
		langs := []string{"pt-br", "en-us", "es", "ja-jp", "fr-fr", "de"}
		if m.languageSelection == 6 && m.customLangCode.Value() != "" {
			// Custom ISO code
			m.config.TargetLang = strings.ToLower(strings.TrimSpace(m.customLangCode.Value()))
		} else if m.languageSelection < len(langs) {
			m.config.TargetLang = langs[m.languageSelection]
		}

		m.config.Temperature = m.tempValue
		m.config.RemoveHITags = m.hiTagsRemoval

		// No default model - user must select from FREE/ALL MODELS tab
		m.config.Model = "" // Empty until user selects in configuration

		// Set bin path
		m.config.BinPath = dependencies.BinDir

		// Save config
		if err := m.config.Save(); err != nil {
			return finishMsg{err: err}
		}

		return finishMsg{err: nil}
	}
}

// Messages

type checkDepsMsg struct {
	status map[string]bool
}

type downloadProgressMsg struct {
	depName  string
	progress float64
}

type downloadCompleteMsg struct {
	err error
}

type finishMsg struct {
	err error
}

// modelsLoadedMsg is sent when models are loaded from the provider
type modelsLoadedMsg struct {
	models []ModelInfo
	err    error
}

// Commands

func checkDependencies() tea.Msg {
	status, _ := dependencies.Check()
	return checkDepsMsg{status: status}
}

// downloadDependenciesCmd creates a download command that reports progress through channel
func downloadDependenciesCmd() tea.Msg {
	// Get missing dependencies
	missing, err := dependencies.GetMissingDependencies()
	if err != nil {
		return downloadCompleteMsg{err: err}
	}

	// Download each missing dependency
	for _, dep := range missing {
		err := dependencies.DownloadAndInstall(dep, func(read, total int64) {
			// Progress callback is called during download
			// Since we can't send tea.Msg from here, the progress is implicit
			// The UI shows a spinner while download is in progress
		})

		if err != nil {
			return downloadCompleteMsg{err: fmt.Errorf("failed to install %s: %w", dep.Name, err)}
		}
	}

	return downloadCompleteMsg{err: nil}
}

func downloadDependencies() tea.Msg {
	return downloadDependenciesCmd()
}

// applyUILanguage applies the currently selected UI language immediately
func (m *Model) applyUILanguage() {
	uiLangs := []string{"en", "pt-br", "es"}
	if m.uiLanguageSelection < len(uiLangs) {
		_ = locales.Load(uiLangs[m.uiLanguageSelection])
	}
}

// fetchModelsCmd creates a command to fetch models from the provider
func (m Model) fetchModelsCmd() tea.Cmd {
	return func() tea.Msg {
		ctx := context.Background()

		var provider ai.LLMProvider
		var err error

		switch m.providerSelection {
		case 0: // OpenRouter
			provider = ai.NewOpenRouterAdapter(m.apiKeyInput.Value(), "", 0.3)
			// Fetch models from provider
			modelStrings, err := provider.ListModels(ctx)
			if err != nil {
				// If fetch fails, return sample data for testing
				freeModels := []ModelInfo{
					{ID: "meta-llama/llama-3.3-70b-instruct:free", Name: "Llama 3.3 70B", ContextSize: "128k", PricePerM: "FREE", IsFree: true},
					{ID: "qwen/qwen-2.5-72b-instruct:free", Name: "Qwen 2.5 72B", ContextSize: "32k", PricePerM: "FREE", IsFree: true},
					{ID: "mistralai/mistral-7b-instruct:free", Name: "Mistral 7B", ContextSize: "32k", PricePerM: "FREE", IsFree: true},
					{ID: "google/gemma-2-9b-it:free", Name: "Gemma 2 9B", ContextSize: "8k", PricePerM: "FREE", IsFree: true},
					{ID: "microsoft/phi-3-medium-128k:free", Name: "Phi-3 Medium", ContextSize: "128k", PricePerM: "FREE", IsFree: true},
				}
				paidModels := []ModelInfo{
					{ID: "anthropic/claude-3-opus", Name: "Claude 3 Opus", ContextSize: "200k", PricePerM: "$15.00", IsFree: false},
					{ID: "anthropic/claude-3-sonnet", Name: "Claude 3 Sonnet", ContextSize: "200k", PricePerM: "$3.00", IsFree: false},
					{ID: "anthropic/claude-3-haiku", Name: "Claude 3 Haiku", ContextSize: "200k", PricePerM: "$0.25", IsFree: false},
					{ID: "openai/gpt-4o", Name: "GPT-4o", ContextSize: "128k", PricePerM: "$2.50", IsFree: false},
					{ID: "openai/gpt-4o-mini", Name: "GPT-4o Mini", ContextSize: "128k", PricePerM: "$0.15", IsFree: false},
					{ID: "google/gemini-1.5-pro", Name: "Gemini 1.5 Pro", ContextSize: "1M", PricePerM: "$1.25", IsFree: false},
					{ID: "google/gemini-2.0-flash-exp:free", Name: "Gemini 2.0 Flash", ContextSize: "1M", PricePerM: "FREE", IsFree: true},
					{ID: "meta-llama/llama-3.1-405b", Name: "Llama 3.1 405B", ContextSize: "128k", PricePerM: "$3.00", IsFree: false},
				}
				allModels := append(freeModels, paidModels...)
				return modelsLoadedMsg{models: allModels, err: nil}
			}
			// Parse real models - assume free if contains ":free"
			models := []ModelInfo{}
			for _, model := range modelStrings {
				isFree := strings.Contains(strings.ToLower(model), ":free")
				models = append(models, parseModelsToModelInfo([]string{model}, isFree)...)
			}
			return modelsLoadedMsg{models: models, err: nil}

		case 1: // Gemini
			provider, err = ai.NewGeminiAdapter(ctx, m.apiKeyInput.Value(), "", 0.3)
			if err != nil {
				return modelsLoadedMsg{models: nil, err: err}
			}
			modelStrings, err := provider.ListModels(ctx)
			if err != nil {
				return modelsLoadedMsg{models: nil, err: err}
			}
			models := parseModelsToModelInfo(modelStrings, false)
			return modelsLoadedMsg{models: models, err: nil}

		case 2: // OpenAI
			// OpenAI - lista padrão
			modelStrings := []string{"gpt-4o", "gpt-4o-mini", "gpt-4-turbo", "gpt-3.5-turbo"}
			return modelsLoadedMsg{
				models: parseModelsToModelInfo(modelStrings, false),
				err:    nil,
			}
		case 3: // Local
			provider = ai.NewLocalLLMAdapter(m.apiEndpointInput.Value(), "", 0.3)
			modelStrings, err := provider.ListModels(ctx)
			if err != nil {
				return modelsLoadedMsg{models: nil, err: err}
			}
			// Local models are free
			models := parseModelsToModelInfo(modelStrings, true)
			return modelsLoadedMsg{models: models, err: nil}
		}

		return modelsLoadedMsg{models: nil, err: fmt.Errorf("provider not configured")}
	}
}

// getLanguageName returns the language name for an ISO code
func getLanguageName(code string) string {
	code = strings.ToLower(strings.TrimSpace(code))
	if name, ok := isoLanguageNames[code]; ok {
		return name
	}
	return ""
}

// filterModels filters available models based on current tab and search query
func (m *Model) filterModels() {
	m.filteredModels = []ModelInfo{}
	searchQuery := strings.ToLower(strings.TrimSpace(m.modelSearchInput.Value()))

	for _, model := range m.availableModels {
		// Filter by tab (FREE vs ALL)
		if m.modelTab == 0 && !model.IsFree {
			continue
		}

		// Filter by search (only in ALL tab)
		if m.modelTab == 1 && searchQuery != "" {
			nameMatch := strings.Contains(strings.ToLower(model.Name), searchQuery)
			idMatch := strings.Contains(strings.ToLower(model.ID), searchQuery)
			if !nameMatch && !idMatch {
				continue
			}
		}

		m.filteredModels = append(m.filteredModels, model)
	}

	// Reset selection if out of bounds
	if m.modelSelection >= len(m.filteredModels) {
		m.modelSelection = 0
		m.modelScrollOffset = 0
	}
}

// getDisplayModels returns the models to display based on current filters
func (m Model) getDisplayModels() []ModelInfo {
	if len(m.filteredModels) > 0 {
		return m.filteredModels
	}

	// If no filters applied, show based on tab
	if m.modelTab == 0 {
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

// parseModelsToModelInfo converts string list to ModelInfo list with metadata
func parseModelsToModelInfo(models []string, isFree bool) []ModelInfo {
	result := []ModelInfo{}

	// Context sizes database
	contextSizes := map[string]string{
		"gpt-4o":           "128k",
		"gpt-4o-mini":      "128k",
		"gpt-4-turbo":      "128k",
		"gpt-3.5-turbo":    "16k",
		"gemini-1.5-pro":   "1M",
		"gemini-1.5-flash": "1M",
		"gemini-2.0-flash": "1M",
		"claude-3-opus":    "200k",
		"claude-3-sonnet":  "200k",
		"claude-3-haiku":   "200k",
		"llama-3":          "8k",
		"llama-3.1":        "128k",
		"llama-3.3":        "128k",
		"mistral":          "32k",
		"mixtral":          "32k",
		"qwen":             "32k",
		"phi-3":            "128k",
		"gemma":            "8k",
	}

	for _, model := range models {
		// Determine context size
		ctx := "varies"
		for key, size := range contextSizes {
			if strings.Contains(strings.ToLower(model), key) {
				ctx = size
				break
			}
		}

		result = append(result, ModelInfo{
			ID:          model,
			Name:        model,
			ContextSize: ctx,
			PricePerM:   getPriceForModel(model, isFree),
			IsFree:      isFree,
		})
	}
	return result
}

// getPriceForModel returns a price string for a given model
func getPriceForModel(modelID string, isFree bool) string {
	if isFree {
		return "FREE"
	}

	// Common model pricing (approximations)
	prices := map[string]string{
		"gpt-4o":           "$2.50",
		"gpt-4o-mini":      "$0.15",
		"gpt-4-turbo":      "$10.00",
		"gpt-3.5-turbo":    "$0.50",
		"gemini-1.5-pro":   "$1.25",
		"gemini-1.5-flash": "$0.07",
		"claude-3-opus":    "$15.00",
		"claude-3-sonnet":  "$3.00",
		"claude-3-haiku":   "$0.25",
	}

	// Check for partial matches
	for key, price := range prices {
		if strings.Contains(strings.ToLower(modelID), key) {
			return price
		}
	}

	return "$?.??"
}

// Finished returns true if the wizard is complete
func (m Model) Finished() bool {
	return m.finished
}

// Quitting returns true if the user quit the wizard
func (m Model) Quitting() bool {
	return m.quitting
}
