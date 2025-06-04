# Basic API key with common permissions
resource "sendgrid_api_key" "basic" {
  name = "my-app-api-key"
  scopes = [
    "mail.send",
    "sender_verification_eligible"
  ]
}

# Output the API key for use in applications
output "api_key" {
  value     = sendgrid_api_key.basic.api_key
  sensitive = true
}
