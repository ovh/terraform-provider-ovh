package ovh

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

const testAccDataSourceVPSSnapshotConfig = `
data "ovh_vps_snapshot" "snap" {
  service_name = "%s"
}
`

// TestAccDataSourceVPSSnapshot_basic is gated by OVH_TESTACC_ORDER_VPS and
// OVH_VPS_SNAPSHOT_SERVICE_NAME (the existing VPS with an existing snapshot).
func TestAccDataSourceVPSSnapshot_basic(t *testing.T) {
	serviceName := os.Getenv("OVH_VPS_SNAPSHOT_SERVICE_NAME")
	config := fmt.Sprintf(testAccDataSourceVPSSnapshotConfig, serviceName)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheckVPSSnapshot(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.ovh_vps_snapshot.snap", "service_name", serviceName),
					resource.TestCheckResourceAttrSet("data.ovh_vps_snapshot.snap", "creation_date"),
					resource.TestCheckResourceAttrSet("data.ovh_vps_snapshot.snap", "region"),
				),
			},
		},
	})
}
