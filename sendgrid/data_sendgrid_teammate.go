package sendgrid

import (
	"context"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	sendgrid "github.com/octoenergy/terraform-provider-sendgrid/sdk"
)

func dataSendgridTeammate() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSendgridTeammateRead,

		Schema: map[string]*schema.Schema{
			"username": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Teammate's username",
			},
			"first_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Teammate's first name",
			},
			"last_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Teammate's last name",
			},
			"email": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Teammate's email",
			},
			"scopes": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Scopes associated to teammate",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"user_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Indicate the type of user: account owner, teammate admin user, or normal teammate",
			},
			"is_admin": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "True if teammate has admin privileges",
			},
		},
	}
}

func dataSendgridTeammateRead(context context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*sendgrid.Client)
	email := d.Get("email").(string)
	tflog.Debug(context, "Reading user", map[string]interface{}{"email": email})

	teammate, err := client.ReadUser(context, email)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(teammate.Email)
	retErr := multierror.Append(
		d.Set("email", teammate.Email),
		d.Set("username", teammate.Username),
		d.Set("user_type", teammate.UserType),
		d.Set("last_name", teammate.LastName),
		d.Set("first_name", teammate.FirstName),
		d.Set("scopes", teammate.Scopes),
		d.Set("is_admin", teammate.IsAdmin),
	)

	return diag.FromErr(retErr.ErrorOrNil())
}
