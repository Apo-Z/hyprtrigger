package hyprtrigger

import "fmt"

var (
	version = "dev"
	commit  = "unknown"
	date    = "unknown"
)

func PrintVersion() {
	fmt.Printf("hyprtrigger %s (commit: %s, built: %s)\n", version, commit, date)
}
