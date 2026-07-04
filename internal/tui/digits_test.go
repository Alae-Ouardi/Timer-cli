package tui

import (
	"strings"
	"testing"
)

func TestGlyphsAreRectangularAndConsistentHeight(t *testing.T) {
	for r, glyph := range glyphs {
		if len(glyph) != glyphHeight {
			t.Errorf("glyph %q has %d rows, want %d", r, len(glyph), glyphHeight)
		}
		width := len([]rune(glyph[0]))
		for i, row := range glyph {
			if len([]rune(row)) != width {
				t.Errorf("glyph %q row %d has width %d, want %d (rectangular)", r, i, len([]rune(row)), width)
			}
		}
	}
}

func TestRenderDigitsProducesEqualLengthLines(t *testing.T) {
	out := RenderDigits("12:34")
	lines := strings.Split(out, "\n")
	if len(lines) != glyphHeight {
		t.Fatalf("RenderDigits() produced %d lines, want %d", len(lines), glyphHeight)
	}

	width := len([]rune(lines[0]))
	for i, line := range lines {
		if len([]rune(line)) != width {
			t.Errorf("line %d has width %d, want %d", i, len([]rune(line)), width)
		}
	}

	if !strings.Contains(out, "█") {
		t.Error("RenderDigits() output does not contain any block characters")
	}
}

func TestRenderDigitsUnknownRuneRendersBlank(t *testing.T) {
	out := RenderDigits("1x")
	if out == "" {
		t.Fatal("RenderDigits() returned empty output")
	}
}

func TestZeroHasNoMiddleBar(t *testing.T) {
	if glyphs['0'][3] != blank {
		t.Errorf("'0' middle row = %q, want blank (0 has no middle segment)", glyphs['0'][3])
	}
}

func TestEightHasMiddleBar(t *testing.T) {
	if glyphs['8'][3] != bar {
		t.Errorf("'8' middle row = %q, want a bar (8 lights every segment)", glyphs['8'][3])
	}
}

func TestOneIsOnlyRightVerticals(t *testing.T) {
	for i, row := range glyphs['1'] {
		if row != blank && row != right {
			t.Errorf("'1' row %d = %q, want blank or right-only (1 only lights the right verticals)", i, row)
		}
	}
}

func TestColonDotsAlignWithDigitVerticals(t *testing.T) {
	colon := glyphs[':']
	if len(colon) != glyphHeight {
		t.Fatalf("colon has %d rows, want %d", len(colon), glyphHeight)
	}
	dotRows := map[int]bool{1: true, 2: true, 4: true, 5: true}
	for i, row := range colon {
		hasDot := strings.Contains(row, "█")
		if hasDot != dotRows[i] {
			t.Errorf("colon row %d hasDot=%v, want %v", i, hasDot, dotRows[i])
		}
	}
}
