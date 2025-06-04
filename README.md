# Terraform Provider for SendGrid

[![Build Status](https://github.com/arslanbekov/terraform-provider-sendgrid/workflows/Tests/badge.svg)](https://github.com/arslanbekov/terraform-provider-sendgrid/actions)
[![Go Report Card](https://goreportcard.com/badge/github.com/arslanbekov/terraform-provider-sendgrid)](https://goreportcard.com/report/github.com/arslanbekov/terraform-provider-sendgrid)
[![codecov](https://codecov.io/gh/arslanbekov/terraform-provider-sendgrid/branch/master/graph/badge.svg)](https://codecov.io/gh/arslanbekov/terraform-provider-sendgrid)

A comprehensive Terraform provider for managing SendGrid resources.

## Key Features

- **Advanced Rate Limiting** - Built-in exponential backoff and retry logic
- **Teammate Management** - Complete lifecycle management including pending invitations
- **Template Management** - Full template and version control
- **Multiple Auth Methods** - Environment variables, Terraform variables, and more
- **95% Test Coverage** - Production-ready with comprehensive testing
- **Rich Documentation** - Extensive examples and troubleshooting guides

## Quick Start

```bash
# Set your API key
export SENDGRID_API_KEY="SG.your-api-key-here"
```

```hcl
terraform {
  required_providers {
    sendgrid = {
      source  = "arslanbekov/sendgrid"
      version = "~> 2.0"
    }
  }
}

provider "sendgrid" {}

resource "sendgrid_teammate" "example" {
  email    = "teammate@example.com"
  is_admin = false
  is_sso   = false
  scopes   = ["mail.send"]
}
```

```bash
terraform init && terraform apply
```

## Documentation

| Topic                                      | Description                                            |
| ------------------------------------------ | ------------------------------------------------------ |
| [Installation](docs/INSTALLATION.md)       | Installation methods and requirements                  |
| [Authentication](docs/AUTHENTICATION.md)   | All authentication methods and security best practices |
| [Resources](docs/RESOURCES.md)             | Complete list of resources and data sources            |
| [Examples](docs/EXAMPLES.md)               | Practical usage examples and patterns                  |
| [Troubleshooting](docs/TROUBLESHOOTING.md) | Common issues and solutions                            |
| [Contributing](docs/CONTRIBUTING.md)       | Development setup and contribution guidelines          |

## Popular Use Cases

- **Team Management**: Invite and manage teammates with specific permissions
- **Email Templates**: Create and version email templates
- **API Key Management**: Secure API key creation with minimal scopes
- **Domain Setup**: Configure domain authentication and link branding
- **Webhook Configuration**: Set up event and parse webhooks

## Quick Links

- [Terraform Registry](https://registry.terraform.io/providers/arslanbekov/sendgrid)
- [Report Issues](https://github.com/arslanbekov/terraform-provider-sendgrid/issues)
- [Discussions](https://github.com/arslanbekov/terraform-provider-sendgrid/discussions)
- [Changelog](CHANGELOG.md)
- [Migration Guide](MIGRATION_GUIDE.md)

## License

This project is licensed under the [Mozilla Public License 2.0](LICENSE).

---

**Disclaimer:** This is an unofficial provider maintained by the community. While it offers enhanced features and comprehensive testing, evaluate thoroughly for production use.
