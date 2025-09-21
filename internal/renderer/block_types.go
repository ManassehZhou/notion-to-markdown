package renderer

import (
	"net/url"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/jomei/notionapi"
)

// block_types contains helpers to convert Notion block types into Markdown.
// The functions here are internal implementation details used by the public
// Renderer type in renderer.go.

// blockToMarkdownWithCache converts a Notion block into Markdown with file caching support.
// childContent is used when a block has pre-rendered child content (for example, for toggles or columns).
// It returns the markdown string and a boolean indicating whether the block is
// a list item (used to control spacing between list items).
func blockToMarkdownWithCache(block notionapi.Block, childContent string, resolve func(string) string, fileCache *FileCache, articlePath string, config *RenderConfig) (string, bool) {
	switch b := block.(type) {
	case *notionapi.ParagraphBlock:
		return paragraphToMarkdown(b, resolve), false
	case *notionapi.Heading1Block:
		return heading1ToMarkdown(b, resolve), false
	case *notionapi.Heading2Block:
		return heading2ToMarkdown(b, resolve), false
	case *notionapi.Heading3Block:
		return heading3ToMarkdown(b, resolve), false
	case *notionapi.BulletedListItemBlock:
		return bulletedListItemToMarkdown(b, childContent, resolve), true
	case *notionapi.NumberedListItemBlock:
		return numberedListItemToMarkdown(b, childContent, resolve), true
	case *notionapi.ToDoBlock:
		return toDoToMarkdown(b, childContent, resolve), true
	case *notionapi.ToggleBlock:
		return toggleToMarkdown(b, childContent, resolve, config), false
	case *notionapi.EquationBlock:
		return equationToMarkdown(b, resolve, config), false
	case *notionapi.CodeBlock:
		return codeToMarkdown(b, resolve), false
	case *notionapi.QuoteBlock:
		return quoteToMarkdown(b, resolve), false
	case *notionapi.CalloutBlock:
		return calloutToMarkdown(b, childContent, resolve, config), false
	case *notionapi.DividerBlock:
		return dividerToMarkdown(b), false
	case *notionapi.ImageBlock:
		return imageToMarkdownWithCache(b, fileCache, articlePath), false
	case *notionapi.BookmarkBlock:
		return bookmarkToMarkdown(b), false
	case *notionapi.EmbedBlock:
		return embedToMarkdown(b, config), false
	case *notionapi.LinkPreviewBlock:
		return linkPreviewToMarkdown(b), false
	case *notionapi.FileBlock:
		return fileToMarkdownWithCache(b, fileCache, articlePath, config), false
	case *notionapi.PdfBlock:
		return pdfToMarkdownWithCache(b, fileCache, articlePath, config), false
	case *notionapi.VideoBlock:
		return videoToMarkdownWithCache(b, fileCache, articlePath, config), false
	case *notionapi.TableBlock:
		return tableToMarkdown(b, childContent), false
	case *notionapi.TableRowBlock:
		return tableRowToMarkdown(b, resolve), false
	case *notionapi.ColumnListBlock:
		return columnListToMarkdown(b, childContent), false
	case *notionapi.ColumnBlock:
		return columnToMarkdown(b, childContent), false
	default:
		return "", false
	}
}

func paragraphToMarkdown(b *notionapi.ParagraphBlock, resolve func(string) string) string {
	return richTextArrToMarkdown(b.Paragraph.RichText, resolve)
}

func heading1ToMarkdown(b *notionapi.Heading1Block, resolve func(string) string) string {
	return "# " + richTextArrToMarkdown(b.Heading1.RichText, resolve)
}

func heading2ToMarkdown(b *notionapi.Heading2Block, resolve func(string) string) string {
	return "## " + richTextArrToMarkdown(b.Heading2.RichText, resolve)
}

func heading3ToMarkdown(b *notionapi.Heading3Block, resolve func(string) string) string {
	return "### " + richTextArrToMarkdown(b.Heading3.RichText, resolve)
}

