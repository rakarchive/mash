// mash
// https://github.com/raklaptudirm/mash
// Copyright (c) 2021 Rak Laptudirm.
// Licensed under the MIT license.

// Package special implements functions which provide the
// functionality for special builtin commands.
//
package special

import "github.com/raklaptudirm/mash/commands"

// CommandMap Commands maps the special builtin command
// functions to the command names.
var Commands = commands.CommandMap{
	"builtin": builtin,
	"local":   local,
}
