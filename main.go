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
	case "attendance":
		cmd.Attendance()
	case "ph":
		cmd.PH()
	case "nmhunt":
		cmd.NMHunter(args[1:])
	case "bkash":
		cmd.BKash(args[1:])
	default:
		cmd.ShowHelp()
	}
}
