/*
Provide a resource to manage an API key.
Example Usage
```hcl

	resource "sendgrid_api_key" "api_key" {
		name   = "my-api-key"
		scopes = [
			"mail.send",
			"sender_verification_eligible",
		]
	}

```
Import
An API key can be imported, e.g.
```hcl
$ terraform import sendgrid_api_key.api_key apiKeyID
```
*/
package sendgrid

import (
	"context"
	"reflect"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	sendgrid "github.com/octoenergy/terraform-provider-sendgrid/sdk"
)

func resourceSendgridAPIKey() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceSendgridAPIKeyCreate,
		ReadContext:   resourceSendgridAPIKeyRead,
		UpdateContext: resourceSendgridAPIKeyUpdate,
		DeleteContext: resourceSendgridAPIKeyDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Description:  "The name you will use to describe this API Key.",
				Required:     true,
				ValidateFunc: validation.StringLenBetween(1, maxStringLength),
			},
			"scopes": {
				Type:        schema.TypeSet,
				Description: "The individual permissions that you are giving to this API Key.",
				Optional:    true,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"api_key": {
				Type:        schema.TypeString,
				Description: "The API key created by the API.",
				Computed:    true,
				Sensitive:   true,
			},
		},
	}
}

func scopeInScopes(scopes []string, scope string) bool {
	for _, v := range scopes {
		if v == scope {
			return true
		}
	}

	return false
}

func resourceSendgridAPIKeyCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var scopes []string

	c := m.(*sendgrid.Client)
	name := d.Get("name").(string)

	for _, scope := range d.Get("scopes").(*schema.Set).List() {
		scopes = append(scopes, scope.(string))
	}

	if ok := scopeInScopes(scopes, "sender_verification_eligible"); !ok {
		scopes = append(scopes, "sender_verification_eligible")
	}

	apiKeyStruct, err := sendgrid.RetryOnRateLimit(ctx, d, func() (interface{}, sendgrid.RequestError) {
		return c.CreateAPIKey(ctx, name, scopes)
	})

	apiKey := apiKeyStruct.(*sendgrid.APIKey)

	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(apiKey.ID)
	//nolint:errcheck
	d.Set("api_key", apiKey.APIKey)

	return resourceSendgridAPIKeyRead(ctx, d, m)
}

func resourceSendgridAPIKeyRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*sendgrid.Client)

	apiKey, err := c.ReadAPIKey(ctx, d.Id())
	if err.Err != nil {
		return diag.FromErr(err.Err)
	}

	//nolint:errcheck
	d.Set("name", apiKey.Name)
	//nolint:errcheck
	d.Set("scopes", apiKey.Scopes)

	return nil
}

func hasDiff(o, n interface{}) bool {
	if eq, ok := o.(schema.Equal); ok {
		return !eq.Equal(n)
	}

	return !reflect.DeepEqual(o, n)
}

func resourceSendgridAPIKeyUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*sendgrid.Client)

	a := sendgrid.APIKey{
		ID:   d.Id(),
		Name: d.Get("name").(string),
	}

	o, n := d.GetChange("scopes")
	n.(*schema.Set).Add("sender_verification_eligible")
	n.(*schema.Set).Add("2fa_required")

	if ok := hasDiff(o, n); ok {
		var scopes []string
		for _, scope := range d.Get("scopes").(*schema.Set).List() {
			scopes = append(scopes, scope.(string))
		}

		a.Scopes = scopes
	}

	_, err := sendgrid.RetryOnRateLimit(ctx, d, func() (interface{}, sendgrid.RequestError) {
		return c.UpdateAPIKey(ctx, d.Id(), a.Name, a.Scopes)
	})
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceSendgridAPIKeyRead(ctx, d, m)
}

func resourceSendgridAPIKeyDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*sendgrid.Client)

	_, err := sendgrid.RetryOnRateLimit(ctx, d, func() (interface{}, sendgrid.RequestError) {
		return c.DeleteAPIKey(ctx, d.Id())
	})
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}
