package usecase

import (
	"content-hub/internal/domain/entity"
	"content-hub/internal/domain/repository"
)

type productUsecase struct {
	repo repository.ProductRepository
}

func NewProductUsecase(r repository.ProductRepository) *productUsecase {
	return &productUsecase{r}
}

func (u *productUsecase) Create(p *entity.Product) error {
	return u.repo.Create(p)
}

func (u *productUsecase) GetAll() ([]entity.Product, error) {
	return u.repo.FindAll()
}
