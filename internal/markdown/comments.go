package markdown

import (
	"fmt"
	"strings"
	"time"

	"github.com/famasya/gdocs-cli/internal/gdocs"
)

// ConvertComments renders a list of comments as a markdown section.
func ConvertComments(comments []gdocs.Comment) string {
	if len(comments) == 0 {
		return ""
	}

	var builder strings.Builder
	builder.WriteString("## Comments\n\n")

	for _, c := range comments {
		if c.QuotedText != "" {
			builder.WriteString("> ")
			builder.WriteString(strings.ReplaceAll(c.QuotedText, "\n", "\n> "))
			builder.WriteString("\n\n")
		}

		author := c.Author
		if author == "" {
			author = "Unknown"
		}
		author = escapeMarkdown(author)
		builder.WriteString(fmt.Sprintf("**%s**", author))
		if ts := formatTime(c.CreatedTime); ts != "" {
			builder.WriteString(fmt.Sprintf(" (%s)", ts))
		}
		if c.Resolved {
			builder.WriteString(" ✓ resolved")
		}
		builder.WriteString(": ")
		builder.WriteString(c.Content)
		builder.WriteString("\n")

		for _, r := range c.Replies {
			rAuthor := r.Author
			if rAuthor == "" {
				rAuthor = "Unknown"
			}
			rAuthor = escapeMarkdown(rAuthor)
			builder.WriteString(fmt.Sprintf("  ↳ **%s**", rAuthor))
			if ts := formatTime(r.CreatedTime); ts != "" {
				builder.WriteString(fmt.Sprintf(" (%s)", ts))
			}
			builder.WriteString(": ")
			builder.WriteString(r.Content)
			builder.WriteString("\n")
		}

		builder.WriteString("\n")
	}

	return builder.String()
}

// formatTime converts an RFC 3339 timestamp to a short date string.
func formatTime(rfc3339 string) string {
	if rfc3339 == "" {
		return ""
	}
	t, err := time.Parse(time.RFC3339, rfc3339)
	if err != nil {
		return ""
	}
	return t.Format("2006-01-02")
}

func escapeMarkdown(s string) string {
	s = strings.ReplaceAll(s, "*", "\\*")
	s = strings.ReplaceAll(s, "_", "\\_")
	return s
}
