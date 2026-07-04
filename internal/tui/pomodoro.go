package tui

import (
	"fmt"
	"time"

	tea "github.com/charmbracelet/bubbletea"

	"timer/internal/notify"
	"timer/internal/store"
	"timer/internal/theme"
)

// Phase identifies which half of a pomodoro round is active.
type Phase string

const (
	PhaseWork  Phase = "work"
	PhaseBreak Phase = "break"
)

// PomodoroModel sequences alternating work/break countdown phases,
// auto-transitioning between them and tracking a round counter. Each
// completed round (one work phase plus one break phase) is persisted to
// history with Mode "pomodoro".
type PomodoroModel struct {
	Title         string
	Theme         theme.Theme
	WorkDuration  time.Duration
	BreakDuration time.Duration
	TotalRounds   int
	CurrentRound  int
	Phase         Phase
	Inner         Model
	Finished      bool

	Now func() time.Time
}

// NewPomodoroModel creates a running pomodoro using the real system clock.
func NewPomodoroModel(title string, work, brk time.Duration, rounds int, th theme.Theme) PomodoroModel {
	return NewPomodoroModelWithClock(title, work, brk, rounds, th, time.Now)
}

// NewPomodoroModelWithClock creates a running pomodoro using the given
// clock function.
func NewPomodoroModelWithClock(title string, work, brk time.Duration, rounds int, th theme.Theme, now func() time.Time) PomodoroModel {
	return PomodoroModel{
		Title:         title,
		Theme:         th,
		WorkDuration:  work,
		BreakDuration: brk,
		TotalRounds:   rounds,
		CurrentRound:  1,
		Phase:         PhaseWork,
		Inner:         NewModelWithClock(phaseTitle(title, PhaseWork, 1, rounds), ModeCountdown, work, th, now),
		Now:           now,
	}
}

func phaseTitle(title string, phase Phase, round, total int) string {
	label := "Work"
	if phase == PhaseBreak {
		label = "Break"
	}
	if title != "" {
		return fmt.Sprintf("%s — %s %d/%d", title, label, round, total)
	}
	return fmt.Sprintf("%s %d/%d", label, round, total)
}

func (m PomodoroModel) clock() time.Time {
	if m.Now != nil {
		return m.Now()
	}
	return time.Now()
}

// Init starts the periodic tick that drives the pomodoro.
func (m PomodoroModel) Init() tea.Cmd {
	return tickCmd()
}

// Update advances the current phase, transitions between work and break
// (recording completed rounds to history), and handles the s/p/c/r/q
// keybindings.
func (m PomodoroModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		inner, _ := m.Inner.Update(msg)
		m.Inner = inner.(Model)
		return m, nil

	case tickMsg:
		if m.Finished {
			return m, nil
		}
		if !m.Inner.IsCountdownFinished() {
			return m, tickCmd()
		}
		return m.advancePhase()

	case tea.KeyMsg:
		if m.Finished {
			return m, tea.Quit
		}
		switch msg.String() {
		case "s":
			_ = store.AppendHistory(m.roundRecord(false))
			return m, tea.Quit
		case "p":
			m.Inner = m.Inner.Pause()
			return m, nil
		case "c":
			m.Inner = m.Inner.Resume()
			return m, nil
		case "r":
			m.Inner = m.Inner.Restart()
			return m, tickCmd()
		case "q", "ctrl+c":
			return m, tea.Quit
		}
	}
	return m, nil
}

// advancePhase transitions from work->break, or from break->(next round or
// finished), notifying and persisting a record once per completed round.
func (m PomodoroModel) advancePhase() (tea.Model, tea.Cmd) {
	if m.Phase == PhaseWork {
		notify.Notify(m.Title, "Work phase complete — break time!")
		m.Phase = PhaseBreak
		m.Inner = NewModelWithClock(phaseTitle(m.Title, PhaseBreak, m.CurrentRound, m.TotalRounds), ModeCountdown, m.BreakDuration, m.Theme, m.Now)
		return m, tickCmd()
	}

	_ = store.AppendHistory(m.roundRecord(true))

	if m.CurrentRound >= m.TotalRounds {
		notify.Notify(m.Title, "Pomodoro complete!")
		m.Finished = true
		return m, nil
	}

	notify.Notify(m.Title, "Break complete — next round!")
	m.CurrentRound++
	m.Phase = PhaseWork
	m.Inner = NewModelWithClock(phaseTitle(m.Title, PhaseWork, m.CurrentRound, m.TotalRounds), ModeCountdown, m.WorkDuration, m.Theme, m.Now)
	return m, tickCmd()
}

func (m PomodoroModel) roundRecord(completed bool) store.HistoryRecord {
	return store.HistoryRecord{
		Title:     m.Title,
		Mode:      "pomodoro",
		Target:    m.WorkDuration + m.BreakDuration,
		Elapsed:   m.Inner.Elapsed(),
		StartedAt: m.Inner.StartedAt,
		EndedAt:   m.clock(),
		Completed: completed,
	}
}

// View delegates rendering to the current phase's inner Model, or shows a
// completion screen once every round has finished.
func (m PomodoroModel) View() string {
	if m.Finished {
		return m.renderFinished()
	}
	return m.Inner.View()
}

func (m PomodoroModel) renderFinished() string {
	done := m.Inner.MarkDone()
	done.Title = fmt.Sprintf("%s — Pomodoro complete (%d/%d rounds)", nonEmpty(m.Title, "Pomodoro"), m.TotalRounds, m.TotalRounds)
	done.Width, done.Height = m.Inner.Width, m.Inner.Height
	return done.View()
}

func nonEmpty(s, fallback string) string {
	if s == "" {
		return fallback
	}
	return s
}
