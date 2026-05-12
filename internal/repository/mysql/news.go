package mysql

import (
	"encoding/json"
	"time"

	"content-hub/internal/domain/entity"
	"content-hub/internal/domain/repository"
	"content-hub/internal/logger"

	"github.com/jmoiron/sqlx"
)

type newsRepo struct {
	db *sqlx.DB
}

func NewNewsRepository(
	db *sqlx.DB,
) repository.NewsRepository {
	return &newsRepo{
		db: db,
	}
}

func (r *newsRepo) Create(
	n *entity.News,
) error {

	start := time.Now()

	query := `
	INSERT INTO news (
		category_id,
		title,
		slug,
		content,
		author,
		status,
		published_at
	)
	VALUES (?, ?, ?, ?, ?, ?, ?)`

	_, err := r.db.Exec(
		query,
		n.CategoryID,
		n.Title,
		n.Slug,
		n.Content,
		n.Author,
		n.Status,
		n.PublishedAt,
	)

	if err != nil {

		logger.Error(err).
			Str("service", "mysql").
			Str("event", "news_create_failed").
			Str("title", n.Title).
			Int64(
				"duration_ms",
				time.Since(start).Milliseconds(),
			).
			Msg("failed create news")

		return err
	}

	logger.Info().
		Str("service", "mysql").
		Str("event", "news_created").
		Str("title", n.Title).
		Dur(
			"duration",
			time.Since(start),
		).
		Int64(
			"duration_ms",
			time.Since(start).Milliseconds(),
		).
		Msg("news created")

	return nil
}

func (r *newsRepo) CreateWithOutbox(
	n *entity.News,
	e *entity.OutboxEvent,
) error {

	start := time.Now()

	tx, err := r.db.Beginx()

	if err != nil {

		logger.Error(err).
			Str("service", "mysql").
			Str("event", "news_create_transaction_begin_failed").
			Int64(
				"duration_ms",
				time.Since(start).Milliseconds(),
			).
			Msg("failed begin transaction create news")

		return err
	}

	res, err := tx.Exec(`
		INSERT INTO news (
			category_id,
			title,
			slug,
			content,
			author,
			status,
			published_at
		)
		VALUES (?, ?, ?, ?, ?, ?, ?)`,
		n.CategoryID,
		n.Title,
		n.Slug,
		n.Content,
		n.Author,
		n.Status,
		n.PublishedAt,
	)

	if err != nil {

		tx.Rollback()

		logger.Error(err).
			Str("service", "mysql").
			Str("event", "news_insert_failed").
			Str("title", n.Title).
			Int64(
				"duration_ms",
				time.Since(start).Milliseconds(),
			).
			Msg("failed insert news")

		return err
	}

	newsID, err := res.LastInsertId()

	if err != nil {

		tx.Rollback()

		logger.Error(err).
			Str("service", "mysql").
			Str("event", "news_last_insert_id_failed").
			Int64(
				"duration_ms",
				time.Since(start).Milliseconds(),
			).
			Msg("failed get news last insert id")

		return err
	}

	n.ID = uint64(newsID)

	payloadBytes, _ := json.Marshal(
		n,
	)

	_, err = tx.Exec(`
		INSERT INTO outbox_events (
			aggregate_type,
			aggregate_id,
			event_type,
			payload,
			status
		)
		VALUES (?, ?, ?, ?, ?)`,
		e.AggregateType,
		newsID,
		e.EventType,
		string(payloadBytes),
		"pending",
	)

	if err != nil {

		tx.Rollback()

		logger.Error(err).
			Str("service", "mysql").
			Str("event", "news_outbox_insert_failed").
			Uint64("news_id", uint64(newsID)).
			Int64(
				"duration_ms",
				time.Since(start).Milliseconds(),
			).
			Msg("failed insert outbox event")

		return err
	}

	err = tx.Commit()

	if err != nil {

		logger.Error(err).
			Str("service", "mysql").
			Str("event", "news_create_commit_failed").
			Uint64("news_id", uint64(newsID)).
			Int64(
				"duration_ms",
				time.Since(start).Milliseconds(),
			).
			Msg("failed commit transaction create news")

		return err
	}

	logger.Info().
		Str("service", "mysql").
		Str("event", "news_created_with_outbox").
		Uint64("news_id", uint64(newsID)).
		Dur(
			"duration",
			time.Since(start),
		).
		Int64(
			"duration_ms",
			time.Since(start).Milliseconds(),
		).
		Msg("news created with outbox event")

	return nil
}

