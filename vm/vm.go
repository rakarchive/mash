// mash
// https://github.com/raklaptudirm/mash
// Copyright (c) 2021 Rak Laptudirm.
// Licensed under the MIT license.

// Package vm provides functions to run an unparsed
// command string by parsing and dispatching it.

package vm

import (
	"errors"
	"fmt"
	"os"
	"os/exec"

	"github.com/raklaptudirm/mash/commands"
	"github.com/raklaptudirm/mash/parser"
)

// Function Run parses a given command string,
// sends it to the dispatcher, and prints error
// messages depending on the returned error.
func Run(command string) {
	cmd, args := parser.Parse(command)
	err := commands.Dispatch(cmd, args)

	// report errors from Dispatch
	switch {
	case errors.Is(err, exec.ErrNotFound):
		fmt.Fprintln(os.Stderr, "mash: "+cmd+": command not found")
	case err != nil:
		fmt.Println(err)
	}
}
