package main

import (
	"fmt"
	"log"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/lsilvatti/bakasub/internal/core/linter"
	"github.com/lsilvatti/bakasub/internal/ui/header"
)

func main() {
	fmt.Println("BakaSub - Phase 3.2: Header Editor & Quality Gate Demo")
	fmt.Println("=======================================================")
	fmt.Println()

	testLinter()
	fmt.Println("\nPress Enter to launch Quality Gate UI (Ctrl+C to skip)...")
	fmt.Scanln()

	launchQualityGateUI()
}

func testLinter() {
	fmt.Println("1. Testing Quality Gate (Linter):")
	fmt.Println()

	// Test cases with various issues
	testLines := []string{
		"Hello, world!",                      // Clean
		"{\\an8}Text without closing",        // Broken ASS tag
		"This is [incomplete bracket",        // Bracket mismatch
		"What is this???!!!",                 // Excessive punctuation
		"Olá mundo, hello friend",            // English residual (if target is not English)
		"Normal line {\\an8}",                // Clean with tag
		"Another line with ((nested problem", // Nested bracket issue
		"Too many dots........",              // Excessive punctuation
	}

	result := linter.Check(testLines, "por") // Portuguese target, so English words are flagged

	fmt.Printf("   Issues Found: %d\n", len(result.Issues))
	fmt.Printf("   Passed All Checks: %v\n\n", result.PassedAll)

	// Display issues
	for i, issue := range result.Issues {
		fmt.Printf("   Issue %d:\n", i+1)
		fmt.Printf("     Line: %d\n", issue.LineID)
		fmt.Printf("     Severity: %s\n", issue.Severity)
		fmt.Printf("     Type: %s\n", issue.IssueType)
		fmt.Printf("     Content: %s\n", issue.Content)
		fmt.Printf("     Auto-Fixable: %v\n\n", issue.AutoFixable)
	}

	// Test auto-fix
	fmt.Println("2. Testing Auto-Fix:")
	fmt.Println()
	fixed := linter.AutoFix(testLines, result.Issues)

	fmt.Println("   Before -> After:")
	for i, line := range testLines {
		if line != fixed[i] {
			fmt.Printf("     [%d] %s\n", i+1, line)
			fmt.Printf("      -> %s\n", fixed[i])
		}
	}

	// Verify fix reduced issues
	recheck := linter.Check(fixed, "por")
	fmt.Printf("\n   Issues After Fix: %d (was %d)\n", len(recheck.Issues), len(result.Issues))
	fmt.Printf("   ✓ Quality Gate validation complete!\n")
}

func launchQualityGateUI() {
	fmt.Println()
	fmt.Println("3. Launching Quality Gate UI...")
	fmt.Println()

	// Create sample issues for UI demo
	testLines := []string{
		"{\\an8}Olá mundo",
		"This is [incomplete",
		"Hello, tudo bem?",
	}

	result := linter.Check(testLines, "por")
	model := header.NewQualityGate(result)

	p := tea.NewProgram(model, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}
