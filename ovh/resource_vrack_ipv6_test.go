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
)

func init() {
	resource.AddTestSweepers("ovh_vrack_ipv6", &resource.Sweeper{
		Name: "ovh_vrack_ipv6",
		F:    testSweepVrackIPv6,
	})
}

func testAccVrackIPv6Config(resourceName, serviceName, block, slaac string) string {
	return fmt.Sprintf(`
	resource "ovh_vrack_ipv6" "%s" {
		service_name = "%s"
		block        = "%s"
		bridged_subrange {
			slaac = "%s"
		}
	}`, resourceName, serviceName, block, slaac)
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
		log.Print("[DEBUG] OVH_IP_V6_BLOCK_TEST is not set. No vrack_ipv6 to sweep")
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
			testAccPreCheckIPv6VRack(t)
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccVrackIPv6Config("test-vrack-ipv6-basic", serviceName, ipBlock, "enabled"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ovh_vrack_ipv6.test-vrack-ipv6-basic", "service_name", serviceName),
					resource.TestCheckResourceAttr("ovh_vrack_ipv6.test-vrack-ipv6-basic", "block", ipBlock),
					resource.TestCheckResourceAttrSet("ovh_vrack_ipv6.test-vrack-ipv6-basic", "bridged_subrange.0.subrange"),
					resource.TestCheckResourceAttrSet("ovh_vrack_ipv6.test-vrack-ipv6-basic", "bridged_subrange.0.gateway"),
					resource.TestCheckResourceAttr("ovh_vrack_ipv6.test-vrack-ipv6-basic", "bridged_subrange.0.slaac", "enabled"),
				),
			},
			{
				// update the Slaac status of the bridged subrange.
				Config: testAccVrackIPv6Config("test-vrack-ipv6-basic", serviceName, ipBlock, "disabled"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ovh_vrack_ipv6.test-vrack-ipv6-basic", "service_name", serviceName),
					resource.TestCheckResourceAttr("ovh_vrack_ipv6.test-vrack-ipv6-basic", "block", ipBlock),
					resource.TestCheckResourceAttrSet("ovh_vrack_ipv6.test-vrack-ipv6-basic", "bridged_subrange.0.subrange"),
					resource.TestCheckResourceAttrSet("ovh_vrack_ipv6.test-vrack-ipv6-basic", "bridged_subrange.0.gateway"),
					resource.TestCheckResourceAttr("ovh_vrack_ipv6.test-vrack-ipv6-basic", "bridged_subrange.0.slaac", "disabled"),
				),
			},
		},
	})
}

func TestAccVrackIPv6_import(t *testing.T) {
	serviceName := os.Getenv("OVH_VRACK_SERVICE_TEST")
	ipBlock := os.Getenv("OVH_IP_V6_BLOCK_IMPORT_TEST")

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckIPv6ImportVRack(t)
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				ResourceName:      "ovh_vrack_ipv6.test-vrack-ipv6-import",
				ImportState:       true,
				ImportStateVerify: true,
				Config:            testAccVrackIPv6Config("test-vrack-ipv6-import", serviceName, ipBlock, "enabled"),
				ImportStateIdFunc: func(state *terraform.State) (string, error) {
					return fmt.Sprintf("%s,%s", serviceName, ipBlock), nil
				},
			},
		},
	})
}
