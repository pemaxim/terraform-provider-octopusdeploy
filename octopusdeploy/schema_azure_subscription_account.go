package octopusdeploy

import (
	"context"
	"time"

	"github.com/pemaxim/go-octopusdeploy/octopusdeploy"
	uuid "github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func expandAzureSubscriptionAccount(d *schema.ResourceData) *octopusdeploy.AzureSubscriptionAccount {
	name := d.Get("name").(string)
	subscriptionID, _ := uuid.Parse(d.Get("subscription_id").(string))

	account, _ := octopusdeploy.NewAzureSubscriptionAccount(name, subscriptionID)
	account.ID = d.Id()

	if v, ok := d.GetOk("azure_environment"); ok {
		account.AzureEnvironment = v.(string)
	}

	if v, ok := d.GetOk("certificate"); ok {
		account.CertificateBytes = octopusdeploy.NewSensitiveValue(v.(string))
	}

	if v, ok := d.GetOk("certificate_thumbprint"); ok {
		account.CertificateThumbprint = v.(string)
	}

	if v, ok := d.GetOk("description"); ok {
		account.Description = v.(string)
	}

	if v, ok := d.GetOk("environments"); ok {
		account.EnvironmentIDs = getSliceFromTerraformTypeList(v)
	}

	if v, ok := d.GetOk("management_endpoint"); ok {
		account.ManagementEndpoint = v.(string)
	}

	if v, ok := d.GetOk("modified_by"); ok {
		account.ModifiedBy = v.(string)
	}

	if v, ok := d.GetOk("modified_on"); ok {
		modifiedOnTime, _ := time.Parse(time.RFC3339, v.(string))
		account.ModifiedOn = &modifiedOnTime
	}

	if v, ok := d.GetOk("name"); ok {
		account.Name = v.(string)
	}

	if v, ok := d.GetOk("space_id"); ok {
		account.SpaceID = v.(string)
	}

	if v, ok := d.GetOk("storage_endpoint_suffix"); ok {
		account.StorageEndpointSuffix = v.(string)
	}

	if v, ok := d.GetOk("subscription_id"); ok {
		subscriptionID, _ := uuid.Parse(v.(string))
		account.SubscriptionID = &subscriptionID
	}

	if v, ok := d.GetOk("tenanted_deployment_participation"); ok {
		account.TenantedDeploymentMode = octopusdeploy.TenantedDeploymentMode(v.(string))
	}

	if v, ok := d.GetOk("tenant_tags"); ok {
		account.TenantTags = getSliceFromTerraformTypeList(v)
	}

	if v, ok := d.GetOk("tenants"); ok {
		account.TenantIDs = getSliceFromTerraformTypeList(v)
	}

	return account
}

func flattenAzureSubscriptionAccount(ctx context.Context, d *schema.ResourceData, account *octopusdeploy.AzureSubscriptionAccount) {
	flattenAccount(ctx, d, account)

	d.Set("account_type", "AzureSubscription")
	d.Set("azure_environment", account.AzureEnvironment)
	d.Set("certificate", account.CertificateBytes)
	d.Set("certificate_thumbprint", account.CertificateThumbprint)
	d.Set("management_endpoint", account.ManagementEndpoint)
	d.Set("storage_endpoint_suffix", account.StorageEndpointSuffix)
	d.Set("subscription_id", account.SubscriptionID.String())

	d.SetId(account.GetID())
}

func getAzureSubscriptionAccountDataSchema() map[string]*schema.Schema {
	schemaMap := getAccountDataSchema()
	schemaMap["account_type"] = &schema.Schema{
		Optional: true,
		Default:  "AzureSubscription",
		Type:     schema.TypeString,
	}
	schemaMap["azure_environment"] = getAzureEnvironmentDataSchema()
	schemaMap["certificate"] = &schema.Schema{
		Computed:  true,
		Sensitive: true,
		Type:      schema.TypeString,
	}
	schemaMap["certificate_thumbprint"] = &schema.Schema{
		Computed:  true,
		Sensitive: true,
		Type:      schema.TypeString,
	}
	schemaMap["management_endpoint"] = &schema.Schema{
		Computed: true,
		Type:     schema.TypeString,
	}
	schemaMap["storage_endpoint_suffix"] = &schema.Schema{
		Computed: true,
		Type:     schema.TypeString,
	}
	schemaMap["subscription_id"] = &schema.Schema{
		Computed: true,
		Type:     schema.TypeString,
	}
	return schemaMap
}

func getAzureSubscriptionAccountSchema() map[string]*schema.Schema {
	schemaMap := getAccountSchema()
	schemaMap["account_type"] = &schema.Schema{
		Optional: true,
		Default:  "AzureSubscription",
		Type:     schema.TypeString,
	}
	schemaMap["azure_environment"] = getAzureEnvironmentSchema()
	schemaMap["certificate"] = &schema.Schema{
		Optional:  true,
		Sensitive: true,
		Type:      schema.TypeString,
	}
	schemaMap["certificate_thumbprint"] = &schema.Schema{
		Optional:  true,
		Sensitive: true,
		Type:      schema.TypeString,
	}
	schemaMap["management_endpoint"] = &schema.Schema{
		Optional: true,
		Type:     schema.TypeString,
	}
	schemaMap["storage_endpoint_suffix"] = &schema.Schema{
		Optional: true,
		Type:     schema.TypeString,
	}
	schemaMap["subscription_id"] = &schema.Schema{
		Required:         true,
		Type:             schema.TypeString,
		ValidateDiagFunc: validateDiagFunc(validation.IsUUID),
	}
	return schemaMap
}
