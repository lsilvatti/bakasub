package main

import (
	"fmt"
	"log"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/lsilvatti/bakasub/internal/ui/picker"
)

func main() {
	fmt.Println("BakaSub - File Picker Demo")
	fmt.Println("===========================")
	fmt.Println()
	fmt.Println("Controls:")
	fmt.Println("  ↑/↓    Navigate")
	fmt.Println("  Enter  Open directory / Select file")
	fmt.Println("  Tab    Toggle mode (Dir/File/Both)")
	fmt.Println("  Esc    Cancel")
	fmt.Println()
	fmt.Println("Press Enter to start...")
	fmt.Scanln()

	// Get starting directory
	startDir, err := os.UserHomeDir()
	if err != nil {
		startDir = "/"
	}

	// Create picker in directory mode
	model := picker.New(startDir, picker.ModeDirectory)

	p := tea.NewProgram(model, tea.WithAltScreen())

	finalModel, err := p.Run()
	if err != nil {
		log.Fatalf("Error running picker: %v", err)
	}

	// Check for selection
	if m, ok := finalModel.(picker.Model); ok {
		fmt.Println()
		fmt.Println("Picker closed.")
		fmt.Println("Selected path:", m.SelectedPath())
		fmt.Println("Current directory:", m.CurrentDirectory())
	}
}
