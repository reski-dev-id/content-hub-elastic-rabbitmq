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
