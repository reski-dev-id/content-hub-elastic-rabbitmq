package usecase

import "content-hub/internal/domain/entity"

type NewsSearchResponse struct {
	Items           []entity.News `json:"items"`
	Recommendations []entity.News `json:"recommendations"`
}

type NewsUsecase interface {
	Create(news *entity.News) error
	GetAll(page, limit int, categoryID *uint64) ([]entity.News, int64, error)
	GetByID(id uint64) (*entity.News, error)
	Search(
		q string,
		categoryID *uint64,
		page int,
		limit int,
	) (*NewsSearchResponse, int64, error)
	Recommend(id uint64, limit int) ([]entity.News, error)
	Update(news *entity.News) error
	Delete(id uint64) error
}
