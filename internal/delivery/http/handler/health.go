package handler

import (
	"net/http"

	"content-hub/internal/delivery/http/response"

	"github.com/gin-gonic/gin"
)

type HealthHandler struct{}

func NewHealthHandler() *HealthHandler {
	return &HealthHandler{}
}

// HealthCheck godoc
// @Summary Health check
// @Description Check API status
// @Tags health
// @Accept json
// @Produce json
// @Success 200 {object} response.Response
// @Router /health [get]
func (h *HealthHandler) Check(c *gin.Context) {

	response.Success(
		c,
		http.StatusOK,
		"service healthy",
		gin.H{
			"status": "UP",
		},
	)
}
