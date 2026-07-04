package tui

import (
	"testing"
	"time"

	tea "github.com/charmbracelet/bubbletea"

	"timer/internal/store"
	"timer/internal/theme"
)

func TestPomodoroStartsInWorkPhaseRoundOne(t *testing.T) {
	fc := newFakeClock()
	m := NewPomodoroModelWithClock("Focus", 5*time.Second, 3*time.Second, 2, theme.Default(), fc.now)

	if m.Phase != PhaseWork {
		t.Errorf("Phase = %q, want %q", m.Phase, PhaseWork)
	}
	if m.CurrentRound != 1 {
		t.Errorf("CurrentRound = %d, want 1", m.CurrentRound)
	}
}

func TestPomodoroAdvancesWorkToBreak(t *testing.T) {
	fc := newFakeClock()
	m := NewPomodoroModelWithClock("Focus", 5*time.Second, 3*time.Second, 2, theme.Default(), fc.now)

	fc.advance(6 * time.Second)
	updated, _ := m.Update(tickMsg(fc.now()))
	nm := updated.(PomodoroModel)

	if nm.Phase != PhaseBreak {
		t.Errorf("Phase = %q, want %q", nm.Phase, PhaseBreak)
	}
	if nm.CurrentRound != 1 {
		t.Errorf("CurrentRound = %d, want 1 (still same round)", nm.CurrentRound)
	}
}

func TestPomodoroAdvancesBreakToNextRound(t *testing.T) {
	t.Setenv("XDG_CONFIG_HOME", t.TempDir())

	fc := newFakeClock()
	m := NewPomodoroModelWithClock("Focus", 5*time.Second, 3*time.Second, 2, theme.Default(), fc.now)

	fc.advance(6 * time.Second)
	updated, _ := m.Update(tickMsg(fc.now()))
	m = updated.(PomodoroModel)

	fc.advance(4 * time.Second)
	updated, _ = m.Update(tickMsg(fc.now()))
	m = updated.(PomodoroModel)

	if m.Phase != PhaseWork {
		t.Errorf("Phase = %q, want %q", m.Phase, PhaseWork)
	}
	if m.CurrentRound != 2 {
		t.Errorf("CurrentRound = %d, want 2", m.CurrentRound)
	}
	if m.Finished {
		t.Error("Finished = true, want false (2 rounds requested)")
	}

	all, _ := store.AllHistory()
	if len(all) != 1 {
		t.Fatalf("expected 1 round recorded, got %d", len(all))
	}
	if all[0].Mode != "pomodoro" || !all[0].Completed {
		t.Errorf("unexpected record: %+v", all[0])
	}
}

func TestPomodoroFinishesAfterLastRound(t *testing.T) {
	t.Setenv("XDG_CONFIG_HOME", t.TempDir())

	fc := newFakeClock()
	m := NewPomodoroModelWithClock("Focus", 5*time.Second, 3*time.Second, 1, theme.Default(), fc.now)

	fc.advance(6 * time.Second)
	updated, _ := m.Update(tickMsg(fc.now()))
	m = updated.(PomodoroModel)

	fc.advance(4 * time.Second)
	updated, cmd := m.Update(tickMsg(fc.now()))
	m = updated.(PomodoroModel)

	if !m.Finished {
		t.Error("expected Finished = true after last round's break completes")
	}
	if cmd != nil {
		t.Error("expected nil cmd once finished")
	}

	all, _ := store.AllHistory()
	if len(all) != 1 {
		t.Fatalf("expected 1 round recorded, got %d", len(all))
	}
}

func TestPomodoroKeyStopFinalizesIncomplete(t *testing.T) {
	t.Setenv("XDG_CONFIG_HOME", t.TempDir())

	fc := newFakeClock()
	m := NewPomodoroModelWithClock("Focus", time.Minute, 30*time.Second, 3, theme.Default(), fc.now)
	fc.advance(10 * time.Second)

	_, cmd := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("s")})
	if !isQuit(cmd) {
		t.Error("expected 's' to quit")
	}

	all, _ := store.AllHistory()
	if len(all) != 1 || all[0].Completed {
		t.Errorf("expected 1 incomplete record, got %+v", all)
	}
}

func TestPomodoroKeyPauseResume(t *testing.T) {
	fc := newFakeClock()
	m := NewPomodoroModelWithClock("Focus", time.Minute, 30*time.Second, 1, theme.Default(), fc.now)

	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("p")})
	m = updated.(PomodoroModel)
	if m.Inner.Running() {
		t.Error("expected paused after 'p'")
	}

	updated, _ = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("c")})
	m = updated.(PomodoroModel)
	if !m.Inner.Running() {
		t.Error("expected running after 'c'")
	}
}

func TestPomodoroFinishedAnyKeyQuits(t *testing.T) {
	fc := newFakeClock()
	m := NewPomodoroModelWithClock("Focus", 5*time.Second, 3*time.Second, 1, theme.Default(), fc.now)
	m.Finished = true

	_, cmd := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("x")})
	if !isQuit(cmd) {
		t.Error("expected any key to quit when finished")
	}
}

func TestPomodoroViewNoPanic(t *testing.T) {
	fc := newFakeClock()
	m := NewPomodoroModelWithClock("Focus", 5*time.Second, 3*time.Second, 2, theme.Default(), fc.now)
	m.Inner.Width, m.Inner.Height = 80, 24

	if out := m.View(); out == "" {
		t.Error("View() returned empty output while running")
	}

	m.Finished = true
	if out := m.View(); out == "" {
		t.Error("View() returned empty output while finished")
	}
}

func TestPomodoroWindowSizePropagatesToInner(t *testing.T) {
	fc := newFakeClock()
	m := NewPomodoroModelWithClock("Focus", 5*time.Second, 3*time.Second, 1, theme.Default(), fc.now)

	updated, _ := m.Update(tea.WindowSizeMsg{Width: 100, Height: 30})
	nm := updated.(PomodoroModel)
	if nm.Inner.Width != 100 || nm.Inner.Height != 30 {
		t.Errorf("Inner size = %d/%d, want 100/30", nm.Inner.Width, nm.Inner.Height)
	}
}
