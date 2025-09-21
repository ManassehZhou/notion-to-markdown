// Package renderer converts Notion pages and blocks into Markdown documents
// suitable for static site generators like Hugo, Hexo, Jekyll, etc. The public
// Renderer type exposes a simple RenderPage method which returns a filename and
// the full file content including YAML front matter.
package renderer

import (
	"path/filepath"
	"strings"
	"time"

	"github.com/jomei/notionapi"
	"gopkg.in/yaml.v3"
)

// Renderer converts Notion pages/blocks into Markdown + frontmatter.
// It is intentionally small and testable: network I/O is provided via an
// injected getChildren callback.
type Renderer struct {
	// resolve maps a Notion page ID to a site-relative path (e.g. "/posts/slug/").
	// If nil, no resolution by ID will occur.
	resolve func(string) string

	// fileCache handles downloading and caching files from Notion
	fileCache *FileCache

	// config controls how non-standard markdown elements are rendered
	config *RenderConfig
}

// New constructs a Renderer with link resolver, file caching and custom config.
func New(resolve func(string) string, basePath string, config *RenderConfig) *Renderer {
	return &Renderer{
		resolve:   resolve,
		fileCache: NewFileCache(basePath),
		config:    config,
	}
}

// RenderPage converts a Notion page and its provided top-level blocks into a
// filename and file content (YAML front matter + Markdown body). The
// getChildren callback is used to lazily fetch block children; this keeps the
// method side-effect free for testing when a mock callback is provided.
func (r *Renderer) RenderPage(page notionapi.Page, blocks []notionapi.Block, getChildren func(notionapi.BlockID) ([]notionapi.Block, error), resolve func(string) string) (string, string, error) {
	meta := r.parseMetadata(page)
	filename := r.buildFilename(meta)

	// render body using recursive helper
	// prefer resolver passed to RenderPage, otherwise use renderer's resolver
	if resolve == nil {
		resolve = r.resolve
	}
	body, err := r.renderBlocksRecursive(blocks, getChildren, resolve, filename)
	if err != nil {
		return "", "", err
	}

	fm, err := r.buildFrontMatter(meta)
	if err != nil {
		return "", "", err
	}
	return filename, fm + body, nil
}

// metadata gathers the common properties used in frontmatter and filename logic.
type metadata struct {
	// Core fields needed for functionality
	Title    string `yaml:"title"`
	Slug     string `yaml:"slug,omitempty"`
	pathType string `yaml:"-"` // Used internally for path generation logic

	// All properties including user-defined ones
	Properties map[string]interface{} `yaml:",inline"`
}

func (r *Renderer) parseMetadata(page notionapi.Page) metadata {
	m := metadata{
		Title:      "untitled",
		Properties: make(map[string]interface{}),
	}

	// Set default timestamps from Notion page metadata
	if !page.CreatedTime.IsZero() {
		m.Properties["date"] = page.CreatedTime.Format("2006-01-02T15:04:05Z07:00")
	}
	if !page.LastEditedTime.IsZero() {
		m.Properties["lastmod"] = page.LastEditedTime.Format("2006-01-02T15:04:05Z07:00")
	}

	// Parse all properties from the Notion page
	for k, prop := range page.Properties {
		lowerKey := strings.ToLower(k)

		// Handle special properties that affect internal logic
		switch lowerKey {
		case "title", "name":
			if tp, ok := prop.(*notionapi.TitleProperty); ok && len(tp.Title) > 0 {
				m.Title = tp.Title[0].PlainText
				m.Properties["title"] = m.Title
			}
		case "slug":
			value := extractPropertyValue(prop)
			if str, ok := value.(string); ok && str != "" {
				m.Slug = str
				m.Properties["slug"] = str
			}
		case "date":
			if dp, ok := prop.(*notionapi.DateProperty); ok && dp.Date != nil && dp.Date.Start != nil {
				dateStr := time.Time(*dp.Date.Start).Format("2006-01-02T15:04:05Z07:00")
				m.Properties["date"] = dateStr // Override default
			}
		case "type":
			value := extractPropertyValue(prop)
			if str, ok := value.(string); ok && str != "" {
				originalType := str
				m.Properties["type"] = originalType

				// Parse type for path generation logic
				if strings.Contains(originalType, ":") {
					parts := strings.SplitN(originalType, ":", 2)
					if len(parts) == 2 {
						pathCategory := strings.ToLower(strings.TrimSpace(parts[0]))
						typeValue := strings.TrimSpace(parts[1])

						if pathCategory == "pages" {
							m.pathType = "pages"
							m.Properties["type"] = typeValue // frontmatter gets "friends"
						} else {
							m.pathType = strings.ToLower(originalType)
						}
					} else {
						m.pathType = strings.ToLower(originalType)
					}
				} else {
					lowerType := strings.ToLower(originalType)
					if lowerType == "post" {
						m.pathType = "posts"
						m.Properties["type"] = "posts" // Normalize in frontmatter too
					} else {
						m.pathType = lowerType
					}
				}
			}
		case "status":
			// Handle status specially to set draft flag
			if sp, ok := prop.(*notionapi.StatusProperty); ok {
				statusName := sp.Status.Name
				m.Properties["status"] = statusName
				if strings.ToLower(statusName) == "draft" {
					m.Properties["draft"] = true
				}
				// Note: We don't set draft: false to allow omitempty behavior
			}
		default:
			// Handle all other properties dynamically
			value := extractPropertyValue(prop)
			if value != nil {
				m.Properties[k] = value
			}
		}
	}

	// Set defaults
	if m.Slug == "" {
		m.Slug = m.Title
	}
	m.Slug = slugify(m.Slug)

	// Set default pathType if not set
	if m.pathType == "" {
		if typeVal, exists := m.Properties["type"]; exists {
			if str, ok := typeVal.(string); ok {
				m.pathType = str
			}
		}
	}

	return m
}

