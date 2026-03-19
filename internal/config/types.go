package config

import "hyprtrigger/internal/events"

type EventConfig struct {
	Events []events.Event `json:"events"`
}
