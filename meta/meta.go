package meta

type Command struct {
	Name        string
	Description string
}

const AppName = "jarvis"
const Version = "0.8.0"
const ShortDescription = "Personal CLI Tool"

var Commands = []Command{
	{"help", "Show help information"},
	{"version", "Show application version"},
	{"lights", "Control RAM lighting"},
	{"lock", "Lock the screen"},
	{"unlock", "Unlock the screen"},
	{"table", "Show a multiplication table"},
	{"observe", "Show live log file of Vault Observer"},
}
