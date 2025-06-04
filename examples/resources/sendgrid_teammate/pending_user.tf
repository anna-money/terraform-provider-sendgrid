# Example: Non-SSO teammate that will be in pending status
resource "sendgrid_teammate" "pending_user" {
  email    = "new.teammate@example.com"
  is_admin = false
  is_sso   = false
  scopes = [
    "mail.send",
    "templates.read"
  ]

  # Recommended: Use timeouts for better reliability
  timeouts {
    create = "30m"
    update = "30m"
    delete = "30m"
  }
}

# Output the user status to see if they're pending or active
output "teammate_status" {
  value       = sendgrid_teammate.pending_user.user_status
  description = "Status of the teammate: 'pending' means invitation sent but not accepted, 'active' means user has accepted invitation"
}

# You can still manage pending users - updates will work
resource "sendgrid_teammate" "pending_user_with_more_scopes" {
  email    = "another.teammate@example.com"
  is_admin = false
  is_sso   = false
  scopes = [
    "mail.send",
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
