package events

import (
	"fmt"
	"time"
)

type Processor struct {
	registry     *Registry
	deduplicator *deduplicationManager
}

type deduplicationManager struct {
	recentExecutions  []EventExecution
	deduplicationTime time.Duration
}

func newDeduplicationManager() *deduplicationManager {
	return &deduplicationManager{
		recentExecutions:  make([]EventExecution, 0),
		deduplicationTime: 2 * time.Second,
	}
}

func (dm *deduplicationManager) wasRecentlyExecuted(windowID, eventName, regex string) bool {
	now := time.Now()

	filtered := make([]EventExecution, 0)
	for _, exec := range dm.recentExecutions {
		if now.Sub(exec.Timestamp) <= dm.deduplicationTime {
			filtered = append(filtered, exec)
		}
	}
	dm.recentExecutions = filtered

	for _, exec := range dm.recentExecutions {
		if exec.WindowID == windowID && exec.EventName == eventName && exec.Regex == regex {
			return true
		}
	}
	return false
}

func (dm *deduplicationManager) record(windowID, eventName, regex string) {
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
		deduplicator: newDeduplicationManager(),
	}
}

func (p *Processor) ProcessEvent(eventName, rawData string) error {
	eventData := ParseEventData(eventName, rawData)
	events := p.registry.GetEventsByName(eventName)

	for _, event := range events {
		if !event.Match(eventData.Content) {
			continue
		}
		if p.deduplicator.wasRecentlyExecuted(eventData.WindowID, eventName, event.Regex) {
			continue
		}
		if err := event.ExecuteCommand(eventData.WindowID); err != nil {
			return fmt.Errorf("command execution failed for %s: %w", event.Name, err)
		}
		p.deduplicator.record(eventData.WindowID, eventName, event.Regex)
	}
	return nil
}

var DefaultProcessor = NewProcessor(DefaultRegistry)

func ProcessEvent(eventName, data string) error {
	return DefaultProcessor.ProcessEvent(eventName, data)
}
