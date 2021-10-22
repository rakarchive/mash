package builtin

import (
	"fmt"
	"os"
)

// Clear command is used to clear the terminal,
// including scroll-back (for now).
func clear(args []string) error {
	if len(args) > 0 {
		fmt.Fprintln(os.Stderr, "clear: too many arguments")
		return &ExitError{1}
	}

	// Escape sequence to preserve scroll-back:
	// fmt.Print("\u001b[2J")
	fmt.Print("\u001bc")
	return nil
}
