# SendGrid Provider Authentication Examples

terraform {
  required_providers {
    sendgrid = {
      source  = "arslanbekov/sendgrid"
      version = "~> 1.1"
    }
  }
}

# =============================================================================
# METHOD 1: Environment Variable (RECOMMENDED)
# =============================================================================
# This is the most secure method for production environments

# Set environment variable before running terraform:
# export SENDGRID_API_KEY="SG.your-actual-api-key-here"

provider "sendgrid" {
  # No api_key specified - will automatically use SENDGRID_API_KEY env var
}

# =============================================================================
# METHOD 2: Terraform Variables (FLEXIBLE)
# =============================================================================
# Good for different environments (dev, staging, prod)

variable "sendgrid_api_key" {
  description = "SendGrid API Key for sending emails"
  type        = string
  sensitive   = true
  # No default value - must be provided
}

provider "sendgrid" {
  api_key = var.sendgrid_api_key
}

# Usage examples:
# terraform apply -var="sendgrid_api_key=SG.your-key-here"
# Or create terraform.tfvars:
# sendgrid_api_key = "SG.your-key-here"

# =============================================================================
# METHOD 3: terraform.tfvars file (CONVENIENT)
# =============================================================================
# Create a terraform.tfvars file with:
# sendgrid_api_key = "SG.your-actual-api-key-here"

# Then use the variable in provider:
# provider "sendgrid" {
#   api_key = var.sendgrid_api_key
# }

# =============================================================================
# METHOD 4: Different providers for different environments
# =============================================================================

variable "environment" {
  description = "Environment name"
  type        = string
  default     = "dev"
}

variable "sendgrid_api_key_dev" {
  description = "SendGrid API Key for development"
  type        = string
  sensitive   = true
  default     = ""
}

variable "sendgrid_api_key_prod" {
  description = "SendGrid API Key for production"
  type        = string
  sensitive   = true
  default     = ""
}

provider "sendgrid" {
  alias   = "dev"
  api_key = var.sendgrid_api_key_dev
}

provider "sendgrid" {
  alias   = "prod"
  api_key = var.sendgrid_api_key_prod
}

# Usage with aliases:
# resource "sendgrid_api_key" "dev_key" {
#   provider = sendgrid.dev
#   name     = "development-key"
#   scopes   = ["mail.send"]
# }

# =============================================================================
# SECURITY BEST PRACTICES
# =============================================================================

# ✅ DO:
# - Use environment variables in CI/CD pipelines
# - Use terraform.tfvars for local development (add to .gitignore)
# - Mark variables as sensitive = true
# - Use different API keys for different environments

# ❌ DON'T:
# - Hardcode API keys in .tf files
# - Commit API keys to version control
# - Use production keys in development
# - Share API keys in plain text

# =============================================================================
# EXAMPLE: Complete setup with environment variables
# =============================================================================

# 1. Set environment variable:
# export SENDGRID_API_KEY="SG.your-sendgrid-api-key-here"

# 2. Simple provider configuration:
# provider "sendgrid" {}

# 3. Create resources:
# resource "sendgrid_teammate" "example" {
#   email    = "teammate@example.com"
#   is_admin = false
#   is_sso   = false
#   scopes   = ["mail.send"]
# }
