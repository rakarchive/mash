// mash
// https://github.com/raklaptudirm/mash
// Copyright (c) 2021 Rak Laptudirm.
// Licensed under the MIT license.

package builtin

import (
	"fmt"
	"os"
)

// Function cd changes the current working
// directory of the shell according to the
// arguments args, which should have 0-1
// arguments, which should be the new
// working directory (defaults to homepath).
func cd(args []string) error {
	var path string
	length := len(args)

	if length < 1 {
		home, err := os.UserHomeDir()
		if err != nil {
			return err
		}
		path = home
	} else if length == 1 {
		path = args[0]
	} else {
		fmt.Fprintln(os.Stderr, "mash: cd: too many arguments")
		return &ExitError{1}
	}

	err := os.Chdir(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "mash: cd: %v: no such file or directory\n", path)
		return &ExitError{1}
	}

	return nil
}
