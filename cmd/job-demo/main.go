package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/lsilvatti/bakasub/internal/config"
	"github.com/lsilvatti/bakasub/internal/ui/job"
)

func main() {
	// Load or create default config
	cfg := &config.Config{
		TargetLang:  "pt-br",
		Model:       "gemini-1.5-flash",
		Temperature: 0.3,
		AIProvider:  "openrouter",
	}

	// Get input path from args or use default
	inputPath := "./test_media"
	if len(os.Args) > 1 {
		inputPath = os.Args[1]
	}

	// Create job setup model
	m := job.New(cfg, inputPath)

	// Run the program
	p := tea.NewProgram(m, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error running job setup: %v\n", err)
		os.Exit(1)
	}
}
