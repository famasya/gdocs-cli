package markdown

import (
	"testing"

	"google.golang.org/api/docs/v1"
)

func TestConvertParagraph(t *testing.T) {
	tests := []struct {
		name  string
		para  *docs.Paragraph
		style *docs.ParagraphStyle
		want  string
	}{
		{
			name:  "nil paragraph",
			para:  nil,
			style: nil,
			want:  "",
		},
		{
			name: "empty paragraph",
			para: &docs.Paragraph{
				Elements: []*docs.ParagraphElement{
					{
						TextRun: &docs.TextRun{Content: "\n"},
					},
				},
			},
			style: &docs.ParagraphStyle{},
			want:  "\n",
		},
		{
			name: "normal paragraph",
			para: &docs.Paragraph{
				Elements: []*docs.ParagraphElement{
					{
						TextRun: &docs.TextRun{Content: "This is a paragraph.\n"},
					},
				},
			},
			style: &docs.ParagraphStyle{NamedStyleType: "NORMAL_TEXT"},
			want:  "This is a paragraph.\n\n",
		},
		{
			name: "heading 1",
			para: &docs.Paragraph{
				Elements: []*docs.ParagraphElement{
					{
						TextRun: &docs.TextRun{Content: "Heading 1\n"},
					},
				},
			},
			style: &docs.ParagraphStyle{NamedStyleType: "HEADING_1"},
			want:  "# Heading 1\n\n",
		},
		{
			name: "heading 2",
			para: &docs.Paragraph{
				Elements: []*docs.ParagraphElement{
					{
						TextRun: &docs.TextRun{Content: "Heading 2\n"},
					},
				},
			},
			style: &docs.ParagraphStyle{NamedStyleType: "HEADING_2"},
			want:  "## Heading 2\n\n",
		},
		{
			name: "heading 3",
			para: &docs.Paragraph{
				Elements: []*docs.ParagraphElement{
					{
						TextRun: &docs.TextRun{Content: "Heading 3\n"},
					},
				},
			},
			style: &docs.ParagraphStyle{NamedStyleType: "HEADING_3"},
			want:  "### Heading 3\n\n",
		},
		{
			name: "title",
			para: &docs.Paragraph{
				Elements: []*docs.ParagraphElement{
					{
						TextRun: &docs.TextRun{Content: "Document Title\n"},
					},
				},
			},
			style: &docs.ParagraphStyle{NamedStyleType: "TITLE"},
			want:  "# Document Title\n\n",
		},
		{
			name: "subtitle",
			para: &docs.Paragraph{
				Elements: []*docs.ParagraphElement{
					{
						TextRun: &docs.TextRun{Content: "Subtitle\n"},
					},
				},
			},
			style: &docs.ParagraphStyle{NamedStyleType: "SUBTITLE"},
			want:  "## Subtitle\n\n",
		},
		{
			name: "bullet list item - level 0",
			para: &docs.Paragraph{
				Elements: []*docs.ParagraphElement{
					{
						TextRun: &docs.TextRun{Content: "First item\n"},
					},
				},
				Bullet: &docs.Bullet{
					ListId:       "list1",
					NestingLevel: 0,
				},
			},
			style: &docs.ParagraphStyle{},
			want:  "- First item\n",
		},
		{
			name: "bullet list item - level 1",
			para: &docs.Paragraph{
				Elements: []*docs.ParagraphElement{
					{
						TextRun: &docs.TextRun{Content: "Nested item\n"},
					},
				},
				Bullet: &docs.Bullet{
					ListId:       "list1",
					NestingLevel: 1,
				},
			},
			style: &docs.ParagraphStyle{},
			want:  "  - Nested item\n",
		},
		{
			name: "bullet list item - level 2",
			para: &docs.Paragraph{
				Elements: []*docs.ParagraphElement{
					{
						TextRun: &docs.TextRun{Content: "Double nested\n"},
					},
				},
				Bullet: &docs.Bullet{
					ListId:       "list1",
					NestingLevel: 2,
				},
			},
			style: &docs.ParagraphStyle{},
			want:  "    - Double nested\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ConvertParagraph(tt.para, tt.style)
			if got != tt.want {
				t.Errorf("ConvertParagraph() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestConvertTable(t *testing.T) {
	tests := []struct {
		name  string
		table *docs.Table
		want  string
	}{
		{
			name:  "nil table",
			table: nil,
			want:  "",
		},
		{
			name: "empty table",
			table: &docs.Table{
				TableRows: []*docs.TableRow{},
			},
			want: "",
		},
		{
			name: "simple 2x2 table",
			table: &docs.Table{
				TableRows: []*docs.TableRow{
					{
						TableCells: []*docs.TableCell{
							{
								Content: []*docs.StructuralElement{
									{
										Paragraph: &docs.Paragraph{
											Elements: []*docs.ParagraphElement{
												{TextRun: &docs.TextRun{Content: "Header 1"}},
											},
										},
									},
								},
							},
							{
								Content: []*docs.StructuralElement{
									{
										Paragraph: &docs.Paragraph{
											Elements: []*docs.ParagraphElement{
												{TextRun: &docs.TextRun{Content: "Header 2"}},
											},
										},
									},
								},
							},
						},
					},
					{
						TableCells: []*docs.TableCell{
							{
								Content: []*docs.StructuralElement{
									{
										Paragraph: &docs.Paragraph{
											Elements: []*docs.ParagraphElement{
												{TextRun: &docs.TextRun{Content: "Cell 1"}},
											},
										},
									},
								},
							},
							{
								Content: []*docs.StructuralElement{
									{
										Paragraph: &docs.Paragraph{
											Elements: []*docs.ParagraphElement{
												{TextRun: &docs.TextRun{Content: "Cell 2"}},
											},
										},
									},
								},
							},
						},
					},
				},
			},
			want: "| Header 1 | Header 2 |\n|---|---|\n| Cell 1 | Cell 2 |\n\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ConvertTable(tt.table)
			if got != tt.want {
				t.Errorf("ConvertTable() = %q, want %q", got, tt.want)
			}
		})
	}
}
