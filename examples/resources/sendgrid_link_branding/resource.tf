# Basic link branding for click tracking
resource "sendgrid_link_branding" "main" {
  domain     = "links.mycompany.com"
  subdomain  = "click"
  is_default = true
}

# Marketing-specific link branding
resource "sendgrid_link_branding" "marketing" {
  domain     = "marketing.mycompany.com"
  subdomain  = "track"
  is_default = false
}
