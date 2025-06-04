#!/bin/bash

# Import an existing parse webhook using its hostname
terraform import sendgrid_parse_webhook.inbound inbound.myapp.com

# You can find parse webhooks in the SendGrid dashboard under Settings > Inbound Parse
