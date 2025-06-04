#!/bin/bash

# Import an existing unsubscribe group using its ID
# Replace '12345' with your actual unsubscribe group ID
terraform import sendgrid_unsubscribe_group.general 12345

# You can find unsubscribe group IDs in the SendGrid dashboard under Settings > Unsubscribe Groups
# Or use the SendGrid API to list existing groups
