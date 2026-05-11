package entity

import "time"

type OutboxEvent struct {
	ID            uint64     `db:"id" json:"id" example:"1"`
	AggregateType string     `db:"aggregate_type" json:"aggregate_type" example:"product"`
	AggregateID   uint64     `db:"aggregate_id" json:"aggregate_id" example:"10"`
	EventType     string     `db:"event_type" json:"event_type" example:"product_created"`
	Payload       string     `db:"payload" json:"payload" example:"{\"id\":10,\"title\":\"MacBook Pro M4\"}"`
	Status        string     `db:"status" json:"status" example:"pending"`
	CreatedAt     time.Time  `db:"created_at" json:"created_at"`
	SentAt        *time.Time `db:"sent_at" json:"sent_at"`
}
