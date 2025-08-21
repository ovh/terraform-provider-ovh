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
	resource.AddTestSweepers("ovh_vrack_dedicated_server_interface", &resource.Sweeper{
		Name: "ovh_vrack_dedicated_server_interface",
		F:    testSweepVrackDedicatedServerInterface,
	})
}

func testSweepVrackDedicatedServerInterface(region string) error {
	config, err := sharedConfigForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting client: %s", err)
	}

	baremetal := os.Getenv("OVH_DEDICATED_SERVER")
	vrack := os.Getenv("OVH_VRACK_SERVICE_TEST")

	// First get vni IDs
	vniIDs := []string{}
	if err := config.OVHClient.Get(
		fmt.Sprintf(
			"/dedicated/server/%s/virtualNetworkInterface",
			url.PathEscape(baremetal),
		),
		&vniIDs,
	); err != nil {
		return fmt.Errorf("error retrieving VNIs for dedicated server %s: %s", baremetal, err)
	}

	if len(vniIDs) == 0 {
		log.Printf("[INFO] No VNI IDs found for dedicated server %s, nothing to sweep", baremetal)
		return nil
	}

	// Fetch vrack interfaces
	vnis := []*DedicatedServerVNI{}
	for _, id := range vniIDs {
		var vni DedicatedServerVNI
		err := config.OVHClient.Get(
			fmt.Sprintf("/dedicated/server/%s/virtualNetworkInterface/%s", url.PathEscape(baremetal), url.PathEscape(id)),
			&vni,
		)

		if err != nil {
			return fmt.Errorf("error retrieving VNI info for dedicated server %s: %s", baremetal, err)
		}

		if vni.Enabled && vni.Mode == "vrack" {
			vnis = append(vnis, &vni)
		}
	}

	if len(vnis) == 0 {
		log.Printf("[INFO] No vrack interfaces found for dedicated server %s, nothing to sweep", baremetal)
		return nil
	}

	// Check if vrack is attached to interface
	endpoint := fmt.Sprintf("/vrack/%s/dedicatedServerInterface/%s",
		url.PathEscape(vrack),
		url.PathEscape(vnis[0].Uuid),
	)

	var vds VrackDedicatedServerInterface
	if err := config.OVHClient.Get(endpoint, &vds); err != nil {
		if errOvh, ok := err.(*ovh.APIError); ok && errOvh.Code == 404 {
			// No vrack interface found, nothing to delete
			return nil
		}
		return fmt.Errorf("error retrieving vrack dedicated server interface %s: %s", vds.DedicatedServerInterface, err)
	}

	// Remove vrack from interface
	var task VrackTask
	if err := config.OVHClient.Delete(endpoint, &task); err != nil {
		return fmt.Errorf("error calling DELETE %s with %s/%s:\n\t %q", endpoint, vrack, vnis[0].Uuid, err)
	}

	if err := waitForVrackTask(&task, config.OVHClient); err != nil {
		return fmt.Errorf("error waiting for vrack (%s) to detach dedicated server (%s): %s", vrack, vnis[0].Uuid, err)
	}

	return nil
}

var testAccVrackDedicatedServerInterfaceConfig = fmt.Sprintf(`
data "ovh_dedicated_server" "server" {
  service_name = "%s"
}

resource "ovh_vrack_dedicated_server_interface" "vdsi" {
  service_name = "%s"
  interface_id = data.ovh_dedicated_server.server.enabled_vrack_vnis[0]
}
`, os.Getenv("OVH_DEDICATED_SERVER"), os.Getenv("OVH_VRACK_SERVICE_TEST"))

func TestAccVrackDedicatedServerInterface_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccCheckVrackDedicatedServerInterfacePreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccVrackDedicatedServerInterfaceConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ovh_vrack_dedicated_server_interface.vdsi", "service_name", os.Getenv("OVH_VRACK_SERVICE_TEST")),
					resource.TestCheckResourceAttrSet("ovh_vrack_dedicated_server_interface.vdsi", "interface_id"),
				),
			},
		},
	})
}

func testAccCheckVrackDedicatedServerInterfacePreCheck(t *testing.T) {
	testAccPreCheckVRack(t)
	testAccCheckVRackExists(t)
	testAccPreCheckDedicatedServer(t)
}
