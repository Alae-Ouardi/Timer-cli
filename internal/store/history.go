package store

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// HistoryRecord captures a single completed or stopped timer session.
type HistoryRecord struct {
	ID        string        `json:"id"`
	Title     string        `json:"title"`
	Mode      string        `json:"mode"` // stopwatch | countdown | pomodoro
	Target    time.Duration `json:"target"`
	Elapsed   time.Duration `json:"elapsed"`
	StartedAt time.Time     `json:"started_at"`
	EndedAt   time.Time     `json:"ended_at"`
	Completed bool          `json:"completed"`
}

// TitleStats aggregates history records sharing the same title.
type TitleStats struct {
	Count int
	Total time.Duration
}

func historyPath() (string, error) {
	dir, err := configDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "history.json"), nil
}

func newID() string {
	buf := make([]byte, 8)
	_, _ = rand.Read(buf)
	return fmt.Sprintf("%x", buf)
}

// AllHistory returns every stored record, or an empty slice if none exist
// yet.
func AllHistory() ([]HistoryRecord, error) {
	path, err := historyPath()
	if err != nil {
		return nil, fmt.Errorf("resolve history path: %w", err)
	}

	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return []HistoryRecord{}, nil
	}
	if err != nil {
		return nil, fmt.Errorf("read history: %w", err)
	}

	var records []HistoryRecord
	if err := json.Unmarshal(data, &records); err != nil {
		return nil, fmt.Errorf("parse history: %w", err)
	}
	return records, nil
}

// AppendHistory adds a record to the history file, assigning it an ID if
// it doesn't already have one.
func AppendHistory(rec HistoryRecord) error {
	if rec.ID == "" {
		rec.ID = newID()
	}

	records, err := AllHistory()
	if err != nil {
		return err
	}
	records = append(records, rec)

	dir, err := ensureConfigDir()
	if err != nil {
		return fmt.Errorf("append history: %w", err)
	}

	data, err := json.MarshalIndent(records, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal history: %w", err)
	}

	if err := os.WriteFile(filepath.Join(dir, "history.json"), data, 0o644); err != nil {
		return fmt.Errorf("write history: %w", err)
	}
	return nil
}

// SearchHistory returns records whose title contains substr, case-insensitive.
func SearchHistory(substr string) ([]HistoryRecord, error) {
	all, err := AllHistory()
	if err != nil {
		return nil, err
	}

	needle := strings.ToLower(substr)
	var matched []HistoryRecord
	for _, r := range all {
		if strings.Contains(strings.ToLower(r.Title), needle) {
			matched = append(matched, r)
		}
	}
	return matched, nil
}

// StatsByTitle aggregates count and total elapsed time per title.
func StatsByTitle() (map[string]TitleStats, error) {
	all, err := AllHistory()
	if err != nil {
		return nil, err
	}

	stats := make(map[string]TitleStats)
	for _, r := range all {
		s := stats[r.Title]
		s.Count++
		s.Total += r.Elapsed
		stats[r.Title] = s
	}
	return stats, nil
}
