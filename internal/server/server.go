// Package server provides dependency injection and server composition.
package server

import (
	"fmt"

	"google.golang.org/grpc"

	"github.com/mandacode-labs/dodream/ent"
	"github.com/mandacode-labs/dodream/internal/config"
	"github.com/mandacode-labs/dodream/internal/nats"
	"github.com/mandacode-labs/dodream/internal/redis"
	"github.com/mandacode-labs/dodream/internal/state"
	"github.com/mandacode-labs/dodream/internal/studyevent"
)

// EngineServer holds all dependencies for the engine server.
type EngineServer struct {
	Config       *config.Config
	Client       *ent.Client
	StateService *state.Service
	RedisClient  *redis.Client
	NATSConsumer *nats.Consumer
	Processor    *EventProcessor
	GRPCSrv      *grpc.Server
}

// Close releases all resources held by the engine server.
func (s *EngineServer) Close() error {
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

// NewEngineServer creates a new engine server with all dependencies wired.
func NewEngineServer(cfg *config.Config) (*EngineServer, error) {
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

	return &EngineServer{
		Config:       cfg,
		Client:       client,
		StateService: stateService,
		RedisClient:  redisClient,
		NATSConsumer: consumer,
		Processor:    processor,
	}, nil
}

// EventProcessor handles study events for the engine.
type EventProcessor struct {
	stateService    *state.Service
	studyEventStore *studyevent.Store
	redisClient     *redis.Client
	calculator      *state.Calculator
	windowSize      int
}

// NewEventProcessor creates a new event processor.
func NewEventProcessor(stateService *state.Service, studyEventStore *studyevent.Store, redisClient *redis.Client, calculator *state.Calculator, windowSize int) *EventProcessor {
	return &EventProcessor{
		stateService:    stateService,
		studyEventStore: studyEventStore,
		redisClient:     redisClient,
		calculator:      calculator,
		windowSize:      windowSize,
	}
}

// HandleEvent processes a single study event message.
func (p *EventProcessor) HandleEvent(msg *nats.StudyEventMessage) error {
	// Simplified - in production would implement full logic
	return nil
}
