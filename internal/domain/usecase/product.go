package usecase

import "content-hub/internal/domain/entity"

type ProductUsecase interface {
	Create(product *entity.Product) error
	GetAll(page, limit int, categoryID *uint64) ([]entity.Product, int64, error)
	GetByID(id uint64) (*entity.Product, error)
	Search(
		q string,
		categoryID *uint64,
		page int,
		limit int,
	) (interface{}, int64, error)
	Recommend(id uint64, limit int) ([]entity.Product, error)
	Update(product *entity.Product) error
	Delete(id uint64) error
}
