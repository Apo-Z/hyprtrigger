package events

import (
	"fmt"
	"time"
)

type Processor struct {
	registry       *Registry
	deduplicator   *DeduplicationManager
}

type DeduplicationManager struct {
	recentExecutions  []EventExecution
	deduplicationTime time.Duration
}

func NewDeduplicationManager() *DeduplicationManager {
	return &DeduplicationManager{
		recentExecutions:  make([]EventExecution, 0),
		deduplicationTime: 2 * time.Second,
	}
}

func (dm *DeduplicationManager) WasRecentlyExecuted(windowID, eventName, regex string) bool {
	now := time.Now()

	// Clean old executions
	filtered := make([]EventExecution, 0)
	for _, exec := range dm.recentExecutions {
		if now.Sub(exec.Timestamp) <= dm.deduplicationTime {
			filtered = append(filtered, exec)
		}
	}
	dm.recentExecutions = filtered

	// Check if recently executed
	for _, exec := range dm.recentExecutions {
		if exec.WindowID == windowID && exec.EventName == eventName && exec.Regex == regex {
			return true
		}
	}

	return false
}

func (dm *DeduplicationManager) RecordExecution(windowID, eventName, regex string) {
	dm.recentExecutions = append(dm.recentExecutions, EventExecution{
		WindowID:  windowID,
		EventName: eventName,
		Regex:     regex,
		Timestamp: time.Now(),
	})
}

func NewProcessor(registry *Registry) *Processor {
	return &Processor{
		registry:     registry,
		deduplicator: NewDeduplicationManager(),
	}
}

func (p *Processor) ProcessEvent(eventName, rawData string) error {
	eventData := ParseEventData(eventName, rawData)
	events := p.registry.GetEventsByName(eventName)

	for _, event := range events {
		if event.Match(eventData.Content) {
			if p.deduplicator.WasRecentlyExecuted(eventData.WindowID, eventName, event.Regex) {
				continue
			}

			if err := event.ExecuteCommand(eventData.WindowID); err != nil {
				return fmt.Errorf("command execution failed for %s: %w", event.Name, err)
			}

			p.deduplicator.RecordExecution(eventData.WindowID, eventName, event.Regex)
		}
	}
	return nil
}

var DefaultProcessor = NewProcessor(DefaultRegistry)

func ProcessEvent(eventName, data string) error {
	return DefaultProcessor.ProcessEvent(eventName, data)
}
