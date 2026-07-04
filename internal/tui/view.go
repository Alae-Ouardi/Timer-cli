package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/lipgloss"

	"timer/internal/timeutil"
)

// edgeBandHeight is the number of rows reserved for the title band at the
// top and the hints band at the bottom. Centering content within a
// 3-row band leaves one blank row of breathing room on either side, so
// the title/hints sit near the edge without touching it.
const edgeBandHeight = 3

// View renders the current frame over the theme's background. While
// running, the title is pinned near the top of the screen, the hints are
// pinned near the bottom, and the big digits (plus a progress bar for
// countdowns) are vertically centered in the space between — all
// re-laid-out automatically as the terminal is resized. The done screen
// is rendered as a single centered block.
func (m Model) View() string {
	bg := lipgloss.Color(m.Theme.Background)

	if m.done {
		body := m.renderDone()
		if m.Width <= 0 || m.Height <= 0 {
			return lipgloss.NewStyle().Background(bg).Render(body)
		}
		return lipgloss.Place(m.Width, m.Height, lipgloss.Center, lipgloss.Center, body,
			lipgloss.WithWhitespaceBackground(bg))
	}

	if m.Width <= 0 || m.Height <= 0 {
		sections := []string{}
		if title := m.renderTitle(); title != "" {
			sections = append(sections, title)
		}
		sections = append(sections, m.renderMiddle(), m.renderHints())
		return lipgloss.NewStyle().Background(bg).Render(lipgloss.JoinVertical(lipgloss.Center, sections...))
	}

	middleHeight := m.Height - 2*edgeBandHeight
	if middleHeight < 1 {
		middleHeight = 1
	}

	top := lipgloss.Place(m.Width, edgeBandHeight, lipgloss.Center, lipgloss.Center, m.renderTitle(),
		lipgloss.WithWhitespaceBackground(bg))
	middle := lipgloss.Place(m.Width, middleHeight, lipgloss.Center, lipgloss.Center, m.renderMiddle(),
		lipgloss.WithWhitespaceBackground(bg))
	bottom := lipgloss.Place(m.Width, edgeBandHeight, lipgloss.Center, lipgloss.Center, m.renderHints(),
		lipgloss.WithWhitespaceBackground(bg))

	return lipgloss.JoinVertical(lipgloss.Left, top, middle, bottom)
}

func (m Model) renderTitle() string {
	if m.Title == "" {
		return ""
	}
	bg := lipgloss.Color(m.Theme.Background)
	titleStyle := lipgloss.NewStyle().Background(bg).Foreground(lipgloss.Color(m.Theme.TitleColor)).Bold(true)
	return titleStyle.Render(m.Title)
}

func (m Model) renderMiddle() string {
	bg := lipgloss.Color(m.Theme.Background)
	digitStyle := lipgloss.NewStyle().Background(bg).Foreground(lipgloss.Color(m.Theme.DigitColor)).Bold(true)

	var clock string
	if m.ModeVal == ModeCountdown {
		clock = timeutil.FormatHMS(m.Remaining())
	} else {
		clock = timeutil.FormatHMS(m.Elapsed())
	}

	sections := []string{digitStyle.Render(RenderDigits(clock))}
	if m.ModeVal == ModeCountdown && m.Target > 0 {
		sections = append(sections, m.renderProgress())
	}
	return lipgloss.JoinVertical(lipgloss.Center, sections...)
}

func (m Model) renderProgress() string {
	pct := float64(m.Elapsed()) / float64(m.Target)
	if pct < 0 {
		pct = 0
	}
	if pct > 1 {
		pct = 1
	}

	bar := progress.New(
		progress.WithSolidFill(m.Theme.ProgressColor),
		progress.WithWidth(30),
	)
	return lipgloss.NewStyle().Background(lipgloss.Color(m.Theme.Background)).Render(bar.ViewAs(pct))
}

func (m Model) renderHints() string {
	bg := lipgloss.Color(m.Theme.Background)
	activeStyle := lipgloss.NewStyle().Background(bg).Foreground(lipgloss.Color(m.Theme.AccentColor)).Bold(true)
	dimStyle := lipgloss.NewStyle().Background(bg).Foreground(lipgloss.Color(m.Theme.HintColor))
	sepStyle := lipgloss.NewStyle().Background(bg)

	pauseLabel, continueLabel := "p=pause", "c=continue"
	if m.running {
		pauseLabel = activeStyle.Render(pauseLabel)
		continueLabel = dimStyle.Render(continueLabel)
	} else {
		pauseLabel = dimStyle.Render(pauseLabel)
		continueLabel = activeStyle.Render(continueLabel)
	}

	stopLabel := activeStyle.Render("s=stop")
	restartLabel := activeStyle.Render("r=restart")

	return strings.Join([]string{stopLabel, pauseLabel, continueLabel, restartLabel}, sepStyle.Render("  "))
}

func (m Model) renderDone() string {
	bg := lipgloss.Color(m.Theme.Background)
	doneStyle := lipgloss.NewStyle().Background(bg).Foreground(lipgloss.Color(m.Theme.DoneColor)).Bold(true)
	hintStyle := lipgloss.NewStyle().Background(bg).Foreground(lipgloss.Color(m.Theme.HintColor))
	plainStyle := lipgloss.NewStyle().Background(bg)

	title := m.Title
	if title == "" {
		title = string(m.ModeVal)
	}

	elapsed := timeutil.FormatHMS(m.Elapsed())
	msg := fmt.Sprintf("%s\n\n%s\n\n%s",
		doneStyle.Render("DONE"),
		plainStyle.Render(fmt.Sprintf("%s — %s", title, elapsed)),
		hintStyle.Render("press any key to continue"),
	)
	return msg
}
