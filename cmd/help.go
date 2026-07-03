package cmd

import (
	"fmt"
	"jarvis/banner"
	"jarvis/config"
	"jarvis/meta"
	"jarvis/support"
)

func ShowHelp() {
	support.ShowBanner(banner.Jarvis)
	fmt.Printf(
		`%s - v%s - %s

Usage:

    %s [command]

Available commands:

`,
		config.AppName,
		config.Version,
		config.ShortDescription,
		config.AppName,
	)

	for _, command := range meta.Commands {
		fmt.Printf("    %-12s %s\n", command.Name, command.Description)
	}

	fmt.Printf(
		`

For more info, Run:

    %s [command] help

`,
		config.AppName,
	)
}
