// mash
// https://github.com/raklaptudirm/mash
// Copyright (c) 2021 Rak Laptudirm.
// Licensed under the MIT license.

// Package commands provides types used by
// the implementations of different types
// of builtin commands.
//
package commands

import "fmt"

// An ExitError is the error returned by
// builtin commands in case of a non-zero
// exit value.
type ExitError struct {
	Code int // Exit code returned by process
}

func (err *ExitError) Error() string {
	return fmt.Sprintf("exit status %v", err.Code)
}
