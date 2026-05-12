package mysql

import (
	"time"

	"content-hub/internal/domain/entity"
	"content-hub/internal/domain/repository"
	"content-hub/internal/logger"

	"github.com/jmoiron/sqlx"
)

type outboxRepo struct {
	db *sqlx.DB
}

func NewOutboxRepository(
	db *sqlx.DB,
) repository.OutboxRepository {
	return &outboxRepo{
		db: db,
	}
}

func (r *outboxRepo) FindPending(
	limit int,
) ([]entity.OutboxEvent, error) {

	start := time.Now()

	var events []entity.OutboxEvent

	query := `
	SELECT * FROM outbox_events
	WHERE status = 'pending'
	ORDER BY id ASC
	LIMIT ?`

	err := r.db.Select(
		&events,
		query,
		limit,
	)

	if err != nil {

		logger.Error(err).
			Str("service", "mysql").
			Str("event", "outbox_pending_fetch_failed").
			Int("limit", limit).
			Int64(
				"duration_ms",
				time.Since(start).Milliseconds(),
			).
			Msg("failed fetch pending outbox events")

		return nil, err
	}

	logger.Info().
		Str("service", "mysql").
		Str("event", "outbox_pending_fetched").
		Int("count", len(events)).
		Int("limit", limit).
		Dur(
			"duration",
			time.Since(start),
		).
		Int64(
			"duration_ms",
			time.Since(start).Milliseconds(),
		).
		Msg("pending outbox events fetched")

	return events, nil
}

func (r *outboxRepo) MarkAsSent(
	id uint64,
) error {

	start := time.Now()

	_, err := r.db.Exec(`
	UPDATE outbox_events
	SET status = 'sent', sent_at = NOW()
	WHERE id = ?`,
		id,
	)

	if err != nil {

		logger.Error(err).
			Str("service", "mysql").
			Str("event", "outbox_mark_sent_failed").
			Uint64("event_id", id).
			Int64(
				"duration_ms",
				time.Since(start).Milliseconds(),
			).
			Msg("failed mark outbox event as sent")

		return err
	}

	logger.Info().
		Str("service", "mysql").
		Str("event", "outbox_marked_sent").
		Uint64("event_id", id).
		Dur(
			"duration",
			time.Since(start),
		).
		Int64(
			"duration_ms",
			time.Since(start).Milliseconds(),
		).
		Msg("outbox event marked as sent")

	return nil
}
