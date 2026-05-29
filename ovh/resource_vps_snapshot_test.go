package ovh

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// testAccPreCheckVPSSnapshot ensures the env vars required for VPS-snapshot
// acceptance tests are present. Tests that hit this helper are skipped unless
// the operator opted in via OVH_TESTACC_ORDER_VPS and supplied an existing
// VPS (OVH_VPS_SNAPSHOT_SERVICE_NAME) with the "snapshot" option subscribed.
func testAccPreCheckVPSSnapshot(t *testing.T) {
	t.Helper()
	testAccPreCheckCredentials(t)
	checkEnvOrSkip(t, "OVH_TESTACC_ORDER_VPS")
	checkEnvOrSkip(t, "OVH_VPS_SNAPSHOT_SERVICE_NAME")
	checkVPSServiceOptionSubscribed(t, os.Getenv("OVH_VPS_SNAPSHOT_SERVICE_NAME"), "snapshot")
}

const testAccResourceVPSSnapshotBasic = `
resource "ovh_vps_snapshot" "snap" {
  service_name = "%s"
  description  = "%s"
}
`

// TestAccResourceVPSSnapshot_basic is an acceptance-test skeleton. It is not
// executed unless OVH_TESTACC_ORDER_VPS and OVH_VPS_SNAPSHOT_SERVICE_NAME
// are set. The targeted VPS must already have the "snapshot" option
// subscribed.
func TestAccResourceVPSSnapshot_basic(t *testing.T) {
	serviceName := os.Getenv("OVH_VPS_SNAPSHOT_SERVICE_NAME")
	description := "tf-acc-test"
	config := fmt.Sprintf(testAccResourceVPSSnapshotBasic, serviceName, description)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheckVPSSnapshot(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ovh_vps_snapshot.snap", "service_name", serviceName),
					resource.TestCheckResourceAttr("ovh_vps_snapshot.snap", "description", description),
					resource.TestCheckResourceAttrSet("ovh_vps_snapshot.snap", "creation_date"),
					resource.TestCheckResourceAttrSet("ovh_vps_snapshot.snap", "region"),
				),
			},
		},
	})
}
