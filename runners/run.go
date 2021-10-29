// mash
// https://github.com/raklaptudirm/mash
// Copyright (c) 2021 Rak Laptudirm.
// Licensed under the MIT license.

package runners

import (
	"os"
	"os/exec"
)

// External runs the provided external command, and returns the
// error returned by the cmd.Run() call.
func External(command string, args []string) error {
	cmd := exec.Command(command, args...)

	// Set File streams to shell default, to
	// allow direct manipulation from the terminal.
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}
