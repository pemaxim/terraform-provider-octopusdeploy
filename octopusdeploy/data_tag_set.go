package octopusdeploy

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/pemaxim/go-octopusdeploy/octopusdeploy"
)

func dataTagSet() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataTagSetReadByName,

		Schema: map[string]*schema.Schema{
			"name": {
				Required: true,
				Type:     schema.TypeString,
			},
			constTag: getTagSchema(),
		},
	}
}

func dataTagSetReadByName(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	name := d.Get("name").(string)

	client := m.(*octopusdeploy.Client)
	resource, err := client.TagSets.GetByName(name)
	if err != nil {
		return diag.FromErr(err)
	}
	if resource == nil {
		return nil
	}

	logResource(constTagSet, m)

	d.SetId(resource.GetID())
	d.Set("name", resource.Name)

	return nil
}
