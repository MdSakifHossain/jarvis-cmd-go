package cmd

import (
	"jarvis/banner"
	"jarvis/output"
	"jarvis/support"
	"os/exec"
)

func Power() {
	support.ShowBanner(banner.Power)
	output.Info("Initializing command...")
	output.Info("Shutting down...")

	cmd := exec.Command("sudo", "shutdown", "now")
	if err := cmd.Run(); err != nil {
		output.Fail("Failed to Power off")
	}
}
