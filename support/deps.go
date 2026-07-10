package support

import (
	"fmt"
	"jarvis/output"
	"os/exec"
)

func requireCommand(cmd, hint string) {
	if _, err := exec.LookPath(cmd); err == nil {
		return
	}

	if hint == "" {
		hint = "No installation instructions provided."
	}

	output.Fail(fmt.Sprintf("Missing dependency: %s\n\n%s", cmd, hint))
}

func RequireOpenRGB() {
	requireCommand(
		"openrgb",
		`Install OpenRGB:

    https://openrgb.org/releases.html

Note:

    Look for -> Linux amd64 (Debian Bookworm .deb)
`,
	)
}

func RequireDBus() {
	requireCommand(
		"dbus-send",
		`Install with:

    sudo apt install dbus-x11
`,
	)
}
