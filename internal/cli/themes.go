package cli

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"

	"timer/internal/store"
	"timer/internal/theme"
)

var themesCmd = &cobra.Command{
	Use:     "themes",
	Aliases: []string{"t"},
	Short:   "List available color themes",
	RunE: func(cmd *cobra.Command, args []string) error {
		listThemes()
		return nil
	},
}

var themesSetCmd = &cobra.Command{
	Use:   "set <name>",
	Short: "Set the default theme used when --theme is not passed",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]
		if !isValidTheme(name) {
			return fmt.Errorf("unknown theme %q (run 'timer themes' to see available themes)", name)
		}

		cfg, err := store.LoadConfig()
		if err != nil {
			return err
		}
		cfg.DefaultTheme = name
		if err := store.SaveConfig(cfg); err != nil {
			return err
		}
		fmt.Printf("Default theme set to %q.\n", name)
		return nil
	},
}

func init() {
	themesCmd.AddCommand(themesSetCmd)
	rootCmd.AddCommand(themesCmd)
}

func isValidTheme(name string) bool {
	for _, th := range theme.List() {
		if th.Name == name {
			return true
		}
	}
	return false
}

func listThemes() {
	cfg, _ := store.LoadConfig()
	for _, th := range theme.List() {
		swatch := lipgloss.NewStyle().
			Background(lipgloss.Color(th.Background)).
			Foreground(lipgloss.Color(th.DigitColor)).
			Render(" ●●●● ")
		marker := "  "
		if th.Name == cfg.DefaultTheme {
			marker = "* "
		}
		fmt.Printf("%s%s %s\n", marker, swatch, th.Name)
	}
}
