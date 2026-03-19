# Builtin Events

Builtin events are defined in `internal/builtin/`. Each app has its own file with a private `register<App>` function.

## Adding a new app

1. Create `internal/builtin/<app>.go`:

```go
package builtin

import "hyprtrigger/internal/events"

func registerMyApp(r *events.Registry) {
    r.RegisterBuiltin(&events.Event{
        Name:     "openwindow",
        Regex:    "myapp",
        Command:  "hyprctl dispatch workspace 2",
        UseShell: false,
    })
}
```

2. Call it from `internal/builtin/register.go`:

```go
func Register(r *events.Registry) {
    // ...existing calls...
    registerMyApp(r)
}
```

No `init()` functions — registration is always explicit via `Register(r)`.

## Supported event names

- `openwindow` — new window opened
- `windowtitlev2` — window title changed
- `activewindow` — window focus changed

Use `{WINDOW_ID}` in commands to target the specific window.
