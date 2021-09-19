// Package builtin provides the builtin
// functionality and commands of mash.
//
package builtin

import (
	"errors"
	"github.com/raklaptudirm/mash/builtin/commands"
)

//
func Run(cmd string, args []string) error {
	switch cmd {
	case "cd":
		return commands.Cd(args)
	case "exit":
		return commands.Exit(args)
	}
	return errors.New("mash: command " + cmd + " not found")
}

func IsCmd(cmd string) bool {
	switch cmd {
	case "cd", "exit":
		return true
	default:
		return false
	}
}
