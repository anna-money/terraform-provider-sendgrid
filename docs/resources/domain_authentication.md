# sendgrid_domain_authentication

Provide a resource to manage a domain authentication.

## Example Usage

```hcl

	resource "sendgrid_domain_authentication" "default" {
		domain = "example.com"
	    ips = [ "10.10.10.10" ]
	    is_default = true
	    automatic_security = false
	}

```

## Argument Reference

The following arguments are supported:

* `domain` - (Required, ForceNew) Domain being authenticated.
* `automatic_security` - (Optional, ForceNew) Whether to allow SendGrid to manage your SPF records, DKIM keys, and DKIM key rotation.
* `custom_dkim_selector` - (Optional, ForceNew) Add a custom DKIM selector. Accepts three letters or numbers.
* `custom_spf` - (Optional) Specify whether to use a custom SPF or allow SendGrid to manage your SPF. This option is only available to authenticated domains set up for manual security.
* `ips` - (Optional, ForceNew) The IP addresses that will be included in the custom SPF record for this.
* `is_default` - (Optional) Whether to use this authenticated domain as the fallback if no authenticated domains match the sender's domain.
* `subdomain` - (Optional, ForceNew) The subdomain to use for this authenticated domain.
* `valid` - (Optional) Indicates if this is a valid authenticated domain or not.

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
$ terraform import sendgrid_domain_authentication.default domainId
```
