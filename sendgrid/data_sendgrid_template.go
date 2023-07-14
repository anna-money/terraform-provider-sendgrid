package sendgrid

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	sendgrid "github.com/octoenergy/terraform-provider-sendgrid/sdk"
)

func dataSendgridTemplate() *schema.Resource {
	s := resourceSendgridTemplate().Schema

	for _, val := range s {
		val.Computed = true
		val.Optional = false
		val.Required = false
		val.Default = nil
		val.ValidateFunc = nil
		val.ForceNew = true
	}

	s["template_id"] = &schema.Schema{
		Type:        schema.TypeString,
		Optional:    true,
		Description: "The ID of the template to retrieve",
	}
	s["name"] = &schema.Schema{
		Type:        schema.TypeString,
		Optional:    true,
		Computed:    true,
		Description: "The name of the template to retrieve",
	}
	s["generation"] = &schema.Schema{
		Type:     schema.TypeString,
		Optional: true,
		Computed: true,
	}

	return &schema.Resource{
		ReadContext: dataSendgridTemplateRead,
		Schema:      s,
	}
}

func dataSendgridTemplateRead(context context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	templateID := d.Get("template_id").(string)
	name := d.Get("name").(string)
	c := m.(*sendgrid.Client)

	switch {
	case templateID != "":
		d.SetId(templateID)

		return resourceSendgridTemplateRead(context, d, m)
	case name != "":
		generation := d.Get("generation").(string)
		if generation == "" {
			generation = "dynamic"
		}

		templates, err := c.ReadTemplates(context, generation)
		if err != nil {
			return diag.FromErr(err)
		}

		names := make([]string, 0)

		for i := range templates {
			template := templates[i]

			if template.Name == name {
				d.SetId(template.ID)

				if err = sendgridTemplateParse(&template, d); err != nil {
					return diag.FromErr(err)
				}

				return nil
			}

			names = append(names, template.Name)
		}

		return diag.Errorf("unable to find a template with name '%s', valid names are %v", name, names)
	default:
		return diag.Errorf("either 'template_id' or 'name' must be specified for data.sendgrid_template")
	}
}
