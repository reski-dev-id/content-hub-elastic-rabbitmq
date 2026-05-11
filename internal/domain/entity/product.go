package entity

import "time"

type Product struct {
	ID          uint64    `db:"id" json:"id" example:"1"`
	CategoryID  uint64    `db:"category_id" json:"category_id" example:"1"`
	Title       string    `db:"title" json:"title" example:"MacBook Pro M4"`
	Slug        string    `db:"slug" json:"slug" example:"macbook-pro-m4"`
	Description string    `db:"description" json:"description" example:"Laptop Apple terbaru"`
	Price       float64   `db:"price" json:"price" example:"42000000"`
	Status      string    `db:"status" json:"status" example:"active"`
	CreatedAt   time.Time `db:"created_at" json:"created_at"`
	UpdatedAt   time.Time `db:"updated_at" json:"updated_at"`
}
