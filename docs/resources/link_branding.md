# sendgrid_link_branding

Provide a resource to manage link branding.

## Example Usage

```hcl

	resource "sendgrid_link_branding" "default" {
		domain = "example.com"
	    is_default = true
	}

```

## Argument Reference

The following arguments are supported:

* `domain` - (Required, ForceNew) Domain being authenticated.
* `is_default` - (Optional) Indicates if this is the default link branding.
* `subdomain` - (Optional, ForceNew) The subdomain to use for this link branding.
* `valid` - (Optional) Indicates if this is a valid link branding or not. Set to `true` to attempt validation on first update.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `dns` - The DNS records used to authenticate the sending domain.
  * `data` - The actual DNS record.
  * `host` - The domain that this CNAME is created for.
  * `type` - The type of DNS record.
  * `valid` - Indicates if this is a valid CNAME.
* `username` - The username associated with this domain.


## Import

An unsubscribe group can be imported, e.g.
```hcl
$ terraform import sendgrid_link_branding.default linkId
```
