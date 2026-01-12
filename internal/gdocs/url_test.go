package gdocs

import (
	"testing"
)

func TestExtractDocumentID(t *testing.T) {
	tests := []struct {
		name    string
		url     string
		want    string
		wantErr bool
	}{
		{
			name:    "valid URL with /edit",
			url:     "https://docs.google.com/document/d/1abc123xyz/edit",
			want:    "1abc123xyz",
			wantErr: false,
		},
		{
			name:    "valid URL with /edit and query params",
			url:     "https://docs.google.com/document/d/1abc123xyz/edit?usp=sharing",
			want:    "1abc123xyz",
			wantErr: false,
		},
		{
			name:    "valid URL without /edit",
			url:     "https://docs.google.com/document/d/1abc123xyz",
			want:    "1abc123xyz",
			wantErr: false,
		},
		{
			name:    "valid URL with trailing slash",
			url:     "https://docs.google.com/document/d/1abc123xyz/",
			want:    "1abc123xyz",
			wantErr: false,
		},
		{
			name:    "valid URL with hyphens and underscores in ID",
			url:     "https://docs.google.com/document/d/1abc-123_xyz/edit",
			want:    "1abc-123_xyz",
			wantErr: false,
		},
		{
			name:    "invalid URL - not Google Docs",
			url:     "https://example.com/document/123",
			want:    "",
			wantErr: true,
		},
		{
			name:    "invalid URL - missing document ID",
			url:     "https://docs.google.com/document/",
			want:    "",
			wantErr: true,
		},
		{
			name:    "invalid URL - Google Sheets",
			url:     "https://docs.google.com/spreadsheets/d/1abc123xyz/edit",
			want:    "",
			wantErr: true,
		},
		{
			name:    "empty URL",
			url:     "",
			want:    "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ExtractDocumentID(tt.url)
			if (err != nil) != tt.wantErr {
				t.Errorf("ExtractDocumentID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ExtractDocumentID() = %v, want %v", got, tt.want)
			}
		})
	}
}
