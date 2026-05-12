package usecase

import (
	"encoding/json"
	"fmt"

	"content-hub/internal/domain/entity"
	"content-hub/internal/domain/repository"
	domain "content-hub/internal/domain/usecase"
	"content-hub/internal/logger"
)

type productUsecase struct {
	repo repository.ProductRepository
}

func NewProductUsecase(
	r repository.ProductRepository,
) domain.ProductUsecase {
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

	err := u.repo.CreateWithOutbox(p, event)

	if err != nil {

		logger.Error(err).
			Str("title", p.Title).
			Msg("failed to create product")

		return err
	}

	logger.Info().
		Uint64("product_id", p.ID).
		Str("title", p.Title).
		Msg("product created")

	return nil
}

func (u *productUsecase) GetAll(
	page,
	limit int,
	categoryID *uint64,
) ([]entity.Product, int64, error) {

	products, err := u.repo.FindAll(
		page,
		limit,
		categoryID,
	)

	if err != nil {

		logger.Error(err).
			Int("page", page).
			Int("limit", limit).
			Msg("failed to fetch products")

		return nil, 0, err
	}

	total, err := u.repo.Count(categoryID)

	if err != nil {

		logger.Error(err).
			Msg("failed to count products")

		return nil, 0, err
	}

	logger.Info().
		Int("page", page).
		Int("limit", limit).
		Int64("total", total).
		Msg("products fetched")

	return products, total, nil
}

func (u *productUsecase) GetByID(id uint64) (*entity.Product, error) {

	product, err := u.repo.FindByID(id)

	if err != nil {

		logger.Error(err).
			Uint64("product_id", id).
			Msg("failed to get product by id")

		return nil, err
	}

	logger.Info().
		Uint64("product_id", id).
		Msg("product fetched by id")

	return product, nil
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

	err := u.repo.UpdateWithOutbox(p, event)

	if err != nil {

		logger.Error(err).
			Uint64("product_id", p.ID).
			Msg("failed to update product")

		return err
	}

	logger.Info().
		Uint64("product_id", p.ID).
		Str("title", p.Title).
		Msg("product updated")

	return nil
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

	err := u.repo.DeleteWithOutbox(id, event)

	if err != nil {

		logger.Error(err).
			Uint64("product_id", id).
			Msg("failed to delete product")

		return err
	}

	logger.Info().
		Uint64("product_id", id).
		Msg("product deleted")

	return nil
}