// renderListItemWithChild renders a list item with base content and optional child content
func renderListItemWithChild(base string, childContent string) string {
	if childContent == "" {
		return base
	}
	return base + "\n" + childContent
}

func bulletedListItemToMarkdown(b *notionapi.BulletedListItemBlock, childContent string, resolve func(string) string) string {
	base := "- " + richTextArrToMarkdown(b.BulletedListItem.RichText, resolve)
	return renderListItemWithChild(base, childContent)
}

func numberedListItemToMarkdown(b *notionapi.NumberedListItemBlock, childContent string, resolve func(string) string) string {
	base := "1. " + richTextArrToMarkdown(b.NumberedListItem.RichText, resolve)
	return renderListItemWithChild(base, childContent)
}

func toDoToMarkdown(b *notionapi.ToDoBlock, childContent string, resolve func(string) string) string {
	checked := " "
	if b.ToDo.Checked {
		checked = "x"
	}
	base := "- [" + checked + "] " + richTextArrToMarkdown(b.ToDo.RichText, resolve)
	return renderListItemWithChild(base, childContent)
}

func toggleToMarkdown(b *notionapi.ToggleBlock, childContent string, resolve func(string) string, config *RenderConfig) string {
	summary := richTextArrToMarkdown(b.Toggle.RichText, resolve)
	if childContent == "" {
		return "> " + summary
	}
	childContent = dedentChildContent(childContent)

	data := map[string]string{
		"Summary": summary,
		"Content": childContent,
	}
	return renderTemplate(config.DetailsTemplate, data)
}

func codeToMarkdown(b *notionapi.CodeBlock, resolve func(string) string) string {
	return "```" + b.Code.Language + "\n" + richTextArrToMarkdown(b.Code.RichText, resolve) + "\n```"
}

func equationToMarkdown(b *notionapi.EquationBlock, resolve func(string) string, config *RenderConfig) string {
	if b.Equation.Expression != "" {
		data := map[string]string{
			"Expression": b.Equation.Expression,
		}
		return renderTemplate(config.MathTemplate, data)
	}
	return ""
}

func quoteToMarkdown(b *notionapi.QuoteBlock, resolve func(string) string) string {
	return "> " + richTextArrToMarkdown(b.Quote.RichText, resolve)
}

func calloutToMarkdown(b *notionapi.CalloutBlock, childContent string, resolve func(string) string, config *RenderConfig) string {
	contentText := richTextArrToMarkdown(b.Callout.RichText, resolve)
	if childContent != "" {
		childContent = dedentChildContent(childContent)
		lines := strings.Split(childContent, "\n")
		addSeparator := false
		for _, l := range lines {
			t := strings.TrimSpace(l)
			if t == "" {
				continue
			}
			if strings.HasPrefix(t, "-") || strings.HasPrefix(t, "1.") || strings.HasPrefix(t, "*") || strings.HasPrefix(t, "<") || strings.HasPrefix(t, "|") {
				addSeparator = false
			} else {
				addSeparator = true
			}
			break
		}
		childLines := make([]string, 0, len(lines)+1)
		if addSeparator {
			childLines = append(childLines, "> ")
		}
		for _, l := range lines {
			if strings.TrimSpace(l) == "" {
				childLines = append(childLines, "> ")
			} else {
				childLines = append(childLines, "> "+l)
			}
		}
		contentText += "\n" + strings.Join(childLines, "\n")
	}

	data := map[string]string{
		"Content": contentText,
	}
	return renderTemplate(config.CalloutTemplate, data)
}

func dividerToMarkdown(b *notionapi.DividerBlock) string {
	_ = b
	return "---"
}

// processFileURL extracts URL and handles caching for Notion file/external blocks
type fileURLExtractor interface {
	getFileURL() (url string, shouldCache bool)
	getCaption() []notionapi.RichText
}

type imageURLExtractor struct{ block *notionapi.ImageBlock }

func (e imageURLExtractor) getFileURL() (string, bool) {
	if e.block.Image.File != nil {
		return e.block.Image.File.URL, true
	} else if e.block.Image.External != nil {
		return e.block.Image.External.URL, false
	}
	return "", false
}
func (e imageURLExtractor) getCaption() []notionapi.RichText { return e.block.Image.Caption }

