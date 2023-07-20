# sendgrid_link_branding_validation

Provide a resource to manage a link branding validation.

## Example Usage

```hcl

	resource "sendgrid_link_branding_validation" "foo" {
		link_branding_id = sendgrid_link_branding.foo.id
	}

```

## Argument Reference

The following arguments are supported:

* `link_branding_id` - (Required) Id of the link branding to validate.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `valid` - Indicates if this is a valid link branding or not.

