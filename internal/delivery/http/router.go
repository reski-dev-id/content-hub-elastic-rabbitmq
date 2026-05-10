package http

import (
	"content-hub/internal/delivery/http/handler"

	"github.com/gin-gonic/gin"
)

func NewRouter(ph *handler.ProductHandler,
	nh *handler.NewsHandler,
	sh *handler.SearchHandler) *gin.Engine {
	r := gin.Default()

	v1 := r.Group("/v1")
	{
		v1.POST("/products", ph.Create)
		v1.GET("/products", ph.GetAll)
		v1.GET("/products/:id", ph.GetByID)
		v1.PUT("/products/:id", ph.Update)
		v1.DELETE("/products/:id", ph.Delete)

		v1.POST("/news", nh.Create)
		v1.GET("/news", nh.GetAll)
		v1.GET("/news/:id", nh.GetByID)
		v1.PUT("/news/:id", nh.Update)
		v1.DELETE("/news/:id", nh.Delete)

		v1.GET("/search", sh.Search)
	}

	return r
}
