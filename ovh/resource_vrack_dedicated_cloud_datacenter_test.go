package ovh

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/ovh/go-ovh/ovh"
	"github.com/ovh/terraform-provider-ovh/v2/ovh/types"
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
	targetServiceName := os.Getenv("OVH_VRACK_TARGET_SERVICE_TEST")
	if targetServiceName == "" {
		log.Print("[DEBUG] OVH_VRACK_TARGET_SERVICE_TEST is not set. No vrack_dedicated_cloud_datacenter to sweep")
		return nil
	}

	dedicatedCloudDatacenter := os.Getenv("OVH_DEDICATED_CLOUD_DATACENTER_TEST")
	if dedicatedCloudDatacenter == "" {
		log.Print("[DEBUG] OVH_DEDICATED_CLOUD_DATACENTER_TEST is not set. No vrack_dedicated_cloud_datacenter to sweep")
		return nil
	}

	// If the datacenter is already on the target_vrack, then move it to the source vrack
	// check if datacenter is already on target_vrack
	endpoint := fmt.Sprintf("/vrack/%s/dedicatedCloudDatacenter/%s",
		url.PathEscape(targetServiceName),
		url.PathEscape(dedicatedCloudDatacenter),
	)

	var task VrackTask
	var result interface{}
	if err := client.Get(endpoint, &result); err != nil {
		if errOvh, ok := err.(*ovh.APIError); ok && errOvh.Code == 404 {
			// datacenter is not on the target_vrack
			return nil
		}
		return err
	}

	// datacenter found on target_vrack, moving it to source vrack
	moveEndpoint := fmt.Sprintf("/vrack/%s/dedicatedCloudDatacenter/%s/move",
		url.PathEscape(targetServiceName),
		url.PathEscape(dedicatedCloudDatacenter))

	toCreatePayload := &VrackDedicatedCloudDatacenterModel{TargetServiceName: types.NewTfStringValue(serviceName)}
	moveErr := client.Post(moveEndpoint, toCreatePayload, &task)
	if moveErr != nil {
		return fmt.Errorf("Error trying to move dedicatedCloudDatacenter (%s) from vrack (%s) to vrack (%s): %s", dedicatedCloudDatacenter, targetServiceName, serviceName, moveErr)
	}

	if err := waitForVrackTask(&task, client); err != nil {
		return fmt.Errorf("Error waiting for vrack (%s) task to move dedicatedCloudDatacenter to vrack (%s): %s", serviceName, targetServiceName, err)
	}

	return nil
}

func TestAccVrackDedicatedCloudDatacenter_basic(t *testing.T) {
	resourceName := "ovh_vrack_dedicated_cloud_datacenter.vrack-dedicatedCloudDatacenter"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckDedicatedCloudDatacenter(t)
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				Config:            testAccDedicatedCloudDatacenterConfig,
				ImportStateIdFunc: func(state *terraform.State) (string, error) {
					return fmt.Sprintf("%s/%s/%s", os.Getenv("OVH_VRACK_SERVICE_TEST"), os.Getenv("OVH_VRACK_DEDICATED_CLOUD_DATACENTER"), os.Getenv("OVH_VRACK_TARGET_SERVICE_TEST")), nil
				},
			},
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
