package cmd

import (
	"fmt"
	"os/exec"
)

func requireCommand(cmd, hint string) {
	if _, err := exec.LookPath(cmd); err != nil {
		fail(fmt.Sprintf(
			"Missing dependency: %s\n\n%s",
			cmd,
			hint,
		))
	}
}

func requireOpenRGB() {
	requireCommand(
		"openrgb",
		`Install OpenRGB from:

https://openrgb.org/releases.html

Recommended package:
Linux amd64 (Debian Bookworm .deb)`,
	)
}

// func requireDBus() {
// 	requireCommand(
// 		"dbus-send",
// 		`Install it with:

// sudo apt install dbus-x11`,
// 	)
// }
