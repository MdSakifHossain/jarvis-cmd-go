package cmd

import (
	"fmt"
	"jarvis/meta"
)

func ShowVersion() {
	fmt.Printf("%v v%v\n", meta.AppName, meta.Version)
}
