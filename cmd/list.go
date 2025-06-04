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

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all socket options, supported by sox. Example: sox list <process pid> <socket fd>",
	Run: func(cmd *cobra.Command, args []string) {
		pid, err := strconv.Atoi(args[0])
		if err != nil {
			slog.Error("strconv.Atoi err", slog.Any("err", err))
		}
		fd, err := strconv.Atoi(args[1])
		if err != nil {
			slog.Error("strconv.Atoi err", slog.Any("err", err))
		}

		sockopt.ListSocketOptions(pid, fd, outputFormat)
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}
