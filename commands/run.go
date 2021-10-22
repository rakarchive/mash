// mash
// https://github.com/raklaptudirm/mash
// Copyright (c) 2021 Rak Laptudirm.
// Licensed under the MIT license.

// Functions to run builtin and external
// commands and to dispatch them to the
// concerned function.

package commands

import (
	"os"
	"os/exec"

	"github.com/raklaptudirm/mash/commands/builtin"
)

// Function dispatch determines wether a command
// is a shell command or an executable file, and
// executes it appropriately.
func Dispatch(command string, args []string) error {
	// I am not falling for that one.
	if command == "" {
		return nil
	}

	function, exists := builtin.Commands[command]
	if exists {
		return function(args)
	}
	return execute(command, args)
}

// Function execute command handles
// executable file executions requested
// by dispatch.
func execute(command string, args []string) error {
	cmd := exec.Command(command, args...)

	// Set File streams to shell default, to
	// allow direct manipulation from the shell.
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}
