package cli

import (
	"github.com/spf13/cobra"

	"timer/internal/timeutil"
	"timer/internal/tui"
)

var countdownCmd = &cobra.Command{
	Use:     "countdown <duration> [title]",
	Aliases: []string{"c"},
	Short:   "Start a countdown for the given duration (e.g. 25m, 1h, 1h30m, 90s, or a bare number of minutes)",
	Args:    cobra.RangeArgs(1, 2),
	RunE: func(cmd *cobra.Command, args []string) error {
		target, err := timeutil.ParseDuration(args[0])
		if err != nil {
			return err
		}

		title := ""
		if len(args) > 1 {
			title = args[1]
		}

		m := tui.NewModel(title, tui.ModeCountdown, target, ResolveTheme())
		return runProgram(m)
	},
}

func init() {
	rootCmd.AddCommand(countdownCmd)
}
