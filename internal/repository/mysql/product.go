package mysql

import (
	"content-hub/internal/domain/entity"
	"content-hub/internal/domain/repository"

	"github.com/jmoiron/sqlx"
)

type productRepo struct {
	db *sqlx.DB
}

func NewProductRepository(db *sqlx.DB) repository.ProductRepository {
	return &productRepo{db}
}

func (r *productRepo) Create(p *entity.Product) error {
	query := `
	INSERT INTO products (category_id, title, slug, description, price, status)
	VALUES (?, ?, ?, ?, ?, ?)`

	_, err := r.db.Exec(query,
		p.CategoryID,
		p.Title,
		p.Slug,
		p.Description,
		p.Price,
		p.Status,
	)

	return err
}

func (r *productRepo) CreateWithOutbox(p *entity.Product, e *entity.OutboxEvent) error {
	tx, err := r.db.Beginx()
	if err != nil {
		return err
	}

	res, err := tx.Exec(`
		INSERT INTO products (category_id, title, slug, description, price, status)
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
		return err
	}

	productID, err := res.LastInsertId()
	if err != nil {
		tx.Rollback()
		return err
	}

	_, err = tx.Exec(`
		INSERT INTO outbox_events (aggregate_type, aggregate_id, event_type, payload, status)
		VALUES (?, ?, ?, ?, ?)`,
		e.AggregateType,
		productID,
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

func (r *productRepo) FindAll(page, limit int, categoryID *uint64) ([]entity.Product, error) {
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
	return products, err
}

func (r *productRepo) FindByID(id uint64) (*entity.Product, error) {
	var product entity.Product

	err := r.db.Get(&product, "SELECT * FROM products WHERE id = ?", id)
	if err != nil {
		return nil, err
	}

	return &product, nil
}

func (r *productRepo) Update(p *entity.Product) error {
	query := `
	UPDATE products
	SET category_id=?, title=?, slug=?, description=?, price=?, status=?
	WHERE id=?`

	_, err := r.db.Exec(query,
		p.CategoryID,
		p.Title,
		p.Slug,
		p.Description,
		p.Price,
		p.Status,
		p.ID,
	)

	return err
}

func (r *productRepo) Delete(id uint64) error {
	_, err := r.db.Exec(
		"DELETE FROM products WHERE id = ?",
		id,
	)

	return err
}

func (r *productRepo) DeleteWithOutbox(id uint64, e *entity.OutboxEvent) error {
	tx, err := r.db.Beginx()
	if err != nil {
		return err
	}

	_, err = tx.Exec(
		"DELETE FROM products WHERE id = ?",
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

func (r *productRepo) UpdateWithOutbox(
	p *entity.Product,
	e *entity.OutboxEvent,
) error {
	tx, err := r.db.Beginx()
	if err != nil {
		return err
	}

	_, err = tx.Exec(`
		UPDATE products
		SET category_id=?, title=?, slug=?, description=?, price=?, status=?
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
		return err
	}

	return tx.Commit()
}
