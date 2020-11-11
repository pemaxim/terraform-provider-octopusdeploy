package octopusdeploy

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/pemaxim/go-octopusdeploy/octopusdeploy"
)

func expandEnvironment(d *schema.ResourceData) *octopusdeploy.Environment {
	name := d.Get("name").(string)

	environment := octopusdeploy.NewEnvironment(name)
	environment.ID = d.Id()

	if v, ok := d.GetOk(constAllowDynamicInfrastructure); ok {
		environment.AllowDynamicInfrastructure = v.(bool)
	}

	if v, ok := d.GetOk("description"); ok {
		environment.Description = v.(string)
	}

	if v, ok := d.GetOk(constSortOrder); ok {
		environment.SortOrder = v.(int)
	}

	if v, ok := d.GetOk(constUseGuidedFailure); ok {
		environment.UseGuidedFailure = v.(bool)
	}

	return environment
}

func flattenEnvironment(ctx context.Context, d *schema.ResourceData, environment *octopusdeploy.Environment) {
	d.Set("allow_dynamic_infrastructure", environment.AllowDynamicInfrastructure)
	d.Set("description", environment.Description)
	d.Set("name", environment.Name)
	d.Set("sort_order", environment.SortOrder)
	d.Set("use_guided_failure", environment.UseGuidedFailure)

	d.SetId(environment.GetID())
}

func getEnvironmentSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"allow_dynamic_infrastructure": {
			Optional: true,
			Type:     schema.TypeBool,
		},
		"description": {
			Optional: true,
			Type:     schema.TypeString,
		},
		"name": {
			Required:     true,
			Type:         schema.TypeString,
			ValidateFunc: validation.StringIsNotEmpty,
		},
		"sort_order": {
			Computed: true,
			Type:     schema.TypeInt,
		},
		"use_guided_failure": {
			Optional: true,
			Type:     schema.TypeBool,
		},
	}
}
