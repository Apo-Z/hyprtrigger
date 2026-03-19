package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"hyprtrigger/internal/builtin"
	"hyprtrigger/internal/events"
)

var eventsCmd = &cobra.Command{
	Use:   "events",
	Short: "Inspect builtin events",
}

var eventsListCmd = &cobra.Command{
	Use:   "list",
	Short: "Print builtin events as JSON",
	RunE: func(cmd *cobra.Command, args []string) error {
		json, err := builtinEventsJSON()
		if err != nil {
			return err
		}
		fmt.Println(json)
		return nil
	},
}

var eventsExportCmd = &cobra.Command{
	Use:   "export <file>",
	Short: "Export builtin events to a JSON file",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		data, err := builtinEventsJSON()
		if err != nil {
			return err
		}
		if err := os.WriteFile(args[0], []byte(data), 0644); err != nil {
			return fmt.Errorf("failed to write file: %w", err)
		}
		fmt.Printf("Builtin events exported to: %s\n", args[0])
		return nil
	},
}

func init() {
	eventsCmd.AddCommand(eventsListCmd)
	eventsCmd.AddCommand(eventsExportCmd)
}

func builtinEventsJSON() (string, error) {
	r := events.NewRegistry()
	builtin.Register(r)

	var all []events.Event
	for _, list := range r.GetBuiltinEvents() {
		for _, ev := range list {
			all = append(all, *ev)
		}
	}

	cfg := struct {
		Events []events.Event `json:"events"`
	}{Events: all}

	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return "", fmt.Errorf("JSON serialization failed: %w", err)
	}
	return string(data), nil
}
