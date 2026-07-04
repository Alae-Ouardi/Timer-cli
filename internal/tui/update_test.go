package tui

import (
	"testing"
	"time"

	tea "github.com/charmbracelet/bubbletea"

	"timer/internal/store"
	"timer/internal/theme"
)

func isQuit(cmd tea.Cmd) bool {
	if cmd == nil {
		return false
	}
	_, ok := cmd().(tea.QuitMsg)
	return ok
}

func TestUpdateWindowSizeMsg(t *testing.T) {
	fc := newFakeClock()
	m := NewModelWithClock("Test", ModeStopwatch, 0, theme.Default(), fc.now)

	updated, _ := m.Update(tea.WindowSizeMsg{Width: 120, Height: 40})
	nm := updated.(Model)
	if nm.Width != 120 || nm.Height != 40 {
		t.Errorf("Width/Height = %d/%d, want 120/40", nm.Width, nm.Height)
	}
}

func TestUpdateKeyPause(t *testing.T) {
	fc := newFakeClock()
	m := NewModelWithClock("Test", ModeStopwatch, 0, theme.Default(), fc.now)

	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("p")})
	nm := updated.(Model)
	if nm.Running() {
		t.Error("expected Running() = false after 'p'")
	}
}

func TestUpdateKeyContinue(t *testing.T) {
	fc := newFakeClock()
	m := NewModelWithClock("Test", ModeStopwatch, 0, theme.Default(), fc.now).Pause()

	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("c")})
	nm := updated.(Model)
	if !nm.Running() {
		t.Error("expected Running() = true after 'c'")
	}
}

func TestUpdateKeyRestart(t *testing.T) {
	fc := newFakeClock()
	m := NewModelWithClock("Test", ModeStopwatch, 0, theme.Default(), fc.now)
	fc.advance(30 * time.Second)

	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("r")})
	nm := updated.(Model)
	if nm.Elapsed() != 0 {
		t.Errorf("Elapsed() after restart = %v, want 0", nm.Elapsed())
	}
}

func TestUpdateKeyStopFinalizesAndQuits(t *testing.T) {
	t.Setenv("XDG_CONFIG_HOME", t.TempDir())

	fc := newFakeClock()
	m := NewModelWithClock("Deep work", ModeStopwatch, 0, theme.Default(), fc.now)
	fc.advance(10 * time.Second)

	updated, cmd := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("s")})
	nm := updated.(Model)
	if !nm.Done() {
		t.Error("expected Done() = true after 's'")
	}
	if !isQuit(cmd) {
		t.Error("expected 's' to return tea.Quit")
	}

	all, err := store.AllHistory()
	if err != nil {
		t.Fatalf("AllHistory() error: %v", err)
	}
	if len(all) != 1 {
		t.Fatalf("expected 1 history record, got %d", len(all))
	}
	if all[0].Title != "Deep work" || !all[0].Completed {
		t.Errorf("unexpected record: %+v", all[0])
	}
}

func TestUpdateCountdownStoppedEarlyIsNotCompleted(t *testing.T) {
	t.Setenv("XDG_CONFIG_HOME", t.TempDir())

	fc := newFakeClock()
	m := NewModelWithClock("Tea", ModeCountdown, time.Minute, theme.Default(), fc.now)
	fc.advance(10 * time.Second)

	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("s")})
	_ = updated

	all, _ := store.AllHistory()
	if len(all) != 1 || all[0].Completed {
		t.Errorf("expected 1 incomplete record, got %+v", all)
	}
}

func TestUpdateTickCountdownFinishAppendsHistoryOnce(t *testing.T) {
	t.Setenv("XDG_CONFIG_HOME", t.TempDir())

	fc := newFakeClock()
	m := NewModelWithClock("Tea", ModeCountdown, 5*time.Second, theme.Default(), fc.now)
	fc.advance(6 * time.Second)

	updated, cmd := m.Update(tickMsg(fc.now()))
	nm := updated.(Model)
	if !nm.Done() {
		t.Fatal("expected Done() = true after countdown tick past target")
	}
	if cmd != nil {
		t.Error("expected nil cmd (no further ticking) once done")
	}

	all, _ := store.AllHistory()
	if len(all) != 1 || !all[0].Completed {
		t.Errorf("expected 1 completed record, got %+v", all)
	}

	// A further tick or keypress while done must not append a second record.
	updated2, quitCmd := nm.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("x")})
	_ = updated2
	if !isQuit(quitCmd) {
		t.Error("expected any key while done to quit")
	}
	all2, _ := store.AllHistory()
	if len(all2) != 1 {
		t.Errorf("expected history to stay at 1 record, got %d", len(all2))
	}
}

func TestUpdateTickRunningCountdownContinuesTicking(t *testing.T) {
	fc := newFakeClock()
	m := NewModelWithClock("Tea", ModeCountdown, time.Minute, theme.Default(), fc.now)
	fc.advance(5 * time.Second)

	updated, cmd := m.Update(tickMsg(fc.now()))
	nm := updated.(Model)
	if nm.Done() {
		t.Error("expected Done() = false, countdown target not reached")
	}
	if cmd == nil {
		t.Error("expected a tick cmd to be scheduled")
	}
}

func TestUpdateKeyQuitDoesNotFinalize(t *testing.T) {
	t.Setenv("XDG_CONFIG_HOME", t.TempDir())

	fc := newFakeClock()
	m := NewModelWithClock("Test", ModeStopwatch, 0, theme.Default(), fc.now)

	_, cmd := m.Update(tea.KeyMsg{Type: tea.KeyCtrlC})
	if !isQuit(cmd) {
		t.Error("expected ctrl+c to quit")
	}

	all, _ := store.AllHistory()
	if len(all) != 0 {
		t.Errorf("expected no history record on ctrl+c, got %d", len(all))
	}
}
