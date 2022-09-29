/*
Provide a resource to manage a subuser.
Example Usage
```hcl

	resource "sendgrid_teammate" "user" {
		email    = arslanbekov@gmail.com
		is_admin = false
		scopes   = [
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
				Type:        schema.TypeBool,
				Description: "Invited user should be admin?.",
				Required:    true,
				Elem:        &schema.Schema{Type: schema.TypeBool},
			},
			"scopes": {
				Type:        schema.TypeSet,
				Description: "Permission scopes, will ignored if parameter is_admin = true.",
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
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

	tflog.Debug(ctx, "Creating teammate", map[string]interface{}{"email": email, "is_admin": is_admin, "scopes": scopes})

	user, err := client.CreateUser(ctx, email, scopes, is_admin)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(user.Email)
	d.Set("email", user.Email)

	return nil
}

func resourceSendgridTeammateRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
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

func resourceSendgridTeammateUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*sendgrid.Client)

	scopesSet := d.Get("scopes").(*schema.Set).List()
	scopes := make([]string, 0)

	for _, scope := range scopesSet {
		scopes = append(scopes, scope.(string))
	}

	_, err := c.UpdateUser(ctx, d.Get("email").(string), scopes, d.Get("is_admin").(bool))
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceSendgridTeammateRead(ctx, d, m)
}

func resourceSendgridTeammateDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*sendgrid.Client)
	var diags diag.Diagnostics
	userEmail := d.Id()
	_, err := c.DeleteUser(ctx, userEmail)

	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}
	return diags
}
