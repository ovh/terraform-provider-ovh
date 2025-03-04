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
	resource.AddTestSweepers("ovh_vrack_ovhcloudconnect", &resource.Sweeper{
		Name: "ovh_vrack_ovhcloudconnect",
		F:    testSweepVrackOCC,
	})
}

func testSweepVrackOCC(region string) error {
	client, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting client: %s", err)
	}
	serviceName := os.Getenv("OVH_VRACK_SERVICE_TEST")
	if serviceName == "" {
		log.Print("[DEBUG] OVH_VRACK_SERVICE_TEST is not set. No vrack_ovhcloudconnect to sweep")
	}

	occ := os.Getenv("OVH_OVH_CLOUD_CONNECT_TEST")
	if occ == "" {
		log.Print("[DEBUG] OVH_OVH_CLOUD_CONNECT_TEST is not set. No vrack_ovhcloudconnect to sweep")
	}

	endpoint := fmt.Sprintf("/vrack/%s/ovhCloudConnect/%s",
		url.PathEscape(serviceName),
		url.PathEscape(occ),
	)

	if err := client.Get(endpoint, nil); err != nil {
		if errOvh, ok := err.(*ovh.APIError); ok && errOvh.Code == 404 {
			return nil
		}
		return err
	}

	task := VrackTask{}
	if err := client.Delete(endpoint, nil); err != nil {
		return fmt.Errorf("Error calling DELETE %s with %s/%s:\n\t %q", endpoint, serviceName, occ, err)
	}
	if err := waitForVrackTask(&task, client); err != nil {
		return fmt.Errorf("Error waiting for vrack (%s) to detach occ (%s): %s", serviceName, occ, err)
	}

	return nil
}

func TestAccVrackOCC_basic(t *testing.T) {
	serviceName := os.Getenv("OVH_VRACK_SERVICE_TEST")
	occ := os.Getenv("OVH_OVH_CLOUD_CONNECT_TEST")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheckCredentials(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
				resource "ovh_vrack_ovhcloudconnect" "vrack-occ" {
				  service_name      = "%s"
				  ovh_cloud_connect = "%s"
				}
				`, serviceName, occ),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ovh_vrack_ovhcloudconnect.vrack-occ", "service_name", serviceName),
					resource.TestCheckResourceAttr("ovh_vrack_ovhcloudconnect.vrack-occ", "ovh_cloud_connect", occ),
				),
			},
		},
	})
}
