package ovh

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"regexp"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/ovh/go-ovh/ovh"
)

var testAccIpReverseConfig = `
resource "ovh_ip_reverse" "reverse" {
    ip = "%s"
    ip_reverse = "%s"
    reverse = "%s"
	readiness_timeout_duration = %q
}
`

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

	reverse := IpReverse{}
	testIp := os.Getenv("OVH_IP_BLOCK_TEST")
	testIpReverse := os.Getenv("OVH_IP_TEST")
	endpoint := fmt.Sprintf("/ip/%s/reverse/%s",
		url.PathEscape(testIp),
		url.PathEscape(testIpReverse),
	)
	if err := client.Get(endpoint, &reverse); err != nil {
		if errOvh, ok := err.(*ovh.APIError); ok && errOvh.Code == 404 {
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
	block := os.Getenv("OVH_IP_BLOCK_TEST")
	ip := os.Getenv("OVH_IP_TEST")
	reverse := os.Getenv("OVH_IP_REVERSE_TEST")

	config := fmt.Sprintf(testAccIpReverseConfig, block, ip, reverse, "60s")

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckIp(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ovh_ip_reverse.reverse", "ip", block),
					resource.TestCheckResourceAttr("ovh_ip_reverse.reverse", "ip_reverse", ip),
					resource.TestCheckResourceAttr("ovh_ip_reverse.reverse", "reverse", reverse),
					resource.TestCheckResourceAttr("ovh_ip_reverse.reverse", "readiness_timeout_duration", "60s"),
				),
			},
			{
				ResourceName:        "ovh_ip_reverse.reverse",
				ImportState:         true,
				ImportStateIdPrefix: block + "|",
				ImportStateVerify:   true,
			},
		},
	})
}

func TestAccIpReverse_invalidDuration(t *testing.T) {
	block := os.Getenv("OVH_IP_BLOCK_TEST")
	ip := os.Getenv("OVH_IP_TEST")
	reverse := os.Getenv("OVH_IP_REVERSE_TEST")

	config := fmt.Sprintf(testAccIpReverseConfig, block, ip, reverse, "1xyz")

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckIp(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      config,
				ExpectError: regexp.MustCompile(`cannot parse readiness_timeout_seconds attribute`),
			},
		},
	})
}

/*
func TestAccIpReverse_importBasic(t *testing.T) {
	block := os.Getenv("OVH_IP_BLOCK_TEST")
	ip := os.Getenv("OVH_IP_TEST")
	reverse := os.Getenv("OVH_IP_REVERSE_TEST")

	config := fmt.Sprintf(testAccIpReverseConfig, block, ip, reverse)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckIp(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: config,
			},
			{
				ResourceName:      "ovh_ip_reverse.reverse",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccIpReverseImportId("ovh_ip_reverse.reverse"),
			},
		},
	})
}

func testAccIpReverseImportId(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		subnet, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("resource not found: %s", resourceName)
		}

		return fmt.Sprintf(
			"%s|%s",
			subnet.Primary.Attributes["ip"],
			subnet.Primary.Attributes["ip_reverse"],
		), nil
	}
}
*/
