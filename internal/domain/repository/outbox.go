package repository

import "content-hub/internal/domain/entity"

type OutboxRepository interface {
	FindPending(limit int) ([]entity.OutboxEvent, error)
	MarkAsSent(id uint64) error
}
