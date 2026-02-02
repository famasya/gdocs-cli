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

// FetchDocument retrieves a Google Docs document by its ID with all tabs included.
func (c *Client) FetchDocument(docID string) (*docs.Document, error) {
	doc, err := c.service.Documents.Get(docID).IncludeTabsContent(true).Do()
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve document: %w\n\nThis could mean:\n1. The document is private and you don't have permission\n2. The document doesn't exist\n3. The document ID is incorrect", err)
	}

	return doc, nil
}

// FindTab searches for a tab by ID in the document's tab tree.
// Returns nil if the tab is not found.
func FindTab(doc *docs.Document, tabID string) *docs.Tab {
	if doc.Tabs == nil {
		return nil
	}

	for _, tab := range doc.Tabs {
		if found := findTabRecursive(tab, tabID); found != nil {
			return found
		}
	}

	return nil
}

// findTabRecursive recursively searches for a tab by ID.
func findTabRecursive(tab *docs.Tab, tabID string) *docs.Tab {
	if tab.TabProperties != nil && tab.TabProperties.TabId == tabID {
		return tab
	}

	for _, child := range tab.ChildTabs {
		if found := findTabRecursive(child, tabID); found != nil {
			return found
		}
	}

	return nil
}

// GetFirstTab returns the first tab in the document.
// Returns nil if the document has no tabs.
func GetFirstTab(doc *docs.Document) *docs.Tab {
	if doc.Tabs == nil || len(doc.Tabs) == 0 {
		return nil
	}
	return doc.Tabs[0]
}
