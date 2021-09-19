package commands

import (
	"errors"
	"os"
)

func Cd(args []string) error {
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
}

func Exit(args []string) error {
	length := len(args)
	if length < 1 {
		os.Exit(0)
	} else if length > 1 {
		return errors.New("mash: exit: too many arguments")
	} else {
		os.Exit(0)
	}
	return nil
}
