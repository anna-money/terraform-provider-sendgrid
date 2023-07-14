# sendgrid_parse_webhook

Provide a resource to manage parse webhook settings.

## Example Usage

```hcl

	resource "sendgrid_parse_webhook" "default" {
		hostname = "parse.foo.bar"
	    url = "https://foo.bar/sendgrid/inbound"
	    spam_check = false
	    send_raw = false
	}

```

## Argument Reference

The following arguments are supported:

* `hostname` - (Required, ForceNew) A specific and unique domain or subdomain that you have created to use exclusively to parse your incoming email. For example, parse.yourdomain.com.
* `url` - (Required, ForceNew) The public URL where you would like SendGrid to POST the data parsed from your email. Any emails sent with the given hostname provided (whose MX records have been updated to point to SendGrid) will be parsed and POSTed to this URL.
* `send_raw` - (Optional) Indicates if you would like SendGrid to post the original MIME-type content of your parsed email. When this parameter is set to "true", SendGrid will send a JSON payload of the content of your email.
* `spam_check` - (Optional) Indicates if you would like SendGrid to check the content parsed from your emails for spam before POSTing them to your domain.


## Import

An unsubscribe webhook can be imported, e.g.
```hcl
$ terraform import sendgrid_parse_webhook.default hostname
```
