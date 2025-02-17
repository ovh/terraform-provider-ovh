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
	resource.AddTestSweepers("ovh_vrack_vrackservices", &resource.Sweeper{
		Name: "ovh_vrack_vrackservices",
		F:    testSweepVrackVrackServices,
	})
}

func testSweepVrackVrackServices(region string) error {
	client, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting client: %s", err)
	}

	vrackId := os.Getenv("OVH_VRACK_SERVICE_TEST")
	if vrackId == "" {
		log.Print("[DEBUG] OVH_VRACK_SERVICE_TEST is not set. No vrack_vrackservices to sweep")
		return nil
	}

	vrackServices := os.Getenv("OVH_VRACKSERVICES_SERVICE_TEST")
	if vrackServices == "" {
		log.Print("[DEBUG] OVH_VRACKSERVICES_SERVICE_TEST is not set. No vrack_vrackservices to sweep")
		return nil
	}

	endpoint := fmt.Sprintf("/vrack/%s/vrackServices/%s",
		url.PathEscape(vrackId),
		url.PathEscape(vrackServices),
	)

	if err := client.Get(endpoint, nil); err != nil {
		if errOvh, ok := err.(*ovh.APIError); ok && errOvh.Code == 404 {
			return nil
		}
		return err
	}

	task := VrackTask{}
	if err := client.Delete(endpoint, &task); err != nil {
		return fmt.Errorf("Error calling DELETE %s with %s/%s:\n\t %q", endpoint, vrackId, vrackServices, err)
	}

	if err := waitForVrackTask(&task, client); err != nil {
		return fmt.Errorf("Error waiting for vrack (%s) to detach vrackServices (%s): %s", vrackId, vrackServices, err)
	}

	return nil
}

func TestAccVrackVrackServices_basic(t *testing.T) {
	serviceName := os.Getenv("OVH_VRACK_SERVICE_TEST")
	vrackServices := os.Getenv("OVH_VRACKSERVICES_SERVICE_TEST")

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckVRack(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
				resource "ovh_vrack_vrackservices" "vrack-vrackServices" {
				  service_name   = "%s"
				  vrack_services = "%s"
				}
				`, serviceName, vrackServices),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ovh_vrack_vrackservices.vrack-vrackServices", "service_name", serviceName),
					resource.TestCheckResourceAttr("ovh_vrack_vrackservices.vrack-vrackServices", "vrack_services", vrackServices),
				),
			},
		},
	})
}
