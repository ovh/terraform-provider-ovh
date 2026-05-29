package ovh

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

const testAccDataSourceVPSSnapshotDownloadBasic = `
data "ovh_vps_snapshot_download" "test" {
  service_name = "%s"
}
`

func TestAccDataSourceVPSSnapshotDownload_basic(t *testing.T) {
	serviceName := os.Getenv("OVH_VPS_SNAPSHOT_SERVICE_NAME")
	config := fmt.Sprintf(testAccDataSourceVPSSnapshotDownloadBasic, serviceName)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckVPS(t)
			checkEnvOrSkip(t, "OVH_VPS_SNAPSHOT_SERVICE_NAME")
			checkVPSServiceOptionSubscribed(t, os.Getenv("OVH_VPS_SNAPSHOT_SERVICE_NAME"), "snapshot")
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"data.ovh_vps_snapshot_download.test", "service_name", serviceName),
					resource.TestCheckResourceAttrSet(
						"data.ovh_vps_snapshot_download.test", "url"),
					resource.TestCheckResourceAttrSet(
						"data.ovh_vps_snapshot_download.test", "size"),
				),
			},
		},
	})
}
