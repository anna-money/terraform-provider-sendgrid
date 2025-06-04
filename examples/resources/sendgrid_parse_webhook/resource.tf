# Basic parse webhook for inbound email processing
resource "sendgrid_parse_webhook" "inbound" {
  hostname   = "inbound.myapp.com"
  url        = "https://api.myapp.com/email/parse"
  spam_check = true
  send_raw   = false
}

# Advanced parse webhook with raw content
resource "sendgrid_parse_webhook" "support" {
  hostname   = "support.myapp.com"
  url        = "https://api.myapp.com/support/parse"
  spam_check = false # Handle spam filtering in application
  send_raw   = true  # Need full MIME content for attachments
}
