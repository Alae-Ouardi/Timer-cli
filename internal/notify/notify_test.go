package notify

import (
	"bytes"
	"image/png"
	"testing"
)

func TestNotifyDoesNotPanic(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("Notify() panicked: %v", r)
		}
	}()
	Notify("Timer done", "Tea is ready")
}

func TestEmbeddedIconIsValidPNG(t *testing.T) {
	if len(icon) == 0 {
		t.Fatal("embedded icon is empty")
	}
	if _, err := png.Decode(bytes.NewReader(icon)); err != nil {
		t.Fatalf("embedded icon is not a valid PNG: %v", err)
	}
}
