package sendgrid

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	sendgrid "github.com/octoenergy/terraform-provider-sendgrid/sdk"
)

func dataSendgridUnsubscribeGroup() *schema.Resource {
	s := resourceSendgridUnsubscribeGroup().Schema

	for _, val := range s {
		val.Computed = true
		val.Optional = false
		val.Required = false
		val.Default = nil
		val.ValidateFunc = nil
		val.ForceNew = true
	}

	s["group_id"] = &schema.Schema{
		Type:        schema.TypeString,
		Optional:    true,
		Description: "The id of the unsubscribe group to retrieve",
	}
	s["name"] = &schema.Schema{
		Type:        schema.TypeString,
		Optional:    true,
		Description: "The name of the unsubscribe group to retrieve",
	}

	return &schema.Resource{
		ReadContext: dataSendgridUnsubscribeGroupRead,
		Schema:      s,
	}
}

func dataSendgridUnsubscribeGroupRead(context context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	groupID := d.Get("group_id").(string)
	name := d.Get("name").(string)
	c := m.(*sendgrid.Client)

	switch {
	case groupID != "":
		d.SetId(groupID)

		return resourceSendgridUnsubscribeGroupRead(context, d, m)
	case name != "":
		groups, err := c.ReadUnsubscribeGroups(context)
		if err.Err != nil {
			return diag.FromErr(err.Err)
		}

		for i := range groups {
			group := groups[i]
			if group.Name == name {
				d.SetId(fmt.Sprint(group.ID))

				if err := sendgridUnsubscribeGroupParse(&group, d); err != nil {
					return diag.FromErr(err)
				}

				return nil
			}
		}

		return diag.Errorf("unable to find a unsubscribe group with name '%s'", name)
	default:
		return diag.Errorf("either 'group_id' or 'name' must be specified for data.sendgrid_unsubscribe_group")
	}
}
