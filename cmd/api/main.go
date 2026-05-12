package main

// @title Content Hub API
// @version 1.0
// @description Product and News API with Elasticsearch Search
// @host localhost:8080
// @BasePath /

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

	start := time.Now()

	logger.Init()

	cfg := config.Load()

	db, err := mysql.NewDB(
		cfg.DBUrl,
	)

	if err != nil {

		logger.Fatal(err).
			Str("service", "api").
			Str("event", "mysql_connection_failed").
			Int64(
				"duration_ms",
				time.Since(start).Milliseconds(),
			).
			Msg("failed to connect mysql")
	}

	esClient, err := esInfra.NewClient(
		cfg.ElasticUrl,
	)

	if err != nil {

		logger.Fatal(err).
			Str("service", "api").
			Str("event", "elasticsearch_connection_failed").
			Int64(
				"duration_ms",
				time.Since(start).Milliseconds(),
			).
			Msg("failed to connect elasticsearch")
	}

	rabbitConn, err := rabbitInfra.NewRabbitMQ(
		cfg.RabbitMQUrl,
	)

	if err != nil {

		logger.Fatal(err).
			Str("service", "api").
			Str("event", "rabbitmq_connection_failed").
			Int64(
				"duration_ms",
				time.Since(start).Milliseconds(),
			).
			Msg("failed to connect rabbitmq")
	}

	healthHandler := handler.NewHealthHandler(
		db,
		esClient,
		rabbitConn,
	)

	// Product

	productRepo := mysqlRepo.NewProductRepository(
		db,
	)

	productUC := usecase.NewProductUsecase(
		productRepo,
	)

	productHandler := handler.NewProductHandler(
		productUC,
	)

	// News

	newsRepo := mysqlRepo.NewNewsRepository(
		db,
	)

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

	logger.Info().
		Str("service", "api").
		Str("event", "bootstrap_completed").
		Dur(
			"duration",
			time.Since(start),
		).
		Int64(
			"duration_ms",
			time.Since(start).Milliseconds(),
		).
		Msg("application bootstrap completed")

	go func() {

		logger.Info().
			Str("service", "api").
			Str("event", "server_started").
			Str("port", cfg.AppPort).
			Msg("server running")

		if err := server.ListenAndServe(); err != nil &&
			err != http.ErrServerClosed {

			logger.Fatal(err).
				Str("service", "api").
				Str("event", "server_start_failed").
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

	shutdownStart := time.Now()

	logger.Info().
		Str("service", "api").
		Str("event", "server_shutdown_started").
		Msg("shutting down server")

	ctx, cancel := context.WithTimeout(
		context.Background(),
		10*time.Second,
	)

	defer cancel()

	if err := server.Shutdown(ctx); err != nil {

		logger.Fatal(err).
			Str("service", "api").
			Str("event", "server_shutdown_failed").
			Int64(
				"duration_ms",
				time.Since(shutdownStart).Milliseconds(),
			).
			Msg("failed to shutdown server")
	}

	if err := db.Close(); err != nil {

		logger.Error(err).
			Str("service", "api").
			Str("event", "mysql_close_failed").
			Msg("failed to close mysql")
	}

	if err := rabbitConn.Close(); err != nil {

		logger.Error(err).
			Str("service", "api").
			Str("event", "rabbitmq_close_failed").
			Msg("failed to close rabbitmq")
	}

	logger.Info().
		Str("service", "api").
		Str("event", "server_stopped").
		Dur(
			"duration",
			time.Since(shutdownStart),
		).
		Int64(
			"duration_ms",
			time.Since(shutdownStart).Milliseconds(),
		).
		Msg("server exited properly")
}
