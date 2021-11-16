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

	"github.com/raklaptudirm/mash/pkg/command"
	"github.com/raklaptudirm/mash/pkg/parser"
)

// Run wraps the tasks of parsing a command string, executing the
// command with the provided args, and reporting errors, if any,
// from the execution.
//
// The errors are classified according to their type, and for some
// types of errors, Run prints them using a custom format. For other
// types of errors, their respective Error method is used for printing.
func Run(command string) {
	cmd, err := parser.Parse(command)
	if err != nil {
		fmt.Println(err)
	}
	err = dispatch(cmd)

	// report errors from Dispatch
	switch {
	case errors.Is(err, exec.ErrNotFound):
		fmt.Fprintf(os.Stderr, "mash: %v: command not found", cmd)
	case err != nil:
		fmt.Println(err)
	}
}

// dispatch searches for the provided command in the builtin command
// maps, and if found, executes the command function and reports any
// returned error. Otherwise, it runs the command as an external command
// and returns any raised error.
func dispatch(cmd command.SimpleCommand) error {
	return cmd.Execute()
}