// extractPropertyValue extracts the value from various Notion property types
func extractPropertyValue(prop notionapi.Property) interface{} {
	switch v := prop.(type) {
	case *notionapi.TitleProperty:
		if len(v.Title) > 0 {
			return v.Title[0].PlainText
		}
	case *notionapi.RichTextProperty:
		if len(v.RichText) > 0 {
			return v.RichText[0].PlainText
		}
	case *notionapi.DateProperty:
		if v.Date != nil && v.Date.Start != nil {
			return time.Time(*v.Date.Start).Format("2006-01-02T15:04:05Z07:00")
		}
	case *notionapi.SelectProperty:
		return v.Select.Name
	case *notionapi.MultiSelectProperty:
		var values []string
		for _, sel := range v.MultiSelect {
			values = append(values, sel.Name)
		}
		return values
	case *notionapi.StatusProperty:
		return v.Status.Name
	}
	return nil
}

// GetPageSlug is a small helper used by callers that need a page's slug
// without rendering the entire page. It mirrors the logic used by parseMetadata
// and returns the final slugified value.
func (r *Renderer) GetPageSlug(page notionapi.Page) string {
	m := r.parseMetadata(page)
	return m.Slug
}

// GetPagePath returns the Hugo site-relative path for a page (e.g. "/posts/slug/")
// without rendering the entire page. This is used for building the resolver map.
func (r *Renderer) GetPagePath(page notionapi.Page) string {
	m := r.parseMetadata(page)
	safeType := slugify(m.pathType)

	// default posts
	if safeType == "" {
		return "/posts/" + m.Slug + "/"
	}
	if safeType == "pages" {
		return "/" + m.Slug + "/"
	}
	return "/" + safeType + "/" + m.Slug + "/"
}

func (r *Renderer) buildFilename(m metadata) string {
	safeType := slugify(m.pathType)
	// default posts
	if safeType == "" {
		return filepath.ToSlash(filepath.Join("posts", m.Slug, "index.md"))
	}
	if safeType == "pages" {
		return filepath.ToSlash(filepath.Join(m.Slug, "index.md"))
	}
	return filepath.ToSlash(filepath.Join(safeType, m.Slug, "index.md"))
}

func (r *Renderer) buildFrontMatter(m metadata) (string, error) {
	// Use the Properties map directly for YAML marshaling
	out, err := yaml.Marshal(m.Properties)
	if err != nil {
		// Fallback to minimal frontmatter on error
		return "", err
	}
	return "---\n" + string(out) + "---\n\n", nil
}

