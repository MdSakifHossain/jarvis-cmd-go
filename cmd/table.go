package cmd

import (
	"fmt"
	"jarvis/banner"
	"jarvis/output"
	"jarvis/support"
	"strconv"
)

func Table(args []string) {
	if len(args) == 0 {
		showTableHelp()
		return
	}

	if args[0] == "help" {
		showTableHelp()
		return
	}

	number, err := strconv.Atoi(args[0])
	if err != nil {
		output.Fail("Number must be an integer.")
	}

	limit := 10

	if len(args) > 1 {
		limit, err = strconv.Atoi(args[1])
		if err != nil {
			output.Fail("Limit must be an integer.")
		}
	}

	showTableBanner()

	for i := 1; i <= limit; i++ {
		fmt.Printf("%4d × %-2d = %d\n", number, i, number*i)
	}
}

func showTableBanner() {
	support.ShowBanner(banner.Table)
}

func showTableHelp() {
	showTableBanner()

	fmt.Println(`Show a multiplication table.

Usage:

    jarvis table <number>
    jarvis table <number> <limit>

Examples:

    jarvis table 9
    jarvis table 9 30`)
}
