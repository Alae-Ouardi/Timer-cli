package cli

import "testing"

func TestResolveThemeFallsBackToDefault(t *testing.T) {
	themeFlag = "does-not-exist"
	th := ResolveTheme()
	if th.Name != "obsidian" {
		t.Errorf("ResolveTheme() = %q, want %q", th.Name, "obsidian")
	}
}

func TestResolveThemeHonorsFlag(t *testing.T) {
	themeFlag = "matrix"
	th := ResolveTheme()
	if th.Name != "matrix" {
		t.Errorf("ResolveTheme() = %q, want %q", th.Name, "matrix")
	}
}
