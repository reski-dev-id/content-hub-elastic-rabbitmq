package usecase

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"content-hub/internal/delivery/http/response"
	"content-hub/internal/domain/entity"
	"content-hub/internal/domain/repository"
	domain "content-hub/internal/domain/usecase"
	"content-hub/internal/logger"

	esRepo "content-hub/internal/repository/elasticsearch"
)

type newsUsecase struct {
	repo   repository.NewsRepository
	esRepo *esRepo.NewsRepository
}

func NewNewsUsecase(
	r repository.NewsRepository,
	es *esRepo.NewsRepository,
) domain.NewsUsecase {
	return &newsUsecase{
		repo:   r,
		esRepo: es,
	}
}

func (u *newsUsecase) Create(
	n *entity.News,
) error {

	start := time.Now()

	payloadBytes, _ := json.Marshal(
		n,
	)

	event := &entity.OutboxEvent{
		AggregateType: "news",
		EventType:     "news_created",
		Payload:       string(payloadBytes),
		Status:        "pending",
	}

	err := u.repo.CreateWithOutbox(
		n,
		event,
	)

	if err != nil {

		logger.Error(err).
			Str("service", "usecase").
			Str("event", "news_create_failed").
			Str("title", n.Title).
			Int64(
				"duration_ms",
				time.Since(start).Milliseconds(),
			).
			Msg("failed to create news")

		return err
	}

	logger.Info().
		Str("service", "usecase").
		Str("event", "news_created").
		Uint64("news_id", n.ID).
		Str("title", n.Title).
		Dur(
			"duration",
			time.Since(start),
		).
		Int64(
			"duration_ms",
			time.Since(start).Milliseconds(),
		).
		Msg("news created")

	return nil
}

func (u *newsUsecase) GetAll(
	page,
	limit int,
	categoryID *uint64,
) ([]entity.News, int64, error) {

	start := time.Now()

	news, err := u.repo.FindAll(
		page,
		limit,
		categoryID,
	)

	if err != nil {

		logger.Error(err).
			Str("service", "usecase").
			Str("event", "news_fetch_failed").
			Int("page", page).
			Int("limit", limit).
			Int64(
				"duration_ms",
				time.Since(start).Milliseconds(),
			).
			Msg("failed to fetch news")

		return nil, 0, err
	}

	total, err := u.repo.Count(
		categoryID,
	)

	if err != nil {

		logger.Error(err).
			Str("service", "usecase").
			Str("event", "news_count_failed").
			Int64(
				"duration_ms",
				time.Since(start).Milliseconds(),
			).
			Msg("failed to count news")

		return nil, 0, err
	}

	logger.Info().
		Str("service", "usecase").
		Str("event", "news_fetched").
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
		Msg("news fetched")

	return news, total, nil
}

func (u *newsUsecase) GetByID(
	id uint64,
) (*entity.News, error) {

	start := time.Now()

	news, err := u.repo.FindByID(
		id,
	)

	if err != nil {

		logger.Error(err).
			Str("service", "usecase").
			Str("event", "news_find_by_id_failed").
			Uint64("news_id", id).
			Int64(
				"duration_ms",
				time.Since(start).Milliseconds(),
			).
			Msg("failed to get news by id")

		return nil, err
	}

	logger.Info().
		Str("service", "usecase").
		Str("event", "news_found").
		Uint64("news_id", id).
		Dur(
			"duration",
			time.Since(start),
		).
		Int64(
			"duration_ms",
			time.Since(start).Milliseconds(),
		).
		Msg("news fetched by id")

	return news, nil
}

func (u *newsUsecase) Recommend(
	id uint64,
	limit int,
) ([]entity.News, error) {

	start := time.Now()

	newsList, err := u.esRepo.Recommend(
		context.Background(),
		strconv.FormatUint(id, 10),
		limit,
	)

	if err != nil {

		logger.Error(err).
			Str("service", "usecase").
			Str("event", "news_recommendation_failed").
			Uint64("news_id", id).
			Int64(
				"duration_ms",
				time.Since(start).Milliseconds(),
			).
			Msg("failed get news recommendations")

		return nil, err
	}

	logger.Info().
		Str("service", "usecase").
		Str("event", "news_recommendations_fetched").
		Uint64("news_id", id).
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

func (u *newsUsecase) Update(
	n *entity.News,
) error {

	start := time.Now()

	payloadBytes, _ := json.Marshal(
		n,
	)

	event := &entity.OutboxEvent{
		AggregateType: "news",
		AggregateID:   n.ID,
		EventType:     "news_updated",
		Payload:       string(payloadBytes),
		Status:        "pending",
	}

	err := u.repo.UpdateWithOutbox(
		n,
		event,
	)

	if err != nil {

		logger.Error(err).
			Str("service", "usecase").
			Str("event", "news_update_failed").
			Uint64("news_id", n.ID).
			Int64(
				"duration_ms",
				time.Since(start).Milliseconds(),
			).
			Msg("failed to update news")

		return err
	}

	logger.Info().
		Str("service", "usecase").
		Str("event", "news_updated").
		Uint64("news_id", n.ID).
		Str("title", n.Title).
		Dur(
			"duration",
			time.Since(start),
		).
		Int64(
			"duration_ms",
			time.Since(start).Milliseconds(),
		).
		Msg("news updated")

	return nil
}

func (u *newsUsecase) Delete(
	id uint64,
) error {

	start := time.Now()

	payload := fmt.Sprintf(
		`{"id": %d}`,
		id,
	)

	event := &entity.OutboxEvent{
		AggregateType: "news",
		AggregateID:   id,
		EventType:     "news_deleted",
		Payload:       payload,
		Status:        "pending",
	}

	err := u.repo.DeleteWithOutbox(
		id,
		event,
	)

	if err != nil {

		logger.Error(err).
			Str("service", "usecase").
			Str("event", "news_delete_failed").
			Uint64("news_id", id).
			Int64(
				"duration_ms",
				time.Since(start).Milliseconds(),
			).
			Msg("failed to delete news")

		return err
	}

	logger.Info().
		Str("service", "usecase").
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
		Msg("news deleted")

	return nil
}

func (u *newsUsecase) Search(
	q string,
	categoryID *uint64,
	page int,
	limit int,
) (interface{}, int64, error) {

	start := time.Now()

	items, total, err := u.esRepo.Search(
		context.Background(),
		q,
		categoryID,
		page,
		limit,
	)

	if err != nil {

		logger.Error(err).
			Str("service", "usecase").
			Str("event", "news_search_failed").
			Str("query", q).
			Int("page", page).
			Int("limit", limit).
			Int64(
				"duration_ms",
				time.Since(start).Milliseconds(),
			).
			Msg("failed search news")

		return nil, 0, err
	}

	result := &response.NewsSearchResponse{
		Items:           items,
		Recommendations: []entity.News{},
	}

	if len(items) > 0 {

		recommendations, err := u.Recommend(
			items[0].ID,
			3,
		)

		if err == nil {
			result.Recommendations = recommendations
		}
	}

	logger.Info().
		Str("service", "usecase").
		Str("event", "news_search_completed").
		Str("query", q).
		Int64("total", total).
		Dur(
			"duration",
			time.Since(start),
		).
		Int64(
			"duration_ms",
			time.Since(start).Milliseconds(),
		).
		Msg("news search completed")

	return result, total, nil
}
