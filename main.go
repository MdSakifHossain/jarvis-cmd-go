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
	case "lock":
		cmd.Lock()
	case "unlock":
		cmd.Unlock()
	case "table":
		cmd.Table(args[1:])
	case "observe":
		cmd.Observe()
	case "power":
		cmd.Power()
	case "tree":
		cmd.Tree()
	case "banner":
		cmd.Banner()
	case "attendance":
		cmd.Attendance()
	default:
		cmd.ShowHelp()
	}
}
