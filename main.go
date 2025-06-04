/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
// Sox is a simple CLI utility for inspecting and modifying TCP socket options.
package main

import (
	"github.com/valexz/sox/cmd"
)

func main() {
	// Execute runs the root command and subcommands provided by the cmd package.
	cmd.Execute()
}

// TODO: LIST OPTIONS in HELP
// TODO: Value range validation for each option
