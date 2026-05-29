package ovh

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

const testAccVPSIpCountryAvailableConfig = `
data "ovh_vps_ip_country_available" "countries" {
  service_name = "%s"
}
`

func TestAccVPSIpCountryAvailableDataSource_basic(t *testing.T) {
	vps := os.Getenv("OVH_VPS")
	if vps == "" {
		t.Skip("OVH_VPS must be set for this acceptance test")
	}
	config := fmt.Sprintf(testAccVPSIpCountryAvailableConfig, vps)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckCredentials(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(
						"data.ovh_vps_ip_country_available.countries", "result.#"),
				),
			},
		},
	})
}
