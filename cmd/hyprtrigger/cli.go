package hyprtrigger

import (
	"flag"
	"fmt"
	"hyprtrigger/pkg/config"
	"hyprtrigger/pkg/daemon"
	"hyprtrigger/pkg/events"
	"hyprtrigger/pkg/socket"
	"hyprtrigger/pkg/utils"
	"os"
	"os/signal"
	"syscall"
	_ "hyprtrigger/events"
)

func printError(message string) {
	fmt.Fprintf(os.Stderr, "Error: %s\n", message)
	fmt.Fprintf(os.Stderr, "Use 'hyprtrigger --help' to see available options\n")
	os.Exit(1)
}

func Execute() {
	// Customize flag error output
	flag.CommandLine.Usage = func() {
		printError("invalid flag")
	}

	var configPath string
	flag.StringVar(&configPath, "config", "", "Path to JSON file or directory containing JSON event files")
	flag.StringVar(&configPath, "c", "", "Path to JSON file or directory (shortcut)")

	var showHelp bool
	flag.BoolVar(&showHelp, "help", false, "Show help")
	flag.BoolVar(&showHelp, "h", false, "Show help (shortcut)")

	var showVersion bool
	flag.BoolVar(&showVersion, "version", false, "Show version information")
	flag.BoolVar(&showVersion, "v", false, "Show version information (shortcut)")

	var skipBuiltinEvents bool
	flag.BoolVar(&skipBuiltinEvents, "no-builtin", false, "Disable builtin events")
	flag.BoolVar(&skipBuiltinEvents, "n", false, "Disable builtin events (shortcut)")

	var configOnly bool
	flag.BoolVar(&configOnly, "config-only", false, "Use only events from configuration file")
	flag.BoolVar(&configOnly, "o", false, "Use only events from configuration file (shortcut)")

	var skipAutoConfig bool
	flag.BoolVar(&skipAutoConfig, "no-auto-config", false, "Skip automatic loading from ~/.config/hyprtrigger/")
	flag.BoolVar(&skipAutoConfig, "s", false, "Skip automatic loading from ~/.config/hyprtrigger/ (shortcut)")

	var createConfig bool
	flag.BoolVar(&createConfig, "create-config", false, "Create ~/.config/hyprtrigger directory and example file")

	var exportBuiltin string
	flag.StringVar(&exportBuiltin, "export-builtin", "", "Export builtin events to JSON file")
	flag.StringVar(&exportBuiltin, "e", "", "Export builtin events to JSON file (shortcut)")

	var printBuiltinJSON bool
	flag.BoolVar(&printBuiltinJSON, "print-builtin", false, "Print builtin events as JSON")
	flag.BoolVar(&printBuiltinJSON, "p", false, "Print builtin events as JSON (shortcut)")

	var reloadConfig bool
	flag.BoolVar(&reloadConfig, "reload", false, "Reload configuration in running instance")
	flag.BoolVar(&reloadConfig, "r", false, "Reload configuration in running instance (shortcut)")

	var showStatus bool
	flag.BoolVar(&showStatus, "status", false, "Show status of running instance")

	var shutdownDaemon bool
	flag.BoolVar(&shutdownDaemon, "shutdown", false, "Shutdown running instance")

	var startDaemon bool
	flag.BoolVar(&startDaemon, "daemon", false, "Start in daemon mode")
	flag.BoolVar(&startDaemon, "d", false, "Start in daemon mode (shortcut)")

	// Parse flags and capture errors
	err := flag.CommandLine.Parse(os.Args[1:])
	if err != nil {
		os.Exit(1)
	}

	// Check for unknown positional arguments
	if flag.NArg() > 0 {
		arg := flag.Arg(0)

		// Specific error messages for common mistakes
		switch arg {
		case "version":
			printError("use '--version' or '-v' to show version")
		case "help":
			printError("use '--help' or '-h' to show help")
		case "config":
			printError("use '--config PATH' or '-c PATH' to specify config file")
		case "reload":
			printError("use '--reload' or '-r' to reload configuration")
		case "status":
			printError("use '--status' to show daemon status")
		case "shutdown":
			printError("use '--shutdown' to shutdown daemon")
		case "daemon":
			printError("use '--daemon' or '-d' to start daemon")
		case "start":
			printError("use '--daemon' or '-d' to start daemon")
		default:
			printError(fmt.Sprintf("unknown argument '%s'", arg))
		}
	}

	// Handle special commands first
	if showHelp {
		printHelp()
		return
	}

	if showVersion {
		PrintVersion()
		return
	}

	if createConfig {
		if err := CreateConfig(); err != nil {
			fmt.Printf("Failed to create config directory: %v\n", err)
			os.Exit(1)
		}
		return
	}

	if exportBuiltin != "" {
		fmt.Println("Exporting builtin events...")
		if err := utils.SaveBuiltinEventsToFile(exportBuiltin); err != nil {
			fmt.Printf("Export failed: %v\n", err)
			os.Exit(1)
		}
		return
	}

	if printBuiltinJSON {
		if err := utils.PrintBuiltinEventsAsJSON(); err != nil {
			fmt.Printf("JSON print failed: %v\n", err)
			os.Exit(1)
		}
		return
	}

	// Handle daemon control commands
	if reloadConfig {
		fmt.Println("Sending reload command to hyprtrigger daemon...")
		if err := daemon.SendReload(); err != nil {
			fmt.Printf("Reload failed: %v\n", err)
			os.Exit(1)
		}
		return
	}

	if showStatus {
		fmt.Println("Checking hyprtrigger daemon status...")
		if err := daemon.SendStatus(); err != nil {
			fmt.Printf("Status check failed: %v\n", err)
			os.Exit(1)
		}
		return
	}

	if shutdownDaemon {
		fmt.Println("Sending shutdown command to hyprtrigger daemon...")
		if err := daemon.SendShutdown(); err != nil {
			fmt.Printf("Shutdown failed: %v\n", err)
			os.Exit(1)
		}
		return
	}

	// Check if daemon is already running
	if daemon.IsDaemonRunning() {
		if startDaemon {
			fmt.Println("Hyprtrigger daemon is already running!")
			fmt.Println("Use 'hyprtrigger --reload' to reload configuration")
			fmt.Println("Use 'hyprtrigger --status' to check status")
			fmt.Println("Use 'hyprtrigger --shutdown' to stop the daemon")
			os.Exit(1)
		} else {
			// If daemon is running and user didn't specify --daemon,
			// inform them about available daemon commands
			fmt.Println("Hyprtrigger daemon is already running!")
			fmt.Println("")
			fmt.Println("Available daemon commands:")
			fmt.Println("  hyprtrigger --reload     # Reload configuration")
			fmt.Println("  hyprtrigger --status     # Check daemon status")
			fmt.Println("  hyprtrigger --shutdown   # Stop daemon")
			fmt.Println("")
			fmt.Println("To start a new daemon instance (not recommended):")
			fmt.Println("  hyprtrigger --shutdown && hyprtrigger --daemon")
			os.Exit(1)
		}
	}

	// Determine if we should start in daemon mode
	shouldStartDaemon := startDaemon || len(os.Args) == 1 // Default behavior or explicit --daemon

	// Start main application
	if shouldStartDaemon {
		fmt.Println("Starting Hyprland event monitor daemon")
	} else {
		fmt.Println("Starting Hyprland event monitor (one-shot mode)")
	}

	// Store configuration for reload capability
	configFlags := struct {
		skipBuiltinEvents bool
		configOnly        bool
		skipAutoConfig    bool
		configPath        string
	}{
		skipBuiltinEvents: skipBuiltinEvents || configOnly,
		configOnly:        configOnly,
		skipAutoConfig:    skipAutoConfig,
		configPath:        configPath,
	}

	// Initial configuration load
	if err := loadConfiguration(configFlags); err != nil {
		fmt.Printf("Configuration loading failed: %v\n", err)
		os.Exit(1)
	}

	// Print summary of loaded events
	utils.PrintEventsSummary()

	// Check if any events are loaded
	if len(events.GetAllEvents()) == 0 {
		fmt.Println("No events loaded, stopping")
		fmt.Println("Use -h to see help")
		os.Exit(1)
	}

	// Start daemon for reload functionality (only in daemon mode)
	var daemonServer *daemon.Daemon
	if shouldStartDaemon {
		daemonServer = daemon.NewDaemon()
		if err := daemonServer.Start(); err != nil {
			fmt.Printf("Failed to start daemon: %v\n", err)
			os.Exit(1)
		}
		defer daemonServer.Stop()
	}

	// Connect to Hyprland socket
	client := socket.NewClient()
	if err := client.Connect(); err != nil {
		fmt.Printf("%v\n", err)
		fmt.Println("Make sure Hyprland is running")
		os.Exit(1)
	}
	defer client.Close()

	// Setup signal handling
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Start listening for events in a goroutine
	listener := socket.NewListener(client)
	listenerDone := make(chan error, 1)

	go func() {
		listenerDone <- listener.Listen()
	}()

	// Main event loop
	if shouldStartDaemon {
		// Daemon mode with hot reload support
		for {
			select {
			case <-daemonServer.GetReloadChannel():
				fmt.Println("🔄 Reloading configuration...")
				if err := config.ReloadConfiguration(
					configFlags.skipBuiltinEvents,
					configFlags.skipAutoConfig,
					configFlags.configPath,
				); err != nil {
					fmt.Printf("Reload failed: %v\n", err)
				} else {
					utils.PrintEventsSummary()
				}

			case <-daemonServer.GetShutdownChannel():
				fmt.Println("Shutdown requested via daemon")
				return

			case sig := <-sigChan:
				fmt.Printf("\nReceived signal: %v\n", sig)
				fmt.Println("Shutting down gracefully...")
				return

			case err := <-listenerDone:
				if err != nil {
					fmt.Printf("Listener error: %v\n", err)
				}
				return
			}
		}
	} else {
		// One-shot mode - just run until interrupted
		select {
		case sig := <-sigChan:
			fmt.Printf("\nReceived signal: %v\n", sig)
			fmt.Println("Shutting down...")
			return

		case err := <-listenerDone:
			if err != nil {
				fmt.Printf("Listener error: %v\n", err)
			}
			return
		}
	}
}

