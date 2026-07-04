// Package theme defines the color palettes used to render the timer TUI.
package theme

import "sort"

// Theme is a named color palette for the TUI. Colors are hex strings
// consumable directly by lipgloss.Color.
type Theme struct {
	Name          string
	DigitColor    string
	TitleColor    string
	HintColor     string
	ProgressColor string
	AccentColor   string
	DoneColor     string
	Background    string
}

const defaultThemeName = "obsidian"

var registry = map[string]Theme{
	"obsidian": {
		Name: "obsidian", DigitColor: "#C9D9E0", TitleColor: "#8FAEC4",
		HintColor: "#4A5C6E", ProgressColor: "#5EC2E0", AccentColor: "#5EC2E0",
		DoneColor: "#7FE0B4", Background: "#0E1725",
	},
	"slate": {
		Name: "slate", DigitColor: "#EAEFEF", TitleColor: "#B8C4D9",
		HintColor: "#5A6B85", ProgressColor: "#7EC8E3", AccentColor: "#7EC8E3",
		DoneColor: "#8FE3B0", Background: "#283048",
	},
	"abyss": {
		Name: "abyss", DigitColor: "#E8F0F5", TitleColor: "#9FB8CC",
		HintColor: "#3E5670", ProgressColor: "#4FD1C5", AccentColor: "#4FD1C5",
		DoneColor: "#7FE0B4", Background: "#0A2947",
	},
	"dracula": {
		Name: "dracula", DigitColor: "#F8F8F2", TitleColor: "#BD93F9",
		HintColor: "#6272A4", ProgressColor: "#BD93F9", AccentColor: "#50FA7B",
		DoneColor: "#FF79C6", Background: "#282A36",
	},
	"midnight": {
		Name: "midnight", DigitColor: "#00E5FF", TitleColor: "#7DD3FC",
		HintColor: "#6B7280", ProgressColor: "#00E5FF", AccentColor: "#22D3EE",
		DoneColor: "#34D399", Background: "#000000",
	},
	"matrix": {
		Name: "matrix", DigitColor: "#39FF14", TitleColor: "#4ADE80",
		HintColor: "#4B5563", ProgressColor: "#39FF14", AccentColor: "#22C55E",
		DoneColor: "#39FF14", Background: "#000000",
	},
	"ember": {
		Name: "ember", DigitColor: "#FF6B35", TitleColor: "#FCA5A5",
		HintColor: "#6B7280", ProgressColor: "#FF6B35", AccentColor: "#F97316",
		DoneColor: "#FDE047", Background: "#000000",
	},
	"mono": {
		Name: "mono", DigitColor: "#FFFFFF", TitleColor: "#D1D5DB",
		HintColor: "#6B7280", ProgressColor: "#FFFFFF", AccentColor: "#9CA3AF",
		DoneColor: "#FFFFFF", Background: "#000000",
	},
	"synthwave": {
		Name: "synthwave", DigitColor: "#FF00FF", TitleColor: "#C084FC",
		HintColor: "#6B7280", ProgressColor: "#FF00FF", AccentColor: "#F472B6",
		DoneColor: "#00FFFF", Background: "#000000",
	},
	"paper": {
		Name: "paper", DigitColor: "#111827", TitleColor: "#374151",
		HintColor: "#9CA3AF", ProgressColor: "#111827", AccentColor: "#4B5563",
		DoneColor: "#059669", Background: "#FFFFFF",
	},
	"daylight": {
		Name: "daylight", DigitColor: "#1D4ED8", TitleColor: "#2563EB",
		HintColor: "#9CA3AF", ProgressColor: "#1D4ED8", AccentColor: "#3B82F6",
		DoneColor: "#059669", Background: "#FFFFFF",
	},
}

// Get returns the theme with the given name, falling back to the default
// theme ("obsidian") if the name is unknown.
func Get(name string) Theme {
	if t, ok := registry[name]; ok {
		return t
	}
	return registry[defaultThemeName]
}

// Default returns the default theme.
func Default() Theme {
	return registry[defaultThemeName]
}

// List returns all themes sorted by name.
func List() []Theme {
	names := make([]string, 0, len(registry))
	for name := range registry {
		names = append(names, name)
	}
	sort.Strings(names)

	themes := make([]Theme, 0, len(names))
	for _, name := range names {
		themes = append(themes, registry[name])
	}
	return themes
}
