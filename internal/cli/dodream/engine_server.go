package dodream

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/spf13/cobra"

	"github.com/mandacode-labs/dodream/internal/config"
	"github.com/mandacode-labs/dodream/internal/server"
)

var engineServerCmd = &cobra.Command{
	Use:   "engine-server",
	Short: "Start the Dodream Engine server",
	Long:  `Starts the Dodream background engine server for processing study events and updating SRS states.`,
}

var engineServerRunCmd = &cobra.Command{
	Use:   "run",
	Short: "Run the engine server",
	RunE:  runEngineServer,
}

func init() {
	rootCmd.AddCommand(engineServerCmd)
	engineServerCmd.AddCommand(engineServerRunCmd)

	engineServerRunCmd.Flags().String("db-url", "", "PostgreSQL connection URL (overrides DATABASE_URL env)")
	engineServerRunCmd.Flags().String("nats-url", "", "NATS server URL (overrides NATS_URL env)")
	engineServerRunCmd.Flags().String("redis-url", "", "Redis server URL (overrides REDIS_URL env)")
	engineServerRunCmd.Flags().Int("window-size", 0, "Number of recent events to keep in the processing window (overrides ENGINE_WINDOW_SIZE env)")
}

func runEngineServer(cmd *cobra.Command, args []string) error {
	cfg := config.Load()

	if dbURL, _ := cmd.Flags().GetString("db-url"); dbURL != "" {
		cfg.Database.URL = dbURL
	}
	if natsURL, _ := cmd.Flags().GetString("nats-url"); natsURL != "" {
		cfg.NATS.URL = natsURL
	}
	if redisURL, _ := cmd.Flags().GetString("redis-url"); redisURL != "" {
		cfg.Redis.URL = redisURL
	}
	if windowSize, _ := cmd.Flags().GetInt("window-size"); windowSize > 0 {
		cfg.Engine.WindowSize = windowSize
	}

	if err := cfg.Validate(); err != nil {
		return fmt.Errorf("invalid configuration: %w", err)
	}

	engServer, err := server.NewEngineServer(cfg)
	if err != nil {
		return fmt.Errorf("failed to create engine server: %w", err)
	}

	log.Printf("Starting Engine server...")
	log.Printf("  NATS: %s", cfg.NATS.URL)
	log.Printf("  Redis: %s", cfg.Redis.URL)
	log.Printf("  Window size: %d", cfg.Engine.WindowSize)

	if err := engServer.NATSConsumer.Start(); err != nil {
		return fmt.Errorf("failed to start nats consumer: %w", err)
	}

	// Start gRPC server
	grpcAddr := cfg.Engine.GRPCAddr
	grpcSrv := server.NewEngineGRPCServer(engServer.Processor)
	grpcServer, err := server.StartGRPCServer(grpcAddr, grpcSrv)
	if err != nil {
		return fmt.Errorf("failed to start grpc server: %w", err)
	}
	engServer.GRPCSrv = grpcServer

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	<-sigCh
	log.Printf("Shutting down engine server...")
	if err := engServer.Close(); err != nil {
		log.Printf("Error during engine server shutdown: %v", err)
	}

	return nil
}
