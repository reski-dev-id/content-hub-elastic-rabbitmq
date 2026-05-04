package repository

import "content-hub/internal/domain/entity"

type ProductRepository interface {
	Create(product *entity.Product) error
	FindAll() ([]entity.Product, error)
}
