package usecase

import (
	"context"

	domain "content-hub/internal/domain/usecase"
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
) (interface{}, error) {

	index := "products"

	if searchType == "news" {
		index = "news"
	}

	result, err := u.repo.Search(
		context.Background(),
		index,
		q,
		categoryID,
		page,
		limit,
	)

	if err != nil {
		return nil, err
	}

	hits := result["hits"].(map[string]interface{})
	items := hits["hits"].([]interface{})

	response := []map[string]interface{}{}

	for _, item := range items {
		hit := item.(map[string]interface{})
		source := hit["_source"].(map[string]interface{})
		source["score"] = hit["_score"]

		response = append(response, source)
	}

	return map[string]interface{}{
		"page":  page,
		"limit": limit,
		"total": hits["total"],
		"items": response,
	}, nil
}
