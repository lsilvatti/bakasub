package utils

import (
	"fmt"
	"os"
	"runtime/debug"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

const (
	Version = "v1.0.0"
	RepoURL = "https://github.com/lsilvatti/bakasub"
)

var (
	bsodStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFFFF")).
			Background(lipgloss.Color("#0000AA")).
			Bold(true)

	errorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FF0000")).
			Bold(true)
)

// RecoverPanic is a global panic handler that displays a BSOD-style error screen
func RecoverPanic() {
	if r := recover(); r != nil {
		// Clear screen
		fmt.Print("\033[2J\033[H")

		// Render BSOD
		renderBSOD(r)

		// Exit with error code
		os.Exit(1)
	}
}

func renderBSOD(panicValue interface{}) {
	width := 80

	// Build the BSOD screen
	var b strings.Builder

	// Top border
	b.WriteString(strings.Repeat("═", width))
	b.WriteString("\n")

	// Title
	title := "CRITICAL SYSTEM ERROR"
	padding := (width - len(title)) / 2
	b.WriteString(strings.Repeat(" ", padding))
	b.WriteString(errorStyle.Render(title))
	b.WriteString("\n\n")

	// Error details
	b.WriteString(centerText("BakaSub has encountered a critical error and needs to close.", width))
	b.WriteString("\n\n")

	// Panic message
	panicMsg := fmt.Sprintf("%v", panicValue)
	b.WriteString(errorStyle.Render("Error Details:"))
	b.WriteString("\n")
	b.WriteString(wrapText(panicMsg, width-4, "  "))
	b.WriteString("\n\n")

	// Stack trace
	stack := string(debug.Stack())
	b.WriteString(errorStyle.Render("Stack Trace:"))
	b.WriteString("\n")
	stackLines := strings.Split(stack, "\n")

	// Show first 10 lines of stack trace
	displayLines := 10
	if len(stackLines) < displayLines {
		displayLines = len(stackLines)
	}

	for i := 0; i < displayLines; i++ {
		if len(stackLines[i]) > width-4 {
			b.WriteString("  " + stackLines[i][:width-7] + "...")
		} else {
			b.WriteString("  " + stackLines[i])
		}
		b.WriteString("\n")
	}

	if len(stackLines) > displayLines {
		b.WriteString(fmt.Sprintf("  ... and %d more lines\n", len(stackLines)-displayLines))
	}

	b.WriteString("\n")

	// Help text
	b.WriteString(centerText("The application has crashed. A log has been saved.", width))
	b.WriteString("\n\n")

	// GitHub link
	issueURL := RepoURL + "/issues/new"
	b.WriteString(centerText("Please report this issue:", width))
	b.WriteString("\n")
	b.WriteString(centerText(issueURL, width))
	b.WriteString("\n\n")

	// Actions
	b.WriteString(centerText("[ g ] Open GitHub Issue     [ q ] Exit Application", width))
	b.WriteString("\n")

	// Bottom border
	b.WriteString(strings.Repeat("═", width))

	// Render with BSOD style
	fmt.Println(bsodStyle.Render(b.String()))

	// Wait for user input
	var input string
	fmt.Scanln(&input)

	if input == "g" || input == "G" {
		fmt.Printf("\nOpening: %s\n", issueURL)
		fmt.Println("(Please open this URL in your browser)")
	}
}

func centerText(text string, width int) string {
	if len(text) >= width {
		return text
	}
	padding := (width - len(text)) / 2
	return strings.Repeat(" ", padding) + text
}

func wrapText(text string, width int, indent string) string {
	words := strings.Fields(text)
	if len(words) == 0 {
		return ""
	}

	var lines []string
	var currentLine string

	for _, word := range words {
		if len(currentLine)+len(word)+1 > width {
			lines = append(lines, indent+currentLine)
			currentLine = word
		} else {
			if currentLine != "" {
				currentLine += " "
			}
			currentLine += word
		}
	}

	if currentLine != "" {
		lines = append(lines, indent+currentLine)
	}

	return strings.Join(lines, "\n")
}

// SafeRun wraps a function with panic recovery
func SafeRun(fn func()) {
	defer RecoverPanic()
	fn()
}
