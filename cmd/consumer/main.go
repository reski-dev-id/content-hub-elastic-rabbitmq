package main

import (
	"context"
	"encoding/json"
	"time"

	"content-hub/config"
	"content-hub/internal/domain/entity"
	esInfra "content-hub/internal/infrastructure/elasticsearch"
	"content-hub/internal/infrastructure/rabbitmq"
	"content-hub/internal/logger"
	esRepo "content-hub/internal/repository/elasticsearch"
)

func main() {

	start := time.Now()

	logger.Init()

	cfg := config.Load()

	esClient, err := esInfra.NewClient(
		cfg.ElasticUrl,
	)

	if err != nil {

		logger.Fatal(err).
			Str("service", "consumer").
			Str("event", "elasticsearch_connection_failed").
			Int64(
				"duration_ms",
				time.Since(start).Milliseconds(),
			).
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
			Str("service", "consumer").
			Str("event", "rabbitmq_connection_failed").
			Int64(
				"duration_ms",
				time.Since(start).Milliseconds(),
			).
			Msg("failed to connect rabbitmq")
	}

	consumer, err := rabbitmq.NewConsumer(
		rmqConn,
	)

	if err != nil {

		logger.Fatal(err).
			Str("service", "consumer").
			Str("event", "rabbitmq_consumer_create_failed").
			Int64(
				"duration_ms",
				time.Since(start).Milliseconds(),
			).
			Msg("failed to create rabbitmq consumer")
	}

	productMsgs, err := consumer.Consume(
		"product",
	)

	if err != nil {

		logger.Fatal(err).
			Str("service", "consumer").
			Str("event", "product_consume_failed").
			Str("queue", "product").
			Int64(
				"duration_ms",
				time.Since(start).Milliseconds(),
			).
			Msg("failed to consume queue")
	}

	newsMsgs, err := consumer.Consume(
		"news",
	)

	if err != nil {

		logger.Fatal(err).
			Str("service", "consumer").
			Str("event", "news_consume_failed").
			Str("queue", "news").
			Int64(
				"duration_ms",
				time.Since(start).Milliseconds(),
			).
			Msg("failed to consume queue")
	}

	logger.Info().
		Str("service", "consumer").
		Str("event", "consumer_started").
		Dur(
			"duration",
			time.Since(start),
		).
		Int64(
			"duration_ms",
			time.Since(start).Milliseconds(),
		).
		Msg("rabbitmq consumer started")

	// PRODUCT CONSUMER

	go func() {

		for msg := range productMsgs {

			processStart := time.Now()

			logger.Info().
				Str("service", "consumer").
				Str("event", "product_message_received").
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
					Str("service", "consumer").
					Str("event", "product_unmarshal_failed").
					Str("queue", "product").
					Bytes("payload", msg.Body).
					Int64(
						"duration_ms",
						time.Since(processStart).Milliseconds(),
					).
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
					Str("service", "consumer").
					Str("event", "product_delete_received").
					Str("queue", "product").
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
						Str("service", "consumer").
						Str("event", "product_delete_missing_id").
						Str("queue", "product").
						Bytes("payload", msg.Body).
						Int64(
							"duration_ms",
							time.Since(processStart).Milliseconds(),
						).
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
						Str("service", "consumer").
						Str("event", "product_delete_failed").
						Str("queue", "product").
						Uint64("product_id", id).
						Int64(
							"duration_ms",
							time.Since(processStart).Milliseconds(),
						).
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
					Str("service", "consumer").
					Str("event", "product_deleted").
					Str("queue", "product").
					Uint64("product_id", id).
					Dur(
						"duration",
						time.Since(processStart),
					).
					Int64(
						"duration_ms",
						time.Since(processStart).Milliseconds(),
					).
					Msg("product deleted from elasticsearch")

				_ = msg.Ack(false)

				continue
			}

			logger.Info().
				Str("service", "consumer").
				Str("event", "product_indexing").
				Str("queue", "product").
				Uint64("product_id", product.ID).
				Msg("indexing product to elasticsearch")

			err = productES.Index(
				context.Background(),
				&product,
			)

			if err != nil {

				logger.Error(err).
					Str("service", "consumer").
					Str("event", "product_index_failed").
					Str("queue", "product").
					Uint64("product_id", product.ID).
					Int64(
						"duration_ms",
						time.Since(processStart).Milliseconds(),
					).
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
				Str("service", "consumer").
				Str("event", "product_indexed").
				Str("queue", "product").
				Uint64("product_id", product.ID).
				Dur(
					"duration",
					time.Since(processStart),
				).
				Int64(
					"duration_ms",
					time.Since(processStart).Milliseconds(),
				).
				Msg("product indexed successfully")

			_ = msg.Ack(false)
		}
	}()

	// NEWS CONSUMER

	go func() {

		for msg := range newsMsgs {

			processStart := time.Now()

			logger.Info().
				Str("service", "consumer").
				Str("event", "news_message_received").
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
					Str("service", "consumer").
					Str("event", "news_unmarshal_failed").
					Str("queue", "news").
					Bytes("payload", msg.Body).
					Int64(
						"duration_ms",
						time.Since(processStart).Milliseconds(),
					).
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
					Str("service", "consumer").
					Str("event", "news_delete_received").
					Str("queue", "news").
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
						Str("service", "consumer").
						Str("event", "news_delete_missing_id").
						Str("queue", "news").
						Bytes("payload", msg.Body).
						Int64(
							"duration_ms",
							time.Since(processStart).Milliseconds(),
						).
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
						Str("service", "consumer").
						Str("event", "news_delete_failed").
						Str("queue", "news").
						Uint64("news_id", id).
						Int64(
							"duration_ms",
							time.Since(processStart).Milliseconds(),
						).
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
					Str("service", "consumer").
					Str("event", "news_deleted").
					Str("queue", "news").
					Uint64("news_id", id).
					Dur(
						"duration",
						time.Since(processStart),
					).
					Int64(
						"duration_ms",
						time.Since(processStart).Milliseconds(),
					).
					Msg("news deleted from elasticsearch")

				_ = msg.Ack(false)

				continue
			}

			logger.Info().
				Str("service", "consumer").
				Str("event", "news_indexing").
				Str("queue", "news").
				Uint64("news_id", news.ID).
				Msg("indexing news to elasticsearch")

			err = newsES.Index(
				context.Background(),
				&news,
			)

			if err != nil {

				logger.Error(err).
					Str("service", "consumer").
					Str("event", "news_index_failed").
					Str("queue", "news").
					Uint64("news_id", news.ID).
					Int64(
						"duration_ms",
						time.Since(processStart).Milliseconds(),
					).
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
				Str("service", "consumer").
				Str("event", "news_indexed").
				Str("queue", "news").
				Uint64("news_id", news.ID).
				Dur(
					"duration",
					time.Since(processStart),
				).
				Int64(
					"duration_ms",
					time.Since(processStart).Milliseconds(),
				).
				Msg("news indexed successfully")

			_ = msg.Ack(false)
		}
	}()

	select {}
}
