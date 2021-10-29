// mash
// https://github.com/raklaptudirm/mash
// Copyright (c) 2021 Rak Laptudirm.
// Licensed under the MIT license.

package commands

// A CommandMap is the map type used to
// connect builtin command functions to
// their names.
type CommandMap map[string]func([]string) error
