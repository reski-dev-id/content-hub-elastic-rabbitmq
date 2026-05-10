package usecase

import "content-hub/internal/domain/entity"

type ProductUsecase interface {
	Create(product *entity.Product) error
	GetAll(page, limit int, categoryID *uint64) ([]entity.Product, error)
	GetByID(id uint64) (*entity.Product, error)
	Update(product *entity.Product) error
	Delete(id uint64) error
}
