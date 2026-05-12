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

		logger.Fatal(err).
			Msg("failed to connect elasticsearch")
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

		logger.Fatal(err).
			Msg("failed to connect rabbitmq")
	}

	consumer, err := rabbitmq.NewConsumer(
		rmqConn,
	)

	if err != nil {

		logger.Fatal(err).
			Msg("failed to create rabbitmq consumer")
	}

	productMsgs, err := consumer.Consume(
		"product",
	)

	if err != nil {

		logger.Fatal(err).
			Str("queue", "product").
			Msg("failed to consume queue")
	}

	newsMsgs, err := consumer.Consume(
		"news",
	)

	if err != nil {

		logger.Fatal(err).
			Str("queue", "news").
			Msg("failed to consume queue")
	}

	logger.Info().
		Msg("rabbitmq consumer started")

	// PRODUCT CONSUMER

	go func() {

		for msg := range productMsgs {

			logger.Info().
				Str("queue", "product").
				Int(
					"retry_count",
					rabbitmq.GetRetryCount(msg),
				).
				Msg("message received")

			var product entity.Product

			err := json.Unmarshal(
				msg.Body,
				&product,
			)

			if err != nil {

				logger.Error(err).
					Bytes("payload", msg.Body).
					Msg("failed to unmarshal product payload")

				_ = rabbitmq.RetryMessage(
					consumer.Channel(),
					"product",
					msg,
				)

				_ = msg.Ack(false)

				continue
			}

			// DELETE EVENT

			if product.ID == 0 {

				logger.Info().
					Bytes("payload", msg.Body).
					Msg("delete product event received")

				var payload map[string]interface{}

				_ = json.Unmarshal(
					msg.Body,
					&payload,
				)

				rawID, ok := payload["id"]

				if !ok {

					logger.Error(nil).
						Bytes("payload", msg.Body).
						Msg("missing id in delete product payload")

					_ = rabbitmq.RetryMessage(
						consumer.Channel(),
						"product",
						msg,
					)

					_ = msg.Ack(false)

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

					logger.Error(err).
						Uint64("product_id", id).
						Msg("failed to delete product from elasticsearch")

					_ = rabbitmq.RetryMessage(
						consumer.Channel(),
						"product",
						msg,
					)

					_ = msg.Ack(false)

					continue
				}

				logger.Info().
					Uint64("product_id", id).
					Msg("product deleted from elasticsearch")

				_ = msg.Ack(false)

				continue
			}

			logger.Info().
				Uint64("product_id", product.ID).
				Msg("indexing product to elasticsearch")

			err = productES.Index(
				context.Background(),
				&product,
			)

			if err != nil {

				logger.Error(err).
					Uint64("product_id", product.ID).
					Msg("failed to index product")

				_ = rabbitmq.RetryMessage(
					consumer.Channel(),
					"product",
					msg,
				)

				_ = msg.Ack(false)

				continue
			}

			logger.Info().
				Uint64("product_id", product.ID).
				Msg("product indexed successfully")

			_ = msg.Ack(false)
		}
	}()

	// NEWS CONSUMER

	go func() {

		for msg := range newsMsgs {

			logger.Info().
				Str("queue", "news").
				Int(
					"retry_count",
					rabbitmq.GetRetryCount(msg),
				).
				Msg("message received")

			var news entity.News

			err := json.Unmarshal(
				msg.Body,
				&news,
			)

			if err != nil {

				logger.Error(err).
					Bytes("payload", msg.Body).
					Msg("failed to unmarshal news payload")

				_ = rabbitmq.RetryMessage(
					consumer.Channel(),
					"news",
					msg,
				)

				_ = msg.Ack(false)

				continue
			}

			// DELETE EVENT

			if news.ID == 0 {

				logger.Info().
					Bytes("payload", msg.Body).
					Msg("delete news event received")

				var payload map[string]interface{}

				_ = json.Unmarshal(
					msg.Body,
					&payload,
				)

				rawID, ok := payload["id"]

				if !ok {

					logger.Error(nil).
						Bytes("payload", msg.Body).
						Msg("missing id in delete news payload")

					_ = rabbitmq.RetryMessage(
						consumer.Channel(),
						"news",
						msg,
					)

					_ = msg.Ack(false)

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

					logger.Error(err).
						Uint64("news_id", id).
						Msg("failed to delete news from elasticsearch")

					_ = rabbitmq.RetryMessage(
						consumer.Channel(),
						"news",
						msg,
					)

					_ = msg.Ack(false)

					continue
				}

				logger.Info().
					Uint64("news_id", id).
					Msg("news deleted from elasticsearch")

				_ = msg.Ack(false)

				continue
			}

			logger.Info().
				Uint64("news_id", news.ID).
				Msg("indexing news to elasticsearch")

			err = newsES.Index(
				context.Background(),
				&news,
			)

			if err != nil {

				logger.Error(err).
					Uint64("news_id", news.ID).
					Msg("failed to index news")

				_ = rabbitmq.RetryMessage(
					consumer.Channel(),
					"news",
					msg,
				)

				_ = msg.Ack(false)

				continue
			}

			logger.Info().
				Uint64("news_id", news.ID).
				Msg("news indexed successfully")

			_ = msg.Ack(false)
		}
	}()

	select {}
}
