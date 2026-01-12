package markdown

import (
	"fmt"
	"strings"

	"google.golang.org/api/docs/v1"
)

// Converter handles the conversion of Google Docs to markdown.
type Converter struct {
	doc *docs.Document
}

// NewConverter creates a new Converter for the given document.
func NewConverter(doc *docs.Document) *Converter {
	return &Converter{doc: doc}
}

// Convert processes the entire document and returns markdown.
func (c *Converter) Convert() (string, error) {
	var builder strings.Builder

	// Generate frontmatter
	frontmatter, err := GenerateFrontmatter(c.doc)
	if err != nil {
		return "", fmt.Errorf("failed to generate frontmatter: %w", err)
	}
	builder.WriteString(frontmatter)
	builder.WriteString("\n")

	// Convert body content
	if c.doc.Body != nil && c.doc.Body.Content != nil {
		body := c.convertBody()
		builder.WriteString(body)
	}

	return builder.String(), nil
}

// convertBody converts the document body to markdown.
func (c *Converter) convertBody() string {
	var builder strings.Builder

	for _, element := range c.doc.Body.Content {
		// Convert based on element type
		if element.Paragraph != nil {
			markdown := ConvertParagraph(element.Paragraph, element.Paragraph.ParagraphStyle)
			builder.WriteString(markdown)
		} else if element.Table != nil {
			markdown := ConvertTable(element.Table)
			builder.WriteString(markdown)
		}
		// Other structural elements can be added here (e.g., SectionBreak, TableOfContents)
	}

	return builder.String()
}
