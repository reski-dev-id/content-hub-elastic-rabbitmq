package entity

import "time"

type Product struct {
	ID          uint64    `db:"id" json:"id"`
	CategoryID  uint64    `db:"category_id" json:"category_id"`
	Title       string    `db:"title" json:"title"`
	Slug        string    `db:"slug" json:"slug"`
	Description string    `db:"description" json:"description"`
	Price       float64   `db:"price" json:"price"`
	Status      string    `db:"status" json:"status"`
	CreatedAt   time.Time `db:"created_at" json:"created_at"`
	UpdatedAt   time.Time `db:"updated_at" json:"updated_at"`
}
