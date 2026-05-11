// @title Content Hub API
// @version 1.0
// @description Product and News API with Elasticsearch Search
// @host localhost:8080
// @BasePath /

package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "content-hub/docs"

	"content-hub/config"
	httpDelivery "content-hub/internal/delivery/http"
	"content-hub/internal/delivery/http/handler"
	esInfra "content-hub/internal/infrastructure/elasticsearch"
	"content-hub/internal/infrastructure/mysql"
	rabbitInfra "content-hub/internal/infrastructure/rabbitmq"
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

	rabbitConn, err := rabbitInfra.NewRabbitMQ(
		cfg.RabbitMQUrl,
	)

	if err != nil {
		panic(err)
	}

	healthHandler := handler.NewHealthHandler(
		db,
		esClient,
		rabbitConn,
	)

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
	router := httpDelivery.NewRouter(
		productHandler,
		newsHandler,
		searchHandler,
		healthHandler,
	)

	server := &http.Server{
		Addr:    ":" + cfg.AppPort,
		Handler: router,
	}

	go func() {

		log.Printf(
			"server running on port %s",
			cfg.AppPort,
		)

		if err := server.ListenAndServe(); err != nil &&
			err != http.ErrServerClosed {

			log.Fatal(err)
		}
	}()

	quit := make(chan os.Signal, 1)

	signal.Notify(
		quit,
		syscall.SIGINT,
		syscall.SIGTERM,
	)

	<-quit

	log.Println("shutting down server...")

	ctx, cancel := context.WithTimeout(
		context.Background(),
		10*time.Second,
	)

	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatal(err)
	}

	if err := db.Close(); err != nil {
		log.Println("failed to close mysql:", err)
	}

	if err := rabbitConn.Close(); err != nil {
		log.Println("failed to close rabbitmq:", err)
	}

	log.Println("server exited properly")
}
