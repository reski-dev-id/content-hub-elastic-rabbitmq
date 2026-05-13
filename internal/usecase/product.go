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

type productUsecase struct {
	repo   repository.ProductRepository
	esRepo *esRepo.ProductRepository
}

func NewProductUsecase(
	r repository.ProductRepository,
	es *esRepo.ProductRepository,
) domain.ProductUsecase {
	return &productUsecase{
		repo:   r,
		esRepo: es,
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

func (u *productUsecase) Recommend(
	id uint64,
	limit int,
) ([]entity.Product, error) {

	start := time.Now()

	products, err := u.esRepo.Recommend(
		context.Background(),
		strconv.FormatUint(id, 10),
		limit,
	)

	if err != nil {

		logger.Error(err).
			Str("service", "usecase").
			Str("event", "product_recommendation_failed").
			Uint64("product_id", id).
			Int64(
				"duration_ms",
				time.Since(start).Milliseconds(),
			).
			Msg("failed get product recommendations")

		return nil, err
	}

	logger.Info().
		Str("service", "usecase").
		Str("event", "product_recommendations_fetched").
		Uint64("product_id", id).
		Int("count", len(products)).
		Dur(
			"duration",
			time.Since(start),
		).
		Int64(
			"duration_ms",
			time.Since(start).Milliseconds(),
		).
		Msg("product recommendations fetched")

	return products, nil
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

func (u *productUsecase) Search(
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
			Str("event", "product_search_failed").
			Str("query", q).
			Int("page", page).
			Int("limit", limit).
			Int64(
				"duration_ms",
				time.Since(start).Milliseconds(),
			).
			Msg("failed search products")

		return nil, 0, err
	}

	result := &response.ProductSearchResponse{
		Items:           items,
		Recommendations: []entity.Product{},
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
		Str("event", "product_search_completed").
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
		Msg("product search completed")

	return result, total, nil
}