func loadConfiguration(configFlags struct {
	skipBuiltinEvents bool
	configOnly        bool
	skipAutoConfig    bool
	configPath        string
}) error {
	// Configure builtin events
	if configFlags.skipBuiltinEvents || configFlags.configOnly {
		fmt.Println("Builtin events disabled")
		events.SetSkipBuiltinEvents(true)
	}

	// Load auto-config first (unless disabled)
	if !configFlags.skipAutoConfig {
		if err := config.LoadAutoConfig(); err != nil {
			fmt.Printf("Auto-config loading failed: %v\n", err)
		}
	}

	// Then load manual config (if specified)
	if configFlags.configPath != "" {
		fmt.Printf("Loading events from: %s\n", configFlags.configPath)
		if err := config.LoadEventsFromPath(configFlags.configPath); err != nil {
			return fmt.Errorf("config loading failed: %w", err)
		}
	} else {
		if configFlags.skipBuiltinEvents || configFlags.configOnly {
			if configFlags.skipAutoConfig {
				fmt.Println("No configuration file specified, auto-config disabled, and builtin events disabled")
			} else {
				fmt.Println("No configuration file specified and builtin events disabled")
			}
			fmt.Println("Use -c to specify a configuration file")
		} else {
			if !configFlags.skipAutoConfig {
				fmt.Println("Using builtin events + auto-config from ~/.config/hyprtrigger/")
			} else {
				fmt.Println("No configuration file specified, using builtin events only")
			}
		}
	}

	return nil
}

