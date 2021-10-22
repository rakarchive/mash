// mash
// https://github.com/raklaptudirm/mash
// Copyright (c) 2021 Rak Laptudirm.
// Licensed under the MIT license.

package builtin

import "fmt"

// An ExitError is the error returned by
// builtin commands in case of a non-zero
// exit value.
type ExitError struct {
	ErrorCode int
}

func (err *ExitError) Error() string {
	return fmt.Sprintf("exit status %v", err.ErrorCode)
}

// Map commands maps the in-built command
// name to the respective go function.
var Commands = map[string]func([]string) error{
	"cd":    cd,
	"exit":  exit,
	"clear": clear,
}
