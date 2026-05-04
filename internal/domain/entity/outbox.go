package entity

import "time"

type OutboxEvent struct {
	ID            uint64
	AggregateType string
	AggregateID   uint64
	EventType     string
	Payload       string
	Status        string
	CreatedAt     time.Time
}
