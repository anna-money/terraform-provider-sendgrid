# sendgrid_event_webhook

Provide a resource to manage an event webhook settings.

## Example Usage

```hcl

	resource "sendgrid_event_webhook" "default" {
		enabled = true
	    url = "https://foo.bar/sendgrid/inbound"
	    group_resubscribe = true
	    delivered = true
	    group_unsubscribe = true
	    spam_report = true
	    bounce = true
	    deferred = true
	    unsubscribe = true
	    processed = true
	    open = true
	    click = true
	    dropped = true
	    oauth_client_id = "a-client-id"
	    oauth_client_secret = "a-client-secret"
	    oauth_token_url = "https://oauth.example.com/token"
	}

```

## Argument Reference

The following arguments are supported:

* `enabled` - (Required) Indicates if the event webhook is enabled.
* `url` - (Required) The public URL where you would like SendGrid to POST the data events from your email. Any emails sent with the given hostname provided (whose MX records have been updated to point to SendGrid) will be eventd and POSTed to this URL.
* `bounce` - (Optional) Receiving server could not or would not accept message.
* `click` - (Optional) Recipient clicked on a link within the message. You need to enable Click Tracking for getting this type of event.
* `deferred` - (Optional) Recipient's email server temporarily rejected message.
* `delivered` - (Optional) Message has been successfully delivered to the receiving server.
* `dropped` - (Optional) You may see the following drop reasons: Invalid SMTPAPI header, Spam Content (if spam checker app enabled), Unsubscribed Address, Bounced Address, Spam Reporting Address, Invalid, Recipient List over Package Quota.
* `group_resubscribe` - (Optional) Recipient resubscribes to specific group by updating preferences. You need to enable Subscription Tracking for getting this type of event.
* `group_unsubscribe` - (Optional) Recipient unsubscribe from specific group, by either direct link or updating preferences. You need to enable Subscription Tracking for getting this type of event.
* `oauth_client_id` - (Optional) The client ID Twilio SendGrid sends to your OAuth server or service provider to generate an OAuth access token.
* `oauth_client_secret` - (Optional) This secret is needed only once to create an access token. SendGrid will store this secret, allowing you to update your Client ID and Token URL without passing the secret to SendGrid again. When passing data in this field, you must also include the oauth_client_id and oauth_token_url fields.
* `oauth_token_url` - (Optional) The URL where Twilio SendGrid sends the Client ID and Client Secret to generate an access token. This should be your OAuth server or service provider. When passing data in this field, you must also include the oauth_client_id field.
* `open` - (Optional) Recipient has opened the HTML message. You need to enable Open Tracking for getting this type of event.
* `processed` - (Optional) Message has been received and is ready to be delivered.
* `signed` - (Optional) Should the event webhook use signing?
* `spam_report` - (Optional) Recipient marked a message as spam.
* `unsubscribe` - (Optional) Recipient clicked on message's subscription management link. You need to enable Subscription Tracking for getting this type of event.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `public_key` - The public key used to sign the event webhook. Only present if 'signed' is true

