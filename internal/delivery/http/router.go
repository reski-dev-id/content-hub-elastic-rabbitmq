package http

import (
	"content-hub/internal/delivery/http/handler"

	"github.com/gin-gonic/gin"
)

func NewRouter(ph *handler.ProductHandler) *gin.Engine {
	r := gin.Default()

	v1 := r.Group("/v1")
	{
		v1.POST("/products", ph.Create)
		v1.GET("/products", ph.GetAll)
		v1.GET("/products/:id", ph.GetByID)
		v1.PUT("/products/:id", ph.Update)
		v1.DELETE("/products/:id", ph.Delete)
	}

	return r
}
