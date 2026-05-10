package repository

import "content-hub/internal/domain/entity"

type ProductRepository interface {
	Create(product *entity.Product) error
	CreateWithOutbox(product *entity.Product, event *entity.OutboxEvent) error
	FindAll(page, limit int, categoryID *uint64) ([]entity.Product, error)
	FindByID(id uint64) (*entity.Product, error)
	Update(product *entity.Product) error
	Delete(id uint64) error
	DeleteWithOutbox(id uint64, event *entity.OutboxEvent) error
	UpdateWithOutbox(product *entity.Product, event *entity.OutboxEvent) error
}
