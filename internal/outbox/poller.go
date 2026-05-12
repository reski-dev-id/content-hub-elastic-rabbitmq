package outbox

import (
	"time"

	"content-hub/internal/domain/repository"
	"content-hub/internal/logger"
)

type Publisher interface {
	Publish(queue string, body []byte) error
}

type Poller struct {
	repo      repository.OutboxRepository
	publisher Publisher
}

func NewPoller(
	r repository.OutboxRepository,
	p Publisher,
) *Poller {
	return &Poller{
		repo:      r,
		publisher: p,
	}
}

func (p *Poller) Start() {

	logger.Info().
		Str("service", "poller").
		Str("event", "poller_started").
		Msg("outbox poller started")

	for {

		cycleStart := time.Now()

		events, err := p.repo.FindPending(10)

		if err != nil {

			logger.Error(err).
				Str("service", "poller").
				Str("event", "pending_events_fetch_failed").
				Int64(
					"duration_ms",
					time.Since(cycleStart).Milliseconds(),
				).
				Msg("failed fetch pending outbox events")

			time.Sleep(2 * time.Second)

			continue
		}

		if len(events) > 0 {

			logger.Info().
				Str("service", "poller").
				Str("event", "pending_events_fetched").
				Int("count", len(events)).
				Dur(
					"duration",
					time.Since(cycleStart),
				).
				Int64(
					"duration_ms",
					time.Since(cycleStart).Milliseconds(),
				).
				Msg("pending outbox events fetched")
		}

		for _, e := range events {

			eventStart := time.Now()

			logger.Info().
				Str("service", "poller").
				Str("event", "outbox_event_publishing").
				Uint64("event_id", e.ID).
				Str("event_type", e.EventType).
				Str("aggregate_type", e.AggregateType).
				Msg("publishing outbox event")

			err := p.publisher.Publish(
				e.AggregateType,
				[]byte(e.Payload),
			)

			if err != nil {

				logger.Error(err).
					Str("service", "poller").
					Str("event", "outbox_publish_failed").
					Uint64("event_id", e.ID).
					Str("event_type", e.EventType).
					Str("aggregate_type", e.AggregateType).
					Int64(
						"duration_ms",
						time.Since(eventStart).Milliseconds(),
					).
					Msg("failed publish rabbitmq message")

				continue
			}

			logger.Info().
				Str("service", "poller").
				Str("event", "outbox_published").
				Uint64("event_id", e.ID).
				Str("event_type", e.EventType).
				Str("aggregate_type", e.AggregateType).
				Dur(
					"duration",
					time.Since(eventStart),
				).
				Int64(
					"duration_ms",
					time.Since(eventStart).Milliseconds(),
				).
				Msg("outbox event published")

			err = p.repo.MarkAsSent(
				e.ID,
			)

			if err != nil {

				logger.Error(err).
					Str("service", "poller").
					Str("event", "outbox_mark_sent_failed").
					Uint64("event_id", e.ID).
					Int64(
						"duration_ms",
						time.Since(eventStart).Milliseconds(),
					).
					Msg("failed mark outbox event as sent")

				continue
			}

			logger.Info().
				Str("service", "poller").
				Str("event", "outbox_marked_sent").
				Uint64("event_id", e.ID).
				Dur(
					"duration",
					time.Since(eventStart),
				).
				Int64(
					"duration_ms",
					time.Since(eventStart).Milliseconds(),
				).
				Msg("outbox event marked as sent")
		}

		logger.Info().
			Str("service", "poller").
			Str("event", "poller_cycle_completed").
			Int("processed_events", len(events)).
			Dur(
				"duration",
				time.Since(cycleStart),
			).
			Int64(
				"duration_ms",
				time.Since(cycleStart).Milliseconds(),
			).
			Msg("poller cycle completed")

		time.Sleep(2 * time.Second)
	}
}
