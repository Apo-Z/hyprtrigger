# Builtin Events

This directory contains the builtin event definitions for popular applications.

## Structure

Each application has its own file with a specific pattern:

- `[app].go` - Contains events for a specific application
- Each file exports a `Register[App]Events()` function
- Each file has an `init()` function that calls the registration function
- `builtin.go` - Orchestrates all builtin events for hot reload

## Files

- `bitwarden.go` - Bitwarden password manager events
- `blender.go` - Blender 3D modeling events
- `discord.go` - Discord messaging app events
- `steam.go` - Steam gaming platform events
- `system.go` - System tools (calculator, volume control, etc.)
- `terminal.go` - Terminal applications (kitty, alacritty, etc.)
- `vscode.go` - Visual Studio Code events
- `builtin.go` - Central registration for hot reload

## Adding New Applications

To add support for a new application:

1. Create `[app].go` file:
```go
package events

import "hyprtrigger/pkg/events"

// Register[App]Events registers all [App]-related events
func Register[App]Events() {
    events.RegisterEvent(&events.Event{
        Name:     "openwindow",
        Regex:    "app-name",
        Command:  "hyprctl dispatch workspace 1",
        UseShell: false,
    })
}

func init() {
    Register[App]Events()
}
```

2. Add the registration call to `builtin.go`:
```go
func RegisterBuiltinEvents() {
    // ... other calls ...
    Register[App]Events()
}
```

## Event Types

Supported Hyprland events:
- `openwindow` - New window opened
- `windowtitlev2` - Window title changed
- `activewindow` - Window focus changed

## Hot Reload Support

The `builtin.go` file provides `RegisterBuiltinEvents()` which:
- Can be called multiple times safely
- Re-registers all builtin events during configuration reload
- Maintains the modular structure while supporting hot reload

This allows the daemon to reload builtin events without restarting.
