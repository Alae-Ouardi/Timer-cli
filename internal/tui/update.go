package tui

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"

	"timer/internal/notify"
	"timer/internal/store"
)

const tickInterval = 250 * time.Millisecond

type tickMsg time.Time

func tickCmd() tea.Cmd {
	return tea.Tick(tickInterval, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

// Init starts the periodic tick that drives timer updates.
func (m Model) Init() tea.Cmd {
	return tickCmd()
}

// Update implements the Bubble Tea update loop: it advances the clock,
// detects countdown completion, handles resize, and dispatches the
// s/p/c/r/q keybindings.
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.Width = msg.Width
		m.Height = msg.Height
		return m, nil

	case tickMsg:
		if m.done {
			return m, nil
		}
		if m.ModeVal == ModeCountdown && m.IsCountdownFinished() {
			m = m.MarkDone()
			notify.Notify(m.Title, "Time's up!")
			_ = store.AppendHistory(m.toHistoryRecord())
			return m, nil
		}
		return m, tickCmd()

	case tea.KeyMsg:
		if m.done {
			return m, tea.Quit
		}
		switch msg.String() {
		case "s":
			m = m.MarkDone()
			_ = store.AppendHistory(m.toHistoryRecord())
			return m, tea.Quit
		case "p":
			return m.Pause(), nil
		case "c":
			return m.Resume(), nil
		case "r":
			return m.Restart(), tickCmd()
		case "q", "ctrl+c":
			return m, tea.Quit
		}
	}
	return m, nil
}

// toHistoryRecord builds the persisted record for this session. A
// countdown is Completed only if it actually reached its target;
// stopwatches are always considered completed once stopped.
func (m Model) toHistoryRecord() store.HistoryRecord {
	completed := true
	if m.ModeVal == ModeCountdown && !m.IsCountdownFinished() {
		completed = false
	}
	return store.HistoryRecord{
		Title:     m.Title,
		Mode:      string(m.ModeVal),
		Target:    m.Target,
		Elapsed:   m.Elapsed(),
		StartedAt: m.StartedAt,
		EndedAt:   m.clock(),
		Completed: completed,
	}
}
