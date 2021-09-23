// mash
// https://github.com/raklaptudirm/mash
// Copyright (c) 2021 Rak Laptudirm.
// Licensed under the MIT license.

package builtin

// Map commands maps the in-built command
// name to the respective go function.
var Commands = map[string]func([]string) error{
	"cd":   cd,
	"exit": exit,
}
