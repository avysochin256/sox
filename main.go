/*
Copyright © 2024 Alexander Vysochin <avyssochin@gmail.com>
Licensed under the Apache License, Version 2.0

SOX is a command-line tool that allows you to view and modify TCP socket options
for any process. It's particularly useful for debugging network issues and tuning
network performance without requiring application restarts.
*/
package main

import (
	"github.com/valexz/sox/cmd"
)

func main() {
	cmd.Execute()
}

// TODO: LIST OPTIONS in HELP
// TODO: Value range validation for each option
