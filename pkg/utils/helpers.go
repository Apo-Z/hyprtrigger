package utils

import (
	"fmt"
	"hyprtrigger/pkg/events"
)

func truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}

func PrintEventsSummary() {
	allEvents := events.GetAllEvents()
	totalEvents := 0

	fmt.Println("\nLoaded events summary:")

	if len(allEvents) == 0 {
		fmt.Println("  No events loaded!")
		if events.DefaultRegistry.GetBuiltinEvents() != nil {
			fmt.Println("  Builtin events are disabled")
		}
		fmt.Println()
		return
	}

	for eventName, eventList := range allEvents {
		fmt.Printf("  %s: %d event(s)\n", eventName, len(eventList))
		for _, event := range eventList {
			fmt.Printf("    - Regex: %s | Command: %s | Shell: %t\n",
				event.Regex,
				truncateString(event.Command, 50),
				event.UseShell)
		}
		totalEvents += len(eventList)
	}

	fmt.Printf("Total: %d event(s) registered", totalEvents)
	fmt.Println("\n")
}
