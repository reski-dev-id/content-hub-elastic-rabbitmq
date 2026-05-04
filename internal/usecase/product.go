package usecase

import (
	"encoding/json"

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

func (u *productUsecase) GetAll(page, limit int, categoryID *uint64) ([]entity.Product, error) {
	return u.repo.FindAll(page, limit, categoryID)
}

func (u *productUsecase) GetByID(id uint64) (*entity.Product, error) {
	return u.repo.FindByID(id)
}

func (u *productUsecase) Update(p *entity.Product) error {
	return u.repo.Update(p)
}
