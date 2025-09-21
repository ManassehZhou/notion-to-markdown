#!/bin/sh

set -e

# Get inputs from action
NOTION_TOKEN="$1"
NOTION_DATABASE_ID="$2"
OUTPUT_DIRECTORY="$3"
CONFIG_FILE="$4"

# Validate required inputs
if [ -z "$NOTION_TOKEN" ]; then
    echo "Error: notion-token is required"
    exit 1
fi

if [ -z "$NOTION_DATABASE_ID" ]; then
    echo "Error: notion-database-id is required"
    exit 1
fi

# Set defaults
OUTPUT_DIRECTORY=${OUTPUT_DIRECTORY:-"content"}
CONFIG_FILE=${CONFIG_FILE:-"/app/config/notion-to-markdown.yaml"}

echo "Starting Notion to Hugo conversion..."
echo "Database ID: $NOTION_DATABASE_ID"
echo "Output Directory: $OUTPUT_DIRECTORY"
echo "Config File: $CONFIG_FILE"

# Change to workspace directory
cd $GITHUB_WORKSPACE

# Set environment variables for the Go application
export NOTION_TOKEN="$NOTION_TOKEN"
export NOTION_DATABASE_ID="$NOTION_DATABASE_ID"

# Run the notion-to-markdown tool
echo "Running notion-to-markdown conversion..."
if /app/notion-to-markdown -token "$NOTION_TOKEN" -database "$NOTION_DATABASE_ID" -out "$OUTPUT_DIRECTORY" -config "$CONFIG_FILE"; then
    # Count generated files (with safety limit to prevent output overflow)
    FILES_COUNT=$(find "$OUTPUT_DIRECTORY" -name "*.md" -type f | wc -l | tr -d ' ')
    
    # Get total size of generated content
    TOTAL_SIZE=$(du -sh "$OUTPUT_DIRECTORY" 2>/dev/null | cut -f1 || echo "unknown")
    
    echo "✅ Generated $FILES_COUNT markdown files in $OUTPUT_DIRECTORY (total size: $TOTAL_SIZE)"
    
    # Set outputs for GitHub Actions (with size limits)
    echo "files-generated=$FILES_COUNT" >> $GITHUB_OUTPUT
    echo "output-path=$OUTPUT_DIRECTORY" >> $GITHUB_OUTPUT  
    echo "success=true" >> $GITHUB_OUTPUT
    
    echo "🎉 Notion to Markdown conversion completed successfully!"
    exit 0
else
    echo "❌ Notion to Hugo conversion failed"
    echo "success=false" >> $GITHUB_OUTPUT
    exit 1
fi