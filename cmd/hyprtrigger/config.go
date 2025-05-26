package hyprtrigger

import (
	"fmt"
	"os"
	"path/filepath"
)

func getConfigDirectory() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	return filepath.Join(homeDir, ".config", "hyprtrigger")
}

func CreateConfig() error {
	configDir := getConfigDirectory()
	if configDir == "" {
		return fmt.Errorf("could not determine home directory")
	}

	if err := os.MkdirAll(configDir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory %s: %w", configDir, err)
	}

	fmt.Printf("Created config directory: %s\n", configDir)

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

	fmt.Printf("Created example configuration: %s\n", exampleFile)
	fmt.Println("Edit this file or create additional *.json files in the same directory")
	fmt.Println("Then run 'hyprtrigger' to use your configuration automatically")

	return nil
}
