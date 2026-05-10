package main

import (
	"content-hub/config"
	"content-hub/internal/delivery/http"
	"content-hub/internal/delivery/http/handler"
	"content-hub/internal/infrastructure/mysql"
	mysqlRepo "content-hub/internal/repository/mysql"
	"content-hub/internal/usecase"
)

func main() {
	cfg := config.Load()

	db, err := mysql.NewDB(cfg.DBUrl)
	if err != nil {
		panic(err)
	}

	productRepo := mysqlRepo.NewProductRepository(db)
	productUC := usecase.NewProductUsecase(productRepo)
	productHandler := handler.NewProductHandler(productUC)

	newsRepo := mysqlRepo.NewNewsRepository(db)
	newsUC := usecase.NewNewsUsecase(newsRepo)
	newsHandler := handler.NewNewsHandler(newsUC)

	router := http.NewRouter(
		productHandler,
		newsHandler,
	)

	router.Run(":" + cfg.AppPort)
}
