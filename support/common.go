package support

import (
	"fmt"
)

func ShowBanner(banner func()) {
	fmt.Println()
	banner()
	fmt.Println()
}
