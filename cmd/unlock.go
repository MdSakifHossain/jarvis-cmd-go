package cmd

import (
	"fmt"
	"jarvis/output"
	"jarvis/support"
	"os/exec"
)

func Unlock() {
	output.Info("Initializing command...")
	output.Info("Unlocking screen...")
	unlockScreen()
	output.Info("Command finished successfully.")
}

func unlockScreen() {
	support.RequireDBus()

	cmd := exec.Command(
		"dbus-send",
		"--session",
		"--dest=org.gnome.ScreenSaver",
		"--type=method_call",
		"--print-reply",
		"/org/gnome/ScreenSaver",
		"org.gnome.ScreenSaver.SetActive",
		"boolean:false",
	)

	if err := cmd.Run(); err != nil {
		output.Fail(fmt.Sprintf("Failed to unlock screen: %v", err))
	}
}
