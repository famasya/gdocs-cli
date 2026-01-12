package markdown

import (
	"fmt"
	"strings"
	"time"

	"google.golang.org/api/docs/v1"
	"gopkg.in/yaml.v3"
)

// Frontmatter represents the YAML frontmatter for a markdown document.
type Frontmatter struct {
	Title        string    `yaml:"title"`
	Author       string    `yaml:"author,omitempty"`
	CreatedDate  time.Time `yaml:"created,omitempty"`
	ModifiedDate time.Time `yaml:"modified,omitempty"`
}

// GenerateFrontmatter creates YAML frontmatter from a Google Docs document.
func GenerateFrontmatter(doc *docs.Document) (string, error) {
	fm := Frontmatter{
		Title: doc.Title,
	}

	// Note: Google Docs API v1 doesn't provide author, created, or modified dates
	// These would need to come from the Drive API
	// For now, we'll leave them empty or use placeholders

	// Marshal to YAML
	data, err := yaml.Marshal(&fm)
	if err != nil {
		return "", fmt.Errorf("failed to marshal frontmatter: %w", err)
	}

	// Format as YAML frontmatter block
	var builder strings.Builder
	builder.WriteString("---\n")
	builder.WriteString(string(data))
	builder.WriteString("---\n")

	return builder.String(), nil
}
