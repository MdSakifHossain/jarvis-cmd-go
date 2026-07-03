package cmd

import "jarvis/banner"

func Banner() {
	showHeader(banner.Jarvis)
	showHeader(banner.Lights)
	showHeader(banner.Lock)
	showHeader(banner.Unlock)
	showHeader(banner.Power)
	showHeader(banner.Tree)
	showHeader(banner.NmHunter)
	showHeader(banner.Bkash)
}
