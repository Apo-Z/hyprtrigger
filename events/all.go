package events

// This file ensures all builtin events are properly loaded by calling their init() functions.
// Each individual event file (bitwarden.go, discord.go, etc.) has an init() function
// that registers its events. By having this file in the same package, we ensure
// all init() functions are called when the events package is imported.

// No additional code needed here - the mere presence of this file in the package
// ensures that when "hyprtrigger/events" is imported, all .go files in this
// directory will have their init() functions executed.

func init() {
	// This init function runs after all other init functions in the package
	// We can use it to verify that events were loaded correctly
	// (Optional: could add debug logging here if needed)
}
