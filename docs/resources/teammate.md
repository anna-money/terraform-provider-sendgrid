# sendgrid_teammate

Provide a resource to manage a subuser.

## Example Usage

```hcl

	resource "sendgrid_teammate" "user" {
		email    = arslanbekov@gmail.com
		is_admin = false
		scopes   = [
			""
		]
	}

```

## Argument Reference

The following arguments are supported:

* `email` - (Required) The email of the user.
* `is_admin` - (Required) Invited user should be admin?
* `is_sso` - (Required) Single Sign-On user?
* `first_name` - (Optional) The first nameof the user.
* `last_name` - (Optional) The last name of the user.
* `scopes` - (Optional) Permission scopes, will ignored if parameter is_admin = true.
* `username` - (Optional) Username.

