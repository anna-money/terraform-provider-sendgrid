# Marketing teammate with extended permissions
resource "sendgrid_teammate" "marketing" {
  email    = "marketing@example.com"
  is_admin = false
  is_sso   = false
  scopes = [
    "mail.send",
    "marketing.read",
    "marketing.automation.read",
    "templates.read",
    "templates.create",
    "templates.update",
    "stats.read"
  ]

  # Recommended: Use timeouts for better reliability
  timeouts {
    create = "30m"
    update = "30m"
    delete = "30m"
  }
}
