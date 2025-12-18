// Package elastic provides an Elasticsearch client for bulk indexing.
package elastic

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// Client is an Elasticsearch client for bulk indexing operations.
type Client struct {
	baseURL  string
	username string
	password string
	client   *http.Client
}

// NewClient creates a new Elasticsearch client.
func NewClient(url, username, password string) *Client {
	return &Client{
		baseURL:  url,
		username: username,
		password: password,
		client:   &http.Client{},
	}
}

// BulkResult contains statistics from a bulk indexing operation.
type BulkResult struct {
	Total      int
	Successful int
	Failed     int
}

// BulkIndex indexes documents using the Elasticsearch bulk API.
// It batches documents and sends them in chunks.
func (c *Client) BulkIndex(ctx context.Context, index string, docs <-chan Document) (*BulkResult, error) {
	var buf bytes.Buffer
	result := &BulkResult{}
	batchSize := 500

	for doc := range docs {
		// Action line
		action := map[string]any{
			"index": map[string]any{
				"_index": index,
				"_id":    doc.ID,
			},
		}
		if err := json.NewEncoder(&buf).Encode(action); err != nil {
			return result, fmt.Errorf("encoding action: %w", err)
		}

		// Document line
		if err := json.NewEncoder(&buf).Encode(doc.Body); err != nil {
			return result, fmt.Errorf("encoding document: %w", err)
		}

		result.Total++

		// Flush batch
		if result.Total%batchSize == 0 {
			if err := c.flush(ctx, &buf); err != nil {
				return result, err
			}
		}
	}

	// Flush remaining
	if buf.Len() > 0 {
		if err := c.flush(ctx, &buf); err != nil {
			return result, err
		}
	}

	result.Successful = result.Total - result.Failed
	return result, nil
}

func (c *Client) flush(ctx context.Context, buf *bytes.Buffer) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+"/_bulk", buf)
	if err != nil {
		return fmt.Errorf("creating request: %w", err)
	}
	req.Header.Set("Content-Type", "application/x-ndjson")
	req.SetBasicAuth(c.username, c.password)

	resp, err := c.client.Do(req)
	if err != nil {
		return fmt.Errorf("executing request: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode >= 400 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("bulk request failed with status %d: %s", resp.StatusCode, body)
	}

	buf.Reset()
	return nil
}

// Document represents a document to be indexed.
type Document struct {
	ID   string
	Body any
}
