package markdown

import (
	"testing"

	"github.com/famasya/gdocs-cli/internal/gdocs"
)

func TestConvertComments(t *testing.T) {
	tests := []struct {
		name     string
		comments []gdocs.Comment
		want     string
	}{
		{
			name:     "nil comments",
			comments: nil,
			want:     "",
		},
		{
			name:     "empty comments",
			comments: []gdocs.Comment{},
			want:     "",
		},
		{
			name: "single comment without quoted text",
			comments: []gdocs.Comment{
				{
					Author:      "Alice",
					Content:     "This needs clarification.",
					CreatedTime: "2025-01-15T10:30:00Z",
				},
			},
			want: "## Comments\n\n**Alice** (2025-01-15): This needs clarification.\n\n",
		},
		{
			name: "single comment with quoted text",
			comments: []gdocs.Comment{
				{
					Author:      "Bob",
					Content:     "Typo here.",
					QuotedText:  "the orignal text",
					CreatedTime: "2025-03-20T14:00:00Z",
				},
			},
			want: "## Comments\n\n> the orignal text\n\n**Bob** (2025-03-20): Typo here.\n\n",
		},
		{
			name: "resolved comment",
			comments: []gdocs.Comment{
				{
					Author:      "Carol",
					Content:     "Fixed now.",
					CreatedTime: "2025-02-01T08:00:00Z",
					Resolved:    true,
				},
			},
			want: "## Comments\n\n**Carol** (2025-02-01) ✓ resolved: Fixed now.\n\n",
		},
		{
			name: "comment with replies",
			comments: []gdocs.Comment{
				{
					Author:      "Dave",
					Content:     "Should we change this?",
					CreatedTime: "2025-04-10T12:00:00Z",
					Replies: []gdocs.Reply{
						{
							Author:      "Eve",
							Content:     "Yes, I agree.",
							CreatedTime: "2025-04-10T13:00:00Z",
						},
						{
							Author:      "Dave",
							Content:     "Done.",
							CreatedTime: "2025-04-10T14:00:00Z",
						},
					},
				},
			},
			want: "## Comments\n\n**Dave** (2025-04-10): Should we change this?\n  ↳ **Eve** (2025-04-10): Yes, I agree.\n  ↳ **Dave** (2025-04-10): Done.\n\n",
		},
		{
			name: "comment with multiline quoted text",
			comments: []gdocs.Comment{
				{
					Author:      "Frank",
					Content:     "This paragraph is too long.",
					QuotedText:  "line one\nline two",
					CreatedTime: "2025-05-01T09:00:00Z",
				},
			},
			want: "## Comments\n\n> line one\n> line two\n\n**Frank** (2025-05-01): This paragraph is too long.\n\n",
		},
		{
			name: "comment with unknown author",
			comments: []gdocs.Comment{
				{
					Content:     "Anonymous feedback.",
					CreatedTime: "2025-06-01T10:00:00Z",
				},
			},
			want: "## Comments\n\n**Unknown** (2025-06-01): Anonymous feedback.\n\n",
		},
		{
			name: "comment with no timestamp",
			comments: []gdocs.Comment{
				{
					Author:  "Grace",
					Content: "No date here.",
				},
			},
			want: "## Comments\n\n**Grace**: No date here.\n\n",
		},
		{
			name: "multiple comments",
			comments: []gdocs.Comment{
				{
					Author:      "Alice",
					Content:     "First comment.",
					CreatedTime: "2025-01-01T00:00:00Z",
				},
				{
					Author:      "Bob",
					Content:     "Second comment.",
					QuotedText:  "some text",
					CreatedTime: "2025-01-02T00:00:00Z",
				},
			},
			want: "## Comments\n\n**Alice** (2025-01-01): First comment.\n\n> some text\n\n**Bob** (2025-01-02): Second comment.\n\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ConvertComments(tt.comments)
			if got != tt.want {
				t.Errorf("ConvertComments() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestFormatTime(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    string
	}{
		{
			name:  "valid RFC3339",
			input: "2025-01-15T10:30:00Z",
			want:  "2025-01-15",
		},
		{
			name:  "valid RFC3339 with offset",
			input: "2025-06-20T14:30:00+07:00",
			want:  "2025-06-20",
		},
		{
			name:  "empty string",
			input: "",
			want:  "",
		},
		{
			name:  "invalid format",
			input: "not-a-date",
			want:  "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := formatTime(tt.input)
			if got != tt.want {
				t.Errorf("formatTime(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}
