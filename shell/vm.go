// mash
// https://github.com/raklaptudirm/mash
// Copyright (c) 2021 Rak Laptudirm.
// Licensed under the MIT license.

// Package shell provides an interface for wrapping
// the tasks of parsing, execution and error
// reporting for simpler and steadier running of
// the tasks.
//
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

// Run wraps the tasks of parsing a command string, executing the
// command with the provided args, and reporting errors, if any,
// from the execution.
//
// The errors are classified according to their type, and for some
// types of errors, Run prints them using a custom format. For other
// types of errors, their respective Error method is used for printing.
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

// dispatch searches for the provided command in the builtin command
// maps, and if found, executes the command function and reports any
// returned error. Otherwise, it runs the command as an external command
// and returns any raised error.
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
