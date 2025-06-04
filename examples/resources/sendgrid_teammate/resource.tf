resource "sendgrid_teammate" "example" {
  email    = "teammate@example.com"
  is_admin = false
  is_sso   = false
  scopes = [
    "mail.send",
    "templates.read"
  ]
}
