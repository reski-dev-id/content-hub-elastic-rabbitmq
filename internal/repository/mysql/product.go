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
	query := `INSERT INTO products (category_id, title, price, status)
			  VALUES (?, ?, ?, ?)`

	_, err := r.db.Exec(query, p.CategoryID, p.Title, p.Price, p.Status)
	return err
}

func (r *productRepo) FindAll() ([]entity.Product, error) {
	var products []entity.Product
	err := r.db.Select(&products, "SELECT * FROM products")
	return products, err
}
