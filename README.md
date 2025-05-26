# HyprTrigger

Automatic window management for Hyprland based on real-time events. React to window creation, title changes, and focus events to automatically organize your workspace.

![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat&logo=go)
![License](https://img.shields.io/badge/license-MIT-blue.svg)

## Features

- **Event-driven automation** - React to Hyprland events in real-time
- **Hot reload configuration** - Update events without restarting the daemon
- **Daemon mode** - Single instance with command-line controls
- **Flexible configuration** - JSON-based event definitions
- **Builtin events** - Pre-configured automation for popular applications
- **Smart deduplication** - Prevents spam execution of the same command
- **Window ID support** - Target specific windows with `{WINDOW_ID}` placeholder
- **Shell command support** - Execute simple commands or complex shell scripts
- **Export/Import** - Export builtin events as templates

## Installation

### From source

```bash
go install github.com/Apo-Z/hyprtrigger@latest
```

### Manual build

```bash
git clone https://github.com/Apo-Z/hyprtrigger.git
cd hyprtrigger

# Build with version info
make build

# Or build manually
go build -o hyprtrigger
```

### Check version
```bash
hyprtrigger --version
```

## Quick Start

### Start the daemon (default behavior)
```bash
# Both commands are equivalent
hyprtrigger
hyprtrigger --daemon
```

### Setup automatic configuration
```bash
# Create config directory with example
hyprtrigger --create-config

# Edit the example file
nano ~/.config/hyprtrigger/example.json

# Reload configuration (daemon must be running)
hyprtrigger --reload
```

### Hot reload workflow
```bash
# 1. Start daemon
hyprtrigger --daemon

# 2. In another terminal, check status
hyprtrigger --status

# 3. Edit configuration files
nano ~/.config/hyprtrigger/my-events.json

# 4. Hot reload without stopping daemon
hyprtrigger --reload

# 5. Stop daemon when done
hyprtrigger --shutdown
```

### Makefile workflow
```bash
make daemon    # Start daemon
make reload    # Hot reload config
make status    # Check daemon status
make shutdown  # Stop daemon
```

## Configuration

### Auto-Configuration

HyprTrigger automatically loads all `*.json` files from `~/.config/hyprtrigger/`. This is the recommended way to manage your events.

```bash
# Setup auto-config directory
hyprtrigger --create-config

# This creates ~/.config/hyprtrigger/example.json
# Edit it or create additional JSON files in the same directory
# Then simply run:
hyprtrigger --daemon
```

### Hot Reload

The daemon supports hot reloading of configuration without interrupting event monitoring:

```bash
# Start daemon (first terminal)
hyprtrigger --daemon

# Reload configuration (second terminal)
hyprtrigger --reload

# Check status
hyprtrigger --status

# Shutdown daemon
hyprtrigger --shutdown
```

### Configuration Priority

Events are loaded in this order:
1. **Builtin events** (unless `--no-builtin`)
2. **Auto-config** from `~/.config/hyprtrigger/*.json` (unless `--no-auto-config`)
3. **Manual config** from `-c path` (if specified)

### Event Types

HyprTrigger supports these Hyprland events:

- `windowtitlev2` - Triggered when window title changes
- `openwindow` - Triggered when a new window opens
- `activewindow` - Triggered when window focus changes

### JSON Configuration Format

```json
{
  "events": [
    {
      "name": "windowtitlev2",
      "regex": "Bitwarden",
      "command": "hyprctl --batch \"dispatch setfloating address:0x{WINDOW_ID}; dispatch centerwindow\"",
      "use_shell": true
    },
    {
      "name": "openwindow",
      "regex": "discord",
      "command": "hyprctl dispatch workspace 2",
      "use_shell": false
    }
  ]
}
```

### Configuration Fields

- `name` - Hyprland event name to listen for
- `regex` - Regular expression to match against event data
- `command` - Command to execute when event matches
- `use_shell` - Whether to execute command through shell (`sh -c`)

### Window ID Placeholder

Use `{WINDOW_ID}` in commands to target the specific window that triggered the event:

```json
{
  "name": "windowtitlev2",
  "regex": "VSCode",
  "command": "hyprctl dispatch focuswindow address:0x{WINDOW_ID}",
  "use_shell": false
}
```

## Examples

### Automatic App Placement

```json
{
  "events": [
    {
      "name": "openwindow",
      "regex": "discord",
      "command": "hyprctl dispatch workspace 2",
      "use_shell": false
    },
    {
      "name": "openwindow",
      "regex": "steam",
      "command": "hyprctl --batch \"dispatch workspace 5; dispatch setfloating address:0x{WINDOW_ID}\"",
      "use_shell": true
    }
  ]
}
```

### Floating Window Management

```json
{
  "events": [
    {
      "name": "openwindow",
      "regex": "calculator|pavucontrol",
      "command": "hyprctl --batch \"dispatch setfloating address:0x{WINDOW_ID}; dispatch centerwindow\"",
      "use_shell": true
    }
  ]
}
```

### Terminal Workspace Assignment

```json
{
  "events": [
    {
      "name": "openwindow",
      "regex": "kitty|alacritty|wezterm",
      "command": "hyprctl dispatch moveworkspacetomonitor 3 current",
      "use_shell": false
    }
  ]
}
```

## Command Line Options

```
Usage: hyprtrigger [options]

Main options:
  -c, --config PATH       Path to JSON file or directory with events
  -n, --no-builtin        Disable builtin events
  -o, --config-only       Use only events from configuration file
  -s, --no-auto-config    Skip automatic loading from ~/.config/hyprtrigger/
  --create-config         Create ~/.config/hyprtrigger directory with example
  -e, --export-builtin    Export builtin events to JSON file
  -p, --print-builtin     Print builtin events as JSON
  -v, --version           Show version information
  -h, --help              Show help

Daemon control:
  -d, --daemon            Start in daemon mode (default behavior)
  -r, --reload            Reload configuration in running daemon
  --status                Show status of running daemon
  --shutdown              Shutdown running daemon
```

## Usage Modes

### 1. Daemon Mode (Default & Recommended)
```bash
# Start daemon (all equivalent)
hyprtrigger
hyprtrigger --daemon

# Control daemon
hyprtrigger --reload    # Hot reload config
hyprtrigger --status    # Check status
hyprtrigger --shutdown  # Stop daemon
```

### 2. Custom Configuration with Hot Reload
```bash
# Start with custom config
hyprtrigger -c my-events.json --daemon

# Hot reload after editing
hyprtrigger --reload
```

### 3. Pure Manual Configuration
```bash
hyprtrigger -c my-events.json --no-builtin --no-auto-config --daemon
```

### 4. Export Builtin Events
```bash
hyprtrigger --export-builtin ~/.config/hyprtrigger/template.json
```

### 5. Development & Testing
```bash
# View builtin events
hyprtrigger --print-builtin

# Create config directory
hyprtrigger --create-config
```

## Hot Reload Workflow

The hot reload feature allows you to update your configuration without interrupting the event monitoring:

### Setup
```bash
# 1. Start the daemon
hyprtrigger --daemon

# 2. Create/edit configuration
mkdir -p ~/.config/hyprtrigger
echo '{"events": [...]}' > ~/.config/hyprtrigger/my-app.json
```

### Development Cycle
```bash
# 3. Edit your event configurations
nano ~/.config/hyprtrigger/my-app.json

# 4. Hot reload (instant, no restart needed)
hyprtrigger --reload

# 5. Test your changes immediately
# 6. Repeat steps 3-5 as needed

# 7. Check daemon status anytime
hyprtrigger --status

# 8. Stop when done
hyprtrigger --shutdown
```

### Makefile Shortcuts
```bash
make daemon    # Start daemon
make reload    # Reload configuration
make status    # Check daemon status
make shutdown  # Stop daemon
```

## Daemon Management

### Check if daemon is running
```bash
hyprtrigger --status
```

### Multiple instance protection
If a daemon is already running, HyprTrigger will inform you:
```
Hyprtrigger daemon is already running!

Available daemon commands:
  hyprtrigger --reload     # Reload configuration
  hyprtrigger --status     # Check daemon status
  hyprtrigger --shutdown   # Stop daemon
```

### Force restart daemon
```bash
hyprtrigger --shutdown && hyprtrigger --daemon
```

## Builtin Events

HyprTrigger comes with several pre-configured events for popular applications. Each application has its own module for easy maintenance:

### Available Applications
- **Bitwarden** - Opens as floating window with custom size
- **Blender** - Preferences window as floating with custom size

### Builtin Events Structure
```
events/
├── bitwarden.go    # Bitwarden password manager events
├── blender.go      # Blender 3D modeling events
├── builtin.go      # Central registration for hot reload
└── README.md       # Documentation for adding new apps
```

### View all builtin events
```bash
hyprtrigger --print-builtin
```

### Export builtin events as template
```bash
hyprtrigger --export-builtin ~/.config/hyprtrigger/template.json
```

### Adding New Builtin Applications
```bash
# Use the generator script
make new-app APP=firefox

# This creates events/firefox.go with:
# - RegisterFirefoxEvents() function
# - Example event configurations
# - init() function for auto-loading
# - Instructions for integration
```

## Advanced Usage

### Directory Configuration
```bash
# Load all JSON files from a directory
hyprtrigger -c ./config/events/ --daemon
```

### Output Redirection
```bash
# Save builtin events to file
hyprtrigger --print-builtin > my-template.json

# Use with jq for filtering
hyprtrigger -p | jq '.events[] | select(.name == "windowtitlev2")'
```

### Daemon Management in Scripts
```bash
# Check if daemon is running
if hyprtrigger --status >/dev/null 2>&1; then
    echo "Daemon is running"
    hyprtrigger --reload
else
    echo "Starting daemon"
    hyprtrigger --daemon &
fi
```

### Systemd Integration

Create a systemd user service for automatic startup:

```bash
# Create service file
mkdir -p ~/.config/systemd/user
cat > ~/.config/systemd/user/hyprtrigger.service << EOF
[Unit]
Description=HyprTrigger Event Monitor
After=graphical-session.target

[Service]
Type=simple
ExecStart=%h/.local/bin/hyprtrigger --daemon
ExecStop=%h/.local/bin/hyprtrigger --shutdown
ExecReload=%h/.local/bin/hyprtrigger --reload
Restart=on-failure
RestartSec=5

[Install]
WantedBy=default.target
EOF

# Enable and start service
systemctl --user daemon-reload
systemctl --user enable hyprtrigger.service
systemctl --user start hyprtrigger.service

# Check status
systemctl --user status hyprtrigger.service

# Reload configuration via systemd
systemctl --user reload hyprtrigger.service
```

## How It Works

1. **Daemon Mode** - Single instance runs as background daemon
2. **Socket Communication** - Control commands sent via Unix domain socket
3. **Hot Reload** - Configuration reloaded without restarting daemon
4. **Event Monitoring** - Connects to Hyprland via Unix socket (`/.socket2.sock`)
5. **Real-time Processing** - Listens for events in real-time
6. **Pattern Matching** - Matches events against configured regex patterns
7. **Command Execution** - Executes commands when patterns match
8. **Deduplication** - Prevents spam execution with built-in 2-second window

## Troubleshooting

### Daemon Issues
```bash
# Check if daemon is running
hyprtrigger --status

# If daemon is unresponsive, find and kill process
ps aux | grep hyprtrigger
kill <pid>

# Clean up stale socket
rm -f ${XDG_RUNTIME_DIR:-/tmp}/hyprtrigger.sock

# Restart daemon
hyprtrigger --daemon
```

### Socket Connection Issues
```bash
# Check if Hyprland is running
ps aux | grep hyprland

# Check Hyprland socket exists
ls /tmp/hypr/$HYPRLAND_INSTANCE_SIGNATURE/.socket2.sock

# Verify HYPRLAND_INSTANCE_SIGNATURE is set
echo $HYPRLAND_INSTANCE_SIGNATURE
```

### Event Not Triggering
```bash
# Test your regex patterns
hyprtrigger --print-builtin | jq '.events[] | select(.regex | test("your-pattern"))'

# Monitor raw events with empty config
hyprtrigger -c /dev/null --no-builtin --no-auto-config --daemon

# Check if events are loaded
hyprtrigger --status
```

### Hot Reload Not Working
```bash
# Check daemon status first
hyprtrigger --status

# If daemon is not running, start it
hyprtrigger --daemon &

# Then try reload again
hyprtrigger --reload

# Check reload succeeded
hyprtrigger --status
```

### Command Not Executing
- Verify the command works manually: `hyprctl dispatch workspace 2`
- Check if `use_shell` is needed for complex commands with pipes/quotes
- Ensure proper escaping of quotes in JSON
- Test with simple commands first

### Configuration Issues
```bash
# Validate JSON syntax
hyprtrigger --print-builtin | jq '.'

# Check auto-config directory
ls -la ~/.config/hyprtrigger/

# Test minimal config
echo '{"events":[]}' > /tmp/test.json
hyprtrigger -c /tmp/test.json --no-builtin --no-auto-config --daemon
```

## Development

### Building

```bash
# Build with version information
make build

# Install system-wide
make install

# Development workflow
make daemon        # Start daemon
make reload        # Reload config
make status        # Check status
make shutdown      # Stop daemon

# Development with source
make dev-daemon    # Start daemon using go run
make dev-reload    # Reload using go run
make dev-status    # Status using go run
make dev-shutdown  # Shutdown using go run

# Cross-platform builds
make build-all     # Build for Linux and macOS (amd64/arm64)
```

### Adding New Builtin Applications

```bash
# 1. Generate new application template
make new-app APP=spotify

# 2. Edit the generated file
nano events/spotify.go

# 3. Add your events
func RegisterSpotifyEvents() {
    events.RegisterEvent(&events.Event{
        Name:     "openwindow",
        Regex:    "spotify",
        Command:  "hyprctl dispatch workspace 4",
        UseShell: false,
    })
}

# 4. Add to builtin registration (as instructed by the script)
nano events/builtin.go

# 5. Test your events
make dev-test      # Show all builtin events
make dev-daemon    # Start daemon to test
```

### Hot Reload Development Cycle

```bash
# Terminal 1: Start development daemon
make dev-daemon

# Terminal 2: Development cycle
make dev-reload    # Test reload functionality
make dev-status    # Check daemon status
make dev-shutdown  # Stop when done
```

### Version Information

Version information is embedded during build:

```bash
# Show version
hyprtrigger --version
# Output: hyprtrigger v1.0.0 (commit: abc1234, built: 2024-01-15T10:30:00Z)

# Makefile automatically detects:
# - Version from git tags
# - Commit hash from git
# - Build timestamp
```

### Testing

```bash
# Run tests
make test

# Run linter
make lint

# Test configuration loading
make dev-config    # Create test config
make dev-daemon    # Start with test config
make dev-reload    # Test hot reload
```

## Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

### Contributing New Builtin Events

1. Use the generator: `make new-app APP=myapp`
2. Edit the generated `events/myapp.go` file
3. Add registration to `events/builtin.go`
4. Test with `make dev-test` and `make dev-daemon`
5. Update documentation if needed
6. Submit a pull request

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Acknowledgments

- [Hyprland](https://hyprland.org/) - The amazing Wayland compositor
- The Hyprland community for event documentation and examples
