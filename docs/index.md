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
* [resource sendgrid_api_key](resources/api_key.md)

### Domain authentication Resource
* [resource sendgrid_domain_authentication](resources/domain_authentication.md)

### Link branding Resource
* [resource sendgrid_link_branding](resources/link_branding.md)

### SSO Resources
* [resource sendgrid_sso_certificate](resources/sso_certificate.md)
* [resource sendgrid sso_integration](resources/sso_integration.md)

### Subuser Resource
* [resource sendgrid_subuser](resources/subuser.md)

### Template Resources
* [resource sendgrid_template](resources/template.md)
* [resource sendgrid_template_version](resources/template_version.md)

### Unsubscribe Group Resource
* [resource sendgrid_unsubscribe_group](resources/unsubscribe_group.md)

### WebHook Resources
* [resource sendgrid_parse_webhook](resources/parse_webhook.md)
* [resource sendgrid_event_webhook](resources/event_webhook.md)
