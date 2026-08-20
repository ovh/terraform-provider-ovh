package ovh

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDataSourceCloudInstanceBackup_basic(t *testing.T) {
	serviceName := os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST")
	region := os.Getenv("OVH_CLOUD_PROJECT_REGION_TEST")
	flavorID := resolveInstanceFlavorID(t, serviceName, region, testAccInstanceFlavorName)
	imageID := resolveInstanceImageID(t, serviceName, region, testAccInstanceImageName)
	instanceName := acctest.RandomWithPrefix(test_prefix)
	backupName := acctest.RandomWithPrefix(test_prefix)

	config := testAccCloudInstanceBackupConfig(serviceName, region, flavorID, imageID, instanceName, backupName) + fmt.Sprintf(`
data "ovh_cloud_instance_backup" "backup" {
  service_name = "%s"
  id           = ovh_cloud_instance_backup.backup.id
}
`, serviceName)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheckCloudInstanceV2(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.ovh_cloud_instance_backup.backup", "service_name", serviceName),
					resource.TestCheckResourceAttrPair(
						"data.ovh_cloud_instance_backup.backup", "id",
						"ovh_cloud_instance_backup.backup", "id",
					),
					resource.TestCheckResourceAttrPair(
						"data.ovh_cloud_instance_backup.backup", "instance_id",
						"ovh_cloud_instance.instance", "id",
					),
					resource.TestCheckResourceAttr("data.ovh_cloud_instance_backup.backup", "name", backupName),
					resource.TestCheckResourceAttr("data.ovh_cloud_instance_backup.backup", "location.region", region),
					resource.TestCheckResourceAttr("data.ovh_cloud_instance_backup.backup", "resource_status", "READY"),
					resource.TestCheckResourceAttrSet("data.ovh_cloud_instance_backup.backup", "status"),
					resource.TestCheckResourceAttrSet("data.ovh_cloud_instance_backup.backup", "visibility"),
					resource.TestCheckResourceAttrPair(
						"data.ovh_cloud_instance_backup.backup", "instance_id",
						"ovh_cloud_instance_backup.backup", "current_state.instance.id",
					),
				),
			},
		},
	})
}
