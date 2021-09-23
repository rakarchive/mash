// mash
// https://github.com/raklaptudirm/mash
// Copyright (c) 2021 Rak Laptudirm.
// Licensed under the MIT license.

// Package builtin provides the builtin
// command functionality of mash.

package commands

import (
	"errors"

	"github.com/raklaptudirm/mash/commands/builtin"
)

// Function run runs a given external
// command cmd with the given set of
// arguments args.
func run(cmd string, args []string) error {
	if isBuiltin(cmd) {
		error := builtin.Commands[cmd](args)
		return error
	} else {
		return errors.New("mash: command " + cmd + " not found")
	}
}

// Function isBuiltin checks if the
// provided command cmd is a builtin
// command of mash.
func isBuiltin(cmd string) bool {
	_, exists := builtin.Commands[cmd]
	return exists
}
