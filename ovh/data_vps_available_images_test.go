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
						"data.ovh_vps_available_images.all", "service_name", vps),
					resource.TestCheckResourceAttrSet(
						"data.ovh_vps_available_images.all", "image_ids.#"),
					resource.TestCheckResourceAttrSet(
						"data.ovh_vps_available_images.all", "images.#"),
				),
			},
		},
	})
}

func TestAccVPSAvailableImagesDataSource_filter(t *testing.T) {
	vps := os.Getenv("OVH_VPS")
	// Match anything — purpose is to exercise the regex code path.
	config := fmt.Sprintf(testAccVPSAvailableImagesDatasourceConfig_Filter, vps, ".*")

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckVPS(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"data.ovh_vps_available_images.filtered", "name_pattern", ".*"),
					resource.TestCheckResourceAttrSet(
						"data.ovh_vps_available_images.filtered", "image_ids.#"),
				),
			},
		},
	})
}

const testAccVPSAvailableImagesDatasourceConfig_Basic = `
data "ovh_vps_available_images" "all" {
  service_name = "%s"
}
`

const testAccVPSAvailableImagesDatasourceConfig_Filter = `
data "ovh_vps_available_images" "filtered" {
  service_name = "%s"
  name_pattern = "%s"
}
`
