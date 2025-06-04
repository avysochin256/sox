/*
Copyright © 2024 Alexander Vysochin <avyssochin@gmail.com>
*/
// Package cmd contains the CLI commands implemented using cobra.
package cmd

import (
	"github.com/spf13/cobra"
	"os"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "sox <command> <process pid> <socket fd> [<option name>] [<option val>]",
	Short: "SOX allows to get/update socket option value for any socket",
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
