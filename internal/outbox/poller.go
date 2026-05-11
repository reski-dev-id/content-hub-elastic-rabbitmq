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
	return &Poller{r, p}
}

func (p *Poller) Start() {

	logger.Log.Info().
		Msg("outbox poller started")

	for {

		events, err := p.repo.FindPending(10)

		if err != nil {

			logger.Log.Error().
				Err(err).
				Msg("failed fetch pending outbox events")

			time.Sleep(2 * time.Second)

			continue
		}

		if len(events) > 0 {

			logger.Log.Info().
				Int("count", len(events)).
				Msg("pending outbox events fetched")
		}

		for _, e := range events {

			logger.Log.Info().
				Uint64("event_id", e.ID).
				Str("event_type", e.EventType).
				Str("aggregate_type", e.AggregateType).
				Msg("publishing outbox event")

			err := p.publisher.Publish(
				e.AggregateType,
				[]byte(e.Payload),
			)

			if err != nil {

				logger.Log.Error().
					Err(err).
					Uint64("event_id", e.ID).
					Msg("failed publish rabbitmq message")

				continue
			}

			logger.Log.Info().
				Uint64("event_id", e.ID).
				Msg("outbox event published")

			err = p.repo.MarkAsSent(e.ID)

			if err != nil {

				logger.Log.Error().
					Err(err).
					Uint64("event_id", e.ID).
					Msg("failed mark outbox event as sent")

				continue
			}

			logger.Log.Info().
				Uint64("event_id", e.ID).
				Msg("outbox event marked as sent")
		}

		time.Sleep(2 * time.Second)
	}
}
