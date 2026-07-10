package cmd

import (
	"fmt"
	"jarvis/banner"
	"jarvis/output"
	"jarvis/support"
	"os"
	"os/exec"
	"path/filepath"
)

func Observe() {
	showObserveHeader()

	home, err := os.UserHomeDir()
	if err != nil {
		output.Fail(fmt.Sprintf("Failed to determine home directory: %v", err))
	}

	logFile := filepath.Join(home, ".local", "logs", "vault-observer.log")

	support.RequireFile(logFile, "Vault observer log file")

	cmd := exec.Command("tail", "-f", logFile)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	if err := cmd.Run(); err != nil {
		output.Fail(fmt.Sprintf("Failed to observe log: %v", err))
	}
}

func showObserveHeader() {
	support.ShowBanner(banner.Vault)
}
