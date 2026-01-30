package wizard

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/lsilvatti/bakasub/internal/config"
	"github.com/lsilvatti/bakasub/internal/core/ai"
	"github.com/lsilvatti/bakasub/internal/core/dependencies"
	"github.com/lsilvatti/bakasub/internal/locales"
	"github.com/lsilvatti/bakasub/internal/ui/components"
	"github.com/lsilvatti/bakasub/internal/ui/components/modelselect"
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
	validatingKey     bool // true while validating key
	keyValidated      bool
	keyValidationErr  string            // validation error message
	modelSelector     modelselect.Model // New reusable component
	loadingModels     bool

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

	fm := focus.NewManager(1) // 1 text input field at a time

	// Create model selector component
	modelSelector := modelselect.New(fm)

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
		modelSelector:       modelSelector,
		loadingModels:       false,
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
	// Request current terminal size
	return tea.Batch(tea.WindowSize(), textinput.Blink, checkDependencies, spinnerCmd)
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
					m.customLangCode.Blur()
					return m, tea.Batch(cmds...)
				}
			}
			// In Step 2 and Step 3, ENTER exits input mode
			if msg.String() == "enter" && (m.step == StepProvider || m.step == StepDefaults) {
				// Exit input mode
				m.focusManager.ExitInput()
				m.apiKeyInput.Blur()
				m.apiEndpointInput.Blur()
				m.customLangCode.Blur()
				m.modelSelector.SetActive(false)
				// Validate key with real API request ONLY if in credentials section
				if m.step == StepProvider && m.activeSection == 1 {
					keyValue := m.apiKeyInput.Value()
					endpointValue := m.apiEndpointInput.Value()
					if (keyValue != "" && m.providerSelection != 3) || (endpointValue != "" && m.providerSelection == 3) {
						// Only validate if key was changed (not already validated)
						if !m.keyValidated || m.keyValidationErr != "" {
							// Start validation
							m.validatingKey = true
							m.keyValidated = false
							m.keyValidationErr = ""
							cmds = append(cmds, m.validateKeyCmd())
						}
					}
				}
				return m, tea.Batch(cmds...)
			}
			// "/" exits search mode
			if msg.String() == "/" && m.step == StepProvider && m.activeSection == 2 {
				m.focusManager.ExitInput()
				m.modelSelector.SetActive(false)
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
		// First, delegate to model selector if it's active (but NOT for global keys like enter, tab, esc)
		if m.step == StepProvider && m.activeSection == 2 && m.modelSelector.IsActive() {
			// These keys should be handled by the wizard, not the model selector
			switch msg.String() {
			case "enter", "tab", "esc", "q", "ctrl+c":
				// Fall through to handleKeyPress
			default:
				// Delegate other keys to model selector
				var selectorCmd tea.Cmd
				m.modelSelector, selectorCmd = m.modelSelector.Update(msg)
				if selectorCmd != nil {
					cmds = append(cmds, selectorCmd)
				}
				return m, tea.Batch(cmds...)
			}
		}

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
			models := []modelselect.ModelInfo{}
			m.modelSelector.SetModels(models)
		} else {
			// Convert ModelInfo to modelselect.ModelInfo
			models := make([]modelselect.ModelInfo, len(msg.models))
			for i, model := range msg.models {
				models[i] = modelselect.ModelInfo{
					ID:          model.ID,
					Name:        model.Name,
					ContextSize: model.ContextSize,
					PricePerM:   model.PricePerM,
					IsFree:      model.IsFree,
				}
			}
			m.modelSelector.SetModels(models)
		}
		return m, tea.Batch(cmds...)

	case keyValidatedMsg:
		m.validatingKey = false
		if msg.isValid {
			m.keyValidated = true
			m.keyValidationErr = ""
			// Start loading models after successful validation
			m.loadingModels = true
			cmds = append(cmds, m.fetchModelsCmd())
		} else {
			m.keyValidated = false
			if msg.err != nil {
				m.keyValidationErr = msg.err.Error()
			} else {
				m.keyValidationErr = "Invalid key"
			}
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
				oldValue := m.apiEndpointInput.Value()
				m.apiEndpointInput, cmd = m.apiEndpointInput.Update(msg)
				// Reset validation if endpoint changed
				if m.apiEndpointInput.Value() != oldValue {
					m.keyValidated = false
					m.keyValidationErr = ""
				}
			} else { // API Key input
				oldValue := m.apiKeyInput.Value()
				m.apiKeyInput, cmd = m.apiKeyInput.Update(msg)
				// Reset validation if key changed
				if m.apiKeyInput.Value() != oldValue {
					m.keyValidated = false
					m.keyValidationErr = ""
				}
			}
		} else if m.activeSection == 2 && m.modelSelector.IsActive() { // Model search
			var selectorCmd tea.Cmd
			m.modelSelector, selectorCmd = m.modelSelector.Update(msg)
			cmd = selectorCmd
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
		// If model selector is active, deactivate it first
		if m.step == StepProvider && m.activeSection == 2 && m.modelSelector.IsActive() {
			m.modelSelector.SetActive(false)
			return m, nil
		}
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
			// Cycle through sections even if in input mode (exit input mode first)
			if m.focusManager.Mode() != focus.ModeNav {
				m.focusManager.ExitInput()
				m.apiKeyInput.Blur()
				m.apiEndpointInput.Blur()
			}
			m.activeSection = (m.activeSection + 1) % 3
			// Reset model selector when navigating away
			if m.activeSection != 2 {
				m.modelSelector.SetActive(false)
			}
		} else if m.step == StepDefaults {
			// Cycle between Language (0) and Preferences (1) sections
			if m.focusManager.Mode() == focus.ModeNav {
				m.activeStep3Section = (m.activeStep3Section + 1) % 2
			}
		}
		return m, nil

	case "e":
		// In Step 2, E enters edit mode in credentials section
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
				return m, textinput.Blink
			}
		}
		return m, nil

	case "enter":
		// ENTER is for step navigation (handled in gatekeeper for exiting input mode)
		return m.handleEnter()

	case "/":
		// "/" key triggers search mode in model section
		if m.step == StepProvider && m.activeSection == 2 && m.focusManager.Mode() == focus.ModeNav {
			// Model section - enter search mode
			m.focusManager.EnterInput(0)
			m.modelSelector.SetActive(true)
			return m, nil
		}
		return m, nil

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
					// Reset validation state when provider changes
					m.keyValidated = false
					m.keyValidationErr = ""
					m.validatingKey = false
					m.loadingModels = false
					m.modelSelector.SetModels(nil)
				}
			case 2: // Model section - delegate to model selector
				m.modelSelector.SetActive(true)
				var selectorCmd tea.Cmd
				m.modelSelector, selectorCmd = m.modelSelector.Update(msg)
				return m, selectorCmd
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
					// Reset validation state when provider changes
					m.keyValidated = false
					m.keyValidationErr = ""
					m.validatingKey = false
					m.loadingModels = false
					m.modelSelector.SetModels(nil)
				}
			case 2: // Model section - delegate to model selector
				m.modelSelector.SetActive(true)
				var selectorCmd tea.Cmd
				m.modelSelector, selectorCmd = m.modelSelector.Update(msg)
				return m, selectorCmd
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
			// Toggle API key visibility when in credentials section (both nav and input mode)
			if m.activeSection == 1 && m.providerSelection != 3 {
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
		if m.step == StepDefaults && m.activeStep3Section == 1 && m.tempValue > 0.1 {
			m.tempValue -= 0.1
			if m.tempValue < 0.0 {
				m.tempValue = 0.0
			}
		}
		return m, nil

	case "right", "l":
		if m.step == StepDefaults && m.activeStep3Section == 1 && m.tempValue < 1.0 {
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
			// Require endpoint validation before proceeding
			if !m.keyValidated && !m.validatingKey {
				// Endpoint not validated yet
				return m, nil
			}
			// Wait for validation to complete
			if m.validatingKey {
				return m, nil
			}
			m.config.LocalEndpoint = m.apiEndpointInput.Value()
			m.config.AIProvider = "local"
		} else {
			// Cloud provider - check API key
			if m.apiKeyInput.Value() == "" {
				return m, nil
			}
			// Require key validation before proceeding
			if !m.keyValidated && !m.validatingKey {
				// Key not validated yet - show error or wait
				return m, nil
			}
			// Wait for validation to complete
			if m.validatingKey {
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
		// Save selected model if available (optional - can proceed without model)
		selectedModel := m.modelSelector.GetSelectedModel()
		if selectedModel != nil {
			m.config.Model = selectedModel.ID
		}
		// Allow proceeding even if no model selected - will use default later
		m.step = StepDefaults
		return m, nil

	case StepDefaults:
		// Step 3: Save config and finish
		m.finished = true
		return m, m.saveAndFinish()
	}

	return m, nil
}

// getReadmeLinkForDependencies returns the appropriate README link with anchor based on current locale
func getReadmeLinkForDependencies() string {
	locale := locales.GetCurrentLocale()
	baseURL := "https://github.com/lsilvatti/bakasub"

	switch locale {
	case "pt-br":
		return baseURL + "/blob/main/README-pt.md#-dependências"
	case "es":
		return baseURL + "/blob/main/README-es.md#-dependencias"
	default:
		return baseURL + "#-dependencies"
	}
}

// View renders the wizard
func (m Model) View() string {
	if m.finished {
		return ""
	}

	// Wait for terminal size
	if layout.IsWaitingForSize(m.width, m.height) {
		return locales.T("common.loading")
	}

	// Check if terminal is too small
	if layout.IsTooSmall(m.width, m.height) {
		return layout.RenderTooSmallWarning(m.width, m.height)
	}

	contentWidth := m.width - 4

	var content string
	switch m.step {
	case StepLanguageDeps:
		content = m.renderLanguageDepsStep(contentWidth)
	case StepProvider:
		content = m.renderProviderStep(contentWidth)
	case StepDefaults:
		content = m.renderDefaultsStep(contentWidth)
	}

	return content
}

// renderLanguageDepsStep renders step 1: UI Language + Dependencies
func (m Model) renderLanguageDepsStep(contentWidth int) string {
	// Progress bar: Step 1/3 = 33%
	progressWidth := contentWidth - len(locales.T("wizard.title")) - 20
	filledWidth := progressWidth / 3
	emptyWidth := progressWidth - filledWidth
	progressBar := strings.Repeat("█", filledWidth) + strings.Repeat("░", emptyWidth)
	header := styles.HeaderBorder.Width(contentWidth).Render(
		fmt.Sprintf(" %s %s [STEP 1/3] ", locales.T("wizard.title"), progressBar),
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
	} else if m.depDownloading {
		// Show download progress bar
		depsContent.WriteString(fmt.Sprintf("   ┌── %s ──────────────────────────────────────\n", m.currentDep))
		depsContent.WriteString("   │   " + locales.T("wizard.step1.status") + ": [" + locales.T("job.downloading") + "]\n")
		depsContent.WriteString("   │   " + locales.T("job.progress") + ":\n")

		progressWidth := 40
		filled := int(m.downloadProgress * float64(progressWidth))
		empty := progressWidth - filled
		progressBar := strings.Repeat("█", filled) + strings.Repeat("░", empty)
		percent := int(m.downloadProgress * 100)
		depsContent.WriteString(fmt.Sprintf("   │   [%s] %d%%\n", progressBar, percent))

		downloadedMB := m.downloadProgress * 48.0 // Assume 48MB total
		depsContent.WriteString(fmt.Sprintf("   │   %.1fMB / 48.0MB\n", downloadedMB))
		depsContent.WriteString("   └──────────────────────────────────────────────\n\n")
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
				depsContent.WriteString("   │   \n")
				depsContent.WriteString("   │   📖 " + locales.T("wizard.step1.see_readme") + ":\n")
				depsContent.WriteString("   │      " + styles.Highlight.Render(getReadmeLinkForDependencies()) + "\n")
			}

			depsContent.WriteString("   └──────────────────────────────────────────────\n\n")
		}
	}

	depsPanel := styles.Panel.Width(contentWidth).Render(depsContent.String())

	var footerText string
	if !m.checkComplete {
		spacer := strings.Repeat(" ", max(1, contentWidth-50))
		footerText = styles.RenderHotkey("Q", locales.T("common.quit")) + spacer +
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
			spacer := strings.Repeat(" ", max(1, contentWidth-50))
			footerText = styles.RenderHotkey("Q", locales.T("common.quit")) + spacer +
				styles.RenderHotkey("ENTER", locales.T("common.next")+" >")
		} else {
			spacer := strings.Repeat(" ", max(1, contentWidth-70))
			footerText = styles.RenderHotkey("Q", locales.T("common.quit")) + " | " +
				styles.RenderHotkey("R", locales.T("wizard.step1.recheck")) + spacer +
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
func (m Model) renderProviderStep(contentWidth int) string {
	// Progress bar: Step 2/3 = 66%
	progressWidth := contentWidth - len(locales.T("wizard.title")) - 20
	filledWidth := (progressWidth * 2) / 3
	emptyWidth := progressWidth - filledWidth
	progressBar := strings.Repeat("█", filledWidth) + strings.Repeat("░", emptyWidth)
	header := styles.HeaderBorder.Width(contentWidth).Render(
		fmt.Sprintf(" %s %s [STEP 2/3] ", locales.T("wizard.title"), progressBar),
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
				editingStatus = "   " + styles.KeyHintStyle.Render("[E] "+locales.T("wizard.step2.start_editing"))
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
				editingStatus = "   " + styles.KeyHintStyle.Render("[E] "+locales.T("wizard.step2.start_editing")+" | [SPACE] "+showLabel)
			}
		}
		// Show validation status
		validationStatus := ""
		if m.apiKeyInput.Value() != "" {
			if m.validatingKey {
				validationStatus = "   " + styles.StatusWarning.Render("⏳ "+locales.T("wizard.step1.validating")+"...")
			} else if m.keyValidated {
				validationStatus = "   " + styles.StatusOK.Render("✓ "+locales.T("wizard.step1.validation_ok"))
			} else if m.keyValidationErr != "" {
				validationStatus = "   " + styles.StatusError.Render("✗ "+m.keyValidationErr)
			}
		}
		credContent = fmt.Sprintf("2. %s\n", locales.T("wizard.step1.credentials_title")) +
			fmt.Sprintf("   %s > ", locales.T("wizard.step1.api_key_label")) + inputView + validationStatus + "\n" +
			editingStatus
	}

	credPanelStyle := styles.Panel.Width(contentWidth)
	if m.activeSection == 1 {
		credPanelStyle = credPanelStyle.BorderForeground(styles.Yellow)
	}
	credPanel := credPanelStyle.Render(credContent)

	// Model selection panel using reusable component
	var modelContent strings.Builder

	if m.loadingModels {
		modelContent.WriteString(fmt.Sprintf("3. %s\n\n", locales.T("wizard.step2.model_selection_title")))
		modelContent.WriteString("   " + styles.StatusWarning.Render("⏳ "+locales.T("wizard.step2.loading_models")+"...") + "\n")
	} else {
		// Set component width and active state
		m.modelSelector.SetWidth(contentWidth)
		// Active when in model section (for navigation)
		m.modelSelector.SetActive(m.activeSection == 2)

		// Render using component
		modelContent.WriteString(m.modelSelector.View())
	}

	modelPanelStyle := styles.Panel.Width(contentWidth)
	if m.activeSection == 2 {
		modelPanelStyle = modelPanelStyle.BorderForeground(styles.Yellow)
	}
	modelPanel := modelPanelStyle.Render(modelContent.String())

	// Footer with proper spacing
	leftText := styles.RenderHotkey("ESC", "< "+locales.T("common.back"))
	// Show appropriate status for next step
	var rightText string
	canProceed := false
	if m.providerSelection == 3 {
		// Local LLM: needs endpoint AND validation
		canProceed = m.apiEndpointInput.Value() != "" && m.keyValidated && !m.validatingKey
	} else {
		canProceed = m.apiKeyInput.Value() != "" && m.keyValidated && !m.validatingKey
	}
	if canProceed {
		rightText = styles.RenderHotkey("ENTER", locales.T("common.next_step")+" >")
	} else if m.validatingKey {
		rightText = styles.StatusWarning.Render("[ " + locales.T("wizard.step1.validating") + "... ]")
	} else if !m.keyValidated && (m.apiKeyInput.Value() != "" || (m.providerSelection == 3 && m.apiEndpointInput.Value() != "")) {
		rightText = styles.KeyHintStyle.Render("[ " + locales.T("wizard.step2.stop_editing") + " " + locales.T("wizard.step1.validation_required") + " ]")
	} else {
		rightText = styles.KeyHintStyle.Render("[ " + locales.T("wizard.step1.credentials_title") + " " + locales.T("common.required") + " ]")
	}
	footerContent := leftText + strings.Repeat(" ", max(1, contentWidth-lipgloss.Width(leftText)-lipgloss.Width(rightText)-4)) + rightText
	footer := styles.FooterBorder.Width(contentWidth).Render(footerContent)

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
func (m Model) renderDefaultsStep(contentWidth int) string {
	// Progress bar: Step 3/3 = 100%
	progressWidth := contentWidth - len(locales.T("wizard.title")) - 20
	progressBar := strings.Repeat("█", progressWidth)
	header := styles.HeaderBorder.Width(contentWidth).Render(
		fmt.Sprintf(" %s %s [STEP 3/3] ", locales.T("wizard.title"), progressBar),
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
			editingHint = " [E] " + locales.T("wizard.step2.start_editing")
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

	spacer := strings.Repeat(" ", max(1, contentWidth-70))
	footer := styles.FooterBorder.Width(contentWidth).Render(
		styles.RenderHotkey("ESC", "< "+locales.T("common.back")) + " | " +
			styles.RenderHotkey("TAB", locales.T("wizard.step3.next_section")) + spacer +
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

		// Model was already set in handleEnter for StepProvider via modelSelector
		// No need to set a default here as the component handles selection

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

// keyValidatedMsg is sent when API key validation completes
type keyValidatedMsg struct {
	isValid bool
	err     error
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

// validateKeyCmd validates the API key by making a real API request
func (m Model) validateKeyCmd() tea.Cmd {
	return func() tea.Msg {
		// Create context with timeout
		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()

		var provider ai.LLMProvider
		var err error

		switch m.providerSelection {
		case 0: // OpenRouter
			provider = ai.NewOpenRouterAdapter(m.apiKeyInput.Value(), "", 0.3)
			// First check if key works with a simple request
			if !provider.ValidateKey(ctx) {
				return keyValidatedMsg{isValid: false, err: fmt.Errorf("Authentication failed. Please check your API key.")}
			}
			// Then try to list models to ensure full access
			_, err = provider.ListModels(ctx)
			if err != nil {
				// If models fail but auth passed, it might be a network issue
				return keyValidatedMsg{isValid: false, err: fmt.Errorf("API key valid but failed to fetch models: %v", err)}
			}
			return keyValidatedMsg{isValid: true, err: nil}
		case 1: // Gemini
			provider, err = ai.NewGeminiAdapter(ctx, m.apiKeyInput.Value(), "", 0.3)
			if err != nil {
				return keyValidatedMsg{isValid: false, err: err}
			}
		case 2: // OpenAI
			provider = ai.NewOpenAIAdapter(m.apiKeyInput.Value(), "", 0.3)
		case 3: // Local LLM
			// For local, just check if endpoint is reachable
			endpoint := m.apiEndpointInput.Value()
			if endpoint == "" {
				return keyValidatedMsg{isValid: false, err: fmt.Errorf("endpoint required")}
			}
			provider = ai.NewLocalLLMAdapter(endpoint, "", 0.3)
		}

		// Try to list models as validation (for non-OpenRouter providers)
		if provider != nil {
			_, err = provider.ListModels(ctx)
			if err != nil {
				return keyValidatedMsg{isValid: false, err: err}
			}
		}

		return keyValidatedMsg{isValid: true, err: nil}
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
			provider = ai.NewOpenAIAdapter(m.apiKeyInput.Value(), "", 0.3)
			modelStrings, err := provider.ListModels(ctx)
			if err != nil {
				// Fallback to static list if API fails
				modelStrings = []string{"gpt-4o", "gpt-4o-mini", "gpt-4-turbo", "gpt-3.5-turbo"}
			}
			models := parseModelsToModelInfo(modelStrings, false)
			return modelsLoadedMsg{models: models, err: nil}
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

// parseModelsToModelInfo converts string list to ModelInfo list with metadata
func parseModelsToModelInfo(models []string, isFree bool) []ModelInfo {
	result := []ModelInfo{}

	for _, modelStr := range models {
		// Check if model string contains pricing data (format: id|price|context)
		parts := strings.Split(modelStr, "|")

		var modelID string
		var pricePerToken string
		var contextLength string

		if len(parts) >= 3 {
			// Format: openai/gpt-4|0.00003|8192
			modelID = parts[0]
			pricePerToken = parts[1]
			contextLength = parts[2]
		} else {
			// Fallback to old format (just ID)
			modelID = modelStr
			pricePerToken = ""
			contextLength = ""
		}

		// Calculate price per 1M tokens if we have pricing data
		var pricePerM string
		if pricePerToken != "" && pricePerToken != "0" {
			// Convert string to float
			var price float64
			fmt.Sscanf(pricePerToken, "%f", &price)
			// Convert to price per 1M tokens
			pricePerMillion := price * 1000000
			if pricePerMillion < 0.01 {
				pricePerM = "FREE"
			} else {
				pricePerM = fmt.Sprintf("$%.2f", pricePerMillion)
			}
		} else {
			// Fallback to old hardcoded prices
			pricePerM = getPriceForModel(modelID, isFree)
		}

		// Determine context size
		ctx := contextLength
		if ctx == "" || ctx == "0" {
			ctx = inferContextSize(modelID)
		} else {
			// Format context nicely (8192 -> 8k, 128000 -> 128k, etc)
			var ctxInt int
			fmt.Sscanf(ctx, "%d", &ctxInt)
			if ctxInt >= 1000000 {
				ctx = fmt.Sprintf("%dM", ctxInt/1000000)
			} else if ctxInt >= 1000 {
				ctx = fmt.Sprintf("%dk", ctxInt/1000)
			}
		}

		// Force free if ID contains :free tag, regardless of price calculation
		modelIsFree := isFree || pricePerM == "FREE" || strings.Contains(strings.ToLower(modelID), ":free")

		// If marked as free, ensure price is "FREE"
		if modelIsFree && pricePerM != "FREE" {
			pricePerM = "FREE"
		}

		result = append(result, ModelInfo{
			ID:          modelID,
			Name:        modelID,
			ContextSize: ctx,
			PricePerM:   pricePerM,
			IsFree:      modelIsFree,
		})
	}
	return result
}

// inferContextSize tries to infer context size from model name
func inferContextSize(modelID string) string {
	contextSizes := map[string]string{
		"gpt-4o":           "128k",
		"gpt-4o-mini":      "128k",
		"gpt-4-turbo":      "128k",
		"gpt-3.5-turbo":    "16k",
		"gemini-1.5-pro":   "1M",
		"gemini-1.5-flash": "1M",
		"gemini-2.0-flash": "1M",
		"gemini-2.5":       "1M",
		"claude-3-opus":    "200k",
		"claude-3-sonnet":  "200k",
		"claude-3-haiku":   "200k",
		"llama-3.3":        "128k",
		"llama-3.1":        "128k",
		"llama-3":          "8k",
		"mistral":          "32k",
		"mixtral":          "32k",
		"qwen":             "32k",
		"phi-3":            "128k",
		"gemma":            "8k",
	}

	lowerID := strings.ToLower(modelID)
	for key, size := range contextSizes {
		if strings.Contains(lowerID, key) {
			return size
		}
	}
	return "varies"
}

// getPriceForModel returns a price string for a given model
func getPriceForModel(modelID string, isFree bool) string {
	if isFree {
		return "FREE"
	}

	lowerID := strings.ToLower(modelID)

	// Common model pricing - ordered from most specific to least specific
	// This order matters to avoid "gpt-4o" matching before "gpt-4o-mini"
	priceList := []struct {
		key   string
		price string
	}{
		// OpenAI - most specific first
		{"gpt-4o-mini", "$0.15"},
		{"gpt-4o", "$2.50"},
		{"gpt-4-turbo", "$10.00"},
		{"gpt-4", "$30.00"},
		{"gpt-3.5-turbo", "$0.50"},
		// Google Gemini - all variations
		{"gemini-2.5-flash", "$0.07"},
		{"gemini-2.5-pro", "$1.25"},
		{"gemini-2-flash", "$0.07"},
		{"gemini-2-pro", "$1.25"},
		{"gemini-1.5-flash", "$0.07"},
		{"gemini-1.5-pro", "$1.25"},
		{"gemini-flash", "$0.07"},
		{"gemini-pro", "$1.25"},
		{"gemini-3-pro", "$1.25"},
		// Anthropic Claude
		{"claude-3.5-sonnet", "$3.00"},
		{"claude-3-opus", "$15.00"},
		{"claude-3-sonnet", "$3.00"},
		{"claude-3-haiku", "$0.25"},
		// Meta Llama
		{"llama-3.3-70b", "$0.59"},
		{"llama-3.1-405b", "$3.00"},
		{"llama-3.1-70b", "$0.59"},
		{"llama-3.1-8b", "$0.06"},
		{"llama-70b", "$0.59"},
		{"llama-8b", "$0.06"},
		// Others
		{"qwen-2.5-72b", "$0.35"},
		{"mistral-7b", "$0.06"},
		{"mixtral-8x7b", "$0.24"},
		{"phi-3-medium", "$0.14"},
		{"gemma-2-9b", "$0.08"},
	}

	// Check matches in order (most specific first)
	for _, p := range priceList {
		if strings.Contains(lowerID, p.key) {
			return p.price
		}
	}

	// Default for unknown models
	return "$0.50"
}

// Finished returns true if the wizard is complete
func (m Model) Finished() bool {
	return m.finished
}

// Quitting returns true if the user quit the wizard
func (m Model) Quitting() bool {
	return m.quitting
}