type fileURLExtractorImpl struct{ block *notionapi.FileBlock }

func (e fileURLExtractorImpl) getFileURL() (string, bool) {
	if e.block.File.File != nil {
		return e.block.File.File.URL, true
	} else if e.block.File.External != nil {
		return e.block.File.External.URL, false
	}
	return "", false
}
func (e fileURLExtractorImpl) getCaption() []notionapi.RichText { return e.block.File.Caption }

type pdfURLExtractor struct{ block *notionapi.PdfBlock }

func (e pdfURLExtractor) getFileURL() (string, bool) {
	if e.block.Pdf.File != nil {
		return e.block.Pdf.File.URL, true
	} else if e.block.Pdf.External != nil {
		return e.block.Pdf.External.URL, false
	}
	return "", false
}
func (e pdfURLExtractor) getCaption() []notionapi.RichText { return e.block.Pdf.Caption }

type videoURLExtractor struct{ block *notionapi.VideoBlock }

func (e videoURLExtractor) getFileURL() (string, bool) {
	if e.block.Video.File != nil {
		return e.block.Video.File.URL, true
	} else if e.block.Video.External != nil {
		return e.block.Video.External.URL, false
	}
	return "", false
}
func (e videoURLExtractor) getCaption() []notionapi.RichText { return e.block.Video.Caption }

func processFileURLWithCache(extractor fileURLExtractor, fileCache *FileCache, articlePath string) (url, text string) {
	var shouldCache bool
	originalURL, shouldCache := extractor.getFileURL()

	if originalURL == "" {
		return "", ""
	}

	// Extract text from caption using original URL
	caption := extractor.getCaption()
	if len(caption) > 0 {
		text = captionFirstParagraph(caption, nil)
	}
	if text == "" {
		text = escapeMarkdown(shortenURLLabel(originalURL))
	}

	// Cache the file only if it's a Notion-hosted file
	url = originalURL
	if shouldCache && fileCache != nil && articlePath != "" {
		if cachedPath, err := fileCache.CacheFile(originalURL, articlePath); err == nil {
			url = cachedPath
		}
		// If caching fails, fall back to original URL
	}

	return url, text
}

func imageToMarkdownWithCache(b *notionapi.ImageBlock, fileCache *FileCache, articlePath string) string {
	url, alt := processFileURLWithCache(imageURLExtractor{b}, fileCache, articlePath)
	if url == "" {
		return ""
	}
	return "![" + escapeMarkdown(alt) + "](" + url + ")"
}

// renderLinkWithCaption creates a markdown link with optional caption text
func renderLinkWithCaption(url string, caption []notionapi.RichText) string {
	if len(caption) > 0 {
		text := captionFirstParagraph(caption, nil)
		if text != "" {
			return "[" + text + "](" + url + ")"
		}
	}
	return "[" + escapeMarkdown(shortenURLLabel(url)) + "](" + url + ")"
}

func bookmarkToMarkdown(b *notionapi.BookmarkBlock) string {
	return renderLinkWithCaption(b.Bookmark.URL, b.Bookmark.Caption)
}

