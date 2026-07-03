package cmd

import (
	"fmt"
	"jarvis/banner"
	"jarvis/console"
	"jarvis/meta"
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
	support.ShowBanner(banner.Lights)
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

	support.ShowBanner(banner.Lights)
	console.Info("Turning lights ON...")
	setLights("ffffff")
	console.Info("Done.")
}

func lightsOff() {
	support.RequireOpenRGB()

	support.ShowBanner(banner.Lights)
	console.Info("Turning lights OFF...")
	setLights("000000")
	console.Info("Done.")
}

func setLights(color string) {
	cmd := exec.Command(
		"openrgb",
		"--mode", "static",
		"--color", color,
	)

	if err := cmd.Run(); err != nil {
		console.Fail(fmt.Sprintf("OpenRGB failed: %v", err))
	}
}
