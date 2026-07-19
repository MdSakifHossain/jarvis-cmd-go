package cmd

import (
	"fmt"
	"jarvis/meta"
	"jarvis/output"
	"jarvis/support"
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
	fmt.Printf(`Change Color of RAM LED

Usage:

    %s lights [command]

Available Commands:

    on        Turn on RAM LED
    off       Turn off RAM LED
    help      Show help

`, meta.AppName)
}

func lightsOn() {
	support.RequireOpenRGB()
	output.Info("Turning lights ON...")
	setLights("ffffff")
	output.Info("Done.")
}

func lightsOff() {
	support.RequireOpenRGB()
	output.Info("Turning lights OFF...")
	setLights("000000")
	output.Info("Done.")
}

func setLights(color string) {
	cmd := exec.Command(
		"openrgb",
		"--mode", "static",
		"--color", color,
	)

	if err := cmd.Run(); err != nil {
		output.Fail(fmt.Sprintf("OpenRGB failed: %v", err))
	}
}
