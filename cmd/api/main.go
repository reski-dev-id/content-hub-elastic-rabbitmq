// @title Content Hub API
// @version 1.0
// @description Product and News API with Elasticsearch Search
// @host localhost:8080
// @BasePath /

package main

import (
	"context"
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
	"content-hub/internal/logger"
	esRepo "content-hub/internal/repository/elasticsearch"
	mysqlRepo "content-hub/internal/repository/mysql"
	"content-hub/internal/usecase"
)

func main() {

	logger.Init()

	cfg := config.Load()

	db, err := mysql.NewDB(cfg.DBUrl)
	if err != nil {

		logger.Fatal(err).
			Msg("failed to connect mysql")
	}

	esClient, err := esInfra.NewClient(
		cfg.ElasticUrl,
	)

	if err != nil {

		logger.Fatal(err).
			Msg("failed to connect elasticsearch")
	}

	rabbitConn, err := rabbitInfra.NewRabbitMQ(
		cfg.RabbitMQUrl,
	)

	if err != nil {

		logger.Fatal(err).
			Msg("failed to connect rabbitmq")
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

		logger.Info().
			Str("port", cfg.AppPort).
			Msg("server running")

		if err := server.ListenAndServe(); err != nil &&
			err != http.ErrServerClosed {

			logger.Fatal(err).
				Msg("failed to start server")
		}
	}()

	quit := make(chan os.Signal, 1)

	signal.Notify(
		quit,
		syscall.SIGINT,
		syscall.SIGTERM,
	)

	<-quit

	logger.Info().
		Msg("shutting down server")

	ctx, cancel := context.WithTimeout(
		context.Background(),
		10*time.Second,
	)

	defer cancel()

	if err := server.Shutdown(ctx); err != nil {

		logger.Fatal(err).
			Msg("failed to shutdown server")
	}

	if err := db.Close(); err != nil {

		logger.Error(err).
			Msg("failed to close mysql")
	}

	if err := rabbitConn.Close(); err != nil {

		logger.Error(err).
			Msg("failed to close rabbitmq")
	}

	logger.Info().
		Msg("server exited properly")
}
