package cmd

import (
	"fmt"
	"jarvis/banner"
	"jarvis/config"
	"os/exec"
)

func Lights(args []string) {
	if len(args) == 0 {
		showLightsHelp()
		return
	}

	switch args[0] {
	case "on":
		lightsOn()
	case "off":
		lightsOff()
	case "help":
		showLightsHelp()
	default:
		showLightsHelp()
	}
}

func showLightsHelp() {
	showHeader(banner.Lights)
	fmt.Printf(`Change Color of RAM LED

Usage:

    %s lights [command]

Available Commands:

    on        Turn on RAM LED
    off       Turn off RAM LED
    help      Show help

`, config.AppName)
}

func lightsOn() {
	requireOpenRGB()

	showHeader(banner.Lights)
	fmt.Println("Turning lights ON...")
	setLights("ffffff")
	fmt.Println("Done.")
}

func lightsOff() {
	requireOpenRGB()

	showHeader(banner.Lights)
	fmt.Println("Turning lights OFF...")
	setLights("000000")
	fmt.Println("Done.")
}

func setLights(color string) {
	cmd := exec.Command(
		"openrgb",
		"--mode", "static",
		"--color", color,
	)

	if err := cmd.Run(); err != nil {
		fail(fmt.Sprintf("OpenRGB failed: %v", err))
	}
}
