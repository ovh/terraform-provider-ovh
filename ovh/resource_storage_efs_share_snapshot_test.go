package ovh

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccStorageEfsShareSnapshot_basic(t *testing.T) {
	serviceName := os.Getenv("OVH_STORAGE_EFS_SERVICE_TEST")

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		ExternalProviders: map[string]resource.ExternalProvider{
			"time": {
				VersionConstraint: "0.13.1",
				Source:            "registry.terraform.io/hashicorp/time",
			},
		},
		PreCheck: func() { testAccPreCheckStorageEfs(t); testAccCheckStorageEfsExists(t) },
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
resource "ovh_storage_efs_share" "share" {
  service_name = "%s"
  name         = "share"
  description  = "My share"
  protocol     = "NFS"
  size         = 100
}

resource "time_sleep" "wait_10_seconds" {
  depends_on = [ovh_storage_efs_share.share]

  destroy_duration = "10s"
}

resource "ovh_storage_efs_share_snapshot" "snapshot" {
  depends_on = [time_sleep.wait_10_seconds]

  service_name = "%s"
  share_id = ovh_storage_efs_share.share.id
  name = "snapshot"
  description = "My snapshot"
}`,
					serviceName,
					serviceName,
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("ovh_storage_efs_share_snapshot.snapshot", "service_name", serviceName),
					resource.TestCheckResourceAttrSet("ovh_storage_efs_share_snapshot.snapshot", "id"),
					resource.TestCheckResourceAttrSet("ovh_storage_efs_share_snapshot.snapshot", "share_id"),
					resource.TestCheckResourceAttrSet("ovh_storage_efs_share_snapshot.snapshot", "created_at"),
					resource.TestCheckResourceAttrSet("ovh_storage_efs_share_snapshot.snapshot", "path"),
					resource.TestCheckResourceAttr("ovh_storage_efs_share_snapshot.snapshot", "name", "snapshot"),
					resource.TestCheckResourceAttr("ovh_storage_efs_share_snapshot.snapshot", "description", "My snapshot"),
					resource.TestCheckResourceAttr("ovh_storage_efs_share_snapshot.snapshot", "status", "available"),
					resource.TestCheckResourceAttr("ovh_storage_efs_share_snapshot.snapshot", "type", "manual"),
				),
			},
			{
				Config: fmt.Sprintf(`
resource "ovh_storage_efs_share" "share" {
  service_name = "%s"
  name         = "share"
  description  = "My share"
  protocol     = "NFS"
  size         = 100
}

resource "time_sleep" "wait_10_seconds" {
  depends_on = [ovh_storage_efs_share.share]

  destroy_duration = "10s"
}

resource "ovh_storage_efs_share_snapshot" "snapshot" {
  depends_on = [time_sleep.wait_10_seconds]

  service_name = "%s"
  share_id = ovh_storage_efs_share.share.id
  name = "snapshot_updated"
  description = "My snapshot updated"
}
`,
					serviceName,
					serviceName,
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("ovh_storage_efs_share_snapshot.snapshot", "service_name", serviceName),
					resource.TestCheckResourceAttrSet("ovh_storage_efs_share_snapshot.snapshot", "id"),
					resource.TestCheckResourceAttrSet("ovh_storage_efs_share_snapshot.snapshot", "share_id"),
					resource.TestCheckResourceAttrSet("ovh_storage_efs_share_snapshot.snapshot", "created_at"),
					resource.TestCheckResourceAttrSet("ovh_storage_efs_share_snapshot.snapshot", "path"),
					resource.TestCheckResourceAttr("ovh_storage_efs_share_snapshot.snapshot", "name", "snapshot_updated"),
					resource.TestCheckResourceAttr("ovh_storage_efs_share_snapshot.snapshot", "description", "My snapshot updated"),
					resource.TestCheckResourceAttr("ovh_storage_efs_share_snapshot.snapshot", "status", "available"),
					resource.TestCheckResourceAttr("ovh_storage_efs_share_snapshot.snapshot", "type", "manual"),
				),
			},
		},
	})

}
