package elasticsearch

import (
	"bytes"
	"context"
	"encoding/json"
	"strconv"
	"time"

	"content-hub/internal/domain/entity"
	"content-hub/internal/logger"

	es8 "github.com/elastic/go-elasticsearch/v8"
)

type ProductRepository struct {
	client *es8.Client
}

func NewProductRepository(
	client *es8.Client,
) *ProductRepository {
	return &ProductRepository{
		client: client,
	}
}

func (r *ProductRepository) Index(
	ctx context.Context,
	product *entity.Product,
) error {

	start := time.Now()

	data, err := json.Marshal(
		product,
	)

	if err != nil {

		logger.Error(err).
			Str("service", "elasticsearch").
			Str("event", "product_marshal_failed").
			Uint64("product_id", product.ID).
			Int64(
				"duration_ms",
				time.Since(start).Milliseconds(),
			).
			Msg("failed marshal product document")

		return err
	}

	_, err = r.client.Index(
		"products",
		bytes.NewReader(data),
		r.client.Index.WithContext(ctx),
		r.client.Index.WithDocumentID(
			strconv.FormatUint(product.ID, 10),
		),
		r.client.Index.WithRefresh("true"),
	)

	if err != nil {

		logger.Error(err).
			Str("service", "elasticsearch").
			Str("event", "product_index_failed").
			Uint64("product_id", product.ID).
			Int64(
				"duration_ms",
				time.Since(start).Milliseconds(),
			).
			Msg("failed index product document")

		return err
	}

	logger.Info().
		Str("service", "elasticsearch").
		Str("event", "product_indexed").
		Uint64("product_id", product.ID).
		Dur(
			"duration",
			time.Since(start),
		).
		Int64(
			"duration_ms",
			time.Since(start).Milliseconds(),
		).
		Msg("product indexed to elasticsearch")

	return nil
}

func (r *ProductRepository) Delete(
	ctx context.Context,
	id uint64,
) error {

	start := time.Now()

	_, err := r.client.Delete(
		"products",
		strconv.FormatUint(id, 10),
		r.client.Delete.WithContext(ctx),
		r.client.Delete.WithRefresh("true"),
	)

	if err != nil {

		logger.Error(err).
			Str("service", "elasticsearch").
			Str("event", "product_delete_failed").
			Uint64("product_id", id).
			Int64(
				"duration_ms",
				time.Since(start).Milliseconds(),
			).
			Msg("failed delete product document")

		return err
	}

	logger.Info().
		Str("service", "elasticsearch").
		Str("event", "product_deleted").
		Uint64("product_id", id).
		Dur(
			"duration",
			time.Since(start),
		).
		Int64(
			"duration_ms",
			time.Since(start).Milliseconds(),
		).
		Msg("product deleted from elasticsearch")

	return nil
}
