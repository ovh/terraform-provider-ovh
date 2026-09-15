package ovh

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDataSourceCloudInstanceSnapshot_basic(t *testing.T) {
	serviceName := os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST")
	region := os.Getenv("OVH_CLOUD_PROJECT_REGION_TEST")
	flavorID := resolveInstanceFlavorID(t, serviceName, region, testAccInstanceFlavorName)
	imageID := resolveInstanceImageID(t, serviceName, region, testAccInstanceImageName)
	instanceName := acctest.RandomWithPrefix(test_prefix)
	snapshotName := acctest.RandomWithPrefix(test_prefix)

	config := testAccCloudInstanceSnapshotConfig(serviceName, region, flavorID, imageID, instanceName, snapshotName) + fmt.Sprintf(`
data "ovh_cloud_instance_snapshot" "snapshot" {
  service_name = "%s"
  id           = ovh_cloud_instance_snapshot.snapshot.id
}
`, serviceName)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheckCloudInstanceV2(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.ovh_cloud_instance_snapshot.snapshot", "service_name", serviceName),
					resource.TestCheckResourceAttrPair(
						"data.ovh_cloud_instance_snapshot.snapshot", "id",
						"ovh_cloud_instance_snapshot.snapshot", "id",
					),
					resource.TestCheckResourceAttrPair(
						"data.ovh_cloud_instance_snapshot.snapshot", "instance_id",
						"ovh_cloud_instance.instance", "id",
					),
					resource.TestCheckResourceAttr("data.ovh_cloud_instance_snapshot.snapshot", "name", snapshotName),
					resource.TestCheckResourceAttr("data.ovh_cloud_instance_snapshot.snapshot", "location.region", region),
					resource.TestCheckResourceAttr("data.ovh_cloud_instance_snapshot.snapshot", "resource_status", "READY"),
					resource.TestCheckResourceAttrSet("data.ovh_cloud_instance_snapshot.snapshot", "status"),
					resource.TestCheckResourceAttrSet("data.ovh_cloud_instance_snapshot.snapshot", "visibility"),
					resource.TestCheckResourceAttrPair(
						"data.ovh_cloud_instance_snapshot.snapshot", "instance_id",
						"ovh_cloud_instance_snapshot.snapshot", "current_state.instance.id",
					),
				),
			},
		},
	})
}
