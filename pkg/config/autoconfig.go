package config

import (
	"fmt"
	"os"
	"path/filepath"
)

func GetConfigDirectory() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	return filepath.Join(homeDir, ".config", "hyprtrigger")
}

func LoadAutoConfig() error {
	configDir := GetConfigDirectory()

	// Check if config directory exists
	if _, err := os.Stat(configDir); os.IsNotExist(err) {
		return nil // No error if directory doesn't exist
	}

	fmt.Printf("Loading auto-config from: %s\n", configDir)

	files, err := filepath.Glob(filepath.Join(configDir, "*.json"))
	if err != nil {
		return fmt.Errorf("failed to scan config directory: %w", err)
	}

	if len(files) == 0 {
		fmt.Printf("No JSON files found in %s\n", configDir)
		return nil
	}

	loadedCount := 0
	for _, file := range files {
		fmt.Printf("Auto-loading: %s\n", filepath.Base(file))
		if err := LoadEventsFromFile(file); err != nil {
			fmt.Printf("Failed to load %s: %v\n", filepath.Base(file), err)
			continue
		}
		loadedCount++
	}

	fmt.Printf("Auto-loaded %d configuration file(s)\n", loadedCount)
	return nil
}
