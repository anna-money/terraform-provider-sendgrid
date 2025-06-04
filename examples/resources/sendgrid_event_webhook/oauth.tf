# Event webhook with OAuth authentication
resource "sendgrid_event_webhook" "oauth" {
  enabled = true
  url     = "https://secure-api.myapp.com/sendgrid/events"

  # OAuth configuration for secure webhook
  oauth_client_id     = "your-oauth-client-id"
  oauth_client_secret = "your-oauth-client-secret" # Will be stored securely
  oauth_token_url     = "https://auth.myapp.com/oauth/token"

  # Event types to track
  delivered         = true
  bounce            = true
  dropped           = true
  spam_report       = true
  unsubscribe       = true
  open              = true
  click             = true
  processed         = true
  deferred          = true
  group_resubscribe = true
  group_unsubscribe = true
}
