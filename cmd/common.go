package cmd

import (
	"fmt"
	"os"
)

func fail(message string) {
	fmt.Fprintln(os.Stderr)
	fmt.Fprintln(os.Stderr, message)
	os.Exit(1)
}

func showHeader(banner func()) {
	fmt.Println()
	banner()
}
