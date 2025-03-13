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
	resource.AddTestSweepers("ovh_vrack_dedicated_cloud", &resource.Sweeper{
		Name: "ovh_vrack_dedicated_cloud",
		F:    testSweepVrackDedicatedCloud,
	})
}

func testSweepVrackDedicatedCloud(region string) error {
	client, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting client: %s", err)
	}
	serviceName := os.Getenv("OVH_VRACK_SERVICE_TEST")
	if serviceName == "" {
		log.Print("[DEBUG] OVH_VRACK_SERVICE_TEST is not set. No vrack_dedicated_cloud to sweep")
		return nil
	}

	dedicatedCloud := os.Getenv("OVH_DEDICATED_CLOUD_TEST")
	if dedicatedCloud == "" {
		log.Print("[DEBUG] OVH_DEDICATED_CLOUD_TEST is not set. No vrack_dedicated_cloud to sweep")
		return nil
	}

	endpoint := fmt.Sprintf("/vrack/%s/dedicatedCloud/%s",
		url.PathEscape(serviceName),
		url.PathEscape(dedicatedCloud),
	)

	if err := client.Get(endpoint, nil); err != nil {
		if errOvh, ok := err.(*ovh.APIError); ok && errOvh.Code == 404 {
			return nil
		}
		return err
	}

	task := VrackTask{}
	if err := client.Delete(endpoint, &task); err != nil {
		return fmt.Errorf("Error calling DELETE %s with %s/%s:\n\t %q", endpoint, serviceName, dedicatedCloud, err)
	}
	if err := waitForVrackTask(&task, client); err != nil {
		return fmt.Errorf("Error waiting for vrack (%s) to detach dedicatedCloud (%s): %s", serviceName, dedicatedCloud, err)
	}

	return nil
}

func TestAccVrackDedicatedCloud_basic(t *testing.T) {
	serviceName := os.Getenv("OVH_VRACK_SERVICE_TEST")
	dedicatedCloud := os.Getenv("OVH_VRACK_DEDICATED_CLOUD")

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckDedicatedCloud(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
				resource "ovh_vrack_dedicated_cloud" "vrack-dedicatedCloud" {
				  service_name      = "%s"
				  dedicated_cloud 	= "%s"
				}
				`, serviceName, dedicatedCloud),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ovh_vrack_dedicated_cloud.vrack-dedicatedCloud", "service_name", serviceName),
					resource.TestCheckResourceAttr("ovh_vrack_dedicated_cloud.vrack-dedicatedCloud", "dedicated_cloud", dedicatedCloud),
				),
			},
		},
	})
}
