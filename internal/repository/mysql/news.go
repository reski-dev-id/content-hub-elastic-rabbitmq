package mysql

import (
	"content-hub/internal/domain/entity"
	"content-hub/internal/domain/repository"
	"encoding/json"

	"github.com/jmoiron/sqlx"
)

type newsRepo struct {
	db *sqlx.DB
}

func NewNewsRepository(db *sqlx.DB) repository.NewsRepository {
	return &newsRepo{db}
}

func (r *newsRepo) Create(n *entity.News) error {
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

	return err
}

func (r *newsRepo) CreateWithOutbox(
	n *entity.News,
	e *entity.OutboxEvent,
) error {
	tx, err := r.db.Beginx()

	if err != nil {
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
		return err
	}

	newsID, err := res.LastInsertId()

	if err != nil {
		tx.Rollback()
		return err
	}

	n.ID = uint64(newsID)

	payloadBytes, _ := json.Marshal(n)

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
		return err
	}

	return tx.Commit()
}

func (r *newsRepo) FindAll(
	page,
	limit int,
	categoryID *uint64,
) ([]entity.News, error) {
	var news []entity.News

	offset := (page - 1) * limit

	query := "SELECT * FROM news WHERE 1=1"
	args := []interface{}{}

	if categoryID != nil {
		query += " AND category_id = ?"
		args = append(args, *categoryID)
	}

	query += " ORDER BY id DESC LIMIT ? OFFSET ?"
	args = append(args, limit, offset)

	err := r.db.Select(&news, query, args...)

	return news, err
}

func (r *newsRepo) Count(categoryID *uint64) (int64, error) {
	var total int64

	query := "SELECT COUNT(*) FROM news WHERE 1=1"
	args := []interface{}{}

	if categoryID != nil {
		query += " AND category_id = ?"
		args = append(args, *categoryID)
	}

	err := r.db.Get(&total, query, args...)

	return total, err
}

func (r *newsRepo) FindByID(id uint64) (*entity.News, error) {
	var news entity.News

	err := r.db.Get(
		&news,
		"SELECT * FROM news WHERE id = ?",
		id,
	)

	if err != nil {
		return nil, err
	}

	return &news, nil
}

func (r *newsRepo) Update(n *entity.News) error {
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

	return err
}

func (r *newsRepo) UpdateWithOutbox(
	n *entity.News,
	e *entity.OutboxEvent,
) error {
	tx, err := r.db.Beginx()

	if err != nil {
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
		return err
	}

	return tx.Commit()
}

func (r *newsRepo) Delete(id uint64) error {
	_, err := r.db.Exec(
		"DELETE FROM news WHERE id = ?",
		id,
	)

	return err
}

func (r *newsRepo) DeleteWithOutbox(
	id uint64,
	e *entity.OutboxEvent,
) error {
	tx, err := r.db.Beginx()

	if err != nil {
		return err
	}

	_, err = tx.Exec(
		"DELETE FROM news WHERE id = ?",
		id,
	)

	if err != nil {
		tx.Rollback()
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
		return err
	}

	return tx.Commit()
}
