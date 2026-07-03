package cmd

import (
	"fmt"
)

func showHeader(banner func()) {
	fmt.Println()
	banner()
}
