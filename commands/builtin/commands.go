// mash
// https://github.com/raklaptudirm/mash
// Copyright (c) 2021 Rak Laptudirm.
// Licensed under the MIT license.

// Package builtin implements functions which
// provide the functionality for normal builtin
// commands.
//
package builtin

import "github.com/raklaptudirm/mash/commands"

// Commands maps the normal builtin command functions
// to their names.
var Commands = commands.CommandMap{
	"cd":    cd,
	"exit":  exit,
	"clear": clear,
	"echo":  echo,
}
