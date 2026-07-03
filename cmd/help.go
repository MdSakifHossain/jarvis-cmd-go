package cmd

import (
	"fmt"
	"jarvis/config"
)

func ShowHelp() {
	fmt.Printf(
		`%s - v%s
%s

Usage:

    %s [command]

Available commands:

    help, version, lights, lock, unlock

For more info, Run:

    %s [command] help

`,
		config.AppName, config.Version, config.ShortDescription, config.AppName, config.AppName)
}
