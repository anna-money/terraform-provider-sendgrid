terraform {
  required_providers {
    sendgrid = {
      source  = "arslanbekov/sendgrid"
      version = "~> 2.0"
    }
  }
}

# Method 1: Using environment variable (recommended for security)
# Set: export SENDGRID_API_KEY="your-api-key-here"
provider "sendgrid" {
  # API key will be automatically read from SENDGRID_API_KEY environment variable
}

# Method 2: Using Terraform variable (for flexibility)
variable "sendgrid_api_key" {
  description = "SendGrid API Key"
  type        = string
  sensitive   = true
}

provider "sendgrid" {
  api_key = var.sendgrid_api_key
}

# Method 3: Direct configuration (NOT recommended for production)
# provider "sendgrid" {
#   api_key = "SG.your-api-key-here"  # Don't do this in production!
# }
