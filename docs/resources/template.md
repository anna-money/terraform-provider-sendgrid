# sendgrid_template

Provide a resource to manage a template of email.

## Example Usage

```hcl

	resource "sendgrid_template" "template" {
		name       = "my-template"
		generation = "dynamic"
	}

```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the template, max length: 100.
* `generation` - (Optional, ForceNew) Defines the generation of the template, allowed values: legacy, dynamic (default).

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `updated_at` - The date and time of the last update of this template.


## Import

A template can be imported, e.g.
```hcl
$ terraform import sendgrid_template.template templateID
```