func (r *newsRepo) FindAll(
	page,
	limit int,
	categoryID *uint64,
) ([]entity.News, error) {

	start := time.Now()

	var news []entity.News

	offset := (page - 1) * limit

	query := "SELECT * FROM news WHERE 1=1"

	args := []interface{}{}

	if categoryID != nil {

		query += " AND category_id = ?"

		args = append(
			args,
			*categoryID,
		)
	}

	query += " ORDER BY id DESC LIMIT ? OFFSET ?"

	args = append(
		args,
		limit,
		offset,
	)

	err := r.db.Select(
		&news,
		query,
		args...,
	)

	if err != nil {

		logger.Error(err).
			Str("service", "mysql").
			Str("event", "news_fetch_failed").
			Int("page", page).
			Int("limit", limit).
			Int64(
				"duration_ms",
				time.Since(start).Milliseconds(),
			).
			Msg("failed fetch news")

		return nil, err
	}

	logger.Info().
		Str("service", "mysql").
		Str("event", "news_fetched").
		Int("count", len(news)).
		Int("page", page).
		Int("limit", limit).
		Dur(
			"duration",
			time.Since(start),
		).
		Int64(
			"duration_ms",
			time.Since(start).Milliseconds(),
		).
		Msg("news fetched successfully")

	return news, nil
}

func (r *newsRepo) Count(
	categoryID *uint64,
) (int64, error) {

	start := time.Now()

	var total int64

	query := "SELECT COUNT(*) FROM news WHERE 1=1"

	args := []interface{}{}

	if categoryID != nil {

		query += " AND category_id = ?"

		args = append(
			args,
			*categoryID,
		)
	}

	err := r.db.Get(
		&total,
		query,
		args...,
	)

	if err != nil {

		logger.Error(err).
			Str("service", "mysql").
			Str("event", "news_count_failed").
			Int64(
				"duration_ms",
				time.Since(start).Milliseconds(),
			).
			Msg("failed count news")

		return 0, err
	}

	logger.Info().
		Str("service", "mysql").
		Str("event", "news_counted").
		Int64("total", total).
		Dur(
			"duration",
			time.Since(start),
		).
		Int64(
			"duration_ms",
			time.Since(start).Milliseconds(),
		).
		Msg("news counted successfully")

	return total, nil
}

func (r *newsRepo) FindByID(
	id uint64,
) (*entity.News, error) {

	start := time.Now()

	var news entity.News

	err := r.db.Get(
		&news,
		"SELECT * FROM news WHERE id = ?",
		id,
	)

	if err != nil {

		logger.Error(err).
			Str("service", "mysql").
			Str("event", "news_find_by_id_failed").
			Uint64("news_id", id).
			Int64(
				"duration_ms",
				time.Since(start).Milliseconds(),
			).
			Msg("failed find news by id")

		return nil, err
	}

	logger.Info().
		Str("service", "mysql").
		Str("event", "news_found").
		Uint64("news_id", id).
		Dur(
			"duration",
			time.Since(start),
		).
		Int64(
			"duration_ms",
			time.Since(start).Milliseconds(),
		).
		Msg("news found successfully")

	return &news, nil
}

func (r *newsRepo) Update(
	n *entity.News,
) error {

	start := time.Now()

	query := `
	UPDATE news
	SET category_id=?,
		title=?,
		slug=?,
		content=?,
		author=?,
		status=?,
		published_at=?
	WHERE id=?`

	_, err := r.db.Exec(
		query,
		n.CategoryID,
		n.Title,
		n.Slug,
		n.Content,
		n.Author,
		n.Status,
		n.PublishedAt,
		n.ID,
	)

	if err != nil {

		logger.Error(err).
			Str("service", "mysql").
			Str("event", "news_update_failed").
			Uint64("news_id", n.ID).
			Int64(
				"duration_ms",
				time.Since(start).Milliseconds(),
			).
			Msg("failed update news")

		return err
	}

	logger.Info().
		Str("service", "mysql").
		Str("event", "news_updated").
		Uint64("news_id", n.ID).
		Dur(
			"duration",
			time.Since(start),
		).
		Int64(
			"duration_ms",
			time.Since(start).Milliseconds(),
		).
		Msg("news updated successfully")

	return nil
}

