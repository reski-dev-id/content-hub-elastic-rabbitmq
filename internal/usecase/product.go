package usecase

import (
	"encoding/json"
	"fmt"

	"content-hub/internal/domain/entity"
	"content-hub/internal/domain/repository"
	domain "content-hub/internal/domain/usecase"
)

type productUsecase struct {
	repo repository.ProductRepository
}

func NewProductUsecase(r repository.ProductRepository) domain.ProductUsecase {
	return &productUsecase{r}
}

func (u *productUsecase) Create(p *entity.Product) error {
	payloadBytes, _ := json.Marshal(p)

	event := &entity.OutboxEvent{
		AggregateType: "product",
		EventType:     "product_created",
		Payload:       string(payloadBytes),
		Status:        "pending",
	}

	return u.repo.CreateWithOutbox(p, event)
}

func (u *productUsecase) GetAll(page, limit int, categoryID *uint64) ([]entity.Product, int64, error) {
	products, err := u.repo.FindAll(page, limit, categoryID)

	if err != nil {
		return nil, 0, err
	}

	total, err := u.repo.Count(categoryID)

	if err != nil {
		return nil, 0, err
	}

	return products, total, nil
}

func (u *productUsecase) GetByID(id uint64) (*entity.Product, error) {
	return u.repo.FindByID(id)
}

func (u *productUsecase) Update(p *entity.Product) error {
	payloadBytes, _ := json.Marshal(p)

	event := &entity.OutboxEvent{
		AggregateType: "product",
		AggregateID:   p.ID,
		EventType:     "product_updated",
		Payload:       string(payloadBytes),
		Status:        "pending",
	}

	return u.repo.UpdateWithOutbox(p, event)
}

func (u *productUsecase) Delete(id uint64) error {
	payload := fmt.Sprintf(`{"id": %d}`, id)

	event := &entity.OutboxEvent{
		AggregateType: "product",
		AggregateID:   id,
		EventType:     "product_deleted",
		Payload:       payload,
		Status:        "pending",
	}

	return u.repo.DeleteWithOutbox(id, event)
}
