# Admin teammate - no scopes needed
resource "sendgrid_teammate" "admin" {
  email    = "admin@example.com"
  is_admin = true
  is_sso   = false
}
