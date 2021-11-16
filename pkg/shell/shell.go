package shell

import (
	"bufio"
	"fmt"
	"io"
	"os"

	"github.com/raklaptudirm/mash/pkg/command"
)

type Runtime struct {
	Input     string
	Commands  []command.SimpleCommand
	SysWriter io.Writer
	SysReader io.Reader
	SysErr    io.Writer
}

// Start initiates the shell command loop, after initializing the
// reader. In the command loop, the prompt is printed, command is
// taken as input and then run using the Run method.
func Start() {
	reader := bufio.NewReader(os.Stdin)

	// Infinite command loop:
	for {
		cwd, err := os.Getwd()
		if err != nil {
			fmt.Fprintln(os.Stderr, "mash: error getting cwd")
			fmt.Fprintln(os.Stderr, err)
			break
		}

		// Prompt
		fmt.Printf("\u001b[32m%v\u001b[0m\nψ ", cwd)

		input, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				fmt.Println("^C")
				continue
			}

			fmt.Fprintln(os.Stderr, "mash: error reading prompt")
			fmt.Fprintln(os.Stderr, err)
			break
		}

		Run(input)
	}
}
