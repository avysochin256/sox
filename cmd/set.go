/*
Copyright © 2024 Alexander Vysochin <avyssochin@gmail.com>
*/
// Package cmd contains the CLI commands implemented using cobra.
package cmd

import (
	"github.com/valexz/sox/pkg/sockopt"
	"log/slog"
	"strconv"

	"github.com/spf13/cobra"
)

// setCmd represents the set command
var setCmd = &cobra.Command{
	Use:   "set",
	Short: "Set value for single socket option. Example: sox set <process pid> <socket fd> <socket option name> <option value>",
	Run: func(cmd *cobra.Command, args []string) {
		pid, err := strconv.Atoi(args[0])
		if err != nil {
			slog.Error("strconv.Atoi err", slog.Any("err", err))
		}
		fd, err := strconv.Atoi(args[1])
		if err != nil {
			slog.Error("strconv.Atoi err", slog.Any("err", err))
		}

		option := args[2]

		val, err := strconv.Atoi(args[3])
		if err != nil {
			slog.Error("strconv.Atoi err", slog.Any("err", err))
		}

		sockopt.SetSocketOption(pid, fd, option, val, outputFormat)
	},
}

func init() {
	rootCmd.AddCommand(setCmd)
}
