/*
Provide a resource to manage parse webhook settings.
Example Usage
```hcl

	resource "sendgrid_parse_webhook" "default" {
		hostname = "parse.foo.bar"
	    url = "https://foo.bar/sendgrid/inbound"
	    spam_check = false
	    send_raw = false
	}

```
Import
An unsubscribe webhook can be imported, e.g.
```hcl
$ terraform import sendgrid_parse_webhook.default hostname
```
*/
package sendgrid

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	sendgrid "github.com/octoenergy/terraform-provider-sendgrid/sdk"
)

func resourceSendgridParseWebhook() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceSendgridParseWebhookCreate,
		ReadContext:   resourceSendgridParseWebhookRead,
		UpdateContext: resourceSendgridParseWebhookUpdate,
		DeleteContext: resourceSendgridParseWebhookDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"hostname": {
				Type: schema.TypeString,
				Description: "A specific and unique domain or subdomain " +
					"that you have created to use exclusively to parse your incoming email. " +
					"For example, parse.yourdomain.com.",
				Required: true,
				ForceNew: true,
			},
			"url": {
				Type: schema.TypeString,
				Description: "The public URL where you would like SendGrid to POST the data parsed from your email. " +
					"Any emails sent with the given hostname provided (whose MX records have been updated to point to SendGrid) " +
					"will be parsed and POSTed to this URL.",
				Required: true,
				ForceNew: true,
			},
			"spam_check": {
				Type: schema.TypeBool,
				Description: "Indicates if you would like SendGrid to check the content parsed from your emails for spam " +
					"before POSTing them to your domain.",
				Optional: true,
			},
			"send_raw": {
				Type: schema.TypeBool,
				Description: "Indicates if you would like SendGrid to post the original MIME-type content of your parsed email. " +
					"When this parameter is set to \"true\", SendGrid will send a JSON payload of the content of your email.",
				Optional: true,
			},
		},
	}
}

func resourceSendgridParseWebhookCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*sendgrid.Client)

	hostname := d.Get("hostname").(string)
	url := d.Get("url").(string)
	spamCheck := d.Get("spam_check").(bool)
	sendRaw := d.Get("send_raw").(bool)

	parseWebhookStruct, err := sendgrid.RetryOnRateLimit(ctx, d, func() (interface{}, sendgrid.RequestError) {
		return c.CreateParseWebhook(ctx, hostname, url, spamCheck, sendRaw)
	})

	webhook := parseWebhookStruct.(*sendgrid.ParseWebhook)

	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(webhook.Hostname)

	return resourceSendgridParseWebhookRead(ctx, d, m)
}

func resourceSendgridParseWebhookRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*sendgrid.Client)

	webhook, err := c.ReadParseWebhook(ctx, d.Id())
	if err.Err != nil {
		return diag.FromErr(err.Err)
	}

	//nolint:errcheck
	d.Set("hostname", webhook.Hostname)
	//nolint:errcheck
	d.Set("url", webhook.URL)
	//nolint:errcheck
	d.Set("spam_check", webhook.SpamCheck)
	//nolint:errcheck
	d.Set("send_raw", webhook.SendRaw)

	return nil
}

func resourceSendgridParseWebhookUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*sendgrid.Client)

	spamCheck := d.Get("spam_check").(bool)
	sendRaw := d.Get("send_raw").(bool)

	_, err := sendgrid.RetryOnRateLimit(ctx, d, func() (interface{}, sendgrid.RequestError) {
		return nil, c.UpdateParseWebhook(ctx, d.Id(), spamCheck, sendRaw)
	})
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceSendgridParseWebhookRead(ctx, d, m)
}

func resourceSendgridParseWebhookDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*sendgrid.Client)

	_, err := sendgrid.RetryOnRateLimit(ctx, d, func() (interface{}, sendgrid.RequestError) {
		return c.DeleteParseWebhook(ctx, d.Id())
	})
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}
