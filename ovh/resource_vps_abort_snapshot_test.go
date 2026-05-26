package ovh

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

const testAccVpsAbortSnapshotConfig = `
resource "ovh_vps_abort_snapshot" "abort" {
  service_name = "%s"

  triggers = {
    nonce = "1"
  }
}
`

func TestAccVpsAbortSnapshot_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckVPS(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccVpsAbortSnapshotConfig, os.Getenv("OVH_VPS")),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(
						"ovh_vps_abort_snapshot.abort", "aborted_at"),
					resource.TestCheckResourceAttrSet(
						"ovh_vps_abort_snapshot.abort", "id"),
				),
			},
		},
	})
}
