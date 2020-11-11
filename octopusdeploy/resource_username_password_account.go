package octopusdeploy

import (
	"context"

	"github.com/pemaxim/go-octopusdeploy/octopusdeploy"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceUsernamePassword() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceUsernamePasswordCreate,
		DeleteContext: resourceAccountDeleteCommon,
		ReadContext:   resourceUsernamePasswordRead,
		Schema:        getUsernamePasswordAccountSchema(),
		UpdateContext: resourceUsernamePasswordUpdate,
	}
}

func resourceUsernamePasswordCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	account := expandUsernamePasswordAccount(d)

	client := m.(*octopusdeploy.Client)
	createdAccount, err := client.Accounts.Add(account)
	if err != nil {
		diag.FromErr(err)
	}

	createdUsernamePasswordAccount := createdAccount.(*octopusdeploy.UsernamePasswordAccount)

	flattenUsernamePasswordAccount(ctx, d, createdUsernamePasswordAccount)
	return nil
}

func resourceUsernamePasswordRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*octopusdeploy.Client)
	accountResource, err := client.Accounts.GetByID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	usernamePasswordAccount := accountResource.(*octopusdeploy.UsernamePasswordAccount)

	flattenUsernamePasswordAccount(ctx, d, usernamePasswordAccount)
	return nil
}

func resourceUsernamePasswordUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	account := expandUsernamePasswordAccount(d)

	client := m.(*octopusdeploy.Client)
	updatedAccount, err := client.Accounts.Update(account)
	if err != nil {
		diag.FromErr(err)
	}

	updatedUsernamePasswordAccount := updatedAccount.(*octopusdeploy.UsernamePasswordAccount)

	flattenUsernamePasswordAccount(ctx, d, updatedUsernamePasswordAccount)
	return nil
}
