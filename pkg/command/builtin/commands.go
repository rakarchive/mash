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
	"cd": command.Command{
		Name:    "cd",
		Args:    []string{},
		Builtin: cd,
	},
	"exit": command.Command{
		Name:    "exit",
		Args:    []string{},
		Builtin: exit,
	},
	"clear": command.Command{
		Name:    "clear",
		Args:    []string{},
		Builtin: cd,
	},
	"echo": command.Command{
		Name:    "echo",
		Args:    []string{},
		Builtin: echo,
	},
}
