package octopusdeploy

import (
	"fmt"
	"testing"

	"github.com/pemaxim/go-octopusdeploy/octopusdeploy"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestSSHKeyBasic(t *testing.T) {
	localName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	prefix := "octopusdeploy_ssh_key_account." + localName

	name := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	passphrase := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	tenantedDeploymentParticipation := octopusdeploy.TenantedDeploymentModeTenantedOrUntenanted
	username := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)

	resource.Test(t, resource.TestCase{
		CheckDestroy: testAccAccountCheckDestroy,
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testSSHKeyBasic(localName, name, username, passphrase, tenantedDeploymentParticipation),
				Check: resource.ComposeTestCheckFunc(
					testAccAccountExists(prefix),
					resource.TestCheckResourceAttr(prefix, "name", name),
					resource.TestCheckResourceAttr(prefix, constPassphrase, passphrase),
					resource.TestCheckResourceAttr(prefix, "tenanted_deployment_participation", string(tenantedDeploymentParticipation)),
					resource.TestCheckResourceAttr(prefix, "username", username),
				),
			},
		},
	})
}

func testSSHKeyBasic(localName string, name string, username string, passphrase string, tenantedDeploymentParticipation octopusdeploy.TenantedDeploymentMode) string {
	return fmt.Sprintf(`resource "octopusdeploy_ssh_key_account" "%s" {
		name = "%s"
		passphrase = "%s"
		tenanted_deployment_participation = "%s"
		username = "%s"
	}`, localName, name, passphrase, tenantedDeploymentParticipation, username)
}
