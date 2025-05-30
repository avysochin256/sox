/*
Copyright © 2024 Alexander Vysochin <avyssochin@gmail.com>
Licensed under the Apache License, Version 2.0
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "sox",
	Short: "SOX - Socket Options eXtractor/setter",
	Long: `SOX is a command-line tool that allows you to view and modify TCP socket options
for any process. It's particularly useful for debugging network issues and tuning
network performance without requiring application restarts.

Examples:
  # List all socket options for process 1234, socket fd 3
  sox list 1234 3

  # Get specific socket option
  sox get 1234 3 TCP_NODELAY

  # Set socket option value
  sox set 1234 3 TCP_NODELAY 1

Requirements:
  - Linux kernel 5.6+
  - CAP_SYS_PTRACE capability (run with sudo)`,
	Version: "1.0.0",
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() {
	if os.Geteuid() != 0 {
		fmt.Println("Error: This program requires root privileges to access socket options.")
		fmt.Println("Please run with sudo or as root.")
		os.Exit(1)
	}

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	// Add persistent flags that will be available to all subcommands
	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "Enable verbose output")
}
