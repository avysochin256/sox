/*
Copyright © 2024 Alexander Vysochin <avyssochin@gmail.com>
Licensed under the Apache License, Version 2.0
*/
package cmd

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
	"github.com/valexz/sox/pkg/sockopt"
)

var listCmd = &cobra.Command{
	Use:   "list <pid> <fd>",
	Short: "List all available socket options and their values",
	Long: `List all available socket options and their current values for a given process and file descriptor.

Example:
  # List all socket options for process 1234, socket fd 3
  sox list 1234 3

The output shows:
  - Option name
  - Current value
  - Description
  - Allowed value range (where applicable)

Note: Some options might not be available depending on your kernel version and socket type.`,
	Args: cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		pid, err := strconv.Atoi(args[0])
		if err != nil {
			return fmt.Errorf("invalid process ID '%s': %w", args[0], err)
		}

		fd, err := strconv.Atoi(args[1])
		if err != nil {
			return fmt.Errorf("invalid file descriptor '%s': %w", args[1], err)
		}

		if err := sockopt.ListSocketOptions(pid, fd); err != nil {
			return fmt.Errorf("failed to list socket options: %w", err)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}
