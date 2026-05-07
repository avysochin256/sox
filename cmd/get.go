/*
Copyright © 2024 Alexander Vysochin <avyssochin@gmail.com>
*/
// Package cmd contains the CLI commands implemented using cobra.
package cmd

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
	"github.com/valexz/sox/pkg/sockopt"
)

// getCmd represents the get command
var getCmd = &cobra.Command{
	Use:          "get <pid> <fd> <option>",
	Short:        "Get a single socket option value",
	Args:         cobra.ExactArgs(3),
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

		option := args[2]
		if _, ok := sockopt.OptionsMap[option]; !ok {
			return fmt.Errorf("unsupported socket option %q (run `sox list <pid> <fd>` for the supported set)", option)
		}

		sockopt.GetSocketOption(pid, fd, option, outputFormat)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(getCmd)
}
