package repository

import "content-hub/internal/domain/entity"

type NewsRepository interface {
	Create(news *entity.News) error
	CreateWithOutbox(news *entity.News, event *entity.OutboxEvent) error
	FindAll(page, limit int, categoryID *uint64) ([]entity.News, error)
	Count(categoryID *uint64) (int64, error)
	FindByID(id uint64) (*entity.News, error)
	Update(news *entity.News) error
	UpdateWithOutbox(news *entity.News, event *entity.OutboxEvent) error
	Delete(id uint64) error
	DeleteWithOutbox(id uint64, event *entity.OutboxEvent) error
}
