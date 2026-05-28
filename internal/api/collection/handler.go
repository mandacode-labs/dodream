package collection

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mandacode-labs/dodream/internal/core"
	collectionservice "github.com/mandacode-labs/dodream/internal/service/collection"
)

// Handler handles collection-related HTTP requests.
type Handler struct {
	service collectionservice.Service
}

// NewHandler creates a new Handler.
func NewHandler(service collectionservice.Service) *Handler {
	return &Handler{service: service}
}

// CreateRequest represents the request body for creating a collection.
type CreateRequest struct {
	Name    string `json:"name" binding:"required"`
	Creator string `json:"creator" binding:"required"`
}

// Create handles POST /collections.
func (h *Handler) Create(c *gin.Context) {
	var req CreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	col, err := h.service.Create(c.Request.Context(), req.Name, core.UserID(req.Creator))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, col)
}

// Get handles GET /collections/:id.
func (h *Handler) Get(c *gin.Context) {
	id := c.Param("id")
	col, err := h.service.GetByID(c.Request.Context(), core.CollectionID(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, col)
}
