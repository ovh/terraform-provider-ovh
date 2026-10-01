package ovh

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccCloudRegionsDataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheckCloud(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					data "ovh_cloud_regions" "regions" {
						service_name = "%s"
					}
				`, os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST")),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.ovh_cloud_regions.regions", "regions.#"),
					resource.TestCheckResourceAttrSet("data.ovh_cloud_regions.regions", "regions.0.name"),
					resource.TestCheckResourceAttrSet("data.ovh_cloud_regions.regions", "regions.0.status"),
				),
			},
		},
	})
}
