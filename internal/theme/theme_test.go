package theme

import "testing"

func TestGetKnownThemes(t *testing.T) {
	names := []string{
		"obsidian", "slate", "abyss", "dracula",
		"midnight", "matrix", "ember", "mono", "synthwave", "paper", "daylight",
	}
	for _, name := range names {
		th := Get(name)
		if th.Name != name {
			t.Errorf("Get(%q).Name = %q, want %q", name, th.Name, name)
		}
	}
}

func TestGetUnknownFallsBackToDefault(t *testing.T) {
	th := Get("does-not-exist")
	if th.Name != "obsidian" {
		t.Errorf("Get(unknown).Name = %q, want %q", th.Name, "obsidian")
	}
}

func TestList(t *testing.T) {
	themes := List()
	if len(themes) != 11 {
		t.Fatalf("List() returned %d themes, want 11", len(themes))
	}
	for i := 1; i < len(themes); i++ {
		if themes[i-1].Name > themes[i].Name {
			t.Errorf("List() not sorted: %q before %q", themes[i-1].Name, themes[i].Name)
		}
	}
}

func TestDefault(t *testing.T) {
	if Default().Name != "obsidian" {
		t.Errorf("Default().Name = %q, want %q", Default().Name, "obsidian")
	}
}

func TestAllThemesHaveNonEmptyColors(t *testing.T) {
	for _, th := range List() {
		if th.DigitColor == "" || th.TitleColor == "" || th.HintColor == "" ||
			th.ProgressColor == "" || th.AccentColor == "" || th.DoneColor == "" || th.Background == "" {
			t.Errorf("theme %q has an empty color field: %+v", th.Name, th)
		}
	}
}
