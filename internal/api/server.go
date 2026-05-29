package api

import (
	"context"
	"fmt"
	"log"

	"github.com/mandacode-labs/dodream/ent"
	"github.com/mandacode-labs/dodream/internal/card"
	"github.com/mandacode-labs/dodream/internal/collection"
	"github.com/mandacode-labs/dodream/internal/config"
	"github.com/mandacode-labs/dodream/internal/deck"
	"github.com/mandacode-labs/dodream/internal/nats"
	"github.com/mandacode-labs/dodream/internal/redis"
	"github.com/mandacode-labs/dodream/internal/studyevent"
	"github.com/mandacode-labs/dodream/internal/user"
	"github.com/mandacode-labs/dodream/pkg/oas"
)

// Server holds the HTTP server and its dependencies for cleanup.
type Server struct {
	Server *oas.Server
	Client *ent.Client
	NATS   *nats.Publisher
	Redis  *redis.Client
}

// Close releases all resources held by the API server.
func (s *Server) Close() error {
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

type redisNextCardService struct {
	client *redis.Client
}

func (s *redisNextCardService) GetNextCards(ctx context.Context, userID, collectionID string, limit int64) ([]oas.CardRecommendation, error) {
	cards, err := s.client.GetNextCards(ctx, userID, collectionID, limit)
	if err != nil {
		return nil, err
	}

	results := make([]oas.CardRecommendation, len(cards))
	for i, c := range cards {
		cardID, _ := c.Member.(string)
		results[i] = oas.CardRecommendation{
			CardID:     oas.NewOptString(cardID),
			FinalScore: oas.NewOptFloat64(c.Score),
		}
	}
	return results, nil
}

// NewServer creates and configures the ogen HTTP server.
func NewServer(cfg *config.Config) (*Server, error) {
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

	var nextCardSvc NextCardService
	if redisClient != nil {
		nextCardSvc = &redisNextCardService{client: redisClient}
	} else {
		nextCardSvc = &noopNextCardService{}
	}

	// Handler
	handler := NewHandler(
		userService,
		cardService,
		deckService,
		collectionService,
		studyEventService,
		nextCardSvc,
	)

	srv, err := oas.NewServer(handler)
	if err != nil {
		return nil, err
	}

	return &Server{
		Server: srv,
		Client: client,
		NATS:   natsPublisher,
		Redis:  redisClient,
	}, nil
}

type noopNextCardService struct{}

func (s *noopNextCardService) GetNextCards(ctx context.Context, userID, collectionID string, limit int64) ([]oas.CardRecommendation, error) {
	return nil, fmt.Errorf("redis not available")
}

var _ NextCardService = (*noopNextCardService)(nil)
var _ NextCardService = (*redisNextCardService)(nil)
