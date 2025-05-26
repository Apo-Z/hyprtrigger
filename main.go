package main

import (
	"hyprtrigger/cmd/hyprtrigger"
	_ "hyprtrigger/events"
)

func main() {
	hyprtrigger.Execute()
}

