package ovh

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

var testAccVrackCloudProjectConfig = fmt.Sprintf(`
resource "ovh_vrack_cloudproject" "vcp" {
  vrack_id = "%s"
  project_id = "%s"
}
`, os.Getenv("OVH_VRACK"), os.Getenv("OVH_PUBLIC_CLOUD"))

var testAccVrackCloudProjectConfig_legacy = fmt.Sprintf(`
resource "ovh_vrack_publiccloud_attachment" "attach" {
  vrack_id = "%s"
  project_id = "%s"
}
`, os.Getenv("OVH_VRACK"), os.Getenv("OVH_PUBLIC_CLOUD"))

func TestAccVrackCloudProject_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccCheckVrackCloudProjectPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccVrackCloudProjectConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ovh_vrack_cloudproject.vcp", "vrack_id", os.Getenv("OVH_VRACK")),
					resource.TestCheckResourceAttr("ovh_vrack_cloudproject.vcp", "project_id", os.Getenv("OVH_PUBLIC_CLOUD")),
				),
			},
		},
	})
}

func TestAccVrackCloudProject_legacy(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccCheckVrackCloudProjectPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccVrackCloudProjectConfig_legacy,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ovh_vrack_publiccloud_attachment.attach", "vrack_id", os.Getenv("OVH_VRACK")),
					resource.TestCheckResourceAttr("ovh_vrack_publiccloud_attachment.attach", "project_id", os.Getenv("OVH_PUBLIC_CLOUD")),
				),
			},
		},
	})
}

func testAccCheckVrackCloudProjectPreCheck(t *testing.T) {
	testAccPreCheckVRack(t)
	testAccCheckVRackExists(t)
	testAccCheckPublicCloudExists(t)
}
