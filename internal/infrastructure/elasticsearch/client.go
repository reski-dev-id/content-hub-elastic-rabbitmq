package elasticsearch

import (
	"github.com/elastic/go-elasticsearch/v8"
)

func NewClient(url string) (*elasticsearch.Client, error) {
	cfg := elasticsearch.Config{
		Addresses: []string{url},
	}

	return elasticsearch.NewClient(cfg)
}