func tableToMarkdown(block *notionapi.TableBlock, childContent string) string {
	childContent = dedentChildContent(childContent)
	s := strings.TrimSpace(childContent)
	if s == "" {
		return ""
	}
	lines := strings.Split(s, "\n")
	if len(lines) == 0 {
		return ""
	}
	parsed := make([][]string, 0, len(lines))
	maxCols := 0
	for _, ln := range lines {
		if strings.TrimSpace(ln) == "" {
			continue
		}
		var parts []string
		if strings.Contains(ln, " | ") {
			parts = strings.Split(ln, " | ")
		} else {
			parts = strings.Split(ln, "|")
		}
		for i := range parts {
			parts[i] = strings.TrimSpace(parts[i])
		}
		if len(parts) > maxCols {
			maxCols = len(parts)
		}
		parsed = append(parsed, parts)
	}
	if len(parsed) == 0 {
		return ""
	}
	normalized := make([]string, 0, len(parsed))
	for _, parts := range parsed {
		if len(parts) < maxCols {
			for i := len(parts); i < maxCols; i++ {
				parts = append(parts, "")
			}
		}
		normalized = append(normalized, "| "+strings.Join(parts, " | ")+" |")
	}
	if block.Table.HasColumnHeader && len(normalized) > 0 {
		sep := make([]string, maxCols)
		for i := range sep {
			sep[i] = "---"
		}
		header := normalized[0]
		rest := ""
		if len(normalized) > 1 {
			rest = "\n" + strings.Join(normalized[1:], "\n")
		}
		return header + "\n| " + strings.Join(sep, " | ") + " |" + rest
	}
	return strings.Join(normalized, "\n")
}

func tableRowToMarkdown(block *notionapi.TableRowBlock, resolve func(string) string) string {
	cells := block.TableRow.Cells
	if len(cells) == 0 {
		return ""
	}
	cols := make([]string, 0, len(cells))
	for _, cell := range cells {
		cols = append(cols, strings.TrimSpace(richTextArrToMarkdown(cell, resolve)))
	}
	return strings.Join(cols, " | ")
}

func embedToMarkdown(b *notionapi.EmbedBlock, config *RenderConfig) string {
	url := b.Embed.URL
	text := ""
	if len(b.Embed.Caption) > 0 {
		text = captionFirstParagraph(b.Embed.Caption, nil)
	}
	if text == "" {
		text = escapeMarkdown(shortenURLLabel(url))
	}

	data := map[string]string{
		"URL":  url,
		"Text": text,
	}
	return renderTemplate(config.EmbedTemplate, data)
}

func columnListToMarkdown(b *notionapi.ColumnListBlock, childContent string) string {
	_ = b
	if strings.TrimSpace(childContent) == "" {
		return ""
	}
	childContent = dedentChildContent(childContent)
	parts := strings.Split(childContent, "__COLUMN_BREAK__")
	cols := make([]string, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p == "" {
			continue
		}
		cols = append(cols, "<td>\n\n"+p+"\n</td>")
	}
	if len(cols) == 0 {
		return ""
	}
	return "<table><tr>" + strings.Join(cols, "") + "</tr></table>"
}

func columnToMarkdown(b *notionapi.ColumnBlock, childContent string) string {
	_ = b
	return dedentChildContent(childContent)
}

func linkPreviewToMarkdown(b *notionapi.LinkPreviewBlock) string {
	text := shortenURLLabel(b.LinkPreview.URL)
	return "[" + escapeMarkdown(text) + "](" + b.LinkPreview.URL + ")"
}

func fileToMarkdownWithCache(b *notionapi.FileBlock, fileCache *FileCache, articlePath string, config *RenderConfig) string {
	url, text := processFileURLWithCache(fileURLExtractorImpl{b}, fileCache, articlePath)
	if url == "" {
		return ""
	}

	data := map[string]string{
		"URL":  url,
		"Text": text,
	}
	return renderTemplate(config.FileTemplate, data)
}

func pdfToMarkdownWithCache(b *notionapi.PdfBlock, fileCache *FileCache, articlePath string, config *RenderConfig) string {
	url, text := processFileURLWithCache(pdfURLExtractor{b}, fileCache, articlePath)
	if url == "" {
		return ""
	}

	data := map[string]string{
		"URL":  url,
		"Text": text,
	}
	return renderTemplate(config.PDFTemplate, data)
}

func videoToMarkdownWithCache(b *notionapi.VideoBlock, fileCache *FileCache, articlePath string, config *RenderConfig) string {
	url, text := processFileURLWithCache(videoURLExtractor{b}, fileCache, articlePath)
	if url == "" {
		return ""
	}

	data := map[string]string{
		"URL":  url,
		"Text": text,
	}
	return renderTemplate(config.VideoTemplate, data)
}

