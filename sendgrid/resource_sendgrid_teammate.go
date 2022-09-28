/*
Provide a resource to manage a subuser.
Example Usage
```hcl

	resource "sendgrid_teammate" "user" {
		email    = arslanbekov@gmail.com
		is_admin = false
		scopes   = ""
	}

```
Import
A teammate can be imported, e.g.
```hcl
$ terraform import sendgrid_teammate.user email
```
*/
package sendgrid

import (
	"context"
	sendgrid "github.com/anna-money/terraform-provider-sendgrid/sdk"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceSendgridTeammate() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceSendgridTeammateCreate,
		ReadContext:   resourceSendgridTeammateRead,
		UpdateContext: resourceSendgridTeammateUpdate,
		DeleteContext: resourceSendgridTeammateDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"email": {
				Type:        schema.TypeString,
				Description: "The email of the user.",
				Required:    true,
			},
			"is_admin": {
				Type:        schema.TypeSet,
				Description: "Invited user should be admin?.",
				Required:    true,
				Elem:        &schema.Schema{Type: schema.TypeBool},
			},
			"scopes": {
				Type:        schema.TypeSet,
				Description: "Permission scopes, will ignored if parameter is_admin = true.",
				Required:    true,
				Elem:        &schema.Schema{Type: schema.TypeMap},
			},
		},
	}
}

func resourceSendgridTeammateCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*sendgrid.Client)

	email := d.Get("email").(string)
	is_admin := d.Get("is_admin").(bool)
	scopesSet := d.Get("scopes").(*schema.Set).List()
	scopes := make([]string, 0)

	for _, scope := range scopesSet {
		scopes = append(scopes, scope.(string))
	}

	tflog.Debug(ctx, "Creating teammate", map[string]interface{}{"email": email, "is_admin": is_admin})

	teammate, err := client.InviteTeammate(ctx, email, scopes, is_admin)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(teammate.Email)

	if d.Get("disabled").(bool) {
		return resourceSendgridTeammateUpdate(ctx, d, meta)
	}

	return resourceSendgridTeammateRead(ctx, d, meta)
}
