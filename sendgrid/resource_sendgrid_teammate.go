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

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	sendgrid "github.com/octoenergy/terraform-provider-sendgrid/sdk"
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
			"is_admin": {
				Type:        schema.TypeBool,
				Description: "Invited user should be admin?",
				Required:    true,
				Elem: &schema.Schema{
					Type: schema.TypeBool,
				},
			},
			"is_sso": {
				Type:        schema.TypeBool,
				Description: "Single Sign-On user?",
				Required:    true,
				Elem: &schema.Schema{
					Type: schema.TypeBool,
				},
			},
			"scopes": {
				Type:        schema.TypeSet,
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
	isAdmin := d.Get("is_admin").(bool)
	isSSO := d.Get("is_sso").(bool)
	firstName := d.Get("first_name").(string)
	lastName := d.Get("last_name").(string)

	scopesSet := d.Get("scopes").(*schema.Set).List()
	var scopes []string
	if !isAdmin {
		for _, scope := range scopesSet {
			scopes = append(scopes, scope.(string))
		}
	}
	tflog.Debug(ctx, "Creating teammate", map[string]interface{}{
		"first_name": firstName, "last_name": lastName,
		"email": email, "is_admin": isAdmin, "scopes": scopes,
	})

	var user *sendgrid.User
	var err error
	if isSSO {
		user, err = client.CreateSSOUser(ctx, firstName, lastName, email, scopes, isAdmin)
	} else {
		user, err = client.CreateUser(ctx, email, scopes, isAdmin)
	}
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(user.Email)
	if err := d.Set("email", user.Email); err != nil {
		return diag.FromErr(err)
	}

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

	// There is no need to track admin scopes since they have full access.
	if teammate.IsAdmin {
		teammate.Scopes = nil
	}

	var filteredScopes []string
	for _, s := range teammate.Scopes {
		// Sendgrid sets these scopes automatically. If you try to set them, you will get a 400 error.
		if s != "2fa_exempt" && s != "2fa_required" {
			filteredScopes = append(filteredScopes, s)
		}
	}

	d.SetId(teammate.Email)
	retErr := multierror.Append(
		d.Set("email", teammate.Email),
		d.Set("username", teammate.Username),
		d.Set("first_name", teammate.FirstName),
		d.Set("last_name", teammate.LastName),
		d.Set("scopes", filteredScopes),
		d.Set("is_admin", teammate.IsAdmin),
	)
	return diag.FromErr(retErr.ErrorOrNil())
}

func resourceSendgridTeammateUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*sendgrid.Client)
	email := d.Get("email").(string)
	isAdmin := d.Get("is_admin").(bool)
	isSSO := d.Get("is_sso").(bool)
	firstName := d.Get("first_name").(string)
	lastName := d.Get("last_name").(string)

	scopesSet := d.Get("scopes").(*schema.Set).List()
	var scopes []string
	if !isAdmin {
		for _, scope := range scopesSet {
			scopes = append(scopes, scope.(string))
		}
	}

	var err error
	if isSSO {
		_, err = client.UpdateSSOUser(ctx, firstName, lastName, email, scopes, isAdmin)
	} else {
		_, err = client.UpdateUser(ctx, email, scopes, isAdmin)
	}

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
