/*
Copyright © 2024 Alexander Vysochin <avyssochin@gmail.com>
*/
// Package cmd contains the CLI commands implemented using cobra.
package cmd

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
	"github.com/avysochin256/sox/pkg/sockopt"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:          "list <pid> <fd>",
	Short:        "List all socket options supported by sox",
	Args:         cobra.ExactArgs(2),
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		pid, err := strconv.Atoi(args[0])
		if err != nil {
			return fmt.Errorf("invalid pid %q: must be an integer", args[0])
		}
		fd, err := strconv.Atoi(args[1])
		if err != nil {
			return fmt.Errorf("invalid fd %q: must be an integer", args[1])
		}

		sockopt.ListSocketOptions(pid, fd, outputFormat)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}
