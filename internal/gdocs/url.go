package gdocs

import (
	"fmt"
	"regexp"
)

// ExtractDocumentID extracts the document ID from a Google Docs URL.
// Supports various URL formats:
// - https://docs.google.com/document/d/{DOC_ID}/edit
// - https://docs.google.com/document/d/{DOC_ID}/edit?usp=sharing
// - https://docs.google.com/document/d/{DOC_ID}/
// - https://docs.google.com/document/d/{DOC_ID}
func ExtractDocumentID(url string) (string, error) {
	// Regex pattern to match Google Docs URL and extract document ID
	pattern := `https://docs\.google\.com/document/d/([a-zA-Z0-9-_]+)`
	re := regexp.MustCompile(pattern)

	matches := re.FindStringSubmatch(url)
	if len(matches) < 2 {
		return "", fmt.Errorf("invalid Google Docs URL: expected format 'https://docs.google.com/document/d/{DOC_ID}...', got '%s'", url)
	}

	return matches[1], nil
}

// ExtractTabID extracts the tab ID from a Google Docs URL if present.
// Tab IDs appear in URLs as ?tab=t.{TAB_ID} or &tab=t.{TAB_ID}
// Returns empty string if no tab ID is found.
func ExtractTabID(url string) string {
	// Regex pattern to match tab parameter (e.g., tab=t.v63b7x227gkk)
	pattern := `[?&]tab=(t\.[a-zA-Z0-9]+)`
	re := regexp.MustCompile(pattern)

	matches := re.FindStringSubmatch(url)
	if len(matches) < 2 {
		return ""
	}

	return matches[1]
}
