package octopusdeploy

import (
	"fmt"
	"testing"

	"github.com/pemaxim/go-octopusdeploy/octopusdeploy"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccOctopusDeployCertificateBasic(t *testing.T) {
	localName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	prefix := constOctopusDeployCertificate + "." + localName

	const certName = "Testing one two three"
	const certNotes = "Cert notes blah blah blah"
	const certData = "MIIDiDCCAnACCQDXHofnqz05ITANBgkqhkiG9w0BAQsFADCBhTELMAkGA1UEBhMCVVMxETAPBgNVBAgMCE9rbGFob21hMQ8wDQYDVQQHDAZOb3JtYW4xEzARBgNVBAoMCk1vb25zd2l0Y2gxGTAXBgNVBAMMEGRlbW8ub2N0b3B1cy5jb20xIjAgBgkqhkiG9w0BCQEWE2plZmZAbW9vbnN3aXRjaC5jb20wHhcNMTkwNjE0MjExMzI1WhcNMjAwNjEzMjExMzI1WjCBhTELMAkGA1UEBhMCVVMxETAPBgNVBAgMCE9rbGFob21hMQ8wDQYDVQQHDAZOb3JtYW4xEzARBgNVBAoMCk1vb25zd2l0Y2gxGTAXBgNVBAMMEGRlbW8ub2N0b3B1cy5jb20xIjAgBgkqhkiG9w0BCQEWE2plZmZAbW9vbnN3aXRjaC5jb20wggEiMA0GCSqGSIb3DQEBAQUAA4IBDwAwggEKAoIBAQDSTiD0OHyFDMH9O+d/h3AiqcuvpvUgRkKjf+whZ6mVlQnGkvPddRTUY48xCEaQ4QD1MAVJcGaJ2PU4NxwhrQgHqWW8TQkAZESL4wfzSwIKO2NX/I2tWqyv7a0uA/WdtlWQye+2oPV5rCnS0kM75X+gjEwOTpFh/ryS6KhMPFDb0zeNGREdg6564FdxWSvN4ppUZMqhvMpfzM7rsDWqEzYsMaQ4CNJDFdWkG89D4j5qk4b4Qb4m+l7QINdmYIXf4qO/0LE1WcfIkCpAS65tjc/hefIHmYtj/E/ijoNJbWKZDK3WLZg3zq99Ipqv/9DFvSiMQFBhZT0jO2B5d5zBUuIHAgMBAAEwDQYJKoZIhvcNAQELBQADggEBAKsa4gaW7GhByu8aq56h99DaIl1LauI5WMVH8Q9Qpapho2VLRIpfwGeI5eENFoXwuKrnJp1ncsCqGnMQnugQHS+SrruS3Yyl0Uog4Zak9GbbK6qn+olx7GNJbsckmD371lqQOaKITLqYzK6kTc7/v8Cv0BwHFCBda1OCrmeVBSaarucPxZhGxzLAielzHHdlkZFQT/oO2VR3thhURIqtni7jVQ2MoeZF1ccvmAfVbzr/QnlNe/jrcmyPYymuF2JyrezzIjlKuiDhalKqwqkCHpOOgzV4y6BFuS+0w3DS8pa07nUudZ6E0kZzvhjjiyAx/sBdX6ZDdUjP9TDJMM4f5YA="
	const tagSetName = "TagSet"
	const tagName = "Tag"
	const envName = "Test Env"
	var tenantTags = fmt.Sprintf("%s/%s", tagSetName, tagName)
	const tenantedDeploymentParticipation = "TenantedOrUntenanted"

	resource.Test(t, resource.TestCase{
		CheckDestroy: testCertificateDestroy,
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testCertificateBasic(localName, tagSetName, tagName, envName, certName, certNotes, certData, tenantedDeploymentParticipation),
				Check: resource.ComposeTestCheckFunc(
					testCertificateExists(prefix),
					resource.TestCheckResourceAttr(prefix, constCertificateData, certData),
					resource.TestCheckResourceAttr(prefix, constName, certName),
					resource.TestCheckResourceAttr(prefix, constNotes, certNotes),
					resource.TestCheckResourceAttr(prefix, "tenant_tags.0", tenantTags),
					resource.TestCheckResourceAttr(prefix, "tenanted_deployment_participation", tenantedDeploymentParticipation),
				),
			},
		},
	})
}

func testCertificateBasic(localName string, tagSetName string, tagName string, environmentName string, certName string, notes string, certificateData string, tenantedDeploymentParticipation string) string {
	return fmt.Sprintf(`
		resource "octopusdeploy_tag_set" "testtagset" {
			name = "%s"
			tag {
				name = "%s"
				color = "#6e6e6f"
			}
		}

		resource "octopusdeploy_environment" "test_env" {
			name = "%s"
		}

		resource "%s" "%s" {
			name = "%s"
			notes = "%s"
			certificate_data = "%s"
			environment_ids = ["${octopusdeploy_environment.test_env.id}"]
			tenanted_deployment_participation = "%s"
			tenant_tags = ["${octopusdeploy_tag_set.testtagset.name}/%s"]
		}`, tagSetName, tagName, environmentName, constOctopusDeployCertificate, localName, certName, notes, certificateData, tenantedDeploymentParticipation, tagName)
}

func testCertificateExists(prefix string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*octopusdeploy.Client)
		certificateID := s.RootModule().Resources[prefix].Primary.ID
		if _, err := client.Certificates.GetByID(certificateID); err != nil {
			return err
		}

		return nil
	}
}

func testCertificateDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*octopusdeploy.Client)
	for _, rs := range s.RootModule().Resources {
		certificateID := rs.Primary.ID
		certificate, err := client.Certificates.GetByID(certificateID)
		if err == nil {
			if certificate != nil {
				return fmt.Errorf("certificate (%s) still exists", rs.Primary.ID)
			}
		}
	}

	return nil
}
