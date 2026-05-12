package usecase

import (
	"encoding/json"
	"fmt"
	"time"

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
	return &productUsecase{
		repo: r,
	}
}

func (u *productUsecase) Create(
	p *entity.Product,
) error {

	start := time.Now()

	payloadBytes, _ := json.Marshal(
		p,
	)

	event := &entity.OutboxEvent{
		AggregateType: "product",
		EventType:     "product_created",
		Payload:       string(payloadBytes),
		Status:        "pending",
	}

	err := u.repo.CreateWithOutbox(
		p,
		event,
	)

	if err != nil {

		logger.Error(err).
			Str("service", "usecase").
			Str("event", "product_create_failed").
			Str("title", p.Title).
			Int64(
				"duration_ms",
				time.Since(start).Milliseconds(),
			).
			Msg("failed to create product")

		return err
	}

	logger.Info().
		Str("service", "usecase").
		Str("event", "product_created").
		Uint64("product_id", p.ID).
		Str("title", p.Title).
		Dur(
			"duration",
			time.Since(start),
		).
		Int64(
			"duration_ms",
			time.Since(start).Milliseconds(),
		).
		Msg("product created")

	return nil
}

func (u *productUsecase) GetAll(
	page,
	limit int,
	categoryID *uint64,
) ([]entity.Product, int64, error) {

	start := time.Now()

	products, err := u.repo.FindAll(
		page,
		limit,
		categoryID,
	)

	if err != nil {

		logger.Error(err).
			Str("service", "usecase").
			Str("event", "products_fetch_failed").
			Int("page", page).
			Int("limit", limit).
			Int64(
				"duration_ms",
				time.Since(start).Milliseconds(),
			).
			Msg("failed to fetch products")

		return nil, 0, err
	}

	total, err := u.repo.Count(
		categoryID,
	)

	if err != nil {

		logger.Error(err).
			Str("service", "usecase").
			Str("event", "products_count_failed").
			Int64(
				"duration_ms",
				time.Since(start).Milliseconds(),
			).
			Msg("failed to count products")

		return nil, 0, err
	}

	logger.Info().
		Str("service", "usecase").
		Str("event", "products_fetched").
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
		Msg("products fetched")

	return products, total, nil
}

func (u *productUsecase) GetByID(
	id uint64,
) (*entity.Product, error) {

	start := time.Now()

	product, err := u.repo.FindByID(
		id,
	)

	if err != nil {

		logger.Error(err).
			Str("service", "usecase").
			Str("event", "product_find_by_id_failed").
			Uint64("product_id", id).
			Int64(
				"duration_ms",
				time.Since(start).Milliseconds(),
			).
			Msg("failed to get product by id")

		return nil, err
	}

	logger.Info().
		Str("service", "usecase").
		Str("event", "product_found").
		Uint64("product_id", id).
		Dur(
			"duration",
			time.Since(start),
		).
		Int64(
			"duration_ms",
			time.Since(start).Milliseconds(),
		).
		Msg("product fetched by id")

	return product, nil
}

func (u *productUsecase) Update(
	p *entity.Product,
) error {

	start := time.Now()

	payloadBytes, _ := json.Marshal(
		p,
	)

	event := &entity.OutboxEvent{
		AggregateType: "product",
		AggregateID:   p.ID,
		EventType:     "product_updated",
		Payload:       string(payloadBytes),
		Status:        "pending",
	}

	err := u.repo.UpdateWithOutbox(
		p,
		event,
	)

	if err != nil {

		logger.Error(err).
			Str("service", "usecase").
			Str("event", "product_update_failed").
			Uint64("product_id", p.ID).
			Int64(
				"duration_ms",
				time.Since(start).Milliseconds(),
			).
			Msg("failed to update product")

		return err
	}

	logger.Info().
		Str("service", "usecase").
		Str("event", "product_updated").
		Uint64("product_id", p.ID).
		Str("title", p.Title).
		Dur(
			"duration",
			time.Since(start),
		).
		Int64(
			"duration_ms",
			time.Since(start).Milliseconds(),
		).
		Msg("product updated")

	return nil
}

func (u *productUsecase) Delete(
	id uint64,
) error {

	start := time.Now()

	payload := fmt.Sprintf(
		`{"id": %d}`,
		id,
	)

	event := &entity.OutboxEvent{
		AggregateType: "product",
		AggregateID:   id,
		EventType:     "product_deleted",
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
			Str("event", "product_delete_failed").
			Uint64("product_id", id).
			Int64(
				"duration_ms",
				time.Since(start).Milliseconds(),
			).
			Msg("failed to delete product")

		return err
	}

	logger.Info().
		Str("service", "usecase").
		Str("event", "product_deleted").
		Uint64("product_id", id).
		Dur(
			"duration",
			time.Since(start),
		).
		Int64(
			"duration_ms",
			time.Since(start).Milliseconds(),
		).
		Msg("product deleted")

	return nil
}
