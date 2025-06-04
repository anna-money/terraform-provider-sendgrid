# Selective event webhook - only track delivery and engagement issues
resource "sendgrid_event_webhook" "selective" {
  enabled = true
  url     = "https://monitoring.myapp.com/email-issues"

  # Only track problematic events
  bounce      = true
  dropped     = true
  spam_report = true
  deferred    = true

  # Don't track successful events to reduce noise
  delivered         = false
  processed         = false
  open              = false
  click             = false
  unsubscribe       = false
  group_resubscribe = false
  group_unsubscribe = false
}
