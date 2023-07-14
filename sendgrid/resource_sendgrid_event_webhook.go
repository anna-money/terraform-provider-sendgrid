/*
Provide a resource to manage an event webhook settings.
Example Usage
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
*/
package sendgrid

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	sendgrid "github.com/octoenergy/terraform-provider-sendgrid/sdk"
)

func resourceSendgridEventWebhook() *schema.Resource { //nolint:funlen
	return &schema.Resource{
		CreateContext: resourceSendgridEventWebhookPatch,
		ReadContext:   resourceSendgridEventWebhookRead,
		UpdateContext: resourceSendgridEventWebhookPatch,
		DeleteContext: resourceSendgridEventWebhookDelete,

		Schema: map[string]*schema.Schema{
			"enabled": {
				Type:        schema.TypeBool,
				Description: "Indicates if the event webhook is enabled.",
				Required:    true,
			},
			"url": {
				Type: schema.TypeString,
				Description: "The public URL where you would like SendGrid to POST the data events from your email. " +
					"Any emails sent with the given hostname provided (whose MX records have been updated to point to SendGrid) " +
					"will be eventd and POSTed to this URL.",
				Required: true,
			},
			"group_resubscribe": {
				Type: schema.TypeBool,
				Description: "Recipient resubscribes to specific group by updating preferences. " +
					"You need to enable Subscription Tracking for getting this type of event.",
				Optional: true,
				Default:  true,
			},
			"delivered": {
				Type:        schema.TypeBool,
				Description: "Message has been successfully delivered to the receiving server.",
				Optional:    true,
				Default:     true,
			},
			"group_unsubscribe": {
				Type: schema.TypeBool,
				Description: "Recipient unsubscribe from specific group, by either direct link or updating preferences. " +
					"You need to enable Subscription Tracking for getting this type of event.",
				Optional: true,
				Default:  true,
			},
			"spam_report": {
				Type:        schema.TypeBool,
				Description: "Recipient marked a message as spam.",
				Optional:    true,
				Default:     true,
			},
			"bounce": {
				Type:        schema.TypeBool,
				Description: "Receiving server could not or would not accept message.",
				Optional:    true,
				Default:     true,
			},
			"deferred": {
				Type:        schema.TypeBool,
				Description: "Recipient's email server temporarily rejected message.",
				Optional:    true,
				Default:     true,
			},
			"unsubscribe": {
				Type: schema.TypeBool,
				Description: "Recipient clicked on message's subscription management link. " +
					"You need to enable Subscription Tracking for getting this type of event.",
				Optional: true,
				Default:  true,
			},
			"processed": {
				Type:        schema.TypeBool,
				Description: "Message has been received and is ready to be delivered.",
				Optional:    true,
				Default:     true,
			},
			"open": {
				Type: schema.TypeBool,
				Description: "Recipient has opened the HTML message. " +
					"You need to enable Open Tracking for getting this type of event.",
				Optional: true,
				Default:  true,
			},
			"click": {
				Type: schema.TypeBool,
				Description: "Recipient clicked on a link within the message. " +
					"You need to enable Click Tracking for getting this type of event.",
				Optional: true,
				Default:  true,
			},
			"dropped": {
				Type: schema.TypeBool,
				Description: "You may see the following drop reasons: " +
					"Invalid SMTPAPI header, Spam Content (if spam checker app enabled), " +
					"Unsubscribed Address, Bounced Address, Spam Reporting Address, Invalid, Recipient List over Package Quota.",
				Optional: true,
				Default:  true,
			},
			"oauth_client_id": {
				Type: schema.TypeString,
				Description: "The client ID Twilio SendGrid sends to your OAuth server or " +
					"service provider to generate an OAuth access token.",
				Optional: true,
			},
			"oauth_client_secret": {
				Type: schema.TypeString,
				Description: "This secret is needed only once to create an access token. SendGrid will store this secret, " +
					"allowing you to update your Client ID and Token URL without passing the secret to SendGrid again. " +
					"When passing data in this field, you must also include the oauth_client_id and oauth_token_url fields.",
				Optional:  true,
				Sensitive: true,
			},
			"oauth_token_url": {
				Type: schema.TypeString,
				Description: "The URL where Twilio SendGrid sends the Client ID and Client Secret to generate an access token. " +
					"This should be your OAuth server or service provider. " +
					"When passing data in this field, you must also include the oauth_client_id field.",
				Optional: true,
			},
			"signed": {
				Type:        schema.TypeBool,
				Description: "Should the event webhook use signing?",
				Optional:    true,
			},
			"public_key": {
				Type:        schema.TypeString,
				Description: "The public key used to sign the event webhook. Only present if 'signed' is true",
				Computed:    true,
			},
		},
	}
}

func resourceSendgridEventWebhookDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// no op as there is no way to delete it
	return nil
}

func resourceSendgridEventWebhookPatch(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*sendgrid.Client)

	enabled := d.Get("enabled").(bool)
	url := d.Get("url").(string)
	groupResubscribe := d.Get("group_resubscribe").(bool)
	delivered := d.Get("delivered").(bool)
	groupUnsubscribe := d.Get("group_unsubscribe").(bool)
	spamReport := d.Get("spam_report").(bool)
	bounce := d.Get("bounce").(bool)
	deferred := d.Get("deferred").(bool)
	unsubscribe := d.Get("unsubscribe").(bool)
	processed := d.Get("processed").(bool)
	open := d.Get("open").(bool)
	click := d.Get("click").(bool)
	dropped := d.Get("dropped").(bool)
	oauthClientID := d.Get("oauth_client_id").(string)
	oauthClientSecret := d.Get("oauth_client_secret").(string)
	oauthTokenURL := d.Get("oauth_token_url").(string)

	_, err := sendgrid.RetryOnRateLimit(ctx, d, func() (interface{}, sendgrid.RequestError) {
		return c.PatchEventWebhook(
			ctx,
			enabled,
			url,
			groupResubscribe,
			delivered,
			groupUnsubscribe,
			spamReport,
			bounce,
			deferred,
			unsubscribe,
			processed,
			open,
			click,
			dropped,
			oauthClientID,
			oauthClientSecret,
			oauthTokenURL,
		)
	})
	if err != nil {
		return diag.FromErr(err)
	}

	if d.HasChange("signed") {
		if _, err := c.ConfigureEventWebhookSigning(ctx, d.Get("signed").(bool)); err.Err != nil {
			return diag.FromErr(err.Err)
		}
	}

	if c.OnBehalfOf != "" {
		d.SetId(c.OnBehalfOf) // since there is only a global event webhook per subuser
	} else {
		d.SetId("default") // or at the parent account level
	}

	return resourceSendgridEventWebhookRead(ctx, d, m)
}

func resourceSendgridEventWebhookRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*sendgrid.Client)

	webhook, err := c.ReadEventWebhook(ctx)
	if err.Err != nil {
		return diag.FromErr(err.Err)
	}

	//nolint:errcheck
	d.Set("enabled", webhook.Enabled)
	//nolint:errcheck
	d.Set("url", webhook.URL)
	//nolint:errcheck
	d.Set("group_resubscribe", webhook.GroupResubscribe)
	//nolint:errcheck
	d.Set("delivered", webhook.Delivered)
	//nolint:errcheck
	d.Set("group_unsubscribe", webhook.GroupUnsubscribe)
	//nolint:errcheck
	d.Set("spam_report", webhook.SpamReport)
	//nolint:errcheck
	d.Set("bounce", webhook.Bounce)
	//nolint:errcheck
	d.Set("deferred", webhook.Delivered)
	//nolint:errcheck
	d.Set("unsubscribe", webhook.Unsubscribe)
	//nolint:errcheck
	d.Set("processed", webhook.Processed)
	//nolint:errcheck
	d.Set("open", webhook.Open)
	//nolint:errcheck
	d.Set("click", webhook.Click)
	//nolint:errcheck
	d.Set("dropped", webhook.Dropped)
	//nolint:errcheck
	d.Set("oauth_client_id", webhook.OAuthClientID)
	//nolint:errcheck
	d.Set("oauth_token_url", webhook.OAuthTokenURL)

	webhookSigning, err := c.ReadEventWebhookSigning(ctx)
	if err.Err != nil {
		return diag.FromErr(err.Err)
	}
	//nolint:errcheck
	d.Set("public_key", webhookSigning.PublicKey)
	//nolint:errcheck
	d.Set("signed", webhookSigning.PublicKey != "")

	return nil
}
