package cmd

import (
	"fmt"
	"jarvis/output"
	"jarvis/support"
	"os/exec"
)

func Lock() {
	output.Info("Initializing command...")
	output.Info("Locking screen...")
	lockScreen()
	output.Info("Command finished successfully.")
}

func lockScreen() {
	support.RequireDBus()

	cmd := exec.Command(
		"dbus-send",
		"--session",
		"--dest=org.gnome.ScreenSaver",
		"--type=method_call",
		"--print-reply",
		"/org/gnome/ScreenSaver",
		"org.gnome.ScreenSaver.Lock",
	)

	if err := cmd.Run(); err != nil {
		output.Fail(fmt.Sprintf("Failed to lock screen: %v", err))
	}
}
