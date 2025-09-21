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

	if summary, ok := meta.Properties["Summary"].(string); !ok || summary != "This is a test summary" {
		t.Errorf("Expected summary 'This is a test summary', got '%v'", meta.Properties["Summary"])
	}

	expectedCategories := []string{"development", "web-development"}
	if categories, ok := meta.Properties["Categories"].([]string); !ok {
		t.Errorf("Expected Categories to be []string, got %T", meta.Properties["Categories"])
	} else if len(categories) != len(expectedCategories) {
		t.Errorf("Expected %d categories, got %d", len(expectedCategories), len(categories))
	} else {
		for i, expected := range expectedCategories {
			if i >= len(categories) || categories[i] != expected {
				t.Errorf("Expected category '%s' at index %d, got '%s'", expected, i, categories[i])
			}
		}
	}

	expectedTags := []string{"hugo", "tutorial"}
	if tags, ok := meta.Properties["Tags"].([]string); !ok {
		t.Errorf("Expected Tags to be []string, got %T", meta.Properties["Tags"])
	} else if len(tags) != len(expectedTags) {
		t.Errorf("Expected %d tags, got %d", len(expectedTags), len(tags))
	} else {
		for i, expected := range expectedTags {
			if i >= len(tags) || tags[i] != expected {
				t.Errorf("Expected tag '%s' at index %d, got '%s'", expected, i, tags[i])
			}
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

	if description, ok := meta.Properties["Description"].(string); !ok || description != "Alternative description field" {
		t.Errorf("Expected summary from 'Description' field, got '%v'", meta.Properties["Description"])
	}

	if categories, ok := meta.Properties["Category"].([]string); !ok || len(categories) != 1 || categories[0] != "single-category" {
		t.Errorf("Expected single category 'single-category', got %v", meta.Properties["Category"])
	}
}

func TestBuildFrontMatter_WithSummaryAndCategories(t *testing.T) {
	renderer := New(nil, "test", nil)

	meta := metadata{
		Title: "Test Article",
		Slug:  "test-article",
		Properties: map[string]interface{}{
			"title":      "Test Article",
			"slug":       "test-article",
			"date":       "2025-01-15T10:00:00Z07:00",
			"lastmod":    "2025-01-15T12:00:00Z07:00",
			"tags":       []string{"hugo", "tutorial"},
			"categories": []string{"development", "web-development"},
			"summary":    "This is a test summary for SEO",
			"type":       "posts",
		},
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
	if date, ok := meta.Properties["date"].(string); !ok || date != expectedCustomDate {
		t.Errorf("Expected custom date '%s', got '%v'", expectedCustomDate, meta.Properties["date"])
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
	if date, ok := meta2.Properties["date"].(string); !ok || date != expectedDefaultDate {
		t.Errorf("Expected fallback date '%s', got '%v'", expectedDefaultDate, meta2.Properties["date"])
	}

	// Test 3: Verify LastMod always uses LastEditedTime
	expectedLastMod := "2025-01-15T12:00:00Z"
	if lastmod, ok := meta.Properties["lastmod"].(string); !ok || lastmod != expectedLastMod {
		t.Errorf("Expected lastmod '%s', got '%v'", expectedLastMod, meta.Properties["lastmod"])
	}
	if lastmod, ok := meta2.Properties["lastmod"].(string); !ok || lastmod != expectedLastMod {
		t.Errorf("Expected lastmod '%s', got '%v'", expectedLastMod, meta2.Properties["lastmod"])
	}
}

func TestParseMetadata_TypeHandling(t *testing.T) {
	renderer := New(nil, "test", nil)
	now := time.Date(2025, 1, 15, 10, 0, 0, 0, time.UTC)

	// Test 1: Normal type
	page1 := notionapi.Page{
		CreatedTime:    now,
		LastEditedTime: now,
		Properties: notionapi.Properties{
			"Title": &notionapi.TitleProperty{
				Title: []notionapi.RichText{
					{PlainText: "Test Article"},
				},
			},
			"Type": &notionapi.SelectProperty{
				Select: notionapi.Option{Name: "blog"},
			},
		},
	}

	meta1 := renderer.parseMetadata(page1)
	if typeVal, ok := meta1.Properties["type"].(string); !ok || typeVal != "blog" {
		t.Errorf("Expected frontmatter type 'blog', got '%v'", meta1.Properties["type"])
	}
	if meta1.pathType != "blog" {
		t.Errorf("Expected path type 'blog', got '%s'", meta1.pathType)
	}

	// Test 2: pages:friends format
	page2 := notionapi.Page{
		CreatedTime:    now,
		LastEditedTime: now,
		Properties: notionapi.Properties{
			"Title": &notionapi.TitleProperty{
				Title: []notionapi.RichText{
					{PlainText: "My Friends"},
				},
			},
			"Type": &notionapi.SelectProperty{
				Select: notionapi.Option{Name: "pages:friends"},
			},
		},
	}

	meta2 := renderer.parseMetadata(page2)
	if typeVal, ok := meta2.Properties["type"].(string); !ok || typeVal != "friends" {
		t.Errorf("Expected frontmatter type 'friends', got '%v'", meta2.Properties["type"])
	}
	if meta2.pathType != "pages" {
		t.Errorf("Expected path type 'pages', got '%s'", meta2.pathType)
	}

	// Test 3: post -> posts normalization
	page3 := notionapi.Page{
		CreatedTime:    now,
		LastEditedTime: now,
		Properties: notionapi.Properties{
			"Title": &notionapi.TitleProperty{
				Title: []notionapi.RichText{
					{PlainText: "Blog Post"},
				},
			},
			"Type": &notionapi.SelectProperty{
				Select: notionapi.Option{Name: "post"},
			},
		},
	}

	meta3 := renderer.parseMetadata(page3)
	if typeVal, ok := meta3.Properties["type"].(string); !ok || typeVal != "posts" {
		t.Errorf("Expected frontmatter type 'posts', got '%v'", meta3.Properties["type"])
	}
	if meta3.pathType != "posts" {
		t.Errorf("Expected path type 'posts', got '%s'", meta3.pathType)
	}

	// Test 4: docs:guides format (non-pages category)
	page4 := notionapi.Page{
		CreatedTime:    now,
		LastEditedTime: now,
		Properties: notionapi.Properties{
			"Title": &notionapi.TitleProperty{
				Title: []notionapi.RichText{
					{PlainText: "User Guide"},
				},
			},
			"Type": &notionapi.SelectProperty{
				Select: notionapi.Option{Name: "docs:guides"},
			},
		},
	}

	meta4 := renderer.parseMetadata(page4)
	if typeVal, ok := meta4.Properties["type"].(string); !ok || typeVal != "docs:guides" {
		t.Errorf("Expected frontmatter type 'docs:guides', got '%v'", meta4.Properties["type"])
	}
	if meta4.pathType != "docs:guides" {
		t.Errorf("Expected path type 'docs:guides', got '%s'", meta4.pathType)
	}
}

func TestGetPagePath_TypeHandling(t *testing.T) {
	renderer := New(nil, "test", nil)
	now := time.Date(2025, 1, 15, 10, 0, 0, 0, time.UTC)

	// Test 1: Normal type
	page1 := notionapi.Page{
		CreatedTime:    now,
		LastEditedTime: now,
		Properties: notionapi.Properties{
			"Title": &notionapi.TitleProperty{
				Title: []notionapi.RichText{
					{PlainText: "Test Article"},
				},
			},
			"Type": &notionapi.SelectProperty{
				Select: notionapi.Option{Name: "blog"},
			},
		},
	}

	path1 := renderer.GetPagePath(page1)
	expected1 := "/blog/test-article/"
	if path1 != expected1 {
		t.Errorf("Expected path '%s', got '%s'", expected1, path1)
	}

	// Test 2: pages:friends format should use pages logic
	page2 := notionapi.Page{
		CreatedTime:    now,
		LastEditedTime: now,
		Properties: notionapi.Properties{
			"Title": &notionapi.TitleProperty{
				Title: []notionapi.RichText{
					{PlainText: "My Friends"},
				},
			},
			"Type": &notionapi.SelectProperty{
				Select: notionapi.Option{Name: "pages:friends"},
			},
		},
	}

	path2 := renderer.GetPagePath(page2)
	expected2 := "/my-friends/"
	if path2 != expected2 {
		t.Errorf("Expected path '%s', got '%s'", expected2, path2)
	}

	// Test 3: Empty type defaults to posts
	page3 := notionapi.Page{
		CreatedTime:    now,
		LastEditedTime: now,
		Properties: notionapi.Properties{
			"Title": &notionapi.TitleProperty{
				Title: []notionapi.RichText{
					{PlainText: "Default Post"},
				},
			},
		},
	}

	path3 := renderer.GetPagePath(page3)
	expected3 := "/posts/default-post/"
	if path3 != expected3 {
		t.Errorf("Expected path '%s', got '%s'", expected3, path3)
	}
}

func TestDynamicProperties(t *testing.T) {
	renderer := New(nil, "test", nil)
	now := time.Date(2025, 1, 15, 10, 0, 0, 0, time.UTC)

	// Test page with various custom properties
	page := notionapi.Page{
		CreatedTime:    now,
		LastEditedTime: now,
		Properties: notionapi.Properties{
			"Title": &notionapi.TitleProperty{
				Title: []notionapi.RichText{
					{PlainText: "Dynamic Properties Test"},
				},
			},
			"Author": &notionapi.RichTextProperty{
				RichText: []notionapi.RichText{
					{PlainText: "John Doe"},
				},
			},
			"Priority": &notionapi.SelectProperty{
				Select: notionapi.Option{Name: "High"},
			},
			"Labels": &notionapi.MultiSelectProperty{
				MultiSelect: []notionapi.Option{
					{Name: "important"},
					{Name: "urgent"},
				},
			},
			"Project Status": &notionapi.StatusProperty{
				Status: notionapi.Option{Name: "In Progress"},
			},
			"Published Date": &notionapi.DateProperty{
				Date: &notionapi.DateObject{
					Start: (*notionapi.Date)(&now),
				},
			},
			"Custom Field": &notionapi.RichTextProperty{
				RichText: []notionapi.RichText{
					{PlainText: "Custom Value"},
				},
			},
		},
	}

	meta := renderer.parseMetadata(page)

	// Test all custom properties are preserved
	if author, ok := meta.Properties["Author"].(string); !ok || author != "John Doe" {
		t.Errorf("Expected Author 'John Doe', got %v", meta.Properties["Author"])
	}

	if priority, ok := meta.Properties["Priority"].(string); !ok || priority != "High" {
		t.Errorf("Expected Priority 'High', got %v", meta.Properties["Priority"])
	}

	if labels, ok := meta.Properties["Labels"].([]string); !ok {
		t.Errorf("Expected Labels to be []string, got %T", meta.Properties["Labels"])
	} else if len(labels) != 2 || labels[0] != "important" || labels[1] != "urgent" {
		t.Errorf("Expected Labels ['important', 'urgent'], got %v", labels)
	}

	if status, ok := meta.Properties["Project Status"].(string); !ok || status != "In Progress" {
		t.Errorf("Expected Project Status 'In Progress', got %v", meta.Properties["Project Status"])
	}

	expectedDate := now.Format("2006-01-02T15:04:05Z07:00")
	if publishedDate, ok := meta.Properties["Published Date"].(string); !ok || publishedDate != expectedDate {
		t.Errorf("Expected Published Date '%s', got %v", expectedDate, meta.Properties["Published Date"])
	}

	if customField, ok := meta.Properties["Custom Field"].(string); !ok || customField != "Custom Value" {
		t.Errorf("Expected Custom Field 'Custom Value', got %v", meta.Properties["Custom Field"])
	}

	// Test frontmatter generation includes all properties
	frontMatter, err := renderer.buildFrontMatter(meta)
	if err != nil {
		t.Fatalf("Unexpected error building front matter: %v", err)
	}

	// Check that all custom properties appear in frontmatter
	expectedInFrontMatter := []string{
		"Author: John Doe",
		"Priority: High",
		"Labels:",
		"- important",
		"- urgent",
		"Project Status: In Progress",
		"Custom Field: Custom Value",
	}

	for _, expected := range expectedInFrontMatter {
		if !strings.Contains(frontMatter, expected) {
			t.Errorf("Expected front matter to contain '%s', but it didn't.\nFull front matter:\n%s", expected, frontMatter)
		}
	}
}
