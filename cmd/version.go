package cmd

import (
	"fmt"
	"jarvis/config"
)

func ShowVersion() {
	fmt.Printf("%v v%v\n", config.AppName, config.Version)
}
