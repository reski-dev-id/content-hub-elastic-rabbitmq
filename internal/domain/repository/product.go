package repository

import "content-hub/internal/domain/entity"

type ProductRepository interface {
	Create(product *entity.Product) error
	FindAll(page, limit int, categoryID *uint64) ([]entity.Product, error)
	FindByID(id uint64) (*entity.Product, error)
	Update(product *entity.Product) error
}
