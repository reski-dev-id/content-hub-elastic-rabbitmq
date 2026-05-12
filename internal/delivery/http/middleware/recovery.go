package middleware

import (
	"net/http"
	"runtime/debug"

	"content-hub/internal/delivery/http/response"
	"content-hub/internal/logger"

	"github.com/gin-gonic/gin"
)

func RecoveryMiddleware() gin.HandlerFunc {

	return func(c *gin.Context) {

		defer func() {

			if err := recover(); err != nil {

				requestID, _ := c.Get(RequestIDKey)

				logger.Error(nil).
					Str("request_id", requestID.(string)).
					Interface("panic", err).
					Str("method", c.Request.Method).
					Str("path", c.Request.URL.Path).
					Bytes("stack", debug.Stack()).
					Msg("panic recovered")

				response.Error(
					c,
					http.StatusInternalServerError,
					"internal server error",
					nil,
				)

				c.Abort()
			}
		}()

		c.Next()
	}
}
