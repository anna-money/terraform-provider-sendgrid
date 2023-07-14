# sendgrid_unsubscribe_group

Provide a resource to manage an unsubscribe group.

## Example Usage

```hcl

	resource "sendgrid_unsubscribe_group" "default" {
		name   = "default-unsubscribe-group"
		description = "The default unsubscribe group"
	    is_default = true
	}

```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name you will use to describe this unsubscribe group.
* `description` - (Optional) The description of the unsubscribe group
* `is_default` - (Optional) Should this unsubscribe group be used as the default group?

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `unsubscribes` - The number of unsubscribes that belong to the group.


## Import

An unsubscribe group can be imported, e.g.
```hcl
$ terraform import sendgrid_unsubscribe_group.default unsubscribeGroupID
```
