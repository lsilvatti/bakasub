package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/lsilvatti/bakasub/internal/config"
	"github.com/lsilvatti/bakasub/internal/locales"
	"github.com/lsilvatti/bakasub/internal/ui/wizard"
)

// wizard-demo: Standalone wizard testing tool
// Usage: go run cmd/wizard-demo/main.go

func main() {
	// Initialize locales first
	if err := locales.Init(); err != nil {
		fmt.Fprintf(os.Stderr, "Warning: Failed to initialize locales: %v\n", err)
	}

	// Use a default config as the base
	cfg := config.Default()

	// Create wizard
	wiz := wizard.New(cfg)

	// Run wizard
	p := tea.NewProgram(wiz, tea.WithAltScreen())
	finalModel, err := p.Run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	// Check result
	wizModel := finalModel.(wizard.Model)
	if wizModel.Quitting() {
		fmt.Println("\n✗ Wizard cancelled by user")
		os.Exit(0)
	}

	if wizModel.Finished() {
		fmt.Println("\n✓ Wizard completed successfully!")
		fmt.Println("  Config saved to config.json")
		fmt.Println("  Run 'bin/bakasub' to start the application")
	} else {
		fmt.Println("\n? Wizard exited unexpectedly")
	}
}
