package entity

import "time"

type OutboxEvent struct {
	ID            uint64     `db:"id"`
	AggregateType string     `db:"aggregate_type"`
	AggregateID   uint64     `db:"aggregate_id"`
	EventType     string     `db:"event_type"`
	Payload       string     `db:"payload"`
	Status        string     `db:"status"`
	CreatedAt     time.Time  `db:"created_at"`
	SentAt        *time.Time `db:"sent_at"`
}
