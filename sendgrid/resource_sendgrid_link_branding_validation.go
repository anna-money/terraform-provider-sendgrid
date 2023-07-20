/*
Provide a resource to manage a link branding validation.
Example Usage
```hcl

	resource "sendgrid_link_branding_validation" "foo" {
		link_branding_id = sendgrid_link_branding.foo.id
	}

```
*/
package sendgrid

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	sendgrid "github.com/octoenergy/terraform-provider-sendgrid/sdk"
)

// https://docs.sendgrid.com/api-reference/link-branding/validate-a-branded-link
func resourceSendgridLinkBrandingValidation() *schema.Resource { //nolint:funlen
	return &schema.Resource{
		CreateContext: resourceSendgridLinkBrandingValidationCreate,
		ReadContext:   resourceSendgridLinkBrandingValidationRead,
		UpdateContext: resourceSendgridLinkBrandingValidationUpdate,
		DeleteContext: resourceSendgridLinkBrandingValidationDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"link_branding_id": {
				Type:        schema.TypeString,
				Description: "Id of the link branding to validate.",
				Required:    true,
			},

			"valid": {
				Type:        schema.TypeBool,
				Description: "Indicates if this is a valid link branding or not.",
				Computed:    true,
			},
		},
	}
}

func resourceSendgridLinkBrandingValidationCreate(
	ctx context.Context,
	d *schema.ResourceData,
	m interface{},
) diag.Diagnostics {
	return validateLinkBranding(ctx, d, m)
}

func validateLinkBranding(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*sendgrid.Client)

	if err := c.ValidateLinkBranding(ctx, d.Get("link_branding_id").(string)); err.Err != nil || err.StatusCode != 200 {
		if err.Err != nil {
			return diag.FromErr(err.Err)
		}
		return diag.Errorf("unable to validate domain DNS configuration")
	}

	return resourceSendgridLinkBrandingValidationRead(ctx, d, m)
}

func resourceSendgridLinkBrandingValidationRead( //nolint:funlen,cyclop
	ctx context.Context,
	d *schema.ResourceData,
	m interface{},
) diag.Diagnostics {
	c := m.(*sendgrid.Client)

	link, err := c.ReadLinkBranding(ctx, d.Get("link_branding_id").(string))
	if err.Err != nil {
		return diag.FromErr(err.Err)
	}

	//nolint:errcheck
	d.Set("valid", link.Valid)
	d.SetId(fmt.Sprint(link.ID))
	return nil
}

func resourceSendgridLinkBrandingValidationUpdate(
	ctx context.Context,
	d *schema.ResourceData,
	m interface{},
) diag.Diagnostics {
	return validateLinkBranding(ctx, d, m)
}

func resourceSendgridLinkBrandingValidationDelete(context.Context, *schema.ResourceData, interface{}) diag.Diagnostics {
	return nil
}
