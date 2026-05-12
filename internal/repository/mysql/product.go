package mysql

import (
	"encoding/json"
	"time"

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
	return &productRepo{
		db: db,
	}
}

func (r *productRepo) Create(
	p *entity.Product,
) error {

	start := time.Now()

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

		logger.Error(err).
			Str("service", "mysql").
			Str("event", "product_create_failed").
			Str("title", p.Title).
			Int64(
				"duration_ms",
				time.Since(start).Milliseconds(),
			).
			Msg("failed create product")

		return err
	}

	logger.Info().
		Str("service", "mysql").
		Str("event", "product_created").
		Str("title", p.Title).
		Dur(
			"duration",
			time.Since(start),
		).
		Int64(
			"duration_ms",
			time.Since(start).Milliseconds(),
		).
		Msg("product created")

	return nil
}

func (r *productRepo) CreateWithOutbox(
	p *entity.Product,
	e *entity.OutboxEvent,
) error {

	start := time.Now()

	tx, err := r.db.Beginx()

	if err != nil {

		logger.Error(err).
			Str("service", "mysql").
			Str("event", "product_create_transaction_begin_failed").
			Int64(
				"duration_ms",
				time.Since(start).Milliseconds(),
			).
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

		logger.Error(err).
			Str("service", "mysql").
			Str("event", "product_insert_failed").
			Str("title", p.Title).
			Int64(
				"duration_ms",
				time.Since(start).Milliseconds(),
			).
			Msg("failed insert product")

		return err
	}

	productID, err := res.LastInsertId()

	if err != nil {

		tx.Rollback()

		logger.Error(err).
			Str("service", "mysql").
			Str("event", "product_last_insert_id_failed").
			Int64(
				"duration_ms",
				time.Since(start).Milliseconds(),
			).
			Msg("failed get product last insert id")

		return err
	}

	p.ID = uint64(productID)

	payloadBytes, _ := json.Marshal(
		p,
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
		productID,
		e.EventType,
		string(payloadBytes),
		"pending",
	)

	if err != nil {

		tx.Rollback()

		logger.Error(err).
			Str("service", "mysql").
			Str("event", "product_outbox_insert_failed").
			Uint64("product_id", uint64(productID)).
			Int64(
				"duration_ms",
				time.Since(start).Milliseconds(),
			).
			Msg("failed insert outbox product event")

		return err
	}

	err = tx.Commit()

	if err != nil {

		logger.Error(err).
			Str("service", "mysql").
			Str("event", "product_create_commit_failed").
			Uint64("product_id", uint64(productID)).
			Int64(
				"duration_ms",
				time.Since(start).Milliseconds(),
			).
			Msg("failed commit transaction create product")

		return err
	}

	logger.Info().
		Str("service", "mysql").
		Str("event", "product_created_with_outbox").
		Uint64("product_id", uint64(productID)).
		Dur(
			"duration",
			time.Since(start),
		).
		Int64(
			"duration_ms",
			time.Since(start).Milliseconds(),
		).
		Msg("product created with outbox event")

	return nil
}

func (r *productRepo) FindAll(
	page,
	limit int,
	categoryID *uint64,
) ([]entity.Product, error) {

	start := time.Now()

	var products []entity.Product

	offset := (page - 1) * limit

	query := "SELECT * FROM products WHERE 1=1"

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
		&products,
		query,
		args...,
	)

	if err != nil {

		logger.Error(err).
			Str("service", "mysql").
			Str("event", "products_fetch_failed").
			Int("page", page).
			Int("limit", limit).
			Int64(
				"duration_ms",
				time.Since(start).Milliseconds(),
			).
			Msg("failed fetch products")

		return nil, err
	}

	logger.Info().
		Str("service", "mysql").
		Str("event", "products_fetched").
		Int("count", len(products)).
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
		Msg("products fetched successfully")

	return products, nil
}

func (r *productRepo) Count(
	categoryID *uint64,
) (int64, error) {

	start := time.Now()

	var total int64

	query := "SELECT COUNT(*) FROM products WHERE 1=1"

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
			Str("event", "products_count_failed").
			Int64(
				"duration_ms",
				time.Since(start).Milliseconds(),
			).
			Msg("failed count products")

		return 0, err
	}

	logger.Info().
		Str("service", "mysql").
		Str("event", "products_counted").
		Int64("total", total).
		Dur(
			"duration",
			time.Since(start),
		).
		Int64(
			"duration_ms",
			time.Since(start).Milliseconds(),
		).
		Msg("products counted successfully")

	return total, nil
}

func (r *productRepo) FindByID(
	id uint64,
) (*entity.Product, error) {

	start := time.Now()

	var product entity.Product

	err := r.db.Get(
		&product,
		"SELECT * FROM products WHERE id = ?",
		id,
	)

	if err != nil {

		logger.Error(err).
			Str("service", "mysql").
			Str("event", "product_find_by_id_failed").
			Uint64("product_id", id).
			Int64(
				"duration_ms",
				time.Since(start).Milliseconds(),
			).
			Msg("failed find product by id")

		return nil, err
	}

	logger.Info().
		Str("service", "mysql").
		Str("event", "product_found").
		Uint64("product_id", id).
		Dur(
			"duration",
			time.Since(start),
		).
		Int64(
			"duration_ms",
			time.Since(start).Milliseconds(),
		).
		Msg("product found successfully")

	return &product, nil
}

