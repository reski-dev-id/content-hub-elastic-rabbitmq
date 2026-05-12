package elasticsearch

import (
	"bytes"
	"context"
	"encoding/json"
	"strconv"

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
	return &ProductRepository{client}
}

func (r *ProductRepository) Index(
	ctx context.Context,
	product *entity.Product,
) error {

	data, err := json.Marshal(product)

	if err != nil {

		logger.Error(err).
			Uint64("product_id", product.ID).
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
			Uint64("product_id", product.ID).
			Msg("failed index product document")

		return err
	}

	logger.Info().
		Uint64("product_id", product.ID).
		Msg("product indexed to elasticsearch")

	return nil
}

func (r *ProductRepository) Delete(
	ctx context.Context,
	id uint64,
) error {

	_, err := r.client.Delete(
		"products",
		strconv.FormatUint(id, 10),
		r.client.Delete.WithContext(ctx),
		r.client.Delete.WithRefresh("true"),
	)

	if err != nil {

		logger.Error(err).
			Uint64("product_id", id).
			Msg("failed delete product document")

		return err
	}

	logger.Info().
		Uint64("product_id", id).
		Msg("product deleted from elasticsearch")

	return nil
}
