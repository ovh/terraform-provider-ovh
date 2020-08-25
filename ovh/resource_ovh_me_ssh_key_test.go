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
	resource.AddTestSweepers("ovh_me_ssh_key", &resource.Sweeper{
		Name: "ovh_me_ssh_key",
		F:    testSweepMeSshKey,
	})
}

func testSweepMeSshKey(region string) error {
	client, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting client: %s", err)
	}

	names := []string{}
	if err := client.Get("/me/sshKey", &names); err != nil {
		return fmt.Errorf("Error calling /me/sshKey:\n\t %q", err)
	}

	if len(names) == 0 {
		log.Print("[DEBUG] No SSH keys to sweep")
		return nil
	}

	for _, keyName := range names {
		if !strings.HasPrefix(keyName, test_prefix) {
			continue
		}

		log.Printf("[DEBUG] SSH key found %v", keyName)
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			log.Printf("[INFO] Deleting SSH key %v", keyName)
			if err := client.Delete(fmt.Sprintf("/me/sshKey/%s", keyName), nil); err != nil {
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

func TestAccMeSshKey_basic(t *testing.T) {
	sshKeyName := acctest.RandomWithPrefix(test_prefix)
	sshKey := "ssh-ed25519 AAAAC3NzaC1yyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyy"
	config := fmt.Sprintf(testAccMeSshKeyConfig, sshKeyName, sshKey)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckCredentials(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"ovh_me_ssh_key.key_1", "key_name", sshKeyName),
					resource.TestCheckResourceAttr(
						"ovh_me_ssh_key.key_1", "key", sshKey),
				),
			},
		},
	})
}

const testAccMeSshKeyConfig = `
resource "ovh_me_ssh_key" "key_1" {
	key_name = "%s"
	key      = "%s"
}
`
