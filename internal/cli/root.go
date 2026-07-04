// Package cli defines the timer command-line interface: the Cobra root
// command, its subcommands, and the glue that launches the Bubble Tea TUI.
package cli

import (
	"github.com/spf13/cobra"

	"timer/internal/store"
	"timer/internal/theme"
)

var themeFlag string

var rootCmd = &cobra.Command{
	Use:           "timer",
	Short:         "A beautiful terminal timer: stopwatch, countdown, and pomodoro.",
	SilenceUsage:  true,
	SilenceErrors: true,
}

// Execute runs the root command; main.go calls this and exits non-zero on
// error.
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	cfg, _ := store.LoadConfig()
	rootCmd.PersistentFlags().StringVar(&themeFlag, "theme", cfg.DefaultTheme, "color theme to use (see: timer themes)")
}

// ResolveTheme returns the theme selected via --theme, falling back to
// the registry default if the name is unknown.
func ResolveTheme() theme.Theme {
	return theme.Get(themeFlag)
}