func (r *productRepo) Update(
	p *entity.Product,
) error {

	start := time.Now()

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

		logger.Error(err).
			Str("service", "mysql").
			Str("event", "product_update_failed").
			Uint64("product_id", p.ID).
			Int64(
				"duration_ms",
				time.Since(start).Milliseconds(),
			).
			Msg("failed update product")

		return err
	}

	logger.Info().
		Str("service", "mysql").
		Str("event", "product_updated").
		Uint64("product_id", p.ID).
		Dur(
			"duration",
			time.Since(start),
		).
		Int64(
			"duration_ms",
			time.Since(start).Milliseconds(),
		).
		Msg("product updated successfully")

	return nil
}

func (r *productRepo) Delete(
	id uint64,
) error {

	start := time.Now()

	_, err := r.db.Exec(
		"DELETE FROM products WHERE id = ?",
		id,
	)

	if err != nil {

		logger.Error(err).
			Str("service", "mysql").
			Str("event", "product_delete_failed").
			Uint64("product_id", id).
			Int64(
				"duration_ms",
				time.Since(start).Milliseconds(),
			).
			Msg("failed delete product")

		return err
	}

	logger.Info().
		Str("service", "mysql").
		Str("event", "product_deleted").
		Uint64("product_id", id).
		Dur(
			"duration",
			time.Since(start),
		).
		Int64(
			"duration_ms",
			time.Since(start).Milliseconds(),
		).
		Msg("product deleted successfully")

	return nil
}

func (r *productRepo) DeleteWithOutbox(
	id uint64,
	e *entity.OutboxEvent,
) error {

	start := time.Now()

	tx, err := r.db.Beginx()

	if err != nil {

		logger.Error(err).
			Str("service", "mysql").
			Str("event", "product_delete_transaction_begin_failed").
			Int64(
				"duration_ms",
				time.Since(start).Milliseconds(),
			).
			Msg("failed begin transaction delete product")

		return err
	}

	_, err = tx.Exec(
		"DELETE FROM products WHERE id = ?",
		id,
	)

	if err != nil {

		tx.Rollback()

		logger.Error(err).
			Str("service", "mysql").
			Str("event", "product_delete_transaction_failed").
			Uint64("product_id", id).
			Int64(
				"duration_ms",
				time.Since(start).Milliseconds(),
			).
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

		logger.Error(err).
			Str("service", "mysql").
			Str("event", "product_delete_outbox_insert_failed").
			Uint64("product_id", id).
			Int64(
				"duration_ms",
				time.Since(start).Milliseconds(),
			).
			Msg("failed insert outbox delete product")

		return err
	}

	err = tx.Commit()

	if err != nil {

		logger.Error(err).
			Str("service", "mysql").
			Str("event", "product_delete_commit_failed").
			Uint64("product_id", id).
			Int64(
				"duration_ms",
				time.Since(start).Milliseconds(),
			).
			Msg("failed commit delete product")

		return err
	}

	logger.Info().
		Str("service", "mysql").
		Str("event", "product_deleted_with_outbox").
		Uint64("product_id", id).
		Dur(
			"duration",
			time.Since(start),
		).
		Int64(
			"duration_ms",
			time.Since(start).Milliseconds(),
		).
		Msg("product deleted with outbox")

	return nil
}

func (r *productRepo) UpdateWithOutbox(
	p *entity.Product,
	e *entity.OutboxEvent,
) error {

	start := time.Now()

	tx, err := r.db.Beginx()

	if err != nil {

		logger.Error(err).
			Str("service", "mysql").
			Str("event", "product_update_transaction_begin_failed").
			Int64(
				"duration_ms",
				time.Since(start).Milliseconds(),
			).
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

		logger.Error(err).
			Str("service", "mysql").
			Str("event", "product_update_transaction_failed").
			Uint64("product_id", p.ID).
			Int64(
				"duration_ms",
				time.Since(start).Milliseconds(),
			).
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

		logger.Error(err).
			Str("service", "mysql").
			Str("event", "product_update_outbox_insert_failed").
			Uint64("product_id", p.ID).
			Int64(
				"duration_ms",
				time.Since(start).Milliseconds(),
			).
			Msg("failed insert outbox update product")

		return err
	}

	err = tx.Commit()

	if err != nil {

		logger.Error(err).
			Str("service", "mysql").
			Str("event", "product_update_commit_failed").
			Uint64("product_id", p.ID).
			Int64(
				"duration_ms",
				time.Since(start).Milliseconds(),
			).
			Msg("failed commit update product")

		return err
	}

	logger.Info().
		Str("service", "mysql").
		Str("event", "product_updated_with_outbox").
		Uint64("product_id", p.ID).
		Dur(
			"duration",
			time.Since(start),
		).
		Int64(
			"duration_ms",
			time.Since(start).Milliseconds(),
		).
		Msg("product updated with outbox")

	return nil
}
