package renderer

import (
	"strings"
	"testing"
	"time"

	"github.com/jomei/notionapi"
)

func TestParseMetadata_SummaryAndCategories(t *testing.T) {
	renderer := New(nil, "test", nil)

	// Create a mock Notion page with Summary and Categories properties
	page := notionapi.Page{
		CreatedTime:    time.Date(2025, 1, 15, 10, 0, 0, 0, time.UTC),
		LastEditedTime: time.Date(2025, 1, 15, 12, 0, 0, 0, time.UTC),
		Properties: notionapi.Properties{
			"Title": &notionapi.TitleProperty{
				Title: []notionapi.RichText{
					{PlainText: "Test Article"},
				},
			},
			"Summary": &notionapi.RichTextProperty{
				RichText: []notionapi.RichText{
					{PlainText: "This is a test summary"},
				},
			},
			"Categories": &notionapi.MultiSelectProperty{
				MultiSelect: []notionapi.Option{
					{Name: "development"},
					{Name: "web-development"},
				},
			},
			"Tags": &notionapi.MultiSelectProperty{
				MultiSelect: []notionapi.Option{
					{Name: "hugo"},
					{Name: "tutorial"},
				},
			},
		},
	}

	meta := renderer.parseMetadata(page)

	// Verify all fields are parsed correctly
	if meta.Title != "Test Article" {
		t.Errorf("Expected title 'Test Article', got '%s'", meta.Title)
	}

	if meta.Summary != "This is a test summary" {
		t.Errorf("Expected summary 'This is a test summary', got '%s'", meta.Summary)
	}

	expectedCategories := []string{"development", "web-development"}
	if len(meta.Categories) != len(expectedCategories) {
		t.Errorf("Expected %d categories, got %d", len(expectedCategories), len(meta.Categories))
	}
	for i, expected := range expectedCategories {
		if i >= len(meta.Categories) || meta.Categories[i] != expected {
			t.Errorf("Expected category '%s' at index %d, got '%s'", expected, i, meta.Categories[i])
		}
	}

	expectedTags := []string{"hugo", "tutorial"}
	if len(meta.Tags) != len(expectedTags) {
		t.Errorf("Expected %d tags, got %d", len(expectedTags), len(meta.Tags))
	}
	for i, expected := range expectedTags {
		if i >= len(meta.Tags) || meta.Tags[i] != expected {
			t.Errorf("Expected tag '%s' at index %d, got '%s'", expected, i, meta.Tags[i])
		}
	}
}

func TestParseMetadata_AlternativeFieldNames(t *testing.T) {
	renderer := New(nil, "test", nil)

	// Test alternative field names (Category vs Categories, Description vs Summary)
	page := notionapi.Page{
		Properties: notionapi.Properties{
			"Title": &notionapi.TitleProperty{
				Title: []notionapi.RichText{
					{PlainText: "Alternative Names Test"},
				},
			},
			"Description": &notionapi.RichTextProperty{
				RichText: []notionapi.RichText{
					{PlainText: "Alternative description field"},
				},
			},
			"Category": &notionapi.MultiSelectProperty{
				MultiSelect: []notionapi.Option{
					{Name: "single-category"},
				},
			},
		},
	}

	meta := renderer.parseMetadata(page)

	if meta.Summary != "Alternative description field" {
		t.Errorf("Expected summary from 'Description' field, got '%s'", meta.Summary)
	}

	if len(meta.Categories) != 1 || meta.Categories[0] != "single-category" {
		t.Errorf("Expected single category 'single-category', got %v", meta.Categories)
	}
}

