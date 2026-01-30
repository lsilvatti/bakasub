package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/lsilvatti/bakasub/internal/config"
	"github.com/lsilvatti/bakasub/internal/locales"
	"github.com/lsilvatti/bakasub/internal/ui/dashboard"
	"github.com/lsilvatti/bakasub/internal/ui/wizard"
	"github.com/lsilvatti/bakasub/pkg/utils"
)

func main() {
	// Handle --version flag
	if len(os.Args) > 1 && (os.Args[1] == "--version" || os.Args[1] == "-v") {
		fmt.Printf("BakaSub %s\n", utils.Version)
		return
	}

	// Wrap entire application with panic recovery (BSOD handler)
	utils.SafeRun(func() {
		// Check if config exists
		if !config.Exists() {
			// Initialize with default English before wizard
			if err := locales.Init(); err != nil {
				fmt.Fprintf(os.Stderr, "Warning: Failed to initialize locales: %v\n", err)
			}

			// Config doesn't exist - run wizard
			fmt.Println("BakaSub - First Run Setup")
			fmt.Println("=========================")
			fmt.Println()
			fmt.Println("No configuration found. Starting setup wizard...")
			fmt.Println()

			runWizard()
			return
		}

		// Config exists - load it first
		cfg, err := config.Load()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error loading config: %v\n", err)
			os.Exit(1)
		}

		// Now initialize locales with the loaded config
		if err := locales.Load(cfg.InterfaceLang); err != nil {
			fmt.Fprintf(os.Stderr, "Warning: Failed to load language '%s': %v\n", cfg.InterfaceLang, err)
		}

		runDashboard(cfg)
	})
}

func runWizard() {
	cfg := config.Default()
	wiz := wizard.New(cfg)

	p := tea.NewProgram(wiz, tea.WithAltScreen())
	finalModel, err := p.Run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	// Check wizard result
	wizModel := finalModel.(wizard.Model)
	if wizModel.Quitting() {
		fmt.Println("\n✗ Setup cancelled. Run 'bakasub' again to retry.")
		os.Exit(0)
	}

	if wizModel.Finished() {
		fmt.Println("\n✓ Setup complete! Launching BakaSub...")
		fmt.Println()

		// Reload config and launch dashboard
		cfg, err := config.Load()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error loading saved config: %v\n", err)
			os.Exit(1)
		}

		// Load the configured language
		if err := locales.Load(cfg.InterfaceLang); err != nil {
			fmt.Fprintf(os.Stderr, "Warning: Failed to load language '%s': %v\n", cfg.InterfaceLang, err)
		}

		runDashboard(cfg)
	}
}

func runDashboard(cfg *config.Config) {
	dashModel := dashboard.New(cfg)

	p := tea.NewProgram(dashModel, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
