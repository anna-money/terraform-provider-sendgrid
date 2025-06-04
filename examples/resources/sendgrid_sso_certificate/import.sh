#!/bin/bash

# Import SSO certificate using its ID
terraform import sendgrid_sso_certificate.okta_cert cert-1234-5678

# Find certificate IDs in SendGrid dashboard under Settings > SSO > Certificates
