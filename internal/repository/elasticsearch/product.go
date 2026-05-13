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

func (r *ProductRepository) Recommend(
	ctx context.Context,
	productID string,
	limit int,
) ([]entity.Product, error) {

	start := time.Now()

	query := map[string]interface{}{
		"query": map[string]interface{}{
			"bool": map[string]interface{}{
				"must": map[string]interface{}{
					"more_like_this": map[string]interface{}{
						"fields": []string{
							"title",
							"description",
						},
						"like": []map[string]interface{}{
							{
								"_index": "products",
								"_id":    productID,
							},
						},
						"min_term_freq": 1,
						"min_doc_freq":  1,
					},
				},
				"must_not": []map[string]interface{}{
					{
						"term": map[string]interface{}{
							"_id": productID,
						},
					},
				},
			},
		},
		"size": limit,
	}

	body, err := json.Marshal(query)

	if err != nil {

		logger.Error(err).
			Str("service", "elasticsearch").
			Str("event", "recommendation_query_marshal_failed").
			Str("product_id", productID).
			Int64(
				"duration_ms",
				time.Since(start).Milliseconds(),
			).
			Msg("failed marshal recommendation query")

		return nil, err
	}

	res, err := r.client.Search(
		r.client.Search.WithContext(ctx),
		r.client.Search.WithIndex("products"),
		r.client.Search.WithBody(bytes.NewReader(body)),
	)

	if err != nil {

		logger.Error(err).
			Str("service", "elasticsearch").
			Str("event", "recommendation_search_failed").
			Str("product_id", productID).
			Int64(
				"duration_ms",
				time.Since(start).Milliseconds(),
			).
			Msg("failed search recommendations")

		return nil, err
	}

	defer res.Body.Close()

	var result struct {
		Hits struct {
			Hits []struct {
				Source entity.Product `json:"_source"`
			} `json:"hits"`
		} `json:"hits"`
	}

	err = json.NewDecoder(res.Body).Decode(&result)

	if err != nil {

		logger.Error(err).
			Str("service", "elasticsearch").
			Str("event", "recommendation_decode_failed").
			Str("product_id", productID).
			Int64(
				"duration_ms",
				time.Since(start).Milliseconds(),
			).
			Msg("failed decode recommendations")

		return nil, err
	}

	var products []entity.Product

	for _, hit := range result.Hits.Hits {

		products = append(
			products,
			hit.Source,
		)
	}

	logger.Info().
		Str("service", "elasticsearch").
		Str("event", "recommendations_fetched").
		Str("product_id", productID).
		Int("count", len(products)).
		Dur(
			"duration",
			time.Since(start),
		).
		Int64(
			"duration_ms",
			time.Since(start).Milliseconds(),
		).
		Msg("product recommendations fetched")

	return products, nil
}

func (r *ProductRepository) Search(
	ctx context.Context,
	q string,
	categoryID *uint64,
	page int,
	limit int,
) ([]entity.Product, int64, error) {

	start := time.Now()

	from := (page - 1) * limit

	must := []interface{}{
		map[string]interface{}{
			"multi_match": map[string]interface{}{
				"query": q,
				"fields": []string{
					"title^5",
				},
				"fuzziness":     "AUTO",
				"prefix_length": 1,
				"operator":      "and",
			},
		},
	}

	filter := []interface{}{}

	if categoryID != nil {

		filter = append(
			filter,
			map[string]interface{}{
				"term": map[string]interface{}{
					"category_id": *categoryID,
				},
			},
		)
	}

	query := map[string]interface{}{
		"query": map[string]interface{}{
			"bool": map[string]interface{}{
				"must":   must,
				"filter": filter,
			},
		},
		"from": from,
		"size": limit,
	}

	body, err := json.Marshal(query)

	if err != nil {

		logger.Error(err).
			Str("service", "elasticsearch").
			Str("event", "product_search_query_marshal_failed").
			Str("query", q).
			Int64(
				"duration_ms",
				time.Since(start).Milliseconds(),
			).
			Msg("failed marshal search query")

		return nil, 0, err
	}

	res, err := r.client.Search(
		r.client.Search.WithContext(ctx),
		r.client.Search.WithIndex("products"),
		r.client.Search.WithBody(bytes.NewReader(body)),
	)

	if err != nil {

		logger.Error(err).
			Str("service", "elasticsearch").
			Str("event", "product_search_failed").
			Str("query", q).
			Int64(
				"duration_ms",
				time.Since(start).Milliseconds(),
			).
			Msg("failed search products")

		return nil, 0, err
	}

	defer res.Body.Close()

	var result struct {
		Hits struct {
			Total struct {
				Value int64 `json:"value"`
			} `json:"total"`
			Hits []struct {
				Source entity.Product `json:"_source"`
			} `json:"hits"`
		} `json:"hits"`
	}

	err = json.NewDecoder(res.Body).Decode(&result)

	if err != nil {

		logger.Error(err).
			Str("service", "elasticsearch").
			Str("event", "product_search_decode_failed").
			Str("query", q).
			Int64(
				"duration_ms",
				time.Since(start).Milliseconds(),
			).
			Msg("failed decode search result")

		return nil, 0, err
	}

	var products []entity.Product

	for _, hit := range result.Hits.Hits {

		products = append(
			products,
			hit.Source,
		)
	}

	logger.Info().
		Str("service", "elasticsearch").
		Str("event", "product_search_completed").
		Str("query", q).
		Int("count", len(products)).
		Int64("total", result.Hits.Total.Value).
		Dur(
			"duration",
			time.Since(start),
		).
		Int64(
			"duration_ms",
			time.Since(start).Milliseconds(),
		).
		Msg("product search completed")

	return products, result.Hits.Total.Value, nil
}
