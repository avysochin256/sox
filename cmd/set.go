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

var setCmd = &cobra.Command{
	Use:   "set <pid> <fd> <option> <value>",
	Short: "Set the value of a specific socket option",
	Long: `Set the value of a socket option for a given process and file descriptor.

Example:
  # Enable TCP_NODELAY for process 1234, socket fd 3
  sox set 1234 3 TCP_NODELAY 1

  # Set TCP keepalive idle time to 60 seconds
  sox set 1234 3 TCP_KEEPIDLE 60

Available Options:
  - TCP_NODELAY       : Disable Nagle's algorithm (0/1)
  - TCP_MAXSEG        : Maximum segment size (536-65535)
  - TCP_KEEPIDLE      : Time before sending keepalive probes (1-32767 seconds)
  - TCP_KEEPINTVL     : Time between keepalive probes (1-32767 seconds)
  - TCP_KEEPCNT       : Number of keepalive probes (1-127)
  - SO_KEEPALIVE      : Enable TCP keepalive (0/1)
  - SO_RCVBUF        : Socket receive buffer size (>= 2048 bytes)
  - SO_SNDBUF        : Socket send buffer size (>= 2048 bytes)
  And more... Use 'sox list' to see all available options`,
	Args: cobra.ExactArgs(4),
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

		value, err := strconv.Atoi(args[3])
		if err != nil {
			return fmt.Errorf("invalid value '%s': %w", args[3], err)
		}

		if err := sockopt.SetSocketOption(pid, fd, option, value); err != nil {
			return fmt.Errorf("failed to set socket option: %w", err)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(setCmd)
}
