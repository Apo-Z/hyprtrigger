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

	if _, err := os.Stat(configDir); os.IsNotExist(err) {
		return nil
	}

	files, err := filepath.Glob(filepath.Join(configDir, "*.json"))
	if err != nil {
		return fmt.Errorf("failed to scan config directory: %w", err)
	}

	if len(files) == 0 {
		return nil
	}

	fmt.Printf("Auto-loading from: %s\n", configDir)

	loaded := 0
	for _, file := range files {
		if err := LoadEventsFromFile(file); err != nil {
			fmt.Printf("Failed to load %s: %v\n", filepath.Base(file), err)
			continue
		}
		loaded++
	}

	fmt.Printf("Auto-loaded %d file(s)\n", loaded)
	return nil
}
