package cmd

import (
	"jarvis/banner"
	"jarvis/support"
)

func Banner() {
	support.ShowBanner(banner.Jarvis)
	support.ShowBanner(banner.Lights)
	support.ShowBanner(banner.Lock)
	support.ShowBanner(banner.Unlock)
	support.ShowBanner(banner.Power)
	support.ShowBanner(banner.Tree)
	support.ShowBanner(banner.NmHunter)
	support.ShowBanner(banner.Bkash)
}
