package entity

import "time"

type News struct {
	ID          uint64     `db:"id" json:"id"`
	CategoryID  uint64     `db:"category_id" json:"category_id"`
	Title       string     `db:"title" json:"title"`
	Slug        string     `db:"slug" json:"slug"`
	Content     string     `db:"content" json:"content"`
	Author      string     `db:"author" json:"author"`
	Status      string     `db:"status" json:"status"`
	PublishedAt *time.Time `db:"published_at" json:"published_at"`
	CreatedAt   time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt   time.Time  `db:"updated_at" json:"updated_at"`
}
