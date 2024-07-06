/*
Copyright © 2024 Alexander Vysochin <avyssochin@gmail.com>
*/
package cmd

import (
	"github.com/spf13/cobra"
	"strconv"
	"log/slog"
	"github.com/valexz/sox/pkg/sockopt"
)

// getCmd represents the get command
var getCmd = &cobra.Command{
	Use:   "get",
	Short: "Get single parameter of socket. Example: sox get <process pid> <socket fd> <socket option name>",
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
	

		sockopt.GetSocketOption(pid, fd, option)
	},
}

func init() {
	rootCmd.AddCommand(getCmd)
}
