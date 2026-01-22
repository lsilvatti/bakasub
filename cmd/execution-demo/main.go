package main

import (
	"fmt"
	"math/rand"
	"os"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/lsilvatti/bakasub/internal/ui/execution"
)

// simulateJobMsg is sent periodically to simulate job progress
type simulateJobMsg struct{}

func simulateJob() tea.Cmd {
	return tea.Tick(500*time.Millisecond, func(t time.Time) tea.Msg {
		return simulateJobMsg{}
	})
}

type model struct {
	execution execution.Model
	simStep   int
}

func main() {
	// Create execution model
	execModel := execution.New("One Piece Batch (Episodes 1-24)", 24)

	m := model{
		execution: execModel,
		simStep:   0,
	}

	p := tea.NewProgram(m, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}

func (m model) Init() tea.Cmd {
	return tea.Batch(m.execution.Init(), simulateJob())
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		// Pass through to execution model
		execModel, cmd := m.execution.Update(msg)
		m.execution = execModel.(execution.Model)
		return m, cmd

	case simulateJobMsg:
		// Simulate different stages of job execution
		m.simStep++

		switch m.simStep {
		case 1:
			cmds = append(cmds, func() tea.Msg {
				return execution.LogMsg{Level: execution.LogInfo, Message: "Preflight check passed. Found 450 dialogue lines."}
			})
		case 2:
			cmds = append(cmds, func() tea.Msg {
				return execution.LogMsg{Level: execution.LogInfo, Message: "Context: Anime Mode active. Temp: 0.7"}
			})
		case 3:
			cmds = append(cmds, func() tea.Msg {
				return execution.LogMsg{Level: execution.LogAI, Message: "Batch 1 (Lines 1-50) sent to OpenRouter (Gemini Flash)."}
			})
			cmds = append(cmds, func() tea.Msg {
				return execution.ProgressMsg{
					FileProgress:  10.0,
					BatchProgress: 4.0,
					CurrentFile:   "One Piece E1080.mkv",
				}
			})
		case 5:
			cmds = append(cmds, func() tea.Msg {
				return execution.LogMsg{Level: execution.LogSuccess, Message: "Batch 1 received. Sanity Check: PASSED."}
			})
		case 6:
			cmds = append(cmds, func() tea.Msg {
				return execution.LogMsg{Level: execution.LogAI, Message: "Batch 2 (Lines 51-100) sent..."}
			})
			cmds = append(cmds, func() tea.Msg {
				return execution.ProgressMsg{
					FileProgress:  22.0,
					BatchProgress: 4.0,
					CurrentFile:   "One Piece E1080.mkv",
				}
			})
		case 10:
			cmds = append(cmds, func() tea.Msg {
				return execution.LogMsg{Level: execution.LogWarn, Message: "Batch 2 Desync (Expected 50 IDs, got 48)."}
			})
		case 11:
			cmds = append(cmds, func() tea.Msg {
				return execution.LogMsg{Level: execution.LogInfo, Message: "└─ Engaging Anti-Desync Protocol (Split Strategy)."}
			})
		case 12:
			cmds = append(cmds, func() tea.Msg {
				return execution.LogMsg{Level: execution.LogInfo, Message: "└─ Split 2a (Lines 51-75) sent... OK."}
			})
		case 14:
			cmds = append(cmds, func() tea.Msg {
				return execution.LogMsg{Level: execution.LogInfo, Message: "└─ Split 2b (Lines 76-100) sent... OK."}
			})
		case 16:
			cmds = append(cmds, func() tea.Msg {
				return execution.LogMsg{Level: execution.LogAI, Message: "Batch 3 (Lines 101-150) sent..."}
			})
			cmds = append(cmds, func() tea.Msg {
				return execution.ProgressMsg{
					FileProgress:  35.0,
					BatchProgress: 4.0,
					CurrentFile:   "One Piece E1080.mkv",
				}
			})
		case 18:
			cmds = append(cmds, func() tea.Msg {
				return execution.LogMsg{Level: execution.LogSuccess, Message: "Batch 3 received."}
			})
		case 20:
			cmds = append(cmds, func() tea.Msg {
				return execution.LogMsg{Level: execution.LogInfo, Message: "Processing Styles: Preserving 'MainDialogue' events."}
			})

		// Rapid fire logs to test circular buffer and auto-scroll
		case 25, 26, 27, 28, 29, 30, 31, 32, 33, 34:
			cmds = append(cmds, func() tea.Msg {
				return execution.LogMsg{
					Level:   execution.LogAI,
					Message: fmt.Sprintf("Batch %d processing... (Simulated rapid logs)", m.simStep-20),
				}
			})
			cmds = append(cmds, func() tea.Msg {
				return execution.ProgressMsg{
					FileProgress:  float64(40 + (m.simStep-25)*5),
					BatchProgress: 4.0,
					CurrentFile:   "One Piece E1080.mkv",
				}
			})

		case 40:
			cmds = append(cmds, func() tea.Msg {
				return execution.LogMsg{Level: execution.LogSuccess, Message: "File 1/24 complete!"}
			})
			cmds = append(cmds, func() tea.Msg {
				return execution.ProgressMsg{
					FileProgress:  100.0,
					BatchProgress: 4.2,
					CurrentFile:   "One Piece E1080.mkv",
				}
			})

		case 45:
			// Start next file
			cmds = append(cmds, func() tea.Msg {
				return execution.LogMsg{Level: execution.LogInfo, Message: "Starting file 2/24: One Piece E1081.mkv"}
			})
			cmds = append(cmds, func() tea.Msg {
				return execution.ProgressMsg{
					FileProgress:  0.0,
					BatchProgress: 8.3,
					CurrentFile:   "One Piece E1081.mkv",
				}
			})

		case 50:
			// Simulate error
			cmds = append(cmds, func() tea.Msg {
				return execution.LogMsg{Level: execution.LogError, Message: "Failed to parse line 125 - malformed ASS tag"}
			})
			cmds = append(cmds, func() tea.Msg {
				return execution.StatsMsg{
					LinesProcessed: 450,
					TokensUsed:     125000,
					CostSoFar:      0.08,
					Errors:         1,
				}
			})

		case 55:
			cmds = append(cmds, func() tea.Msg {
				return execution.LogMsg{Level: execution.LogInfo, Message: "Auto-recovery: Skipping malformed line"}
			})

		case 60:
			// Generate many logs to test buffer limit
			for i := 0; i < 20; i++ {
				logLevel := execution.LogInfo
				if rand.Intn(5) == 0 {
					logLevel = execution.LogAI
				}
				cmds = append(cmds, func() tea.Msg {
					return execution.LogMsg{
						Level:   logLevel,
						Message: fmt.Sprintf("Processing batch segment %d...", i),
					}
				})
			}

		case 70:
			cmds = append(cmds, func() tea.Msg {
				return execution.ProgressMsg{
					FileProgress:  50.0,
					BatchProgress: 8.3,
					CurrentFile:   "One Piece E1081.mkv",
				}
			})
			cmds = append(cmds, func() tea.Msg {
				return execution.StatsMsg{
					LinesProcessed: 900,
					TokensUsed:     250000,
					CostSoFar:      0.15,
					Errors:         1,
				}
			})
		}

		// Update stats periodically
		if m.simStep%10 == 0 {
			cmds = append(cmds, func() tea.Msg {
				return execution.StatsMsg{
					LinesProcessed: m.simStep * 10,
					TokensUsed:     m.simStep * 3000,
					CostSoFar:      float64(m.simStep) * 0.002,
					Errors:         rand.Intn(2),
				}
			})
		}

		// Continue simulation unless we've gone far enough
		if m.simStep < 80 {
			cmds = append(cmds, simulateJob())
		} else {
			// Job complete
			cmds = append(cmds, func() tea.Msg {
				return execution.StatusMsg{Status: execution.StatusComplete}
			})
			cmds = append(cmds, func() tea.Msg {
				return execution.LogMsg{Level: execution.LogSuccess, Message: "All files processed successfully!"}
			})
		}

	default:
		// Pass other messages to execution model
		execModel, cmd := m.execution.Update(msg)
		m.execution = execModel.(execution.Model)
		return m, cmd
	}

	// Update execution model
	var cmd tea.Cmd
	execModel, cmd := m.execution.Update(msg)
	m.execution = execModel.(execution.Model)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m model) View() string {
	return m.execution.View()
}
