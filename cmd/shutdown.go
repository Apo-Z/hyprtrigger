package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"hyprtrigger/internal/daemon"
)

var shutdownCmd = &cobra.Command{
	Use:   "shutdown",
	Short: "Stop the running daemon",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := daemon.SendShutdown(); err != nil {
			return fmt.Errorf("shutdown failed: %w", err)
		}
		return nil
	},
}
