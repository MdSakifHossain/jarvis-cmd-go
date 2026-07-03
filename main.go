package main

import (
	"jarvis/cmd"
	"os"
)

func main() {
	args := os.Args[1:]

	if len(args) == 0 {
		cmd.ShowHelp()
		return
	}

	switch args[0] {
	case "help":
		cmd.ShowHelp()
	case "version":
		cmd.ShowVersion()
	case "lights":
		cmd.Lights(args[1:])
	case "banner":
		cmd.Banner()
	default:
		cmd.ShowHelp()
	}
}
