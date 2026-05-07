/*
Copyright © 2024 Alexander Vysochin <avyssochin@gmail.com>
*/
// Package cmd contains the CLI commands implemented using cobra.
package cmd

import (
	"github.com/spf13/cobra"
	"github.com/avysochin256/sox/pkg/sockopt"
)

// explainCmd represents the explain command
var explainCmd = &cobra.Command{
	Use:          "explain <option>",
	Short:        "Explain a socket option in detail (level, range, summary, details)",
	Args:         cobra.ExactArgs(1),
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		return sockopt.ExplainSocketOption(args[0], outputFormat)
	},
}

func init() {
	rootCmd.AddCommand(explainCmd)
}