func richTextArrToMarkdown(arr []notionapi.RichText, resolve func(string) string) string {
	result := ""
	for _, t := range arr {
		txt := t.PlainText
		if t.Href != "" {
			url := t.Href
			// If the link points to a Notion page, convert it to a Hugo site link.
			// Use resolver when available.
			if resolve != nil {
				// the notionURLToHugoLink will attempt to extract an ID and call resolve
				url = notionURLToHugoLink(url, resolve)
			} else {
				url = notionURLToHugoLink(url, nil)
			}
			txt = "[" + escapeMarkdown(richTextAnnotationsToMarkdown(t)) + "](" + url + ")"
			result += txt
			continue
		}
		if t.Annotations.Code {
			result += "`" + escapeBackticks(txt) + "`"
			continue
		}
		wrapped := txt
		if t.Annotations.Bold {
			wrapped = "**" + wrapped + "**"
		}
		if t.Annotations.Italic {
			wrapped = "*" + wrapped + "*"
		}
		if t.Annotations.Strikethrough {
			wrapped = "~~" + wrapped + "~~"
		}
		if t.Annotations.Underline {
			wrapped = "<u>" + wrapped + "</u>"
		}
		result += wrapped
	}
	return result
}

// notionURLToHugoLink converts a Notion page URL to a site-relative link
// for static site generators when possible. Example: https://www.notion.so/Workspace-Page-Title-<uuid>
// becomes the appropriate path based on the page type (posts, gallery, etc.).
// If the URL does not look like a Notion page link it is returned unchanged.
func notionURLToHugoLink(raw string, resolve func(string) string) string {
	if raw == "" {
		return raw
	}

	// Parse URL using net/url for proper handling
	parsedURL, err := url.Parse(raw)
	if err != nil || parsedURL.Host == "" {
		return raw
	}

	// Quick check for notion domain
	if !strings.Contains(parsedURL.Host, "notion.so") {
		return raw
	}

	// Get the last path segment
	pathSegments := strings.Split(strings.Trim(parsedURL.Path, "/"), "/")
	if len(pathSegments) == 0 {
		return raw
	}
	lastSegment := pathSegments[len(pathSegments)-1]

	// Extract UUID from the last segment using regex
	// Notion URLs can have UUID with or without dashes, and may be directly concatenated with title
	// Pattern 1: UUID with dashes (36 chars): 8-4-4-4-12
	// Pattern 2: UUID without dashes (32 chars): all together
	var uuid, titlePart string

	// First try to find UUID with dashes
	reDashed := regexp.MustCompile(`(?i)([0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12})$`)
	if matches := reDashed.FindStringSubmatch(lastSegment); len(matches) > 0 {
		uuid = matches[1]
		// Remove UUID from the end to get title part
		titlePart = strings.TrimSuffix(lastSegment, "-"+uuid)
		if titlePart == lastSegment {
			// Try without dash separator
			titlePart = strings.TrimSuffix(lastSegment, uuid)
		}
	} else {
		// Try to find UUID without dashes (32 hex chars)
		reNoDash := regexp.MustCompile(`(?i)([0-9a-f]{32})$`)
		if matches := reNoDash.FindStringSubmatch(lastSegment); len(matches) > 0 {
			uuid = matches[1]
			// Remove UUID from the end to get title part
			titlePart = strings.TrimSuffix(lastSegment, uuid)
			// If there was a dash before the UUID, remove it too
			titlePart = strings.TrimSuffix(titlePart, "-")
		}
	}

	// If no UUID found, return original URL
	if uuid == "" {
		return raw
	}
	// Normalize UUID by removing dashes to match pageMap key format
	normalizedUUID := strings.ReplaceAll(uuid, "-", "")

	// If we have a resolver, try to resolve the UUID to the correct path
	if resolve != nil {
		if resolvedPath := resolve(normalizedUUID); resolvedPath != "" {
			return resolvedPath
		}
	}

	// Fallback: extract title and create a generic posts link
	if titlePart == "" {
		titlePart = lastSegment
	}
	if uuid != "" {
		// Title part was already extracted above during UUID detection
		// No additional processing needed here
	}

	// Convert title to slug
	slug := slugify(strings.ReplaceAll(titlePart, "-", " "))
	if slug == "" {
		// Fallback to UUID if no title found
		slug = normalizedUUID
	}

	// Default fallback to posts path if resolver failed
	return "/posts/" + slug + "/"
}

