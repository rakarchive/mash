// mash
// https://github.com/raklaptudirm/mash
// Copyright (c) 2021 Rak Laptudirm.
// Licensed under the MIT license.

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
	"strings"
)

// Function Parse accepts user command as
// input and separates the command string
// from the argument string, and then
// dispatches the argument string to the
// Args function. It returns the command
// and the argument slice.
func Parse(input string) (string, []string) {
	input = strings.Trim(input, " \t\n\r")
	cmdTill := strings.Index(input, " ")
	if cmdTill == -1 {
		return input, []string{}
	}
	cmd := input[:cmdTill]
	args := Args(input[cmdTill:])
	return cmd, args
}

// Function Args parses an argument string
// input into an argument slice.
func Args(input string) []string {
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
				str += "\""
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
