package renderer

import (
	"crypto/sha256"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// FileCache handles downloading and caching files from Notion
type FileCache struct {
	// basePath is the root content directory (e.g., "content")
	basePath string
	// httpClient for downloading files
	httpClient *http.Client
}

// NewFileCache creates a new file cache instance
func NewFileCache(basePath string) *FileCache {
	return &FileCache{
		basePath: basePath,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// CacheFile downloads a file from Notion and saves it to the article directory.
// Returns the relative path that should be used in markdown (e.g., "./image.jpg")
// This method assumes the caller has already determined the file should be cached.
func (fc *FileCache) CacheFile(notionURL, articlePath string) (string, error) {
	// Get the directory where the article will be saved
	articleDir := filepath.Dir(articlePath)
	fullArticleDir := filepath.Join(fc.basePath, articleDir)

	// Ensure the directory exists
	if err := os.MkdirAll(fullArticleDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create directory %s: %w", fullArticleDir, err)
	}

	// Generate a filename for the cached file
	filename, err := fc.generateFilename(notionURL)
	if err != nil {
		return "", fmt.Errorf("failed to generate filename: %w", err)
	}

	// Full path where the file will be saved
	localPath := filepath.Join(fullArticleDir, filename)

	// Check if file already exists
	if _, err := os.Stat(localPath); err == nil {
		// File already exists, return relative path
		return "./" + filename, nil
	}

	// Download the file
	if err := fc.downloadFile(notionURL, localPath); err != nil {
		return "", fmt.Errorf("failed to download file: %w", err)
	}

	// Return relative path for markdown
	return "./" + filename, nil
}

// generateFilename creates a unique filename based on the URL
func (fc *FileCache) generateFilename(notionURL string) (string, error) {
	// Extract file extension from URL
	ext := fc.extractExtension(notionURL)

	// Extract the file identifier (without signed parameters) for consistent caching
	fileId := fc.extractFileIdentifier(notionURL)

	// Create a hash of the file identifier for uniqueness using SHA-256
	hasher := sha256.New()
	hasher.Write([]byte(fileId))
	hash := fmt.Sprintf("%x", hasher.Sum(nil))

	// Use first 8 chars of hash + extension
	filename := hash[:8] + ext

	return filename, nil
}

// extractFileIdentifier extracts a stable identifier from the Notion file URL
// This removes signed parameters to ensure consistent caching
func (fc *FileCache) extractFileIdentifier(notionURL string) string {
	parsed, err := url.Parse(notionURL)
	if err != nil {
		// If URL parsing fails, use the full URL as fallback
		return notionURL
	}

	// For AWS S3 URLs (both old and new Notion formats), use the path without query parameters
	// Old format: https://s3.us-west-2.amazonaws.com/secure.notion-static.com/abc123/image.jpg?X-Amz-...
	// New format: https://prod-files-secure.s3.us-west-2.amazonaws.com/workspace-id/file-id/filename.pdf?X-Amz-...
	// We want the path part: /secure.notion-static.com/abc123/image.jpg or /workspace-id/file-id/filename.pdf
	if strings.Contains(parsed.Host, "amazonaws.com") {
		return parsed.Path
	}

	// For other Notion URLs, use host + path
	// Example: https://www.notion.so/workspace/file-id
	// We want: notion.so/workspace/file-id
	if strings.Contains(parsed.Host, "notion.so") {
		return parsed.Host + parsed.Path
	}

	// For other URLs, use full URL without query parameters as fallback
	return parsed.Host + parsed.Path
}

// extractExtension tries to extract file extension from URL
func (fc *FileCache) extractExtension(u string) string {
	// Remove query parameters

	parsed, err := url.Parse(u)
	if err != nil {
		return ".bin"
	}
	path := parsed.Path

	ext := filepath.Ext(path)
	if ext == "" {
		// Try to guess from URL patterns
		if strings.Contains(u, "image") {
			return ".jpg"
		} else if strings.Contains(u, "video") {
			return ".mp4"
		} else if strings.Contains(u, "pdf") {
			return ".pdf"
		}
		// Default to .bin if can't determine
		return ".bin"
	}

	return ext
}

// downloadFile downloads a file from URL and saves it to localPath
func (fc *FileCache) downloadFile(url, localPath string) error {
	resp, err := fc.httpClient.Get(url)
	if err != nil {
		return fmt.Errorf("failed to fetch URL %s: %w", url, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("HTTP %d when fetching %s", resp.StatusCode, url)
	}

	file, err := os.Create(localPath)
	if err != nil {
		return fmt.Errorf("failed to create file %s: %w", localPath, err)
	}
	defer file.Close()

	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return fmt.Errorf("failed to write file %s: %w", localPath, err)
	}

	return nil
}
