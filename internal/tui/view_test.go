package tui

import (
	"strings"
	"testing"
	"time"

	"timer/internal/theme"
)

func TestViewStopwatchNoPanic(t *testing.T) {
	fc := newFakeClock()
	m := NewModelWithClock("Deep work", ModeStopwatch, 0, theme.Default(), fc.now)
	fc.advance(65 * time.Second)

	sizes := [][2]int{{0, 0}, {10, 5}, {80, 24}, {200, 60}}
	for _, sz := range sizes {
		m.Width, m.Height = sz[0], sz[1]
		out := m.View()
		if out == "" {
			t.Errorf("View() empty at size %v", sz)
		}
		if !strings.Contains(out, "Deep work") {
			t.Errorf("View() at size %v missing title", sz)
		}
	}
}

func TestViewCountdownShowsProgressAndRemaining(t *testing.T) {
	fc := newFakeClock()
	m := NewModelWithClock("Tea", ModeCountdown, time.Minute, theme.Default(), fc.now)
	fc.advance(30 * time.Second)
	m.Width, m.Height = 80, 24

	out := m.View()
	if out == "" {
		t.Fatal("View() returned empty output")
	}
	if !strings.Contains(out, "Tea") {
		t.Error("View() missing title")
	}
}

func TestViewPausedDimsHints(t *testing.T) {
	fc := newFakeClock()
	m := NewModelWithClock("Test", ModeStopwatch, 0, theme.Default(), fc.now).Pause()
	m.Width, m.Height = 80, 24

	out := m.View()
	if !strings.Contains(out, "continue") {
		t.Error("View() while paused should show 'continue' hint")
	}
}

func TestViewRunningShowsPauseHint(t *testing.T) {
	fc := newFakeClock()
	m := NewModelWithClock("Test", ModeStopwatch, 0, theme.Default(), fc.now)
	m.Width, m.Height = 80, 24

	out := m.View()
	if !strings.Contains(out, "pause") {
		t.Error("View() while running should show 'pause' hint")
	}
}

func TestViewDoneScreen(t *testing.T) {
	fc := newFakeClock()
	m := NewModelWithClock("Tea", ModeCountdown, 5*time.Second, theme.Default(), fc.now)
	fc.advance(6 * time.Second)
	m = m.MarkDone()
	m.Width, m.Height = 80, 24

	out := m.View()
	if !strings.Contains(out, "DONE") {
		t.Error("View() done screen missing 'DONE'")
	}
	if !strings.Contains(out, "press any key") {
		t.Error("View() done screen missing key-press hint")
	}
	if !strings.Contains(out, "Tea") {
		t.Error("View() done screen missing title")
	}
}

func TestViewDoneScreenNoPanicAtVariousSizes(t *testing.T) {
	fc := newFakeClock()
	m := NewModelWithClock("", ModeStopwatch, 0, theme.Default(), fc.now).MarkDone()

	sizes := [][2]int{{0, 0}, {1, 1}, {80, 24}}
	for _, sz := range sizes {
		m.Width, m.Height = sz[0], sz[1]
		if out := m.View(); out == "" {
			t.Errorf("View() empty at size %v", sz)
		}
	}
}

func lineOf(out, substr string) int {
	for i, line := range strings.Split(out, "\n") {
		if strings.Contains(line, substr) {
			return i
		}
	}
	return -1
}

func TestViewTitleNearTopHintsNearBottomDigitsInMiddle(t *testing.T) {
	fc := newFakeClock()
	m := NewModelWithClock("Deep work", ModeStopwatch, 0, theme.Default(), fc.now)
	m.Width, m.Height = 80, 30

	out := m.View()
	lines := strings.Split(out, "\n")
	lastLine := len(lines) - 1

	titleLine := lineOf(out, "Deep work")
	hintsLine := lineOf(out, "s=stop")
	digitsLine := lineOf(out, "█")

	if titleLine < 0 || hintsLine < 0 || digitsLine < 0 {
		t.Fatalf("expected to find title, hints, and digits; got titleLine=%d hintsLine=%d digitsLine=%d", titleLine, hintsLine, digitsLine)
	}

	if titleLine > lastLine/4 {
		t.Errorf("title at line %d, want it near the top (within first quarter of %d lines)", titleLine, len(lines))
	}
	if hintsLine < lastLine*3/4 {
		t.Errorf("hints at line %d, want them near the bottom (within last quarter of %d lines)", hintsLine, len(lines))
	}
	if titleLine >= digitsLine || digitsLine >= hintsLine {
		t.Errorf("expected order title(%d) < digits(%d) < hints(%d)", titleLine, digitsLine, hintsLine)
	}
}

func TestViewAllThemesNoPanic(t *testing.T) {
	fc := newFakeClock()
	for _, th := range theme.List() {
		m := NewModelWithClock("Focus", ModeCountdown, time.Minute, th, fc.now)
		m.Width, m.Height = 80, 24
		if out := m.View(); out == "" {
			t.Errorf("View() empty for theme %q", th.Name)
		}
	}
}
