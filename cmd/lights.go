package cmd

import (
	"fmt"
	"jarvis/banner"
	"jarvis/config"
)

func Lights(args []string) {

	if len(args) == 0 {
		showLightsHelp()
		return
	}

	switch args[0] {
	case "on":
		fmt.Println("Trun on the RAM lights")
	case "off":
		fmt.Println("Trun off the RAM lights")
	case "help":
		showLightsHelp()

	default:
		showLightsHelp()

	}
}

func showLightsHelp() {
	fmt.Println("")
	banner.Lights()
	fmt.Printf(`Change Color of RAM LED

Usage:

    %s lights [command]

Available Commands:

    on        Turn on RAM LED
    off       Turn off RAM LED
    help      Show help

`, config.AppName)
}
