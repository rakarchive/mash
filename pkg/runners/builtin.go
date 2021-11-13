// mash
// https://github.com/raklaptudirm/mash
// Copyright (c) 2021 Rak Laptudirm.
// Licensed under the MIT license.

// Package runners provides runners, i.e.
// command executors for normal builtin,
// and external commands.
//
package runners

import (
	"fmt"

	"github.com/raklaptudirm/mash/pkg/command/builtin"
)

// A NotBuiltinCmd is the type of error returned
// by the Builtin runner if the provided command
// is not a valid builtin command.
type NotBuiltinCmd struct {
	name string // The command name
}

func (e *NotBuiltinCmd) Error() string {
	return fmt.Sprintf("mash: %v: not a builtin command", e.name)
}

// Builtin runs the provided builtin command, or returns a
// NotBuiltinCmd error if it is not one.
func Builtin(command string, args []string) error {
	cmd, exists := builtin.Commands[command]
	if !exists {
		return &NotBuiltinCmd{command}
	}
	cmd.Args = args
	return cmd.Execute()
}
