package dodream

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "dodream",
	Short: "Dodream flashcard SRS engine",
	Long: `Dodream is a flashcard application with an event-sourced SRS engine.
It provides an HTTP API server and a background engine server for processing study events.`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
