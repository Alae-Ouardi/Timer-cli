package tui

import (
	"time"

	"timer/internal/theme"
)

// Mode identifies which kind of timer a Model represents.
type Mode string

const (
	ModeStopwatch Mode = "stopwatch"
	ModeCountdown Mode = "countdown"
	ModePomodoro  Mode = "pomodoro"
)

// Model holds all state for a single timer session. Its methods return
// new Model values rather than mutating the receiver, matching Bubble
// Tea's value-semantics Update loop.
type Model struct {
	Title     string
	ModeVal   Mode
	Theme     theme.Theme
	Target    time.Duration // 0 means no target (stopwatch)
	StartedAt time.Time

	accumulated  time.Duration
	segmentStart time.Time
	running      bool
	done         bool

	Width  int
	Height int

	// Now returns the current time. Overridable in tests for
	// deterministic, sleep-free timing assertions.
	Now func() time.Time
}

// NewModel creates a running Model using the real system clock.
func NewModel(title string, mode Mode, target time.Duration, th theme.Theme) Model {
	return NewModelWithClock(title, mode, target, th, time.Now)
}

// NewModelWithClock creates a running Model using the given clock function.
func NewModelWithClock(title string, mode Mode, target time.Duration, th theme.Theme, now func() time.Time) Model {
	start := now()
	return Model{
		Title:        title,
		ModeVal:      mode,
		Theme:        th,
		Target:       target,
		StartedAt:    start,
		segmentStart: start,
		running:      true,
		Now:          now,
	}
}

func (m Model) clock() time.Time {
	if m.Now != nil {
		return m.Now()
	}
	return time.Now()
}

// Elapsed returns the total time counted so far, across all run segments.
func (m Model) Elapsed() time.Duration {
	if m.running {
		return m.accumulated + m.clock().Sub(m.segmentStart)
	}
	return m.accumulated
}

// Remaining returns the time left until Target, clamped to zero. Only
// meaningful when Target > 0.
func (m Model) Remaining() time.Duration {
	r := m.Target - m.Elapsed()
	if r < 0 {
		return 0
	}
	return r
}

// Running reports whether the timer is actively counting.
func (m Model) Running() bool { return m.running }

// Done reports whether the timer has been finalized (stopped, or a
// countdown that reached zero).
func (m Model) Done() bool { return m.done }

// IsCountdownFinished reports whether a countdown's target has been
// reached. Always false for stopwatches (Target == 0).
func (m Model) IsCountdownFinished() bool {
	return m.Target > 0 && m.Remaining() <= 0
}

// Pause freezes the elapsed time. No-op if already paused or done.
func (m Model) Pause() Model {
	if !m.running || m.done {
		return m
	}
	m.accumulated = m.Elapsed()
	m.running = false
	return m
}

// Resume continues counting from the accumulated elapsed time. No-op if
// already running or done.
func (m Model) Resume() Model {
	if m.running || m.done {
		return m
	}
	m.segmentStart = m.clock()
	m.running = true
	return m
}

// Restart zeroes the accumulated time and starts a fresh run segment.
func (m Model) Restart() Model {
	m.accumulated = 0
	m.segmentStart = m.clock()
	m.running = true
	m.done = false
	return m
}

// MarkDone finalizes the elapsed time and stops the timer.
func (m Model) MarkDone() Model {
	m.accumulated = m.Elapsed()
	m.running = false
	m.done = true
	return m
}
