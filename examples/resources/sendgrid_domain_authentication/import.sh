#!/bin/bash

# Import domain authentication using its ID
terraform import sendgrid_domain_authentication.main 12345

# Find domain authentication IDs in SendGrid dashboard under Settings > Sender Authentication
