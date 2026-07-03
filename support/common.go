package support

import (
	"fmt"
)

func ShowHeader(banner func()) {
	fmt.Println()
	banner()
}
