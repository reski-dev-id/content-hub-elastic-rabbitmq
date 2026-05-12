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

type NewsRepository struct {
	client *es8.Client
}

func NewNewsRepository(
	client *es8.Client,
) *NewsRepository {
	return &NewsRepository{client}
}

func (r *NewsRepository) Index(
	ctx context.Context,
	news *entity.News,
) error {

	data, err := json.Marshal(news)

	if err != nil {

		logger.Error(err).
			Uint64("news_id", news.ID).
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
			Uint64("news_id", news.ID).
			Msg("failed index news document")

		return err
	}

	logger.Info().
		Uint64("news_id", news.ID).
		Msg("news indexed to elasticsearch")

	return nil
}

func (r *NewsRepository) Delete(
	ctx context.Context,
	id uint64,
) error {

	_, err := r.client.Delete(
		"news",
		strconv.FormatUint(id, 10),
		r.client.Delete.WithContext(ctx),
		r.client.Delete.WithRefresh("true"),
	)

	if err != nil {

		logger.Error(err).
			Uint64("news_id", id).
			Msg("failed delete news document")

		return err
	}

	logger.Info().
		Uint64("news_id", id).
		Msg("news deleted from elasticsearch")

	return nil
}
