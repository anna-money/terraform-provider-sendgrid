terraform {
  required_providers {
    sendgrid = {
      source  = "arslanbekov/sendgrid"
      version = "~> 2.0"
    }
  }
}

provider "sendgrid" {
  api_key = var.sendgrid_api_key
}

variable "sendgrid_api_key" {
  description = "SendGrid API Key"
  type        = string
  sensitive   = true
}

# Example: Basic teammate with minimal permissions
resource "sendgrid_teammate" "developer" {
  email    = "developer@example.com"
  is_admin = false
  is_sso   = false
  scopes = [
    "mail.send",
    "templates.read"
  ]

  # Recommended: Increase timeouts for better reliability
  timeouts {
    create = "30m"
    update = "30m"
    delete = "30m"
  }
}

# Example: Admin user (no scopes needed)
resource "sendgrid_teammate" "admin" {
  email    = "admin@example.com"
  is_admin = true
  is_sso   = false
  # Note: scopes are ignored for admin users

  timeouts {
    create = "30m"
    update = "30m"
    delete = "30m"
  }
}

# Example: Marketing team member with extended permissions
resource "sendgrid_teammate" "marketing" {
  email    = "marketing@example.com"
  is_admin = false
  is_sso   = false
  scopes = [
    "mail.send",
    "marketing.send",
    "marketing.campaigns",
    "marketing.contacts",
    "templates.read",
    "templates.create",
    "stats.read"
  ]

  timeouts {
    create = "30m"
    update = "30m"
    delete = "30m"
  }
}

# Example: SSO user with first and last name
resource "sendgrid_teammate" "sso_user" {
  email      = "sso.user@example.com"
  first_name = "SSO"
  last_name  = "User"
  is_admin   = false
  is_sso     = true
  scopes = [
    "mail.send",
    "stats.read"
  ]

  timeouts {
    create = "30m"
    update = "30m"
    delete = "30m"
  }
}

# Example: Multiple teammates using for_each
locals {
  developers = {
    "dev1" = {
      email  = "dev1@example.com"
      scopes = ["mail.send", "templates.read"]
    }
    "dev2" = {
      email  = "dev2@example.com"
      scopes = ["mail.send", "templates.read", "stats.read"]
    }
    "dev3" = {
      email  = "dev3@example.com"
      scopes = ["mail.send"]
    }
  }
}

resource "sendgrid_teammate" "developers" {
  for_each = local.developers

  email    = each.value.email
  is_admin = false
  is_sso   = false
  scopes   = each.value.scopes

  # Important: Use higher timeouts when creating multiple resources
  timeouts {
    create = "45m"
    update = "45m"
    delete = "45m"
  }
}

# Example: Read-only user for monitoring/reporting
resource "sendgrid_teammate" "readonly" {
  email    = "readonly@example.com"
  is_admin = false
  is_sso   = false
  scopes = [
    "stats.read",
    "stats.global",
    "templates.read"
  ]

  timeouts {
    create = "30m"
    update = "30m"
    delete = "30m"
  }
}

# Outputs for reference
output "teammate_emails" {
  description = "List of all teammate emails"
  value = [
    sendgrid_teammate.developer.email,
    sendgrid_teammate.admin.email,
    sendgrid_teammate.marketing.email,
    sendgrid_teammate.sso_user.email,
    sendgrid_teammate.readonly.email
  ]
}

output "developer_team_emails" {
  description = "List of developer team emails"
  value       = [for teammate in sendgrid_teammate.developers : teammate.email]
}

# Best Practices Notes:
#
# 1. ALWAYS use timeouts for teammate resources to handle rate limiting
# 2. NEVER include automatic scopes like '2fa_exempt' or '2fa_required'
# 3. Use 'terraform plan' before 'terraform apply' to catch validation errors
# 4. For bulk operations, use reduced parallelism: terraform apply -parallelism=2
# 5. Admin users don't need scopes - they have full access automatically
# 6. Check your SendGrid plan supports the scopes you're trying to assign
# 7. Marketing scopes (marketing.*) require a Marketing plan
# 8. Advanced scopes (teammates.*, subusers.*) require Pro+ plan
