package builtin

import (
	"errors"
	"os"
)

var Commands = map[string]func([]string) error{
	"cd": func(args []string) error {
		var path string
		length := len(args)

		if length < 1 {
			path, _ = os.UserHomeDir()
		} else if length == 1 {
			path = args[0]
		} else {
			return errors.New("mash: cd: too many arguments")
		}
		return os.Chdir(path)
	},
	"exit": func(args []string) error {
		length := len(args)
		if length < 1 {
			os.Exit(0)
		} else if length > 1 {
			return errors.New("mash: exit: too many arguments")
		} else {
			os.Exit(0)
		}
		return nil
	},
}
