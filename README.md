# Terraform Provider for Sendgrid

> **⚠️ Unofficial Provider Notice**
> This is an **unofficial** Terraform provider for SendGrid that extends the original functionality with additional features:
>
> - **Teammate Management**: Full support for `sendgrid_teammate` resource and data source
> - **Enhanced Rate Limiting**: Comprehensive rate limiting support for all resources with automatic retry on HTTP 429 errors
> - **Improved Reliability**: Built-in exponential backoff strategy for better API interaction
>
> This provider is maintained independently and is not affiliated with or officially supported by SendGrid/Twilio.

## Usage

Detailed documentation is available on the [Terraform provider registry](https://registry.terraform.io/providers/anna-money/sendgrid/latest).

## Features

### Rate Limiting Support

This provider includes comprehensive rate limiting support for all SendGrid resources. When the SendGrid API returns HTTP 429 "too many requests" errors, the provider automatically retries with exponential backoff.

**Key features:**

- Automatic retry on HTTP 429 errors
- Exponential backoff strategy
- Configurable timeouts per resource
- Support for all SendGrid resources

For detailed information, see [Rate Limiting Documentation](docs/rate_limiting.md).

**Quick tips:**

- Use `-parallelism=1` for API key creation to avoid rate limits
- Configure custom timeouts for operations that may need more retry time
- Monitor SendGrid API usage in your dashboard

### Teammate Management

Unlike the official provider, this version includes full support for managing SendGrid teammates:

```hcl
resource "sendgrid_teammate" "example" {
  email    = "teammate@example.com"
  is_admin = false
  scopes   = ["mail.send", "marketing.read"]
}

data "sendgrid_teammate" "existing" {
  email = "existing@example.com"
}
```

## Build

To compile the provider, run `make build`. This will build the provider and put the provider binary in the `$GOPATH/bin` directory.

```sh
make build
```

In order to test the provider, you can simply run `make test`.

```sh
make test
```

In order to run the full suite of Acceptance tests, run `make testacc`.

```sh
make testacc
```

## Known issues

The API KEY API is not completely documented: when you don't set scopes, you get all scopes. This is managed by the provider.

When you set one or multiple scopes, even if you don't set the scopes `sender_verification_eligible` and `2fa_required`, you will get them in the end. It's managed by the provider: if you don't add these scopes to the list of scopes, the provider does it for you.
