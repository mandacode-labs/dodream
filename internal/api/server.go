package api

import (
	"github.com/gin-gonic/gin"
	"github.com/mandacode-labs/dodream/ent"
	"github.com/mandacode-labs/dodream/internal/api/card"
	"github.com/mandacode-labs/dodream/internal/api/collection"
	"github.com/mandacode-labs/dodream/internal/api/deck"
	"github.com/mandacode-labs/dodream/internal/api/health"
	"github.com/mandacode-labs/dodream/internal/api/user"
	cardservice "github.com/mandacode-labs/dodream/internal/service/card"
	collectionservice "github.com/mandacode-labs/dodream/internal/service/collection"
	deckservice "github.com/mandacode-labs/dodream/internal/service/deck"
	userservice "github.com/mandacode-labs/dodream/internal/service/user"
	"github.com/mandacode-labs/dodream/internal/store"
)

// Server represents the HTTP server.
type Server struct {
	addr   string
	router *gin.Engine
	client *ent.Client
}

// NewServer creates a new HTTP server with the given address and ent client.
func NewServer(addr string, client *ent.Client) *Server {
	router := gin.Default()
	s := &Server{
		addr:   addr,
		router: router,
		client: client,
	}
	s.setupRoutes()
	return s
}

// Run starts the HTTP server.
func (s *Server) Run() error {
	return s.router.Run(s.addr)
}

func (s *Server) setupRoutes() {
	// Stores
	userStore := store.NewUserStore(s.client)
	cardStore := store.NewCardStore(s.client)
	deckStore := store.NewDeckStore(s.client)
	collectionStore := store.NewCollectionStore(s.client)

	// Services
	userService := userservice.NewUserService(userStore)
	cardService := cardservice.NewCardService(cardStore)
	deckService := deckservice.NewDeckService(deckStore)
	collectionService := collectionservice.NewCollectionService(collectionStore)

	// Handlers
	healthHandler := health.NewHandler()
	userHandler := user.NewHandler(userService)
	cardHandler := card.NewHandler(cardService)
	deckHandler := deck.NewHandler(deckService)
	collectionHandler := collection.NewHandler(collectionService)

	// Health
	s.router.GET("/health", healthHandler.Get)

	// API v1
	v1 := s.router.Group("/api/v1")
	{
		// Users
		v1.POST("/users", userHandler.Create)
		v1.GET("/users/:id", userHandler.Get)

		// Cards
		v1.POST("/cards", cardHandler.Create)
		v1.GET("/cards/:id", cardHandler.Get)

		// Decks
		v1.POST("/decks", deckHandler.Create)
		v1.GET("/decks/:id", deckHandler.Get)

		// Collections
		v1.POST("/collections", collectionHandler.Create)
		v1.GET("/collections/:id", collectionHandler.Get)
	}
}
