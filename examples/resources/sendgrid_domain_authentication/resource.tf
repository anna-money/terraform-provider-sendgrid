# Basic domain authentication setup
resource "sendgrid_domain_authentication" "main" {
  domain             = "mycompany.com"
  subdomain          = "em"
  is_default         = true
  automatic_security = true
  custom_spf         = false
}

# Domain with custom SPF record
resource "sendgrid_domain_authentication" "marketing" {
  domain             = "marketing.mycompany.com"
  subdomain          = "mail"
  is_default         = false
  automatic_security = false
  custom_spf         = true
}
