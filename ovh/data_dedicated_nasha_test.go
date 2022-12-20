package ovh

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDedicatedNashaData(t *testing.T) {
	serviceName := os.Getenv("OVH_NASHA_SERVICE_TEST")
	nashaIp := os.Getenv("OVH_NASHA_SERVICE_IP_TEST")

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			checkEnvOrFail(t, "OVH_NASHA_SERVICE_TEST")
			checkEnvOrFail(t, "OVH_NASHA_SERVICE_IP_TEST")
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
					resource.TestCheckResourceAttr("data.ovh_dedicated_nasha.testacc", "zpool_size", "3000"),
					resource.TestCheckResourceAttr("data.ovh_dedicated_nasha.testacc", "ip", nashaIp),
				),
			},
		},
	})
}
