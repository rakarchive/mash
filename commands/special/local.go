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

// Function local implements the special function
// local that is used to run only external commands.
//
// The first argument is used as the external command
// while the rest are used as the proper args.
func local(args []string) error {
	length := len(args)
	switch {
	case length == 1:
		return runners.External(args[0], []string{})
	case length > 1:
		return runners.External(args[0], args[1:])
	case length < 1:
		fmt.Fprintln(os.Stderr, "mash: local: expected local command")
		return &commands.ExitError{Code: 1}
	}
	return nil
}
