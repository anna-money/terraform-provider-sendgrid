#!/bin/bash

# Import an existing template using its ID
# Replace 'd-template-id-here' with your actual template ID
terraform import sendgrid_template.welcome_email d-template-id-here

# You can find template IDs in the SendGrid dashboard under Email API > Dynamic Templates
# Or use the SendGrid API to list existing templates
