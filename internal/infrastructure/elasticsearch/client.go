package elasticsearch

import (
	"time"

	"content-hub/internal/logger"

	es8 "github.com/elastic/go-elasticsearch/v8"
)

func NewClient(
	url string,
) (*es8.Client, error) {

	start := time.Now()

	cfg := es8.Config{
		Addresses: []string{
			url,
		},
	}

	client, err := es8.NewClient(
		cfg,
	)

	duration := time.Since(start).Milliseconds()

	if err != nil {

		logger.Error(err).
			Str("service", "elasticsearch").
			Str("event", "connection_failed").
			Str("url", url).
			Int64("duration_ms", duration).
			Msg("failed to connect elasticsearch")

		return nil, err
	}

	logger.Info().
		Str("service", "elasticsearch").
		Str("event", "connected").
		Str("url", url).
		Int64("duration_ms", duration).
		Msg("elasticsearch connected")

	return client, nil
}
