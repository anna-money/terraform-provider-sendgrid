# Default unsubscribe group for general emails
resource "sendgrid_unsubscribe_group" "general" {
  name        = "General Communications"
  description = "General company communications and updates"
  is_default  = true
}

# Marketing emails unsubscribe group
resource "sendgrid_unsubscribe_group" "marketing" {
  name        = "Marketing Emails"
  description = "Promotional offers, newsletters, and marketing content"
  is_default  = false
}

# Transactional emails unsubscribe group
resource "sendgrid_unsubscribe_group" "transactional" {
  name        = "Account Notifications"
  description = "Account-related notifications and alerts"
  is_default  = false
}
