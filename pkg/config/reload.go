package config

import (
	"fmt"
	"hyprtrigger/events"
	eventsPkg "hyprtrigger/pkg/events"
)

func ReloadConfiguration(skipBuiltinEvents, skipAutoConfig bool, configPath string) error {
	fmt.Println("🔄 Reloading configuration...")

	// Clear current events registry
	eventsPkg.ClearRegistry()

	// Re-register builtin events if needed
	if !skipBuiltinEvents {
		events.RegisterBuiltinEvents()
		fmt.Println("Builtin events reloaded")
	} else {
		fmt.Println("Builtin events disabled")
		eventsPkg.SetSkipBuiltinEvents(true)
	}

	// Load auto-config first (unless disabled)
	if !skipAutoConfig {
		if err := LoadAutoConfig(); err != nil {
			fmt.Printf("Auto-config loading failed: %v\n", err)
		}
	} else {
		fmt.Println("Auto-config disabled")
	}

	// Then load manual config (if specified)
	if configPath != "" {
		fmt.Printf("Loading events from: %s\n", configPath)
		if err := LoadEventsFromPath(configPath); err != nil {
			return fmt.Errorf("config loading failed: %w", err)
		}
	}

	// Check if any events are loaded
	totalEvents := 0
	allEvents := eventsPkg.GetAllEvents()
	for _, eventList := range allEvents {
		totalEvents += len(eventList)
	}

	if totalEvents == 0 {
		return fmt.Errorf("no events loaded after reload")
	}

	fmt.Printf("✅ Configuration reloaded successfully (%d events)\n", totalEvents)
	return nil
}
