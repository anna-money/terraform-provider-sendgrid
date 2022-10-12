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
	"github.com/hashicorp/go-multierror"
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
			"first_name": {
				Type:        schema.TypeString,
				Description: "The first nameof the user.",
				Optional:    true,
			},
			"last_name": {
				Type:        schema.TypeString,
				Description: "The last name of the user.",
				Optional:    true,
			},
			"user_type": {
				Type:        schema.TypeString,
				Description: "Type of the user.",
				Optional:    true,
			},
			"is_admin": {
				Type:        schema.TypeBool,
				Description: "Invited user should be admin?.",
				Required:    true,
				Elem: &schema.Schema{
					Type: schema.TypeBool,
				},
			},
			"scopes": {
				Type:        schema.TypeList,
				Description: "Permission scopes, will ignored if parameter is_admin = true.",
				Optional:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"username": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Username.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
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

func resourceSendgridTeammateRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*sendgrid.Client)

	var diags diag.Diagnostics
	email := d.Id()

	teammate, err := client.ReadUser(ctx, email)
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}

	d.SetId(teammate.Email)
	retErr := multierror.Append(
		d.Set("email", teammate.Email),
		d.Set("username", teammate.Username),
		d.Set("user_type", teammate.UserType),
		d.Set("first_name", teammate.FirstName),
		d.Set("last_name", teammate.LastName),
		d.Set("scopes", teammate.Scopes),
		d.Set("is_admin", teammate.IsAdmin),
	)
	return diag.FromErr(retErr.ErrorOrNil())
}

func resourceSendgridTeammateUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*sendgrid.Client)

	scopesSet := d.Get("scopes").(*schema.Set).List()
	scopes := make([]string, 0)

	for _, scope := range scopesSet {
		scopes = append(scopes, scope.(string))
	}

	_, err := client.UpdateUser(ctx, d.Get("email").(string), scopes, d.Get("is_admin").(bool))
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceSendgridTeammateRead(ctx, d, meta)
}

func resourceSendgridTeammateDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*sendgrid.Client)

	var diags diag.Diagnostics
	userEmail := d.Id()
	_, err := client.DeleteUser(ctx, userEmail)

	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}
	return diags
}
