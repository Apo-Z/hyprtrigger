package events

import (
	"regexp"
	"time"
)

type Event struct {
	Name     string `json:"name"`
	Regex    string `json:"regex"`
	Command  string `json:"command"`
	UseShell bool   `json:"use_shell"`
	compiled *regexp.Regexp
}

type EventData struct {
	WindowID string
	Content  string
}

type EventExecution struct {
	WindowID  string
	EventName string
	Regex     string
	Timestamp time.Time
}
