package markdown

import (
	"strings"

	"google.golang.org/api/docs/v1"
)

// ConvertTextRun converts a Google Docs TextRun to markdown with formatting.
func ConvertTextRun(textRun *docs.TextRun) string {
	if textRun == nil || textRun.Content == "" {
		return ""
	}

	text := textRun.Content
	style := textRun.TextStyle

	return ApplyTextStyle(text, style)
}

// ApplyTextStyle applies markdown formatting to text based on TextStyle.
func ApplyTextStyle(text string, style *docs.TextStyle) string {
	if style == nil {
		return text
	}

	// Handle links
	if style.Link != nil && style.Link.Url != "" {
		text = formatLink(text, style.Link.Url)
	}

	// Handle bold and italic
	// Check both combinations to apply correct markdown syntax
	if style.Bold && style.Italic {
		text = "***" + text + "***"
	} else if style.Bold {
		text = "**" + text + "**"
	} else if style.Italic {
		text = "*" + text + "*"
	}

	// Handle strikethrough
	if style.Strikethrough {
		text = "~~" + text + "~~"
	}

	return text
}

// formatLink creates a markdown link from text and URL.
func formatLink(text string, url string) string {
	// Remove trailing newlines from link text for cleaner markdown
	text = strings.TrimRight(text, "\n")
	return "[" + text + "](" + url + ")"
}

// ConvertParagraphElements converts all paragraph elements to markdown text.
func ConvertParagraphElements(elements []*docs.ParagraphElement) string {
	var builder strings.Builder

	for _, element := range elements {
		if element.TextRun != nil {
			builder.WriteString(ConvertTextRun(element.TextRun))
		}
		// Handle other element types if needed (e.g., InlineObject, PageBreak)
	}

	return builder.String()
}
