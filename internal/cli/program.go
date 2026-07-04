package cli

import tea "github.com/charmbracelet/bubbletea"

// runProgram launches a Bubble Tea program for the given model in the
// terminal's alternate screen buffer.
func runProgram(m tea.Model) error {
	p := tea.NewProgram(m, tea.WithAltScreen())
	_, err := p.Run()
	return err
}
