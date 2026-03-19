package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"hyprtrigger/internal/daemon"
)

var reloadCmd = &cobra.Command{
	Use:   "reload",
	Short: "Reload configuration in running daemon",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := daemon.SendReload(); err != nil {
			return fmt.Errorf("reload failed: %w", err)
		}
		return nil
	},
}
