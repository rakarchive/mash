// mash
// https://github.com/raklaptudirm/mash
// Copyright (c) 2021 Rak Laptudirm.
// Licensed under the MIT license.

package special

import (
	"fmt"
	"os"

	"github.com/raklaptudirm/mash/commands"
	"github.com/raklaptudirm/mash/runners"
)

// Function builtin implements the special function builtin that
// is used to run only normal builtin commands.
//
// The first argument is used as the builtin command name and the
// rest are the proper args.
func builtin(args []string) error {
	length := len(args)
	switch {
	case length == 1:
		return runners.Builtin(args[0], []string{})
	case length > 1:
		return runners.Builtin(args[0], args[1:])
	case length < 1:
		fmt.Fprintln(os.Stderr, "mash: builtin: expected builtin command")
		return &commands.ExitError{Code: 1}
	}
	return nil
}
