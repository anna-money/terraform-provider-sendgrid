# Authentication Guide

This guide covers all authentication methods for the SendGrid Terraform provider with security best practices.

## Quick Start

The fastest way to get started:

```bash
export SENDGRID_API_KEY="SG.your-api-key-here"
```

```hcl
provider "sendgrid" {}
```

## Authentication Methods

### Method 1: Environment Variable (Recommended)

**Best for:** Production, CI/CD, security-conscious environments

```bash
# Set the environment variable
export SENDGRID_API_KEY="SG.your-actual-sendgrid-api-key-here"
```

```hcl
terraform {
  required_providers {
    sendgrid = {
      source  = "arslanbekov/sendgrid"
      version = "~> 1.1"
    }
  }
}

provider "sendgrid" {
  # API key automatically read from SENDGRID_API_KEY environment variable
}
```

**Advantages:**

- API key never stored in code
- Works with CI/CD systems
- Easy to rotate keys
- No risk of accidental commits

### Method 2: Terraform Variables

**Best for:** Multi-environment setups, flexible configurations

```hcl
variable "sendgrid_api_key" {
  description = "SendGrid API Key for email operations"
  type        = string
  sensitive   = true
  # No default value - must be provided
}

provider "sendgrid" {
  api_key = var.sendgrid_api_key
}
```

**Usage options:**

```bash
# Option A: Command line
terraform apply -var="sendgrid_api_key=SG.your-key-here"

# Option B: terraform.tfvars file (add to .gitignore!)
echo 'sendgrid_api_key = "SG.your-key-here"' > terraform.tfvars

# Option C: Environment variable for Terraform
export TF_VAR_sendgrid_api_key="SG.your-key-here"

# Option D: Interactive prompt
terraform apply  # Will prompt for the variable
```

### Method 3: Multiple Environments

**Best for:** Organizations with dev/staging/prod environments

```hcl
variable "environment" {
  description = "Environment name (dev/staging/prod)"
  type        = string
  default     = "dev"
}

variable "sendgrid_api_keys" {
  description = "SendGrid API keys per environment"
  type        = map(string)
  sensitive   = true
  default = {
    dev     = ""  # Set via terraform.tfvars or env vars
    staging = ""
    prod    = ""
  }
}

provider "sendgrid" {
  api_key = var.sendgrid_api_keys[var.environment]
}
```

**terraform.tfvars example:**

```hcl
environment = "prod"
sendgrid_api_keys = {
  dev     = "SG.development-key-here"
  staging = "SG.staging-key-here"
  prod    = "SG.production-key-here"
}
```

### Method 4: Provider Aliases

**Best for:** Managing multiple SendGrid accounts

```hcl
# Production account
provider "sendgrid" {
  alias   = "prod"
  api_key = var.sendgrid_prod_key
}

# Development account
provider "sendgrid" {
  alias   = "dev"
  api_key = var.sendgrid_dev_key
}

# Use with resources
resource "sendgrid_teammate" "prod_user" {
  provider = sendgrid.prod
  email    = "user@company.com"
  is_admin = false
  scopes   = ["mail.send"]
}

resource "sendgrid_teammate" "dev_user" {
  provider = sendgrid.dev
  email    = "dev@company.com"
  is_admin = true
}
```

## Getting Your API Key

