package gdocs

import (
	"context"
	"fmt"
	"net/http"

	"google.golang.org/api/docs/v1"
	"google.golang.org/api/option"
)

// Client wraps the Google Docs API service.
type Client struct {
	service *docs.Service
}

// NewClient creates a new Google Docs API client using the provided authenticated HTTP client.
func NewClient(ctx context.Context, httpClient *http.Client) (*Client, error) {
	service, err := docs.NewService(ctx, option.WithHTTPClient(httpClient))
	if err != nil {
		return nil, fmt.Errorf("unable to create Docs service: %w", err)
	}

	return &Client{service: service}, nil
}

// FetchDocument retrieves a Google Docs document by its ID.
func (c *Client) FetchDocument(docID string) (*docs.Document, error) {
	doc, err := c.service.Documents.Get(docID).Do()
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve document: %w\n\nThis could mean:\n1. The document is private and you don't have permission\n2. The document doesn't exist\n3. The document ID is incorrect", err)
	}

	return doc, nil
}
