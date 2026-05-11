package usecase

import (
	"encoding/json"
	"fmt"

	"content-hub/internal/domain/entity"
	"content-hub/internal/domain/repository"
	domain "content-hub/internal/domain/usecase"
	"content-hub/internal/logger"
)

type newsUsecase struct {
	repo repository.NewsRepository
}

func NewNewsUsecase(
	r repository.NewsRepository,
) domain.NewsUsecase {
	return &newsUsecase{r}
}

func (u *newsUsecase) Create(n *entity.News) error {

	payloadBytes, _ := json.Marshal(n)

	event := &entity.OutboxEvent{
		AggregateType: "news",
		EventType:     "news_created",
		Payload:       string(payloadBytes),
		Status:        "pending",
	}

	err := u.repo.CreateWithOutbox(n, event)

	if err != nil {

		logger.Log.Error().
			Err(err).
			Str("title", n.Title).
			Msg("failed to create news")

		return err
	}

	logger.Log.Info().
		Uint64("news_id", n.ID).
		Str("title", n.Title).
		Msg("news created")

	return nil
}

func (u *newsUsecase) GetAll(
	page,
	limit int,
	categoryID *uint64,
) ([]entity.News, int64, error) {

	news, err := u.repo.FindAll(
		page,
		limit,
		categoryID,
	)

	if err != nil {

		logger.Log.Error().
			Err(err).
			Int("page", page).
			Int("limit", limit).
			Msg("failed to fetch news")

		return nil, 0, err
	}

	total, err := u.repo.Count(categoryID)

	if err != nil {

		logger.Log.Error().
			Err(err).
			Msg("failed to count news")

		return nil, 0, err
	}

	logger.Log.Info().
		Int("page", page).
		Int("limit", limit).
		Int64("total", total).
		Msg("news fetched")

	return news, total, nil
}

func (u *newsUsecase) GetByID(id uint64) (*entity.News, error) {

	news, err := u.repo.FindByID(id)

	if err != nil {

		logger.Log.Error().
			Err(err).
			Uint64("news_id", id).
			Msg("failed to get news by id")

		return nil, err
	}

	logger.Log.Info().
		Uint64("news_id", id).
		Msg("news fetched by id")

	return news, nil
}

func (u *newsUsecase) Update(n *entity.News) error {

	payloadBytes, _ := json.Marshal(n)

	event := &entity.OutboxEvent{
		AggregateType: "news",
		AggregateID:   n.ID,
		EventType:     "news_updated",
		Payload:       string(payloadBytes),
		Status:        "pending",
	}

	err := u.repo.UpdateWithOutbox(n, event)

	if err != nil {

		logger.Log.Error().
			Err(err).
			Uint64("news_id", n.ID).
			Msg("failed to update news")

		return err
	}

	logger.Log.Info().
		Uint64("news_id", n.ID).
		Str("title", n.Title).
		Msg("news updated")

	return nil
}

func (u *newsUsecase) Delete(id uint64) error {

	payload := fmt.Sprintf(`{"id": %d}`, id)

	event := &entity.OutboxEvent{
		AggregateType: "news",
		AggregateID:   id,
		EventType:     "news_deleted",
		Payload:       payload,
		Status:        "pending",
	}

	err := u.repo.DeleteWithOutbox(id, event)

	if err != nil {

		logger.Log.Error().
			Err(err).
			Uint64("news_id", id).
			Msg("failed to delete news")

		return err
	}

	logger.Log.Info().
		Uint64("news_id", id).
		Msg("news deleted")

	return nil
}
