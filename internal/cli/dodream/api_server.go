package dodream

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/lib/pq"
	"github.com/spf13/cobra"

	"github.com/mandacode-labs/dodream/internal/config"
	"github.com/mandacode-labs/dodream/internal/server"
)

var apiServerCmd = &cobra.Command{
	Use:   "api-server",
	Short: "Start the Dodream API server",
	Long:  `Starts the Dodream HTTP API server with all REST endpoints.`,
}

var apiServerRunCmd = &cobra.Command{
	Use:   "run",
	Short: "Run the API server",
	RunE:  runAPIServer,
}

func init() {
	rootCmd.AddCommand(apiServerCmd)
	apiServerCmd.AddCommand(apiServerRunCmd)

	apiServerRunCmd.Flags().String("addr", ":8080", "Server address to listen on")
	apiServerRunCmd.Flags().String("db-url", "", "PostgreSQL connection URL (overrides DATABASE_URL env)")
}

func runAPIServer(cmd *cobra.Command, args []string) error {
	cfg := config.Load()

	if addr, _ := cmd.Flags().GetString("addr"); addr != "" {
		cfg.Server.Addr = addr
	}
	if dbURL, _ := cmd.Flags().GetString("db-url"); dbURL != "" {
		cfg.Database.URL = dbURL
	}

	if err := cfg.Validate(); err != nil {
		return fmt.Errorf("invalid configuration: %w", err)
	}

	apiSrv, err := server.NewAPIServer(cfg)
	if err != nil {
		return fmt.Errorf("failed to create api server: %w", err)
	}
	defer apiSrv.Close()

	httpSrv := &http.Server{
		Addr:    cfg.Server.Addr,
		Handler: apiSrv.Server,
	}

	go func() {
		log.Printf("Starting Dodream API server on %s", cfg.Server.Addr)
		if err := httpSrv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Printf("HTTP server error: %v", err)
		}
	}()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	<-sigCh

	log.Printf("Shutting down API server...")
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := httpSrv.Shutdown(shutdownCtx); err != nil {
		log.Printf("HTTP server shutdown error: %v", err)
	}

	return nil
}
