package ovh

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccCloudRegionDataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheckCloudRegion(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					data "ovh_cloud_region" "region" {
						service_name = "%s"
						name         = "%s"
					}
				`,
					os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST"),
					os.Getenv("OVH_CLOUD_PROJECT_REGION_TEST"),
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(
						"data.ovh_cloud_region.region", "name", os.Getenv("OVH_CLOUD_PROJECT_REGION_TEST")),
					resource.TestCheckResourceAttrSet("data.ovh_cloud_region.region", "status"),
					resource.TestCheckResourceAttrSet("data.ovh_cloud_region.region", "continent"),
					resource.TestCheckResourceAttrSet("data.ovh_cloud_region.region", "services.#"),
				),
			},
		},
	})
}
