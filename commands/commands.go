// Package builtin provides the builtin
// functionality and commands of mash.
//
package commands

import (
	"errors"

	"github.com/raklaptudirm/mash/commands/builtin"
)

//
func Run(cmd string, args []string) error {
	if IsBuiltin(cmd) {
		error := builtin.Commands[cmd](args)
		return error
	} else {
		return errors.New("mash: command " + cmd + " not found")
	}
}

func IsBuiltin(cmd string) bool {
	_, exists := builtin.Commands[cmd]
	return exists
}
