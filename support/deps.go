package support

import (
	"fmt"
	"jarvis/output"
	"os"
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

func RequireTree() {
	requireCommand(
		"tree",
		`Install Tree:

    sudo apt install tree -y
`)
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

func RequireFile(path, description string) {
	if _, err := os.Stat(path); err == nil {
		return
	}

	if description == "" {
		description = "Required file"
	}

	output.Fail(fmt.Sprintf(
		`%s not found.

Expected location:
    %s`,
		description,
		path,
	))
}
