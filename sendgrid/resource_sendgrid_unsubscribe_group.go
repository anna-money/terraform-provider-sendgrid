/*
Provide a resource to manage an unsubscribe group.
Example Usage
```hcl
resource "sendgrid_unsubscribe_group" "default" {
	name   = "default-unsubscribe-group"
	description = "The default unsubscribe group"
    is_default = true
}
```
Import
An unsubscribe group can be imported, e.g.
```hcl
$ terraform import sendgrid_unsubscribe_group.default unsubscribeGroupID
```
*/
package sendgrid

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	sendgrid "github.com/trois-six/terraform-provider-sendgrid/sdk"
)

func resourceSendgridUnsubscribeGroup() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceSendgridUnsubscribeGroupCreate,
		ReadContext:   resourceSendgridUnsubscribeGroupRead,
		UpdateContext: resourceSendgridUnsubscribeGroupUpdate,
		DeleteContext: resourceSendgridUnsubscribeGroupDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Description:  "The name you will use to describe this unsubscribe group.",
				Required:     true,
				ValidateFunc: validation.StringLenBetween(1, unsubscribeGroupLength),
			},
			"description": {
				Type:         schema.TypeString,
				Description:  "The description of the unsubscribe group",
				Optional:     true,
				ValidateFunc: validation.StringLenBetween(0, maxStringLength),
			},
			"is_default": {
				Type:        schema.TypeBool,
				Description: "Should this unsubscribe group be used as the default group?",
				Optional:    true,
			},
			"unsubscribes": {
				Type:        schema.TypeInt,
				Description: "The number of unsubscribes that belong to the group.",
				Computed:    true,
			},
		},
	}
}

func resourceSendgridUnsubscribeGroupCreate(
	ctx context.Context,
	d *schema.ResourceData,
	m interface{}) diag.Diagnostics {
	c := m.(*sendgrid.Client)

	name := d.Get("name").(string)
	description := d.Get("description").(string)
	isDefault := d.Get("is_default").(bool)

	apiKeyStruct, err := sendgrid.RetryOnRateLimit(ctx, d, func() (interface{}, sendgrid.RequestError) {
		return c.CreateUnsubscribeGroup(name, description, isDefault)
	})

	group := apiKeyStruct.(*sendgrid.UnsubscribeGroup)

	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(fmt.Sprint(group.ID))

	return resourceSendgridUnsubscribeGroupRead(ctx, d, m)
}

func resourceSendgridUnsubscribeGroupRead(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*sendgrid.Client)

	group, err := c.ReadUnsubscribeGroup(d.Id())
	if err.Err != nil {
		return diag.FromErr(err.Err)
	}

	if err := sendgridUnsubscribeGroupParse(group, d); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func sendgridUnsubscribeGroupParse(group *sendgrid.UnsubscribeGroup, d *schema.ResourceData) error {
	if err := d.Set("name", group.Name); err != nil {
		return ErrSetUnsubscribeGroupName
	}

	if err := d.Set("description", group.Description); err != nil {
		return ErrSetUnsubscribeGroupDesc
	}

	if err := d.Set("is_default", group.IsDefault); err != nil {
		return ErrSetUnsubscribeGroupIsDefault
	}

	if err := d.Set("unsubscribes", group.Unsubscribes); err != nil {
		return ErrSetUnsubscribeGroupUnsuscribes
	}

	return nil
}

func resourceSendgridUnsubscribeGroupUpdate(
	ctx context.Context,
	d *schema.ResourceData,
	m interface{}) diag.Diagnostics {
	c := m.(*sendgrid.Client)

	name := d.Get("name").(string)
	description := d.Get("description").(string)
	isDefault := d.Get("is_default").(bool)

	_, err := sendgrid.RetryOnRateLimit(ctx, d, func() (interface{}, sendgrid.RequestError) {
		return c.UpdateUnsubscribeGroup(d.Id(), name, description, isDefault)
	})
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceSendgridUnsubscribeGroupRead(ctx, d, m)
}

func resourceSendgridUnsubscribeGroupDelete(
	ctx context.Context,
	d *schema.ResourceData,
	m interface{}) diag.Diagnostics {
	c := m.(*sendgrid.Client)

	_, err := sendgrid.RetryOnRateLimit(ctx, d, func() (interface{}, sendgrid.RequestError) {
		return c.DeleteUnsubscribeGroup(d.Id())
	})
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}
