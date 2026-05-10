package usecase

import "content-hub/internal/domain/entity"

type NewsUsecase interface {
	Create(news *entity.News) error

	GetAll(
		page,
		limit int,
		categoryID *uint64,
	) ([]entity.News, error)

	GetByID(id uint64) (*entity.News, error)

	Update(news *entity.News) error

	Delete(id uint64) error
}
