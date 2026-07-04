package store

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func newRecord(title, mode string, elapsed time.Duration, when time.Time) HistoryRecord {
	return HistoryRecord{
		Title:     title,
		Mode:      mode,
		Elapsed:   elapsed,
		StartedAt: when,
		EndedAt:   when.Add(elapsed),
		Completed: true,
	}
}

func TestAppendThenAll(t *testing.T) {
	t.Setenv("XDG_CONFIG_HOME", t.TempDir())

	now := time.Now()
	if err := AppendHistory(newRecord("Deep work", "stopwatch", 10*time.Minute, now)); err != nil {
		t.Fatalf("AppendHistory() unexpected error: %v", err)
	}
	if err := AppendHistory(newRecord("Tea", "countdown", 5*time.Minute, now)); err != nil {
		t.Fatalf("AppendHistory() unexpected error: %v", err)
	}

	all, err := AllHistory()
	if err != nil {
		t.Fatalf("AllHistory() unexpected error: %v", err)
	}
	if len(all) != 2 {
		t.Fatalf("AllHistory() returned %d records, want 2", len(all))
	}
	if all[0].ID == "" || all[1].ID == "" {
		t.Errorf("expected non-empty IDs, got %q and %q", all[0].ID, all[1].ID)
	}
}

func TestAllHistoryEmptyWhenMissing(t *testing.T) {
	t.Setenv("XDG_CONFIG_HOME", t.TempDir())

	all, err := AllHistory()
	if err != nil {
		t.Fatalf("AllHistory() unexpected error: %v", err)
	}
	if len(all) != 0 {
		t.Fatalf("AllHistory() returned %d records, want 0", len(all))
	}
}

func TestSearchFiltersByTitleCaseInsensitive(t *testing.T) {
	t.Setenv("XDG_CONFIG_HOME", t.TempDir())

	now := time.Now()
	_ = AppendHistory(newRecord("Deep work", "stopwatch", time.Minute, now))
	_ = AppendHistory(newRecord("Deep focus", "stopwatch", time.Minute, now))
	_ = AppendHistory(newRecord("Tea break", "countdown", time.Minute, now))

	results, err := SearchHistory("deep")
	if err != nil {
		t.Fatalf("SearchHistory() unexpected error: %v", err)
	}
	if len(results) != 2 {
		t.Fatalf("SearchHistory(\"deep\") returned %d records, want 2", len(results))
	}
}

func TestStatsByTitle(t *testing.T) {
	t.Setenv("XDG_CONFIG_HOME", t.TempDir())

	now := time.Now()
	_ = AppendHistory(newRecord("Deep work", "stopwatch", 10*time.Minute, now))
	_ = AppendHistory(newRecord("Deep work", "stopwatch", 20*time.Minute, now))
	_ = AppendHistory(newRecord("Tea break", "countdown", 5*time.Minute, now))

	stats, err := StatsByTitle()
	if err != nil {
		t.Fatalf("StatsByTitle() unexpected error: %v", err)
	}

	dw, ok := stats["Deep work"]
	if !ok {
		t.Fatalf("expected stats for %q", "Deep work")
	}
	if dw.Count != 2 {
		t.Errorf("Deep work count = %d, want 2", dw.Count)
	}
	if dw.Total != 30*time.Minute {
		t.Errorf("Deep work total = %v, want %v", dw.Total, 30*time.Minute)
	}

	tb, ok := stats["Tea break"]
	if !ok || tb.Count != 1 || tb.Total != 5*time.Minute {
		t.Errorf("Tea break stats = %+v, want count=1 total=5m", tb)
	}
}

func TestAllHistoryCorruptFileReturnsError(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("XDG_CONFIG_HOME", dir)

	cfgDir := filepath.Join(dir, "timer-cli")
	if err := os.MkdirAll(cfgDir, 0o755); err != nil {
		t.Fatalf("MkdirAll: %v", err)
	}
	if err := os.WriteFile(filepath.Join(cfgDir, "history.json"), []byte("not json"), 0o644); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}

	if _, err := AllHistory(); err == nil {
		t.Error("AllHistory() expected error for corrupt file, got nil")
	}
}
