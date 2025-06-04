# Basic SSO integration
resource "sendgrid_sso_integration" "okta" {
  name                  = "Okta SSO"
  enabled               = true
  signin_url            = "https://dev-12345.okta.com/app/sendgrid/abcd1234/sso/saml"
  signout_url           = "https://dev-12345.okta.com/login/signout"
  entity_id             = "http://www.okta.com/abcd1234"
  completed_integration = true
}

# SSO integration for Azure AD
resource "sendgrid_sso_integration" "azure_ad" {
  name                  = "Azure Active Directory"
  enabled               = true
  signin_url            = "https://login.microsoftonline.com/tenant-id/saml2"
  signout_url           = "https://login.microsoftonline.com/tenant-id/saml2"
  entity_id             = "https://sts.windows.net/tenant-id/"
  completed_integration = false # Will be completed after certificate upload
}
