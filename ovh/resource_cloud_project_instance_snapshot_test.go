package ovh

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccCloudProjectInstanceSnapshot_basic(t *testing.T) {
	serviceName := os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST")
	instanceID := os.Getenv("OVH_CLOUD_PROJECT_INSTANCE_TEST")
	snapshotName := acctest.RandomWithPrefix(test_prefix)

	config := fmt.Sprintf(`
	resource "ovh_cloud_project_instance_snapshot" "snapshot" {
		service_name = "%s"
		instance_id  = "%s"
		name         = "%s"
	  }
	`, serviceName, instanceID, snapshotName)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheckCloud(t); testAccCheckCloudProjectExists(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ovh_cloud_project_instance_snapshot.snapshot", "service_name", serviceName),
					resource.TestCheckResourceAttr("ovh_cloud_project_instance_snapshot.snapshot", "instance_id", instanceID),
					resource.TestCheckResourceAttr("ovh_cloud_project_instance_snapshot.snapshot", "name", snapshotName),
					resource.TestCheckResourceAttr("ovh_cloud_project_instance_snapshot.snapshot", "visibility", "private"),
					resource.TestCheckResourceAttrSet("ovh_cloud_project_instance_snapshot.snapshot", "status"),
					resource.TestCheckResourceAttrSet("ovh_cloud_project_instance_snapshot.snapshot", "size"),
					resource.TestCheckResourceAttrSet("ovh_cloud_project_instance_snapshot.snapshot", "type"),
				),
			},
			{
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateIdPrefix:     serviceName + "/",
				ResourceName:            "ovh_cloud_project_instance_snapshot.snapshot",
				ImportStateVerifyIgnore: []string{"instance_id"},
			},
		},
	})
}
