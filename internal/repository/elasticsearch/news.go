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

type NewsRepository struct {
	client *es8.Client
}

func NewNewsRepository(
	client *es8.Client,
) *NewsRepository {
	return &NewsRepository{
		client: client,
	}
}

func (r *NewsRepository) Index(
	ctx context.Context,
	news *entity.News,
) error {

	start := time.Now()

	data, err := json.Marshal(
		news,
	)

	if err != nil {

		logger.Error(err).
			Str("service", "elasticsearch").
			Str("event", "news_marshal_failed").
			Uint64("news_id", news.ID).
			Int64(
				"duration_ms",
				time.Since(start).Milliseconds(),
			).
			Msg("failed marshal news document")

		return err
	}

	_, err = r.client.Index(
		"news",
		bytes.NewReader(data),
		r.client.Index.WithContext(ctx),
		r.client.Index.WithDocumentID(
			strconv.FormatUint(news.ID, 10),
		),
		r.client.Index.WithRefresh("true"),
	)

	if err != nil {

		logger.Error(err).
			Str("service", "elasticsearch").
			Str("event", "news_index_failed").
			Uint64("news_id", news.ID).
			Int64(
				"duration_ms",
				time.Since(start).Milliseconds(),
			).
			Msg("failed index news document")

		return err
	}

	logger.Info().
		Str("service", "elasticsearch").
		Str("event", "news_indexed").
		Uint64("news_id", news.ID).
		Dur(
			"duration",
			time.Since(start),
		).
		Int64(
			"duration_ms",
			time.Since(start).Milliseconds(),
		).
		Msg("news indexed to elasticsearch")

	return nil
}

func (r *NewsRepository) Delete(
	ctx context.Context,
	id uint64,
) error {

	start := time.Now()

	_, err := r.client.Delete(
		"news",
		strconv.FormatUint(id, 10),
		r.client.Delete.WithContext(ctx),
		r.client.Delete.WithRefresh("true"),
	)

	if err != nil {

		logger.Error(err).
			Str("service", "elasticsearch").
			Str("event", "news_delete_failed").
			Uint64("news_id", id).
			Int64(
				"duration_ms",
				time.Since(start).Milliseconds(),
			).
			Msg("failed delete news document")

		return err
	}

	logger.Info().
		Str("service", "elasticsearch").
		Str("event", "news_deleted").
		Uint64("news_id", id).
		Dur(
			"duration",
			time.Since(start),
		).
		Int64(
			"duration_ms",
			time.Since(start).Milliseconds(),
		).
		Msg("news deleted from elasticsearch")

	return nil
}

func (r *NewsRepository) Recommend(
	ctx context.Context,
	newsID string,
	limit int,
) ([]entity.News, error) {

	start := time.Now()

	query := map[string]interface{}{
		"query": map[string]interface{}{
			"bool": map[string]interface{}{
				"must": map[string]interface{}{
					"more_like_this": map[string]interface{}{
						"fields": []string{
							"title",
							"content",
						},
						"like": []map[string]interface{}{
							{
								"_index": "news",
								"_id":    newsID,
							},
						},
						"min_term_freq": 1,
						"min_doc_freq":  1,
					},
				},
				"must_not": []map[string]interface{}{
					{
						"term": map[string]interface{}{
							"_id": newsID,
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
			Str("event", "news_recommendation_query_marshal_failed").
			Str("news_id", newsID).
			Int64(
				"duration_ms",
				time.Since(start).Milliseconds(),
			).
			Msg("failed marshal news recommendation query")

		return nil, err
	}

	res, err := r.client.Search(
		r.client.Search.WithContext(ctx),
		r.client.Search.WithIndex("news"),
		r.client.Search.WithBody(bytes.NewReader(body)),
	)

	if err != nil {

		logger.Error(err).
			Str("service", "elasticsearch").
			Str("event", "news_recommendation_search_failed").
			Str("news_id", newsID).
			Int64(
				"duration_ms",
				time.Since(start).Milliseconds(),
			).
			Msg("failed search news recommendations")

		return nil, err
	}

	defer res.Body.Close()

	var result struct {
		Hits struct {
			Hits []struct {
				Source entity.News `json:"_source"`
			} `json:"hits"`
		} `json:"hits"`
	}

	err = json.NewDecoder(res.Body).Decode(&result)

	if err != nil {

		logger.Error(err).
			Str("service", "elasticsearch").
			Str("event", "news_recommendation_decode_failed").
			Str("news_id", newsID).
			Int64(
				"duration_ms",
				time.Since(start).Milliseconds(),
			).
			Msg("failed decode news recommendations")

		return nil, err
	}

	var newsList []entity.News

	for _, hit := range result.Hits.Hits {

		newsList = append(
			newsList,
			hit.Source,
		)
	}

	logger.Info().
		Str("service", "elasticsearch").
		Str("event", "news_recommendations_fetched").
		Str("news_id", newsID).
		Int("count", len(newsList)).
		Dur(
			"duration",
			time.Since(start),
		).
		Int64(
			"duration_ms",
			time.Since(start).Milliseconds(),
		).
		Msg("news recommendations fetched")

	return newsList, nil
}

func (r *NewsRepository) Search(
	ctx context.Context,
	q string,
	categoryID *uint64,
	page int,
	limit int,
) ([]entity.News, int64, error) {

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
			Str("event", "news_search_query_marshal_failed").
			Str("query", q).
			Int64(
				"duration_ms",
				time.Since(start).Milliseconds(),
			).
			Msg("failed marshal news search query")

		return nil, 0, err
	}

	res, err := r.client.Search(
		r.client.Search.WithContext(ctx),
		r.client.Search.WithIndex("news"),
		r.client.Search.WithBody(bytes.NewReader(body)),
	)

	if err != nil {

		logger.Error(err).
			Str("service", "elasticsearch").
			Str("event", "news_search_failed").
			Str("query", q).
			Int64(
				"duration_ms",
				time.Since(start).Milliseconds(),
			).
			Msg("failed search news")

		return nil, 0, err
	}

	defer res.Body.Close()

	var result struct {
		Hits struct {
			Total struct {
				Value int64 `json:"value"`
			} `json:"total"`
			Hits []struct {
				Source entity.News `json:"_source"`
			} `json:"hits"`
		} `json:"hits"`
	}

	err = json.NewDecoder(res.Body).Decode(&result)

	if err != nil {

		logger.Error(err).
			Str("service", "elasticsearch").
			Str("event", "news_search_decode_failed").
			Str("query", q).
			Int64(
				"duration_ms",
				time.Since(start).Milliseconds(),
			).
			Msg("failed decode news search result")

		return nil, 0, err
	}

	var newsList []entity.News

	for _, hit := range result.Hits.Hits {

		newsList = append(
			newsList,
			hit.Source,
		)
	}

	logger.Info().
		Str("service", "elasticsearch").
		Str("event", "news_search_completed").
		Str("query", q).
		Int("count", len(newsList)).
		Int64("total", result.Hits.Total.Value).
		Dur(
			"duration",
			time.Since(start),
		).
		Int64(
			"duration_ms",
			time.Since(start).Milliseconds(),
		).
		Msg("news search completed")

	return newsList, result.Hits.Total.Value, nil
}
