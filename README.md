# Notion to Markdown GitHub Action

[![GitHub release](https://img.shields.io/github/release/ManassehZhou/notion-to-hugo.svg)](https://GitHub.com/ManassehZhou/notion-to-hugo/releases/)
[![GitHub license](https://img.shields.io/badge/license-Apache%202.0-blue.svg)](https://github.com/ManassehZhou/notion-to-hugo/blob/main/LICENSE)
[![GitHub stars](https://img.shields.io/github/stars/ManassehZhou/notion-to-hugo.svg)](https://GitHub.com/ManassehZhou/notion-to-hugo/stargazers/)

🚀 A GitHub Action that automatically fetches content from your Notion database and converts it into Markdown files with front matter. Compatible with popular static site generators like Hugo, Hexo, Jekyll, and more.

## ✨ Features

- 📝 **Markdown Conversion**: Transform Notion pages into clean Markdown with front matter
- 🔄 **Automated Workflow**: Set up CI/CD pipelines for content publishing
- 🎨 **Multi-Platform Support**: Compatible with Hugo, Hexo, Jekyll, and other static site generators
- 📁 **File Organization**: Automatic file organization based on Notion properties
- 🔗 **Internal Links**: Convert Notion page links to relative URLs
- 🖼️ **Media Handling**: Download and organize images, PDFs, and other attachments
- ⚙️ **Flexible Configuration**: Customize rendering with YAML configuration files
- 🏷️ **Front Matter**: Generate YAML front matter from Notion properties

## 📋 Prerequisites

1. **Notion Integration**: Create a Notion integration and obtain an API token
2. **Notion Database**: Set up a database with your content and share it with your integration
3. **Static Site Generator**: Have a site repository (Hugo, Hexo, Jekyll, etc.) where you want to publish content

### Setting up Notion Integration

1. Go to [Notion Developers](https://www.notion.so/my-integrations)
2. Click "Create new integration"
3. Give it a name and select your workspace
4. Copy the "Internal Integration Token"

## 🚀 Quick Start

### Basic Usage (Hugo)

Create a workflow file (e.g., `.github/workflows/notion-to-hugo.yml`):

```yaml
name: Update Hugo Content from Notion

on:
  schedule:
    # Run every day at 6 AM UTC
    - cron: '0 6 * * *'
  workflow_dispatch:

jobs:
  notion-to-markdown:
    runs-on: ubuntu-latest
    steps:
      - name: Checkout repository
        uses: actions/checkout@v4

      - name: Fetch content from Notion
        uses: ManassehZhou/notion-to-hugo@v1
        with:
          notion-token: ${{ secrets.NOTION_TOKEN }}
          notion-database-id: ${{ secrets.NOTION_DATABASE_ID }}
          output-directory: 'content/posts'

      - name: Setup Hugo
        uses: peaceiris/actions-hugo@v2
        with:
          hugo-version: '0.119.0'
          extended: true

      - name: Build Hugo site
        run: hugo --minify

      - name: Deploy to GitHub Pages
        uses: peaceiris/actions-gh-pages@v3
        with:
          github_token: ${{ secrets.GITHUB_TOKEN }}
          publish_dir: ./public
```

### Basic Usage (Hexo)

```yaml
name: Update Hexo Content from Notion

on:
  schedule:
    - cron: '0 6 * * *'
  workflow_dispatch:

jobs:
  notion-to-markdown:
    runs-on: ubuntu-latest
    steps:
      - name: Checkout repository
        uses: actions/checkout@v4

      - name: Fetch content from Notion
        uses: ManassehZhou/notion-to-hugo@v1
        with:
          notion-token: ${{ secrets.NOTION_TOKEN }}
          notion-database-id: ${{ secrets.NOTION_DATABASE_ID }}
          output-directory: 'source/_posts'
          config-file: 'config/hexo-notion.yaml'

      - name: Setup Node.js
        uses: actions/setup-node@v3
        with:
          node-version: '18'
          cache: 'npm'

      - name: Install dependencies
        run: npm install

      - name: Generate Hexo site
        run: npx hexo generate

      - name: Deploy to GitHub Pages
        uses: peaceiris/actions-gh-pages@v3
        with:
          github_token: ${{ secrets.GITHUB_TOKEN }}
          publish_dir: ./public
```

### Complete CI/CD Pipeline Examples

#### Hugo + GitHub Pages
```yaml
name: Hugo Site CI/CD

on:
  push:
    branches: [ main ]
  schedule:
    - cron: '0 6 * * *'
  workflow_dispatch:

jobs:
  build-and-deploy:
    runs-on: ubuntu-latest
    steps:
      - name: Checkout
        uses: actions/checkout@v4
        with:
          submodules: true
          fetch-depth: 0

      - name: Fetch latest content from Notion
        uses: ManassehZhou/notion-to-hugo@v1
        with:
          notion-token: ${{ secrets.NOTION_TOKEN }}
          notion-database-id: ${{ secrets.NOTION_DATABASE_ID }}
          output-directory: 'content/posts'
          config-file: 'config/hugo-notion.yaml'

      - name: Setup Hugo
        uses: peaceiris/actions-hugo@v2
        with:
          hugo-version: '0.119.0'  # Use specific version for reproducible builds
          extended: true

      - name: Build Hugo site
        run: hugo --minify

      - name: Deploy to GitHub Pages
        uses: peaceiris/actions-gh-pages@v3
        if: github.ref == 'refs/heads/main'
        with:
          github_token: ${{ secrets.GITHUB_TOKEN }}
          publish_dir: ./public
```

#### Hexo + GitHub Pages
```yaml
name: Hexo Site CI/CD

on:
  push:
    branches: [ main ]
  schedule:
    - cron: '0 6 * * *'
  workflow_dispatch:

jobs:
  build-and-deploy:
    runs-on: ubuntu-latest
    steps:
      - name: Checkout
        uses: actions/checkout@v4

      - name: Fetch latest content from Notion
        uses: ManassehZhou/notion-to-hugo@v1
        with:
          notion-token: ${{ secrets.NOTION_TOKEN }}
          notion-database-id: ${{ secrets.NOTION_DATABASE_ID }}
          output-directory: 'source/_posts'
          config-file: 'config/hexo-notion.yaml'

      - name: Setup Node.js
        uses: actions/setup-node@v3
        with:
          node-version: '18'
          cache: 'npm'

      - name: Install dependencies
        run: npm install

      - name: Generate Hexo site
        run: npx hexo generate

      - name: Deploy to GitHub Pages
        uses: peaceiris/actions-gh-pages@v3
        if: github.ref == 'refs/heads/main'
        with:
          github_token: ${{ secrets.GITHUB_TOKEN }}
          publish_dir: ./public
```

#### Jekyll + GitHub Pages
```yaml
name: Jekyll Site CI/CD

on:
  push:
    branches: [ main ]
  schedule:
    - cron: '0 6 * * *'
  workflow_dispatch:

jobs:
  build-and-deploy:
    runs-on: ubuntu-latest
    steps:
      - name: Checkout
        uses: actions/checkout@v4

      - name: Fetch latest content from Notion
        uses: ManassehZhou/notion-to-hugo@v1
        with:
          notion-token: ${{ secrets.NOTION_TOKEN }}
          notion-database-id: ${{ secrets.NOTION_DATABASE_ID }}
          output-directory: '_posts'
          config-file: 'config/jekyll-notion.yaml'

      - name: Setup Ruby
        uses: ruby/setup-ruby@v1
        with:
          ruby-version: '3.1'
          bundler-cache: true

      - name: Build Jekyll site
        run: bundle exec jekyll build

      - name: Deploy to GitHub Pages
        uses: peaceiris/actions-gh-pages@v3
        if: github.ref == 'refs/heads/main'
        with:
          github_token: ${{ secrets.GITHUB_TOKEN }}
          publish_dir: ./_site
```

## ⚙️ Configuration

### Action Inputs

| Input | Description | Required | Default |
|-------|-------------|----------|---------|
| `notion-token` | Notion integration token | ✅ | - |
| `notion-database-id` | Notion database ID to fetch content from | ✅ | - |
| `output-directory` | Output directory for generated markdown files | ❌ | `content` |
| `config-file` | Path to YAML configuration file for customizing Markdown output | ❌ | `config/notion-to-markdown.yaml` |

### Action Outputs

| Output | Description |
|--------|-------------|
| `files-generated` | Number of markdown files generated |
| `output-path` | Path where files were generated |
| `success` | Whether the operation completed successfully (true/false) |

### Configuration File

Create a configuration file to customize the Markdown output for your static site generator:

#### Universal Configuration (`notion-to-markdown.yaml`)
```yaml
# Compatible with most static site generators
math_template: |
  $$
  {{.Expression}}
  $$

details_template: |
  <details>
  <summary>{{.Summary}}</summary>
  {{.Content}}
  </details>

video_template: '<iframe src="{{.URL}}" frameborder="0" allowfullscreen></iframe>'
pdf_template: "[📄 {{.Text}}]({{.URL}})"
embed_template: '<iframe src="{{.URL}}" width="100%" height="400"></iframe>'
callout_template: "> **Note:** {{.Content}}"
file_template: "[📁 {{.Text}}]({{.URL}})"
```

#### Hugo-Specific Configuration
```yaml
# Optimized for Hugo shortcodes
math_template: |
  {{< math >}}
  {{.Expression}}
  {{< /math >}}

details_template: |
  {{< details "{{.Summary}}" >}}
  {{.Content}}
  {{< /details >}}

video_template: "{{< youtube {{.VideoID}} >}}"
pdf_template: "[📄 {{.Text}}]({{.URL}})"
embed_template: '<iframe src="{{.URL}}" width="100%" height="400"></iframe>'
callout_template: "> **{{.Icon}}** {{.Content}}"
file_template: "[📁 {{.Text}}]({{.URL}})"
```

#### Hexo-Specific Configuration  
```yaml
# Optimized for Hexo tags
math_template: |
  {% math %}
  {{.Expression}}
  {% endmath %}

details_template: |
  <details>
  <summary>{{.Summary}}</summary>
  {{.Content}}
  </details>

video_template: '{% video "{{.URL}}" %}'
pdf_template: "[📄 {{.Text}}]({{.URL}})"
embed_template: '<iframe src="{{.URL}}" width="100%" height="400"></iframe>'
callout_template: |
  {% note info %}
  {{.Content}}
  {% endnote %}
file_template: "[📁 {{.Text}}]({{.URL}})"
```

#### Jekyll-Specific Configuration
```yaml
# Optimized for Jekyll liquid tags
math_template: |
  $$
  {{.Expression}}
  $$

details_template: |
  <details>
  <summary>{{.Summary}}</summary>
  {{.Content}}
  </details>

video_template: '<iframe src="{{.URL}}" frameborder="0" allowfullscreen></iframe>'
pdf_template: "[📄 {{.Text}}]({{.URL}})"
embed_template: '<iframe src="{{.URL}}" width="100%" height="400"></iframe>'
callout_template: |
  {% include callout.html content="{{.Content}}" type="info" %}
file_template: "[📁 {{.Text}}]({{.URL}})"
```


## 📁 Notion Database Structure

Your Notion database should include these properties for optimal results. The action automatically maps these properties to front matter for static site generators:

### Required Properties

| Property Name | Notion Type | Front Matter | Description | Required |
|---------------|-------------|--------------|-------------|----------|
| `Title` or `Name` | Title | `title` | Page title (becomes the main title) | ✅ |

### Optional Properties

| Property Name | Notion Type | Front Matter | Description | Default Behavior |
|---------------|-------------|--------------|-------------|------------------|
| `Slug` | Rich Text or Title | `slug` | URL-safe identifier for the page | Auto-generated from title |
| `Date` | Date | `date` | Publication date in ISO format | Uses custom date if set, otherwise page creation time |
| `Tags` or `Tag` | Multi-select | `tags` | Content tags as array | Empty array |
| `Categories` or `Category` | Multi-select | `categories` | Content categories as array | Empty array |
| `Summary` or `Description` | Rich Text or Title | `summary` | Page summary/description | Empty if not provided |
| `Status` | Status | `draft` | Publication status | `draft: false` unless status is "Draft" |
| `Type` | Select | `type` + path | Content type affecting file path | Defaults to "posts" |

### Auto-Generated Properties

| Property Name | Front Matter | Description | Source |
|---------------|--------------|-------------|--------|
| `LastMod` | `lastmod` | Last modification time | Notion page's last edit time (automatic) |

### Property Details

#### 📝 **Title/Name** (Required)
- **Type**: Title
- **Usage**: Main page title, used for H1 and file organization
- **Example**: "Getting Started with Hugo"

#### 🔗 **Slug** (Optional)
- **Type**: Rich Text or Title
- **Usage**: URL-friendly version of the title
- **Auto-generation**: If not provided, automatically created from title
- **Example**: "getting-started-with-hugo"

#### 📅 **Date** (Optional)
- **Type**: Date
- **Usage**: Publication date for posts
- **Format**: ISO 8601 format in front matter
- **Date Priority Logic**:
  1. **Primary**: Uses custom `Date` property if set in Notion database
  2. **Fallback**: Uses Notion page's creation time if no custom date is provided
- **Behavior**: 
  - If you add a `Date` property to your database and set a value, it will override the page creation time
  - If no `Date` property exists or is empty, defaults to when the page was created in Notion
- **Example**: "2025-01-15" → outputs as "2025-01-15T00:00:00Z07:00"

#### 🏷️ **Tags/Tag** (Optional)
- **Type**: Multi-select
- **Usage**: Content categorization and filtering
- **Output**: Array in front matter
- **Example**: `["hugo", "tutorial", "getting-started"]`

#### � **Categories/Category** (Optional)
- **Type**: Multi-select
- **Usage**: Content categorization, often used for site navigation
- **Output**: Array in front matter
- **Example**: `["development", "web-development"]`

#### 📝 **Summary/Description** (Optional)
- **Type**: Rich Text or Title
- **Usage**: Brief description or summary of the content
- **Output**: String in front matter, useful for SEO and previews
- **Example**: "A comprehensive guide to getting started with Hugo"

#### �📊 **Status** (Optional)
- **Type**: Status
- **Usage**: Controls publication state
- **Values**: 
  - "Draft" → `draft: true` (hidden from site)
  - Any other value → `draft: false` (published)
- **Example**: "Published", "In Review", "Draft"

#### 📂 **Type** (Optional)
- **Type**: Select
- **Usage**: Determines content type and file path
- **Path generation**:
  - `"post"` → automatically converted to `"posts"` 
  - `"posts"` or empty → `/content/posts/slug/`
  - `"pages"` → `/content/slug/`
  - Other values → `/content/type/slug/`
- **Examples**: "posts", "docs", "blog", "tutorials"

#### ⏰ **LastMod** (Auto-Generated)
- **Type**: Automatic (not configurable)
- **Usage**: Static site generator's last modification timestamp for SEO and caching
- **Source**: Notion page's `last_edited_time` (automatically tracked by Notion)
- **Format**: ISO 8601 format in front matter
- **Behavior**: 
  - Always reflects when the Notion page was last edited
  - Updates automatically when you modify the page content in Notion
  - Cannot be manually overridden
- **Example**: "2025-01-15T10:30:00Z"

### Example Database Setup

#### Step 1: Create Database Properties
```
Title (Title) - Required
Slug (Rich Text) - Optional
Date (Date) - Optional  
Tags (Multi-select) - Optional
Categories (Multi-select) - Optional
Summary (Rich Text) - Optional
Status (Status) - Optional
Type (Select) - Optional
```

#### Step 2: Configure Status Options
```
📝 Draft (will set draft: true)
✅ Published
🔄 In Review  
⏸️ Archived
```

#### Step 3: Configure Type Options
```
posts (default if empty)
pages
...
```

#### Step 4: Configure Tag Options (Examples)
```
hugo, markdown, tutorial, guide, 
getting-started, advanced, tips, 
deployment, themes, shortcodes
```

#### Step 5: Configure Category Options (Examples)
```
development, web-development, tutorials,
documentation, guides, blog, news,
frontend, backend, devops
```

### Generated Front Matter Example

With the following Notion properties:
- **Title**: "Getting Started with Hugo"
- **Slug**: "hugo-getting-started"  
- **Date**: "2025-01-15"
- **Tags**: ["hugo", "tutorial"]
- **Categories**: ["development", "web-development"]
- **Summary**: "A comprehensive guide to getting started with Hugo"
- **Status**: "Published"
- **Type**: "posts"

The action generates this front matter:
```yaml
---
title: Getting Started with Hugo
slug: hugo-getting-started
date: 2025-01-15T00:00:00Z
lastmod: 2025-01-15T10:30:00Z
tags:
  - hugo
  - tutorial
categories:
  - development
  - web-development
summary: A comprehensive guide to getting started with Hugo
type: posts
---
```

> **Note**: The `draft` field is only included when set to `true`. When the status is not "Draft", the field is omitted from the front matter (equivalent to `draft: false`).

### File Path Generation

Based on the `Type` property:

| Type Value | Generated Path | Example |
|------------|----------------|---------|
| `posts` or empty | `content/posts/slug/index.md` | `content/posts/hugo-getting-started/index.md` |
| `pages` | `content/slug/index.md` | `content/about/index.md` |
| `docs` | `content/docs/slug/index.md` | `content/docs/installation/index.md` |
| `blog` | `content/blog/slug/index.md` | `content/blog/my-story/index.md` |

### Database Sharing Setup

1. **Create Integration**: Go to [Notion Developers](https://www.notion.so/my-integrations)
2. **Copy Token**: Save the integration token securely
3. **Share Database**: 
   - Open your database
   - Click "Share" → "Invite" 
   - Add your integration
   - Grant "Read" permission
4. **Get Database ID**: Copy from the database URL
   ```
   https://notion.so/workspace/DATABASE_ID?v=...
   ```

## 🔧 Advanced Usage

### Environment Variables

You can also use environment variables instead of action inputs:

```yaml
- name: Fetch content from Notion
  uses: ManassehZhou/notion-to-hugo@v1
  env:
    NOTION_TOKEN: ${{ secrets.NOTION_TOKEN }}
    NOTION_DATABASE_ID: ${{ secrets.NOTION_DATABASE_ID }}
  with:
    output-directory: 'content/posts'
```

### Multiple Databases

To fetch from multiple Notion databases:

```yaml
- name: Fetch blog posts
  uses: ManassehZhou/notion-to-hugo@v1
  with:
    notion-token: ${{ secrets.NOTION_TOKEN }}
    notion-database-id: ${{ secrets.BLOG_DATABASE_ID }}
    output-directory: 'content/posts'

- name: Fetch documentation
  uses: ManassehZhou/notion-to-hugo@v1
  with:
    notion-token: ${{ secrets.NOTION_TOKEN }}
    notion-database-id: ${{ secrets.DOCS_DATABASE_ID }}
    output-directory: 'content/docs'
```

### Custom Post-Processing

Add custom scripts after content generation:

```yaml
- name: Fetch content from Notion
  uses: ManassehZhou/notion-to-hugo@v1
  with:
    notion-token: ${{ secrets.NOTION_TOKEN }}
    notion-database-id: ${{ secrets.NOTION_DATABASE_ID }}

- name: Post-process content
  run: |
    # Add custom processing here
    find content -name "*.md" -exec sed -i 's/old-text/new-text/g' {} \;
```

## 🔧 Troubleshooting

### Common Issues

#### Files Not Generated
```yaml
- name: Check output directory
  run: |
    ls -la ${{ steps.notion-sync.outputs.output-path }}
    echo "Files generated: ${{ steps.notion-sync.outputs.files-generated }}"
```

#### Large File Handling
```yaml
- name: Handle large files
  run: |
    # Check for large files (>50MB)
    find content -type f -size +50M -ls
    # Compress large images
    find content -name "*.png" -o -name "*.jpg" | xargs -I {} sh -c 'echo "Processing: {}"'
```

#### Memory Issues
```yaml
- name: Fetch content with retry
  uses: ManassehZhou/notion-to-hugo@v1
  with:
    notion-token: ${{ secrets.NOTION_TOKEN }}
    notion-database-id: ${{ secrets.NOTION_DATABASE_ID }}
  continue-on-error: true
  
- name: Retry on failure
  if: failure()
  run: echo "Action failed, check for memory or rate limit issues"
```

### Debug Mode
Enable verbose logging for troubleshooting:

```yaml
- name: Debug Notion sync
  uses: ManassehZhou/notion-to-hugo@v1
  with:
    notion-token: ${{ secrets.NOTION_TOKEN }}
    notion-database-id: ${{ secrets.NOTION_DATABASE_ID }}
  env:
    VERBOSE: "true"
```

## 🔒 Security

- Store your Notion token in GitHub Secrets, never in your code
- Limit your Notion integration permissions to only the databases you need
- Regularly rotate your Notion integration tokens
- Use the principle of least privilege for GitHub Action permissions

## ⚠️ Limitations and Considerations

### File Size and Storage
- **Large attachments**: Very large files may cause performance issues or timeouts
- **Total content**: Be mindful of repository size limits when generating many files
- **GitHub Actions**: Action outputs are limited to 1MB; large file lists may be truncated

### Performance
- **API rate limits**: Notion API has rate limits that may affect large databases
- **Processing time**: Large databases with many pages/blocks may take longer to process
- **Memory usage**: Very large pages with many blocks may consume significant memory

### Recommendations
- **Start small**: Begin with a small subset of your database to test the process
- **Monitor sizes**: Regularly check generated content size
- **Optimize images**: Consider compressing images before uploading to Notion
- **Test incrementally**: Gradually increase the number of pages being processed

## 🤝 Contributing

We welcome contributions! Please see [CONTRIBUTING.md](CONTRIBUTING.md) for details.

### Development Setup

1. Clone the repository
2. Install Go 1.21 or later
3. Install dependencies: `go mod download`
4. Run tests: `go test ./...`
5. Build: `go build -o notion-to-hugo main.go`

## 📄 License

This project is licensed under the Apache License 2.0 - see the [LICENSE](LICENSE) file for details.

## 🐛 Issues and Support

- 🐛 [Report a bug](https://github.com/ManassehZhou/notion-to-hugo/issues/new?template=bug_report.md)
- 💡 [Request a feature](https://github.com/ManassehZhou/notion-to-hugo/issues/new?template=feature_request.md)
- 💬 [Ask a question](https://github.com/ManassehZhou/notion-to-hugo/discussions)

## 🌟 Acknowledgments

- Built with [jomei/notionapi](https://github.com/jomei/notionapi) for Notion API integration
- Inspired by the static site generator community
- Thanks to all contributors and users

## 📚 Related Projects

- [Hugo](https://gohugo.io/) - The world's fastest framework for building websites
- [Hexo](https://hexo.io/) - A fast, simple & powerful blog framework
- [Jekyll](https://jekyllrb.com/) - Transform your plain text into static websites and blogs
- [Notion API](https://developers.notion.com/) - Official Notion API documentation
- [GitHub Actions](https://github.com/features/actions) - CI/CD platform

---

Made with ❤️ by [ManassehZhou](https://github.com/ManassehZhou)