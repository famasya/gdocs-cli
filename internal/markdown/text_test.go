package markdown

import (
	"testing"

	"google.golang.org/api/docs/v1"
)

func TestApplyTextStyle(t *testing.T) {
	tests := []struct {
		name  string
		text  string
		style *docs.TextStyle
		want  string
	}{
		{
			name:  "no style",
			text:  "plain text",
			style: nil,
			want:  "plain text",
		},
		{
			name:  "bold text",
			text:  "bold text",
			style: &docs.TextStyle{Bold: true},
			want:  "**bold text**",
		},
		{
			name:  "italic text",
			text:  "italic text",
			style: &docs.TextStyle{Italic: true},
			want:  "*italic text*",
		},
		{
			name:  "bold and italic",
			text:  "bold italic",
			style: &docs.TextStyle{Bold: true, Italic: true},
			want:  "***bold italic***",
		},
		{
			name:  "strikethrough",
			text:  "strikethrough",
			style: &docs.TextStyle{Strikethrough: true},
			want:  "~~strikethrough~~",
		},
		{
			name:  "bold and strikethrough",
			text:  "text",
			style: &docs.TextStyle{Bold: true, Strikethrough: true},
			want:  "~~**text**~~",
		},
		{
			name:  "link",
			text:  "click here",
			style: &docs.TextStyle{Link: &docs.Link{Url: "https://example.com"}},
			want:  "[click here](https://example.com)",
		},
		{
			name:  "bold link",
			text:  "click here",
			style: &docs.TextStyle{Bold: true, Link: &docs.Link{Url: "https://example.com"}},
			want:  "**[click here](https://example.com)**",
		},
		{
			name:  "link with trailing newline",
			text:  "click here\n",
			style: &docs.TextStyle{Link: &docs.Link{Url: "https://example.com"}},
			want:  "[click here](https://example.com)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ApplyTextStyle(tt.text, tt.style)
			if got != tt.want {
				t.Errorf("ApplyTextStyle() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestConvertTextRun(t *testing.T) {
	tests := []struct {
		name    string
		textRun *docs.TextRun
		want    string
	}{
		{
			name:    "nil text run",
			textRun: nil,
			want:    "",
		},
		{
			name:    "empty content",
			textRun: &docs.TextRun{Content: ""},
			want:    "",
		},
		{
			name: "plain text",
			textRun: &docs.TextRun{
				Content:   "Hello World",
				TextStyle: &docs.TextStyle{},
			},
			want: "Hello World",
		},
		{
			name: "bold text",
			textRun: &docs.TextRun{
				Content:   "Bold",
				TextStyle: &docs.TextStyle{Bold: true},
			},
			want: "**Bold**",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ConvertTextRun(tt.textRun)
			if got != tt.want {
				t.Errorf("ConvertTextRun() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestConvertParagraphElements(t *testing.T) {
	tests := []struct {
		name     string
		elements []*docs.ParagraphElement
		want     string
	}{
		{
			name:     "empty elements",
			elements: []*docs.ParagraphElement{},
			want:     "",
		},
		{
			name: "single text run",
			elements: []*docs.ParagraphElement{
				{
					TextRun: &docs.TextRun{
						Content:   "Hello",
						TextStyle: &docs.TextStyle{},
					},
				},
			},
			want: "Hello",
		},
		{
			name: "multiple text runs",
			elements: []*docs.ParagraphElement{
				{
					TextRun: &docs.TextRun{
						Content:   "Hello ",
						TextStyle: &docs.TextStyle{},
					},
				},
				{
					TextRun: &docs.TextRun{
						Content:   "World",
						TextStyle: &docs.TextStyle{Bold: true},
					},
				},
			},
			want: "Hello **World**",
		},
		{
			name: "mixed formatting",
			elements: []*docs.ParagraphElement{
				{
					TextRun: &docs.TextRun{
						Content:   "This is ",
						TextStyle: &docs.TextStyle{},
					},
				},
				{
					TextRun: &docs.TextRun{
						Content:   "bold",
						TextStyle: &docs.TextStyle{Bold: true},
					},
				},
				{
					TextRun: &docs.TextRun{
						Content:   " and ",
						TextStyle: &docs.TextStyle{},
					},
				},
				{
					TextRun: &docs.TextRun{
						Content:   "italic",
						TextStyle: &docs.TextStyle{Italic: true},
					},
				},
			},
			want: "This is **bold** and *italic*",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ConvertParagraphElements(tt.elements)
			if got != tt.want {
				t.Errorf("ConvertParagraphElements() = %q, want %q", got, tt.want)
			}
		})
	}
}
