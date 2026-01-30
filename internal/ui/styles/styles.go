package styles

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// Native Neon Design System
// Inspired by btop and lazygit aesthetics

// Color Definitions
var (
	// Primary Neon Colors
	NeonPink = lipgloss.Color("#F700FF")
	Cyan     = lipgloss.Color("#00FFFF")
	Yellow   = lipgloss.Color("#FFFF00")

	// Neutral/Terminal Colors
	Gray     = lipgloss.Color("#808080")
	DarkGray = lipgloss.Color("#404040")
)

// Border Styles
var (
	// Standard rounded border for panels
	RoundedBorder = lipgloss.RoundedBorder()

	// Double border for main container/window
	DoubleBorder = lipgloss.DoubleBorder()

	// Normal border for secondary elements
	NormalBorder = lipgloss.NormalBorder()
)

// Base Styles
var (
	// MainWindow - Double border, no background (terminal default)
	MainWindow = lipgloss.NewStyle().
			Border(DoubleBorder).
			BorderForeground(NeonPink).
			Padding(0, 1)

	// Panel - Rounded border, transparent background
	Panel = lipgloss.NewStyle().
		Border(RoundedBorder).
		BorderForeground(Cyan).
		Padding(0, 1)

	// PanelTitle - Accent header for panels
	PanelTitle = lipgloss.NewStyle().
			Foreground(Yellow).
			Bold(true)

	// Hotkey - Styled key hints like [ k ]
	Hotkey = lipgloss.NewStyle().
		Foreground(NeonPink).
		Bold(true)

	// HotkeyText - Description text next to hotkeys
	HotkeyText = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFFFF"))

	// StatusOK - Success/Online indicators
	StatusOK = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#00FF00")).
			Bold(true)

	// StatusError - Error/Offline indicators
	StatusError = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FF0000")).
			Bold(true)

	// StatusWarning - Warning indicators
	StatusWarning = lipgloss.NewStyle().
			Foreground(Yellow).
			Bold(true)

	// Highlight - Selected/Active items
	Highlight = lipgloss.NewStyle().
			Foreground(NeonPink).
			Background(DarkGray).
			Bold(true)

	// HighlightStyle - Alternative name for Highlight
	HighlightStyle = Highlight

	// HeaderBorder - Styled header for wizards/modals
	HeaderBorder = lipgloss.NewStyle().
			Border(DoubleBorder).
			BorderForeground(NeonPink).
			Padding(0, 1).
			Bold(true)

	// FooterBorder - Styled footer for wizards/modals
	FooterBorder = lipgloss.NewStyle().
			Border(NormalBorder).
			BorderForeground(Cyan).
			Padding(0, 1)

	// ModalStyle - Modal/Dialog styling with double border
	ModalStyle = lipgloss.NewStyle().
			Border(DoubleBorder).
			BorderForeground(NeonPink).
			Padding(1, 2)

	// ProgressBarFilled - Block character █ style
	ProgressBarFilled = lipgloss.NewStyle().
				Foreground(Cyan)

	// ProgressBarEmpty - Block character ░ style
	ProgressBarEmpty = lipgloss.NewStyle().
				Foreground(DarkGray)

	// Logo - ASCII art title styling
	Logo = lipgloss.NewStyle().
		Foreground(NeonPink).
		Bold(true)

	// Footer - Bottom status bar
	Footer = lipgloss.NewStyle().
		Foreground(Gray)

	// Dimmed - Secondary/helper text
	Dimmed = lipgloss.NewStyle().
		Foreground(Gray)

	// CodeBlock - Monospace content (paths, logs)
	CodeBlock = lipgloss.NewStyle().
			Foreground(Cyan)

	// TitleStyle - Modal/Screen titles
	TitleStyle = lipgloss.NewStyle().
			Foreground(NeonPink).
			Bold(true)

	// SectionStyle - Section headers within panels
	SectionStyle = lipgloss.NewStyle().
			Foreground(Yellow).
			Bold(true)

	// KeyHintStyle - Keyboard shortcut hints
	KeyHintStyle = lipgloss.NewStyle().
			Foreground(NeonPink).
			Bold(true)

	// SuccessStyle - Success messages
	SuccessStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#00FF00"))

	// ErrorStyle - Error messages
	ErrorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FF0000"))

	// WarningStyle - Warning messages
	WarningStyle = lipgloss.NewStyle().
			Foreground(Yellow)

	// AppStyle - Main application container
	AppStyle = lipgloss.NewStyle().
			Padding(1, 2)

	// NeonCyan - Neon cyan accent color
	NeonCyan = Cyan

	// SubtleStyle - Subtle/dimmed text
	SubtleStyle = lipgloss.NewStyle().
			Foreground(Gray)

	// InfoBoxStyle - Information box styling
	InfoBoxStyle = lipgloss.NewStyle().
			Border(RoundedBorder).
			BorderForeground(Cyan).
			Padding(0, 1)

	// DisabledStyle - Disabled UI elements
	DisabledStyle = lipgloss.NewStyle().
			Foreground(DarkGray)

	// AccentStyle - Accent/highlight for interactive elements
	AccentStyle = lipgloss.NewStyle().
			Foreground(Cyan).
			Bold(true)
)

