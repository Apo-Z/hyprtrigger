package events

import (
	"fmt"
	"os/exec"
	"regexp"
	"strings"
)

func (ev *Event) ExecuteCommand(windowID string) error {
	command := strings.ReplaceAll(ev.Command, "{WINDOW_ID}", windowID)

	var cmd *exec.Cmd

	if ev.UseShell {
		cmd = exec.Command("sh", "-c", command)
	} else {
		parts := strings.Fields(command)
		if len(parts) == 0 {
			return fmt.Errorf("empty command")
		}
		cmd = exec.Command(parts[0], parts[1:]...)
	}

	return cmd.Run()
}

func (ev *Event) Match(input string) bool {
	if ev.compiled == nil {
		var err error
		ev.compiled, err = regexp.Compile(ev.Regex)
		if err != nil {
			return false
		}
	}
	return ev.compiled.MatchString(input)
}
