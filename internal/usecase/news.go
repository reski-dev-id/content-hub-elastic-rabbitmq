package usecase

import (
	"encoding/json"
	"fmt"

	"content-hub/internal/domain/entity"
	"content-hub/internal/domain/repository"
	domain "content-hub/internal/domain/usecase"
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

	return u.repo.CreateWithOutbox(n, event)
}

func (u *newsUsecase) GetAll(
	page,
	limit int,
	categoryID *uint64,
) ([]entity.News, error) {
	return u.repo.FindAll(page, limit, categoryID)
}

func (u *newsUsecase) GetByID(id uint64) (*entity.News, error) {
	return u.repo.FindByID(id)
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

	return u.repo.UpdateWithOutbox(n, event)
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

	return u.repo.DeleteWithOutbox(id, event)
}
