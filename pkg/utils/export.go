package utils

import (
	"encoding/json"
	"fmt"
	"hyprtrigger/pkg/events"
	"os"
)

func ExportBuiltinEventsToJSON() (string, error) {
	builtinEvents := events.DefaultRegistry.GetBuiltinEvents()

	var allEvents []events.Event
	for _, eventList := range builtinEvents {
		for _, event := range eventList {
			allEvents = append(allEvents, *event)
		}
	}

	// Create inline config structure to avoid circular import
	eventConfig := struct {
		Events []events.Event `json:"events"`
	}{
		Events: allEvents,
	}

	jsonData, err := json.MarshalIndent(eventConfig, "", "  ")
	if err != nil {
		return "", fmt.Errorf("JSON serialization failed: %w", err)
	}

	return string(jsonData), nil
}

func SaveBuiltinEventsToFile(filename string) error {
	jsonData, err := ExportBuiltinEventsToJSON()
	if err != nil {
		return err
	}

	if err := os.WriteFile(filename, []byte(jsonData), 0644); err != nil {
		return fmt.Errorf("failed to write file %s: %w", filename, err)
	}

	fmt.Printf("Builtin events exported to: %s\n", filename)
	return nil
}

func PrintBuiltinEventsAsJSON() error {
	jsonData, err := ExportBuiltinEventsToJSON()
	if err != nil {
		return err
	}

	fmt.Println(jsonData)
	return nil
}
