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
	for {
		events, err := p.repo.FindPending(10)
		if err != nil {
			log.Println("error fetch outbox:", err)
			time.Sleep(2 * time.Second)
			continue
		}

		for _, e := range events {
			err := p.publisher.Publish(e.AggregateType, []byte(e.Payload))
			if err != nil {
				log.Println("publish error:", err)
				continue
			}

			_ = p.repo.MarkAsSent(e.ID)
		}

		time.Sleep(2 * time.Second)
	}
}
