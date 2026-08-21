package ovh

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func testAccCloudInstanceBackupConfig(serviceName, region, flavorID, imageID, instanceName, backupName string) string {
	return fmt.Sprintf(`
resource "ovh_cloud_instance" "instance" {
  service_name = "%s"
  region       = "%s"
  name         = "%s"
  flavor_id    = "%s"
  image_id     = "%s"

  networks = [
    { auto_assign_public_ip = true },
  ]
}

resource "ovh_cloud_instance_backup" "backup" {
  service_name = "%s"
  name         = "%s"
  region       = "%s"
  instance_id  = ovh_cloud_instance.instance.id
}
`, serviceName, region, instanceName, flavorID, imageID, serviceName, backupName, region)
}

func TestAccCloudInstanceBackup_basic(t *testing.T) {
	serviceName := os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST")
	region := os.Getenv("OVH_CLOUD_PROJECT_REGION_TEST")
	flavorID := resolveInstanceFlavorID(t, serviceName, region, testAccInstanceFlavorName)
	imageID := resolveInstanceImageID(t, serviceName, region, testAccInstanceImageName)
	instanceName := acctest.RandomWithPrefix(test_prefix)
	backupName := acctest.RandomWithPrefix(test_prefix)
	backupNameUpdated := acctest.RandomWithPrefix(test_prefix)

	config := testAccCloudInstanceBackupConfig(serviceName, region, flavorID, imageID, instanceName, backupName)
	configUpdated := testAccCloudInstanceBackupConfig(serviceName, region, flavorID, imageID, instanceName, backupNameUpdated)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheckCloudInstanceV2(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ovh_cloud_instance_backup.backup", "service_name", serviceName),
					resource.TestCheckResourceAttr("ovh_cloud_instance_backup.backup", "name", backupName),
					resource.TestCheckResourceAttr("ovh_cloud_instance_backup.backup", "region", region),
					resource.TestCheckResourceAttr("ovh_cloud_instance_backup.backup", "resource_status", "READY"),
					resource.TestCheckResourceAttrSet("ovh_cloud_instance_backup.backup", "id"),
					resource.TestCheckResourceAttrSet("ovh_cloud_instance_backup.backup", "checksum"),
					resource.TestCheckResourceAttrSet("ovh_cloud_instance_backup.backup", "created_at"),
					resource.TestCheckResourceAttrSet("ovh_cloud_instance_backup.backup", "instance_id"),
					resource.TestCheckResourceAttrSet("ovh_cloud_instance_backup.backup", "current_state.name"),
					resource.TestCheckResourceAttrSet("ovh_cloud_instance_backup.backup", "current_state.instance.id"),
				),
			},
			{
				Config: configUpdated,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("ovh_cloud_instance_backup.backup", plancheck.ResourceActionReplace),
					},
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ovh_cloud_instance_backup.backup", "name", backupNameUpdated),
					resource.TestCheckResourceAttr("ovh_cloud_instance_backup.backup", "resource_status", "READY"),
					resource.TestCheckResourceAttrSet("ovh_cloud_instance_backup.backup", "id"),
					resource.TestCheckResourceAttrSet("ovh_cloud_instance_backup.backup", "checksum"),
				),
			},
			{
				ResourceName:      "ovh_cloud_instance_backup.backup",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: func(state *terraform.State) (string, error) {
					return fmt.Sprintf(
						"%s/%s",
						state.RootModule().Resources["ovh_cloud_instance_backup.backup"].Primary.Attributes["service_name"],
						state.RootModule().Resources["ovh_cloud_instance_backup.backup"].Primary.Attributes["id"],
					), nil
				},
			},
		},
	})
}
