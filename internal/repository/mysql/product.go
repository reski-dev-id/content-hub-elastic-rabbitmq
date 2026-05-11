package mysql

import (
	"encoding/json"

	"content-hub/internal/domain/entity"
	"content-hub/internal/domain/repository"
	"content-hub/internal/logger"

	"github.com/jmoiron/sqlx"
)

type productRepo struct {
	db *sqlx.DB
}

func NewProductRepository(
	db *sqlx.DB,
) repository.ProductRepository {
	return &productRepo{db}
}

func (r *productRepo) Create(
	p *entity.Product,
) error {

	query := `
	INSERT INTO products (
		category_id,
		title,
		slug,
		description,
		price,
		status
	)
	VALUES (?, ?, ?, ?, ?, ?)`

	_, err := r.db.Exec(
		query,
		p.CategoryID,
		p.Title,
		p.Slug,
		p.Description,
		p.Price,
		p.Status,
	)

	if err != nil {

		logger.Log.Error().
			Err(err).
			Str("title", p.Title).
			Msg("failed create product")

		return err
	}

	return nil
}

func (r *productRepo) CreateWithOutbox(
	p *entity.Product,
	e *entity.OutboxEvent,
) error {

	tx, err := r.db.Beginx()

	if err != nil {

		logger.Log.Error().
			Err(err).
			Msg("failed begin transaction create product")

		return err
	}

	res, err := tx.Exec(`
		INSERT INTO products (
			category_id,
			title,
			slug,
			description,
			price,
			status
		)
		VALUES (?, ?, ?, ?, ?, ?)`,
		p.CategoryID,
		p.Title,
		p.Slug,
		p.Description,
		p.Price,
		p.Status,
	)

	if err != nil {

		tx.Rollback()

		logger.Log.Error().
			Err(err).
			Str("title", p.Title).
			Msg("failed insert product")

		return err
	}

	productID, err := res.LastInsertId()

	if err != nil {

		tx.Rollback()

		logger.Log.Error().
			Err(err).
			Msg("failed get product last insert id")

		return err
	}

	p.ID = uint64(productID)

	payloadBytes, _ := json.Marshal(p)

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
		productID,
		e.EventType,
		string(payloadBytes),
		"pending",
	)

	if err != nil {

		tx.Rollback()

		logger.Log.Error().
			Err(err).
			Uint64("product_id", uint64(productID)).
			Msg("failed insert outbox product event")

		return err
	}

	err = tx.Commit()

	if err != nil {

		logger.Log.Error().
			Err(err).
			Uint64("product_id", uint64(productID)).
			Msg("failed commit transaction create product")

		return err
	}

	return nil
}

func (r *productRepo) FindAll(
	page,
	limit int,
	categoryID *uint64,
) ([]entity.Product, error) {

	var products []entity.Product

	offset := (page - 1) * limit

	query := "SELECT * FROM products WHERE 1=1"
	args := []interface{}{}

	if categoryID != nil {
		query += " AND category_id = ?"
		args = append(args, *categoryID)
	}

	query += " ORDER BY id DESC LIMIT ? OFFSET ?"
	args = append(args, limit, offset)

	err := r.db.Select(&products, query, args...)

	if err != nil {

		logger.Log.Error().
			Err(err).
			Int("page", page).
			Int("limit", limit).
			Msg("failed fetch products")

		return nil, err
	}

	return products, nil
}

func (r *productRepo) Count(
	categoryID *uint64,
) (int64, error) {

	var total int64

	query := "SELECT COUNT(*) FROM products WHERE 1=1"
	args := []interface{}{}

	if categoryID != nil {
		query += " AND category_id = ?"
		args = append(args, *categoryID)
	}

	err := r.db.Get(&total, query, args...)

	if err != nil {

		logger.Log.Error().
			Err(err).
			Msg("failed count products")

		return 0, err
	}

	return total, nil
}

func (r *productRepo) FindByID(
	id uint64,
) (*entity.Product, error) {

	var product entity.Product

	err := r.db.Get(
		&product,
		"SELECT * FROM products WHERE id = ?",
		id,
	)

	if err != nil {

		logger.Log.Error().
			Err(err).
			Uint64("product_id", id).
			Msg("failed find product by id")

		return nil, err
	}

	return &product, nil
}

func (r *productRepo) Update(
	p *entity.Product,
) error {

	query := `
	UPDATE products
	SET category_id=?,
		title=?,
		slug=?,
		description=?,
		price=?,
		status=?
	WHERE id=?`

	_, err := r.db.Exec(
		query,
		p.CategoryID,
		p.Title,
		p.Slug,
		p.Description,
		p.Price,
		p.Status,
		p.ID,
	)

	if err != nil {

		logger.Log.Error().
			Err(err).
			Uint64("product_id", p.ID).
			Msg("failed update product")

		return err
	}

	return nil
}

func (r *productRepo) Delete(
	id uint64,
) error {

	_, err := r.db.Exec(
		"DELETE FROM products WHERE id = ?",
		id,
	)

	if err != nil {

		logger.Log.Error().
			Err(err).
			Uint64("product_id", id).
			Msg("failed delete product")

		return err
	}

	return nil
}

func (r *productRepo) DeleteWithOutbox(
	id uint64,
	e *entity.OutboxEvent,
) error {

	tx, err := r.db.Beginx()

	if err != nil {

		logger.Log.Error().
			Err(err).
			Msg("failed begin transaction delete product")

		return err
	}

	_, err = tx.Exec(
		"DELETE FROM products WHERE id = ?",
		id,
	)

	if err != nil {

		tx.Rollback()

		logger.Log.Error().
			Err(err).
			Uint64("product_id", id).
			Msg("failed delete product transaction")

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

		logger.Log.Error().
			Err(err).
			Uint64("product_id", id).
			Msg("failed insert outbox delete product")

		return err
	}

	err = tx.Commit()

	if err != nil {

		logger.Log.Error().
			Err(err).
			Uint64("product_id", id).
			Msg("failed commit delete product")

		return err
	}

	return nil
}

func (r *productRepo) UpdateWithOutbox(
	p *entity.Product,
	e *entity.OutboxEvent,
) error {

	tx, err := r.db.Beginx()

	if err != nil {

		logger.Log.Error().
			Err(err).
			Msg("failed begin transaction update product")

		return err
	}

	_, err = tx.Exec(`
		UPDATE products
		SET category_id=?,
			title=?,
			slug=?,
			description=?,
			price=?,
			status=?
		WHERE id=?`,
		p.CategoryID,
		p.Title,
		p.Slug,
		p.Description,
		p.Price,
		p.Status,
		p.ID,
	)

	if err != nil {

		tx.Rollback()

		logger.Log.Error().
			Err(err).
			Uint64("product_id", p.ID).
			Msg("failed update product transaction")

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
		p.ID,
		e.EventType,
		e.Payload,
		"pending",
	)

	if err != nil {

		tx.Rollback()

		logger.Log.Error().
			Err(err).
			Uint64("product_id", p.ID).
			Msg("failed insert outbox update product")

		return err
	}

	err = tx.Commit()

	if err != nil {

		logger.Log.Error().
			Err(err).
			Uint64("product_id", p.ID).
			Msg("failed commit update product")

		return err
	}

	return nil
}
