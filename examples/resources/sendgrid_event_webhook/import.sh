#!/bin/bash

# Import an existing event webhook using its ID
terraform import sendgrid_event_webhook.main webhook-id-12345

# You can find webhook IDs in the SendGrid dashboard under Settings > Mail Settings > Event Webhooks
