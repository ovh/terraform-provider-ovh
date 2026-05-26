package ovh

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccVPSAvailableImageDataSource_basic(t *testing.T) {
	vps := os.Getenv("OVH_VPS")
	imageID := os.Getenv("OVH_VPS_IMAGE_ID")
	config := fmt.Sprintf(testAccVPSAvailableImageDatasourceConfig_Basic, vps, imageID)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckVPS(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"data.ovh_vps_available_image.img", "service_name", vps),
					resource.TestCheckResourceAttr(
						"data.ovh_vps_available_image.img", "image_id", imageID),
					resource.TestCheckResourceAttrSet(
						"data.ovh_vps_available_image.img", "name"),
				),
			},
		},
	})
}

const testAccVPSAvailableImageDatasourceConfig_Basic = `
data "ovh_vps_available_image" "img" {
  service_name = "%s"
  image_id     = "%s"
}
`
