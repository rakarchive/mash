// mash
// https://github.com/raklaptudirm/mash
// Copyright (c) 2021 Rak Laptudirm.
// Licensed under the MIT license.

// mash is a simple shell written in go.
//
// Usage:
//
//  mash
//  mash [ command ]
//
// "mash" starts the shell process in the current
// terminal. Currently, customization of any of the
// features of mash is unavailable, but will be added
// soon.
//
// "mash command" runs the provided command and exits.
//
// The shell consists of a command loop, which can
// take user input and execute the commands accordingly.
//
// Mash provides support for both builtin and external
// commands, both of which are used the same way.
//
// The shell reports the exit code of all types of
// commands in the case that it is non-zero. The shell
// goes into another iteration of the command loop in
// the case of such errors, and only exits in case of
// serious errors like unable to read input or failing
// to fetch the working directory.
//
// Process exit signals like ^C and SIGTERM are caught
// and the shell does not exit as a result of them.
//
package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/raklaptudirm/mash/pkg/shell"
)

func main() {
	// Catch ctrl+c and SIGTERM events so as not to
	// interrupt the shell input, unlike normal
	// processes.
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)
	signal.Notify(interrupt, syscall.SIGTERM)

	args := os.Args[1:]
	if len(args) < 1 {
		// Start shell instance if no args provided.
		shell.Start()
	} else {
		// If an argument is provided, run it as a command.
		shell.Run(args[0])
	}
}
