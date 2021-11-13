// mash
// https://github.com/raklaptudirm/mash
// Copyright (c) 2021 Rak Laptudirm.
// Licensed under the MIT license.

// This whole thing needs to be a hell lot
// more powerful to do the stuff I need it
// to do later in this journey.

// Package parser provides functions to
// parse the user command input into
// command and arguments.
//
// Specification:
// - Leading and trailing whitespace are
//  considered void.
// - Command is the string up to the first
//  internal whitespace.
// - Arguments are seperated by whitespace.
// - A collection of whitespace is considered
//  to be one.
// - Whitespace is taken literally inside
//  quotes.
// - \ is used to escape control characters.

package parser

import (
	"fmt"
	"strings"

	"github.com/raklaptudirm/mash/pkg/command"
	"github.com/raklaptudirm/mash/pkg/command/builtin"
)

var NoCommandError = fmt.Errorf("No commands provided")

// Function Parse accepts user command as
// input and separates the command string
// from the argument string, and then
// dispatches the argument string to the
// Args function. It returns the command
// and the argument slice.
func Parse(input string) (command.Command, error) {
	input = strings.Trim(input, " \t\n\r")
	words := Words(input)
	if len(words) == 0 {
		return command.Command{}, NoCommandError
	}
	cmd, exists := builtin.Commands[words[0]]
	if !exists {
		return command.Command{
			Name: words[0],
			Args: words[1:],
		}, nil
	}
	cmd.Args = words[1:]
	return cmd, nil
}

// Function Args parses an argument string
// input into an argument slice.
func Words(input string) []string {
	length := len(input)
	current := 0
	args := []string{}

	str := ""
	inDoubleQuote := false
	inSingleQuote := false
	for current < length {
		char := input[current]
		switch char {
		case '"':
			if inDoubleQuote {
				inDoubleQuote = false
			} else if inSingleQuote {
				str += "\""
			} else {
				inDoubleQuote = true
			}
		case '\'':
			if inSingleQuote {
				inSingleQuote = false
			} else if inDoubleQuote {
				str += "'"
			} else {
				inSingleQuote = true
			}
		case ' ', '\t', '\r':
			if inDoubleQuote || inSingleQuote {
				str += string(char)
				break
			}
			if str != "" {
				args = append(args, str)
				str = ""
			}
		case '\\':
			if current == length-1 {
				break
			}
			str += escape(input[current+1])
			current++
		default:
			str += string(char)
		}
		current++
	}
	if str != "" {
		args = append(args, str)
	}
	return args
}

func escape(char byte) string {
	switch char {
	case 'n':
		return "\n"
	case 'r':
		return "\r"
	case 't':
		return "\t"
	case '\\', '"', '\'':
		return string(char)
	default:
		return ""
	}
}
