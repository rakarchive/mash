// mash
// https://github.com/raklaptudirm/mash
// Copyright (c) 2021 Rak Laptudirm.
// Licensed under the MIT license.

// Package builtin implements functions which
// provide the functionality for normal builtin
// commands.
//
package builtin

import "github.com/raklaptudirm/mash/pkg/command"

// Commands maps the normal builtin command functions
// to their names.
var Commands = command.CommandMap{
	"cd": command.SimpleCommand{
		Name:    "cd",
		Args:    []string{},
		Builtin: cd,
	},
	"exit": command.SimpleCommand{
		Name:    "exit",
		Args:    []string{},
		Builtin: exit,
	},
	"clear": command.SimpleCommand{
		Name:    "clear",
		Args:    []string{},
		Builtin: clear,
	},
	"echo": command.SimpleCommand{
		Name:    "echo",
		Args:    []string{},
		Builtin: echo,
	},
}
