# SSO certificate for Okta integration
resource "sendgrid_sso_certificate" "okta_cert" {
  integration_id     = sendgrid_sso_integration.okta.id
  public_certificate = <<EOF
-----BEGIN CERTIFICATE-----
MIIDpDCCAoygAwIBAgIGAV2ka+55MA0GCSqGSIb3DQEBCwUAMIGSMQswCQYDVQQG
EwJVUzETMBEGA1UECAwKQ2FsaWZvcm5pYTEWMBQGA1UEBwwNU2FuIEZyYW5jaXNj
bzENMAsGA1UECgwET2t0YTEUMBIGA1UECwwLU1NPUHJvdmlkZXIxEzARBgNVBAMM
Ck15T2t0YURvbWFpbjEcMBoGCSqGSIb3DQEJARYNaW5mb0Bva3RhLmNvbTAeFw0x
... (certificate content) ...
-----END CERTIFICATE-----
EOF
}

# SSO certificate from file
resource "sendgrid_sso_certificate" "azure_cert" {
  integration_id     = sendgrid_sso_integration.azure_ad.id
  public_certificate = file("${path.module}/azure-ad-certificate.pem")
}
