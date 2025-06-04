#!/bin/bash

# Import an existing subuser using its username
terraform import sendgrid_subuser.app_subuser app-emails

# You can find subuser usernames in the SendGrid dashboard under Settings > Subuser Management
