package cmd

import (
	"fmt"
	"jarvis/banner"
	"jarvis/output"
	"jarvis/support"
	"os"
	"os/exec"
)

func Tree() {
	support.ShowBanner(banner.Tree)
	support.RequireTree()
	cmd := exec.Command("tree", "--gitignore", "--dirsfirst")
	// telling command to output the commands output and error on the terminal
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		output.Fail(fmt.Sprintf("Tree failed: %v", err))
	}
}
