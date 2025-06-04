# API key with extensive permissions for admin operations
resource "sendgrid_api_key" "admin" {
  name = "admin-api-key"
  scopes = [
    "api_keys.create",
    "api_keys.read",
    "api_keys.update",
    "api_keys.delete",
    "mail.send",
    "sender_verification_eligible",
    "stats.read",
    "stats.global.read",
    "templates.create",
    "templates.read",
    "templates.update",
    "templates.delete",
    "templates.versions.create",
    "templates.versions.read",
    "templates.versions.update",
    "templates.versions.delete",
    "templates.versions.activate.create",
    "teammates.read",
    "subusers.read",
    "suppression.read",
    "user.profile.read",
    "user.settings.enforced_tls.read"
  ]
}
