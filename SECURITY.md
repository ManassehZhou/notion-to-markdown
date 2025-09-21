# Security Policy

## Supported Versions

We actively support the following versions of the Notion to Hugo Action:

| Version | Supported          |
| ------- | ------------------ |
| 1.x.x   | :white_check_mark: |
| < 1.0   | :x:                |

## Reporting a Vulnerability

We take security vulnerabilities seriously. If you discover a security vulnerability in the Notion to Hugo Action, please follow these steps:

### 🔒 Private Disclosure

**DO NOT** create a public GitHub issue for security vulnerabilities.

Instead, please report security vulnerabilities by emailing:
- **Email**: security@example.com
- **Subject**: [SECURITY] Notion to Hugo Action Vulnerability

### 📧 What to Include

Please include the following information in your report:

1. **Description**: Clear description of the vulnerability
2. **Impact**: What could an attacker accomplish with this vulnerability?
3. **Reproduction**: Step-by-step instructions to reproduce the issue
4. **Environment**: Version of the action, operating system, etc.
5. **Mitigation**: Any temporary workarounds you've discovered

### 🔍 Example Report

```
Subject: [SECURITY] Notion to Hugo Action Vulnerability

Description:
The action is vulnerable to [specific vulnerability type] which could allow [impact].

Impact:
An attacker could [describe what they could achieve].

Reproduction:
1. Set up action with [specific configuration]
2. [Step-by-step instructions]
3. Observe [vulnerability manifestation]

Environment:
- Action Version: v1.0.0
- Runner OS: ubuntu-latest
- Hugo Version: 0.120.0

Mitigation:
[Any temporary workarounds]
```

### ⏱️ Response Timeline

- **Acknowledgment**: We will acknowledge receipt of your report within 48 hours
- **Initial Assessment**: We will provide an initial assessment within 5 business days
- **Updates**: We will provide regular updates on our progress
- **Resolution**: We aim to resolve critical vulnerabilities within 30 days

### 🛡️ Security Measures

#### For Users

1. **Keep Updated**: Always use the latest version of the action
2. **Secure Secrets**: Store Notion tokens and other sensitive data in GitHub Secrets
3. **Principle of Least Privilege**: Grant minimal necessary permissions to your Notion integration
4. **Regular Audits**: Regularly review your Notion integrations and access logs

#### For the Action

1. **Input Validation**: All user inputs are validated and sanitized
2. **Secure API Calls**: All Notion API calls use HTTPS
3. **No Sensitive Data Logging**: Tokens and sensitive data are never logged
4. **Container Security**: Docker images are regularly updated and scanned

### 🚨 Common Security Considerations

#### Notion Token Security
- Never commit your Notion integration token to version control
- Use GitHub Secrets to store your token
- Regularly rotate your Notion integration tokens
- Monitor Notion integration usage logs

#### Action Permissions
- Use minimal necessary permissions for your GitHub Actions workflow
- Review what data your Notion integration has access to
- Consider using separate integrations for different purposes

#### Content Security
- Be aware that public repositories will expose generated content
- Consider using private repositories for sensitive content
- Review generated markdown files for any unintended data exposure

### 📋 Security Checklist for Users

- [ ] Notion token stored in GitHub Secrets
- [ ] Notion integration has minimal required permissions
- [ ] Repository permissions are appropriately configured
- [ ] Generated content reviewed for sensitive information
- [ ] Action version pinned to specific tag (not `@main`)
- [ ] Workflow permissions follow principle of least privilege

### 🔗 Additional Resources

- [GitHub Security Documentation](https://docs.github.com/en/actions/security-guides)
- [Notion Security Documentation](https://developers.notion.com/docs/security)
- [Hugo Security Considerations](https://gohugo.io/about/security-model/)

### 📞 Contact Information

For non-security related issues, please use:
- GitHub Issues: https://github.com/ManassehZhou/notion-to-hugo/issues
- GitHub Discussions: https://github.com/ManassehZhou/notion-to-hugo/discussions

---

**Thank you for helping keep the Notion to Hugo Action secure!** 🛡️