package main

import (
	"fmt"
	"log"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/lsilvatti/bakasub/internal/ui/attachments"
	"github.com/lsilvatti/bakasub/internal/ui/remuxer"
)

func main() {
	fmt.Println("BakaSub - Phase 3.3: Toolbox Expansion Demo")
	fmt.Println("============================================")
	fmt.Println()

	if len(os.Args) < 2 {
		fmt.Println("Usage: ./mkv-tools-demo <path-to-mkv-file>")
		fmt.Println("\nThis demo requires an actual MKV file to demonstrate:")
		fmt.Println("  1. Attachment Manager - View/Add/Delete embedded fonts/images")
		fmt.Println("  2. Quick Remuxer - Select tracks to keep/remove")
		fmt.Println("\nExample: ./mkv-tools-demo ~/Videos/anime_episode.mkv")
		return
	}

	mkvPath := os.Args[1]

	// Validate file exists
	if _, err := os.Stat(mkvPath); os.IsNotExist(err) {
		fmt.Printf("Error: File not found: %s\n", mkvPath)
		return
	}

	showMenu(mkvPath)
}

func showMenu(mkvPath string) {
	for {
		fmt.Println("\n╔════════════════════════════════════════════════════════╗")
		fmt.Println("║  TOOLBOX DEMO - Select Tool:                          ║")
		fmt.Println("╠════════════════════════════════════════════════════════╣")
		fmt.Println("║  [1] Attachment Manager                                ║")
		fmt.Println("║      • View embedded fonts/images                      ║")
		fmt.Println("║      • Add new attachments                             ║")
		fmt.Println("║      • Extract all to folder                           ║")
		fmt.Println("║      • Delete attachments                              ║")
		fmt.Println("║                                                        ║")
		fmt.Println("║  [2] Quick Remuxer                                     ║")
		fmt.Println("║      • Select tracks to keep                           ║")
		fmt.Println("║      • Remove unwanted audio/subtitle tracks           ║")
		fmt.Println("║      • Create new MKV with selected tracks only        ║")
		fmt.Println("║                                                        ║")
		fmt.Println("║  [q] Quit                                              ║")
		fmt.Println("╚════════════════════════════════════════════════════════╝")
		fmt.Print("\nChoice: ")

		var choice string
		fmt.Scanln(&choice)

		switch choice {
		case "1":
			launchAttachmentManager(mkvPath)
		case "2":
			launchRemuxer(mkvPath)
		case "q", "Q":
			fmt.Println("Goodbye!")
			return
		default:
			fmt.Println("Invalid choice. Try again.")
		}
	}
}

func launchAttachmentManager(mkvPath string) {
	fmt.Println()
	fmt.Println("🔧 Launching Attachment Manager...")
	fmt.Println()

	model, err := attachments.New(mkvPath)
	if err != nil {
		log.Fatalf("Error loading attachments: %v", err)
	}

	p := tea.NewProgram(model, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}

func launchRemuxer(mkvPath string) {
	fmt.Println()
	fmt.Println("🔧 Launching Quick Remuxer...")
	fmt.Println()

	model, err := remuxer.New(mkvPath)
	if err != nil {
		log.Fatalf("Error loading tracks: %v", err)
	}

	p := tea.NewProgram(model, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}
