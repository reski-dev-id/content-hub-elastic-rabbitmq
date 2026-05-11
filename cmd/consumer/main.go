package main

import (
	"context"
	"encoding/json"

	"content-hub/config"
	"content-hub/internal/domain/entity"
	esInfra "content-hub/internal/infrastructure/elasticsearch"
	"content-hub/internal/infrastructure/rabbitmq"
	"content-hub/internal/logger"
	esRepo "content-hub/internal/repository/elasticsearch"
)

func main() {

	logger.Init()

	cfg := config.Load()

	esClient, err := esInfra.NewClient(
		cfg.ElasticUrl,
	)

	if err != nil {

		logger.Log.Fatal().
			Err(err).
			Msg("failed connect elasticsearch")
	}

	productES := esRepo.NewProductRepository(
		esClient,
	)

	newsES := esRepo.NewNewsRepository(
		esClient,
	)

	rmqConn, err := rabbitmq.NewRabbitMQ(
		cfg.RabbitMQUrl,
	)

	if err != nil {

		logger.Log.Fatal().
			Err(err).
			Msg("failed connect rabbitmq")
	}

	consumer, err := rabbitmq.NewConsumer(
		rmqConn,
	)

	if err != nil {

		logger.Log.Fatal().
			Err(err).
			Msg("failed create rabbitmq consumer")
	}

	productMsgs, err := consumer.Consume(
		"product",
	)

	if err != nil {

		logger.Log.Fatal().
			Err(err).
			Str("queue", "product").
			Msg("failed consume product queue")
	}

	newsMsgs, err := consumer.Consume(
		"news",
	)

	if err != nil {

		logger.Log.Fatal().
			Err(err).
			Str("queue", "news").
			Msg("failed consume news queue")
	}

	logger.Log.Info().
		Msg("rabbitmq consumer started")

	go func() {

		for msg := range productMsgs {

			logger.Log.Info().
				Str("queue", "product").
				Msg("product message received")

			var product entity.Product

			err := json.Unmarshal(
				msg.Body,
				&product,
			)

			if err != nil {

				logger.Log.Error().
					Err(err).
					Bytes("payload", msg.Body).
					Msg("failed unmarshal product payload")

				continue
			}

			if product.ID == 0 {

				logger.Log.Info().
					Bytes("payload", msg.Body).
					Msg("delete product event received")

				var payload map[string]interface{}

				_ = json.Unmarshal(
					msg.Body,
					&payload,
				)

				rawID, ok := payload["id"]

				if !ok {

					logger.Log.Error().
						Bytes("payload", msg.Body).
						Msg("missing id in delete product payload")

					continue
				}

				id := uint64(
					rawID.(float64),
				)

				err = productES.Delete(
					context.Background(),
					id,
				)

				if err != nil {

					logger.Log.Error().
						Err(err).
						Uint64("product_id", id).
						Msg("failed delete product from elasticsearch")

					continue
				}

				logger.Log.Info().
					Uint64("product_id", id).
					Msg("product deleted from elasticsearch")

				continue
			}

			logger.Log.Info().
				Uint64("product_id", product.ID).
				Msg("indexing product to elasticsearch")

			err = productES.Index(
				context.Background(),
				&product,
			)

			if err != nil {

				logger.Log.Error().
					Err(err).
					Uint64("product_id", product.ID).
					Msg("failed index product")

				continue
			}

			logger.Log.Info().
				Uint64("product_id", product.ID).
				Msg("product indexed successfully")
		}
	}()

	go func() {

		for msg := range newsMsgs {

			logger.Log.Info().
				Str("queue", "news").
				Msg("news message received")

			var news entity.News

			err := json.Unmarshal(
				msg.Body,
				&news,
			)

			if err != nil {

				logger.Log.Error().
					Err(err).
					Bytes("payload", msg.Body).
					Msg("failed unmarshal news payload")

				continue
			}

			if news.ID == 0 {

				logger.Log.Info().
					Bytes("payload", msg.Body).
					Msg("delete news event received")

				var payload map[string]interface{}

				_ = json.Unmarshal(
					msg.Body,
					&payload,
				)

				rawID, ok := payload["id"]

				if !ok {

					logger.Log.Error().
						Bytes("payload", msg.Body).
						Msg("missing id in delete news payload")

					continue
				}

				id := uint64(
					rawID.(float64),
				)

				err = newsES.Delete(
					context.Background(),
					id,
				)

				if err != nil {

					logger.Log.Error().
						Err(err).
						Uint64("news_id", id).
						Msg("failed delete news from elasticsearch")

					continue
				}

				logger.Log.Info().
					Uint64("news_id", id).
					Msg("news deleted from elasticsearch")

				continue
			}

			logger.Log.Info().
				Uint64("news_id", news.ID).
				Msg("indexing news to elasticsearch")

			err = newsES.Index(
				context.Background(),
				&news,
			)

			if err != nil {

				logger.Log.Error().
					Err(err).
					Uint64("news_id", news.ID).
					Msg("failed index news")

				continue
			}

			logger.Log.Info().
				Uint64("news_id", news.ID).
				Msg("news indexed successfully")
		}
	}()

	select {}
}
