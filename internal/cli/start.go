package cli

import (
	"github.com/spf13/cobra"

	"timer/internal/tui"
)

var startCmd = &cobra.Command{
	Use:     "start [title]",
	Aliases: []string{"s"},
	Short:   "Start a stopwatch, counting up from 00:00:00",
	Args:    cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		title := ""
		if len(args) > 0 {
			title = args[0]
		}

		m := tui.NewModel(title, tui.ModeStopwatch, 0, ResolveTheme())
		return runProgram(m)
	},
}

func init() {
	rootCmd.AddCommand(startCmd)
}
