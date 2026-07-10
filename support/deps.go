package support

import (
	"fmt"
	"jarvis/output"
	"os/exec"
)

func RequireCommand(cmd, hint string) {
	if _, err := exec.LookPath(cmd); err == nil {
		return
	}

	if hint == "" {
		hint = "No installation instructions provided."
	}

	output.Fail(fmt.Sprintf(`Missing dependency: %s\n\n%s`, cmd, hint))
}

func RequireOpenRGB() {
	RequireCommand(
		"openrgb",
		`Install OpenRGB from:

https://openrgb.org/releases.html

Recommended package:
Linux amd64 (Debian Bookworm .deb)`,
	)
}

func RequireDBus() {
	RequireCommand(
		"dbus-send",
		`Install with:

sudo apt install dbus-x11`,
	)
}
