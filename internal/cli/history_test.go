package cli

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"io"
	"os"
	"testing"
	"time"

	"timer/internal/store"
)

func sampleRecords() []store.HistoryRecord {
	now := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)
	return []store.HistoryRecord{
		{ID: "1", Title: "Deep work", Mode: "stopwatch", Elapsed: 10 * time.Minute, StartedAt: now, EndedAt: now.Add(10 * time.Minute), Completed: true},
		{ID: "2", Title: "Deep work", Mode: "stopwatch", Elapsed: 20 * time.Minute, StartedAt: now, EndedAt: now.Add(20 * time.Minute), Completed: true},
		{ID: "3", Title: "Tea", Mode: "countdown", Target: 5 * time.Minute, Elapsed: 5 * time.Minute, StartedAt: now, EndedAt: now.Add(5 * time.Minute), Completed: true},
	}
}

func TestAggregateStats(t *testing.T) {
	stats := aggregateStats(sampleRecords())

	dw, ok := stats["Deep work"]
	if !ok || dw.Count != 2 || dw.Total != 30*time.Minute {
		t.Errorf("Deep work stats = %+v, want count=2 total=30m", dw)
	}

	tea, ok := stats["Tea"]
	if !ok || tea.Count != 1 || tea.Total != 5*time.Minute {
		t.Errorf("Tea stats = %+v, want count=1 total=5m", tea)
	}
}

func captureStdout(t *testing.T, fn func() error) (string, error) {
	t.Helper()
	old := os.Stdout
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("os.Pipe: %v", err)
	}
	os.Stdout = w

	fnErr := fn()

	w.Close()
	os.Stdout = old

	out, _ := io.ReadAll(r)
	return string(out), fnErr
}

func TestExportHistoryJSON(t *testing.T) {
	out, err := captureStdout(t, func() error {
		return exportHistoryJSON(sampleRecords())
	})
	if err != nil {
		t.Fatalf("exportHistoryJSON() error: %v", err)
	}

	var decoded []store.HistoryRecord
	if err := json.Unmarshal([]byte(out), &decoded); err != nil {
		t.Fatalf("output is not valid JSON: %v\n%s", err, out)
	}
	if len(decoded) != 3 {
		t.Errorf("decoded %d records, want 3", len(decoded))
	}
}

func TestExportHistoryCSV(t *testing.T) {
	out, err := captureStdout(t, func() error {
		return exportHistoryCSV(sampleRecords())
	})
	if err != nil {
		t.Fatalf("exportHistoryCSV() error: %v", err)
	}

	reader := csv.NewReader(bytes.NewReader([]byte(out)))
	rows, err := reader.ReadAll()
	if err != nil {
		t.Fatalf("output is not valid CSV: %v\n%s", err, out)
	}
	if len(rows) != 4 { // header + 3 records
		t.Fatalf("got %d rows, want 4 (1 header + 3 records)", len(rows))
	}
	if rows[0][1] != "title" {
		t.Errorf("header row = %v, want title column at index 1", rows[0])
	}
}

func TestPrintHistoryTableEmptyNoPanic(t *testing.T) {
	_, err := captureStdout(t, func() error {
		printHistoryTable(nil)
		return nil
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestPrintHistoryTableWithRecordsNoPanic(t *testing.T) {
	out, err := captureStdout(t, func() error {
		printHistoryTable(sampleRecords())
		return nil
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out == "" {
		t.Error("expected non-empty table output")
	}
}
