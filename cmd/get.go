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

var getCmd = &cobra.Command{
	Use:   "get <pid> <fd> <option>",
	Short: "Get the value of a specific socket option",
	Long: `Get the current value of a socket option for a given process and file descriptor.

Example:
  # Get TCP_NODELAY value for process 1234, socket fd 3
  sox get 1234 3 TCP_NODELAY

Available Options:
  - TCP_NODELAY       : Disable Nagle's algorithm
  - TCP_MAXSEG        : Maximum segment size
  - TCP_KEEPIDLE      : Time before sending keepalive probes
  - TCP_KEEPINTVL     : Time between keepalive probes
  - TCP_KEEPCNT       : Number of keepalive probes
  - SO_KEEPALIVE      : Enable TCP keepalive
  - SO_RCVBUF        : Socket receive buffer size
  - SO_SNDBUF        : Socket send buffer size
  And more... Use 'sox list' to see all available options`,
	Args: cobra.ExactArgs(3),
	RunE: func(cmd *cobra.Command, args []string) error {
		pid, err := strconv.Atoi(args[0])
		if err != nil {
			return fmt.Errorf("invalid process ID '%s': %w", args[0], err)
		}

		fd, err := strconv.Atoi(args[1])
		if err != nil {
			return fmt.Errorf("invalid file descriptor '%s': %w", args[1], err)
		}

		option := args[2]

		if err := sockopt.GetSocketOption(pid, fd, option); err != nil {
			return fmt.Errorf("failed to get socket option: %w", err)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(getCmd)
}
