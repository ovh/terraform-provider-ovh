package ovh

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDedicatedNashaData(t *testing.T) {
	serviceName := os.Getenv("OVH_NASHA_SERVICE_TEST")

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			checkEnvOrSkip(t, "OVH_NASHA_SERVICE_TEST")
			testAccPreCheckCredentials(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					data "ovh_dedicated_nasha" "testacc" {
						service_name = "%s"
					}
				`, serviceName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.ovh_dedicated_nasha.testacc", "disk_type", "ssd"),
					resource.TestCheckResourceAttr("data.ovh_dedicated_nasha.testacc", "service_name", serviceName),
					resource.TestCheckResourceAttrSet("data.ovh_dedicated_nasha.testacc", "urn"),
				),
			},
		},
	})
}
