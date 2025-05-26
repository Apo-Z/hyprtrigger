package socket

import (
	"bufio"
	"fmt"
	"hyprtrigger/pkg/events"
	"strings"
)

type Listener struct {
	client *Client
}

func NewListener(client *Client) *Listener {
	return &Listener{
		client: client,
	}
}

func (l *Listener) Listen() error {
	conn := l.client.GetConnection()
	scanner := bufio.NewScanner(conn)

	for scanner.Scan() {
		line := scanner.Text()

		parts := strings.SplitN(line, ">>", 2)
		if len(parts) != 2 {
			continue
		}

		eventName := parts[0]
		eventData := parts[1]

		fmt.Printf("Event received: %s -> %s\n", eventName, eventData)

		if err := events.ProcessEvent(eventName, eventData); err != nil {
			fmt.Printf("Processing failed: %v\n", err)
		}
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("socket read error: %w", err)
	}

	return nil
}
