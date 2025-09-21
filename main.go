package main

import (
	"flag"
	"log/slog"
	"os"
	"strings"

	"github.com/ManassehZhou/notion-to-markdown/internal/notionclient"
	"github.com/ManassehZhou/notion-to-markdown/internal/renderer"
	"github.com/ManassehZhou/notion-to-markdown/internal/writer"

	"github.com/jomei/notionapi"
)

// main is the CLI entrypoint. It reads configuration from flags or environment
// variables, queries a Notion database for pages, converts each page to a
// Markdown file (with YAML front matter), and writes the resulting files to
// disk. Compatible with Hugo, Hexo, Jekyll, and other static site generators.
func newLogger(level slog.Level) *slog.Logger {
	return slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: level,
	}))
}

func main() {
	// Setup structured logging
	logger := newLogger(slog.LevelInfo)
	slog.SetDefault(logger)

	slog.Info("🚀 Notion to Markdown Converter")

	// CLI flags with environment fallbacks
	tokenFlag := flag.String("token", "", "Notion integration token (or set NOTION_TOKEN)")
	dbFlag := flag.String("database", "", "Notion database ID (or set NOTION_DATABASE_ID)")
	outFlag := flag.String("out", "content", "Output directory for generated markdown files")
	configFlag := flag.String("config", "config/notion-to-markdown.yaml", "Path to YAML configuration file")
	verboseFlag := flag.Bool("verbose", false, "Enable verbose logging")
	flag.Parse()

	notionToken := *tokenFlag
	if notionToken == "" {
		notionToken = os.Getenv("NOTION_TOKEN")
	}
	databaseID := *dbFlag
	if databaseID == "" {
		databaseID = os.Getenv("NOTION_DATABASE_ID")
	}
	outDir := *outFlag
	configPath := *configFlag
	verbose := *verboseFlag

	// Enable verbose logging in GitHub Actions environment
	if os.Getenv("GITHUB_ACTIONS") == "true" || os.Getenv("VERBOSE") == "true" {
		verbose = true
	}

	// Update log level based on verbose flag
	if verbose {
		logger = newLogger(slog.LevelDebug)
		slog.SetDefault(logger)
	}

	if notionToken == "" || databaseID == "" {
		slog.Error("❌ Error: Missing required parameters")
		slog.Info("Usage: notion-to-markdown -token TOKEN -database DATABASE_ID [-out DIR] [-config CONFIG.yaml]")
		slog.Info("You can also set NOTION_TOKEN and NOTION_DATABASE_ID environment variables.")
		os.Exit(1)
	}

	if verbose {
		slog.Debug("📂 Output directory", "path", outDir)
		slog.Debug("⚙️ Config file", "path", configPath)
		slog.Debug("🗄️ Database ID", "id", databaseID)
	}

	nc := notionclient.New(notionToken)
	// We'll build a resolver map from the database pages so internal Notion links
	// can be converted to site-relative Hugo links.
	w := writer.New()

	// Load render configuration from YAML file
	if verbose {
		slog.Debug("📄 Loading configuration", "path", configPath)
	}
	config := renderer.LoadConfigWithFallback(configPath)

	if verbose {
		slog.Info("🔄 Fetching pages from Notion database...")
	}
	pages, err := nc.FetchPages(databaseID)
	if err != nil {
		slog.Error("❌ Failed to query Notion database", "error", err)
		os.Exit(1)
	}

	if verbose {
		slog.Info("📊 Found pages in database", "count", len(pages))
		if len(pages) > 100 {
			slog.Warn("Large number of pages detected, processing may take time", "count", len(pages))
		}
	}

	// build page ID -> path map for resolver
	if verbose {
		slog.Info("🔗 Building page resolver map...")
	}
	pageMap := map[string]string{}
	resolve := func(pageID string) string {
		if path, ok := pageMap[pageID]; ok {
			return path
		}
		return ""
	}
	r := renderer.New(resolve, outDir, config)

	for _, p := range pages {
		// Get the full path including content type, not just slug
		path := r.GetPagePath(p)

		// Use page ID and normalize it by removing dashes
		normalizedID := strings.ReplaceAll(string(p.ID), "-", "")
		pageMap[normalizedID] = path
	}

	// Update renderer with the resolver
	filesGenerated := 0

	if verbose {
		slog.Info("📝 Converting pages to Markdown...")
	}

	for i, p := range pages {
		if verbose {
			slog.Debug("Processing page", "current", i+1, "total", len(pages))
		}

		// Fetch top-level blocks for the page (convert ObjectID to BlockID)
		blocks, err := nc.GetChildren(notionapi.BlockID(p.ID))
		if err != nil {
			slog.Error("❌ Failed to fetch page blocks", "error", err)
			os.Exit(1)
		}
		filename, content, err := r.RenderPage(p, blocks, nc.GetChildren, resolve)
		if err != nil {
			slog.Error("❌ Failed to render page", "error", err)
			os.Exit(1)
		}
		// ensure we write into the requested output directory
		// if filename already contains a top-level path like "posts/..." we keep it,
		// otherwise prefix with outDir
		finalPath := filename
		if outDir != "" && !strings.HasPrefix(filename, outDir+"/") {
			finalPath = outDir + "/" + filename
		}

		if err := w.WriteFile(finalPath, content); err != nil {
			slog.Error("❌ Failed to write file", "error", err)
			os.Exit(1)
		}

		if verbose {
			slog.Info("✅ Generated file", "path", finalPath)
		} else {
			// Print progress dot for non-verbose mode
			print(".")
		}
		filesGenerated++
	}

	if !verbose {
		println() // New line after dots
	}

	slog.Info("🎉 Successfully generated markdown files", "count", filesGenerated, "directory", outDir)

	// Warn about large numbers of files
	if filesGenerated > 50 {
		slog.Warn("Large number of files generated, check repository size limits", "count", filesGenerated)
	}
}