func (r *newsRepo) UpdateWithOutbox(
	n *entity.News,
	e *entity.OutboxEvent,
) error {

	start := time.Now()

	tx, err := r.db.Beginx()

	if err != nil {

		logger.Error(err).
			Str("service", "mysql").
			Str("event", "news_update_transaction_begin_failed").
			Int64(
				"duration_ms",
				time.Since(start).Milliseconds(),
			).
			Msg("failed begin transaction update news")

		return err
	}

	_, err = tx.Exec(`
		UPDATE news
		SET category_id=?,
			title=?,
			slug=?,
			content=?,
			author=?,
			status=?,
			published_at=?
		WHERE id=?`,
		n.CategoryID,
		n.Title,
		n.Slug,
		n.Content,
		n.Author,
		n.Status,
		n.PublishedAt,
		n.ID,
	)

	if err != nil {

		tx.Rollback()

		logger.Error(err).
			Str("service", "mysql").
			Str("event", "news_update_transaction_failed").
			Uint64("news_id", n.ID).
			Int64(
				"duration_ms",
				time.Since(start).Milliseconds(),
			).
			Msg("failed update news transaction")

		return err
	}

	_, err = tx.Exec(`
		INSERT INTO outbox_events (
			aggregate_type,
			aggregate_id,
			event_type,
			payload,
			status
		)
		VALUES (?, ?, ?, ?, ?)`,
		e.AggregateType,
		n.ID,
		e.EventType,
		e.Payload,
		"pending",
	)

	if err != nil {

		tx.Rollback()

		logger.Error(err).
			Str("service", "mysql").
			Str("event", "news_update_outbox_insert_failed").
			Uint64("news_id", n.ID).
			Int64(
				"duration_ms",
				time.Since(start).Milliseconds(),
			).
			Msg("failed insert outbox update news")

		return err
	}

	err = tx.Commit()

	if err != nil {

		logger.Error(err).
			Str("service", "mysql").
			Str("event", "news_update_commit_failed").
			Uint64("news_id", n.ID).
			Int64(
				"duration_ms",
				time.Since(start).Milliseconds(),
			).
			Msg("failed commit update news")

		return err
	}

	logger.Info().
		Str("service", "mysql").
		Str("event", "news_updated_with_outbox").
		Uint64("news_id", n.ID).
		Dur(
			"duration",
			time.Since(start),
		).
		Int64(
			"duration_ms",
			time.Since(start).Milliseconds(),
		).
		Msg("news updated with outbox")

	return nil
}

func (r *newsRepo) Delete(
	id uint64,
) error {

	start := time.Now()

	_, err := r.db.Exec(
		"DELETE FROM news WHERE id = ?",
		id,
	)

	if err != nil {

		logger.Error(err).
			Str("service", "mysql").
			Str("event", "news_delete_failed").
			Uint64("news_id", id).
			Int64(
				"duration_ms",
				time.Since(start).Milliseconds(),
			).
			Msg("failed delete news")

		return err
	}

	logger.Info().
		Str("service", "mysql").
		Str("event", "news_deleted").
		Uint64("news_id", id).
		Dur(
			"duration",
			time.Since(start),
		).
		Int64(
			"duration_ms",
			time.Since(start).Milliseconds(),
		).
		Msg("news deleted successfully")

	return nil
}

func (r *newsRepo) DeleteWithOutbox(
	id uint64,
	e *entity.OutboxEvent,
) error {

	start := time.Now()

	tx, err := r.db.Beginx()

	if err != nil {

		logger.Error(err).
			Str("service", "mysql").
			Str("event", "news_delete_transaction_begin_failed").
			Int64(
				"duration_ms",
				time.Since(start).Milliseconds(),
			).
			Msg("failed begin transaction delete news")

		return err
	}

	_, err = tx.Exec(
		"DELETE FROM news WHERE id = ?",
		id,
	)

	if err != nil {

		tx.Rollback()

		logger.Error(err).
			Str("service", "mysql").
			Str("event", "news_delete_transaction_failed").
			Uint64("news_id", id).
			Int64(
				"duration_ms",
				time.Since(start).Milliseconds(),
			).
			Msg("failed delete news transaction")

		return err
	}

	_, err = tx.Exec(`
		INSERT INTO outbox_events (
			aggregate_type,
			aggregate_id,
			event_type,
			payload,
			status
		)
		VALUES (?, ?, ?, ?, ?)`,
		e.AggregateType,
		id,
		e.EventType,
		e.Payload,
		"pending",
	)

	if err != nil {

		tx.Rollback()

		logger.Error(err).
			Str("service", "mysql").
			Str("event", "news_delete_outbox_insert_failed").
			Uint64("news_id", id).
			Int64(
				"duration_ms",
				time.Since(start).Milliseconds(),
			).
			Msg("failed insert outbox delete news")

		return err
	}

	err = tx.Commit()

	if err != nil {

		logger.Error(err).
			Str("service", "mysql").
			Str("event", "news_delete_commit_failed").
			Uint64("news_id", id).
			Int64(
				"duration_ms",
				time.Since(start).Milliseconds(),
			).
			Msg("failed commit delete news")

		return err
	}

	logger.Info().
		Str("service", "mysql").
		Str("event", "news_deleted_with_outbox").
		Uint64("news_id", id).
		Dur(
			"duration",
			time.Since(start),
		).
		Int64(
			"duration_ms",
			time.Since(start).Milliseconds(),
		).
		Msg("news deleted with outbox")

	return nil
}
