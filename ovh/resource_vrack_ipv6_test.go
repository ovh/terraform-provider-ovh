package ovh

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/ovh/go-ovh/ovh"
)

func init() {
	resource.AddTestSweepers("ovh_vrack_ipv6", &resource.Sweeper{
		Name: "ovh_vrack_ipv6",
		F:    testSweepVrackIPv6,
	})
}

func testSweepVrackIPv6(region string) error {
	client, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting client: %s", err)
	}

	vrackId := os.Getenv("OVH_VRACK_SERVICE_TEST")
	if vrackId == "" {
		log.Print("[DEBUG] OVH_VRACK_SERVICE_TEST is not set. No vrack_ipv6 to sweep")
		return nil
	}

	ipBlock := os.Getenv("OVH_IP_V6_BLOCK_TEST")
	if ipBlock == "" {
		log.Print("[DEBUG] OVH_CLOUD_PROJECT_SERVICE_TEST is not set. No vrack_ipv6 to sweep")
		return nil
	}

	endpoint := fmt.Sprintf("/vrack/%s/ipv6/%s",
		url.PathEscape(vrackId),
		url.PathEscape(ipBlock),
	)

	if err := client.Get(endpoint, nil); err != nil {
		if errOvh, ok := err.(*ovh.APIError); ok && errOvh.Code == 404 {
			return nil
		}
		return err
	}

	task := VrackTask{}
	if err := client.Delete(endpoint, &task); err != nil {
		return fmt.Errorf("Error calling DELETE %s with %s/%s:\n\t %q", endpoint, vrackId, ipBlock, err)
	}

	if err := waitForVrackTask(&task, client); err != nil {
		return fmt.Errorf("Error waiting for vrack (%s) to detach ipv6 (%s): %s", vrackId, ipBlock, err)
	}

	return nil
}

func TestAccVrackIPv6_basic(t *testing.T) {
	serviceName := os.Getenv("OVH_VRACK_SERVICE_TEST")
	ipBlock := os.Getenv("OVH_IP_V6_BLOCK_TEST")

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckOCCVRack(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
				resource "ovh_vrack_ipv6" "vrack-ipv6" {
				  service_name = "%s"
				  block        = "%s"
				}
				`, serviceName, ipBlock),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ovh_vrack_ipv6.vrack-ipv6", "service_name", serviceName),
					resource.TestCheckResourceAttr("ovh_vrack_ipv6.vrack-ipv6", "block", ipBlock),
				),
			},
		},
	})
}
