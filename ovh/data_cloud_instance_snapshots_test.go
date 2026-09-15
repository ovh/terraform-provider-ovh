package ovh

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDataSourceCloudInstanceSnapshots_basic(t *testing.T) {
	serviceName := os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST")
	region := os.Getenv("OVH_CLOUD_PROJECT_REGION_TEST")
	flavorID := resolveInstanceFlavorID(t, serviceName, region, testAccInstanceFlavorName)
	imageID := resolveInstanceImageID(t, serviceName, region, testAccInstanceImageName)
	instanceName := acctest.RandomWithPrefix(test_prefix)
	snapshotName := acctest.RandomWithPrefix(test_prefix)

	config := testAccCloudInstanceSnapshotConfig(serviceName, region, flavorID, imageID, instanceName, snapshotName) + fmt.Sprintf(`
data "ovh_cloud_instance_snapshots" "snapshots" {
  service_name = "%s"
  region       = "%s"
  instance_id  = ovh_cloud_instance.instance.id

  depends_on = [ovh_cloud_instance_snapshot.snapshot]
}
`, serviceName, region)

	configNoMatch := testAccCloudInstanceSnapshotConfig(serviceName, region, flavorID, imageID, instanceName, snapshotName) + fmt.Sprintf(`
data "ovh_cloud_instance_snapshots" "empty" {
  service_name = "%s"
  region       = "%s"
  instance_id  = "00000000-0000-0000-0000-000000000000"

  depends_on = [ovh_cloud_instance_snapshot.snapshot]
}
`, serviceName, region)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheckCloudInstanceV2(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.ovh_cloud_instance_snapshots.snapshots", "service_name", serviceName),
					resource.TestCheckResourceAttr("data.ovh_cloud_instance_snapshots.snapshots", "region", region),
					resource.TestCheckResourceAttrPair(
						"data.ovh_cloud_instance_snapshots.snapshots", "instance_id",
						"ovh_cloud_instance.instance", "id",
					),
					resource.TestCheckResourceAttrPair(
						"data.ovh_cloud_instance_snapshots.snapshots", "snapshots.0.id",
						"ovh_cloud_instance_snapshot.snapshot", "id",
					),
					resource.TestCheckResourceAttr("data.ovh_cloud_instance_snapshots.snapshots", "snapshots.0.name", snapshotName),
					resource.TestCheckResourceAttr("data.ovh_cloud_instance_snapshots.snapshots", "snapshots.0.location.region", region),
					resource.TestCheckResourceAttr("data.ovh_cloud_instance_snapshots.snapshots", "snapshots.0.resource_status", "READY"),
					resource.TestCheckResourceAttrSet("data.ovh_cloud_instance_snapshots.snapshots", "snapshots.0.status"),
					resource.TestCheckResourceAttrSet("data.ovh_cloud_instance_snapshots.snapshots", "snapshots.0.visibility"),
					resource.TestCheckResourceAttrPair(
						"data.ovh_cloud_instance_snapshots.snapshots", "snapshots.0.instance_id",
						"ovh_cloud_instance.instance", "id",
					),
				),
			},
			{
				Config: configNoMatch,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.ovh_cloud_instance_snapshots.empty", "snapshots.#", "0"),
				),
			},
		},
	})
}
