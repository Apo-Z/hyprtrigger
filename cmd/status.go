package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"hyprtrigger/internal/daemon"
)

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show status of running daemon",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := daemon.SendStatus(); err != nil {
			return fmt.Errorf("status check failed: %w", err)
		}
		return nil
	},
}
