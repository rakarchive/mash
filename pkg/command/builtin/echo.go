// mash
// https://github.com/raklaptudirm/mash
// Copyright (c) 2021 Rak Laptudirm.
// Licensed under the MIT license.

package builtin

import "fmt"

// Function echo implements the echo command which is
// used to print output to the console.
func echo(args []string) error {
	// Print all the args seperated by a space
	for i, str := range args {
		if i != 0 {
			fmt.Print(" ")
		}
		fmt.Print(str)
	}
	fmt.Print("\n")
	return nil
}
