# SSO teammate with first and last name
resource "sendgrid_teammate" "sso_user" {
  email      = "sso.user@example.com"
  first_name = "John"
  last_name  = "Doe"
  is_admin   = false
  is_sso     = true
  scopes = [
    "mail.send",
    "stats.read"
  ]
}
