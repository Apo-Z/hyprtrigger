package config

import "hyprtrigger/pkg/events"

type EventConfig struct {
	Events []events.Event `json:"events"`
}
