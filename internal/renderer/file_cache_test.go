package renderer

import (
	"path/filepath"
	"testing"
)

func TestFileCache_CacheFile(t *testing.T) {
	// Create a temporary directory for testing
	tempDir := t.TempDir()

	// Initialize file cache
	fc := NewFileCache(tempDir)

	// Test filename generation for consistency
	notionURL := "https://s3.us-west-2.amazonaws.com/secure.notion-static.com/test-image.jpg?X-Amz-Algorithm=AWS4-HMAC-SHA256"
	filename, err := fc.generateFilename(notionURL)
	if err != nil {
		t.Fatalf("Expected no error generating filename, got %v", err)
	}
	if filename == "" {
		t.Error("Expected non-empty filename")
	}
	if filepath.Ext(filename) == "" {
		t.Error("Expected filename to have an extension")
	}
}

func TestFileCache_GenerateFilename(t *testing.T) {
	fc := NewFileCache("test")

	testCases := []struct {
		url      string
		expected string // Just check the extension
	}{
		{"https://example.com/image.jpg", ".jpg"},
		{"https://example.com/document.pdf", ".pdf"},
		{"https://example.com/video.mp4", ".mp4"},
		{"https://example.com/file", ".bin"}, // Default when no extension
	}

	for _, tc := range testCases {
		filename, err := fc.generateFilename(tc.url)
		if err != nil {
			t.Errorf("Unexpected error for URL %s: %v", tc.url, err)
			continue
		}

		ext := filepath.Ext(filename)
		if ext != tc.expected {
			t.Errorf("For URL %s, expected extension %s, got %s", tc.url, tc.expected, ext)
		}

		// Check that filename is not empty and has reasonable length
		if len(filename) < 4 { // At least hash chars + extension
			t.Errorf("Filename too short: %s", filename)
		}
	}
}

func TestFileCache_ExtractFileIdentifier(t *testing.T) {
	fc := NewFileCache("test")

	testCases := []struct {
		name       string
		url        string
		expectedID string
	}{
		{
			name:       "AWS S3 signed URL (old format)",
			url:        "https://s3.us-west-2.amazonaws.com/secure.notion-static.com/abc123/image.jpg?X-Amz-Algorithm=AWS4-HMAC-SHA256&X-Amz-Credential=test",
			expectedID: "/secure.notion-static.com/abc123/image.jpg",
		},
		{
			name:       "AWS S3 signed URL (new format)",
			url:        "https://prod-files-secure.s3.us-west-2.amazonaws.com/d9d52f73-bbd3-47db-96fe-27b0615621ac/20ecd1ff-e4de-4471-b86f-6f14ec891fc0/%E6%B5%8B%E8%AF%95%E5%B5%8C%E5%85%A5PDF.pdf?X-Amz-Algorithm=AWS4-HMAC-SHA256&X-Amz-Date=20250920T045407Z",
			expectedID: "/d9d52f73-bbd3-47db-96fe-27b0615621ac/20ecd1ff-e4de-4471-b86f-6f14ec891fc0/测试嵌入PDF.pdf",
		},
		{
			name:       "Another AWS S3 signed URL with different params",
			url:        "https://s3.us-west-2.amazonaws.com/secure.notion-static.com/abc123/image.jpg?X-Amz-Algorithm=AWS4-HMAC-SHA256&X-Amz-Date=20230101T000000Z",
			expectedID: "/secure.notion-static.com/abc123/image.jpg",
		},
		{
			name:       "Notion.so URL",
			url:        "https://www.notion.so/workspace/file-id?v=abc123",
			expectedID: "www.notion.so/workspace/file-id",
		},
		{
			name:       "External URL",
			url:        "https://example.com/image.jpg?param=value",
			expectedID: "example.com/image.jpg",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := fc.extractFileIdentifier(tc.url)
			if result != tc.expectedID {
				t.Errorf("Expected %s, got %s", tc.expectedID, result)
			}
		})
	}
}

func TestFileCache_ConsistentCaching(t *testing.T) {
	fc := NewFileCache("test")

	// Test that the same file with different signed parameters generates the same filename
	url1 := "https://s3.us-west-2.amazonaws.com/secure.notion-static.com/abc123/image.jpg?X-Amz-Algorithm=AWS4-HMAC-SHA256&X-Amz-Date=20230101T000000Z"
	url2 := "https://s3.us-west-2.amazonaws.com/secure.notion-static.com/abc123/image.jpg?X-Amz-Algorithm=AWS4-HMAC-SHA256&X-Amz-Date=20230102T000000Z"

	filename1, err1 := fc.generateFilename(url1)
	filename2, err2 := fc.generateFilename(url2)

	if err1 != nil || err2 != nil {
		t.Fatalf("Unexpected errors: %v, %v", err1, err2)
	}

	if filename1 != filename2 {
		t.Errorf("Expected same filename for same file with different signatures, got %s and %s", filename1, filename2)
	}

	// Test new format consistency
	url3 := "https://prod-files-secure.s3.us-west-2.amazonaws.com/d9d52f73-bbd3-47db-96fe-27b0615621ac/20ecd1ff-e4de-4471-b86f-6f14ec891fc0/test.pdf?X-Amz-Date=20250920T045407Z&X-Amz-Signature=abc123"
	url4 := "https://prod-files-secure.s3.us-west-2.amazonaws.com/d9d52f73-bbd3-47db-96fe-27b0615621ac/20ecd1ff-e4de-4471-b86f-6f14ec891fc0/test.pdf?X-Amz-Date=20250921T045407Z&X-Amz-Signature=def456"

	filename3, err3 := fc.generateFilename(url3)
	filename4, err4 := fc.generateFilename(url4)

	if err3 != nil || err4 != nil {
		t.Fatalf("Unexpected errors for new format: %v, %v", err3, err4)
	}

	if filename3 != filename4 {
		t.Errorf("Expected same filename for same file (new format) with different signatures, got %s and %s", filename3, filename4)
	}
}
