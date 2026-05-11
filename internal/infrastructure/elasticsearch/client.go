package elasticsearch

import (
	"content-hub/internal/logger"

	es8 "github.com/elastic/go-elasticsearch/v8"
)

func NewClient(url string) (*es8.Client, error) {

	cfg := es8.Config{
		Addresses: []string{url},
	}

	client, err := es8.NewClient(cfg)

	if err != nil {

		logger.Log.Error().
			Err(err).
			Msg("failed connect elasticsearch")

		return nil, err
	}

	logger.Log.Info().
		Str("url", url).
		Msg("elasticsearch connected")

	return client, nil
}
