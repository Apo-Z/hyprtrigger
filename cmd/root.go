package cmd

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/spf13/cobra"
	"hyprtrigger/internal/builtin"
	"hyprtrigger/internal/config"
	"hyprtrigger/internal/daemon"
	"hyprtrigger/internal/events"
	"hyprtrigger/internal/hyprland"
)

var (
	configPath   string
	noBuiltin    bool
	noAutoConfig bool
)

var rootCmd = &cobra.Command{
	Use:           "hyprtrigger",
	Short:         "Event-driven automation daemon for Hyprland",
	SilenceUsage:  true,
	SilenceErrors: true,
	RunE:          runDaemon,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&configPath, "config", "c", "", "Path to JSON config file or directory")
	rootCmd.PersistentFlags().BoolVarP(&noBuiltin, "no-builtin", "n", false, "Disable builtin events")
	rootCmd.PersistentFlags().BoolVarP(&noAutoConfig, "no-auto-config", "s", false, "Skip auto-loading from ~/.config/hyprtrigger/")

	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(reloadCmd)
	rootCmd.AddCommand(statusCmd)
	rootCmd.AddCommand(shutdownCmd)
	rootCmd.AddCommand(initConfigCmd)
	rootCmd.AddCommand(eventsCmd)
}

func runDaemon(cmd *cobra.Command, args []string) error {
	if daemon.IsDaemonRunning() {
		fmt.Println("Hyprtrigger is already running!")
		fmt.Println("  hyprtrigger reload    # Reload configuration")
		fmt.Println("  hyprtrigger status    # Check status")
		fmt.Println("  hyprtrigger shutdown  # Stop daemon")
		return fmt.Errorf("daemon already running")
	}

	fmt.Println("Starting Hyprland event monitor")

	if err := loadConfig(); err != nil {
		return err
	}

	printEventsSummary()

	if len(events.GetAllEvents()) == 0 {
		return fmt.Errorf("no events loaded. Use -c to specify a config file or run 'hyprtrigger init-config'")
	}

	daemonServer := daemon.NewDaemon()
	if err := daemonServer.Start(); err != nil {
		return fmt.Errorf("failed to start daemon: %w", err)
	}
	defer daemonServer.Stop()

	client := hyprland.NewClient()
	if err := client.Connect(); err != nil {
		return fmt.Errorf("%v\nMake sure Hyprland is running", err)
	}
	defer client.Close()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	listenerDone := make(chan error, 1)
	go func() {
		listenerDone <- hyprland.NewListener(client).Listen()
	}()

	for {
		select {
		case <-daemonServer.GetReloadChannel():
			fmt.Println("Reloading configuration...")
			if err := reloadConfig(); err != nil {
				fmt.Printf("Reload failed: %v\n", err)
			} else {
				printEventsSummary()
			}

		case <-daemonServer.GetShutdownChannel():
			fmt.Println("Shutdown requested")
			return nil

		case sig := <-sigChan:
			fmt.Printf("\nReceived signal: %v\n", sig)
			return nil

		case err := <-listenerDone:
			return err
		}
	}
}

func loadConfig() error {
	if noBuiltin {
		fmt.Println("Builtin events disabled")
		events.DefaultRegistry.SetSkipBuiltinEvents(true)
	} else {
		builtin.Register(events.DefaultRegistry)
	}

	if !noAutoConfig {
		if err := config.LoadAutoConfig(); err != nil {
			fmt.Printf("Auto-config loading failed: %v\n", err)
		}
	}

	if configPath != "" {
		fmt.Printf("Loading config: %s\n", configPath)
		if err := config.LoadEventsFromPath(configPath); err != nil {
			return fmt.Errorf("config loading failed: %w", err)
		}
	}

	return nil
}

func reloadConfig() error {
	events.DefaultRegistry.Clear()
	if err := loadConfig(); err != nil {
		return err
	}
	if len(events.GetAllEvents()) == 0 {
		return fmt.Errorf("no events loaded after reload")
	}
	return nil
}

func printEventsSummary() {
	allEvents := events.GetAllEvents()
	if len(allEvents) == 0 {
		fmt.Println("No events loaded")
		return
	}

	fmt.Println("\nLoaded events:")
	total := 0
	for name, list := range allEvents {
		fmt.Printf("  %s: %d event(s)\n", name, len(list))
		for _, ev := range list {
			cmd := ev.Command
			if len(cmd) > 50 {
				cmd = cmd[:47] + "..."
			}
			fmt.Printf("    - regex: %-30s  cmd: %s\n", ev.Regex, cmd)
		}
		total += len(list)
	}
	fmt.Printf("Total: %d event(s)\n\n", total)
}
