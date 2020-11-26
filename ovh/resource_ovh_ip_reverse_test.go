package ovh

import (
	"fmt"
	"log"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/ovh/go-ovh/ovh"
)

var testAccIpReverseConfig = fmt.Sprintf(`
resource "ovh_ip_reverse" "reverse" {
    ip = "%s"
    ipreverse = "%s"
    reverse = "%s"
}
`, os.Getenv("OVH_IP_BLOCK"), os.Getenv("OVH_IP"), os.Getenv("OVH_IP_REVERSE"))

func init() {
	resource.AddTestSweepers("ovh_ip_reverse", &resource.Sweeper{
		Name: "ovh_ip_reverse",
		F:    testSweepIpReverse,
	})
}

func testSweepIpReverse(region string) error {
	client, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting client: %s", err)
	}

	reverse := OvhIpReverse{}
	testIp := os.Getenv("OVH_IP_BLOCK")
	testIpReverse := os.Getenv("OVH_IP")
	endpoint := fmt.Sprintf("/ip/%s/reverse/%s", strings.Replace(testIp, "/", "%2F", 1), testIpReverse)
	if err := client.Get(endpoint, &reverse); err != nil {
		if err.(*ovh.APIError).Code == 404 {
			// no ip reverse set, nothing to sweep
			return nil
		}

		return fmt.Errorf("Error calling %s:\n\t %q", endpoint, err)
	}

	log.Printf("[DEBUG] ip reverse found %v", reverse)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		log.Printf("[INFO] Deleting reverse %v", reverse)
		if err := client.Delete(endpoint, nil); err != nil {
			return resource.RetryableError(err)
		}
		// Successful delete
		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

func TestAccIpReverse_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckIp(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckIpReverseDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccIpReverseConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIpReverseExists("ovh_ip_reverse.reverse", t),
				),
			},
		},
	})
}

func testAccCheckIpReverseExists(n string, t *testing.T) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		config := testAccProvider.Meta().(*Config)

		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.Attributes["ip"] == "" {
			return fmt.Errorf("No IP block is set")
		}

		if rs.Primary.Attributes["ipreverse"] == "" {
			return fmt.Errorf("No IP is set")
		}

		return resourceOvhIpReverseExists(rs.Primary.Attributes["ip"], rs.Primary.Attributes["ipreverse"], config.OVHClient)
	}
}

func testAccCheckIpReverseDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "ovh_ip_reverse" {
			continue
		}

		err := resourceOvhIpReverseExists(rs.Primary.Attributes["ip"], rs.Primary.Attributes["ipreverse"], config.OVHClient)
		if err == nil {
			return fmt.Errorf("IP Reverse still exists")
		}
	}
	return nil
}
