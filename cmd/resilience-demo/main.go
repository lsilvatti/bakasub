package main

import (
	"fmt"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/lsilvatti/bakasub/pkg/utils"
)

type model struct {
	updateStatus    string
	updateAvailable bool
	latestVersion   string
	releaseURL      string
	checking        bool
}

func initialModel() model {
	return model{
		updateStatus: "Checking for updates...",
		checking:     true,
	}
}

func (m model) Init() tea.Cmd {
	// Start update check
	return utils.CheckForUpdates(utils.Version)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "q" || msg.String() == "ctrl+c" {
			return m, tea.Quit
		}
		if msg.String() == "1" {
			// Trigger panic
			panic("User requested panic test")
		}

	case utils.MsgUpdateAvailable:
		m.checking = false
		m.updateAvailable = true
		m.latestVersion = msg.LatestVersion
		m.releaseURL = msg.ReleaseURL
		m.updateStatus = fmt.Sprintf("✓ Update available: %s → %s", msg.CurrentVersion, msg.LatestVersion)

	case utils.MsgUpdateCheckFailed:
		m.checking = false
		m.updateStatus = fmt.Sprintf("✗ Update check failed: %v", msg.Err)
	}

	return m, nil
}

func (m model) View() string {
	s := "BakaSub - Phase 3.4: Panic Handler & Update Checker Demo\n"
	s += "=========================================================\n\n"

	// Update status
	s += "📡 UPDATE CHECKER:\n"
	s += "   Current Version: " + utils.Version + "\n"

	if m.checking {
		s += "   Status: " + m.updateStatus + "\n"
	} else if m.updateAvailable {
		s += "   Status: " + m.updateStatus + "\n"
		s += "   Download: " + m.releaseURL + "\n"
	} else {
		s += "   Status: " + m.updateStatus + "\n"
	}

	s += "\n"

	// Panic test
	s += "💥 PANIC HANDLER:\n"
	s += "   Press [1] to trigger a test panic\n"
	s += "   This will demonstrate the BSOD (Blue Screen of Death) error screen\n"

	s += "\n"
	s += "[q] Quit\n"

	return s
}

func main() {
	fmt.Println("Initializing demo with panic recovery...")
	fmt.Println()

	// Wrap main execution with panic handler
	utils.SafeRun(func() {
		p := tea.NewProgram(initialModel())
		if _, err := p.Run(); err != nil {
			fmt.Printf("Error: %v\n", err)
		}
	})
}

// Test functions for panic scenarios
func testPanicScenarios() {
	fmt.Println()
	fmt.Println("Testing Panic Scenarios:")
	fmt.Println("========================")
	fmt.Println()

	// Test 1: Nil pointer dereference
	fmt.Println("1. Testing nil pointer dereference...")
	utils.SafeRun(func() {
		var ptr *string
		fmt.Println(*ptr) // This will panic
	})
	time.Sleep(1 * time.Second)

	// Test 2: Index out of range
	fmt.Println("\n2. Testing index out of range...")
	utils.SafeRun(func() {
		arr := []int{1, 2, 3}
		fmt.Println(arr[10]) // This will panic
	})
	time.Sleep(1 * time.Second)

	// Test 3: Explicit panic
	fmt.Println("\n3. Testing explicit panic...")
	utils.SafeRun(func() {
		panic("This is a test panic!")
	})
}
