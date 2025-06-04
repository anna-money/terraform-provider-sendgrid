# Rate Limiting Handling

The SendGrid Terraform Provider includes built-in rate limiting handling to manage HTTP 429 "too many requests" errors from the SendGrid API.

## How it works

All resource operations (create, read, update, delete) automatically retry when they encounter HTTP 429 rate limit errors using exponential backoff. The provider will:

1. Detect HTTP 429 responses from the SendGrid API
2. Wait and retry the request with exponential backoff
3. Continue retrying until the timeout is reached
4. Use the timeout configured in your Terraform resource

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
```

## Best Practices

### 1. Use Parallelism Limits

Reduce the number of concurrent operations to avoid hitting rate limits:

```bash
terraform apply -parallelism=2
```

### 2. Batch Operations

When creating multiple teammates, consider spacing out the operations:

```hcl
# Instead of creating many teammates at once, consider using for_each
# with a smaller set or using terraform apply with lower parallelism

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

### Monitoring Retries

Enable Terraform logging to see retry attempts:

```bash
export TF_LOG=INFO
terraform apply
```

Look for log entries containing "retry" or "rate limit" to monitor the retry behavior.

## API Rate Limits

SendGrid API rate limits vary by endpoint and plan. Common limits include:

- **Web API v3**: 10,000 requests per hour (Pro plan and above)
- **Teammate management**: Lower limits may apply
- **Burst limits**: Short-term burst allowances

Check the [SendGrid API documentation](https://docs.sendgrid.com/api-reference/how-to-use-the-sendgrid-v3-api/rate-limits) for current rate limit information.
