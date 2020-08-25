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
	resource.AddTestSweepers("ovh_me_ipxe_script", &resource.Sweeper{
		Name: "ovh_me_ipxe_script",
		F:    testSweepMeIpxeScript,
	})
}

func testSweepMeIpxeScript(region string) error {
	client, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting client: %s", err)
	}

	names := []string{}
	if err := client.Get("/me/ipxeScript", &names); err != nil {
		return fmt.Errorf("Error calling /me/ipxeScript:\n\t %q", err)
	}

	if len(names) == 0 {
		log.Print("[DEBUG] No IpxeScript to sweep")
		return nil
	}

	for _, scriptName := range names {
		if !strings.HasPrefix(scriptName, test_prefix) {
			continue
		}

		log.Printf("[DEBUG] IpxeScript found %v", scriptName)
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			log.Printf("[INFO] Deleting IpxeScript %v", scriptName)
			if err := client.Delete(fmt.Sprintf("/me/ipxeScript/%s", scriptName), nil); err != nil {
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

func TestAccMeIpxeScript_basic(t *testing.T) {
	scriptName := acctest.RandomWithPrefix(test_prefix)
	config := fmt.Sprintf(testAccMeIpxeScriptConfig, scriptName)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckCredentials(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"ovh_me_ipxe_script.script", "name", scriptName),
					resource.TestCheckResourceAttr(
						"ovh_me_ipxe_script.script", "script", "test"),
				),
			},
		},
	})
}

const testAccMeIpxeScriptConfig = `
resource "ovh_me_ipxe_script" "script" {
  name        = "%s"
  script      = "test"
  description = "fake description"
}
`
