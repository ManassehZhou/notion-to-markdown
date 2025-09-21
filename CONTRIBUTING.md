# Contributing to Notion to Hugo Action

Thank you for your interest in contributing to the Notion to Hugo Action! We welcome contributions from everyone.

## 📋 Table of Contents

- [Code of Conduct](#code-of-conduct)
- [Getting Started](#getting-started)
- [Development Setup](#development-setup)
- [Making Changes](#making-changes)
- [Testing](#testing)
- [Submitting Changes](#submitting-changes)
- [Style Guide](#style-guide)
- [Release Process](#release-process)

## 📜 Code of Conduct

This project and everyone participating in it is governed by our Code of Conduct. By participating, you are expected to uphold this code. Please report unacceptable behavior to [manasseh.zhou@example.com](mailto:manasseh.zhou@example.com).

## 🚀 Getting Started

### Prerequisites

- Go 1.21 or later
- Docker (for testing the action locally)
- Git
- A Notion account with integration token
- Basic understanding of Hugo and GitHub Actions

### Development Setup

1. **Fork the repository** on GitHub
2. **Clone your fork** locally:
   ```bash
   git clone https://github.com/YOUR_USERNAME/notion-to-markdown.git
   cd notion-to-markdown
   ```

3. **Add the upstream remote**:
   ```bash
   git remote add upstream https://github.com/ManassehZhou/notion-to-markdown.git
   ```

4. **Install dependencies**:
   ```bash
   go mod download
   ```

5. **Create a feature branch**:
   ```bash
   git checkout -b feature/your-feature-name
   ```

## 🔧 Making Changes

### Project Structure

```
notion-to-markdown/
├── action.yml              # GitHub Action metadata
├── Dockerfile              # Container definition
├── entrypoint.sh           # Action entrypoint script
├── main.go                 # CLI application
├── go.mod                  # Go module definition
├── config/                 
│   └── notion-to-markdown.yaml # Default configuration
├── internal/               # Internal packages
│   ├── notionclient/       # Notion API client
│   ├── renderer/           # Markdown rendering
│   └── writer/             # File writing utilities
└── .github/                # GitHub templates and workflows
```

### Areas for Contribution

1. **New Hugo Themes**: Add support for additional Hugo themes
2. **Notion Block Types**: Extend support for more Notion block types
3. **Configuration Options**: Add new customization options
4. **Performance**: Optimize API calls and rendering speed
5. **Documentation**: Improve docs, examples, and guides
6. **Bug Fixes**: Fix reported issues

### Making Code Changes

1. **Follow Go conventions**: Use `gofmt`, `go vet`, and follow effective Go practices
2. **Use structured logging**: Use `log/slog` for all logging instead of `fmt.Print*` functions
3. **Add tests**: Include unit tests for new functionality
4. **Update documentation**: Update README.md and code comments
5. **Test locally**: Ensure your changes work with real Notion content

### Environment Variables for Testing

Create a `.env` file (not committed to git):

```bash
NOTION_TOKEN=your_notion_integration_token
NOTION_DATABASE_ID=your_database_id
```

## 🧪 Testing

### Unit Tests

Run the test suite:

```bash
go test ./...
```

Run tests with coverage:

```bash
go test -cover ./...
```

### Integration Testing

Test with a real Notion database:

```bash
go run main.go -token $NOTION_TOKEN -database $NOTION_DATABASE_ID -out test-output
```

### Docker Testing

Build and test the Docker image:

```bash
docker build -t notion-to-markdown .
docker run -e NOTION_TOKEN=$NOTION_TOKEN -e NOTION_DATABASE_ID=$NOTION_DATABASE_ID notion-to-markdown
```

### Testing GitHub Action Locally

Use [act](https://github.com/nektos/act) to test actions locally:

```bash
# Install act
brew install act

# Test the action
act -j test-action --secret NOTION_TOKEN=$NOTION_TOKEN --secret NOTION_DATABASE_ID=$NOTION_DATABASE_ID
```

## 📤 Submitting Changes

### Pull Request Process

1. **Update documentation** for any user-facing changes
2. **Add tests** for new functionality
3. **Ensure all tests pass**
4. **Update CHANGELOG.md** if applicable
5. **Create a pull request** with:
   - Clear title and description
   - Reference to related issues
   - Screenshots for UI changes
   - Testing instructions

### Pull Request Template

```markdown
## Description
Brief description of changes

## Type of Change
- [ ] Bug fix
- [ ] New feature
- [ ] Breaking change
- [ ] Documentation update

## Testing
- [ ] Unit tests pass
- [ ] Integration tests pass
- [ ] Tested with real Notion database

## Checklist
- [ ] Code follows style guidelines
- [ ] Self-review completed
- [ ] Documentation updated
- [ ] Tests added/updated
```

## 🎨 Style Guide

### Go Code Style

- Follow the [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)
- Use `gofmt` to format code
- Run `go vet` to check for issues
- Use meaningful variable and function names
- Add comments for exported functions and complex logic
- Use structured logging with `log/slog` instead of `fmt.Print*` functions
- Log levels: use `Info` for normal operation, `Debug` for verbose mode, `Error` for errors, `Warn` for warnings

### Commit Messages

Use conventional commit format:

```
type(scope): description

body (optional)

footer (optional)
```

Types:
- `feat`: New feature
- `fix`: Bug fix
- `docs`: Documentation changes
- `style`: Code style changes
- `refactor`: Code refactoring
- `test`: Test additions/changes
- `chore`: Maintenance tasks

Examples:
```
feat(renderer): add support for toggle blocks
fix(docker): resolve permission issues in container
docs(readme): update configuration examples
```

### Documentation Style

- Use clear, concise language
- Include code examples
- Add emojis for visual appeal (but don't overuse)
- Use proper markdown formatting
- Test all code examples

## 🚢 Release Process

### Versioning

We use [Semantic Versioning](https://semver.org/):
- **MAJOR**: Breaking changes
- **MINOR**: New features (backward compatible)
- **PATCH**: Bug fixes (backward compatible)

### Creating a Release

1. Update version in relevant files
2. Update CHANGELOG.md
3. Create a release tag
4. Publish release notes
5. Update marketplace listing

## 🆘 Getting Help

- 💬 [Start a discussion](https://github.com/ManassehZhou/notion-to-markdown/discussions)
- 🐛 [Report an issue](https://github.com/ManassehZhou/notion-to-markdown/issues)
- 📧 Email: manasseh.zhou@example.com

## 📚 Resources

- [Notion API Documentation](https://developers.notion.com/)
- [Hugo Documentation](https://gohugo.io/documentation/)
- [GitHub Actions Documentation](https://docs.github.com/en/actions)
- [Go Documentation](https://golang.org/doc/)

Thank you for contributing! 🎉