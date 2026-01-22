package main

import (
	"fmt"
	"math/rand"
	"os"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/lsilvatti/bakasub/internal/ui/execution"
)

// Simulate translation job with realistic anime dialogue
var samplePairs = []struct {
	original   string
	translated string
}{
	{"Don't underestimate the sea, kid!", "Não subestime o mar, garoto!"},
	{"I'm going to be the King of the Pirates!", "Eu vou ser o Rei dos Piratas!"},
	{"Believe it! That's my ninja way!", "Acredite! Esse é o meu jeito ninja!"},
	{"Plus Ultra!", "Plus Ultra!"},
	{"This is the power of friendship!", "Este é o poder da amizade!"},
	{"I'll take a potato chip... and EAT IT!", "Vou pegar uma batata frita... e COMER!"},
	{"You're already dead.", "Você já está morto."},
	{"Kamehameha!", "Kamehameha!"},
	{"I want to be the very best!", "Eu quero ser o melhor!"},
	{"Notice me, senpai!", "Me note, senpai!"},
	{"This isn't even my final form!", "Esta nem é minha forma final!"},
	{"The cake is a lie!", "O bolo é uma mentira!"},
	{"I am the one who knocks!", "Eu sou aquele que bate!"},
	{"Winter is coming.", "O inverno está chegando."},
	{"May the Force be with you.", "Que a Força esteja com você."},
	{"To infinity and beyond!", "Ao infinito e além!"},
	{"Why so serious?", "Por que tão sério?"},
	{"I'll be back.", "Eu voltarei."},
	{"Show me what you got!", "Mostre-me o que você tem!"},
	{"Get schwifty!", "Fique schwifty!"},
	{"Wubba lubba dub dub!", "Wubba lubba dub dub!"},
	{"That's all folks!", "É isso pessoal!"},
	{"What's up, doc?", "E aí, doutor?"},
	{"Cowabunga!", "Cowabunga!"},
	{"It's morphin' time!", "É hora de transformar!"},
	{"By the power of Grayskull!", "Pelo poder de Grayskull!"},
	{"Hasta la vista, baby.", "Hasta la vista, baby."},
	{"I see dead people.", "Eu vejo pessoas mortas."},
	{"You shall not pass!", "Você não passará!"},
	{"My precious...", "Meu precioso..."},
	{"Do or do not, there is no try.", "Faça ou não faça, não há tentativa."},
	{"I am your father.", "Eu sou seu pai."},
	{"Houston, we have a problem.", "Houston, temos um problema."},
	{"E.T. phone home.", "E.T. ligar para casa."},
	{"Here's Johnny!", "Aqui está o Johnny!"},
	{"You talking to me?", "Está falando comigo?"},
	{"Say hello to my little friend!", "Diga olá ao meu pequeno amigo!"},
	{"Life is like a box of chocolates.", "A vida é como uma caixa de chocolates."},
	{"I'll have what she's having.", "Eu vou querer o que ela está tomando."},
	{"Keep your friends close, but your enemies closer.", "Mantenha seus amigos perto, mas seus inimigos mais perto."},
}

type model struct {
	execModel  execution.Model
	index      int
	ticker     *time.Ticker
	quitting   bool
	progress   float64
	fileIndex  int
	totalFiles int
}

func initialModel() model {
	execModel := execution.New("One Piece - Season 1", 24)
	execModel.SetSize(100, 35)

	return model{
		execModel:  execModel,
		index:      0,
		ticker:     time.NewTicker(800 * time.Millisecond), // Add pair every 800ms
		progress:   0.0,
		fileIndex:  1,
		totalFiles: 24,
	}
}

func (m model) Init() tea.Cmd {
	return tea.Batch(
		m.execModel.Init(),
		tickCmd(),
	)
}

type tickMsg time.Time

func tickCmd() tea.Cmd {
	return tea.Tick(800*time.Millisecond, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			if m.execModel.View() == "" {
				m.quitting = true
				return m, tea.Quit
			}
		case "esc":
			m.quitting = true
			return m, tea.Quit
		}

	case tickMsg:
		if m.index < len(samplePairs) {
			// Add new translation pair
			pair := samplePairs[m.index]
			updated, tCmd := m.execModel.Update(execution.TranslationMsg{
				ID:           m.index + 1,
				OriginalText: pair.original,
				Translated:   pair.translated,
			})
			m.execModel = updated.(execution.Model)
			cmds = append(cmds, tCmd)

			// Update progress
			m.progress = float64(m.index+1) / float64(len(samplePairs)) * 100
			updated, tCmd = m.execModel.Update(execution.ProgressMsg{
				FileProgress:  m.progress,
				BatchProgress: float64(m.fileIndex) / float64(m.totalFiles) * 100,
				CurrentFile:   fmt.Sprintf("episode_%03d.mkv", m.fileIndex),
			})
			m.execModel = updated.(execution.Model)
			cmds = append(cmds, tCmd)

			// Add random log messages
			if rand.Intn(3) == 0 {
				logMsgs := []struct {
					level   execution.LogLevel
					message string
				}{
					{execution.LogInfo, fmt.Sprintf("Batch %d sent to AI provider...", m.index/10+1)},
					{execution.LogAI, fmt.Sprintf("Gemini Flash 1.5 responded in 1.2s")},
					{execution.LogSuccess, "Quality gate passed: All ASS tags preserved"},
					{execution.LogWarn, "Minor glossary mismatch detected (auto-fixed)"},
				}
				logMsg := logMsgs[rand.Intn(len(logMsgs))]
				updated, tCmd = m.execModel.Update(execution.LogMsg{
					Level:   logMsg.level,
					Message: logMsg.message,
				})
				m.execModel = updated.(execution.Model)
				cmds = append(cmds, tCmd)
			}

			// Update stats
			updated, tCmd = m.execModel.Update(execution.StatsMsg{
				LinesProcessed: m.index + 1,
				TokensUsed:     (m.index + 1) * 250,
				CostSoFar:      float64(m.index+1) * 0.000175,
				Errors:         0,
			})
			m.execModel = updated.(execution.Model)
			cmds = append(cmds, tCmd)

			m.index++

			// Continue ticking
			if m.index < len(samplePairs) {
				cmds = append(cmds, tickCmd())
			} else {
				// Job complete
				updated, tCmd = m.execModel.Update(execution.StatusMsg{
					Status: execution.StatusComplete,
				})
				m.execModel = updated.(execution.Model)
				cmds = append(cmds, tCmd)

				updated, tCmd = m.execModel.Update(execution.LogMsg{
					Level:   execution.LogSuccess,
					Message: "Translation completed successfully! All files processed.",
				})
				m.execModel = updated.(execution.Model)
				cmds = append(cmds, tCmd)
			}
		}

	case tea.WindowSizeMsg:
		m.execModel.SetSize(msg.Width, msg.Height)
	}

	// Update execution model
	updated, tCmd := m.execModel.Update(msg)
	m.execModel = updated.(execution.Model)
	cmds = append(cmds, tCmd)

	return m, tea.Batch(cmds...)
}

func (m model) View() string {
	if m.quitting {
		return "👋 Demo finished. Thanks for watching!\n"
	}

	return m.execModel.View()
}

func main() {
	// Seed random for varied log messages
	rand.Seed(time.Now().UnixNano())

	p := tea.NewProgram(initialModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error running demo: %v\n", err)
		os.Exit(1)
	}
}
