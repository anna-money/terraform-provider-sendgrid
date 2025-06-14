---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "sendgrid_api_key Resource - sendgrid"
subcategory: ""
description: |-
  
---

# sendgrid_api_key (Resource)



## Example Usage

```terraform
# Basic API key with common permissions
resource "sendgrid_api_key" "basic" {
  name = "my-app-api-key"
  scopes = [
    "mail.send",
    "sender_verification_eligible"
  ]
}

# Output the API key for use in applications
output "api_key" {
  value     = sendgrid_api_key.basic.api_key
  sensitive = true
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `name` (String) The name you will use to describe this API Key.

### Optional

- `scopes` (Set of String) The individual permissions that you are giving to this API Key.

### Read-Only

- `api_key` (String, Sensitive) The API key created by the API.
- `id` (String) The ID of this resource.

## Import

Import is supported using the following syntax:

```shell
#!/bin/bash

# Import an existing API key using its ID
# Replace 'SG.example_api_key_id' with your actual API key ID
terraform import sendgrid_api_key.basic SG.example_api_key_id

# You can find API key IDs in the SendGrid dashboard under Settings > API Keys
# Or use the SendGrid API to list existing keys
```
