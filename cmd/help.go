package cmd

import (
	"fmt"
	"jarvis/meta"
)

func ShowHelp() {
	fmt.Printf(
		`%s - v%s - %s

Usage:

    %s [command]

Available commands:

`,
		meta.AppName,
		meta.Version,
		meta.ShortDescription,
		meta.AppName,
	)

	for _, command := range meta.Commands {
		fmt.Printf("    %-12s %s\n", command.Name, command.Description)
	}

	fmt.Printf(
		`
For more info, Run:

    %s [command] help

`,
		meta.AppName,
	)
}
