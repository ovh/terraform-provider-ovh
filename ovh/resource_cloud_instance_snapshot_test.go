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

func testAccCloudInstanceSnapshotConfig(serviceName, region, flavorID, imageID, instanceName, snapshotName string) string {
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

resource "ovh_cloud_instance_snapshot" "snapshot" {
  service_name = "%s"
  name         = "%s"
  region       = "%s"
  instance_id  = ovh_cloud_instance.instance.id
}
`, serviceName, region, instanceName, flavorID, imageID, serviceName, snapshotName, region)
}

func TestAccCloudInstanceSnapshot_basic(t *testing.T) {
	serviceName := os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST")
	region := os.Getenv("OVH_CLOUD_PROJECT_REGION_TEST")
	flavorID := resolveInstanceFlavorID(t, serviceName, region, testAccInstanceFlavorName)
	imageID := resolveInstanceImageID(t, serviceName, region, testAccInstanceImageName)
	instanceName := acctest.RandomWithPrefix(test_prefix)
	snapshotName := acctest.RandomWithPrefix(test_prefix)
	snapshotNameUpdated := acctest.RandomWithPrefix(test_prefix)

	config := testAccCloudInstanceSnapshotConfig(serviceName, region, flavorID, imageID, instanceName, snapshotName)
	configUpdated := testAccCloudInstanceSnapshotConfig(serviceName, region, flavorID, imageID, instanceName, snapshotNameUpdated)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheckCloudInstanceV2(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ovh_cloud_instance_snapshot.snapshot", "service_name", serviceName),
					resource.TestCheckResourceAttr("ovh_cloud_instance_snapshot.snapshot", "name", snapshotName),
					resource.TestCheckResourceAttr("ovh_cloud_instance_snapshot.snapshot", "region", region),
					resource.TestCheckResourceAttr("ovh_cloud_instance_snapshot.snapshot", "resource_status", "READY"),
					resource.TestCheckResourceAttrSet("ovh_cloud_instance_snapshot.snapshot", "id"),
					resource.TestCheckResourceAttrSet("ovh_cloud_instance_snapshot.snapshot", "checksum"),
					resource.TestCheckResourceAttrSet("ovh_cloud_instance_snapshot.snapshot", "created_at"),
					resource.TestCheckResourceAttrSet("ovh_cloud_instance_snapshot.snapshot", "instance_id"),
					resource.TestCheckResourceAttrSet("ovh_cloud_instance_snapshot.snapshot", "current_state.name"),
					resource.TestCheckResourceAttrSet("ovh_cloud_instance_snapshot.snapshot", "current_state.instance.id"),
				),
			},
			{
				Config: configUpdated,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("ovh_cloud_instance_snapshot.snapshot", plancheck.ResourceActionReplace),
					},
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ovh_cloud_instance_snapshot.snapshot", "name", snapshotNameUpdated),
					resource.TestCheckResourceAttr("ovh_cloud_instance_snapshot.snapshot", "resource_status", "READY"),
					resource.TestCheckResourceAttrSet("ovh_cloud_instance_snapshot.snapshot", "id"),
					resource.TestCheckResourceAttrSet("ovh_cloud_instance_snapshot.snapshot", "checksum"),
				),
			},
			{
				ResourceName:      "ovh_cloud_instance_snapshot.snapshot",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: func(state *terraform.State) (string, error) {
					return fmt.Sprintf(
						"%s/%s",
						state.RootModule().Resources["ovh_cloud_instance_snapshot.snapshot"].Primary.Attributes["service_name"],
						state.RootModule().Resources["ovh_cloud_instance_snapshot.snapshot"].Primary.Attributes["id"],
					), nil
				},
			},
		},
	})
}
