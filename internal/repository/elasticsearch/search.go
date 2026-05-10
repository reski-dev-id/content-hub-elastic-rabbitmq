package elasticsearch

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"

	es8 "github.com/elastic/go-elasticsearch/v8"
)

type SearchRepository struct {
	client *es8.Client
}

func NewSearchRepository(
	client *es8.Client,
) *SearchRepository {
	return &SearchRepository{
		client: client,
	}
}

func (r *SearchRepository) Search(
	ctx context.Context,
	index string,
	q string,
	categoryID *uint64,
	page int,
	limit int,
) (map[string]interface{}, error) {

	from := (page - 1) * limit

	must := []map[string]interface{}{
		{
			"multi_match": map[string]interface{}{
				"query": q,
				"fields": []string{
					"title",
					"description",
					"content",
				},
			},
		},
	}

	if categoryID != nil {
		must = append(must,
			map[string]interface{}{
				"term": map[string]interface{}{
					"category_id": *categoryID,
				},
			},
		)
	}

	query := map[string]interface{}{
		"from": from,
		"size": limit,
		"query": map[string]interface{}{
			"bool": map[string]interface{}{
				"must": must,
			},
		},
	}

	body, _ := json.Marshal(query)

	res, err := r.client.Search(
		r.client.Search.WithContext(ctx),
		r.client.Search.WithIndex(index),
		r.client.Search.WithBody(
			bytes.NewReader(body),
		),
		r.client.Search.WithTrackTotalHits(true),
	)

	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	if res.IsError() {
		return nil, fmt.Errorf(res.String())
	}

	var result map[string]interface{}

	err = json.NewDecoder(res.Body).
		Decode(&result)

	return result, err
}
