// Package tui implements the Bubble Tea model, update loop, and view for
// the timer's terminal UI.
package tui

import "strings"

const glyphHeight = 7

// Row fragments used to build each digit's seven-segment shape:
// a horizontal bar (top/middle/bottom), a blank row, and the three states
// of the paired vertical segments (both lit, left only, right only).
const (
	bar   = " █████ "
	blank = "       "
	both  = "█     █"
	left  = "█      "
	right = "      █"
)

// glyphs maps each supported rune to its seven-row, seven-segment-style
// shape: a top bar, two rows of upper verticals, a middle bar, two rows
// of lower verticals, and a bottom bar.
var glyphs = map[rune][]string{
	'0': {bar, both, both, blank, both, both, bar},
	'1': {blank, right, right, blank, right, right, blank},
	'2': {bar, right, right, bar, left, left, bar},
	'3': {bar, right, right, bar, right, right, bar},
	'4': {blank, both, both, bar, right, right, blank},
	'5': {bar, left, left, bar, right, right, bar},
	'6': {bar, left, left, bar, both, both, bar},
	'7': {bar, right, right, blank, right, right, blank},
	'8': {bar, both, both, bar, both, both, bar},
	'9': {bar, both, both, bar, right, right, bar},
	':': {"   ", " █ ", " █ ", "   ", " █ ", " █ ", "   "},
}

var blankGlyph = []string{blank, blank, blank, blank, blank, blank, blank}

// RenderDigits renders s (expected to contain only 0-9 and ':') as a
// multi-line string of large seven-segment-style glyphs, joined
// horizontally with a single-space gap between characters. Unknown runes
// render as blank space of consistent height.
func RenderDigits(s string) string {
	rows := make([]strings.Builder, glyphHeight)

	runes := []rune(s)
	for i, r := range runes {
		glyph, ok := glyphs[r]
		if !ok {
			glyph = blankGlyph
		}
		for row := 0; row < glyphHeight; row++ {
			rows[row].WriteString(glyph[row])
			if i != len(runes)-1 {
				rows[row].WriteString(" ")
			}
		}
	}

	lines := make([]string, glyphHeight)
	for i := range rows {
		lines[i] = rows[i].String()
	}
	return strings.Join(lines, "\n")
}
