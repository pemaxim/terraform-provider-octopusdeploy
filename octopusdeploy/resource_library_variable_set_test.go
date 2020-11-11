package octopusdeploy

import (
	"fmt"
	"testing"

	"github.com/pemaxim/go-octopusdeploy/octopusdeploy"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccOctopusDeployLibraryVariableSetBasic(t *testing.T) {
	localName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	prefix := constOctopusDeployLibraryVariableSet + "." + localName

	name := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)

	resource.Test(t, resource.TestCase{
		CheckDestroy: testLibraryVariableSetDestroy,
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testLibraryVariableSetBasic(localName, name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOctopusDeployLibraryVariableSetExists(prefix),
					resource.TestCheckResourceAttr(prefix, "name", name),
				),
			},
		},
	})
}

func TestAccOctopusDeployLibraryVariableSetComplex(t *testing.T) {
	localName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	prefix := constOctopusDeployLibraryVariableSet + "." + localName

	name := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)

	resource.Test(t, resource.TestCase{
		CheckDestroy: testLibraryVariableSetDestroy,
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testLibraryVariableSetComplex(localName, name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOctopusDeployLibraryVariableSetExists(prefix),
					resource.TestCheckResourceAttr(prefix, "name", name),
				),
			},
		},
	})
}

func TestAccOctopusDeployLibraryVariableSetWithUpdate(t *testing.T) {
	localName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	prefix := constOctopusDeployLibraryVariableSet + "." + localName

	description := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	name := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)

	resource.Test(t, resource.TestCase{
		CheckDestroy: testLibraryVariableSetDestroy,
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		Steps: []resource.TestStep{
			// create variable set with no description
			{
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOctopusDeployLibraryVariableSetExists(prefix),
					resource.TestCheckResourceAttr(prefix, "name", name),
				),
				Config: testLibraryVariableSetBasic(localName, name),
			},
			// create update it with a description
			{
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOctopusDeployLibraryVariableSetExists(prefix),
					resource.TestCheckResourceAttr(prefix, "name", name),
					resource.TestCheckResourceAttr(prefix, "description", description),
				),
				Config: testAccLibraryVariableSetWithDescription(localName, name, description),
			},
			// update again by remove its description
			{
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOctopusDeployLibraryVariableSetExists(prefix),
					resource.TestCheckResourceAttr(prefix, "name", name),
					resource.TestCheckResourceAttr(prefix, "description", ""),
				),
				Config: testLibraryVariableSetBasic(localName, name),
			},
		},
	})
}

func testLibraryVariableSetBasic(localName string, name string) string {
	return fmt.Sprintf(`resource "%s" "%s" {
		name = "%s"
	}`, constOctopusDeployLibraryVariableSet, localName, name)
}

func testAccLibraryVariableSetWithDescription(localName string, name string, description string) string {
	return fmt.Sprintf(`resource "%s" "%s" {
		name        = "%s"
		description = "%s"
	}`, constOctopusDeployLibraryVariableSet, localName, name, description)
}

func testLibraryVariableSetComplex(localName string, name string) string {
	return fmt.Sprintf(`resource "%s" "%s" {
		description     = "This is the description."
		name            = "%s"
		template {
			default_value    = "Default Value???"
			display_settings = {
				"Octopus.ControlType" = "SingleLineText"
			}
			help_text        = "This is the help text"
			label            = "Test Label"
			name             = "Test Template"
		}
		template {
			default_value    = "wjehqwjkehwqkejhqwe"
			display_settings = {
				"Octopus.ControlType" = "MultiLineText"
			}
			help_text        = "jhasdkjashdaksjhd"
			label            = "alsdjhaldksh"
			name             = "Another Variable???"
		}
		template {
			default_value    = "qweq|qwe"
			display_settings = {
				"Octopus.ControlType" = "MultiLineText"
			}
			help_text        = "qwe"
			label            = "qwe"
			name             = "weq"
		}
	}`, constOctopusDeployLibraryVariableSet, localName, name)
}

func testLibraryVariableSetDestroy(s *terraform.State) error {
	if err := destroyHelperLibraryVariableSet(s); err != nil {
		return err
	}
	return nil
}

func testAccCheckOctopusDeployLibraryVariableSetExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*octopusdeploy.Client)
		if err := existsHelperLibraryVariableSet(s, client); err != nil {
			return err
		}
		return nil
	}
}

func destroyHelperLibraryVariableSet(s *terraform.State) error {
	client := testAccProvider.Meta().(*octopusdeploy.Client)
	for _, rs := range s.RootModule().Resources {
		libraryVariableSetID := rs.Primary.ID
		libraryVariableSet, err := client.LibraryVariableSets.GetByID(libraryVariableSetID)
		if err == nil {
			if libraryVariableSet != nil {
				return fmt.Errorf("library variable set (%s) still exists", rs.Primary.ID)
			}
		}
	}

	return nil
}

func existsHelperLibraryVariableSet(s *terraform.State, client *octopusdeploy.Client) error {
	for _, r := range s.RootModule().Resources {
		if r.Type == constOctopusDeployLibraryVariableSet {
			if _, err := client.LibraryVariableSets.GetByID(r.Primary.ID); err != nil {
				return fmt.Errorf("received an error retrieving library variable set %s", err)
			}
		}
	}
	return nil
}
