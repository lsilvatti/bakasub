package langselector

import (
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/lsilvatti/bakasub/internal/locales"
	"github.com/lsilvatti/bakasub/internal/ui/focus"
	"github.com/lsilvatti/bakasub/internal/ui/styles"
)

// Language represents a selectable language option
type Language struct {
	Code    string
	Display string
}

// Mode determines what the selector is for
type Mode int

const (
	// ModeUILanguage for interface language selection (en, pt-br, es)
	ModeUILanguage Mode = iota
	// ModeTargetLanguage for translation target language
	ModeTargetLanguage
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

// Model represents the language selector state
type Model struct {
	mode            Mode
	selectedIndex   int
	languages       []Language
	customInput     textinput.Model
	hasCustomOption bool
	focusManager    *focus.Manager
	isActive        bool
	width           int
}

// NewUILanguageSelector creates a selector for interface language
func NewUILanguageSelector(focusManager *focus.Manager) Model {
	return Model{
		mode:          ModeUILanguage,
		selectedIndex: 0,
		languages: []Language{
			{Code: "en", Display: "ENGLISH (Default)"},
			{Code: "pt-br", Display: "PORTUGUÊS (Brasil)"},
			{Code: "es", Display: "ESPAÑOL"},
		},
		focusManager:    focusManager,
		hasCustomOption: false,
		width:           70,
	}
}

// NewTargetLanguageSelector creates a selector for translation target language
func NewTargetLanguageSelector(focusManager *focus.Manager) Model {
	customInput := textinput.New()
	customInput.Placeholder = "pt-br, en-us, ja-jp, etc."
	customInput.CharLimit = 10
	customInput.Width = 30

	return Model{
		mode:          ModeTargetLanguage,
		selectedIndex: 0,
		languages: []Language{
			{Code: "pt-br", Display: "PT-BR (Português)"},
			{Code: "en-us", Display: "EN-US (English)"},
			{Code: "es", Display: "ES (Español)"},
			{Code: "ja-jp", Display: "JA-JP (Japanese)"},
			{Code: "fr-fr", Display: "FR-FR (Français)"},
			{Code: "de", Display: "DE (Deutsch)"},
		},
		customInput:     customInput,
		focusManager:    focusManager,
		hasCustomOption: true,
		width:           70,
	}
}

// SetWidth sets the component width
func (m *Model) SetWidth(width int) {
	m.width = width
	if m.hasCustomOption {
		m.customInput.Width = min(width-30, 40)
	}
}

// SetActive sets whether this component is focused
func (m *Model) SetActive(active bool) {
	m.isActive = active
}

// IsActive returns whether this component is focused
func (m Model) IsActive() bool {
	return m.isActive
}

// SetSelectedByCode selects a language by its code
func (m *Model) SetSelectedByCode(code string) {
	code = strings.ToLower(code)
	for i, lang := range m.languages {
		if strings.ToLower(lang.Code) == code {
			m.selectedIndex = i
			return
		}
	}
	// If not found and has custom option, select "OTHER" and set custom value
	if m.hasCustomOption {
		m.selectedIndex = len(m.languages) // OTHER option
		m.customInput.SetValue(code)
	}
}

// SetCustomValue sets the custom input value
func (m *Model) SetCustomValue(value string) {
	m.customInput.SetValue(value)
}

// GetSelectedCode returns the selected language code
func (m Model) GetSelectedCode() string {
	if m.hasCustomOption && m.selectedIndex == len(m.languages) {
		// OTHER selected - return custom value
		return strings.ToLower(strings.TrimSpace(m.customInput.Value()))
	}
	if m.selectedIndex >= 0 && m.selectedIndex < len(m.languages) {
		return m.languages[m.selectedIndex].Code
	}
	return ""
}

// GetSelectedIndex returns the selected index
func (m Model) GetSelectedIndex() int {
	return m.selectedIndex
}

// IsCustomSelected returns true if "OTHER" is selected
func (m Model) IsCustomSelected() bool {
	return m.hasCustomOption && m.selectedIndex == len(m.languages)
}

// Update handles messages
func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd

	if !m.isActive {
		return m, nil
	}

	// Handle custom input if active
	if m.IsCustomSelected() && m.focusManager != nil && m.focusManager.Mode() == focus.ModeInput {
		if keyMsg, ok := msg.(tea.KeyMsg); ok {
			switch keyMsg.String() {
			case "enter", "esc":
				// Exit input mode
				m.customInput.Blur()
				m.focusManager.ExitInput()
				return m, nil
			}
		}
		m.customInput, cmd = m.customInput.Update(msg)
		return m, cmd
	}

	// Navigation mode
	if keyMsg, ok := msg.(tea.KeyMsg); ok {
		maxIndex := len(m.languages) - 1
		if m.hasCustomOption {
			maxIndex = len(m.languages) // Include OTHER option
		}

		switch keyMsg.String() {
		case "up", "k":
			if m.selectedIndex > 0 {
				m.selectedIndex--
			}
		case "down", "j":
			if m.selectedIndex < maxIndex {
				m.selectedIndex++
			}
		case "e":
			// If OTHER selected, enter input mode with E key
			if m.IsCustomSelected() && m.focusManager != nil {
				m.focusManager.EnterInput(0)
				m.customInput.Focus()
				return m, textinput.Blink
			}
		}
	}

	return m, cmd
}

