// mash
// https://github.com/raklaptudirm/mash
// Copyright (c) 2021 Rak Laptudirm.
// Licensed under the MIT license.

// mash is a simple shell written in go.
// Features:
// - cd command
// - exit command
// - run executable files

package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"syscall"

	"github.com/raklaptudirm/mash/commands"
	"github.com/raklaptudirm/mash/parser"
)

func main() {
	reader := bufio.NewReader(os.Stdin)

	// Catch ctrl+c and SIGTERM events
	ctrlC := make(chan os.Signal, 1)
	signal.Notify(ctrlC, os.Interrupt)
	signal.Notify(ctrlC, syscall.SIGTERM)

	// Command loop
	for {
		cwd, _ := os.Getwd()

		// Prompt
		fmt.Printf("\u001b[32m%v\u001b[0m\n$ ", cwd)

		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		} else {
			cmd, args := parser.Parse(input)
			if err := commands.Dispatch(cmd, args); err != nil {
				if errors.Is(err, exec.ErrNotFound) {
					fmt.Fprintln(os.Stderr, "mash: "+cmd+": command not found")
				} else {
					fmt.Println(err)
				}
			}
		}
	}
}
