package config

import (
	"encoding/json"
	"fmt"
	"hyprtrigger/pkg/events"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

func LoadEventsFromFile(filename string) error {
	data, err := os.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("failed to read file %s: %w", filename, err)
	}

	var config EventConfig
	if err := json.Unmarshal(data, &config); err != nil {
		return fmt.Errorf("failed to parse JSON %s: %w", filename, err)
	}

	fmt.Printf("Loading %d event(s) from %s\n", len(config.Events), filename)

	for _, event := range config.Events {
		if event.Name == "" || event.Regex == "" || event.Command == "" {
			fmt.Printf("Invalid event ignored in %s\n", filename)
			continue
		}

		eventCopy := events.Event{
			Name:     event.Name,
			Regex:    event.Regex,
			Command:  event.Command,
			UseShell: event.UseShell,
		}

		events.RegisterEventExplicit(&eventCopy)
		fmt.Printf("Event loaded: %s -> %s\n", event.Name, event.Regex)
	}

	return nil
}

func LoadEventsFromDirectory(dirPath string) error {
	fmt.Printf("Scanning directory: %s\n", dirPath)

	var loadedFiles int

	err := filepath.WalkDir(dirPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() || !strings.HasSuffix(strings.ToLower(path), ".json") {
			return nil
		}

		fmt.Printf("Processing: %s\n", path)

		if err := LoadEventsFromFile(path); err != nil {
			fmt.Printf("Failed to load %s: %v\n", path, err)
			return nil
		}

		loadedFiles++
		return nil
	})

	if err != nil {
		return fmt.Errorf("failed to scan directory %s: %w", dirPath, err)
	}

	fmt.Printf("Summary: %d JSON file(s) processed\n", loadedFiles)
	return nil
}

func LoadEventsFromPath(path string) error {
	info, err := os.Stat(path)
	if err != nil {
		return fmt.Errorf("path does not exist: %s", path)
	}

	if info.IsDir() {
		return LoadEventsFromDirectory(path)
	} else {
		return LoadEventsFromFile(path)
	}
}
