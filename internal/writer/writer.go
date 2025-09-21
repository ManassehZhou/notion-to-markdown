package writer

import (
	"os"
	"path/filepath"
)

// Package writer contains a small helper to write generated files to disk.
// It is intentionally minimal so tests can replace it with a mock writer.

// Writer writes content to disk. Kept small so we can swap with a mock in tests.
type Writer struct{}

// New constructs a Writer instance.
func New() *Writer { return &Writer{} }

// WriteFile ensures the parent directory exists and writes content to filename.
// It returns any error from directory creation or file writing.
func (w *Writer) WriteFile(filename, content string) error {
	dir := filepath.Dir(filename)
	if dir != "" && dir != "." {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return err
		}
	}
	return os.WriteFile(filename, []byte(content), 0644)
}
