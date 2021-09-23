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
)

// Function dispatch determines wether
// a command is a shell command or an
// exe file, and executes it appropriately.
func Dispatch(command string, args []string) error {
	if isBuiltin(command) {
		return run(command, args)
	}
	return execute(command, args)
}

// Function execute command handles
// exe file executions requested by
// dispatch.
func execute(command string, args []string) error {
	cmd := exec.Command(command, args...)

	// Set command streams to shell default
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}
