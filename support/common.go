package support

import (
	"fmt"
	"jarvis/output"
	"os"
)

func ShowBanner(banner func()) {
	fmt.Println()
	banner()
	fmt.Println()
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
