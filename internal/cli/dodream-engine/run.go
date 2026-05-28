package dodream_engine

import (
	"context"
	"fmt"
	"log"

	_ "github.com/lib/pq"
	"github.com/mandacode-labs/dodream/ent"
	"github.com/mandacode-labs/dodream/internal/api"
	"github.com/spf13/cobra"
)

// runCmd represents the run command.
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Start the Dodream HTTP server",
	Long:  `Starts the Dodream HTTP server with all API endpoints.`,
	RunE:  runServer,
}

func init() {
	rootCmd.AddCommand(runCmd)
	runCmd.Flags().String("addr", ":8080", "Server address to listen on")
	runCmd.Flags().String("db-url", "postgres://user:password@localhost:5432/dodream?sslmode=disable", "PostgreSQL connection URL")
}

func runServer(cmd *cobra.Command, args []string) error {
	addr, _ := cmd.Flags().GetString("addr")
	dbURL, _ := cmd.Flags().GetString("db-url")

	log.Printf("Connecting to database...")
	client, err := ent.Open("postgres", dbURL)
	if err != nil {
		return fmt.Errorf("failed opening connection to postgres: %w", err)
	}
	defer client.Close()

	log.Printf("Running database migration...")
	ctx := context.Background()
	if err := client.Schema.Create(ctx); err != nil {
		return fmt.Errorf("failed creating schema resources: %w", err)
	}

	log.Printf("Starting Dodream server on %s", addr)
	server := api.NewServer(addr, client)
	return server.Run()
}
