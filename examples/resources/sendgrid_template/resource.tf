# Basic dynamic template for transactional emails
resource "sendgrid_template" "welcome_email" {
  name       = "Welcome Email Template"
  generation = "dynamic"
}

# Legacy template (for backward compatibility)
resource "sendgrid_template" "legacy_newsletter" {
  name       = "Legacy Newsletter Template"
  generation = "legacy"
}
