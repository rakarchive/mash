package command

import (
	"fmt"
	"os"
	"os/exec"
)

type Command struct {
	Name    string
	Args    []string
	Builtin Builtin
}

type Builtin func([]string) error

func (c *Command) String() string {
	return c.Name
}

func (c *Command) Execute() error {
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

func (c *Command) Run() error {
	if c.Builtin == nil {
		return fmt.Errorf("Not a builtin command, check that your parser is working properly.")
	}
	return c.Builtin(c.Args)
}
