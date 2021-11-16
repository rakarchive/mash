package command

import (
	"fmt"
	"io"
	"os"
	"os/exec"
)

type Command struct {
	SimpleCommands []SimpleCommand
	// the following will make it easy to redirect using pipe and > or <
	Output               io.Writer
	Input                io.Reader
	ErrOutput            io.Writer
	CurrentCommand       *Command
	CurrentSimpleCommand *SimpleCommand
}

type SimpleCommand struct {
	Name    string
	Args    []string
	Builtin Builtin
}

type Builtin func([]string) error

func (c *SimpleCommand) String() string {
	return c.Name
}

func (c *SimpleCommand) Execute() error {
	if c.Builtin != nil {
		err := c.Run()
		if err != nil {
			return err
		}
		return nil
	}
	cmd := exec.Command(c.Name, c.Args...)
	// Set File streams to shell default, to
	// allow direct manipulation from the terminal.
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func (c *SimpleCommand) Run() error {
	if c.Builtin == nil {
		return fmt.Errorf("Not a builtin command, check that your parser is working properly.")
	}
	return c.Builtin(c.Args)
}