// Helper Functions

// RenderHotkey creates a styled hotkey indicator like "[ k ] KEY"
func RenderHotkey(key, description string) string {
	return Hotkey.Render("[ "+key+" ]") + " " + HotkeyText.Render(description)
}

// RenderProgressBar creates a block-based progress bar
// Example: [████████░░░░] 65%
func RenderProgressBar(current, total int, width int) string {
	if total == 0 {
		return ""
	}

	percentage := float64(current) / float64(total)
	if percentage > 1.0 {
		percentage = 1.0
	}
	filled := int(percentage * float64(width))

	var bar string
	for i := 0; i < width; i++ {
		if i < filled {
			bar += ProgressBarFilled.Render("█")
		} else {
			bar += ProgressBarEmpty.Render("░")
		}
	}

	percent := int(percentage * 100)
	percentStr := fmt.Sprintf("%3d%%", percent)

	return fmt.Sprintf("[%s] %s", bar, lipgloss.NewStyle().Foreground(Yellow).Render(percentStr))
}

// AdaptToSize adjusts a style's width/height dynamically
func AdaptToSize(style lipgloss.Style, width, height int) lipgloss.Style {
	return style.Width(width).Height(height)
}

// RenderDimmedOverlay creates a dimmed background effect using block characters
// This is used for modal overlays to indicate background is inactive
func RenderDimmedOverlay(width, height int) string {
	var lines []string
	dimChar := "▒"
	line := strings.Repeat(dimChar, width)
	dimStyle := lipgloss.NewStyle().Foreground(DarkGray)

	for i := 0; i < height; i++ {
		lines = append(lines, dimStyle.Render(line))
	}

	return strings.Join(lines, "\n")
}

// RenderModalWithOverlay renders a modal centered over a dimmed background
func RenderModalWithOverlay(modal string, width, height int) string {
	overlay := RenderDimmedOverlay(width, height)
	modalLines := strings.Split(modal, "\n")

	// Calculate center position
	modalHeight := len(modalLines)
	modalWidth := 0
	for _, line := range modalLines {
		if len(line) > modalWidth {
			modalWidth = len(line)
		}
	}

	startY := (height - modalHeight) / 2
	startX := (width - modalWidth) / 2

	if startY < 0 {
		startY = 0
	}
	if startX < 0 {
		startX = 0
	}

	// Overlay the modal on the dimmed background
	overlayLines := strings.Split(overlay, "\n")
	for i, line := range modalLines {
		lineIdx := startY + i
		if lineIdx >= 0 && lineIdx < len(overlayLines) {
			// Replace portion of overlay with modal line
			padding := strings.Repeat(" ", startX)
			overlayLines[lineIdx] = padding + line
		}
	}

	return strings.Join(overlayLines, "\n")
}

// DenseTableStyles returns table styles with minimal padding for btop aesthetic
func DenseTableStyles() (lipgloss.Style, lipgloss.Style, lipgloss.Style) {
	// Header style: bold, colored, minimal padding
	headerStyle := lipgloss.NewStyle().
		Foreground(Yellow).
		Bold(true).
		Padding(0, 1)

	// Cell style: minimal padding
	cellStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FFFFFF")).
		Padding(0, 1)

	// Selected cell style: highlighted
	selectedStyle := lipgloss.NewStyle().
		Foreground(NeonPink).
		Background(DarkGray).
		Bold(true).
		Padding(0, 1)

	return headerStyle, cellStyle, selectedStyle
}

// RenderSection creates a styled section box with title using Panel border
func RenderSection(title string, content string) string {
	titleLine := SectionStyle.Render(title)
	fullContent := "\n" + titleLine + "\n\n" + content + "\n"
	return Panel.Render(fullContent)
}

// RenderInfoBox creates a simple info box with rounded border
func RenderInfoBox(content string) string {
	return Panel.Render(content)
}
