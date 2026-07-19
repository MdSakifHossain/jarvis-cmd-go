package meta

import "sort"

type Command struct {
	Name        string
	Description string
}

const AppName = "jarvis"
const Version = "0.15.0"
const ShortDescription = "Personal CLI Tool"

var Commands = []Command{
	{"help", "Show help information"},
	{"version", "Show application version"},
	{"lights", "Control RAM lighting"},
	{"lock", "Lock the screen"},
	{"unlock", "Unlock the screen"},
	{"table", "Show a multiplication table"},
	{"observe", "Show live log file of Vault Observer"},
	{"power", "Turn off PC's Power"},
	{"tree", "Same as original but with extra flags"},
	{"ph", "Scaffold new module of PH with correct Connection"},
	{"attendance", "Create an Attendance Sheet on current dir"},
	{"nmhunt", "Runs Node_Modules hunter"},
	{"bkash", "Bkash and other MFS Calculations"},
}

func init() {
	sort.Slice(Commands, func(i, j int) bool {
		return Commands[i].Name < Commands[j].Name
	})
}
