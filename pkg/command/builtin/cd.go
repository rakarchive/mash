// mash
// https://github.com/raklaptudirm/mash
// Copyright (c) 2021 Rak Laptudirm.
// Licensed under the MIT license.

package builtin

import (
	"fmt"
	"os"

	"github.com/raklaptudirm/mash/pkg/command"
)

// cd changes the current working directory of the shell
// according to the arguments args, which should have 0-1
// arguments, which should be the new working directory
// (defaults to homepath).
func cd(args []string) error {
	var path string
	length := len(args)

	switch {
	case length < 1:
		// If no arguments provided, go to homepath
		home, err := os.UserHomeDir()
		if err != nil {
			return err
		}
		path = home
	case length > 1:
		fmt.Fprintln(os.Stderr, "mash: cd: too many arguments")
		return &command.ExitError{Code: 1}
	default:
		path = args[0]
	}

	err := os.Chdir(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "mash: cd: %v: no such file or directory\n", path)
		return &command.ExitError{Code: 1}
	}

	return nil
}