package main

import (
	"context"
	"encoding/json"
	"log"

	"content-hub/config"
	"content-hub/internal/domain/entity"
	esInfra "content-hub/internal/infrastructure/elasticsearch"
	"content-hub/internal/infrastructure/rabbitmq"
	esRepo "content-hub/internal/repository/elasticsearch"
)

func main() {
	cfg := config.Load()

	esClient, err := esInfra.NewClient(cfg.ElasticUrl)
	if err != nil {
		log.Fatal(err)
	}

	productES := esRepo.NewProductRepository(esClient)
	newsES := esRepo.NewNewsRepository(esClient)

	rmqConn, err := rabbitmq.NewRabbitMQ(
		cfg.RabbitMQUrl,
	)

	if err != nil {
		log.Fatal(err)
	}

	consumer, err := rabbitmq.NewConsumer(rmqConn)
	if err != nil {
		log.Fatal(err)
	}

	productMsgs, err := consumer.Consume("product")
	if err != nil {
		log.Fatal(err)
	}

	newsMsgs, err := consumer.Consume("news")
	if err != nil {
		log.Fatal(err)
	}

	log.Println("rabbitmq consumer started")

	go func() {
		for msg := range productMsgs {
			var product entity.Product

			err := json.Unmarshal(msg.Body, &product)
			if err != nil {
				log.Println(err)
				continue
			}

			if product.ID == 0 {
				log.Println("delete product", string(msg.Body))

				var payload map[string]interface{}

				_ = json.Unmarshal(msg.Body, &payload)

				rawID, ok := payload["id"]
				if !ok {
					log.Println("missing id in delete payload")
					continue
				}

				id := uint64(rawID.(float64))

				err = productES.Delete(
					context.Background(),
					id,
				)

				if err != nil {
					log.Println(err)
				}

				continue
			}

			log.Println("index product:", product.ID)

			err = productES.Index(
				context.Background(),
				&product,
			)

			if err != nil {
				log.Println(err)
			}
		}
	}()

	go func() {
		for msg := range newsMsgs {
			var news entity.News

			err := json.Unmarshal(msg.Body, &news)
			if err != nil {
				log.Println(err)
				continue
			}

			if news.ID == 0 {
				log.Println("delete news", string(msg.Body))

				var payload map[string]interface{}

				_ = json.Unmarshal(msg.Body, &payload)

				rawID, ok := payload["id"]
				if !ok {
					log.Println("missing id in delete payload")
					continue
				}

				id := uint64(rawID.(float64))

				err = newsES.Delete(
					context.Background(),
					id,
				)

				if err != nil {
					log.Println(err)
				}

				continue
			}

			log.Println("index news:", news.ID)

			err = newsES.Index(
				context.Background(),
				&news,
			)

			if err != nil {
				log.Println(err)
			}
		}
	}()

	select {}
}
