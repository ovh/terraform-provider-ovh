package ovh

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccVPSAvailableImagesDataSource_basic(t *testing.T) {
	vps := os.Getenv("OVH_VPS")
	config := fmt.Sprintf(testAccVPSAvailableImagesDatasourceConfig_Basic, vps)

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
						"data.ovh_vps_available_images.available_images", "image_ids"),
				),
			},
		},
	})
}

const testAccVPSAvailableImagesDatasourceConfig_Basic = `
data "ovh_vps_available_images" "available_images" {
  service_name  = "%s"
}
`
