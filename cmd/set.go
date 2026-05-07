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

// setCmd represents the set command
var setCmd = &cobra.Command{
	Use:          "set <pid> <fd> <option> <value>",
	Short:        "Set the value of a single socket option",
	Args:         cobra.ExactArgs(4),
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

		val, err := strconv.Atoi(args[3])
		if err != nil {
			return fmt.Errorf("invalid value %q: must be an integer", args[3])
		}

		sockopt.SetSocketOption(pid, fd, option, val, outputFormat)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(setCmd)
}
