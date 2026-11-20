package ovh

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccVPSAvailableImageDataSource_basic(t *testing.T) {
	vps := os.Getenv("OVH_VPS")
	config := fmt.Sprintf(testAccVPSAvailableImageDatasourceConfig_Basic, vps, vps)
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckVPS(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"data.ovh_vps_available_images.available_images", "service_name", vps),
					resource.TestCheckResourceAttrSet(
						"data.ovh_vps_available_images.available_images", "image_ids.#"),
				),
			},
		},
	})
}

const testAccVPSAvailableImageDatasourceConfig_Basic = `
data "ovh_vps_available_images" "available_images" {
  service_name = "%s"
}

data "ovh_vps_available_image" "all_images" {
  for_each     = toset(data.ovh_vps_available_images.available_images.image_ids)
  service_name = "%s"
  image_id     = each.value
}
`