1. **Log in to SendGrid**: <https://app.sendgrid.com/>
2. **Navigate to Settings** → **API Keys**
3. **Create API Key**:
   - **Name**: Give it a descriptive name (e.g., "Terraform Production")
   - **Permissions**: Select specific scopes (see [Required Scopes](#required-scopes))
4. **Copy the key**: Save it securely (you won't see it again!)

## Required Scopes

### Minimum Scopes for Basic Operations

```bash
# Teammate management
teammates.create
teammates.read
teammates.update
teammates.delete

# API key management
api_keys.create
api_keys.read
api_keys.update
api_keys.delete
```

### Recommended Scopes for Full Functionality

```bash
# Templates
templates.create
templates.read
templates.update
templates.delete
templates.versions.create
templates.versions.read
templates.versions.update
templates.versions.delete

# Domain management
whitelabel.create
whitelabel.read
whitelabel.update
whitelabel.delete

# Webhooks
user.webhooks.event.settings.create
user.webhooks.event.settings.read
user.webhooks.event.settings.update
user.webhooks.event.settings.delete

# Unsubscribe groups
asm.groups.create
asm.groups.read
asm.groups.update
asm.groups.delete

# SSO (if needed)
sso.settings.create
sso.settings.read
sso.settings.update
sso.settings.delete
```

### Full Access (Not Recommended)

```bash
# Only use if you need everything
mail.send
```

## Security Best Practices

### DO

- **Use environment variables** in production
- **Mark variables as sensitive** in Terraform
- **Add terraform.tfvars to .gitignore**
- **Use different API keys** for different environments
- **Rotate API keys regularly** (every 90 days)
- **Use minimal required scopes** for each key
- **Store keys in secure vaults** (AWS Secrets Manager, HashiCorp Vault, etc.)

### DON'T

- **Hardcode API keys** in `.tf` files
- **Commit API keys** to version control
- **Use production keys** in development
- **Share API keys** in plain text (Slack, email, etc.)
- **Use overly broad scopes** (like full access)
- **Leave old keys active** after rotation

## CI/CD Integration

### GitHub Actions

```yaml
name: Terraform
on: [push]

jobs:
  terraform:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3

      - name: Setup Terraform
        uses: hashicorp/setup-terraform@v2

      - name: Terraform Init
        run: terraform init
        env:
          SENDGRID_API_KEY: ${{ secrets.SENDGRID_API_KEY }}

      - name: Terraform Plan
        run: terraform plan
        env:
          SENDGRID_API_KEY: ${{ secrets.SENDGRID_API_KEY }}
```

### GitLab CI

```yaml
terraform:
  image: hashicorp/terraform:latest
  variables:
    SENDGRID_API_KEY: $SENDGRID_API_KEY
  script:
    - terraform init
    - terraform plan
    - terraform apply -auto-approve
```

### Jenkins

```groovy
pipeline {
    agent any
    environment {
        SENDGRID_API_KEY = credentials('sendgrid-api-key')
    }
    stages {
        stage('Terraform') {
            steps {
                sh 'terraform init'
                sh 'terraform plan'
                sh 'terraform apply -auto-approve'
            }
        }
    }
}
```

## Docker Integration

```dockerfile
FROM hashicorp/terraform:latest

# Set environment variable
ENV SENDGRID_API_KEY=""

WORKDIR /workspace
COPY . .

# Run terraform
RUN terraform init
CMD ["terraform", "apply", "-auto-approve"]
```

```bash
# Run with API key
docker run -e SENDGRID_API_KEY="SG.your-key" your-terraform-image
```

## Advanced Configuration

### Custom API Endpoint

```hcl
provider "sendgrid" {
  api_key = var.sendgrid_api_key
  host    = "https://api.sendgrid.com"  # Default
}
```

### Timeout Configuration

```hcl
provider "sendgrid" {
  api_key = var.sendgrid_api_key

  # Custom timeouts for all operations
  timeout = "30s"
}
```

## Testing Your Configuration

### Validate API Key

```bash
# Test API key manually
curl -X GET \
  https://api.sendgrid.com/v3/user/profile \
  -H "Authorization: Bearer SG.your-api-key"
```

### Terraform Validation

```bash
# Validate configuration
terraform validate

# Test provider connection
terraform plan
```

### Quick Test Resource

```hcl
# Create a test API key to verify connection
resource "sendgrid_api_key" "test" {
  name   = "terraform-test-key"
  scopes = ["mail.send"]
}

output "test_key_id" {
  value = sendgrid_api_key.test.id
}
```

## Troubleshooting

### Invalid API Key

```bash
Error: request failed: api response: HTTP 401: {"errors":[{"message":"The provided authorization grant is invalid, expired, or revoked"}]}
```

**Solutions:**

1. Verify API key is correct
2. Check key hasn't expired
3. Ensure key has required scopes

### Insufficient Permissions

```bash
Error: request failed: api response: HTTP 403: {"errors":[{"message":"access forbidden"}]}
```

**Solutions:**

1. Add required scopes to API key
2. Check SendGrid plan limitations
3. Verify account permissions

### Environment Variable Not Set

```bash
Error: Missing required argument: "api_key"
```

**Solutions:**

```bash
# Check if variable is set
echo $SENDGRID_API_KEY

# Set the variable
export SENDGRID_API_KEY="SG.your-key"

# Or use terraform variable
terraform apply -var="sendgrid_api_key=SG.your-key"
```

## Next Steps

After setting up authentication:

1. [Check Examples](EXAMPLES.md)
2. [Browse Resources](RESOURCES.md)
3. [Troubleshooting Guide](TROUBLESHOOTING.md)

## Need Help?

- [Report Authentication Issues](https://github.com/arslanbekov/terraform-provider-sendgrid/issues)
- [Community Discussions](https://github.com/arslanbekov/terraform-provider-sendgrid/discussions)
- [SendGrid API Documentation](https://docs.sendgrid.com/api-reference)
