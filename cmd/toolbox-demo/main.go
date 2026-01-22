package main

import (
	"fmt"
	"log"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/lsilvatti/bakasub/internal/ui/glossary"
)

func main() {
	fmt.Println("BakaSub - Phase 3.1: Toolbox Demo")
	fmt.Println("==================================")
	fmt.Println()

	testGlossary()
	fmt.Println("\nPress Enter to launch Glossary UI (Ctrl+C to skip)...")
	fmt.Scanln()

	launchGlossaryUI()
}

func testGlossary() {
	fmt.Println("1. Testing Auto-Detection:")

	testSub := `1
00:00:01,000 --> 00:00:03,000
Hello, my name is Luffy!

2
00:00:04,000 --> 00:00:06,000
I'm going to become the Pirate King!

3
00:00:07,000 --> 00:00:09,000
Nami and Zoro are my nakama.`

	tmpFile := "/tmp/test_sub.srt"
	os.WriteFile(tmpFile, []byte(testSub), 0644)
	defer os.Remove(tmpFile)

	entries, err := glossary.AutoDetectTerms(tmpFile)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("   ✓ Detected %d unique terms:\n", len(entries))
	for _, e := range entries {
		fmt.Printf("     - %s (%s)\n", e.Original, e.Type)
	}
}

func launchGlossaryUI() {
	fmt.Println("\n2. Launching Glossary UI...")

	glossaryPath := "/tmp/bakasub_glossary.json"
	model := glossary.New(glossaryPath)

	p := tea.NewProgram(model, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}
