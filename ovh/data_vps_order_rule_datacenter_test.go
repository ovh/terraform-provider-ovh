package ovh

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccVPSOrderRuleDatacenterDataSource_basic(t *testing.T) {
	// Allow the test runner to pick a subsidiary/planCode pair that matches
	// their account region. Defaults to a US/2025-model1 pair if unset; users
	// running against EU/CA accounts should set OVH_TESTACC_VPS_SUBSIDIARY
	// (e.g. FR, CA, DE) and OVH_TESTACC_VPS_PLAN_CODE accordingly.
	subsidiary := os.Getenv("OVH_TESTACC_VPS_SUBSIDIARY")
	if subsidiary == "" {
		subsidiary = "US"
	}
	planCode := os.Getenv("OVH_TESTACC_VPS_PLAN_CODE")
	if planCode == "" {
		planCode = "vps-2025-model1"
	}

	config := fmt.Sprintf(`
data "ovh_vps_order_rule_datacenter" "test" {
  ovh_subsidiary = "%s"
  plan_code      = "%s"
}
`, subsidiary, planCode)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckCredentials(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(
						"data.ovh_vps_order_rule_datacenter.test", "datacenters.#"),
				),
			},
		},
	})
}
