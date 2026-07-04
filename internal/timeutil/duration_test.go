package timeutil

import (
	"testing"
	"time"
)

func TestParseDuration(t *testing.T) {
	cases := []struct {
		in      string
		want    time.Duration
		wantErr bool
	}{
		{"25m", 25 * time.Minute, false},
		{"1h", time.Hour, false},
		{"1h30m", time.Hour + 30*time.Minute, false},
		{"90s", 90 * time.Second, false},
		{"90", 90 * time.Minute, false},
		{"5", 5 * time.Minute, false},
		{"", 0, true},
		{"abc", 0, true},
		{"-5m", 0, true},
		{"0", 0, true},
		{"0m", 0, true},
	}

	for _, c := range cases {
		got, err := ParseDuration(c.in)
		if c.wantErr {
			if err == nil {
				t.Errorf("ParseDuration(%q): expected error, got %v", c.in, got)
			}
			continue
		}
		if err != nil {
			t.Errorf("ParseDuration(%q): unexpected error: %v", c.in, err)
			continue
		}
		if got != c.want {
			t.Errorf("ParseDuration(%q) = %v, want %v", c.in, got, c.want)
		}
	}
}

func TestFormatHMS(t *testing.T) {
	cases := []struct {
		in   time.Duration
		want string
	}{
		{0, "00:00:00"},
		{5 * time.Second, "00:00:05"},
		{90 * time.Second, "00:01:30"},
		{time.Hour, "01:00:00"},
		{25*time.Hour + time.Minute + time.Second, "25:01:01"},
		{-5 * time.Second, "00:00:00"},
	}

	for _, c := range cases {
		got := FormatHMS(c.in)
		if got != c.want {
			t.Errorf("FormatHMS(%v) = %q, want %q", c.in, got, c.want)
		}
	}
}
