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

	db, _ := mysql.NewDB(cfg.DBUrl)

	productRepo := mysqlRepo.NewProductRepository(db)
	productUC := usecase.NewProductUsecase(productRepo)
	productHandler := handler.NewProductHandler(productUC)

	router := http.NewRouter(productHandler)
	router.Run(":" + cfg.AppPort)
}
