package hyprland

import (
	"bufio"
	"fmt"
	"hyprtrigger/internal/events"
	"strings"
)

type Listener struct {
	client *Client
}

func NewListener(client *Client) *Listener {
	return &Listener{client: client}
}

func (l *Listener) Listen() error {
	scanner := bufio.NewScanner(l.client.GetConnection())

	for scanner.Scan() {
		line := scanner.Text()

		parts := strings.SplitN(line, ">>", 2)
		if len(parts) != 2 {
			continue
		}

		eventName, eventData := parts[0], parts[1]
		fmt.Printf("Event: %s -> %s\n", eventName, eventData)

		if err := events.ProcessEvent(eventName, eventData); err != nil {
			fmt.Printf("Processing error: %v\n", err)
		}
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("socket read error: %w", err)
	}
	return nil
}
