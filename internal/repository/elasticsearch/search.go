package elasticsearch

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"time"

	"content-hub/internal/logger"

	es8 "github.com/elastic/go-elasticsearch/v8"
)

type SearchRepository struct {
	client *es8.Client
}

func NewSearchRepository(
	client *es8.Client,
) *SearchRepository {
	return &SearchRepository{
		client: client,
	}
}

func (r *SearchRepository) Search(
	ctx context.Context,
	index string,
	q string,
	categoryID *uint64,
	page int,
	limit int,
) (map[string]interface{}, error) {

	start := time.Now()

	from := (page - 1) * limit

	logger.Info().
		Str("service", "elasticsearch").
		Str("event", "search_started").
		Str("index", index).
		Str("query", q).
		Int("page", page).
		Int("limit", limit).
		Msg("starting elasticsearch search")

	must := []map[string]interface{}{
		{
			"multi_match": map[string]interface{}{
				"query": q,
				"fields": []string{
					"title",
					"description",
					"content",
				},
			},
		},
	}

	if categoryID != nil {

		logger.Info().
			Str("service", "elasticsearch").
			Str("event", "search_category_filter").
			Str("index", index).
			Uint64("category_id", *categoryID).
			Msg("applying category filter")

		must = append(
			must,
			map[string]interface{}{
				"term": map[string]interface{}{
					"category_id": *categoryID,
				},
			},
		)
	}

	query := map[string]interface{}{
		"from": from,
		"size": limit,
		"query": map[string]interface{}{
			"bool": map[string]interface{}{
				"must": must,
			},
		},
	}

	body, err := json.Marshal(
		query,
	)

	if err != nil {

		logger.Error(err).
			Str("service", "elasticsearch").
			Str("event", "search_query_marshal_failed").
			Str("index", index).
			Int64(
				"duration_ms",
				time.Since(start).Milliseconds(),
			).
			Msg("failed marshal elasticsearch query")

		return nil, err
	}

	res, err := r.client.Search(
		r.client.Search.WithContext(ctx),
		r.client.Search.WithIndex(index),
		r.client.Search.WithBody(
			bytes.NewReader(body),
		),
		r.client.Search.WithTrackTotalHits(true),
	)

	if err != nil {

		logger.Error(err).
			Str("service", "elasticsearch").
			Str("event", "search_request_failed").
			Str("index", index).
			Str("query", q).
			Int64(
				"duration_ms",
				time.Since(start).Milliseconds(),
			).
			Msg("failed execute elasticsearch search")

		return nil, err
	}

	defer res.Body.Close()

	if res.IsError() {

		logger.Error(nil).
			Str("service", "elasticsearch").
			Str("event", "search_response_error").
			Str("index", index).
			Str("query", q).
			Str("response", res.String()).
			Int64(
				"duration_ms",
				time.Since(start).Milliseconds(),
			).
			Msg("elasticsearch returned error response")

		return nil, fmt.Errorf(
			res.String(),
		)
	}

	var result map[string]interface{}

	err = json.NewDecoder(
		res.Body,
	).Decode(
		&result,
	)

	if err != nil {

		logger.Error(err).
			Str("service", "elasticsearch").
			Str("event", "search_decode_failed").
			Str("index", index).
			Int64(
				"duration_ms",
				time.Since(start).Milliseconds(),
			).
			Msg("failed decode elasticsearch response")

		return nil, err
	}

	hits := result["hits"].(map[string]interface{})

	totalMap := hits["total"].(map[string]interface{})

	total := int64(
		totalMap["value"].(float64),
	)

	logger.Info().
		Str("service", "elasticsearch").
		Str("event", "search_completed").
		Str("index", index).
		Str("query", q).
		Int("page", page).
		Int("limit", limit).
		Int64("total_hits", total).
		Dur(
			"duration",
			time.Since(start),
		).
		Int64(
			"duration_ms",
			time.Since(start).Milliseconds(),
		).
		Msg("elasticsearch search completed")

	return result, nil
}
