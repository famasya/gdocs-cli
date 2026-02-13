package markdown

import (
	"fmt"
	"strings"

	"github.com/famasya/gdocs-cli/internal/gdocs"
	"google.golang.org/api/docs/v1"
)

// Converter handles the conversion of Google Docs to markdown.
type Converter struct {
	doc      *docs.Document
	body     *docs.Body
	title    string
	tabName  string
	comments []gdocs.Comment
}

// NewConverter creates a new Converter for the given document.
// Uses the first tab's content by default.
func NewConverter(doc *docs.Document) *Converter {
	c := &Converter{doc: doc, title: doc.Title}

	// Use tab content if available, otherwise fall back to legacy doc.Body
	if doc.Tabs != nil && len(doc.Tabs) > 0 {
		tab := doc.Tabs[0]
		if tab.DocumentTab != nil {
			c.body = tab.DocumentTab.Body
		}
		if tab.TabProperties != nil {
			c.tabName = tab.TabProperties.Title
		}
	} else if doc.Body != nil {
		c.body = doc.Body
	}

	return c
}

// NewConverterFromTab creates a new Converter for a specific tab.
func NewConverterFromTab(doc *docs.Document, tab *docs.Tab) *Converter {
	c := &Converter{doc: doc, title: doc.Title}

	if tab != nil && tab.DocumentTab != nil {
		c.body = tab.DocumentTab.Body
		if tab.TabProperties != nil {
			c.tabName = tab.TabProperties.Title
		}
	}

	return c
}

// SetComments sets the comments to be appended to the markdown output.
func (c *Converter) SetComments(comments []gdocs.Comment) {
	c.comments = comments
}

// Convert processes the entire document and returns markdown.
func (c *Converter) Convert() (string, error) {
	var builder strings.Builder

	// Generate frontmatter
	frontmatter, err := c.generateFrontmatter()
	if err != nil {
		return "", fmt.Errorf("failed to generate frontmatter: %w", err)
	}
	builder.WriteString(frontmatter)
	builder.WriteString("\n")

	// Convert body content
	if c.body != nil && c.body.Content != nil {
		body := c.convertBody()
		builder.WriteString(body)
	}

	// Append comments if present
	if len(c.comments) > 0 {
		builder.WriteString(ConvertComments(c.comments))
	}

	return builder.String(), nil
}

// generateFrontmatter creates frontmatter including tab info if present.
func (c *Converter) generateFrontmatter() (string, error) {
	// Use the existing GenerateFrontmatter for the base, but we'll
	// add tab info if we have it
	frontmatter, err := GenerateFrontmatter(c.doc)
	if err != nil {
		return "", err
	}

	// If we have a tab name that differs from the doc title, include it
	if c.tabName != "" && c.tabName != c.title {
		// Insert tab info before the closing ---
		frontmatter = strings.TrimSuffix(frontmatter, "---\n")
		frontmatter += fmt.Sprintf("tab: %s\n---\n", c.tabName)
	}

	return frontmatter, nil
}

// convertBody converts the document body to markdown.
func (c *Converter) convertBody() string {
	var builder strings.Builder

	for _, element := range c.body.Content {
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
