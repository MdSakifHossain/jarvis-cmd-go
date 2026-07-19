package cmd

import (
	"jarvis/output"
	"os/exec"
)

func Power() {
	output.Info("Initializing command...")
	output.Info("Shutting down...")

	cmd := exec.Command("sudo", "shutdown", "now")
	if err := cmd.Run(); err != nil {
		output.Fail("Failed to Power off")
	}
}
