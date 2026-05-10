package elasticsearch

import (
	es8 "github.com/elastic/go-elasticsearch/v8"
)

func NewClient(url string) (*es8.Client, error) {
	cfg := es8.Config{
		Addresses: []string{url},
	}

	return es8.NewClient(cfg)
}
