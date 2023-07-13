# sendgrid_subuser

Provide a resource to manage a subuser.

## Example Usage

```hcl

	resource "sendgrid_subuser" "subuser" {
		username = "my-subuser"
		email    = "subuser@example.org"
		password = "Passw0rd!"
		ips      = [
			"127.0.0.1"
		]
	}

```

## Argument Reference

The following arguments are supported:

* `email` - (Required) The email of the subuser.
* `ips` - (Required) The IP addresses that should be assigned to this subuser.
* `password` - (Required) The password the subuser will use when logging into SendGrid.
* `username` - (Required) The name of the subuser.


## Import

A subuser can be imported, e.g.
```hcl
$ terraform import sendgrid_subuser.subuser userName
```
