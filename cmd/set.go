/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
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
			slog.Error("strconv.Atoi err:", err)
		}
		fd, err := strconv.Atoi(args[1])
		if err != nil {
			slog.Error("strconv.Atoi err:", err)
		}

		option := args[2]

		val, err := strconv.Atoi(args[3])
		if err != nil {
			slog.Error("strconv.Atoi err:", err)
		}

		sockopt.SetSocketOption(pid, fd, option, val)
	},
}

func init() {
	rootCmd.AddCommand(setCmd)
}