// FocusCustomInput focuses the custom input field
func (m *Model) FocusCustomInput() tea.Cmd {
	if m.hasCustomOption {
		m.customInput.Focus()
		return textinput.Blink
	}
	return nil
}

// BlurCustomInput removes focus from the custom input
func (m *Model) BlurCustomInput() {
	if m.hasCustomOption {
		m.customInput.Blur()
	}
}

// View renders the component
func (m Model) View() string {
	var content strings.Builder

	for i, lang := range m.languages {
		icon := "( )"
		if i == m.selectedIndex {
			icon = "(o)"
			content.WriteString("   " + styles.Highlight.Render(icon) + " " + styles.Highlight.Render(lang.Display) + "\n")
		} else {
			content.WriteString("   " + icon + " " + lang.Display + "\n")
		}
	}

	// Render OTHER option if applicable
	if m.hasCustomOption {
		icon := "( )"
		otherLabel := locales.T("wizard.step3.other") + " (" + locales.T("wizard.step3.custom_iso") + ")"

		if m.selectedIndex == len(m.languages) {
			icon = "(o)"
			content.WriteString("   " + styles.Highlight.Render(icon) + " " + styles.Highlight.Render(otherLabel) + "\n")
			// Show custom input
			content.WriteString("\n   " + locales.T("wizard.step3.iso_code") + " > " + m.customInput.View())

			// Show language name hint
			langName := GetLanguageName(m.customInput.Value())
			if langName != "" {
				content.WriteString(" → " + styles.StatusOK.Render(langName))
			} else if m.customInput.Value() != "" {
				content.WriteString(" → " + styles.StatusWarning.Render("?"))
			}

			content.WriteString("\n")
			// Show editing hint
			if m.focusManager != nil {
				if m.focusManager.Mode() == focus.ModeInput {
					content.WriteString("   " + styles.StatusOK.Render("["+locales.T("wizard.step2.editing")+"]") + " [ENTER] " + locales.T("wizard.step2.stop_editing") + "\n")
				} else {
					content.WriteString("   " + styles.KeyHintStyle.Render("[E] "+locales.T("wizard.step2.start_editing")) + "\n")
				}
			}
		} else {
			content.WriteString("   " + icon + " " + otherLabel + "\n")
		}
	}

	return content.String()
}

// GetLanguageName returns the language name for an ISO code
func GetLanguageName(code string) string {
	code = strings.ToLower(strings.TrimSpace(code))
	if name, ok := isoLanguageNames[code]; ok {
		return name
	}
	return ""
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
