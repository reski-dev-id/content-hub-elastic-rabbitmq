package main

import (
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

	productRepo := mysqlRepo.NewProductRepository(db)
	productUC := usecase.NewProductUsecase(productRepo)
	productHandler := handler.NewProductHandler(productUC)

	newsRepo := mysqlRepo.NewNewsRepository(db)
	newsUC := usecase.NewNewsUsecase(newsRepo)
	newsHandler := handler.NewNewsHandler(newsUC)

	searchRepo := esRepo.NewSearchRepository(
		esClient,
	)

	searchUC := usecase.NewSearchUsecase(
		searchRepo,
	)

	searchHandler := handler.NewSearchHandler(
		searchUC,
	)

	router := http.NewRouter(
		productHandler,
		newsHandler,
		searchHandler,
	)

	router.Run(":" + cfg.AppPort)
}
