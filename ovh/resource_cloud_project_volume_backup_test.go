package ovh

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccCloudProjectVolumeBackup_basic(t *testing.T) {
	serviceName := os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST")
	regionName := os.Getenv("OVH_CLOUD_PROJECT_REGION_TEST")

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		PreCheck:                 func() { testAccPreCheckCloud(t); testAccCheckCloudProjectExists(t) },
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					resource "ovh_cloud_project_volume" "volume"  {
						region_name = "%s"
						service_name = "%s"
						description = "test"
						name = "test"
						size = 15
						type = "classic"
					}

					resource "ovh_cloud_project_volume_backup" "backup"  {
						region_name = "%s"
						service_name = "%s"
						name = "test_backup"
						volume_id = ovh_cloud_project_volume.volume.id
					}
				`,
					regionName,
					serviceName,
					regionName,
					serviceName,
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("ovh_cloud_project_volume_backup.backup", "region_name", regionName),
					resource.TestCheckResourceAttr("ovh_cloud_project_volume_backup.backup", "service_name", serviceName),
					resource.TestCheckResourceAttrSet("ovh_cloud_project_volume_backup.backup", "size"),
					resource.TestCheckResourceAttr("ovh_cloud_project_volume_backup.backup", "status", "ok"),
				),
			},
			{
				ImportState:         true,
				ImportStateVerify:   true,
				ResourceName:        "ovh_cloud_project_volume_backup.backup",
				ImportStateIdPrefix: serviceName + "/" + regionName + "/",
			},
		},
	})
}
