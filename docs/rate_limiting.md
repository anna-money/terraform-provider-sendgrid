# Rate Limiting Handling

The SendGrid Terraform Provider includes built-in rate limiting handling to manage HTTP 429 "too many requests" errors from the SendGrid API.

## How it works

All resource operations (create, read, update, delete) automatically retry when they encounter HTTP 429 rate limit errors using exponential backoff. The provider will:

1. Detect HTTP 429 responses from the SendGrid API
2. Wait and retry the request with exponential backoff
3. Continue retrying until the timeout is reached
4. Use the timeout configured in your Terraform resource

## Resources with Rate Limiting Support

All SendGrid resources now include rate limiting protection:

- **sendgrid_api_key** - API key management
- **sendgrid_teammate** - Teammate/user management
- **sendgrid_template** - Email template management
- **sendgrid_template_version** - Template version management
- **sendgrid_subuser** - Subuser management
- **sendgrid_unsubscribe_group** - Unsubscribe group management
- **sendgrid_domain_authentication** - Domain authentication
- **sendgrid_link_branding** - Link branding
- **sendgrid_parse_webhook** - Parse webhook configuration
- **sendgrid_event_webhook** - Event webhook configuration
- **sendgrid_sso_integration** - SSO integration
- **sendgrid_sso_certificate** - SSO certificate management

## Configuration

### Default Timeouts

By default, Terraform resources use these timeouts:

- **Create**: 20 minutes
- **Update**: 20 minutes
- **Delete**: 20 minutes

### Custom Timeouts

You can configure custom timeouts to allow for more retry attempts:

```hcl
resource "sendgrid_teammate" "example" {
  email    = "user@example.com"
  is_admin = false
  scopes   = ["mail.send"]

  timeouts {
    create = "30m"
    update = "30m"
    delete = "30m"
  }
}

resource "sendgrid_api_key" "example" {
  name   = "my-api-key"
  scopes = ["mail.send"]

  timeouts {
    create = "30m"
    update = "30m"
    delete = "30m"
  }
}
```

## Best Practices

### 1. Use Parallelism Limits

Reduce the number of concurrent operations to avoid hitting rate limits:

```bash
terraform apply -parallelism=2
```

For creating multiple API keys or teammates, consider even lower parallelism:

```bash
terraform apply -parallelism=1
```

### 2. Batch Operations

When creating multiple resources, consider spacing out the operations:

```hcl
# Instead of creating many resources at once, consider using for_each
# with a smaller set or using terraform apply with lower parallelism

resource "sendgrid_api_key" "keys" {
  for_each = {
    "key1" = ["mail.send"]
    "key2" = ["mail.send", "marketing.send"]
    # Limit to smaller batches when applying
  }

  name   = each.key
  scopes = each.value
}

resource "sendgrid_teammate" "users" {
  for_each = {
    "user1" = "user1@example.com"
    "user2" = "user2@example.com"
    # Limit to smaller batches
  }

  email    = each.value
  is_admin = false
  scopes   = ["mail.send"]
}
```

### 3. Monitor API Usage

Keep track of your SendGrid API usage in the SendGrid dashboard to understand your rate limits.

### 4. Environment Variables

Set environment variables for better control:

```bash
export SENDGRID_API_KEY="your-api-key"
export TF_LOG=INFO  # To see retry attempts in logs
```

## Troubleshooting

### Still Getting Rate Limit Errors?

1. **Increase timeouts**: Allow more time for retries
2. **Reduce parallelism**: Use `-parallelism=1` for sequential operations
3. **Check API key permissions**: Ensure your API key has proper scopes
4. **Review SendGrid plan limits**: Verify your plan's API rate limits
5. **Batch your operations**: Apply changes in smaller groups

### Monitoring Retries

Enable Terraform logging to see retry attempts:

```bash
export TF_LOG=INFO
terraform apply
```

Look for log entries containing "retry" or "rate limit" to monitor the retry behavior.

### Common Rate Limit Scenarios

#### Creating Multiple API Keys

When creating 6+ API keys, you may hit rate limits. SendGrid allows approximately one API key creation every 1-3 seconds.

**Solution:**

```bash
# Apply with sequential processing
terraform apply -parallelism=1

# Or apply in batches
terraform apply -target=sendgrid_api_key.key1 -target=sendgrid_api_key.key2
terraform apply -target=sendgrid_api_key.key3 -target=sendgrid_api_key.key4
# etc.
```

#### Managing Many Teammates

Similar restrictions apply when creating multiple teammates.

**Solution:**

```bash
# Use reduced parallelism
terraform apply -parallelism=2
```

## API Rate Limits

SendGrid API rate limits vary by endpoint and plan. Common limits include:

- **Web API v3**: 10,000 requests per hour (Pro plan and above)
- **API Key creation**: ~1 request per 1-3 seconds
- **Teammate management**: Lower limits may apply
- **Template operations**: Moderate limits
- **Burst limits**: Short-term burst allowances

Check the [SendGrid API documentation](https://docs.sendgrid.com/api-reference/how-to-use-the-sendgrid-v3-api/rate-limits) for current rate limit information.

## Error Messages

When rate limits are exceeded, you may see errors like:

```
Error: api response: HTTP 429: {"errors":[{"field":null,"message":"too many requests"}]}
```

These errors will now be automatically retried with exponential backoff until the configured timeout is reached.
