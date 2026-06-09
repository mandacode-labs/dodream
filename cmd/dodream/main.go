package main

import (
	"log"

	dodreamcli "github.com/mandacode-labs/dodream/internal/cli/dodream"
)

var (
	version   = "dev"
	gitCommit = "unknown"
)

func main() {
	log.Printf("Dodream %s (%s)", version, gitCommit)
	dodreamcli.Execute()
}
