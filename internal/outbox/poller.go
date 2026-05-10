package outbox

import (
	"log"
	"time"

	"content-hub/internal/domain/repository"
)

type Publisher interface {
	Publish(queue string, body []byte) error
}

type Poller struct {
	repo      repository.OutboxRepository
	publisher Publisher
}

func NewPoller(r repository.OutboxRepository, p Publisher) *Poller {
	return &Poller{r, p}
}

func (p *Poller) Start() {
	log.Println("outbox poller started")

	for {
		events, err := p.repo.FindPending(10)
		if err != nil {
			log.Println("error fetch outbox:", err)
			time.Sleep(2 * time.Second)
			continue
		}

		if len(events) > 0 {
			log.Printf("found %d pending events\n", len(events))
		}

		for _, e := range events {
			log.Printf(
				"publishing event id=%d type=%s aggregate=%s\n",
				e.ID,
				e.EventType,
				e.AggregateType,
			)

			err := p.publisher.Publish(
				e.AggregateType,
				[]byte(e.Payload),
			)

			if err != nil {
				log.Println("publish error:", err)
				continue
			}

			log.Printf("event id=%d published\n", e.ID)

			err = p.repo.MarkAsSent(e.ID)
			if err != nil {
				log.Println("mark as sent error:", err)
				continue
			}

			log.Printf("event id=%d marked as sent\n", e.ID)
		}

		time.Sleep(2 * time.Second)
	}
}
