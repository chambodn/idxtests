package cmd

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/elastic/go-elasticsearch/v7"
	"github.com/elastic/go-elasticsearch/v7/esapi"
)

// ResultsConfig configures the results.
type ResultsConfig struct {
	Client    *elasticsearch.Client
	IndexName string
}

// Results allows to index test results.
type Results struct {
	es        *elasticsearch.Client
	indexName string
}

// NewResults returns a new instance of the test results.
func NewResults(c ResultsConfig) (*Results, error) {
	indexName := c.IndexName
	if indexName == "" {
		indexName = "test-suites"
	}

	r := Results{es: c.Client, indexName: indexName}
	return &r, nil
}

// CreateIndex creates a new index with mapping
func (r *Results) CreateIndex(mapping string) error {
	res, err := r.es.Indices.Create(r.indexName, r.es.Indices.Create.WithBody(strings.NewReader((mapping))))
	if err != nil {
		return err
	}
	if res.IsError() {
		return fmt.Errorf("error: %s", res)
	}

	return nil
}

// Create indexes a new test result into associated index.
func (r *Results) Create(item *Test) error {
	payload, err := json.Marshal(item)
	if err != nil {
		return err
	}
	ctx := context.Background()
	res, err := esapi.CreateRequest{
		Index: r.indexName,
		Body:  bytes.NewReader(payload),
	}.Do(ctx, r.es)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.IsError() {
		var e map[string]interface{}
		if err := json.NewDecoder(res.Body).Decode(&e); err != nil {
			return err
		}
		return fmt.Errorf("[%s] %s: %s", res.Status(), e["error"].(map[string]interface{})["type"], e["error"].(map[string]interface{})["reason"])
	}
	return nil
}
