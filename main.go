// mash
// https://github.com/raklaptudirm/mash
// Copyright (c) 2021 Rak Laptudirm.
// Licensed under the MIT license.

// mash is a simple shell written in go.
// Features:
// - cd command
// - exit command
// - run executable files
//
package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
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
			if err := dispatch(cmd, args); err != nil {
				fmt.Fprintln(os.Stderr, err)
			}
		}
	}
}

// Function dispatch determines wether
// a command is a shell command or an
// exe file, and executes it appropriately.
func dispatch(command string, args []string) error {
	switch command {
	case "cd":
		var path string
		if len(args) < 1 {
			path, _ = os.UserHomeDir()
		} else {
			path = args[0]
		}
		return os.Chdir(path)
	case "exit":
		os.Exit(0)
	default:
		return execute(command, args)
	}
	return nil
}

// Function execute command handles
// exe file executions requested by
// dispatch.
func execute(command string, args []string) error {
	cmd := exec.Command(command, args...)

	// Set command streams to shell default
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}
