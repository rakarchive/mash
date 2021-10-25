// mash
// https://github.com/raklaptudirm/mash
// Copyright (c) 2021 Rak Laptudirm.
// Licensed under the MIT license.

// Package vm provides functions to run an unparsed
// command string by parsing and dispatching it.

package shell

import (
	"errors"
	"fmt"
	"os"
	"os/exec"

	"github.com/raklaptudirm/mash/commands/special"
	"github.com/raklaptudirm/mash/parser"
	"github.com/raklaptudirm/mash/runners"
)

// Function Run parses a given command string,
// sends it to the dispatcher, and prints error
// messages depending on the returned error.
func Run(command string) {
	cmd, args := parser.Parse(command)
	err := dispatch(cmd, args)

	// report errors from Dispatch
	switch {
	case errors.Is(err, exec.ErrNotFound):
		fmt.Fprintln(os.Stderr, "mash: "+cmd+": command not found")
	case err != nil:
		fmt.Println(err)
	}
}

// Function discpatch dispatches the execution
// of the provided command, depending on the
// type of the command.
func dispatch(command string, args []string) error {
	function, exists := special.Commands[command]
	if exists {
		return function(args)
	}

	err := runners.Builtin(command, args)
	if _, is := err.(*runners.NotBuiltinCmd); !is {
		return err
	}

	return runners.External(command, args)
}
