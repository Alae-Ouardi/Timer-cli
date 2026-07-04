package cli

import (
	"fmt"

	"github.com/spf13/cobra"

	"timer/internal/timeutil"
	"timer/internal/tui"
)

var (
	pomodoroWork   string
	pomodoroBreak  string
	pomodoroRounds int
)

var pomodoroCmd = &cobra.Command{
	Use:     "pomodoro [title]",
	Aliases: []string{"p"},
	Short:   "Run a pomodoro session: alternating work/break countdowns with a round counter",
	Args:    cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		work, err := timeutil.ParseDuration(pomodoroWork)
		if err != nil {
			return fmt.Errorf("invalid --work: %w", err)
		}
		brk, err := timeutil.ParseDuration(pomodoroBreak)
		if err != nil {
			return fmt.Errorf("invalid --break: %w", err)
		}
		if pomodoroRounds < 1 {
			return fmt.Errorf("--rounds must be at least 1, got %d", pomodoroRounds)
		}

		title := ""
		if len(args) > 0 {
			title = args[0]
		}

		m := tui.NewPomodoroModel(title, work, brk, pomodoroRounds, ResolveTheme())
		return runProgram(m)
	},
}

func init() {
	pomodoroCmd.Flags().StringVar(&pomodoroWork, "work", "25m", "work phase duration")
	pomodoroCmd.Flags().StringVar(&pomodoroBreak, "break", "5m", "break phase duration")
	pomodoroCmd.Flags().IntVar(&pomodoroRounds, "rounds", 4, "number of work/break rounds")
	rootCmd.AddCommand(pomodoroCmd)
}
