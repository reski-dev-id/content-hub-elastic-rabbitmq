package handler

import (
	"context"
	"net/http"
	"time"

	"content-hub/internal/delivery/http/response"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	amqp "github.com/rabbitmq/amqp091-go"
)

type HealthHandler struct {
	db  *sqlx.DB
	es  *elasticsearch.Client
	rmq *amqp.Connection
}

func NewHealthHandler(
	db *sqlx.DB,
	es *elasticsearch.Client,
	rmq *amqp.Connection,
) *HealthHandler {
	return &HealthHandler{
		db:  db,
		es:  es,
		rmq: rmq,
	}
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

	mysqlStatus := "UP"
	elasticStatus := "UP"
	rabbitStatus := "UP"

	if err := h.db.Ping(); err != nil {
		mysqlStatus = "DOWN"
	}

	ctx, cancel := context.WithTimeout(
		context.Background(),
		5*time.Second,
	)

	defer cancel()

	_, err := h.es.Info(
		h.es.Info.WithContext(ctx),
	)

	if err != nil {
		elasticStatus = "DOWN"
	}

	if h.rmq == nil || h.rmq.IsClosed() {
		rabbitStatus = "DOWN"
	}

	response.Success(
		c,
		http.StatusOK,
		"service healthy",
		gin.H{
			"app":           "UP",
			"mysql":         mysqlStatus,
			"elasticsearch": elasticStatus,
			"rabbitmq":      rabbitStatus,
		},
	)
}
