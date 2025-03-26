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

var testAccDedicatedCloudDatacenterConfig = fmt.Sprintf(`
resource "ovh_vrack_dedicated_cloud_datacenter" "vrack-dedicatedCloudDatacenter" {
  service_name     = "%s"
  datacenter = "%s"
}
`, os.Getenv("OVH_VRACK_SERVICE_TEST"), os.Getenv("OVH_VRACK_DEDICATED_CLOUD_DATACENTER"))

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

	return nil
}

func TestAccVrackDedicatedCloudDatacenter_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckDedicatedCloudDatacenter(t)
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDedicatedCloudDatacenterConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ovh_vrack_dedicated_cloud_datacenter.vrack-dedicatedCloudDatacenter", "service_name", os.Getenv("OVH_VRACK_TARGET_SERVICE_TEST")),
					resource.TestCheckResourceAttr("ovh_vrack_dedicated_cloud_datacenter.vrack-dedicatedCloudDatacenter", "datacenter", os.Getenv("OVH_VRACK_DEDICATED_CLOUD_DATACENTER")),
				),
			},
		},
	})
}
