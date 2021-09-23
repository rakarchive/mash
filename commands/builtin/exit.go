// mash
// https://github.com/raklaptudirm/mash
// Copyright (c) 2021 Rak Laptudirm.
// Licensed under the MIT license.

package builtin

import (
	"errors"
	"os"
	"strconv"
)

// Function exit terminates the shell process
// with a return value according to args,
// which should have 0 or 1 item, which should
// be the exit code (default 0).
func exit(args []string) error {
	length := len(args)
	if length < 1 {
		os.Exit(0)
	}
	if length > 1 {
		return errors.New("exit: too many arguments")
	}
	if num, err := strconv.Atoi(args[0]); err == nil {
		os.Exit(num)
	} else {
		return errors.New("exit: expected numeric argument")
	}
	os.Exit(0)
	return nil
}
