# Basic event webhook configuration
resource "sendgrid_event_webhook" "basic" {
  enabled = true
  url     = "https://api.myapp.com/sendgrid/events"

  # Basic email events
  delivered   = true
  bounce      = true
  dropped     = true
  spam_report = true
  unsubscribe = true

  # Engagement events (requires tracking enabled)
  open  = true
  click = true

  # Processing events
  processed = true
  deferred  = true

  # Group events (requires subscription tracking)
  group_resubscribe = true
  group_unsubscribe = true
}
