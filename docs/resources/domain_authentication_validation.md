# sendgrid_domain_authentication_validation

Provide a resource to manage a domain authentication validation.

## Example Usage

```hcl

	resource "sendgrid_domain_authentication_validation" "foo" {
		domain_authentication_id = sendgrid_domain_authentication.foo.id
	}

```

## Argument Reference

The following arguments are supported:

* `domain_authentication_id` - (Required) Id of the domain authentication to validate.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `valid` - Indicates if this is a valid authenticated domain or not.