func TestBuildFrontMatter_WithSummaryAndCategories(t *testing.T) {
	renderer := New(nil, "test", nil)

	meta := metadata{
		Title:      "Test Article",
		Slug:       "test-article",
		Date:       "2025-01-15T10:00:00Z07:00",
		LastMod:    "2025-01-15T12:00:00Z07:00",
		Tags:       []string{"hugo", "tutorial"},
		Categories: []string{"development", "web-development"},
		Summary:    "This is a test summary for SEO",
		Draft:      false,
		Type:       "posts",
	}

	frontMatter, err := renderer.buildFrontMatter(meta)
	if err != nil {
		t.Fatalf("Unexpected error building front matter: %v", err)
	}

	// Check that front matter contains our new fields
	expectedSubstrings := []string{
		"title: Test Article",
		"slug: test-article",
		"categories:",
		"- development",
		"- web-development",
		"summary: This is a test summary for SEO",
		"tags:",
		"- hugo",
		"- tutorial",
		"type: posts",
	}

	for _, expected := range expectedSubstrings {
		if !strings.Contains(frontMatter, expected) {
			t.Errorf("Expected front matter to contain '%s', but it didn't.\nFull front matter:\n%s", expected, frontMatter)
		}
	}

	// Test that draft: false is omitted (due to omitempty) when false
	if strings.Contains(frontMatter, "draft: false") {
		t.Error("Expected 'draft: false' to be omitted due to omitempty YAML tag")
	}

	// Verify it starts and ends with proper YAML delimiters
	if !strings.HasPrefix(frontMatter, "---\n") {
		t.Error("Front matter should start with '---\\n'")
	}
	if !strings.HasSuffix(frontMatter, "---\n\n") {
		t.Error("Front matter should end with '---\\n\\n'")
	}
}

func TestParseMetadata_DatePriorityLogic(t *testing.T) {
	renderer := New(nil, "test", nil)

	// Test 1: Page with custom Date property should override CreatedTime
	pageWithCustomDate := notionapi.Page{
		CreatedTime:    time.Date(2025, 1, 15, 10, 0, 0, 0, time.UTC),
		LastEditedTime: time.Date(2025, 1, 15, 12, 0, 0, 0, time.UTC),
		Properties: notionapi.Properties{
			"Title": &notionapi.TitleProperty{
				Title: []notionapi.RichText{
					{PlainText: "Custom Date Test"},
				},
			},
			"Date": &notionapi.DateProperty{
				Date: &notionapi.DateObject{
					Start: (*notionapi.Date)(&[]time.Time{time.Date(2025, 2, 20, 0, 0, 0, 0, time.UTC)}[0]),
				},
			},
		},
	}

	meta := renderer.parseMetadata(pageWithCustomDate)
	expectedCustomDate := "2025-02-20T00:00:00Z"
	if meta.Date != expectedCustomDate {
		t.Errorf("Expected custom date '%s', got '%s'", expectedCustomDate, meta.Date)
	}

	// Test 2: Page without Date property should use CreatedTime
	pageWithoutCustomDate := notionapi.Page{
		CreatedTime:    time.Date(2025, 1, 15, 10, 0, 0, 0, time.UTC),
		LastEditedTime: time.Date(2025, 1, 15, 12, 0, 0, 0, time.UTC),
		Properties: notionapi.Properties{
			"Title": &notionapi.TitleProperty{
				Title: []notionapi.RichText{
					{PlainText: "Default Date Test"},
				},
			},
		},
	}

	meta2 := renderer.parseMetadata(pageWithoutCustomDate)
	expectedDefaultDate := "2025-01-15T10:00:00Z"
	if meta2.Date != expectedDefaultDate {
		t.Errorf("Expected fallback date '%s', got '%s'", expectedDefaultDate, meta2.Date)
	}

	// Test 3: Verify LastMod always uses LastEditedTime
	expectedLastMod := "2025-01-15T12:00:00Z"
	if meta.LastMod != expectedLastMod {
		t.Errorf("Expected lastmod '%s', got '%s'", expectedLastMod, meta.LastMod)
	}
	if meta2.LastMod != expectedLastMod {
		t.Errorf("Expected lastmod '%s', got '%s'", expectedLastMod, meta2.LastMod)
	}
}
