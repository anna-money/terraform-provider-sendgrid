/*
Provide a resource to manage a sso user.
Example Usage
```hcl

	resource "sendgrid_teammate_sso" "user" {
		first_name   = "Denis"
		last_name    = "Arslanbekov"
		email        = arslanbekov@gmail.com
		is_read_only = falst
		is_admin     = false
		scopes       = [
			""
		]
	}

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

func resourceSendgridTeammateSSO() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceSendgridTeammateSSOCreate,
		ReadContext:   resourceSendgridTeammateSSORead,
		UpdateContext: resourceSendgridTeammateSSOUpdate,
		DeleteContext: resourceSendgridTeammateSSODelete,
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
				Type:        schema.TypeBool,
				Description: "Invited user should be admin?.",
				Required:    true,
				Elem:        &schema.Schema{Type: schema.TypeBool},
			},
			"scopes": {
				Type:        schema.TypeSet,
				Description: "Permission scopes, will ignored if parameter is_admin = true.",
				Required:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func resourceSendgridTeammateSSOCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*sendgrid.Client)

	first_name := d.Get("first_name").(string)
	last_name := d.Get("last_name").(string)
	email := d.Get("email").(string)
	is_read_only := d.Get("is_read_only").(bool)
	is_admin := d.Get("is_admin").(bool)
	scopesSet := d.Get("scopes").(*schema.Set).List()
	scopes := make([]string, 0)

	for _, scope := range scopesSet {
		scopes = append(scopes, scope.(string))
	}

	tflog.Debug(ctx, "Creating teammate", map[string]interface{}{"email": email, "is_admin": is_admin, "scopes": scopes})

	user, err := client.CreateUserSSO(ctx, first_name, last_name, email, scopes, is_read_only, is_admin)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(user.Email)

	return nil
}

func resourceSendgridTeammateSSORead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*sendgrid.Client)

	var diags diag.Diagnostics
	userEmail := d.Id()

	u, err := c.ReadUser(ctx, userEmail)
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}
	if u == nil {
		d.SetId("")
	} else {
		d.Set("email", u.Email)
	}
	return diags
}

func resourceSendgridTeammateSSOUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*sendgrid.Client)

	scopesSet := d.Get("scopes").(*schema.Set).List()
	scopes := make([]string, 0)

	for _, scope := range scopesSet {
		scopes = append(scopes, scope.(string))
	}

	_, err := c.UpdateUserSSO(ctx, d.Get("first_name").(string), d.Get("last_name").(string), d.Get("email").(string), scopes, d.Get("is_read_only").(bool), d.Get("is_admin").(bool))
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceSendgridTeammateRead(ctx, d, m)
}

func resourceSendgridTeammateSSODelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*sendgrid.Client)
	var diags diag.Diagnostics
	userEmail := d.Id()
	_, err := c.DeleteUser(ctx, userEmail)

	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}
	return diags
}
