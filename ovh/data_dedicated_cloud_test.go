package ovh

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDedicatedCloudData(t *testing.T) {
	serviceName := os.Getenv("OVH_DEDICATED_CLOUD")

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			checkEnvOrSkip(t, "OVH_DEDICATED_CLOUD")
			testAccPreCheckCredentials(t)
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					data "ovh_dedicated_cloud" "pcc" {
						service_name = "%s"
					}
				`, serviceName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.ovh_dedicated_cloud.pcc", "service_name", serviceName),
					resource.TestCheckResourceAttrSet("data.ovh_dedicated_cloud.pcc", "state"),
					resource.TestCheckResourceAttrSet("data.ovh_dedicated_cloud.pcc", "v_scope_url"),
					resource.TestCheckResourceAttrSet("data.ovh_dedicated_cloud.pcc", "web_interface_url"),
					resource.TestCheckResourceAttrSet("data.ovh_dedicated_cloud.pcc", "certified_interface_url"),
					resource.TestCheckResourceAttrSet("data.ovh_dedicated_cloud.pcc", "commercial_range"),
					resource.TestCheckResourceAttrSet("data.ovh_dedicated_cloud.pcc", "iam.urn"),
				),
			},
		},
	})
}
