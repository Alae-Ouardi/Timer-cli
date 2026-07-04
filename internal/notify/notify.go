// Package notify sends best-effort desktop notifications and terminal
// bells when a timer needs the user's attention. Failures are swallowed:
// a missing notification daemon must never crash the TUI.
package notify

import (
	_ "embed"
	"fmt"
	"os"

	"github.com/gen2brain/beeep"
)

//go:embed icon.png
var icon []byte

// Notify sends a desktop notification (macOS/Linux) with the timer's icon
// and rings the terminal bell. Errors from either are ignored.
//
// On Linux the icon renders via D-Bus or notify-send/kdialog. On macOS it
// requires terminal-notifier (`brew install terminal-notifier`); without
// it, notifications fall back to a plain osascript alert using the
// system's default icon.
func Notify(title, message string) {
	_ = beeep.Notify(title, message, icon)
	fmt.Fprint(os.Stdout, "\a")
}
