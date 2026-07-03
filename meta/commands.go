package meta

type Command struct {
	Name        string
	Description string
}

var Commands = []Command{
	{"help", "Show help information"},
	{"version", "Show application version"},
	{"lights", "Control RAM lighting"},
	{"lock", "Lock the screen"},
	{"unlock", "Unlock the screen"},
}