func richTextAnnotationsToMarkdown(t notionapi.RichText) string {
	txt := t.PlainText
	if t.Annotations.Code {
		return "`" + escapeBackticks(txt) + "`"
	}
	wrapped := txt
	if t.Annotations.Bold {
		wrapped = "**" + wrapped + "**"
	}
	if t.Annotations.Italic {
		wrapped = "*" + wrapped + "*"
	}
	if t.Annotations.Strikethrough {
		wrapped = "~~" + wrapped + "~~"
	}
	if t.Annotations.Underline {
		wrapped = "<u>" + wrapped + "</u>"
	}
	return wrapped
}

func escapeBackticks(s string) string {
	return strings.ReplaceAll(s, "`", "\\`")
}

func escapeMarkdown(s string) string {
	s = strings.ReplaceAll(s, "[", "\\[")
	s = strings.ReplaceAll(s, "]", "\\]")
	return s
}

func dedentChildContent(childContent string) string {
	if childContent == "" {
		return childContent
	}
	lines := strings.Split(childContent, "\n")
	for i, l := range lines {
		if strings.HasPrefix(l, "    ") {
			lines[i] = strings.TrimPrefix(l, "    ")
		}
	}
	return strings.Join(lines, "\n")
}

func shortenURLLabel(raw string) string {
	if raw == "" {
		return ""
	}

	max := 40

	// Try to parse as URL first
	if u, err := url.Parse(raw); err == nil && u.Scheme != "" {
		// Extract filename from path
		filename := filepath.Base(u.Path)

		// Check if it's a meaningful filename (has extension and not just "/")
		if filename != "." && filename != "/" && strings.Contains(filename, ".") {
			// If filename is short enough, show domain + filename
			if len(filename) <= max-10 && u.Host != "" {
				if len(u.Host)+len(filename)+5 <= max { // +5 for "/.../""
					return u.Host + "/.../" + filename
				}
			}

			// If just filename fits, show it
			if len(filename) <= max {
				return filename
			}

			// Show truncated filename, keeping extension if possible
			if len(filename) > max-3 {
				ext := filepath.Ext(filename)
				if len(ext) <= 8 && len(ext) > 0 {
					nameLen := max - len(ext) - 3
					if nameLen > 0 {
						base := filename[:len(filename)-len(ext)]
						if len(base) > nameLen {
							return base[:nameLen] + "..." + ext
						}
					}
				}
				return filename[:max-3] + "..."
			}
		}

		// Fallback to host + path logic
		if u.Host != "" {
			// Remove protocol, show clean URL
			cleanURL := u.Host + u.Path
			if u.RawQuery != "" {
				cleanURL += "?" + u.RawQuery
			}

			if len(cleanURL) <= max {
				return cleanURL
			}

			// Try host + first path segment
			pathParts := strings.Split(strings.Trim(u.Path, "/"), "/")
			if len(pathParts) > 0 && pathParts[0] != "" {
				candidate := u.Host + "/" + pathParts[0]
				if len(candidate) <= max {
					return candidate + "..."
				}
			}

			// Just host if it fits
			if len(u.Host) <= max-3 {
				return u.Host + "..."
			}
		}
	}

	// Fallback for non-URLs or unparseable URLs
	if len(raw) <= max {
		return raw
	}

	return raw[:max-3] + "..."
}

func captionFirstParagraph(arr []notionapi.RichText, resolve func(string) string) string {
	if len(arr) == 0 {
		return ""
	}
	full := richTextArrToMarkdown(arr, resolve)
	parts := strings.Split(full, "\n\n")
	if len(parts) == 0 {
		return strings.TrimSpace(full)
	}
	return strings.TrimSpace(parts[0])
}
