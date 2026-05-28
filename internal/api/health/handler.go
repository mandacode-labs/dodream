package health

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Handler handles health check requests.
type Handler struct{}

// NewHandler creates a new Handler.
func NewHandler() *Handler {
	return &Handler{}
}

// Get handles GET /health.
func (h *Handler) Get(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}
