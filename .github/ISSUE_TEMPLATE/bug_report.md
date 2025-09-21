---
name: Bug Report
about: Create a report to help us improve
title: '[BUG] '
labels: ['bug', 'triage']
assignees: ''
---

## 🐛 Bug Description

A clear and concise description of what the bug is.

## 🔄 Steps to Reproduce

Steps to reproduce the behavior:

1. Set up Notion database with '...'
2. Configure action with '...'
3. Run workflow '...'
4. See error

## ✅ Expected Behavior

A clear and concise description of what you expected to happen.

## ❌ Actual Behavior

A clear and concise description of what actually happened.

## 📸 Screenshots

If applicable, add screenshots to help explain your problem.

## 🔧 Configuration

**Action Configuration:**
```yaml
- uses: ManassehZhou/notion-to-markdown@v1
  with:
    notion-token: ${{ secrets.NOTION_TOKEN }}
    notion-database-id: ${{ secrets.NOTION_DATABASE_ID }}
    # Add your full configuration here
```

**Notion Database Schema:**
- Describe your database properties and structure

## 🌍 Environment

- OS: [e.g. ubuntu-latest, windows-latest]
- Action Version: [e.g. v1.0.0]
- Hugo Version: [e.g. 0.120.0]
- Go Version: [if running locally, e.g. 1.21]

## 📋 Additional Context

Add any other context about the problem here.

## 🔍 Error Logs

<details>
<summary>Click to expand error logs</summary>

```
Paste your error logs here
```

</details>

## ✔️ Checklist

- [ ] I have searched existing issues for duplicates
- [ ] I have provided all the requested information
- [ ] I have tested with the latest version of the action
- [ ] I have included relevant configuration and error logs