// renderBlocksRecursive renders top-level blocks and recursively fetches children
// via getChildren. It returns the combined markdown body.
func (r *Renderer) renderBlocksRecursive(blocks []notionapi.Block, getChildren func(notionapi.BlockID) ([]notionapi.Block, error), resolve func(string) string, articlePath string) (string, error) {
	// helper to detect ID/HasChildren
	getBlockIDAndHasChildren := func(block notionapi.Block) (notionapi.BlockID, bool) {
		switch b := block.(type) {
		case *notionapi.ParagraphBlock:
			return notionapi.BlockID(b.ID), b.HasChildren
		case *notionapi.Heading1Block:
			return notionapi.BlockID(b.ID), b.HasChildren
		case *notionapi.Heading2Block:
			return notionapi.BlockID(b.ID), b.HasChildren
		case *notionapi.Heading3Block:
			return notionapi.BlockID(b.ID), b.HasChildren
		case *notionapi.BulletedListItemBlock:
			return notionapi.BlockID(b.ID), b.HasChildren
		case *notionapi.NumberedListItemBlock:
			return notionapi.BlockID(b.ID), b.HasChildren
		case *notionapi.ToDoBlock:
			return notionapi.BlockID(b.ID), b.HasChildren
		case *notionapi.ToggleBlock:
			return notionapi.BlockID(b.ID), b.HasChildren
		case *notionapi.EquationBlock:
			return notionapi.BlockID(b.ID), b.HasChildren
		case *notionapi.CodeBlock:
			return notionapi.BlockID(b.ID), b.HasChildren
		case *notionapi.QuoteBlock:
			return notionapi.BlockID(b.ID), b.HasChildren
		case *notionapi.CalloutBlock:
			return notionapi.BlockID(b.ID), b.HasChildren
		case *notionapi.DividerBlock:
			return notionapi.BlockID(b.ID), b.HasChildren
		case *notionapi.ImageBlock:
			return notionapi.BlockID(b.ID), b.HasChildren
		case *notionapi.BookmarkBlock:
			return notionapi.BlockID(b.ID), b.HasChildren
		case *notionapi.EmbedBlock:
			return notionapi.BlockID(b.ID), b.HasChildren
		case *notionapi.FileBlock:
			return notionapi.BlockID(b.ID), b.HasChildren
		case *notionapi.VideoBlock:
			return notionapi.BlockID(b.ID), b.HasChildren
		case *notionapi.TableBlock:
			return notionapi.BlockID(b.ID), b.HasChildren
		case *notionapi.TableRowBlock:
			return notionapi.BlockID(b.ID), b.HasChildren
		case *notionapi.ColumnListBlock:
			return notionapi.BlockID(b.ID), b.HasChildren
		case *notionapi.ColumnBlock:
			return notionapi.BlockID(b.ID), b.HasChildren
		default:
			return "", false
		}
	}

	var renderBlock func(notionapi.Block) (string, bool, error)
	renderBlock = func(block notionapi.Block) (string, bool, error) {
		childContent := ""
		if id, has := getBlockIDAndHasChildren(block); has && getChildren != nil {
			children, err := getChildren(id)
			if err != nil {
				return "", false, err
			}
			prevChildIsList := false
			_, isColumnList := block.(*notionapi.ColumnListBlock)
			for _, cb := range children {
				cstr, childIsList, err := renderBlock(cb)
				if err != nil {
					return "", false, err
				}
				indent := ""
				switch block.(type) {
				case *notionapi.BulletedListItemBlock, *notionapi.NumberedListItemBlock, *notionapi.ToDoBlock:
					indent = strings.Repeat(" ", 4)
				}
				lines := strings.Split(strings.TrimRight(cstr, "\n"), "\n")
				for i, l := range lines {
					if strings.TrimSpace(l) == "" {
						continue
					}
					lines[i] = indent + l
				}
				rendered := strings.Join(lines, "\n")
				sep := "\n\n"
				if prevChildIsList && childIsList {
					sep = "\n"
				}
				if childContent == "" {
					childContent = rendered
				} else {
					childContent += sep + rendered
				}
				prevChildIsList = childIsList
				if isColumnList {
					childContent += "\n__COLUMN_BREAK__\n"
				}
			}
			childContent = strings.TrimRight(childContent, "\n")
		}
		s, isList := blockToMarkdownWithCache(block, childContent, resolve, r.fileCache, articlePath, r.config)
		return strings.TrimRight(s, "\n"), isList, nil
	}

	markdown := ""
	prevIsList := false
	for _, block := range blocks {
		s, isList, err := renderBlock(block)
		if err != nil {
			return "", err
		}

		// Add separator before current block (except for first block)
		if markdown != "" {
			if prevIsList && isList {
				markdown += "\n"
			} else {
				markdown += "\n\n"
			}
		}

		// Add the block content
		markdown += s
		prevIsList = isList
	}
	return markdown, nil
}

// helper: simple slugifier for file names
func slugify(s string) string {
	s = strings.ToLower(s)
	s = strings.ReplaceAll(s, " ", "-")
	safe := make([]rune, 0, len(s))
	for _, r := range s {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '-' || r == '_' {
			safe = append(safe, r)
		}
	}
	return string(safe)
}
