package http

import (
	"content-hub/internal/delivery/http/handler"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/gin-gonic/gin"
)

func NewRouter(
	ph *handler.ProductHandler,
	nh *handler.NewsHandler,
	sh *handler.SearchHandler,
	hh *handler.HealthHandler,
) *gin.Engine {

	r := gin.Default()

	// Health Check
	r.GET("/health", hh.Check)

	v1 := r.Group("/v1")
	{
		// Swagger
		v1.GET(
			"/swagger/*any",
			ginSwagger.WrapHandler(
				swaggerFiles.Handler,
			),
		)

		// Product
		v1.POST("/products", ph.Create)
		v1.GET("/products", ph.GetAll)
		v1.GET("/products/:id", ph.GetByID)
		v1.PUT("/products/:id", ph.Update)
		v1.DELETE("/products/:id", ph.Delete)

		// News
		v1.POST("/news", nh.Create)
		v1.GET("/news", nh.GetAll)
		v1.GET("/news/:id", nh.GetByID)
		v1.PUT("/news/:id", nh.Update)
		v1.DELETE("/news/:id", nh.Delete)

		// Search
		v1.GET("/search", sh.Search)
	}

	return r
}
