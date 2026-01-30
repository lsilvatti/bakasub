package layout

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// Terminal size constraints
const (
	MinWidth  = 80
	MinHeight = 24
)

// CalculateHalf returns half the available width minus padding
func CalculateHalf(totalWidth int, padding int) int {
	return (totalWidth / 2) - padding
}

// CalculateThird returns one third of the available width minus padding
func CalculateThird(totalWidth int, padding int) int {
	return (totalWidth / 3) - padding
}

// CalculateTwoThirds returns two thirds of the available width minus padding
func CalculateTwoThirds(totalWidth int, padding int) int {
	return (totalWidth * 2 / 3) - padding
}

// CalculateQuarter returns one quarter of the available width minus padding
func CalculateQuarter(totalWidth int, padding int) int {
	return (totalWidth / 4) - padding
}

// PlaceInPanel uses lipgloss.Place to center content within a fixed-size panel
// and ensures overflow is gracefully truncated
func PlaceInPanel(content string, width, height int, hAlign, vAlign lipgloss.Position) string {
	return lipgloss.Place(width, height, hAlign, vAlign, content)
}

// TruncateToWidth truncates a string to fit within a specific width,
// adding ellipsis if needed
func TruncateToWidth(s string, maxWidth int) string {
	if maxWidth <= 0 {
		return ""
	}

	// Simple character-based truncation (does not handle multi-width runes)
	if len(s) <= maxWidth {
		return s
	}

	if maxWidth <= 3 {
		return s[:maxWidth]
	}

	return s[:maxWidth-3] + "..."
}

// TruncateLines truncates content to a maximum number of lines
func TruncateLines(content string, maxLines int) string {
	if maxLines <= 0 {
		return ""
	}

	lines := strings.Split(content, "\n")
	if len(lines) <= maxLines {
		return content
	}

	return strings.Join(lines[:maxLines], "\n")
}

// IsTooSmall checks if the terminal is below minimum dimensions
// Returns false if dimensions are 0 (still waiting for WindowSizeMsg)
func IsTooSmall(width, height int) bool {
	// Don't show "too small" warning while waiting for size - return false
	// so that View functions can show their own loading state
	if width == 0 || height == 0 {
		return false
	}
	return width < MinWidth || height < MinHeight
}

// IsWaitingForSize checks if we're still waiting for terminal size
func IsWaitingForSize(width, height int) bool {
	return width == 0 || height == 0
}

// RenderTooSmallWarning renders a warning screen when terminal is too small
func RenderTooSmallWarning(width, height int) string {
	// Define colors inline to avoid circular dependency
	neonPink := lipgloss.Color("#F700FF")
	cyan := lipgloss.Color("#00FFFF")
	gray := lipgloss.Color("#808080")

	warning := lipgloss.NewStyle().
		Foreground(neonPink).
		Bold(true).
		Render("TERMINAL TOO SMALL")

	msg := lipgloss.NewStyle().
		Foreground(gray).
		Render("Please resize to at least 80x24")

	currentSize := lipgloss.NewStyle().
		Foreground(cyan).
		Render(lipgloss.JoinVertical(
			lipgloss.Center,
			"",
			"Current Size:",
			lipgloss.NewStyle().Faint(true).Render(
				fmt.Sprintf("Width: %d | Height: %d", width, height),
			),
			"",
			lipgloss.NewStyle().Faint(true).Render("(Minimum: 80x24)"),
		))

	content := lipgloss.JoinVertical(
		lipgloss.Center,
		"",
		warning,
		"",
		msg,
		"",
		currentSize,
		"",
	)

	// Center everything in the available space
	return lipgloss.Place(
		width,
		height,
		lipgloss.Center,
		lipgloss.Center,
		content,
	)
}

// PanelStyle returns a bordered panel style with dynamic width and height
func PanelStyle(width, height int) lipgloss.Style {
	cyan := lipgloss.Color("#00FFFF")

	return lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(cyan).
		Width(width).
		Height(height).
		Padding(0, 1)
}

// TitleBar creates a horizontal title bar with dynamic width
func TitleBar(title string, width int) string {
	neonPink := lipgloss.Color("#F700FF")

	titleStyle := lipgloss.NewStyle().
		Foreground(neonPink).
		Bold(true).
		Width(width).
		Align(lipgloss.Center)

	return titleStyle.Render(title)
}

// ProgressBar creates a simple ASCII progress bar
func ProgressBar(current, total, width int) string {
	neonPink := lipgloss.Color("#F700FF")

	if width <= 0 || total <= 0 {
		return ""
	}

	percentage := float64(current) / float64(total)
	filled := int(percentage * float64(width))

	if filled > width {
		filled = width
	}

	bar := strings.Repeat("█", filled) + strings.Repeat("░", width-filled)

	return lipgloss.NewStyle().
		Foreground(neonPink).
		Render(bar)
}

// SafeWidth ensures width is never negative or below a minimum
func SafeWidth(width, minWidth int) int {
	if width < minWidth {
		return minWidth
	}
	return width
}

// SafeHeight ensures height is never negative or below a minimum
func SafeHeight(height, minHeight int) int {
	if height < minHeight {
		return minHeight
	}
	return height
}

// CalculateContentArea calculates the available content area after accounting for borders and padding
func CalculateContentArea(totalWidth, totalHeight int, borderWidth, padding int) (int, int) {
	// Border width is typically 2 (left + right)
	// Padding is typically 1 on each side (2 total)
	contentWidth := totalWidth - (borderWidth * 2) - (padding * 2)
	contentHeight := totalHeight - (borderWidth * 2) - (padding * 2)

	return SafeWidth(contentWidth, 10), SafeHeight(contentHeight, 5)
}
