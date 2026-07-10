package output

import (
	"fmt"
	"os"
)

func Info(message string) {
	fmt.Println(message)
}

func Fail(message string) {
	fmt.Fprintln(os.Stderr, message)
	os.Exit(1)
}