func printHelp() {
	fmt.Println("Hyprland Event Monitor")
	fmt.Println("")
	fmt.Println("Usage:")
	fmt.Println("  hyprtrigger [options]")
	fmt.Println("")
	fmt.Println("Options:")
	fmt.Println("  -c, --config PATH       Path to JSON file or directory with events")
	fmt.Println("  -n, --no-builtin        Disable builtin events")
	fmt.Println("  -o, --config-only       Use only events from configuration file")
	fmt.Println("  --create-config         Create ~/.config/hyprtrigger directory with example")
	fmt.Println("  -s, --no-auto-config    Skip automatic loading from ~/.config/hyprtrigger/")
	fmt.Println("  -e, --export-builtin    Export builtin events to JSON file")
	fmt.Println("  -p, --print-builtin     Print builtin events as JSON")
	fmt.Println("  -r, --reload            Reload configuration in running instance")
	fmt.Println("  --status                Show status of running instance")
	fmt.Println("  --shutdown              Shutdown running instance")
	fmt.Println("  -v, --version           Show version information")
	fmt.Println("  -h, --help              Show this help")
	fmt.Println("")
	fmt.Println("Examples:")
	fmt.Println("  hyprtrigger                                      # Start daemon (default)")
	fmt.Println("  hyprtrigger --daemon                             # Start daemon (explicit)")
	fmt.Println("  hyprtrigger -c events.json                      # Start daemon with custom config")
	fmt.Println("  hyprtrigger --reload                            # Reload configuration in running daemon")
	fmt.Println("  hyprtrigger --status                            # Check if daemon is running")
	fmt.Println("  hyprtrigger --shutdown                          # Stop running daemon")
	fmt.Println("  hyprtrigger -c events.json --no-builtin         # Auto-config + manual config only")
	fmt.Println("  hyprtrigger --no-auto-config                    # Builtin events only")
	fmt.Println("  hyprtrigger -s -c events.json                   # Builtin + manual config only")
	fmt.Println("  hyprtrigger -o -c events.json                   # Manual config only")
	fmt.Println("")
	fmt.Println("Daemon control:")
	fmt.Println("  The first hyprtrigger instance runs as a daemon and listens for reload commands.")
	fmt.Println("  Subsequent calls to 'hyprtrigger --reload' will refresh the configuration")
	fmt.Println("  without restarting the daemon or interrupting event monitoring.")
	fmt.Println("")
	fmt.Println("Auto-configuration:")
	fmt.Println("  hyprtrigger automatically loads *.json files from ~/.config/hyprtrigger/")
	fmt.Println("  Use --no-auto-config to disable this behavior")
	fmt.Println("  mkdir -p ~/.config/hyprtrigger && echo '{...}' > ~/.config/hyprtrigger/my-events.json")
	fmt.Println("")
	fmt.Println("Hot reload workflow:")
	fmt.Println("  1. Start daemon: hyprtrigger")
	fmt.Println("  2. Edit configs: ~/.config/hyprtrigger/*.json")
	fmt.Println("  3. Reload: hyprtrigger --reload")
	fmt.Println("  4. Check status: hyprtrigger --status")
	fmt.Println("  5. Stop daemon: hyprtrigger --shutdown")
	fmt.Println("")
	fmt.Println("Export and inspection:")
	fmt.Println("  hyprtrigger --print-builtin                     # Show builtin events")
	fmt.Println("  hyprtrigger --export-builtin builtin.json       # Export to file")
	fmt.Println("  hyprtrigger -p > ~/.config/hyprtrigger/template.json  # Create auto-config template")
	fmt.Println("  hyprtrigger -e ~/.config/hyprtrigger/builtin.json     # Export to auto-config")
	fmt.Println("")
	fmt.Println("Usage modes:")
	fmt.Println("  1. Daemon mode (default and recommended):")
	fmt.Println("     hyprtrigger")
	fmt.Println("     hyprtrigger --daemon")
	fmt.Println("     → Starts daemon with builtin events + auto-loads ~/.config/hyprtrigger/*.json")
	fmt.Println("     → Supports hot reload with --reload")
	fmt.Println("")
	fmt.Println("  2. Custom configuration:")
	fmt.Println("     hyprtrigger -c config.json")
	fmt.Println("     → Daemon with builtin + auto-config + manual config")
	fmt.Println("")
	fmt.Println("  3. Pure custom configuration:")
	fmt.Println("     hyprtrigger -c config.json --no-builtin --no-auto-config")
	fmt.Println("     → Daemon with only manual config file")
	fmt.Println("")
	fmt.Println("  4. Daemon control:")
	fmt.Println("     hyprtrigger --reload")
	fmt.Println("     → Hot reloads all configurations in running daemon")
	fmt.Println("")
	fmt.Println("Expected JSON format:")
	fmt.Println(`  {
    "events": [
      {
        "name": "windowtitlev2",
        "regex": "Bitwarden",
        "command": "hyprctl dispatch exec firefox",
        "use_shell": false
      }
    ]
  }`)
}

