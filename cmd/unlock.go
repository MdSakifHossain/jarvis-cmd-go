package cmd

import (
	"fmt"
	"jarvis/banner"
	"jarvis/console"
	"jarvis/support"
	"os/exec"
)

func Unlock() {
	showUnlockHeader()

	console.Info("Initializing command...")
	console.Info("Unlocking screen...")

	unlockScreen()

	console.Info("Command finished successfully.")
}

func showUnlockHeader() {
	support.ShowHeader(banner.Unlock)
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
		console.Fail(fmt.Sprintf("Failed to unlock screen: %v", err))
	}
}
