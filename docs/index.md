# Terraform Sendgrid provider

The Terraform Sendgrid provider is used to interact with many resources supported by Sendgrid.
The provider needs to be configured with the proper credentials before it can be used.

## Example Usage

```hcl
# Configure the provider
provider "sendgrid" {
    api_key = "SG.XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX"
}

# Create a template
resource "sendgrid_template" "template"{
    name       = "my-template"
    generation = "dynamic"
}

# Create a template version
resource "sendgrid_template_version" "template_version" {
    name                   = "my-template-version"
    template_id            = sendgrid_template.template.id
    active                 = 1
    html_content           = "<%body%>"
    generate_plain_content = true
    subject                = "subject"
}
```

## Authentication

The Sendgrid provider offers a flexible means of providing credentials for authentication.
The following methods are supported, and explained below in this order:

- Static credentials
- Environment variables

### Static credentials

Static credentials can be provided by adding `api_key` to the Sendgrid provider block, you can configure `host` and `subuser` too.

Usage:

```hcl
provider "sendgrid" {
    api_key = "SG.XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX"
    host    = "https://api.sendgrid.com/v3"
    subuser = "username"
}
```

### Environment variables

You can provide your credentials via `SENDGRID_API_KEY`. You also can set `SENDGRID_HOST` and `SENDGRID_SUBUSER` environment variables.

```hcl
provider "sendgrid" {}
```

Usage:

```shell
$ export SENDGRID_API_KEY="SG.XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX"
$ terraform plan
```

## Testing

Credentials must be provided via the `SENDGRID_API_KEY` environment variable in order to run acceptance tests.

## Datasources/Resources reference

### API key Resource

### sendgrid_api_key

### Domain authentication Resource

### sendgrid_domain_authentication

### Link branding Resource

### sendgrid_link_branding

### SSO Resources

### sendgrid_sso_certificate

### sendgrid sso_integration

### Subuser Resource

### sendgrid_subuser

### Template Resources

### sendgrid_template

### sendgrid_template_version

### Unsubscribe Group Resource

### sendgrid_unsubscribe_group

### WebHook Resources

### sendgrid_parse_webhook

### sendgrid_event_webhook
