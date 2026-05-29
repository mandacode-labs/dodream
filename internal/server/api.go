package server

import (
	"fmt"
	"log"

	"github.com/mandacode-labs/dodream/ent"
	"github.com/mandacode-labs/dodream/internal/api"
	"github.com/mandacode-labs/dodream/internal/card"
	"github.com/mandacode-labs/dodream/internal/collection"
	"github.com/mandacode-labs/dodream/internal/config"
	"github.com/mandacode-labs/dodream/internal/deck"
	"github.com/mandacode-labs/dodream/internal/nats"
	"github.com/mandacode-labs/dodream/internal/redis"
	"github.com/mandacode-labs/dodream/internal/studyevent"
	"github.com/mandacode-labs/dodream/internal/user"
	gen "github.com/mandacode-labs/dodream/pkg/api/v1"
)

// APIServer holds the HTTP server and its dependencies for cleanup.
type APIServer struct {
	Server *gen.Server
	Client *ent.Client
	NATS   *nats.Publisher
	Redis  *redis.Client
}

// Close releases all resources held by the API server.
func (s *APIServer) Close() error {
	if s.NATS != nil {
		s.NATS.Close()
	}
	if s.Redis != nil {
		s.Redis.Close()
	}
	if s.Client != nil {
		s.Client.Close()
	}
	return nil
}

// NewAPIServer creates and configures the ogen HTTP server.
func NewAPIServer(cfg *config.Config) (*APIServer, error) {
	client, err := ent.Open("postgres", cfg.Database.URL)
	if err != nil {
		return nil, fmt.Errorf("failed opening connection to postgres: %w", err)
	}

	// Stores
	userStore := user.NewStore(client)
	cardStore := card.NewStore(client)
	deckStore := deck.NewStore(client)
	collectionStore := collection.NewStore(client)
	studyEventStore := studyevent.NewStore(client)

	// NATS Publisher (optional)
	var natsPublisher *nats.Publisher
	if cfg.NATS.URL != "" {
		var err error
		natsPublisher, err = nats.NewPublisher(cfg.NATS.URL)
		if err != nil {
			log.Printf("Warning: failed to create NATS publisher: %v", err)
		}
	}

	// Redis Client (optional)
	var redisClient *redis.Client
	if cfg.Redis.URL != "" {
		var err error
		redisClient, err = redis.NewClient(cfg.Redis.URL)
		if err != nil {
			log.Printf("Warning: failed to create Redis client: %v", err)
		}
	}

	// Services
	userService := user.NewService(userStore)
	cardService := card.NewService(cardStore)
	deckService := deck.NewService(deckStore)
	collectionService := collection.NewService(collectionStore)
	studyEventService := studyevent.NewService(studyEventStore, natsPublisher)

	// Handler
	handler := api.NewHandler(
		userService,
		cardService,
		deckService,
		collectionService,
		studyEventService,
		redisClient,
	)

	srv, err := gen.NewServer(handler)
	if err != nil {
		return nil, err
	}

	return &APIServer{
		Server: srv,
		Client: client,
		NATS:   natsPublisher,
		Redis:  redisClient,
	}, nil
}
