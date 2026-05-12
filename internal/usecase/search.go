package usecase

import (
	"context"
	"time"

	domain "content-hub/internal/domain/usecase"
	"content-hub/internal/logger"
	esRepo "content-hub/internal/repository/elasticsearch"
)

type searchUsecase struct {
	repo *esRepo.SearchRepository
}

func NewSearchUsecase(
	r *esRepo.SearchRepository,
) domain.SearchUsecase {
	return &searchUsecase{
		repo: r,
	}
}

func (u *searchUsecase) Search(
	q string,
	searchType string,
	categoryID *uint64,
	page int,
	limit int,
) (interface{}, int64, error) {

	start := time.Now()

	index := "products"

	if searchType == "news" {
		index = "news"
	}

	logger.Info().
		Str("service", "usecase").
		Str("event", "search_started").
		Str("query", q).
		Str("type", searchType).
		Int("page", page).
		Int("limit", limit).
		Msg("starting search")

	result, err := u.repo.Search(
		context.Background(),
		index,
		q,
		categoryID,
		page,
		limit,
	)

	if err != nil {

		logger.Error(err).
			Str("service", "usecase").
			Str("event", "search_failed").
			Str("query", q).
			Str("type", searchType).
			Int("page", page).
			Int("limit", limit).
			Int64(
				"duration_ms",
				time.Since(start).Milliseconds(),
			).
			Msg("failed to search content")

		return nil, 0, err
	}

	hits := result["hits"].(map[string]interface{})

	items := hits["hits"].([]interface{})

	totalMap := hits["total"].(map[string]interface{})

	total := int64(
		totalMap["value"].(float64),
	)

	response := []map[string]interface{}{}

	for _, item := range items {

		hit := item.(map[string]interface{})

		source := hit["_source"].(map[string]interface{})

		source["score"] = hit["_score"]

		response = append(
			response,
			source,
		)
	}

	logger.Info().
		Str("service", "usecase").
		Str("event", "search_completed").
		Str("query", q).
		Str("type", searchType).
		Int("page", page).
		Int("limit", limit).
		Int64("total", total).
		Dur(
			"duration",
			time.Since(start),
		).
		Int64(
			"duration_ms",
			time.Since(start).Milliseconds(),
		).
		Msg("search executed")

	return response, total, nil
}
