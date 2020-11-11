package octopusdeploy

import (
	"fmt"
	"testing"

	"github.com/pemaxim/go-octopusdeploy/octopusdeploy"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccOctopusDeployVariableBasic(t *testing.T) {
	localName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	prefix := constOctopusDeployVariable + "." + localName

	description := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	name := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	value := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)

	resource.Test(t, resource.TestCase{
		CheckDestroy: testVariableDestroy,
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testVariableBasic(localName, name, description, value),
				Check: resource.ComposeTestCheckFunc(
					testOctopusDeployVariableExists(prefix),
					resource.TestCheckResourceAttr(prefix, "name", name),
					resource.TestCheckResourceAttr(prefix, "description", description),
					resource.TestCheckResourceAttr(prefix, "value", value),
				),
			},
		},
	})
}

func testVariableBasic(localName string, name string, description string, value string) string {
	projectDescription := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	projectLocalName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	projectName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)

	config := fmt.Sprintf(testAccProjectBasic(projectLocalName, projectName, projectDescription)+"\n"+
		`resource "%s" "%s" {
			description = "%s"
			name        = "%s"
			project_id  = "${%s.%s.id}"
			type        = "String"
			value       = "%s"
		}`, constOctopusDeployVariable, localName, description, name, constOctopusDeployProject, projectLocalName, value)
	return config
}

func testOctopusDeployVariableExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		var projectID string
		var variableID string

		for _, r := range s.RootModule().Resources {
			if r.Type == constOctopusDeployProject {
				projectID = r.Primary.ID
			}

			if r.Type == constOctopusDeployVariable {
				variableID = r.Primary.ID
			}
		}

		client := testAccProvider.Meta().(*octopusdeploy.Client)
		if _, err := client.Variables.GetByID(projectID, variableID); err != nil {
			return fmt.Errorf("Received an error retrieving variable %s", err)
		}

		return nil
	}
}

func testVariableDestroy(s *terraform.State) error {
	var projectID string
	var variableID string

	for _, r := range s.RootModule().Resources {
		if r.Type == constOctopusDeployProject {
			projectID = r.Primary.ID
		}

		if r.Type == constOctopusDeployVariable {
			variableID = r.Primary.ID
		}
	}

	client := testAccProvider.Meta().(*octopusdeploy.Client)
	variable, err := client.Variables.GetByID(projectID, variableID)
	if err == nil {
		if variable != nil {
			return fmt.Errorf("variable (%s) still exists", variableID)
		}
	}

	return nil
}
