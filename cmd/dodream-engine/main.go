package main

import (
	"log"

	dodreamengine "github.com/mandacode-labs/dodream/internal/cli/dodream-engine"
)

var (
	version   = "dev"
	gitCommit = "unknown"
)

func main() {
	log.Printf("Dodream Engine %s (%s)", version, gitCommit)
	dodreamengine.Execute()
}
