package card

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/mandacode-labs/dodream/internal/core"
	cardservice "github.com/mandacode-labs/dodream/internal/service/card"
)

// Handler handles card-related HTTP requests.
type Handler struct {
	service *cardservice.Service
}

// NewHandler creates a new Handler.
func NewHandler(service *cardservice.Service) *Handler {
	return &Handler{service: service}
}

// CreateRequest represents the request body for creating a card.
type CreateRequest struct {
	Hint    string `json:"hint"`
	Content string `json:"content" binding:"required"`
	Creator string `json:"creator" binding:"required"`
}

// Create handles POST /cards.
func (h *Handler) Create(c *gin.Context) {
	var req CreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	card, err := h.service.Create(c.Request.Context(), req.Hint, req.Content, core.UserID(req.Creator))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, card)
}

// Get handles GET /cards/:id.
func (h *Handler) Get(c *gin.Context) {
	id := c.Param("id")
	card, err := h.service.GetByID(c.Request.Context(), core.CardID(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, card)
}
