package middleware

import (
	"time"

	"content-hub/internal/logger"

	"github.com/gin-gonic/gin"
)

func Logger() gin.HandlerFunc {

	return func(c *gin.Context) {

		start := time.Now()

		c.Next()

		latency := time.Since(start)

		requestID, _ := c.Get(RequestIDKey)

		logger.Log.Info().
			Str("request_id", requestID.(string)).
			Str("method", c.Request.Method).
			Str("path", c.Request.URL.Path).
			Int("status", c.Writer.Status()).
			Str("ip", c.ClientIP()).
			Dur("latency", latency).
			Msg("incoming request")
	}
}
