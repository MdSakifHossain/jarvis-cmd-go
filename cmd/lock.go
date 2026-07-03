package cmd

import (
	"fmt"
	"jarvis/banner"
	"jarvis/console"
	"jarvis/support"
	"os/exec"
)

func Lock() {
	support.ShowBanner(banner.Lock)
	console.Info("Initializing command...")
	console.Info("Locking screen...")
	lockScreen()
	console.Info("Command finished successfully.")
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
		console.Fail(fmt.Sprintf("Failed to lock screen: %v", err))
	}
}
