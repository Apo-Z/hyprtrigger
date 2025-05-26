package events

// RegisterBuiltinEvents registers all builtin events
// This function calls individual registration functions from each module
// It can be called multiple times safely (e.g., during reload)
func RegisterBuiltinEvents() {
	// Register all builtin events by calling individual registration functions
	RegisterBitwardenEvents()
	RegisterBlenderEvents()
	// RegisterDiscordEvents()
	// RegisterVSCodeEvents()
	// RegisterSteamEvents()
	// RegisterTerminalEvents()
	// RegisterSystemToolEvents()
}
