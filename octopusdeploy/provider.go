package octopusdeploy

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Provider is the plugin entry point
func Provider() *schema.Provider {
	return &schema.Provider{
		DataSourcesMap: map[string]*schema.Resource{
			"octopusdeploy_account":                    dataAccount(),
			"octopusdeploy_aws_account":                dataSourceAmazonWebServicesAccount(),
			"octopusdeploy_azure_service_principal":    dataSourceAzureServicePrincipalAccount(),
			"octopusdeploy_azure_subscription_account": dataSourceAzureSubscriptionAccount(),
			"octopusdeploy_certificate":                dataSourceCertificate(),
			"octopusdeploy_channel":                    dataSourceChannel(),
			"octopusdeploy_deployment_target":          dataSourceDeploymentTarget(),
			"octopusdeploy_environment":                dataSourceEnvironment(),
			"octopusdeploy_feed":                       dataFeed(),
			"octopusdeploy_library_variable_set":       dataSourceLibraryVariableSet(),
			"octopusdeploy_lifecycle":                  dataSourceLifecycle(),
			"octopusdeploy_machine_policy":             dataMachinePolicy(),
			"octopusdeploy_nuget_feed":                 dataSourceNuGetFeed(),
			"octopusdeploy_project":                    dataSourceProject(),
			"octopusdeploy_ssh_key_account":            dataSourceSSHKeyAccount(),
			"octopusdeploy_space":                      dataSourceSpace(),
			"octopusdeploy_token_account":              dataSourceTokenAccount(),
			"octopusdeploy_user":                       dataSourceUser(),
			"octopusdeploy_username_password_account":  dataSourceUsernamePasswordAccount(),
			"octopusdeploy_variable":                   dataVariable(),
		},
		ResourcesMap: map[string]*schema.Resource{
			"octopusdeploy_account":                           resourceAccount(),
			"octopusdeploy_aws_account":                       resourceAmazonWebServicesAccount(),
			"octopusdeploy_azure_service_principal":           resourceAzureServicePrincipalAccount(),
			"octopusdeploy_azure_subscription_account":        resourceAzureSubscriptionAccount(),
			"octopusdeploy_certificate":                       resourceCertificate(),
			"octopusdeploy_channel":                           resourceChannel(),
			"octopusdeploy_deployment_target":                 resourceDeploymentTarget(),
			"octopusdeploy_deployment_process":                resourceDeploymentProcess(),
			"octopusdeploy_environment":                       resourceEnvironment(),
			"octopusdeploy_feed":                              resourceFeed(),
			"octopusdeploy_library_variable_set":              resourceLibraryVariableSet(),
			"octopusdeploy_lifecycle":                         resourceLifecycle(),
			"octopusdeploy_nuget_feed":                        resourceNuGetFeed(),
			"octopusdeploy_project":                           resourceProject(),
			"octopusdeploy_project_deployment_target_trigger": resourceProjectDeploymentTargetTrigger(),
			"octopusdeploy_project_group":                     resourceProjectGroup(),
			"octopusdeploy_space":                             resourceSpace(),
			"octopusdeploy_ssh_key_account":                   resourceSSHKey(),
			"octopusdeploy_tag_set":                           resourceTagSet(),
			"octopusdeploy_token_account":                     resourceTokenAccount(),
			"octopusdeploy_user":                              resourceUser(),
			"octopusdeploy_username_password_account":         resourceUsernamePassword(),
			"octopusdeploy_variable":                          resourceVariable(),
		},
		Schema: map[string]*schema.Schema{
			constAddress: {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("OCTOPUS_URL", nil),
				Description: "The URL of the Octopus Deploy server",
			},
			constAPIKey: {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("OCTOPUS_APIKEY", nil),
				Description: "The API to use with the Octopus Deploy server.",
			},
			constSpaceID: {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("OCTOPUS_SPACE", constEmptyString),
				Description: "The name of the Space in Octopus Deploy server",
			},
		},

		ConfigureContextFunc: providerConfigure,
	}
}

func providerConfigure(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	config := Config{
		Address: d.Get(constAddress).(string),
		APIKey:  d.Get(constAPIKey).(string),
		Space:   d.Get(constSpaceID).(string),
	}

	return config.Client()
}
