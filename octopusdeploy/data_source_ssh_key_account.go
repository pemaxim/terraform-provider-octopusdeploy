package octopusdeploy

import (
	"context"

	"github.com/pemaxim/go-octopusdeploy/octopusdeploy"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceSSHKeyAccount() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceSSHKeyAccountRead,
		Schema:      getSSHKeyAccountDataSchema(),
	}
}

func dataSourceSSHKeyAccountRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*octopusdeploy.Client)
	name := d.Get("name").(string)
	query := octopusdeploy.AccountsQuery{
		PartialName: name,
		Take:        1,
	}

	accounts, err := client.Accounts.Get(query)
	if err != nil {
		return diag.FromErr(err)
	}
	if accounts == nil || len(accounts.Items) == 0 {
		d.SetId("")
		return diag.Errorf("unable to retrieve account (partial name: %s)", name)
	}

	sshKeyAccount := accounts.Items[0].(*octopusdeploy.SSHKeyAccount)

	flattenSSHKeyAccount(ctx, d, sshKeyAccount)
	return nil
}
