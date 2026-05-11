package entity

import "time"

type News struct {
	ID          uint64     `db:"id" json:"id" example:"1"`
	CategoryID  uint64     `db:"category_id" json:"category_id" example:"2"`
	Title       string     `db:"title" json:"title" example:"OpenAI Rilis GPT Baru"`
	Slug        string     `db:"slug" json:"slug" example:"openai-rilis-gpt-baru"`
	Content     string     `db:"content" json:"content" example:"Model terbaru membawa peningkatan reasoning dan performa."`
	Author      string     `db:"author" json:"author" example:"Reski"`
	Status      string     `db:"status" json:"status" example:"published"`
	PublishedAt *time.Time `db:"published_at" json:"published_at"`
	CreatedAt   time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt   time.Time  `db:"updated_at" json:"updated_at"`
}
