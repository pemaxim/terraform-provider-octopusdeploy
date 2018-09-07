package octopusdeploy

import (
	"fmt"
	"testing"

	"github.com/MattHodge/go-octopusdeploy/octopusdeploy"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccOctopusDeployEnvironmentBasic(t *testing.T) {
	const envPrefix = "octopusdeploy_environment.foo"
	const envName = "Testing one two three"
	const envDesc = "Terraform testing module environment"
	const envGuided = "false"
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckOctopusDeployProjectDestroy,
		Steps: []resource.TestStep{
			{
				Config: testEnvironmenttBasic(envName, envDesc, envGuided),
				Check: resource.ComposeTestCheckFunc(
					testOctopusDeployEnvironmentExists(envPrefix),
					resource.TestCheckResourceAttr(
						envPrefix, "name", envName),
					resource.TestCheckResourceAttr(
						envPrefix, "description", envDesc),
					resource.TestCheckResourceAttr(
						envPrefix, "useguidedfailure", envGuided),
				),
			},
		},
	})
}

func testEnvironmenttBasic(name, description, useguided string) string {
	return fmt.Sprintf(`
		resource "octopusdeploy_project" "foo" {
			name           = "%s"
			description    = "%s"
			useguidedfailure = "%s"
		}
		`,
		name, description, useguided,
	)
}

func testOctopusDeployEnvironmentExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*octopusdeploy.Client)
		return existsEnvHelper(s, client)
	}
}

func existsEnvHelper(s *terraform.State, client *octopusdeploy.Client) error {
	for _, r := range s.RootModule().Resources {
		if _, err := client.Environment.Get(r.Primary.ID); err != nil {
			return fmt.Errorf("Received an error retrieving environment %s", err)
		}
	}
	return nil
}
