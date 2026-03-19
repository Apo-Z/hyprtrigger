package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"hyprtrigger/internal/config"
)

var initConfigCmd = &cobra.Command{
	Use:   "init-config",
	Short: "Create ~/.config/hyprtrigger/ with an example config",
	RunE: func(cmd *cobra.Command, args []string) error {
		configDir := config.GetConfigDirectory()
		if configDir == "" {
			return fmt.Errorf("could not determine home directory")
		}

		if err := os.MkdirAll(configDir, 0755); err != nil {
			return fmt.Errorf("failed to create config directory: %w", err)
		}
		fmt.Printf("Config directory: %s\n", configDir)

		exampleFile := filepath.Join(configDir, "example.json")
		exampleConfig := `{
  "events": [
    {
      "name": "windowtitlev2",
      "regex": "Firefox",
      "command": "hyprctl dispatch workspace 1",
      "use_shell": false
    },
    {
      "name": "openwindow",
      "regex": "calculator",
      "command": "hyprctl --batch \"dispatch setfloating address:0x{WINDOW_ID}; dispatch centerwindow\"",
      "use_shell": true
    }
  ]
}`
		if err := os.WriteFile(exampleFile, []byte(exampleConfig), 0644); err != nil {
			return fmt.Errorf("failed to create example file: %w", err)
		}

		fmt.Printf("Example config: %s\n", exampleFile)
		fmt.Println("Edit this file or add more *.json files in the same directory.")
		fmt.Println("Then run 'hyprtrigger' to start.")
		return nil
	},
}
