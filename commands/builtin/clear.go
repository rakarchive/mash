package builtin

import (
	"fmt"
	"os"
)

func clear(args []string) error {
	if len(args) > 0 {
		fmt.Fprintln(os.Stderr, "clear: too many arguments")
		return &ExitError{1}
	}
	//fmt.Print("\u001b[2J")
	fmt.Print("\u001bc")
	return nil
}
