package builtin

import "hyprtrigger/internal/events"

// Register adds all builtin events to the given registry.
// Safe to call multiple times; always call r.Clear() before re-registering on reload.
func Register(r *events.Registry) {
	registerBitwarden(r)
	registerBlender(r)
}
