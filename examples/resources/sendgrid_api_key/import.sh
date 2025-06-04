#!/bin/bash

# Import an existing API key using its ID
# Replace 'SG.example_api_key_id' with your actual API key ID
terraform import sendgrid_api_key.basic SG.example_api_key_id

# You can find API key IDs in the SendGrid dashboard under Settings > API Keys
# Or use the SendGrid API to list existing keys
