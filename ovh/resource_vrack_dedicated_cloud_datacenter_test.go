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
	resource.AddTestSweepers("ovh_vrack_dedicated_cloud_datacenter", &resource.Sweeper{
		Name: "ovh_vrack_dedicated_cloud_datacenter",
		F:    testSweepVrackDedicatedCloudDatacenter,
	})
}

func testSweepVrackDedicatedCloudDatacenter(region string) error {
	client, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting client: %s", err)
	}
	serviceName := os.Getenv("OVH_VRACK_SERVICE_TEST")
	if serviceName == "" {
		log.Print("[DEBUG] OVH_VRACK_SERVICE_TEST is not set. No vrack_dedicated_cloud_datacenter to sweep")
		return nil
	}

	dedicatedCloudDatacenter := os.Getenv("OVH_DEDICATED_CLOUD_DATACENTER_TEST")
	if dedicatedCloudDatacenter == "" {
		log.Print("[DEBUG] OVH_DEDICATED_CLOUD_DATACENTER_TEST is not set. No vrack_dedicated_cloud_datacenter to sweep")
		return nil
	}

	endpoint := fmt.Sprintf("/vrack/%s/dedicatedCloudDatacenter/%s",
		url.PathEscape(serviceName),
		url.PathEscape(dedicatedCloudDatacenter),
	)

	var result interface{}
	if err := client.Get(endpoint, &result); err != nil {
		if errOvh, ok := err.(*ovh.APIError); ok && errOvh.Code == 404 {
			return nil
		}
		return err
	}

	task := VrackTask{}
	if err := client.Delete(endpoint, &task); err != nil {
		return fmt.Errorf("Error calling DELETE %s with %s/%s:\n\t %q", endpoint, serviceName, dedicatedCloudDatacenter, err)
	}
	if err := waitForVrackTask(&task, client); err != nil {
		return fmt.Errorf("Error waiting for vrack (%s) to detach dedicatedCloudDatacenter (%s): %s", serviceName, dedicatedCloudDatacenter, err)
	}

	return nil
}

func TestAccVrackDedicatedCloudDatacenter_basic(t *testing.T) {
	serviceName := os.Getenv("OVH_VRACK_SERVICE_TEST")
	dedicatedCloudDatacenter := os.Getenv("OVH_VRACK_DEDICATED_CLOUD_DATACENTER")

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckDedicatedCloudDatacenter(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
				resource "ovh_vrack_dedicated_cloud_datacenter" "vrack-dedicatedCloudDatacenter" {
				  service_name      	= "%s"
				  datacenter 			= "%s"
				}
				`, serviceName, dedicatedCloudDatacenter),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ovh_vrack_dedicated_cloud_datacenter.vrack-dedicatedCloudDatacenter", "service_name", serviceName),
					resource.TestCheckResourceAttr("ovh_vrack_dedicated_cloud_datacenter.vrack-dedicatedCloudDatacenter", "datacenter", dedicatedCloudDatacenter),
				),
			},
		},
	})
}
