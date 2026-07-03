package cmd

import (
	"jarvis/banner"
	"jarvis/support"
)

func Banner() {
	support.ShowHeader(banner.Jarvis)
	support.ShowHeader(banner.Lights)
	support.ShowHeader(banner.Lock)
	support.ShowHeader(banner.Unlock)
	support.ShowHeader(banner.Power)
	support.ShowHeader(banner.Tree)
	support.ShowHeader(banner.NmHunter)
	support.ShowHeader(banner.Bkash)
}
