package tui

import (
	"testing"
	"time"

	"timer/internal/theme"
)

type fakeClock struct {
	t time.Time
}

func (f *fakeClock) now() time.Time { return f.t }

func (f *fakeClock) advance(d time.Duration) { f.t = f.t.Add(d) }

func newFakeClock() *fakeClock {
	return &fakeClock{t: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)}
}

func TestElapsedAccumulatesWhileRunning(t *testing.T) {
	fc := newFakeClock()
	m := NewModelWithClock("Test", ModeStopwatch, 0, theme.Default(), fc.now)

	fc.advance(5 * time.Second)
	if got := m.Elapsed(); got != 5*time.Second {
		t.Errorf("Elapsed() = %v, want %v", got, 5*time.Second)
	}
}

func TestPauseFreezesElapsed(t *testing.T) {
	fc := newFakeClock()
	m := NewModelWithClock("Test", ModeStopwatch, 0, theme.Default(), fc.now)

	fc.advance(5 * time.Second)
	m = m.Pause()
	fc.advance(10 * time.Second)

	if got := m.Elapsed(); got != 5*time.Second {
		t.Errorf("Elapsed() after pause = %v, want %v (paused time should not count)", got, 5*time.Second)
	}
	if m.Running() {
		t.Error("Running() = true after Pause(), want false")
	}
}

func TestResumeContinuesAccumulating(t *testing.T) {
	fc := newFakeClock()
	m := NewModelWithClock("Test", ModeStopwatch, 0, theme.Default(), fc.now)

	fc.advance(5 * time.Second)
	m = m.Pause()
	fc.advance(10 * time.Second) // should not count
	m = m.Resume()
	fc.advance(3 * time.Second)

	if got := m.Elapsed(); got != 8*time.Second {
		t.Errorf("Elapsed() after resume = %v, want %v", got, 8*time.Second)
	}
	if !m.Running() {
		t.Error("Running() = false after Resume(), want true")
	}
}

func TestRestartZeroesAccumulated(t *testing.T) {
	fc := newFakeClock()
	m := NewModelWithClock("Test", ModeStopwatch, 0, theme.Default(), fc.now)

	fc.advance(30 * time.Second)
	m = m.Restart()

	if got := m.Elapsed(); got != 0 {
		t.Errorf("Elapsed() after restart = %v, want 0", got)
	}

	fc.advance(2 * time.Second)
	if got := m.Elapsed(); got != 2*time.Second {
		t.Errorf("Elapsed() after restart+advance = %v, want %v", got, 2*time.Second)
	}
}

func TestCountdownCrossingZeroBecomesFinished(t *testing.T) {
	fc := newFakeClock()
	m := NewModelWithClock("Tea", ModeCountdown, 5*time.Second, theme.Default(), fc.now)

	fc.advance(4 * time.Second)
	if m.IsCountdownFinished() {
		t.Error("IsCountdownFinished() = true before target reached")
	}
	if got := m.Remaining(); got != 1*time.Second {
		t.Errorf("Remaining() = %v, want %v", got, 1*time.Second)
	}

	fc.advance(2 * time.Second) // total 6s > 5s target
	if !m.IsCountdownFinished() {
		t.Error("IsCountdownFinished() = false after target reached")
	}
	if got := m.Remaining(); got != 0 {
		t.Errorf("Remaining() = %v, want 0 (clamped)", got)
	}
}

func TestMarkDoneFinalizes(t *testing.T) {
	fc := newFakeClock()
	m := NewModelWithClock("Tea", ModeCountdown, 5*time.Second, theme.Default(), fc.now)
	fc.advance(6 * time.Second)

	m = m.MarkDone()
	if !m.Done() {
		t.Error("Done() = false after MarkDone()")
	}
	if m.Running() {
		t.Error("Running() = true after MarkDone()")
	}

	fc.advance(100 * time.Second)
	if got := m.Elapsed(); got != 6*time.Second {
		t.Errorf("Elapsed() after done = %v, want frozen at %v", got, 6*time.Second)
	}
}

func TestPauseResumeNoOpWhenDone(t *testing.T) {
	fc := newFakeClock()
	m := NewModelWithClock("Test", ModeStopwatch, 0, theme.Default(), fc.now)
	m = m.MarkDone()

	if p := m.Pause(); p.Running() != m.Running() {
		t.Error("Pause() on done model should be a no-op")
	}
	if r := m.Resume(); r.Running() != m.Running() {
		t.Error("Resume() on done model should be a no-op")
	}
}

func TestStopwatchNeverFinished(t *testing.T) {
	fc := newFakeClock()
	m := NewModelWithClock("Test", ModeStopwatch, 0, theme.Default(), fc.now)
	fc.advance(1000 * time.Hour)

	if m.IsCountdownFinished() {
		t.Error("IsCountdownFinished() = true for a stopwatch (Target=0)")
	}
}
