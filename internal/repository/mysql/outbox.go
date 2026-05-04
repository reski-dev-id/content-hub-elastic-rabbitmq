package mysql

import (
	"content-hub/internal/domain/entity"
	"content-hub/internal/domain/repository"

	"github.com/jmoiron/sqlx"
)

type outboxRepo struct {
	db *sqlx.DB
}

func NewOutboxRepository(db *sqlx.DB) repository.OutboxRepository {
	return &outboxRepo{db}
}

func (r *outboxRepo) FindPending(limit int) ([]entity.OutboxEvent, error) {
	var events []entity.OutboxEvent

	query := `
	SELECT * FROM outbox_events
	WHERE status = 'pending'
	ORDER BY id ASC
	LIMIT ?`

	err := r.db.Select(&events, query, limit)
	return events, err
}

func (r *outboxRepo) MarkAsSent(id uint64) error {
	_, err := r.db.Exec(`
	UPDATE outbox_events
	SET status = 'sent', sent_at = NOW()
	WHERE id = ?`, id)

	return err
}
