// @title Content Hub API
// @version 1.0
// @description Product and News API with Elasticsearch Search
// @host localhost:8080
// @BasePath /

package main

import (
	_ "content-hub/docs"

	"content-hub/config"
	"content-hub/internal/delivery/http"
	"content-hub/internal/delivery/http/handler"
	esInfra "content-hub/internal/infrastructure/elasticsearch"
	"content-hub/internal/infrastructure/mysql"
	esRepo "content-hub/internal/repository/elasticsearch"
	mysqlRepo "content-hub/internal/repository/mysql"
	"content-hub/internal/usecase"
)

func main() {

	cfg := config.Load()

	db, err := mysql.NewDB(cfg.DBUrl)
	if err != nil {
		panic(err)
	}

	esClient, err := esInfra.NewClient(
		cfg.ElasticUrl,
	)

	if err != nil {
		panic(err)
	}

	healthHandler := handler.NewHealthHandler()

	// Product
	productRepo := mysqlRepo.NewProductRepository(db)

	productUC := usecase.NewProductUsecase(
		productRepo,
	)

	productHandler := handler.NewProductHandler(
		productUC,
	)

	// News
	newsRepo := mysqlRepo.NewNewsRepository(db)

	newsUC := usecase.NewNewsUsecase(
		newsRepo,
	)

	newsHandler := handler.NewNewsHandler(
		newsUC,
	)

	// Search
	searchRepo := esRepo.NewSearchRepository(
		esClient,
	)

	searchUC := usecase.NewSearchUsecase(
		searchRepo,
	)

	searchHandler := handler.NewSearchHandler(
		searchUC,
	)

	// Router
	router := http.NewRouter(
		productHandler,
		newsHandler,
		searchHandler,
		healthHandler,
	)

	router.Run(":" + cfg.AppPort)
}
