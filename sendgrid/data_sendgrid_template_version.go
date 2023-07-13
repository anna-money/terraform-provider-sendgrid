package sendgrid

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	sendgrid "github.com/octoenergy/terraform-provider-sendgrid/sdk"
)

func dataSendgridTemplateVersion() *schema.Resource {
	s := resourceSendgridTemplateVersion().Schema

	for key, val := range s {
		if key != "template_id" {
			val.Computed = true
			val.Optional = false
			val.Required = false
			val.Default = nil
			val.ValidateFunc = nil
		}
	}

	return &schema.Resource{
		ReadContext: dataSendgridTemplateVersionRead,
		Schema:      s,
	}
}

func dataSendgridTemplateVersionRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	templateID := d.Get("template_id").(string)
	c := m.(*sendgrid.Client)

	template, err := c.ReadTemplate(ctx, templateID)
	if err != nil {
		return diag.FromErr(err)
	}

	var activeVersion *sendgrid.TemplateVersion

	for i := range template.Versions {
		if template.Versions[i].Active == 1 {
			activeVersion = &template.Versions[i]

			break
		}
	}

	if activeVersion == nil {
		return diag.FromErr(ErrNoNewVersionFoundForTemplate)
	}

	d.SetId(activeVersion.ID)

	if err := parseTemplateVersion(d, activeVersion); err != nil {
		return diag.FromErr(err)
	}

	return nil
}
