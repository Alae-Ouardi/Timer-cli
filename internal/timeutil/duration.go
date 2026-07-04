// Package timeutil parses user-supplied duration strings and formats
// durations for display as hh:mm:ss.
package timeutil

import (
	"fmt"
	"strconv"
	"time"
)

// ParseDuration parses a duration string. It accepts Go duration syntax
// ("25m", "1h", "1h30m", "90s") as well as a bare integer, which is
// interpreted as a number of minutes ("90" -> 90m). The result must be
// strictly positive.
func ParseDuration(s string) (time.Duration, error) {
	if s == "" {
		return 0, fmt.Errorf("duration is empty")
	}

	var d time.Duration
	if n, err := strconv.Atoi(s); err == nil {
		d = time.Duration(n) * time.Minute
	} else {
		parsed, err := time.ParseDuration(s)
		if err != nil {
			return 0, fmt.Errorf("invalid duration %q: %w", s, err)
		}
		d = parsed
	}

	if d <= 0 {
		return 0, fmt.Errorf("duration must be positive, got %q", s)
	}
	return d, nil
}

// FormatHMS formats a duration as zero-padded hh:mm:ss. Durations longer
// than 24h roll into the hours field (e.g. "25:00:00"). Negative durations
// are clamped to zero.
func FormatHMS(d time.Duration) string {
	if d < 0 {
		d = 0
	}
	total := int64(d.Seconds())
	h := total / 3600
	m := (total % 3600) / 60
	s := total % 60
	return fmt.Sprintf("%02d:%02d:%02d", h, m, s)
}
