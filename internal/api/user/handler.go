package user

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/mandacode-labs/dodream/internal/core"
	userservice "github.com/mandacode-labs/dodream/internal/service/user"
)

// Handler handles user-related HTTP requests.
type Handler struct {
	service *userservice.Service
}

// NewHandler creates a new Handler.
func NewHandler(service *userservice.Service) *Handler {
	return &Handler{service: service}
}

// CreateRequest represents the request body for creating a user.
type CreateRequest struct {
	Nickname   string `json:"nickname" binding:"required"`
	ProviderID string `json:"provider_id" binding:"required"`
}

// Create handles POST /users.
func (h *Handler) Create(c *gin.Context) {
	var req CreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.service.Create(c.Request.Context(), req.Nickname, req.ProviderID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, user)
}

// Get handles GET /users/:id.
func (h *Handler) Get(c *gin.Context) {
	id := c.Param("id")
	user, err := h.service.GetByID(c.Request.Context(), core.UserID(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, user)
}
