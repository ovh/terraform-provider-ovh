package ovh

import (
	"fmt"
	"log"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func init() {
	resource.AddTestSweepers("ovh_me_installation_template", &resource.Sweeper{
		Name: "ovh_me_installation_template",
		F:    testSweepMeInstallationTemplate,
	})
}

func testSweepMeInstallationTemplate(region string) error {
	client, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting client: %s", err)
	}

	ids := []string{}
	if err := client.Get("/me/installationTemplate", &ids); err != nil {
		return fmt.Errorf("Error calling /me/installationTemplate:\n\t %q", err)
	}

	if len(ids) == 0 {
		log.Print("[DEBUG] No installation template to sweep")
		return nil
	}

	for _, id := range ids {
		if !strings.HasPrefix(id, test_prefix) {
			continue
		}

		log.Printf("[DEBUG] template found %v", id)
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			log.Printf("[INFO] Deleting template %v", id)
			if err := client.Delete(fmt.Sprintf("/me/installationTemplate/%s", id), nil); err != nil {
				return resource.RetryableError(err)
			}
			// Successful delete
			return nil
		})
		if err != nil {
			return err
		}
	}
	return nil
}

func TestAccMeInstallationTemplateResource_basic(t *testing.T) {
	templateName := acctest.RandomWithPrefix(test_prefix)
	config := fmt.Sprintf(
		testAccMeInstallationTemplateResourceConfig_Basic,
		templateName,
	)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckCredentials(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"ovh_me_installation_template.template",
						"template_name",
						templateName,
					),
					resource.TestCheckResourceAttrSet(
						"ovh_me_installation_template.template",
						"family",
					),
				),
			},
		},
	})
}

func TestAccMeInstallationTemplateResource_customization(t *testing.T) {
	templateName := acctest.RandomWithPrefix(test_prefix)
	config := fmt.Sprintf(
		testAccMeInstallationTemplateResourceConfig_Customization,
		templateName,
	)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckCredentials(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"ovh_me_installation_template.template",
						"template_name",
						templateName,
					),

					resource.TestCheckResourceAttr(
						"ovh_me_installation_template.template",
						"customization.0.custom_hostname",
						"mytest",
					),
					resource.TestCheckResourceAttr(
						"ovh_me_installation_template.template",
						"customization.0.ssh_key_name",
						"test",
					),
					resource.TestCheckResourceAttr(
						"ovh_me_installation_template.template",
						"customization.0.post_installation_script_link",
						"http://mylink.org",
					),
					resource.TestCheckResourceAttr(
						"ovh_me_installation_template.template",
						"customization.0.post_installation_script_return",
						"returned_string",
					),
					resource.TestCheckResourceAttrSet(
						"ovh_me_installation_template.template",
						"family",
					),
				),
			},
		},
	})
}

const testAccMeInstallationTemplateResourceConfig_Basic = `
resource "ovh_me_installation_template" "template" {
  base_template_name = "centos7_64"
  template_name      = "%s"
  default_language   = "en"
}
`
const testAccMeInstallationTemplateResourceConfig_Customization = `
resource "ovh_me_installation_template" "template" {
  base_template_name = "centos7_64"
  template_name      = "%s"
  default_language   = "en"

  customization {
     custom_hostname                 = "mytest"
     ssh_key_name                    = "test"
     post_installation_script_link   = "http://mylink.org"
     post_installation_script_return = "returned_string"
  }
}
`
