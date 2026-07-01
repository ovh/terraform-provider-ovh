package ovh

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func init() {
	resource.AddTestSweepers("ovh_vrack_public_routing_priority", &resource.Sweeper{
		Name: "ovh_vrack_public_routing_priority",
		F:    testSweepVrackPublicRoutingPrioirity,
	})
}

var testAccPublicRoutingPriorityConfig = fmt.Sprintf(`
resource "ovh_vrack_public_routing_priority" "vrack_publicRoutingPriority" {
  service_name       = "%s"
  region             = "%s"
  availability_zones = [
		{
			priority: 1,
			name: "%s"
		},
		{
			priority: 2,
			name: "%s"
		},
		{
			priority: 3,
			name: "%s"
		},
	]
}
`, os.Getenv("OVH_VRACK_SERVICE_TEST"),
	os.Getenv("OVH_VRACK_PUBLIC_ROUTING_PRIORITY_REGION_TEST"),
	os.Getenv("OVH_VRACK_PUBLIC_ROUTING_PRIORITY_AZ_A_TEST"),
	os.Getenv("OVH_VRACK_PUBLIC_ROUTING_PRIORITY_AZ_B_TEST"),
	os.Getenv("OVH_VRACK_PUBLIC_ROUTING_PRIORITY_AZ_C_TEST"))

var testAccPublicRoutingPriorityImportConfig = fmt.Sprintf(`
import {
  to = ovh_vrack_public_routing_priority.vrack_publicRoutingPriority
  id = "%s/%s"
}

%s`, os.Getenv("OVH_VRACK_SERVICE_TEST"), os.Getenv("OVH_VRACK_PUBLIC_ROUTING_PRIORITY_ID_TEST"), testAccPublicRoutingPriorityConfig)

func testSweepVrackPublicRoutingPrioirity(region string) error {
	client, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting client: %s", err)
	}
	serviceName := os.Getenv("OVH_VRACK_SERVICE_TEST")
	if serviceName == "" {
		log.Print("[DEBUG] OVH_VRACK_SERVICE_TEST is not set. No ovh_vrack_public_routing_priority to sweep")
		return nil
	}
	publicRoutingPriorityRegion := os.Getenv("OVH_VRACK_PUBLIC_ROUTING_PRIORITY_REGION_TEST")
	if publicRoutingPriorityRegion == "" {
		log.Print("[DEBUG] OVH_VRACK_PUBLIC_ROUTING_PRIORITY_REGION_TEST is not set. No ovh_vrack_public_routing_priority to sweep")
		return nil
	}

	// If the priority is found for the vRack, the delete it
	// list routing priorities for the vRack
	var responseDatas []string
	endpoint := fmt.Sprintf("/vrack/%s/publicRoutingPriority",
		url.PathEscape(serviceName),
	)
	if err := client.Get(endpoint, &responseDatas); err != nil {
		return err
	}

	var priorityID string
	for i := range responseDatas {
		ID := responseDatas[i]
		endpoint := fmt.Sprintf("/vrack/%s/publicRoutingPriority/%s",
			url.PathEscape(serviceName),
			url.PathEscape(ID))

		var resp VrackPublicRoutingPriorityModel
		if err := client.Get(endpoint, &resp); err != nil {
			return err
		}

		if resp.Region.ValueString() == publicRoutingPriorityRegion {
			// publicRoutingPriority is found from vrack/region
			priorityID = resp.PriorityId.ValueString()
			break
		}
	}

	if priorityID == "" {
		// publicRoutingPriority is not on the vrack
		return nil
	}

	deleteEndpoint := fmt.Sprintf("/vrack/%s/publicRoutingPriority/%s",
		url.PathEscape(serviceName),
		url.PathEscape(priorityID),
	)
	task := VrackTask{}
	if err := client.Delete(deleteEndpoint, &task); err != nil {
		return fmt.Errorf("Error calling %s:\n\t %q", deleteEndpoint, err)
	}

	if err := waitForVrackTask(&task, client); err != nil {
		return fmt.Errorf("Error waiting for vrack (%s) to delete publicRoutingPriority (%s): %s", serviceName, priorityID, err)
	}

	return nil
}

// This test should be the simplest possible test for this resource, and will generally set only required fields on the resource.
func TestAccResourceVrackPublicRoutingPriority_basic(t *testing.T) {
	resourceName := "ovh_vrack_public_routing_priority.vrack_publicRoutingPriority"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckPublicRoutingPriority(t)
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				ResourceName: resourceName,
				Config:       testAccPublicRoutingPriorityConfig, // a resource is created and its representation is written to state.
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ovh_vrack_public_routing_priority.vrack_publicRoutingPriority", "service_name", os.Getenv("OVH_VRACK_SERVICE_TEST")),
					resource.TestCheckResourceAttr("ovh_vrack_public_routing_priority.vrack_publicRoutingPriority", "region", os.Getenv("OVH_VRACK_PUBLIC_ROUTING_PRIORITY_REGION_TEST")),
					resource.TestCheckResourceAttr("ovh_vrack_public_routing_priority.vrack_publicRoutingPriority", "availability_zones.0.name", os.Getenv("OVH_VRACK_PUBLIC_ROUTING_PRIORITY_AZ_A_TEST")),
					resource.TestCheckResourceAttr("ovh_vrack_public_routing_priority.vrack_publicRoutingPriority", "availability_zones.1.name", os.Getenv("OVH_VRACK_PUBLIC_ROUTING_PRIORITY_AZ_B_TEST")),
					resource.TestCheckResourceAttr("ovh_vrack_public_routing_priority.vrack_publicRoutingPriority", "availability_zones.2.name", os.Getenv("OVH_VRACK_PUBLIC_ROUTING_PRIORITY_AZ_C_TEST")),
				),
			},
		},
	})
}

func TestAccPublicRoutingPriority_import(t *testing.T) {
	resourceName := "ovh_vrack_public_routing_priority.vrack_publicRoutingPriority"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckPublicRoutingPriority(t)
			checkEnvOrSkip(t, "OVH_VRACK_PUBLIC_ROUTING_PRIORITY_ID_TEST")
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Apply config that imports the existing resource
			{
				Config: testAccPublicRoutingPriorityImportConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "service_name", os.Getenv("OVH_VRACK_SERVICE_TEST")),
					resource.TestCheckResourceAttr(resourceName, "region", os.Getenv("OVH_VRACK_PUBLIC_ROUTING_PRIORITY_REGION_TEST")),
				),
			},
			// Run import again via test framework and verify state matches Read
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				Config:            testAccPublicRoutingPriorityConfig,
				ImportStateIdFunc: func(state *terraform.State) (string, error) {
					rs, ok := state.RootModule().Resources[resourceName]
					if !ok {
						return "", fmt.Errorf("resource %s not found in state", resourceName)
					}
					return rs.Primary.ID, nil
				},
			},
		},
	})
}
