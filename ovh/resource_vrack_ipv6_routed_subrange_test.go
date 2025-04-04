package ovh

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"

	"github.com/ovh/go-ovh/ovh"
)

func init() {
	resource.AddTestSweepers("ovh_vrack_ipv6_routed_subrange", &resource.Sweeper{
		Name: "ovh_vrack_ipv6_routed_subrange",
		F:    testSweepVrackIPv6RoutedSubrange,
	})
}

func testAccVrackIPv6RoutedSubrangeConfig(resourceName, serviceName, block, routedSubrange, nexthop string) string {
	rn := strings.Split(resourceName, ".")
	return fmt.Sprintf(`
	resource %s %s {
		service_name 	= "%s"
		block        	= "%s"
		routed_subrange = "%s"
		nexthop 		= "%s"
	}`, rn[0], rn[1], serviceName, block, routedSubrange, nexthop)
}

func testSweepVrackIPv6RoutedSubrange(region string) error {
	client, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting client: %s", err)
	}

	vrackId := os.Getenv("OVH_VRACK_SERVICE_TEST")
	if vrackId == "" {
		log.Print("[DEBUG] OVH_VRACK_SERVICE_TEST is not set. No vrack_ipv6_routed_subrange to sweep")
		return nil
	}

	ipBlock := os.Getenv("OVH_IP_V6_BLOCK_TEST")
	if ipBlock == "" {
		log.Print("[DEBUG] OVH_IP_V6_BLOCK_TEST is not set. No vrack_ipv6_routed_subrange to sweep")
		return nil
	}

	routedSubrange := os.Getenv("OVH_IP_V6_ROUTED_SUBRANGE_TEST")
	if ipBlock == "" {
		log.Print("[DEBUG] OVH_IP_V6_ROUTED_SUBRANGE_TEST is not set. No vrack_ipv6_routed_subrange to sweep")
		return nil
	}

	endpoint := fmt.Sprintf("/vrack/%s/ipv6/%s/routedSubrange/%s",
		url.PathEscape(vrackId),
		url.PathEscape(ipBlock),
		url.PathEscape(routedSubrange),
	)

	if err := client.Get(endpoint, nil); err != nil {
		if errOvh, ok := err.(*ovh.APIError); ok && errOvh.Code == 404 {
			return nil
		}
		return err
	}

	task := VrackTask{}
	if err := client.Delete(endpoint, &task); err != nil {
		return fmt.Errorf("Error calling %s:\n\t %q", endpoint, err)
	}

	if err := waitForVrackTask(&task, client); err != nil {
		return fmt.Errorf("Error waiting for vrack (%s) to delete ipv6 routed subrange (%s): %s", vrackId, routedSubrange, err)
	}

	return nil
}

func TestAccVrackIPv6RoutedSubrange_basic(t *testing.T) {
	resourceName := "ovh_vrack_ipv6_routed_subrange.test-vrack-ipv6-routed-subrange-basic"
	serviceName := os.Getenv("OVH_VRACK_SERVICE_TEST")
	ipBlock := os.Getenv("OVH_IP_V6_BLOCK_TEST")
	routedSubrange := os.Getenv("OVH_IP_V6_ROUTED_SUBRANGE_TEST")
	nexhop := os.Getenv("OVH_IP_V6_ROUTED_SUBRANGE_NEXTHOP_TEST")

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckIPv6RoutedSubrangeVRack(t)
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccVrackIPv6RoutedSubrangeConfig(resourceName, serviceName, ipBlock, routedSubrange, nexhop),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "service_name", serviceName),
					resource.TestCheckResourceAttr(resourceName, "block", ipBlock),
					resource.TestCheckResourceAttr(resourceName, "routed_subrange", routedSubrange),
					resource.TestCheckResourceAttr(resourceName, "nexthop", nexhop),
				),
			},
		},
	})
}

func TestAccVrackIPv6RoutedSubrange_import(t *testing.T) {
	resourceName := "ovh_vrack_ipv6_routed_subrange.test-vrack-ipv6-routed-subrange-import"
	serviceName := os.Getenv("OVH_VRACK_SERVICE_TEST")
	ipBlock := os.Getenv("OVH_IP_V6_BLOCK_IMPORT_TEST")
	routedSubrange := os.Getenv("OVH_IP_V6_ROUTED_SUBRANGE_IMPORT_TEST")
	nexhop := os.Getenv("OVH_IP_V6_ROUTED_SUBRANGE_NEXTHOP_TEST")

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckIPv6RoutedSubrangeImportVRack(t)
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				Config:            testAccVrackIPv6RoutedSubrangeConfig(resourceName, serviceName, ipBlock, routedSubrange, nexhop),
				ImportStateIdFunc: func(state *terraform.State) (string, error) {
					return fmt.Sprintf("%s,%s,%s", serviceName, ipBlock, routedSubrange), nil
				},
			},
		},
	})
}
