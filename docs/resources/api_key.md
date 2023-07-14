# sendgrid_api_key

Provide a resource to manage an API key.

## Example Usage

```hcl

	resource "sendgrid_api_key" "api_key" {
		name   = "my-api-key"
		scopes = [
			"mail.send",
			"sender_verification_eligible",
		]
	}

```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name you will use to describe this API Key.
* `scopes` - (Optional) The individual permissions that you are giving to this API Key.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `api_key` - The API key created by the API.


## Import

An API key can be imported, e.g.
```hcl
$ terraform import sendgrid_api_key.api_key apiKeyID
```
