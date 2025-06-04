#!/bin/bash

# Import an existing template version using template_id/version_id format
# Replace 'd-template-id' with your template ID and 'version-id' with version ID
terraform import sendgrid_template_version.welcome_v1 d-template-id/version-id

# Example with actual IDs:
# terraform import sendgrid_template_version.welcome_v1 d-123456789/v-987654321

# You can find template and version IDs in the SendGrid dashboard under Email API > Dynamic Templates
# Or use the SendGrid API to list templates and their versions
