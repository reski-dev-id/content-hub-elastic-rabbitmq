package usecase

import "content-hub/internal/domain/entity"

type ProductUsecase interface {
	Create(product *entity.Product) error
	GetAll() ([]entity.Product, error)
}
