package octopusdeploy

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/pemaxim/go-octopusdeploy/octopusdeploy"
)

func dataSourceTokenAccount() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceTokenAccountRead,
		Schema:      getTokenAccountDataSchema(),
	}
}

func dataSourceTokenAccountRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
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

	tokenAccount := accounts.Items[0].(*octopusdeploy.TokenAccount)

	flattenTokenAccount(ctx, d, tokenAccount)
	return nil
}
