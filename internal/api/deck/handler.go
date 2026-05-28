package deck

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mandacode-labs/dodream/internal/core"
	deckservice "github.com/mandacode-labs/dodream/internal/service/deck"
)

// Handler handles deck-related HTTP requests.
type Handler struct {
	service deckservice.Service
}

// NewHandler creates a new Handler.
func NewHandler(service deckservice.Service) *Handler {
	return &Handler{service: service}
}

// CreateRequest represents the request body for creating a deck.
type CreateRequest struct {
	Name    string `json:"name" binding:"required"`
	Creator string `json:"creator" binding:"required"`
}

// Create handles POST /decks.
func (h *Handler) Create(c *gin.Context) {
	var req CreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	d, err := h.service.Create(c.Request.Context(), req.Name, core.UserID(req.Creator))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, d)
}

// Get handles GET /decks/:id.
func (h *Handler) Get(c *gin.Context) {
	id := c.Param("id")
	d, err := h.service.GetByID(c.Request.Context(), core.DeckID(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, d)
}
