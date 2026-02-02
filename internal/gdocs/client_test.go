package gdocs

import (
	"testing"

	"google.golang.org/api/docs/v1"
)

func TestFindTab(t *testing.T) {
	// Create a mock document with nested tabs
	doc := &docs.Document{
		Title: "Test Document",
		Tabs: []*docs.Tab{
			{
				TabProperties: &docs.TabProperties{
					TabId: "t.tab1",
					Title: "Tab 1",
				},
				DocumentTab: &docs.DocumentTab{
					Body: &docs.Body{},
				},
				ChildTabs: []*docs.Tab{
					{
						TabProperties: &docs.TabProperties{
							TabId: "t.tab1child1",
							Title: "Tab 1 Child 1",
						},
						DocumentTab: &docs.DocumentTab{
							Body: &docs.Body{},
						},
					},
				},
			},
			{
				TabProperties: &docs.TabProperties{
					TabId: "t.tab2",
					Title: "Tab 2",
				},
				DocumentTab: &docs.DocumentTab{
					Body: &docs.Body{},
				},
			},
		},
	}

	tests := []struct {
		name      string
		tabID     string
		wantTitle string
		wantNil   bool
	}{
		{
			name:      "find top-level tab",
			tabID:     "t.tab1",
			wantTitle: "Tab 1",
			wantNil:   false,
		},
		{
			name:      "find second top-level tab",
			tabID:     "t.tab2",
			wantTitle: "Tab 2",
			wantNil:   false,
		},
		{
			name:      "find nested child tab",
			tabID:     "t.tab1child1",
			wantTitle: "Tab 1 Child 1",
			wantNil:   false,
		},
		{
			name:    "tab not found",
			tabID:   "t.nonexistent",
			wantNil: true,
		},
		{
			name:    "empty tab ID",
			tabID:   "",
			wantNil: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := FindTab(doc, tt.tabID)
			if tt.wantNil {
				if got != nil {
					t.Errorf("FindTab() = %v, want nil", got)
				}
				return
			}
			if got == nil {
				t.Errorf("FindTab() = nil, want tab with title %q", tt.wantTitle)
				return
			}
			if got.TabProperties.Title != tt.wantTitle {
				t.Errorf("FindTab() title = %q, want %q", got.TabProperties.Title, tt.wantTitle)
			}
		})
	}
}

func TestFindTab_NilTabs(t *testing.T) {
	doc := &docs.Document{
		Title: "Test Document",
		Tabs:  nil,
	}

	got := FindTab(doc, "t.any")
	if got != nil {
		t.Errorf("FindTab() with nil tabs = %v, want nil", got)
	}
}

func TestGetFirstTab(t *testing.T) {
	tests := []struct {
		name      string
		doc       *docs.Document
		wantTitle string
		wantNil   bool
	}{
		{
			name: "document with tabs",
			doc: &docs.Document{
				Tabs: []*docs.Tab{
					{
						TabProperties: &docs.TabProperties{
							TabId: "t.first",
							Title: "First Tab",
						},
					},
				},
			},
			wantTitle: "First Tab",
			wantNil:   false,
		},
		{
			name: "document with no tabs",
			doc: &docs.Document{
				Tabs: []*docs.Tab{},
			},
			wantNil: true,
		},
		{
			name: "document with nil tabs",
			doc: &docs.Document{
				Tabs: nil,
			},
			wantNil: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetFirstTab(tt.doc)
			if tt.wantNil {
				if got != nil {
					t.Errorf("GetFirstTab() = %v, want nil", got)
				}
				return
			}
			if got == nil {
				t.Errorf("GetFirstTab() = nil, want tab with title %q", tt.wantTitle)
				return
			}
			if got.TabProperties.Title != tt.wantTitle {
				t.Errorf("GetFirstTab() title = %q, want %q", got.TabProperties.Title, tt.wantTitle)
			}
		})
	}
}
