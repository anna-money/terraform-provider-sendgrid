# Read-only API key for monitoring and analytics
resource "sendgrid_api_key" "readonly" {
  name = "monitoring-readonly-key"
  scopes = [
    "stats.read",
    "stats.global.read",
    "templates.read",
    "templates.versions.read",
    "teammates.read",
    "subusers.read",
    "suppression.read",
    "user.profile.read",
    "user.account.read",
    "categories.read",
    "categories.stats.read",
    "mailbox_providers.stats.read"
  ]
}
