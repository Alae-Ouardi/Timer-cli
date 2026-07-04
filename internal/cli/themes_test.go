package cli

import "testing"

func TestIsValidTheme(t *testing.T) {
	if !isValidTheme("midnight") {
		t.Error("isValidTheme(\"midnight\") = false, want true")
	}
	if isValidTheme("does-not-exist") {
		t.Error("isValidTheme(\"does-not-exist\") = true, want false")
	}
}

func TestListThemesNoPanic(t *testing.T) {
	out, err := captureStdout(t, func() error {
		listThemes()
		return nil
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out == "" {
		t.Error("listThemes() produced no output")
	}
}
