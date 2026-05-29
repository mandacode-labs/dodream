// Package engine provides the background engine server for processing study events.
package engine

import (
	"fmt"

	"google.golang.org/grpc"

	"github.com/mandacode-labs/dodream/ent"
	"github.com/mandacode-labs/dodream/internal/config"
	"github.com/mandacode-labs/dodream/internal/engine/state"
	"github.com/mandacode-labs/dodream/internal/nats"
	"github.com/mandacode-labs/dodream/internal/redis"
	"github.com/mandacode-labs/dodream/internal/studyevent"
)

// Server holds all dependencies for the engine server.
type Server struct {
	Config       *config.Config
	Client       *ent.Client
	StateService *state.Service
	RedisClient  *redis.Client
	NATSConsumer *nats.Consumer
	Processor    *EventProcessor
	GRPCSrv      *grpc.Server
}

// Close releases all resources held by the engine server.
func (s *Server) Close() error {
	if s.GRPCSrv != nil {
		s.GRPCSrv.GracefulStop()
	}
	if s.NATSConsumer != nil {
		_ = s.NATSConsumer.Stop()
	}
	if s.RedisClient != nil {
		s.RedisClient.Close()
	}
	if s.Client != nil {
		s.Client.Close()
	}
	return nil
}

// NewServer creates a new engine server with all dependencies wired.
func NewServer(cfg *config.Config) (*Server, error) {
	client, err := ent.Open("postgres", cfg.Database.URL)
	if err != nil {
		return nil, fmt.Errorf("failed opening connection to postgres: %w", err)
	}

	// Redis
	redisClient, err := redis.NewClient(cfg.Redis.URL)
	if err != nil {
		client.Close()
		return nil, fmt.Errorf("failed connecting to redis: %w", err)
	}

	// Stores
	stateStore := state.NewStore(client)
	studyEventStore := studyevent.NewStore(client)

	// Services
	stateService := state.NewService(stateStore)

	// NATS Consumer
	calculator := state.NewCalculator()
	processor := NewEventProcessor(stateService, studyEventStore, redisClient, calculator, cfg.Engine.WindowSize)
	consumer, err := nats.NewConsumer(cfg.NATS.URL, "engine-workers", processor.HandleEvent)
	if err != nil {
		redisClient.Close()
		client.Close()
		return nil, fmt.Errorf("failed to create nats consumer: %w", err)
	}

	return &Server{
		Config:       cfg,
		Client:       client,
		StateService: stateService,
		RedisClient:  redisClient,
		NATSConsumer: consumer,
		Processor:    processor,
	}, nil
}
